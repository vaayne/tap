package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	DefaultProxyListenAddr = "127.0.0.1:9401"
	pendingRequestTTL      = 60 * time.Second
	pendingCleanupInterval = 30 * time.Second
	reconnectInterval      = 2 * time.Second
	cdpRequestTimeout      = 10 * time.Second
)

var blockedClientMethods = map[string]struct{}{
	"Target.activateTarget": {},
	"Page.bringToFront":     {},
}

// ProxyConfig configures the local CDP proxy.
type ProxyConfig struct {
	ListenAddr string
	Upstream   string
}

type Proxy struct {
	cfg ProxyConfig

	httpServer *http.Server
	upgrader   websocket.Upgrader
	dialer     websocket.Dialer

	nextID atomic.Int64

	mu            sync.RWMutex
	upstreamConn  *websocket.Conn
	upstreamReady bool
	reconnecting  bool

	pending       map[int64]*pendingRequest
	sessionOwners map[string]*proxyClient
	targetOwners  map[string]*proxyClient
	clientState   map[*proxyClient]*proxyState
}

type pendingRequest struct {
	client     *proxyClient
	originalID int64
	method     string
	createdAt  time.Time
	internal   chan map[string]any
}

type proxyState struct {
	tabs     map[string]struct{}
	sessions map[string]struct{}
	proxyIDs map[int64]struct{}
}

type proxyClient struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func NewProxy(cfg ProxyConfig) *Proxy {
	if strings.TrimSpace(cfg.ListenAddr) == "" {
		cfg.ListenAddr = DefaultProxyListenAddr
	}
	return &Proxy{
		cfg: cfg,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(_ *http.Request) bool { return true },
		},
		pending:       make(map[int64]*pendingRequest),
		sessionOwners: make(map[string]*proxyClient),
		targetOwners:  make(map[string]*proxyClient),
		clientState:   make(map[*proxyClient]*proxyState),
	}
}

func (p *Proxy) Serve(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/json/version", p.handleJSONVersion)
	mux.HandleFunc("/json", p.handleJSONList)
	mux.HandleFunc("/json/list", p.handleJSONList)
	mux.HandleFunc("/json/new", p.handleJSONNew)
	mux.HandleFunc("/proxy/status", p.handleStatus)
	mux.HandleFunc("/devtools/browser/proxy", p.handleBrowserWS)

	p.httpServer = &http.Server{
		Addr:    p.cfg.ListenAddr,
		Handler: mux,
	}

	if err := p.connectUpstream(ctx); err != nil {
		return err
	}

	go p.cleanupPendingLoop(ctx)
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = p.httpServer.Shutdown(shutdownCtx)
		p.closeUpstream()
	}()

	err := p.httpServer.ListenAndServe()
	if err == nil || err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (p *Proxy) connectUpstream(ctx context.Context) error {
	resolvedURL, err := ResolveDebugURL(ctx, p.cfg.Upstream)
	if err != nil {
		return fmt.Errorf("resolve upstream debug endpoint: %w", err)
	}

	conn, _, err := p.dialer.DialContext(ctx, resolvedURL, nil)
	if err != nil {
		return fmt.Errorf("connect upstream debug endpoint: %w", err)
	}

	p.mu.Lock()
	if p.upstreamConn != nil {
		_ = p.upstreamConn.Close()
	}
	p.upstreamConn = conn
	p.upstreamReady = true
	p.reconnecting = false
	p.mu.Unlock()

	go p.readUpstreamLoop(ctx, conn)
	return nil
}

func (p *Proxy) scheduleReconnect(ctx context.Context) {
	p.mu.Lock()
	if p.reconnecting {
		p.mu.Unlock()
		return
	}
	p.reconnecting = true
	p.upstreamReady = false
	p.upstreamConn = nil
	p.resetOwnershipStateLocked()
	p.mu.Unlock()

	go func() {
		defer func() {
			p.mu.Lock()
			p.reconnecting = false
			p.mu.Unlock()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(reconnectInterval):
			}
			if err := p.connectUpstream(ctx); err == nil {
				return
			}
		}
	}()
}

func (p *Proxy) closeUpstream() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.upstreamConn != nil {
		_ = p.upstreamConn.Close()
		p.upstreamConn = nil
	}
	p.upstreamReady = false
}

func (p *Proxy) cleanupPendingLoop(ctx context.Context) {
	ticker := time.NewTicker(pendingCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()

			p.mu.Lock()
			for id, req := range p.pending {
				if now.Sub(req.createdAt) <= pendingRequestTTL {
					continue
				}
				p.removePendingRequestLocked(id, req)
			}
			p.mu.Unlock()
		}
	}
}

func (p *Proxy) readUpstreamLoop(ctx context.Context, conn *websocket.Conn) {
	for {
		var msg map[string]any
		if err := conn.ReadJSON(&msg); err != nil {
			p.scheduleReconnect(ctx)
			return
		}
		p.handleUpstreamMessage(msg)
	}
}

