import { l as require_jsx_runtime } from "./react-dom-BccjEFPc.js";
import { H as Link } from "./button-BAIw3jLv.js";
import { a as CardFooter, c as TrendingUp, i as CardDescription, n as Card, o as CardHeader, s as CardTitle, t as Badge } from "./badge-B7SmroTr.js";
//#region src/components/ScriptCard.tsx
var import_jsx_runtime = require_jsx_runtime();
function ScriptCard({ script }) {
	const [site, command] = script.name.split("/");
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Link, {
		to: "/scripts/$site/$command",
		params: {
			site,
			command
		},
		className: "block focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded-xl",
		children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Card, {
			className: "h-full transition-shadow hover:shadow-md",
			children: [/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(CardHeader, { children: [/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "flex items-start justify-between gap-2",
				children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(CardTitle, {
					className: "font-mono text-sm",
					children: script.name
				}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Badge, {
					variant: "secondary",
					children: script.site
				})]
			}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(CardDescription, {
				className: "line-clamp-2",
				children: script.description || "No description"
			})] }), /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(CardFooter, {
				className: "text-xs text-muted-foreground",
				children: [
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(TrendingUp, { className: "size-3.5 mr-1" }),
					script.usageCount,
					" ",
					script.usageCount === 1 ? "use" : "uses",
					script.domain && /* @__PURE__ */ (0, import_jsx_runtime.jsx)("span", {
						className: "ml-auto truncate text-muted-foreground/60",
						children: script.domain
					})
				]
			})]
		})
	});
}
//#endregion
export { ScriptCard as t };
