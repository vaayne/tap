import { N as escapeHtml, Y as __toESM, a as useRouter, c as useHydrated, j as deepEqual, l as require_jsx_runtime, n as useStore, q as require_react } from "./react-dom-BccjEFPc.js";
import { d as RouterCore, h as createNonReactiveReadonlyStore, l as getAssetCrossOrigin, m as createNonReactiveMutableStore, n as Outlet, u as resolveManifestAssetLink } from "./Match-Bo4tCy4J.js";
import { H as Link, t as Button } from "./button-BAIw3jLv.js";
import { o as createFileRoute, s as createRootRoute } from "./server-fns-BWiKIn9S.js";
import { a as listSites, i as listScripts, n as getScriptContent, o as reportUsage, r as getSyncManifest, t as getScript } from "./db-BIba6voY.js";
import { n as Terminal, r as Github, t as Route$6 } from "./routes-B2iTOTyi.js";
import { n as Search, t as Route$7 } from "./scripts-B6d2OddL.js";
import { t as Route$8 } from "./_site._command-wEfuCnhP.js";
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/routerStores.js
var getStoreFactory = (opts) => {
	return {
		createMutableStore: createNonReactiveMutableStore,
		createReadonlyStore: createNonReactiveReadonlyStore,
		batch: (fn) => fn()
	};
};
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/router.js
/**
* Creates a new Router instance for React.
*
* Pass the returned router to `RouterProvider` to enable routing.
* Notable options: `routeTree` (your route definitions) and `context`
* (required if the root route was created with `createRootRouteWithContext`).
*
* @param options Router options used to configure the router.
* @returns A Router instance to be provided to `RouterProvider`.
* @link https://tanstack.com/router/latest/docs/framework/react/api/router/createRouterFunction
*/
var createRouter = (options) => {
	return new Router(options);
};
var Router = class extends RouterCore {
	constructor(options) {
		super(options, getStoreFactory);
	}
};
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/Asset.js
var import_react = /* @__PURE__ */ __toESM(require_react(), 1);
var import_jsx_runtime = require_jsx_runtime();
function Asset({ tag, attrs, children, nonce }) {
	switch (tag) {
		case "title": return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("title", {
			...attrs,
			suppressHydrationWarning: true,
			children
		});
		case "meta": return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("meta", {
			...attrs,
			suppressHydrationWarning: true
		});
		case "link": return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("link", {
			...attrs,
			nonce,
			suppressHydrationWarning: true
		});
		case "style": return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("style", {
			...attrs,
			dangerouslySetInnerHTML: { __html: children },
			nonce
		});
		case "script": return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Script, {
			attrs,
			children
		});
		default: return null;
	}
}
function Script({ attrs, children }) {
	useRouter();
	useHydrated();
	const dataScript = typeof attrs?.type === "string" && attrs.type !== "" && attrs.type !== "text/javascript" && attrs.type !== "module";
	import_react.useEffect(() => {
		if (dataScript) return;
		if (attrs?.src) {
			const normSrc = (() => {
				try {
					const base = document.baseURI || window.location.href;
					return new URL(attrs.src, base).href;
				} catch {
					return attrs.src;
				}
			})();
			if (Array.from(document.querySelectorAll("script[src]")).find((el) => el.src === normSrc)) return;
			const script = document.createElement("script");
			for (const [key, value] of Object.entries(attrs)) if (key !== "suppressHydrationWarning" && value !== void 0 && value !== false) script.setAttribute(key, typeof value === "boolean" ? "" : String(value));
			document.head.appendChild(script);
			return () => {
				if (script.parentNode) script.parentNode.removeChild(script);
			};
		}
		if (typeof children === "string") {
			const typeAttr = typeof attrs?.type === "string" ? attrs.type : "text/javascript";
			const nonceAttr = typeof attrs?.nonce === "string" ? attrs.nonce : void 0;
			if (Array.from(document.querySelectorAll("script:not([src])")).find((el) => {
				if (!(el instanceof HTMLScriptElement)) return false;
				const sType = el.getAttribute("type") ?? "text/javascript";
				const sNonce = el.getAttribute("nonce") ?? void 0;
				return el.textContent === children && sType === typeAttr && sNonce === nonceAttr;
			})) return;
			const script = document.createElement("script");
			script.textContent = children;
			if (attrs) {
				for (const [key, value] of Object.entries(attrs)) if (key !== "suppressHydrationWarning" && value !== void 0 && value !== false) script.setAttribute(key, typeof value === "boolean" ? "" : String(value));
			}
			document.head.appendChild(script);
			return () => {
				if (script.parentNode) script.parentNode.removeChild(script);
			};
		}
	}, [
		attrs,
		children,
		dataScript
	]);
	if (attrs?.src) return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("script", {
		...attrs,
		suppressHydrationWarning: true
	});
	if (typeof children === "string") return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("script", {
		...attrs,
		dangerouslySetInnerHTML: { __html: children },
		suppressHydrationWarning: true
	});
	return null;
}
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/headContentUtils.js
function buildTagsFromMatches(router, nonce, matches, assetCrossOrigin) {
	const routeMeta = matches.map((match) => match.meta).filter(Boolean);
	const resultMeta = [];
	const metaByAttribute = {};
	let title;
	for (let i = routeMeta.length - 1; i >= 0; i--) {
		const metas = routeMeta[i];
		for (let j = metas.length - 1; j >= 0; j--) {
			const m = metas[j];
			if (!m) continue;
			if (m.title) {
				if (!title) title = {
					tag: "title",
					children: m.title
				};
			} else if ("script:ld+json" in m) try {
				const json = JSON.stringify(m["script:ld+json"]);
				resultMeta.push({
					tag: "script",
					attrs: { type: "application/ld+json" },
					children: escapeHtml(json)
				});
			} catch {}
			else {
				const attribute = m.name ?? m.property;
				if (attribute) if (metaByAttribute[attribute]) continue;
				else metaByAttribute[attribute] = true;
				resultMeta.push({
					tag: "meta",
					attrs: {
						...m,
						nonce
					}
				});
			}
		}
	}
	if (title) resultMeta.push(title);
	if (nonce) resultMeta.push({
		tag: "meta",
		attrs: {
			property: "csp-nonce",
			content: nonce
		}
	});
	resultMeta.reverse();
	const constructedLinks = matches.map((match) => match.links).filter(Boolean).flat(1).map((link) => ({
		tag: "link",
		attrs: {
			...link,
			nonce
		}
	}));
	const manifest = router.ssr?.manifest;
	const assetLinks = matches.map((match) => manifest?.routes[match.routeId]?.assets ?? []).filter(Boolean).flat(1).filter((asset) => asset.tag === "link").map((asset) => ({
		tag: "link",
		attrs: {
			...asset.attrs,
			crossOrigin: getAssetCrossOrigin(assetCrossOrigin, "stylesheet") ?? asset.attrs?.crossOrigin,
			suppressHydrationWarning: true,
			nonce
		}
	}));
	const preloadLinks = [];
	matches.map((match) => router.looseRoutesById[match.routeId]).forEach((route) => router.ssr?.manifest?.routes[route.id]?.preloads?.filter(Boolean).forEach((preload) => {
		const preloadLink = resolveManifestAssetLink(preload);
		preloadLinks.push({
			tag: "link",
			attrs: {
				rel: "modulepreload",
				href: preloadLink.href,
				crossOrigin: getAssetCrossOrigin(assetCrossOrigin, "modulepreload") ?? preloadLink.crossOrigin,
				nonce
			}
		});
	}));
	const styles = matches.map((match) => match.styles).flat(1).filter(Boolean).map(({ children, ...attrs }) => ({
		tag: "style",
		attrs: {
			...attrs,
			nonce
		},
		children
	}));
	const headScripts = matches.map((match) => match.headScripts).flat(1).filter(Boolean).map(({ children, ...script }) => ({
		tag: "script",
		attrs: {
			...script,
			nonce
		},
		children
	}));
	return uniqBy([
		...resultMeta,
		...preloadLinks,
		...constructedLinks,
		...assetLinks,
		...styles,
		...headScripts
	], (d) => JSON.stringify(d));
}
/**
* Build the list of head/link/meta/script tags to render for active matches.
* Used internally by `HeadContent`.
*/
var useTags = (assetCrossOrigin) => {
	const router = useRouter();
	const nonce = router.options.ssr?.nonce;
	return buildTagsFromMatches(router, nonce, router.stores.activeMatchesSnapshot.state, assetCrossOrigin);
};
function uniqBy(arr, fn) {
	const seen = /* @__PURE__ */ new Set();
	return arr.filter((item) => {
		const key = fn(item);
		if (seen.has(key)) return false;
		seen.add(key);
		return true;
	});
}
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/HeadContent.js
/**
* Render route-managed head tags (title, meta, links, styles, head scripts).
* Place inside the document head of your app shell.
* @link https://tanstack.com/router/latest/docs/framework/react/guide/document-head-management
*/
function HeadContent(props) {
	const tags = useTags(props.assetCrossOrigin);
	const nonce = useRouter().options.ssr?.nonce;
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(import_jsx_runtime.Fragment, { children: tags.map((tag) => /* @__PURE__ */ (0, import_react.createElement)(Asset, {
		...tag,
		key: `tsr-meta-${JSON.stringify(tag)}`,
		nonce
	})) });
}
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/Scripts.js
/**
* Render body script tags collected from route matches and SSR manifests.
* Should be placed near the end of the document body.
*/
var Scripts = () => {
	const router = useRouter();
	const nonce = router.options.ssr?.nonce;
	const getAssetScripts = (matches) => {
		const assetScripts = [];
		const manifest = router.ssr?.manifest;
		if (!manifest) return [];
		matches.map((match) => router.looseRoutesById[match.routeId]).forEach((route) => manifest.routes[route.id]?.assets?.filter((d) => d.tag === "script").forEach((asset) => {
			assetScripts.push({
				tag: "script",
				attrs: {
					...asset.attrs,
					nonce
				},
				children: asset.children
			});
		}));
		return assetScripts;
	};
	const getScripts = (matches) => matches.map((match) => match.scripts).flat(1).filter(Boolean).map(({ children, ...script }) => ({
		tag: "script",
		attrs: {
			...script,
			suppressHydrationWarning: true,
			nonce
		},
		children
	}));
	{
		const assetScripts = getAssetScripts(router.stores.activeMatchesSnapshot.state);
		return renderScripts(router, getScripts(router.stores.activeMatchesSnapshot.state), assetScripts);
	}
	const assetScripts = useStore(router.stores.activeMatchesSnapshot, getAssetScripts, deepEqual);
	return renderScripts(router, useStore(router.stores.activeMatchesSnapshot, getScripts, deepEqual), assetScripts);
};
function renderScripts(router, scripts, assetScripts) {
	let serverBufferedScript = void 0;
	if (router.serverSsr) serverBufferedScript = router.serverSsr.takeBufferedScripts();
	const allScripts = [...scripts, ...assetScripts];
	if (serverBufferedScript) allScripts.unshift(serverBufferedScript);
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(import_jsx_runtime.Fragment, { children: allScripts.map((asset, i) => /* @__PURE__ */ (0, import_react.createElement)(Asset, {
		...asset,
		key: `tsr-scripts-${asset.tag}-${i}`
	})) });
}
//#endregion
//#region src/components/Header.tsx
function Header() {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("header", {
		className: "sticky top-0 z-40 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60",
		children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
			className: "mx-auto flex h-14 max-w-6xl items-center gap-6 px-4",
			children: [
				/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Link, {
					to: "/",
					className: "flex items-center gap-2 font-semibold",
					children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Terminal, { className: "size-5" }), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("span", {
						className: "text-base",
						children: "tap"
					})]
				}),
				/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("nav", {
					className: "flex items-center gap-1",
					children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Button, {
						variant: "ghost",
						size: "sm",
						render: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Link, { to: "/" }),
						children: "Home"
					}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Button, {
						variant: "ghost",
						size: "sm",
						render: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Link, { to: "/scripts" }),
						children: "Scripts"
					})]
				}),
				/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
					className: "ml-auto flex items-center gap-2",
					children: [/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Button, {
						variant: "outline",
						size: "sm",
						className: "text-muted-foreground",
						render: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Link, {
							to: "/scripts",
							search: { q: "" }
						}),
						children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Search, { className: "size-3.5" }), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("span", {
							className: "hidden sm:inline",
							children: "Search scripts…"
						})]
					}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Button, {
						variant: "ghost",
						size: "icon-sm",
						render: /* @__PURE__ */ (0, import_jsx_runtime.jsx)("a", {
							href: "https://github.com/vaayne/tap",
							target: "_blank",
							rel: "noopener noreferrer"
						}),
						children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Github, { className: "size-4" })
					})]
				})
			]
		})
	});
}
//#endregion
//#region src/styles.css?url
var styles_default = "/assets/styles-Dn3CWhTX.css";
//#endregion
//#region src/routes/__root.tsx
var Route$5 = createRootRoute({
	head: () => ({
		meta: [
			{ charSet: "utf-8" },
			{
				name: "viewport",
				content: "width=device-width, initial-scale=1"
			},
			{ title: "tap — Tap into any website from your terminal" },
			{
				name: "description",
				content: "A Go CLI and library that runs JavaScript scripts against real websites. 100+ scripts across 40+ sites. QuickJS first, browser fallback."
			}
		],
		links: [{
			rel: "stylesheet",
			href: styles_default
		}]
	}),
	component: RootComponent
});
function RootComponent() {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("html", {
		lang: "en",
		children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)("head", { children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(HeadContent, {}) }), /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("body", {
			className: "min-h-screen bg-background font-sans antialiased",
			children: [
				/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Header, {}),
				/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Outlet, {}),
				/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Scripts, {})
			]
		})]
	});
}
//#endregion
//#region src/routes/api/usage.ts
var Route$4 = createFileRoute("/api/usage")({ server: { handlers: { POST: async ({ request }) => {
	try {
		const scriptName = (await request.json())?.script;
		if (!scriptName || typeof scriptName !== "string") return Response.json({ error: "Missing required field: script" }, { status: 400 });
		await reportUsage(scriptName);
		return Response.json({ ok: true }, { status: 201 });
	} catch (error) {
		console.error("Error reporting usage:", error);
		return Response.json({ error: "Failed to report usage" }, { status: 500 });
	}
} } } });
//#endregion
//#region src/routes/api/sync.ts
var Route$3 = createFileRoute("/api/sync")({ server: { handlers: { GET: async () => {
	try {
		const scripts = await getSyncManifest();
		return Response.json({ scripts });
	} catch (error) {
		console.error("Error getting sync manifest:", error);
		return Response.json({ error: "Failed to get sync manifest" }, { status: 500 });
	}
} } } });
//#endregion
//#region src/routes/api/search.ts
var Route$2 = createFileRoute("/api/search")({ server: { handlers: { GET: async ({ request }) => {
	try {
		const url = new URL(request.url);
		const query = url.searchParams.get("q") ?? void 0;
		const sort = url.searchParams.get("sort") ?? "name";
		const site = url.searchParams.get("site") ?? void 0;
		const [scripts, sites] = await Promise.all([listScripts({
			query,
			sort,
			site
		}), listSites()]);
		return Response.json({
			scripts,
			sites
		});
	} catch (error) {
		console.error("Error listing scripts:", error);
		return Response.json({ error: "Failed to list scripts" }, { status: 500 });
	}
} } } });
//#endregion
//#region src/routes/api/scripts.ts
var Route$1 = createFileRoute("/api/scripts")({});
//#endregion
//#region src/routes/api/scripts/$.ts
var Route = createFileRoute("/api/scripts/$")({ server: { handlers: { GET: async ({ request, params }) => {
	try {
		const splat = params._splat;
		if (!splat) return Response.json({ error: "Script name is required" }, { status: 400 });
		if (splat.endsWith("/content")) {
			const name = splat.replace(/\/content$/, "");
			const content = await getScriptContent(name);
			if (!content) return Response.json({ error: `Script '${name}' not found` }, { status: 404 });
			return new Response(content, { headers: {
				"Content-Type": "application/javascript",
				"Cache-Control": "public, max-age=300"
			} });
		}
		const script = await getScript(splat);
		if (!script) return Response.json({ error: `Script '${splat}' not found` }, { status: 404 });
		return Response.json(script);
	} catch (error) {
		console.error("Error getting script:", error);
		return Response.json({ error: "Failed to get script" }, { status: 500 });
	}
} } } });
//#endregion
//#region src/routeTree.gen.ts
var IndexRoute = Route$6.update({
	id: "/",
	path: "/",
	getParentRoute: () => Route$5
});
var ScriptsIndexRoute = Route$7.update({
	id: "/scripts/",
	path: "/scripts/",
	getParentRoute: () => Route$5
});
var ApiUsageRoute = Route$4.update({
	id: "/api/usage",
	path: "/api/usage",
	getParentRoute: () => Route$5
});
var ApiSyncRoute = Route$3.update({
	id: "/api/sync",
	path: "/api/sync",
	getParentRoute: () => Route$5
});
var ApiSearchRoute = Route$2.update({
	id: "/api/search",
	path: "/api/search",
	getParentRoute: () => Route$5
});
var ApiScriptsRoute = Route$1.update({
	id: "/api/scripts",
	path: "/api/scripts",
	getParentRoute: () => Route$5
});
var ScriptsSiteCommandRoute = Route$8.update({
	id: "/scripts/$site/$command",
	path: "/scripts/$site/$command",
	getParentRoute: () => Route$5
});
var ApiScriptsRouteChildren = { ApiScriptsSplatRoute: Route.update({
	id: "/$",
	path: "/$",
	getParentRoute: () => ApiScriptsRoute
}) };
var rootRouteChildren = {
	IndexRoute,
	ApiScriptsRoute: ApiScriptsRoute._addFileChildren(ApiScriptsRouteChildren),
	ApiSearchRoute,
	ApiSyncRoute,
	ApiUsageRoute,
	ScriptsIndexRoute,
	ScriptsSiteCommandRoute
};
var routeTree = Route$5._addFileChildren(rootRouteChildren)._addFileTypes();
//#endregion
//#region src/router.tsx
function getRouter() {
	return createRouter({
		routeTree,
		scrollRestoration: true,
		defaultPreload: "intent",
		defaultPreloadStaleTime: 0
	});
}
//#endregion
export { getRouter };
