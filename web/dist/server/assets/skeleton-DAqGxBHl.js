import { l as require_jsx_runtime } from "./react-dom-BccjEFPc.js";
import { n as cn } from "./button-BAIw3jLv.js";
//#region src/components/ui/skeleton.tsx
var import_jsx_runtime = require_jsx_runtime();
function Skeleton({ className, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("div", {
		"data-slot": "skeleton",
		className: cn("animate-pulse rounded-md bg-muted", className),
		...props
	});
}
//#endregion
export { Skeleton as t };
