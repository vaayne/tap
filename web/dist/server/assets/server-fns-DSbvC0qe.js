import { t as createServerFn, u as TSS_SERVER_FUNCTION } from "./createServerFn-CjUtNNCZ.js";
import { a as listSites, i as listScripts, t as getScript } from "./db-BIba6voY.js";
//#region node_modules/.pnpm/@tanstack+start-server-core@1.167.9/node_modules/@tanstack/start-server-core/dist/esm/createServerRpc.js
var createServerRpc = (serverFnMeta, splitImportFn) => {
	const url = "/_serverFn/" + serverFnMeta.id;
	return Object.assign(splitImportFn, {
		url,
		serverFnMeta,
		[TSS_SERVER_FUNCTION]: true
	});
};
//#endregion
//#region src/lib/server-fns.ts?tss-serverfn-split
var fetchScriptList_createServerFn_handler = createServerRpc({
	id: "11e6d5cbda00c9c7830ab83720749b9fd45c0539ef22248dd225cc68b5ba34f3",
	name: "fetchScriptList",
	filename: "src/lib/server-fns.ts"
}, (opts) => fetchScriptList.__executeServer(opts));
var fetchScriptList = createServerFn({ method: "GET" }).inputValidator((input) => input).handler(fetchScriptList_createServerFn_handler, async ({ data }) => {
	const sort = data.sort || "popular";
	const [scripts, sites] = await Promise.all([listScripts({
		query: data.query || void 0,
		sort,
		site: data.site || void 0
	}), listSites()]);
	return {
		scripts,
		sites
	};
});
var fetchPopularScripts_createServerFn_handler = createServerRpc({
	id: "1067cc4e8368411c2dbbab60f8fb43022fa9ae4e9fdd161914bb58b5257d657d",
	name: "fetchPopularScripts",
	filename: "src/lib/server-fns.ts"
}, (opts) => fetchPopularScripts.__executeServer(opts));
var fetchPopularScripts = createServerFn({ method: "GET" }).inputValidator((input) => input).handler(fetchPopularScripts_createServerFn_handler, async ({ data }) => {
	const limit = data.limit ?? 6;
	const [scripts, sites] = await Promise.all([listScripts({ sort: "popular" }), listSites()]);
	return {
		scripts: scripts.slice(0, limit),
		totalScripts: scripts.length,
		totalSites: sites.length
	};
});
var fetchScriptDetail_createServerFn_handler = createServerRpc({
	id: "4ea0c7376908e4abdb13362c68558068218d9c46ed698af35f4ee2689c157e47",
	name: "fetchScriptDetail",
	filename: "src/lib/server-fns.ts"
}, (opts) => fetchScriptDetail.__executeServer(opts));
var fetchScriptDetail = createServerFn({ method: "GET" }).inputValidator((input) => input).handler(fetchScriptDetail_createServerFn_handler, async ({ data }) => {
	return getScript(data.name);
});
//#endregion
export { fetchPopularScripts_createServerFn_handler, fetchScriptDetail_createServerFn_handler, fetchScriptList_createServerFn_handler };
