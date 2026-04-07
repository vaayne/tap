//go:build !windows

package browser

import (
	"context"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// scanElectronProcesses uses ps to list all processes and filters those that
// carry --remote-debugging-port=PORT in their command line.
func scanElectronProcesses(ctx context.Context) ([]ElectronProcess, error) {
	out, err := exec.CommandContext(ctx, "ps", "-ax", "-o", "pid=,args=").Output()
	if err != nil {
		return nil, nil // ps failure is non-fatal; return empty list
	}

	var procs []ElectronProcess
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		port := extractDebugPort(strings.Join(fields[1:], " "))
		if port == 0 {
			continue
		}
		procs = append(procs, ElectronProcess{
			PID:  pid,
			Name: filepath.Base(fields[1]),
			Port: port,
		})
	}
	return procs, nil
}

// extractDebugPort parses --remote-debugging-port=PORT from a command line
// string. Returns 0 if not found or not a valid port.
func extractDebugPort(cmdline string) int {
	const prefix = "--remote-debugging-port="
	idx := strings.Index(cmdline, prefix)
	if idx < 0 {
		return 0
	}
	rest := cmdline[idx+len(prefix):]
	if end := strings.IndexAny(rest, " \t"); end >= 0 {
		rest = rest[:end]
	}
	port, _ := strconv.Atoi(rest)
	if port <= 0 || port > 65535 {
		return 0
	}
	return port
}
