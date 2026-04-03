import { a as lazyRouteComponent, i as createLucideIcon, o as createFileRoute, t as fetchPopularScripts } from "./server-fns-BWiKIn9S.js";
var Github = createLucideIcon("github", [["path", {
	d: "M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4",
	key: "tonef"
}], ["path", {
	d: "M9 18c-4.51 2-5-2-7-2",
	key: "9comsn"
}]]);
var Terminal = createLucideIcon("terminal", [["path", {
	d: "M12 19h8",
	key: "baeox8"
}], ["path", {
	d: "m4 17 6-6-6-6",
	key: "1yngyt"
}]]);
//#endregion
//#region src/routes/index.tsx
var $$splitComponentImporter = () => import("./routes-C92EMzdB.js");
var Route = createFileRoute("/")({
	loader: () => fetchPopularScripts({ data: { limit: 6 } }),
	component: lazyRouteComponent($$splitComponentImporter, "component")
});
//#endregion
export { Terminal as n, Github as r, Route as t };
