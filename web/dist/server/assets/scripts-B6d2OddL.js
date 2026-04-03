import { l as require_jsx_runtime } from "./react-dom-BccjEFPc.js";
import { a as lazyRouteComponent, i as createLucideIcon, o as createFileRoute, r as fetchScriptList } from "./server-fns-BWiKIn9S.js";
import { t as Skeleton } from "./skeleton-DAqGxBHl.js";
var Search = createLucideIcon("search", [["path", {
	d: "m21 21-4.34-4.34",
	key: "14j7rj"
}], ["circle", {
	cx: "11",
	cy: "11",
	r: "8",
	key: "4ej97u"
}]]);
//#endregion
//#region src/routes/scripts/index.tsx
var import_jsx_runtime = require_jsx_runtime();
var $$splitComponentImporter = () => import("./scripts-B9BkAvOk.js");
var Route = createFileRoute("/scripts/")({
	validateSearch: (search) => ({
		q: search.q || void 0,
		sort: search.sort || void 0,
		site: search.site || void 0
	}),
	loaderDeps: ({ search }) => search,
	loader: ({ deps }) => fetchScriptList({ data: {
		query: deps.q,
		sort: deps.sort || "popular",
		site: deps.site
	} }),
	component: lazyRouteComponent($$splitComponentImporter, "component"),
	pendingComponent: BrowsePending
});
function BrowsePending() {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("main", {
		className: "mx-auto max-w-6xl px-4 py-8",
		children: [
			/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "mb-8",
				children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Skeleton, { className: "h-8 w-48" }), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Skeleton, { className: "mt-2 h-5 w-32" })]
			}),
			/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "mb-6 flex gap-3",
				children: [
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Skeleton, { className: "h-8 w-64" }),
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Skeleton, { className: "h-8 w-32" }),
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Skeleton, { className: "h-8 w-40" })
				]
			}),
			/* @__PURE__ */ (0, import_jsx_runtime.jsx)("div", {
				className: "grid gap-4 sm:grid-cols-2 lg:grid-cols-3",
				children: Array.from({ length: 6 }).map((_, i) => /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Skeleton, { className: "h-36 rounded-xl" }, i))
			})
		]
	});
}
//#endregion
export { Search as n, Route as t };
