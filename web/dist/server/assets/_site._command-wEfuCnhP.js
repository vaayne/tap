import { l as require_jsx_runtime } from "./react-dom-BccjEFPc.js";
import { a as lazyRouteComponent, n as fetchScriptDetail, o as createFileRoute } from "./server-fns-BWiKIn9S.js";
import { t as Skeleton } from "./skeleton-DAqGxBHl.js";
//#region src/routes/scripts/$site.$command.tsx
var import_jsx_runtime = require_jsx_runtime();
var $$splitNotFoundComponentImporter = () => import("./_site._command-YzaSJbPV.js");
var $$splitComponentImporter = () => import("./_site._command-CxpzknxF.js");
var Route = createFileRoute("/scripts/$site/$command")({
	loader: ({ params }) => fetchScriptDetail({ data: { name: `${params.site}/${params.command}` } }),
	component: lazyRouteComponent($$splitComponentImporter, "component"),
	pendingComponent: DetailPending,
	notFoundComponent: lazyRouteComponent($$splitNotFoundComponentImporter, "notFoundComponent")
});
function DetailPending() {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("main", {
		className: "mx-auto max-w-4xl px-4 py-8",
		children: [
			/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Skeleton, { className: "mb-4 h-7 w-16" }),
			/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Skeleton, { className: "mb-2 h-8 w-64" }),
			/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Skeleton, { className: "mb-6 h-5 w-96" }),
			/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "mb-6 grid grid-cols-3 gap-3",
				children: [
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Skeleton, { className: "h-20 rounded-xl" }),
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Skeleton, { className: "h-20 rounded-xl" }),
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Skeleton, { className: "h-20 rounded-xl" })
				]
			}),
			/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Skeleton, { className: "h-64 rounded-xl" })
		]
	});
}
//#endregion
export { Route as t };
