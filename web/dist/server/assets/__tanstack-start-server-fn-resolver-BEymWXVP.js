import { Y as __toESM, q as require_react } from "./react-dom-BccjEFPc.js";
//#region node_modules/.pnpm/@tanstack+router-core@1.168.9/node_modules/@tanstack/router-core/dist/esm/not-found.js
/** Determine if a value is a TanStack Router not-found error. */
function isNotFound(obj) {
	return !!obj?.isNotFound;
}
//#endregion
//#region node_modules/.pnpm/@tanstack+router-core@1.168.9/node_modules/@tanstack/router-core/dist/esm/root.js
/** Stable identifier used for the root route in a route tree. */
var rootRouteId = "__root__";
//#endregion
//#region node_modules/.pnpm/@tanstack+react-router@1.168.10_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@tanstack/react-router/dist/esm/matchContext.js
var import_react = /* @__PURE__ */ __toESM(require_react(), 1);
var matchContext = import_react.createContext(void 0);
var dummyMatchContext = import_react.createContext(void 0);
//#endregion
//#region \0#tanstack-start-server-fn-resolver
var manifest = {
	"11e6d5cbda00c9c7830ab83720749b9fd45c0539ef22248dd225cc68b5ba34f3": {
		functionName: "fetchScriptList_createServerFn_handler",
		importer: () => import("./server-fns-DSbvC0qe.js")
	},
	"1067cc4e8368411c2dbbab60f8fb43022fa9ae4e9fdd161914bb58b5257d657d": {
		functionName: "fetchPopularScripts_createServerFn_handler",
		importer: () => import("./server-fns-DSbvC0qe.js")
	},
	"4ea0c7376908e4abdb13362c68558068218d9c46ed698af35f4ee2689c157e47": {
		functionName: "fetchScriptDetail_createServerFn_handler",
		importer: () => import("./server-fns-DSbvC0qe.js")
	}
};
async function getServerFnById(id) {
	const serverFnInfo = manifest[id];
	if (!serverFnInfo) throw new Error("Server function info not found for " + id);
	const fnModule = await serverFnInfo.importer();
	if (!fnModule) {
		console.info("serverFnInfo", serverFnInfo);
		throw new Error("Server function module not resolved for " + id);
	}
	const action = fnModule[serverFnInfo.functionName];
	if (!action) {
		console.info("serverFnInfo", serverFnInfo);
		console.info("fnModule", fnModule);
		throw new Error(`Server function module export not resolved for serverFn ID: ${id}`);
	}
	return action;
}
//#endregion
export { isNotFound as a, rootRouteId as i, dummyMatchContext as n, matchContext as r, getServerFnById as t };
