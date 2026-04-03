import { E as invariant, H as reactUse, I as isModuleNotFoundError, V as replaceEqualDeep, Y as __toESM, a as useRouter, l as require_jsx_runtime, m as joinPaths, n as useStore, q as require_react, v as trimPathLeft, y as trimPathRight } from "./react-dom-BccjEFPc.js";
import { i as rootRouteId, n as dummyMatchContext, r as matchContext, t as getServerFnById } from "./__tanstack-start-server-fn-resolver-BEymWXVP.js";
import { g as redirect, t as createServerFn, u as TSS_SERVER_FUNCTION } from "./createServerFn-CjUtNNCZ.js";
import { H as Link } from "./button-BAIw3jLv.js";
//#region node_modules/.pnpm/@tanstack+router-core@1.168.9/node_modules/@tanstack/router-core/dist/esm/route.js
var BaseRoute = class {
	get to() {
		return this._to;
	}
	get id() {
		return this._id;
	}
	get path() {
		return this._path;
	}
	get fullPath() {
		return this._fullPath;
	}
	constructor(options) {
		this.init = (opts) => {
			this.originalIndex = opts.originalIndex;
			const options = this.options;
			const isRoot = !options?.path && !options?.id;
			this.parentRoute = this.options.getParentRoute?.();
			if (isRoot) this._path = rootRouteId;
			else if (!this.parentRoute) invariant();
			let path = isRoot ? rootRouteId : options?.path;
			if (path && path !== "/") path = trimPathLeft(path);
			const customId = options?.id || path;
			let id = isRoot ? rootRouteId : joinPaths([this.parentRoute.id === "__root__" ? "" : this.parentRoute.id, customId]);
			if (path === "__root__") path = "/";
			if (id !== "__root__") id = joinPaths(["/", id]);
			const fullPath = id === "__root__" ? "/" : joinPaths([this.parentRoute.fullPath, path]);
			this._path = path;
			this._id = id;
			this._fullPath = fullPath;
			this._to = trimPathRight(fullPath);
		};
		this.addChildren = (children) => {
			return this._addFileChildren(children);
		};
		this._addFileChildren = (children) => {
			if (Array.isArray(children)) this.children = children;
			if (typeof children === "object" && children !== null) this.children = Object.values(children);
			return this;
		};
		this._addFileTypes = () => {
			return this;
		};
		this.updateLoader = (options) => {
			Object.assign(this.options, options);
			return this;
		};
		this.update = (options) => {
			Object.assign(this.options, options);
			return this;
		};
		this.lazy = (lazyFn) => {
			this.lazyFn = lazyFn;
			return this;
		};
		this.redirect = (opts) => redirect({
			from: this.fullPath,
			...opts
		});
		this.options = options || {};
		this.isRoot = !options?.getParentRoute;
		if (options?.id && options?.path) throw new Error(`Route cannot have both an 'id' and a 'path' option.`);
	}
};
var BaseRootRoute = class extends BaseRoute {
	constructor(options) {
		super(options);
	}
};
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/useMatch.js
var import_react = /* @__PURE__ */ __toESM(require_react(), 1);
var dummyStore = {
	state: void 0,
	get: () => void 0,
	subscribe: () => () => {}
};
/**
* Read and select the nearest or targeted route match.
* @link https://tanstack.com/router/latest/docs/framework/react/api/router/useMatchHook
*/
function useMatch(opts) {
	const router = useRouter();
	const nearestMatchId = import_react.useContext(opts.from ? dummyMatchContext : matchContext);
	const key = opts.from ?? nearestMatchId;
	const matchStore = key ? opts.from ? router.stores.getMatchStoreByRouteId(key) : router.stores.activeMatchStoresById.get(key) : void 0;
	{
		const match = matchStore?.state;
		if ((opts.shouldThrow ?? true) && !match) invariant();
		if (match === void 0) return;
		return opts.select ? opts.select(match) : match;
	}
	const previousResult = import_react.useRef(void 0);
	return useStore(matchStore ?? dummyStore, (match) => {
		if ((opts.shouldThrow ?? true) && !match) invariant();
		if (match === void 0) return;
		const selected = opts.select ? opts.select(match) : match;
		if (opts.structuralSharing ?? router.options.defaultStructuralSharing) {
			const shared = replaceEqualDeep(previousResult.current, selected);
			previousResult.current = shared;
			return shared;
		}
		return selected;
	});
}
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/useLoaderData.js
/**
* Read and select the current route's loader data with type‑safety.
*
* Options:
* - `from`/`strict`: Choose which route's data to read and strictness
* - `select`: Map the loader data to a derived value
* - `structuralSharing`: Enable structural sharing for stable references
*
* @returns The loader data (or selected value) for the matched route.
* @link https://tanstack.com/router/latest/docs/framework/react/api/router/useLoaderDataHook
*/
function useLoaderData(opts) {
	return useMatch({
		from: opts.from,
		strict: opts.strict,
		structuralSharing: opts.structuralSharing,
		select: (s) => {
			return opts.select ? opts.select(s.loaderData) : s.loaderData;
		}
	});
}
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/useLoaderDeps.js
/**
* Read and select the current route's loader dependencies object.
*
* Options:
* - `from`: Choose which route's loader deps to read
* - `select`: Map the deps to a derived value
* - `structuralSharing`: Enable structural sharing for stable references
*
* @returns The loader deps (or selected value) for the matched route.
* @link https://tanstack.com/router/latest/docs/framework/react/api/router/useLoaderDepsHook
*/
function useLoaderDeps(opts) {
	const { select, ...rest } = opts;
	return useMatch({
		...rest,
		select: (s) => {
			return select ? select(s.loaderDeps) : s.loaderDeps;
		}
	});
}
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/useParams.js
/**
* Access the current route's path parameters with type-safety.
*
* Options:
* - `from`/`strict`: Specify the matched route and whether to enforce strict typing
* - `select`: Project the params object to a derived value for memoized renders
* - `structuralSharing`: Enable structural sharing for stable references
* - `shouldThrow`: Throw if the route is not found in strict contexts
*
* @returns The params object (or selected value) for the matched route.
* @link https://tanstack.com/router/latest/docs/framework/react/api/router/useParamsHook
*/
function useParams(opts) {
	return useMatch({
		from: opts.from,
		shouldThrow: opts.shouldThrow,
		structuralSharing: opts.structuralSharing,
		strict: opts.strict,
		select: (match) => {
			const params = opts.strict === false ? match.params : match._strictParams;
			return opts.select ? opts.select(params) : params;
		}
	});
}
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/useSearch.js
/**
* Read and select the current route's search parameters with type-safety.
*
* Options:
* - `from`/`strict`: Control which route's search is read and how strictly it's typed
* - `select`: Map the search object to a derived value for render optimization
* - `structuralSharing`: Enable structural sharing for stable references
* - `shouldThrow`: Throw when the route is not found (strict contexts)
*
* @returns The search object (or selected value) for the matched route.
* @link https://tanstack.com/router/latest/docs/framework/react/api/router/useSearchHook
*/
function useSearch(opts) {
	return useMatch({
		from: opts.from,
		strict: opts.strict,
		shouldThrow: opts.shouldThrow,
		structuralSharing: opts.structuralSharing,
		select: (match) => {
			return opts.select ? opts.select(match.search) : match.search;
		}
	});
}
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/useNavigate.js
/**
* Imperative navigation hook.
*
* Returns a stable `navigate(options)` function to change the current location
* programmatically. Prefer the `Link` component for user-initiated navigation,
* and use this hook from effects, callbacks, or handlers where imperative
* navigation is required.
*
* Options:
* - `from`: Optional route base used to resolve relative `to` paths.
*
* @returns A function that accepts `NavigateOptions`.
* @link https://tanstack.com/router/latest/docs/framework/react/api/router/useNavigateHook
*/
function useNavigate(_defaultOpts) {
	const router = useRouter();
	return import_react.useCallback((options) => {
		return router.navigate({
			...options,
			from: options.from ?? _defaultOpts?.from
		});
	}, [_defaultOpts?.from, router]);
}
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/useRouteContext.js
function useRouteContext(opts) {
	return useMatch({
		...opts,
		select: (match) => opts.select ? opts.select(match.context) : match.context
	});
}
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/route.js
var import_jsx_runtime = require_jsx_runtime();
var Route = class extends BaseRoute {
	/**
	* @deprecated Use the `createRoute` function instead.
	*/
	constructor(options) {
		super(options);
		this.useMatch = (opts) => {
			return useMatch({
				select: opts?.select,
				from: this.id,
				structuralSharing: opts?.structuralSharing
			});
		};
		this.useRouteContext = (opts) => {
			return useRouteContext({
				...opts,
				from: this.id
			});
		};
		this.useSearch = (opts) => {
			return useSearch({
				select: opts?.select,
				structuralSharing: opts?.structuralSharing,
				from: this.id
			});
		};
		this.useParams = (opts) => {
			return useParams({
				select: opts?.select,
				structuralSharing: opts?.structuralSharing,
				from: this.id
			});
		};
		this.useLoaderDeps = (opts) => {
			return useLoaderDeps({
				...opts,
				from: this.id
			});
		};
		this.useLoaderData = (opts) => {
			return useLoaderData({
				...opts,
				from: this.id
			});
		};
		this.useNavigate = () => {
			return useNavigate({ from: this.fullPath });
		};
		this.Link = import_react.forwardRef((props, ref) => {
			return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Link, {
				ref,
				from: this.fullPath,
				...props
			});
		});
	}
};
/**
* Creates a non-root Route instance for code-based routing.
*
* Use this to define a route that will be composed into a route tree
* (typically via a parent route's `addChildren`). If you're using file-based
* routing, prefer `createFileRoute`.
*
* @param options Route options (path, component, loader, context, etc.).
* @returns A Route instance to be attached to the route tree.
* @link https://tanstack.com/router/latest/docs/framework/react/api/router/createRouteFunction
*/
function createRoute(options) {
	return new Route(options);
}
var RootRoute = class extends BaseRootRoute {
	/**
	* @deprecated `RootRoute` is now an internal implementation detail. Use `createRootRoute()` instead.
	*/
	constructor(options) {
		super(options);
		this.useMatch = (opts) => {
			return useMatch({
				select: opts?.select,
				from: this.id,
				structuralSharing: opts?.structuralSharing
			});
		};
		this.useRouteContext = (opts) => {
			return useRouteContext({
				...opts,
				from: this.id
			});
		};
		this.useSearch = (opts) => {
			return useSearch({
				select: opts?.select,
				structuralSharing: opts?.structuralSharing,
				from: this.id
			});
		};
		this.useParams = (opts) => {
			return useParams({
				select: opts?.select,
				structuralSharing: opts?.structuralSharing,
				from: this.id
			});
		};
		this.useLoaderDeps = (opts) => {
			return useLoaderDeps({
				...opts,
				from: this.id
			});
		};
		this.useLoaderData = (opts) => {
			return useLoaderData({
				...opts,
				from: this.id
			});
		};
		this.useNavigate = () => {
			return useNavigate({ from: this.fullPath });
		};
		this.Link = import_react.forwardRef((props, ref) => {
			return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Link, {
				ref,
				from: this.fullPath,
				...props
			});
		});
	}
};
/**
* Creates a root Route instance used to build your route tree.
*
* Typically paired with `createRouter({ routeTree })`. If you need to require
* a typed router context, use `createRootRouteWithContext` instead.
*
* @param options Root route options (component, error, pending, etc.).
* @returns A root route instance.
* @link https://tanstack.com/router/latest/docs/framework/react/api/router/createRootRouteFunction
*/
function createRootRoute(options) {
	return new RootRoute(options);
}
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/fileRoute.js
/**
* Creates a file-based Route factory for a given path.
*
* Used by TanStack Router's file-based routing to associate a file with a
* route. The returned function accepts standard route options. In normal usage
* the `path` string is inserted and maintained by the `tsr` generator.
*
* @param path File path literal for the route (usually auto-generated).
* @returns A function that accepts Route options and returns a Route instance.
* @link https://tanstack.com/router/latest/docs/framework/react/api/router/createFileRouteFunction
*/
function createFileRoute(path) {
	if (typeof path === "object") return new FileRoute(path, { silent: true }).createRoute(path);
	return new FileRoute(path, { silent: true }).createRoute;
}
/** 
@deprecated It's no longer recommended to use the `FileRoute` class directly.
Instead, use `createFileRoute('/path/to/file')(options)` to create a file route.
*/
var FileRoute = class {
	constructor(path, _opts) {
		this.path = path;
		this.createRoute = (options) => {
			const route = createRoute(options);
			route.isRoot = false;
			return route;
		};
		this.silent = _opts?.silent;
	}
};
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/lazyRouteComponent.js
/**
* Wrap a dynamic import to create a route component that supports
* `.preload()` and friendly reload-on-module-missing behavior.
*
* @param importer Function returning a module promise
* @param exportName Named export to use (default: `default`)
* @returns A lazy route component compatible with TanStack Router
* @link https://tanstack.com/router/latest/docs/framework/react/api/router/lazyRouteComponentFunction
*/
function lazyRouteComponent(importer, exportName) {
	let loadPromise;
	let comp;
	let error;
	let reload;
	const load = () => {
		if (!loadPromise) loadPromise = importer().then((res) => {
			loadPromise = void 0;
			comp = res[exportName ?? "default"];
		}).catch((err) => {
			error = err;
			if (isModuleNotFoundError(error)) {
				if (error instanceof Error && typeof window !== "undefined" && typeof sessionStorage !== "undefined") {
					const storageKey = `tanstack_router_reload:${error.message}`;
					if (!sessionStorage.getItem(storageKey)) {
						sessionStorage.setItem(storageKey, "1");
						reload = true;
					}
				}
			}
		});
		return loadPromise;
	};
	const lazyComp = function Lazy(props) {
		if (reload) {
			window.location.reload();
			throw new Promise(() => {});
		}
		if (error) throw error;
		if (!comp) if (reactUse) reactUse(load());
		else throw load();
		return import_react.createElement(comp, props);
	};
	lazyComp.preload = load;
	return lazyComp;
}
//#endregion
//#region node_modules/.pnpm/lucide-react@0.545.0_react@19.2.4/node_modules/lucide-react/dist/esm/shared/src/utils.js
/**
* @license lucide-react v0.545.0 - ISC
*
* This source code is licensed under the ISC license.
* See the LICENSE file in the root directory of this source tree.
*/
var toKebabCase = (string) => string.replace(/([a-z0-9])([A-Z])/g, "$1-$2").toLowerCase();
var toCamelCase = (string) => string.replace(/^([A-Z])|[\s-_]+(\w)/g, (match, p1, p2) => p2 ? p2.toUpperCase() : p1.toLowerCase());
var toPascalCase = (string) => {
	const camelCase = toCamelCase(string);
	return camelCase.charAt(0).toUpperCase() + camelCase.slice(1);
};
var mergeClasses = (...classes) => classes.filter((className, index, array) => {
	return Boolean(className) && className.trim() !== "" && array.indexOf(className) === index;
}).join(" ").trim();
var hasA11yProp = (props) => {
	for (const prop in props) if (prop.startsWith("aria-") || prop === "role" || prop === "title") return true;
};
//#endregion
//#region node_modules/.pnpm/lucide-react@0.545.0_react@19.2.4/node_modules/lucide-react/dist/esm/defaultAttributes.js
/**
* @license lucide-react v0.545.0 - ISC
*
* This source code is licensed under the ISC license.
* See the LICENSE file in the root directory of this source tree.
*/
var defaultAttributes = {
	xmlns: "http://www.w3.org/2000/svg",
	width: 24,
	height: 24,
	viewBox: "0 0 24 24",
	fill: "none",
	stroke: "currentColor",
	strokeWidth: 2,
	strokeLinecap: "round",
	strokeLinejoin: "round"
};
//#endregion
//#region node_modules/.pnpm/lucide-react@0.545.0_react@19.2.4/node_modules/lucide-react/dist/esm/Icon.js
/**
* @license lucide-react v0.545.0 - ISC
*
* This source code is licensed under the ISC license.
* See the LICENSE file in the root directory of this source tree.
*/
var Icon = (0, import_react.forwardRef)(({ color = "currentColor", size = 24, strokeWidth = 2, absoluteStrokeWidth, className = "", children, iconNode, ...rest }, ref) => (0, import_react.createElement)("svg", {
	ref,
	...defaultAttributes,
	width: size,
	height: size,
	stroke: color,
	strokeWidth: absoluteStrokeWidth ? Number(strokeWidth) * 24 / Number(size) : strokeWidth,
	className: mergeClasses("lucide", className),
	...!children && !hasA11yProp(rest) && { "aria-hidden": "true" },
	...rest
}, [...iconNode.map(([tag, attrs]) => (0, import_react.createElement)(tag, attrs)), ...Array.isArray(children) ? children : [children]]));
//#endregion
//#region node_modules/.pnpm/lucide-react@0.545.0_react@19.2.4/node_modules/lucide-react/dist/esm/createLucideIcon.js
/**
* @license lucide-react v0.545.0 - ISC
*
* This source code is licensed under the ISC license.
* See the LICENSE file in the root directory of this source tree.
*/
var createLucideIcon = (iconName, iconNode) => {
	const Component = (0, import_react.forwardRef)(({ className, ...props }, ref) => (0, import_react.createElement)(Icon, {
		ref,
		iconNode,
		className: mergeClasses(`lucide-${toKebabCase(toPascalCase(iconName))}`, `lucide-${iconName}`, className),
		...props
	}));
	Component.displayName = toPascalCase(iconName);
	return Component;
};
//#endregion
//#region node_modules/.pnpm/@tanstack+start-server-core@1.167.9/node_modules/@tanstack/start-server-core/dist/esm/createSsrRpc.js
var createSsrRpc = (functionId, importer) => {
	const url = "/_serverFn/" + functionId;
	const serverFnMeta = { id: functionId };
	const fn = async (...args) => {
		return (importer ? await importer() : await getServerFnById(functionId))(...args);
	};
	return Object.assign(fn, {
		url,
		serverFnMeta,
		[TSS_SERVER_FUNCTION]: true
	});
};
//#endregion
//#region src/lib/server-fns.ts
var fetchScriptList = createServerFn({ method: "GET" }).inputValidator((input) => input).handler(createSsrRpc("11e6d5cbda00c9c7830ab83720749b9fd45c0539ef22248dd225cc68b5ba34f3"));
var fetchPopularScripts = createServerFn({ method: "GET" }).inputValidator((input) => input).handler(createSsrRpc("1067cc4e8368411c2dbbab60f8fb43022fa9ae4e9fdd161914bb58b5257d657d"));
var fetchScriptDetail = createServerFn({ method: "GET" }).inputValidator((input) => input).handler(createSsrRpc("4ea0c7376908e4abdb13362c68558068218d9c46ed698af35f4ee2689c157e47"));
//#endregion
export { lazyRouteComponent as a, useNavigate as c, createLucideIcon as i, fetchScriptDetail as n, createFileRoute as o, fetchScriptList as r, createRootRoute as s, fetchPopularScripts as t };