func (p *Proxy) handleUpstreamMessage(msg map[string]any) {
	if id, ok := messageID(msg); ok {
		p.mu.Lock()
		req := p.pending[id]
		if req != nil {
			delete(p.pending, id)
			if req.client != nil {
				if state := p.clientState[req.client]; state != nil {
					delete(state.proxyIDs, id)
				}
			}
		}
		p.mu.Unlock()

		if req == nil {
			return
		}

		p.recordOwnership(msg, req)

		msg["id"] = req.originalID
		if req.internal != nil {
			select {
			case req.internal <- msg:
			default:
			}
			return
		}
		if req.client != nil {
			_ = req.client.writeJSON(msg)
		}
		return
	}

	method, _ := msg["method"].(string)
	if method == "" {
		p.broadcast(msg)
		return
	}

	if sessionID, _ := msg["sessionId"].(string); sessionID != "" {
		p.mu.RLock()
		owner := p.sessionOwners[sessionID]
		p.mu.RUnlock()
		if owner != nil {
			_ = owner.writeJSON(msg)
			return
		}
	}

	if strings.HasPrefix(method, "Target.") {
		if targetID := extractTargetID(msg["params"]); targetID != "" {
			p.mu.RLock()
			owner := p.targetOwners[targetID]
			p.mu.RUnlock()
			if owner != nil {
				_ = owner.writeJSON(msg)
				return
			}
		}
	}

	p.broadcast(msg)
}

func (p *Proxy) broadcast(msg map[string]any) {
	p.mu.RLock()
	clients := make([]*proxyClient, 0, len(p.clientState))
	for client := range p.clientState {
		clients = append(clients, client)
	}
	p.mu.RUnlock()

	for _, client := range clients {
		_ = client.writeJSON(msg)
	}
}

func (p *Proxy) handleBrowserWS(w http.ResponseWriter, r *http.Request) {
	conn, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &proxyClient{conn: conn}

	p.mu.Lock()
	p.clientState[client] = newProxyState()
	p.mu.Unlock()

	defer p.removeClient(client)

	for {
		var msg map[string]any
		if err := conn.ReadJSON(&msg); err != nil {
			return
		}
		if err := p.handleClientMessage(client, msg); err != nil {
			_ = client.writeJSON(map[string]any{
				"id":    msg["id"],
				"error": map[string]any{"code": -1, "message": err.Error()},
			})
		}
	}
}

func (p *Proxy) handleClientMessage(client *proxyClient, msg map[string]any) error {
	method, _ := msg["method"].(string)
	if _, blocked := blockedClientMethods[method]; blocked {
		resp := map[string]any{"id": msg["id"], "result": map[string]any{}}
		if sessionID, _ := msg["sessionId"].(string); sessionID != "" {
			resp["sessionId"] = sessionID
		}
		return client.writeJSON(resp)
	}

	if method == "Target.createTarget" {
		params, _ := msg["params"].(map[string]any)
		if params == nil {
			params = make(map[string]any)
		}
		params["background"] = true
		params["focus"] = false
		msg["params"] = params
	}

	if id, ok := messageID(msg); ok {
		proxyID := p.nextID.Add(1)
		p.mu.Lock()
		state := p.clientState[client]
		if state != nil {
			state.proxyIDs[proxyID] = struct{}{}
		}
		p.pending[proxyID] = &pendingRequest{
			client:     client,
			originalID: id,
			method:     method,
			createdAt:  time.Now(),
		}
		p.mu.Unlock()
		msg["id"] = proxyID
	}

	return p.sendUpstream(msg)
}

func (p *Proxy) sendUpstream(msg map[string]any) error {
	p.mu.RLock()
	conn := p.upstreamConn
	ready := p.upstreamReady
	p.mu.RUnlock()
	if !ready || conn == nil {
		return fmt.Errorf("upstream browser not connected")
	}
	return conn.WriteJSON(msg)
}

func (p *Proxy) cdpRequest(method string, params map[string]any) (map[string]any, error) {
	proxyID := p.nextID.Add(1)
	ch := make(chan map[string]any, 1)

	p.mu.Lock()
	p.pending[proxyID] = &pendingRequest{
		originalID: proxyID,
		method:     method,
		createdAt:  time.Now(),
		internal:   ch,
	}
	p.mu.Unlock()

	msg := map[string]any{
		"id":     proxyID,
		"method": method,
	}
	if params != nil {
		msg["params"] = params
	}
	if err := p.sendUpstream(msg); err != nil {
		p.mu.Lock()
		delete(p.pending, proxyID)
		p.mu.Unlock()
		return nil, err
	}

	select {
	case resp := <-ch:
		if errField, ok := resp["error"]; ok && errField != nil {
			return nil, fmt.Errorf("cdp error: %v", errField)
		}
		return resp, nil
	case <-time.After(cdpRequestTimeout):
		p.mu.Lock()
		delete(p.pending, proxyID)
		p.mu.Unlock()
		return nil, fmt.Errorf("cdp request timeout")
	}
}

