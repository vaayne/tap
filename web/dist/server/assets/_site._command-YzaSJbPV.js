import { l as require_jsx_runtime } from "./react-dom-BccjEFPc.js";
import { H as Link, t as Button } from "./button-BAIw3jLv.js";
//#region src/routes/scripts/$site.$command.tsx?tsr-split=notFoundComponent
var import_jsx_runtime = require_jsx_runtime();
var SplitNotFoundComponent = () => /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("main", {
	className: "mx-auto max-w-4xl px-4 py-16 text-center",
	children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)("p", {
		className: "text-lg text-muted-foreground",
		children: "Script not found"
	}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Button, {
		variant: "ghost",
		className: "mt-4",
		render: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Link, { to: "/scripts" }),
		children: "← Back to scripts"
	})]
});
//#endregion
export { SplitNotFoundComponent as notFoundComponent };
