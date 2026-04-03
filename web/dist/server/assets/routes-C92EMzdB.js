import { l as require_jsx_runtime } from "./react-dom-BccjEFPc.js";
import { H as Link, t as Button } from "./button-BAIw3jLv.js";
import { i as createLucideIcon } from "./server-fns-BWiKIn9S.js";
import { n as Terminal, r as Github, t as Route } from "./routes-B2iTOTyi.js";
import { n as Card, o as CardHeader, r as CardContent, s as CardTitle, t as Badge } from "./badge-B7SmroTr.js";
import { t as ScriptCard } from "./ScriptCard-B0Ds6M5t.js";
var ArrowRight = createLucideIcon("arrow-right", [["path", {
	d: "M5 12h14",
	key: "1ays0h"
}], ["path", {
	d: "m12 5 7 7-7 7",
	key: "xquz4c"
}]]);
var Code = createLucideIcon("code", [["path", {
	d: "m16 18 6-6-6-6",
	key: "eg8j8"
}], ["path", {
	d: "m8 6-6 6 6 6",
	key: "ppft3o"
}]]);
var Download = createLucideIcon("download", [
	["path", {
		d: "M12 15V3",
		key: "m9g1x1"
	}],
	["path", {
		d: "M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4",
		key: "ih7n3h"
	}],
	["path", {
		d: "m7 10 5 5 5-5",
		key: "brsn70"
	}]
]);
var ExternalLink = createLucideIcon("external-link", [
	["path", {
		d: "M15 3h6v6",
		key: "1q9fwt"
	}],
	["path", {
		d: "M10 14 21 3",
		key: "gplh6r"
	}],
	["path", {
		d: "M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6",
		key: "a6xqqp"
	}]
]);
var FileText = createLucideIcon("file-text", [
	["path", {
		d: "M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z",
		key: "1rqfz7"
	}],
	["path", {
		d: "M14 2v4a2 2 0 0 0 2 2h4",
		key: "tnqrlb"
	}],
	["path", {
		d: "M10 9H8",
		key: "b1mrlr"
	}],
	["path", {
		d: "M16 13H8",
		key: "t4e002"
	}],
	["path", {
		d: "M16 17H8",
		key: "z1uh3a"
	}]
]);
var Layers = createLucideIcon("layers", [
	["path", {
		d: "M12.83 2.18a2 2 0 0 0-1.66 0L2.6 6.08a1 1 0 0 0 0 1.83l8.58 3.91a2 2 0 0 0 1.66 0l8.58-3.9a1 1 0 0 0 0-1.83z",
		key: "zw3jo"
	}],
	["path", {
		d: "M2 12a1 1 0 0 0 .58.91l8.6 3.91a2 2 0 0 0 1.65 0l8.58-3.9A1 1 0 0 0 22 12",
		key: "1wduqc"
	}],
	["path", {
		d: "M2 17a1 1 0 0 0 .58.91l8.6 3.91a2 2 0 0 0 1.65 0l8.58-3.9A1 1 0 0 0 22 17",
		key: "kqbvx6"
	}]
]);
var Zap = createLucideIcon("zap", [["path", {
	d: "M4 14a1 1 0 0 1-.78-1.63l9.9-10.2a.5.5 0 0 1 .86.46l-1.92 6.02A1 1 0 0 0 13 10h7a1 1 0 0 1 .78 1.63l-9.9 10.2a.5.5 0 0 1-.86-.46l1.92-6.02A1 1 0 0 0 11 14z",
	key: "1xq2db"
}]]);
//#endregion
//#region src/routes/index.tsx?tsr-split=component
var import_jsx_runtime = require_jsx_runtime();
function HomePage() {
	const { scripts, totalScripts, totalSites } = Route.useLoaderData();
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("main", { children: [
		/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("section", {
			className: "relative overflow-hidden border-b",
			children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)("div", { className: "absolute inset-0 bg-gradient-to-b from-muted/50 to-background" }), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("div", {
				className: "relative mx-auto max-w-6xl px-4 py-20 sm:py-32",
				children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
					className: "mx-auto max-w-3xl text-center",
					children: [
						/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
							className: "mb-6 inline-flex items-center gap-2 rounded-full border bg-background px-4 py-1.5 text-sm",
							children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Zap, { className: "size-3.5 text-yellow-500" }), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("span", {
								className: "text-muted-foreground",
								children: "QuickJS first, browser fallback when needed"
							})]
						}),
						/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("h1", {
							className: "mb-4 text-4xl font-bold tracking-tight sm:text-6xl",
							children: [
								"🚰 Tap into any website",
								/* @__PURE__ */ (0, import_jsx_runtime.jsx)("br", {}),
								/* @__PURE__ */ (0, import_jsx_runtime.jsx)("span", {
									className: "text-muted-foreground",
									children: "from your terminal"
								})
							]
						}),
						/* @__PURE__ */ (0, import_jsx_runtime.jsx)("p", {
							className: "mx-auto mb-8 max-w-2xl text-lg text-muted-foreground",
							children: "A Go library and CLI that runs JavaScript scripts against real websites — fast via QuickJS, with full browser fallback. Also extracts clean content from any URL."
						}),
						/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
							className: "flex flex-wrap items-center justify-center gap-3",
							children: [/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Button, {
								size: "lg",
								render: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Link, { to: "/scripts" }),
								children: [
									/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Terminal, { className: "size-4" }),
									"Browse ",
									totalScripts,
									" Scripts"
								]
							}), /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Button, {
								variant: "outline",
								size: "lg",
								render: /* @__PURE__ */ (0, import_jsx_runtime.jsx)("a", {
									href: "https://github.com/vaayne/tap",
									target: "_blank",
									rel: "noopener noreferrer"
								}),
								children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Github, { className: "size-4" }), "GitHub"]
							})]
						})
					]
				})
			})]
		}),
		/* @__PURE__ */ (0, import_jsx_runtime.jsx)("section", {
			className: "border-b",
			children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "mx-auto max-w-6xl px-4 py-16",
				children: [/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
					className: "mb-10 text-center",
					children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)("h2", {
						className: "mb-2 text-2xl font-bold sm:text-3xl",
						children: "Get started in seconds"
					}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("p", {
						className: "text-muted-foreground",
						children: "Install the CLI or use as a Go library"
					})]
				}), /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
					className: "mx-auto grid max-w-3xl gap-4 sm:grid-cols-2",
					children: [/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Card, { children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(CardHeader, { children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(CardTitle, {
						className: "flex items-center gap-2 text-base",
						children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Download, { className: "size-4" }), "CLI"]
					}) }), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(CardContent, { children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)("pre", {
						className: "overflow-x-auto rounded-lg bg-muted p-3 text-sm font-mono",
						children: "go install github.com/vaayne/tap/cmd/tap@latest"
					}) })] }), /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Card, { children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(CardHeader, { children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(CardTitle, {
						className: "flex items-center gap-2 text-base",
						children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Code, { className: "size-4" }), "Library"]
					}) }), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(CardContent, { children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)("pre", {
						className: "overflow-x-auto rounded-lg bg-muted p-3 text-sm font-mono",
						children: "go get github.com/vaayne/tap"
					}) })] })]
				})]
			})
		}),
		/* @__PURE__ */ (0, import_jsx_runtime.jsx)("section", {
			className: "border-b",
			children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "mx-auto max-w-6xl px-4 py-16",
				children: [
					/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
						className: "mb-10 text-center",
						children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)("h2", {
							className: "mb-2 text-2xl font-bold sm:text-3xl",
							children: "How it works"
						}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("p", {
							className: "text-muted-foreground",
							children: "Two commands, one shared transport layer"
						})]
					}),
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)("div", {
						className: "mx-auto max-w-4xl",
						children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)("pre", {
							className: "overflow-x-auto rounded-xl border bg-muted/50 p-6 text-xs sm:text-sm font-mono leading-relaxed text-center",
							children: `                ┌─────────────────────────────────┐
                │       Shared Transport Layer     │
                │  Level 1: HTTP  │  Level 2: CDP  │
                └────────┬────────┴───────┬────────┘
                         │                │
          ┌──────────────┴──┐    ┌────────┴──────────┐
          │    tap site     │    │    tap fetch       │
          │  QuickJS → CDP  │    │  HTTP → CDP        │
          │  → structured   │    │  → defuddle        │
          │    JSON         │    │  → markdown/HTML   │
          └─────────────────┘    └───────────────────┘`
						})
					}),
					/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
						className: "mx-auto mt-10 grid max-w-4xl gap-6 sm:grid-cols-3",
						children: [
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(FeatureCard, {
								icon: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Layers, { className: "size-5" }),
								title: "Transport Layer",
								description: "Shared HTTP client and headless Chrome (CDP), configured once and used by all consumers."
							}),
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(FeatureCard, {
								icon: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Terminal, { className: "size-5" }),
								title: "Site Scripts",
								description: "Predefined recipes that fetch structured data. QuickJS first, Chrome fallback for pages needing cookies or DOM."
							}),
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(FeatureCard, {
								icon: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(FileText, { className: "size-5" }),
								title: "Content Extraction",
								description: "Generic clean content from any URL via go-defuddle. Direct HTTP first, browser for JS-rendered pages."
							})
						]
					})
				]
			})
		}),
		/* @__PURE__ */ (0, import_jsx_runtime.jsx)("section", {
			className: "border-b",
			children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "mx-auto max-w-6xl px-4 py-16",
				children: [/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
					className: "mb-10 text-center",
					children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)("h2", {
						className: "mb-2 text-2xl font-bold sm:text-3xl",
						children: "Usage examples"
					}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("p", {
						className: "text-muted-foreground",
						children: "Pipe to jq, combine with other tools"
					})]
				}), /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
					className: "mx-auto max-w-3xl space-y-4",
					children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(UsageExample, {
						label: "Run site scripts",
						commands: [
							"tap site list",
							"tap site v2ex/hot",
							"tap site twitter/search query=\"claude code\"",
							"tap site hackernews/top | jq '.stories[:3]'"
						]
					}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(UsageExample, {
						label: "Extract clean content",
						commands: ["tap fetch https://example.com/article", "tap fetch --json https://example.com/article"]
					})]
				})]
			})
		}),
		/* @__PURE__ */ (0, import_jsx_runtime.jsx)("section", {
			className: "border-b bg-muted/30",
			children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "mx-auto grid max-w-6xl grid-cols-3 divide-x px-4",
				children: [
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(StatItem, {
						value: totalScripts,
						label: "Scripts"
					}),
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(StatItem, {
						value: totalSites,
						label: "Sites"
					}),
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(StatItem, {
						value: "Go",
						label: "Language"
					})
				]
			})
		}),
		/* @__PURE__ */ (0, import_jsx_runtime.jsx)("section", {
			className: "border-b",
			children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "mx-auto max-w-6xl px-4 py-16",
				children: [/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
					className: "mb-8 flex items-center justify-between",
					children: [/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", { children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)("h2", {
						className: "text-2xl font-bold sm:text-3xl",
						children: "Popular scripts"
					}), /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("p", {
						className: "mt-1 text-muted-foreground",
						children: [
							totalScripts,
							" scripts across ",
							totalSites,
							" platforms"
						]
					})] }), /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Button, {
						variant: "outline",
						render: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Link, { to: "/scripts" }),
						children: ["View all", /* @__PURE__ */ (0, import_jsx_runtime.jsx)(ArrowRight, { className: "size-3.5" })]
					})]
				}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("div", {
					className: "grid gap-4 sm:grid-cols-2 lg:grid-cols-3",
					children: scripts.map((script) => /* @__PURE__ */ (0, import_jsx_runtime.jsx)(ScriptCard, { script }, script.name))
				})]
			})
		}),
		/* @__PURE__ */ (0, import_jsx_runtime.jsx)("section", {
			className: "border-b bg-muted/30",
			children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "mx-auto max-w-6xl px-4 py-12 text-center",
				children: [
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)("h2", {
						className: "mb-3 text-xl font-bold",
						children: "Standing on the shoulders of giants"
					}),
					/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("p", {
						className: "mb-6 text-muted-foreground",
						children: [
							"Based on",
							" ",
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)("a", {
								href: "https://github.com/epiral/bb-sites",
								target: "_blank",
								rel: "noopener noreferrer",
								className: "underline underline-offset-4 hover:text-foreground",
								children: "bb-sites"
							}),
							" ",
							"and fully compatible with",
							" ",
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)("a", {
								href: "https://github.com/epiral/bb-browser",
								target: "_blank",
								rel: "noopener noreferrer",
								className: "underline underline-offset-4 hover:text-foreground",
								children: "bb-browser"
							})
						]
					}),
					/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
						className: "flex flex-wrap items-center justify-center gap-3",
						children: [
							/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Badge, {
								variant: "secondary",
								className: "gap-1.5 px-3 py-1",
								children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(ExternalLink, { className: "size-3" }), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("a", {
									href: "https://github.com/epiral/bb-sites",
									target: "_blank",
									rel: "noopener noreferrer",
									children: "epiral/bb-sites"
								})]
							}),
							/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Badge, {
								variant: "secondary",
								className: "gap-1.5 px-3 py-1",
								children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(ExternalLink, { className: "size-3" }), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("a", {
									href: "https://github.com/epiral/bb-browser",
									target: "_blank",
									rel: "noopener noreferrer",
									children: "epiral/bb-browser"
								})]
							}),
							/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Badge, {
								variant: "secondary",
								className: "gap-1.5 px-3 py-1",
								children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Github, { className: "size-3" }), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("a", {
									href: "https://github.com/vaayne/tap",
									target: "_blank",
									rel: "noopener noreferrer",
									children: "vaayne/tap"
								})]
							})
						]
					})
				]
			})
		}),
		/* @__PURE__ */ (0, import_jsx_runtime.jsx)("footer", {
			className: "py-8 text-center text-sm text-muted-foreground",
			children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("p", { children: [
				"MIT License ·",
				" ",
				/* @__PURE__ */ (0, import_jsx_runtime.jsx)("a", {
					href: "https://github.com/vaayne/tap",
					target: "_blank",
					rel: "noopener noreferrer",
					className: "underline underline-offset-4 hover:text-foreground",
					children: "GitHub"
				})
			] })
		})
	] });
}
function FeatureCard({ icon, title, description }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Card, { children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(CardContent, {
		className: "pt-6",
		children: [
			/* @__PURE__ */ (0, import_jsx_runtime.jsx)("div", {
				className: "mb-3 flex size-10 items-center justify-center rounded-lg bg-muted",
				children: icon
			}),
			/* @__PURE__ */ (0, import_jsx_runtime.jsx)("h3", {
				className: "mb-1 font-semibold",
				children: title
			}),
			/* @__PURE__ */ (0, import_jsx_runtime.jsx)("p", {
				className: "text-sm text-muted-foreground",
				children: description
			})
		]
	}) });
}
function UsageExample({ label, commands }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
		className: "rounded-xl border bg-muted/30 overflow-hidden",
		children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)("div", {
			className: "border-b bg-muted/50 px-4 py-2 text-xs font-medium text-muted-foreground",
			children: label
		}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("pre", {
			className: "overflow-x-auto p-4 text-sm font-mono leading-relaxed",
			children: commands.map((cmd, i) => /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", { children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)("span", {
				className: "select-none text-muted-foreground",
				children: "$ "
			}), cmd] }, i))
		})]
	});
}
function StatItem({ value, label }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
		className: "flex items-center justify-center gap-2 py-5 text-sm",
		children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)("span", {
			className: "text-lg font-bold",
			children: value
		}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("span", {
			className: "text-muted-foreground",
			children: label
		})]
	});
}
//#endregion
export { HomePage as component };