func (p *Proxy) handleJSONVersion(w http.ResponseWriter, r *http.Request) {
	wsURL := proxyWSURL(r)
	writeJSON(w, http.StatusOK, map[string]any{
		"Browser":              "Chrome (via tap proxy)",
		"webSocketDebuggerUrl": wsURL,
	})
}

func (p *Proxy) handleJSONList(w http.ResponseWriter, r *http.Request) {
	resp, err := p.cdpRequest("Target.getTargets", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	result, _ := resp["result"].(map[string]any)
	targetInfos, _ := result["targetInfos"].([]any)
	out := make([]map[string]any, 0, len(targetInfos))
	for _, raw := range targetInfos {
		info, _ := raw.(map[string]any)
		if info == nil {
			continue
		}
		if info["type"] != TargetTypePage {
			continue
		}
		out = append(out, map[string]any{
			"id":    info["targetId"],
			"title": info["title"],
			"url":   info["url"],
			"type":  info["type"],
		})
	}
	writeJSON(w, http.StatusOK, out)
}

func (p *Proxy) handleJSONNew(w http.ResponseWriter, r *http.Request) {
	resp, err := p.cdpRequest("Target.createTarget", map[string]any{
		"url":        "about:blank",
		"background": true,
		"focus":      false,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	writeJSON(w, http.StatusOK, resp["result"])
}

func (p *Proxy) handleStatus(w http.ResponseWriter, _ *http.Request) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	writeJSON(w, http.StatusOK, map[string]any{
		"upstream":        p.cfg.Upstream,
		"upstreamReady":   p.upstreamReady,
		"clients":         len(p.clientState),
		"sessions":        len(p.sessionOwners),
		"targets":         len(p.targetOwners),
		"pendingRequests": len(p.pending),
	})
}

func (p *Proxy) removeClient(client *proxyClient) {
	p.mu.Lock()
	defer p.mu.Unlock()

	state := p.clientState[client]
	for sessionID, owner := range p.sessionOwners {
		if owner == client {
			delete(p.sessionOwners, sessionID)
		}
	}
	for targetID, owner := range p.targetOwners {
		if owner == client {
			delete(p.targetOwners, targetID)
		}
	}
	if state != nil {
		for proxyID := range state.proxyIDs {
			delete(p.pending, proxyID)
		}
	}
	delete(p.clientState, client)
	_ = client.conn.Close()
}

func (p *Proxy) resetOwnershipStateLocked() {
	for sessionID := range p.sessionOwners {
		delete(p.sessionOwners, sessionID)
	}
	for targetID := range p.targetOwners {
		delete(p.targetOwners, targetID)
	}
	for _, state := range p.clientState {
		state.tabs = make(map[string]struct{})
		state.sessions = make(map[string]struct{})
	}
}

func (p *Proxy) removePendingRequestLocked(id int64, req *pendingRequest) {
	if req.client != nil {
		if state := p.clientState[req.client]; state != nil {
			delete(state.proxyIDs, id)
		}
	}
	delete(p.pending, id)
}

func (p *Proxy) recordOwnership(msg map[string]any, req *pendingRequest) {
	if req == nil || req.client == nil {
		return
	}

	result, _ := msg["result"].(map[string]any)
	if result == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	state := p.clientState[req.client]
	switch req.method {
	case "Target.createTarget":
		targetID, _ := result["targetId"].(string)
		if targetID == "" {
			return
		}
		p.targetOwners[targetID] = req.client
		if state != nil {
			state.tabs[targetID] = struct{}{}
		}
	case "Target.attachToTarget":
		sessionID, _ := result["sessionId"].(string)
		if sessionID == "" {
			return
		}
		p.sessionOwners[sessionID] = req.client
		if state != nil {
			state.sessions[sessionID] = struct{}{}
		}
	}
}

func newProxyState() *proxyState {
	return &proxyState{
		tabs:     make(map[string]struct{}),
		sessions: make(map[string]struct{}),
		proxyIDs: make(map[int64]struct{}),
	}
}

func (c *proxyClient) writeJSON(msg map[string]any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(msg)
}

func messageID(msg map[string]any) (int64, bool) {
	switch v := msg["id"].(type) {
	case float64:
		return int64(v), true
	case int64:
		return v, true
	case int:
		return int64(v), true
	case json.Number:
		i, err := v.Int64()
		return i, err == nil
	default:
		return 0, false
	}
}

func extractTargetID(raw any) string {
	params, _ := raw.(map[string]any)
	if params == nil {
		return ""
	}
	if targetID, _ := params["targetId"].(string); targetID != "" {
		return targetID
	}
	if info, _ := params["targetInfo"].(map[string]any); info != nil {
		if targetID, _ := info["targetId"].(string); targetID != "" {
			return targetID
		}
	}
	return ""
}

func proxyWSURL(r *http.Request) string {
	scheme := "ws"
	if r.TLS != nil {
		scheme = "wss"
	}
	return fmt.Sprintf("%s://%s/devtools/browser/proxy", scheme, r.Host)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
