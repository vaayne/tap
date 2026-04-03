import { Y as __toESM, i as require_shim, l as require_jsx_runtime, q as require_react, r as require_with_selector, t as require_react_dom } from "./react-dom-BccjEFPc.js";
import { A as getParentNode, B as isTopLayer, C as getComputedStyle$1, D as getNodeName, E as getFrameElement, F as isLastTraversableNode, I as isNode, L as isOverflowElement, M as isContainingBlock, N as isElement, O as getNodeScroll, P as isHTMLElement, R as isShadowRoot, S as useRefWithInit, T as getDocumentElement, V as isWebKit, _ as formatErrorMessage, c as ownerVisuallyHidden, d as NOOP, f as isReactVersionAtLeast, i as useRenderElement, j as getWindow, k as getOverflowAncestors, l as EMPTY_ARRAY, m as useButton, n as cn, o as DISABLED_TRANSITIONS_STYLE, p as useMergedRefs, s as DROPDOWN_COLLISION_AVOIDANCE, t as Button, u as EMPTY_OBJECT, v as mergeProps$1, w as getContainingBlock, x as useStableCallback, y as useIsoLayoutEffect, z as isTableElement } from "./button-BAIw3jLv.js";
import { c as useNavigate, i as createLucideIcon } from "./server-fns-BWiKIn9S.js";
import { $ as round, A as getMaxListIndex, At as listNavigation, B as getAlignmentAxis, Bt as ownerDocument, C as inertValue, Ct as isSafari, D as getGridCellIndexOfCorner, Dt as escapeKey, E as findNonDisabledListIndex, Et as cancelOpen, F as clamp$1, Ft as windowResize, G as getOppositeAxisPlacements, H as getAxisLength, Ht as CompositeList, I as createCoords, It as IndexGuessBehavior, J as getSide, K as getOppositePlacement, L as evaluate, Lt as useCompositeListItem, M as isElementVisible, Mt as outsidePress, N as isIndexOutOfListBounds, O as getGridCellIndices, Ot as focusOut, P as isListIndexDisabled, Pt as triggerPress, Q as rectToClientRect, R as floor, Rt as useBaseUiId, S as transitionStatusMapping, St as isIOS, T as createGridCellMap, Tt as createChangeEventDetails, U as getExpandedPlacements, Ut as useControlled, V as getAlignmentSides, Vt as useDirection, W as getOppositeAxis, Wt as Check, X as max, Y as getSideAxis, Z as min, _ as useOpenChangeComplete, _t as ARROW_LEFT, at as isVirtualPointerEvent, b as useAnimationFrame, bt as FOCUSABLE_ATTRIBUTE, ct as contains, dt as isEventTargetWithin, et as sides, ft as isRootElement, g as useTransitionStatus, gt as ARROW_DOWN, it as isVirtualClick, j as getMinListIndex, jt as none, k as getGridNavigatedIndex, kt as itemPress, lt as getFloatingFocusElement, mt as isTypeableElement, nt as isMouseLikePointerType, ot as stopEvent, pt as isTypeableCombobox, q as getPaddingObject, rt as isReactEvent, s as COMPOSITE_KEYS, st as activeElement, tt as isClickLikeEvent, ut as getTarget, v as resolveRef, vt as ARROW_RIGHT, w as useOnMount, wt as isWebKit$1, x as TransitionStatusDataAttributes, y as AnimationFrame, yt as ARROW_UP, z as getAlignment, zt as useId } from "./composite-CSdsaCvZ.js";
import { n as Search, t as Route } from "./scripts-B6d2OddL.js";
import { t as ScriptCard } from "./ScriptCard-B0Ds6M5t.js";
var ChevronDown = createLucideIcon("chevron-down", [["path", {
	d: "m6 9 6 6 6-6",
	key: "qrunsl"
}]]);
var ChevronUp = createLucideIcon("chevron-up", [["path", {
	d: "m18 15-6-6-6 6",
	key: "153udz"
}]]);
var X = createLucideIcon("x", [["path", {
	d: "M18 6 6 18",
	key: "1bl5f8"
}], ["path", {
	d: "m6 6 12 12",
	key: "d8bk6v"
}]]);
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/utils/nodes.js
function getNodeChildren(nodes, id, onlyOpenChildren = true) {
	return nodes.filter((node) => node.parentId === id && (!onlyOpenChildren || node.context?.open)).flatMap((child) => [child, ...getNodeChildren(nodes, child.id, onlyOpenChildren)]);
}
function getNodeAncestors(nodes, id) {
	let allAncestors = [];
	let currentParentId = nodes.find((node) => node.id === id)?.parentId;
	while (currentParentId) {
		const currentNode = nodes.find((node) => node.id === currentParentId);
		currentParentId = currentNode?.parentId;
		if (currentNode) allAncestors = allAncestors.concat(currentNode);
	}
	return allAncestors;
}
var candidateSelector = /* @__PURE__ */ [
	"input:not([inert]):not([inert] *)",
	"select:not([inert]):not([inert] *)",
	"textarea:not([inert]):not([inert] *)",
	"a[href]:not([inert]):not([inert] *)",
	"button:not([inert]):not([inert] *)",
	"[tabindex]:not(slot):not([inert]):not([inert] *)",
	"audio[controls]:not([inert]):not([inert] *)",
	"video[controls]:not([inert]):not([inert] *)",
	"[contenteditable]:not([contenteditable=\"false\"]):not([inert]):not([inert] *)",
	"details>summary:first-of-type:not([inert]):not([inert] *)",
	"details:not([inert]):not([inert] *)"
].join(",");
var NoElement = typeof Element === "undefined";
var matches = NoElement ? function() {} : Element.prototype.matches || Element.prototype.msMatchesSelector || Element.prototype.webkitMatchesSelector;
var getRootNode = !NoElement && Element.prototype.getRootNode ? function(element) {
	var _element$getRootNode;
	return element === null || element === void 0 ? void 0 : (_element$getRootNode = element.getRootNode) === null || _element$getRootNode === void 0 ? void 0 : _element$getRootNode.call(element);
} : function(element) {
	return element === null || element === void 0 ? void 0 : element.ownerDocument;
};
/**
* Determines if a node is inert or in an inert ancestor.
* @param {Node} [node]
* @param {boolean} [lookUp] If true and `node` is not inert, looks up at ancestors to
*  see if any of them are inert. If false, only `node` itself is considered.
* @returns {boolean} True if inert itself or by way of being in an inert ancestor.
*  False if `node` is falsy.
*/
var _isInert = function isInert(node, lookUp) {
	var _node$getAttribute;
	if (lookUp === void 0) lookUp = true;
	var inertAtt = node === null || node === void 0 ? void 0 : (_node$getAttribute = node.getAttribute) === null || _node$getAttribute === void 0 ? void 0 : _node$getAttribute.call(node, "inert");
	return inertAtt === "" || inertAtt === "true" || lookUp && node && (typeof node.closest === "function" ? node.closest("[inert]") : _isInert(node.parentNode));
};
/**
* Determines if a node's content is editable.
* @param {Element} [node]
* @returns True if it's content-editable; false if it's not or `node` is falsy.
*/
var isContentEditable = function isContentEditable(node) {
	var _node$getAttribute2;
	var attValue = node === null || node === void 0 ? void 0 : (_node$getAttribute2 = node.getAttribute) === null || _node$getAttribute2 === void 0 ? void 0 : _node$getAttribute2.call(node, "contenteditable");
	return attValue === "" || attValue === "true";
};
/**
* @param {Element} el container to check in
* @param {boolean} includeContainer add container to check
* @param {(node: Element) => boolean} filter filter candidates
* @returns {Element[]}
*/
var getCandidates = function getCandidates(el, includeContainer, filter) {
	if (_isInert(el)) return [];
	var candidates = Array.prototype.slice.apply(el.querySelectorAll(candidateSelector));
	if (includeContainer && matches.call(el, candidateSelector)) candidates.unshift(el);
	candidates = candidates.filter(filter);
	return candidates;
};
/**
* @callback GetShadowRoot
* @param {Element} element to check for shadow root
* @returns {ShadowRoot|boolean} ShadowRoot if available or boolean indicating if a shadowRoot is attached but not available.
*/
/**
* @callback ShadowRootFilter
* @param {Element} shadowHostNode the element which contains shadow content
* @returns {boolean} true if a shadow root could potentially contain valid candidates.
*/
/**
* @typedef {Object} CandidateScope
* @property {Element} scopeParent contains inner candidates
* @property {Element[]} candidates list of candidates found in the scope parent
*/
/**
* @typedef {Object} IterativeOptions
* @property {GetShadowRoot|boolean} getShadowRoot true if shadow support is enabled; falsy if not;
*  if a function, implies shadow support is enabled and either returns the shadow root of an element
*  or a boolean stating if it has an undisclosed shadow root
* @property {(node: Element) => boolean} filter filter candidates
* @property {boolean} flatten if true then result will flatten any CandidateScope into the returned list
* @property {ShadowRootFilter} shadowRootFilter filter shadow roots;
*/
/**
* @param {Element[]} elements list of element containers to match candidates from
* @param {boolean} includeContainer add container list to check
* @param {IterativeOptions} options
* @returns {Array.<Element|CandidateScope>}
*/
var _getCandidatesIteratively = function getCandidatesIteratively(elements, includeContainer, options) {
	var candidates = [];
	var elementsToCheck = Array.from(elements);
	while (elementsToCheck.length) {
		var element = elementsToCheck.shift();
		if (_isInert(element, false)) continue;
		if (element.tagName === "SLOT") {
			var assigned = element.assignedElements();
			var nestedCandidates = _getCandidatesIteratively(assigned.length ? assigned : element.children, true, options);
			if (options.flatten) candidates.push.apply(candidates, nestedCandidates);
			else candidates.push({
				scopeParent: element,
				candidates: nestedCandidates
			});
		} else {
			if (matches.call(element, candidateSelector) && options.filter(element) && (includeContainer || !elements.includes(element))) candidates.push(element);
			var shadowRoot = element.shadowRoot || typeof options.getShadowRoot === "function" && options.getShadowRoot(element);
			var validShadowRoot = !_isInert(shadowRoot, false) && (!options.shadowRootFilter || options.shadowRootFilter(element));
			if (shadowRoot && validShadowRoot) {
				var _nestedCandidates = _getCandidatesIteratively(shadowRoot === true ? element.children : shadowRoot.children, true, options);
				if (options.flatten) candidates.push.apply(candidates, _nestedCandidates);
				else candidates.push({
					scopeParent: element,
					candidates: _nestedCandidates
				});
			} else elementsToCheck.unshift.apply(elementsToCheck, element.children);
		}
	}
	return candidates;
};
/**
* @private
* Determines if the node has an explicitly specified `tabindex` attribute.
* @param {HTMLElement} node
* @returns {boolean} True if so; false if not.
*/
var hasTabIndex = function hasTabIndex(node) {
	return !isNaN(parseInt(node.getAttribute("tabindex"), 10));
};
/**
* Determine the tab index of a given node.
* @param {HTMLElement} node
* @returns {number} Tab order (negative, 0, or positive number).
* @throws {Error} If `node` is falsy.
*/
var getTabIndex = function getTabIndex(node) {
	if (!node) throw new Error("No node provided");
	if (node.tabIndex < 0) {
		if ((/^(AUDIO|VIDEO|DETAILS)$/.test(node.tagName) || isContentEditable(node)) && !hasTabIndex(node)) return 0;
	}
	return node.tabIndex;
};
/**
* Determine the tab index of a given node __for sort order purposes__.
* @param {HTMLElement} node
* @param {boolean} [isScope] True for a custom element with shadow root or slot that, by default,
*  has tabIndex -1, but needs to be sorted by document order in order for its content to be
*  inserted into the correct sort position.
* @returns {number} Tab order (negative, 0, or positive number).
*/
var getSortOrderTabIndex = function getSortOrderTabIndex(node, isScope) {
	var tabIndex = getTabIndex(node);
	if (tabIndex < 0 && isScope && !hasTabIndex(node)) return 0;
	return tabIndex;
};
var sortOrderedTabbables = function sortOrderedTabbables(a, b) {
	return a.tabIndex === b.tabIndex ? a.documentOrder - b.documentOrder : a.tabIndex - b.tabIndex;
};
var isInput = function isInput(node) {
	return node.tagName === "INPUT";
};
var isHiddenInput = function isHiddenInput(node) {
	return isInput(node) && node.type === "hidden";
};
var isDetailsWithSummary = function isDetailsWithSummary(node) {
	return node.tagName === "DETAILS" && Array.prototype.slice.apply(node.children).some(function(child) {
		return child.tagName === "SUMMARY";
	});
};
var getCheckedRadio = function getCheckedRadio(nodes, form) {
	for (var i = 0; i < nodes.length; i++) if (nodes[i].checked && nodes[i].form === form) return nodes[i];
};
var isTabbableRadio = function isTabbableRadio(node) {
	if (!node.name) return true;
	var radioScope = node.form || getRootNode(node);
	var queryRadios = function queryRadios(name) {
		return radioScope.querySelectorAll("input[type=\"radio\"][name=\"" + name + "\"]");
	};
	var radioSet;
	if (typeof window !== "undefined" && typeof window.CSS !== "undefined" && typeof window.CSS.escape === "function") radioSet = queryRadios(window.CSS.escape(node.name));
	else try {
		radioSet = queryRadios(node.name);
	} catch (err) {
		console.error("Looks like you have a radio button with a name attribute containing invalid CSS selector characters and need the CSS.escape polyfill: %s", err.message);
		return false;
	}
	var checked = getCheckedRadio(radioSet, node.form);
	return !checked || checked === node;
};
var isRadio = function isRadio(node) {
	return isInput(node) && node.type === "radio";
};
var isNonTabbableRadio = function isNonTabbableRadio(node) {
	return isRadio(node) && !isTabbableRadio(node);
};
var isNodeAttached = function isNodeAttached(node) {
	var _nodeRoot;
	var nodeRoot = node && getRootNode(node);
	var nodeRootHost = (_nodeRoot = nodeRoot) === null || _nodeRoot === void 0 ? void 0 : _nodeRoot.host;
	var attached = false;
	if (nodeRoot && nodeRoot !== node) {
		var _nodeRootHost, _nodeRootHost$ownerDo, _node$ownerDocument;
		attached = !!((_nodeRootHost = nodeRootHost) !== null && _nodeRootHost !== void 0 && (_nodeRootHost$ownerDo = _nodeRootHost.ownerDocument) !== null && _nodeRootHost$ownerDo !== void 0 && _nodeRootHost$ownerDo.contains(nodeRootHost) || node !== null && node !== void 0 && (_node$ownerDocument = node.ownerDocument) !== null && _node$ownerDocument !== void 0 && _node$ownerDocument.contains(node));
		while (!attached && nodeRootHost) {
			var _nodeRoot2, _nodeRootHost2, _nodeRootHost2$ownerD;
			nodeRoot = getRootNode(nodeRootHost);
			nodeRootHost = (_nodeRoot2 = nodeRoot) === null || _nodeRoot2 === void 0 ? void 0 : _nodeRoot2.host;
			attached = !!((_nodeRootHost2 = nodeRootHost) !== null && _nodeRootHost2 !== void 0 && (_nodeRootHost2$ownerD = _nodeRootHost2.ownerDocument) !== null && _nodeRootHost2$ownerD !== void 0 && _nodeRootHost2$ownerD.contains(nodeRootHost));
		}
	}
	return attached;
};
var isZeroArea = function isZeroArea(node) {
	var _node$getBoundingClie = node.getBoundingClientRect(), width = _node$getBoundingClie.width, height = _node$getBoundingClie.height;
	return width === 0 && height === 0;
};
var isHidden = function isHidden(node, _ref) {
	var displayCheck = _ref.displayCheck, getShadowRoot = _ref.getShadowRoot;
	if (displayCheck === "full-native") {
		if ("checkVisibility" in node) return !node.checkVisibility({
			checkOpacity: false,
			opacityProperty: false,
			contentVisibilityAuto: true,
			visibilityProperty: true,
			checkVisibilityCSS: true
		});
	}
	if (getComputedStyle(node).visibility === "hidden") return true;
	var nodeUnderDetails = matches.call(node, "details>summary:first-of-type") ? node.parentElement : node;
	if (matches.call(nodeUnderDetails, "details:not([open]) *")) return true;
	if (!displayCheck || displayCheck === "full" || displayCheck === "full-native" || displayCheck === "legacy-full") {
		if (typeof getShadowRoot === "function") {
			var originalNode = node;
			while (node) {
				var parentElement = node.parentElement;
				var rootNode = getRootNode(node);
				if (parentElement && !parentElement.shadowRoot && getShadowRoot(parentElement) === true) return isZeroArea(node);
				else if (node.assignedSlot) node = node.assignedSlot;
				else if (!parentElement && rootNode !== node.ownerDocument) node = rootNode.host;
				else node = parentElement;
			}
			node = originalNode;
		}
		if (isNodeAttached(node)) return !node.getClientRects().length;
		if (displayCheck !== "legacy-full") return true;
	} else if (displayCheck === "non-zero-area") return isZeroArea(node);
	return false;
};
var isDisabledFromFieldset = function isDisabledFromFieldset(node) {
	if (/^(INPUT|BUTTON|SELECT|TEXTAREA)$/.test(node.tagName)) {
		var parentNode = node.parentElement;
		while (parentNode) {
			if (parentNode.tagName === "FIELDSET" && parentNode.disabled) {
				for (var i = 0; i < parentNode.children.length; i++) {
					var child = parentNode.children.item(i);
					if (child.tagName === "LEGEND") return matches.call(parentNode, "fieldset[disabled] *") ? true : !child.contains(node);
				}
				return true;
			}
			parentNode = parentNode.parentElement;
		}
	}
	return false;
};
var isNodeMatchingSelectorFocusable = function isNodeMatchingSelectorFocusable(options, node) {
	if (node.disabled || isHiddenInput(node) || isHidden(node, options) || isDetailsWithSummary(node) || isDisabledFromFieldset(node)) return false;
	return true;
};
var isNodeMatchingSelectorTabbable = function isNodeMatchingSelectorTabbable(options, node) {
	if (isNonTabbableRadio(node) || getTabIndex(node) < 0 || !isNodeMatchingSelectorFocusable(options, node)) return false;
	return true;
};
var isShadowRootTabbable = function isShadowRootTabbable(shadowHostNode) {
	var tabIndex = parseInt(shadowHostNode.getAttribute("tabindex"), 10);
	if (isNaN(tabIndex) || tabIndex >= 0) return true;
	return false;
};
/**
* @param {Array.<Element|CandidateScope>} candidates
* @returns Element[]
*/
var _sortByOrder = function sortByOrder(candidates) {
	var regularTabbables = [];
	var orderedTabbables = [];
	candidates.forEach(function(item, i) {
		var isScope = !!item.scopeParent;
		var element = isScope ? item.scopeParent : item;
		var candidateTabindex = getSortOrderTabIndex(element, isScope);
		var elements = isScope ? _sortByOrder(item.candidates) : element;
		if (candidateTabindex === 0) isScope ? regularTabbables.push.apply(regularTabbables, elements) : regularTabbables.push(element);
		else orderedTabbables.push({
			documentOrder: i,
			tabIndex: candidateTabindex,
			item,
			isScope,
			content: elements
		});
	});
	return orderedTabbables.sort(sortOrderedTabbables).reduce(function(acc, sortable) {
		sortable.isScope ? acc.push.apply(acc, sortable.content) : acc.push(sortable.content);
		return acc;
	}, []).concat(regularTabbables);
};
var tabbable = function tabbable(container, options) {
	options = options || {};
	var candidates;
	if (options.getShadowRoot) candidates = _getCandidatesIteratively([container], options.includeContainer, {
		filter: isNodeMatchingSelectorTabbable.bind(null, options),
		flatten: false,
		getShadowRoot: options.getShadowRoot,
		shadowRootFilter: isShadowRootTabbable
	});
	else candidates = getCandidates(container, options.includeContainer, isNodeMatchingSelectorTabbable.bind(null, options));
	return _sortByOrder(candidates);
};
var focusable = function focusable(container, options) {
	options = options || {};
	var candidates;
	if (options.getShadowRoot) candidates = _getCandidatesIteratively([container], options.includeContainer, {
		filter: isNodeMatchingSelectorFocusable.bind(null, options),
		flatten: true,
		getShadowRoot: options.getShadowRoot
	});
	else candidates = getCandidates(container, options.includeContainer, isNodeMatchingSelectorFocusable.bind(null, options));
	return candidates;
};
var isTabbable = function isTabbable(node, options) {
	options = options || {};
	if (!node) throw new Error("No node provided");
	if (matches.call(node, candidateSelector) === false) return false;
	return isNodeMatchingSelectorTabbable(options, node);
};
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/utils/tabbable.js
var getTabbableOptions = () => ({
	getShadowRoot: true,
	displayCheck: typeof ResizeObserver === "function" && ResizeObserver.toString().includes("[native code]") ? "full" : "none"
});
function getTabbableIn(container, dir) {
	const list = tabbable(container, getTabbableOptions());
	const len = list.length;
	if (len === 0) return;
	const active = activeElement(ownerDocument(container));
	const index = list.indexOf(active);
	return list[index === -1 ? dir === 1 ? 0 : len - 1 : index + dir];
}
function getNextTabbable(referenceElement) {
	return getTabbableIn(ownerDocument(referenceElement).body, 1) || referenceElement;
}
function getPreviousTabbable(referenceElement) {
	return getTabbableIn(ownerDocument(referenceElement).body, -1) || referenceElement;
}
function isOutsideEvent(event, container) {
	const containerElement = container || event.currentTarget;
	const relatedTarget = event.relatedTarget;
	return !relatedTarget || !contains(containerElement, relatedTarget);
}
function disableFocusInside(container) {
	tabbable(container, getTabbableOptions()).forEach((element) => {
		element.dataset.tabindex = element.getAttribute("tabindex") || "";
		element.setAttribute("tabindex", "-1");
	});
}
function enableFocusInside(container) {
	container.querySelectorAll("[data-tabindex]").forEach((element) => {
		const tabindex = element.dataset.tabindex;
		delete element.dataset.tabindex;
		if (tabindex) element.setAttribute("tabindex", tabindex);
		else element.removeAttribute("tabindex");
	});
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/csp-provider/CSPContext.js
var import_react = /* @__PURE__ */ __toESM(require_react());
/**
* @internal
*/
var CSPContext = /* @__PURE__ */ import_react.createContext(void 0);
var DEFAULT_CSP_CONTEXT_VALUE = { disableStyleElements: false };
/**
* @internal
*/
function useCSPContext() {
	return import_react.useContext(CSPContext) ?? DEFAULT_CSP_CONTEXT_VALUE;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/field/control/FieldControlDataAttributes.js
var FieldControlDataAttributes = /* @__PURE__ */ function(FieldControlDataAttributes) {
	/**
	* Present when the field is disabled.
	*/
	FieldControlDataAttributes["disabled"] = "data-disabled";
	/**
	* Present when the field is in valid state.
	*/
	FieldControlDataAttributes["valid"] = "data-valid";
	/**
	* Present when the field is in invalid state.
	*/
	FieldControlDataAttributes["invalid"] = "data-invalid";
	/**
	* Present when the field has been touched.
	*/
	FieldControlDataAttributes["touched"] = "data-touched";
	/**
	* Present when the field's value has changed.
	*/
	FieldControlDataAttributes["dirty"] = "data-dirty";
	/**
	* Present when the field is filled.
	*/
	FieldControlDataAttributes["filled"] = "data-filled";
	/**
	* Present when the field control is focused.
	*/
	FieldControlDataAttributes["focused"] = "data-focused";
	return FieldControlDataAttributes;
}({});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/field/utils/constants.js
var DEFAULT_VALIDITY_STATE = {
	badInput: false,
	customError: false,
	patternMismatch: false,
	rangeOverflow: false,
	rangeUnderflow: false,
	stepMismatch: false,
	tooLong: false,
	tooShort: false,
	typeMismatch: false,
	valid: null,
	valueMissing: false
};
var DEFAULT_FIELD_STATE_ATTRIBUTES = {
	valid: null,
	touched: false,
	dirty: false,
	filled: false,
	focused: false
};
var DEFAULT_FIELD_ROOT_STATE = {
	disabled: false,
	...DEFAULT_FIELD_STATE_ATTRIBUTES
};
var fieldValidityMapping = { valid(value) {
	if (value === null) return null;
	if (value) return { [FieldControlDataAttributes.valid]: "" };
	return { [FieldControlDataAttributes.invalid]: "" };
} };
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/field/root/FieldRootContext.js
var FieldRootContext = /* @__PURE__ */ import_react.createContext({
	invalid: void 0,
	name: void 0,
	validityData: {
		state: DEFAULT_VALIDITY_STATE,
		errors: [],
		error: "",
		value: "",
		initialValue: null
	},
	setValidityData: NOOP,
	disabled: void 0,
	touched: DEFAULT_FIELD_STATE_ATTRIBUTES.touched,
	setTouched: NOOP,
	dirty: DEFAULT_FIELD_STATE_ATTRIBUTES.dirty,
	setDirty: NOOP,
	filled: DEFAULT_FIELD_STATE_ATTRIBUTES.filled,
	setFilled: NOOP,
	focused: DEFAULT_FIELD_STATE_ATTRIBUTES.focused,
	setFocused: NOOP,
	validate: () => null,
	validationMode: "onSubmit",
	validationDebounceTime: 0,
	shouldValidateOnChange: () => false,
	state: DEFAULT_FIELD_ROOT_STATE,
	markedDirtyRef: { current: false },
	validation: {
		getValidationProps: (props = EMPTY_OBJECT) => props,
		getInputValidationProps: (props = EMPTY_OBJECT) => props,
		inputRef: { current: null },
		commit: async () => {}
	}
});
function useFieldRootContext(optional = true) {
	const context = import_react.useContext(FieldRootContext);
	if (context.setValidityData === NOOP && !optional) throw new Error(formatErrorMessage(28));
	return context;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/form/FormContext.js
var FormContext = /* @__PURE__ */ import_react.createContext({
	formRef: { current: { fields: /* @__PURE__ */ new Map() } },
	errors: {},
	clearErrors: NOOP,
	validationMode: "onSubmit",
	submitAttemptedRef: { current: false }
});
function useFormContext() {
	return import_react.useContext(FormContext);
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/labelable-provider/LabelableContext.js
/**
* A context for providing [labelable elements](https://html.spec.whatwg.org/multipage/forms.html#category-label)\
* with an accessible name (label) and description.
*/
var LabelableContext = /* @__PURE__ */ import_react.createContext({
	controlId: void 0,
	registerControlId: NOOP,
	labelId: void 0,
	setLabelId: NOOP,
	messageIds: [],
	setMessageIds: NOOP,
	getDescriptionProps: (externalProps) => externalProps
});
function useLabelableContext() {
	return import_react.useContext(LabelableContext);
}
//#endregion
//#region node_modules/.pnpm/@base-ui+utils@0.2.6_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/utils/esm/useTimeout.js
var EMPTY = 0;
var Timeout = class Timeout {
	static create() {
		return new Timeout();
	}
	currentId = EMPTY;
	/**
	* Executes `fn` after `delay`, clearing any previously scheduled call.
	*/
	start(delay, fn) {
		this.clear();
		this.currentId = setTimeout(() => {
			this.currentId = EMPTY;
			fn();
		}, delay);
	}
	isStarted() {
		return this.currentId !== EMPTY;
	}
	clear = () => {
		if (this.currentId !== EMPTY) {
			clearTimeout(this.currentId);
			this.currentId = EMPTY;
		}
	};
	disposeEffect = () => {
		return this.clear;
	};
};
/**
* A `setTimeout` with automatic cleanup and guard.
*/
function useTimeout() {
	const timeout = useRefWithInit(Timeout.create).current;
	useOnMount(timeout.disposeEffect);
	return timeout;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/field/utils/getCombinedFieldValidityData.js
/**
* Combines the field's client-side, stateful validity data with the external invalid state to
* determine the field's true validity.
*/
function getCombinedFieldValidityData(validityData, invalid) {
	return {
		...validityData,
		state: {
			...validityData.state,
			valid: !invalid && validityData.state.valid
		}
	};
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/labelable-provider/useLabelableId.js
function useLabelableId(params = {}) {
	const { id, implicit = false, controlRef } = params;
	const { controlId, registerControlId } = useLabelableContext();
	const defaultId = useBaseUiId(id);
	const controlIdForEffect = implicit ? controlId : void 0;
	const controlSourceRef = useRefWithInit(() => Symbol("labelable-control"));
	const hasRegisteredRef = import_react.useRef(false);
	const hadExplicitIdRef = import_react.useRef(id != null);
	const unregisterControlId = useStableCallback(() => {
		if (!hasRegisteredRef.current || registerControlId === NOOP) return;
		hasRegisteredRef.current = false;
		registerControlId(controlSourceRef.current, void 0);
	});
	useIsoLayoutEffect(() => {
		if (registerControlId === NOOP) return;
		let nextId;
		if (implicit) {
			const elem = controlRef?.current;
			if (isElement(elem) && elem.closest("label") != null) nextId = id ?? null;
			else nextId = controlIdForEffect ?? defaultId;
		} else if (id != null) {
			hadExplicitIdRef.current = true;
			nextId = id;
		} else if (hadExplicitIdRef.current) nextId = defaultId;
		else {
			unregisterControlId();
			return;
		}
		if (nextId === void 0) {
			unregisterControlId();
			return;
		}
		hasRegisteredRef.current = true;
		registerControlId(controlSourceRef.current, nextId);
	}, [
		id,
		controlRef,
		controlIdForEffect,
		registerControlId,
		implicit,
		defaultId,
		controlSourceRef,
		unregisterControlId
	]);
	import_react.useEffect(() => {
		return unregisterControlId;
	}, [unregisterControlId]);
	return controlId ?? defaultId;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/field/useField.js
var import_react_dom = /* @__PURE__ */ __toESM(require_react_dom());
function useField(params) {
	const { enabled = true, value, id, name, controlRef, commit } = params;
	const { formRef } = useFormContext();
	const { invalid, markedDirtyRef, validityData, setValidityData } = useFieldRootContext();
	const getValue = useStableCallback(params.getValue);
	useIsoLayoutEffect(() => {
		if (!enabled) return;
		let initialValue = value;
		if (initialValue === void 0) initialValue = getValue();
		if (validityData.initialValue === null && initialValue !== null) setValidityData((prev) => ({
			...prev,
			initialValue
		}));
	}, [
		enabled,
		setValidityData,
		value,
		validityData.initialValue,
		getValue
	]);
	useIsoLayoutEffect(() => {
		if (!enabled || !id) return;
		formRef.current.fields.set(id, {
			getValue,
			name,
			controlRef,
			validityData: getCombinedFieldValidityData(validityData, invalid),
			validate(flushSync = true) {
				let nextValue = value;
				if (nextValue === void 0) nextValue = getValue();
				markedDirtyRef.current = true;
				if (!flushSync) commit(nextValue);
				else import_react_dom.flushSync(() => commit(nextValue));
			}
		});
	}, [
		commit,
		controlRef,
		enabled,
		formRef,
		getValue,
		id,
		invalid,
		markedDirtyRef,
		name,
		validityData,
		value
	]);
	useIsoLayoutEffect(() => {
		const fields = formRef.current.fields;
		return () => {
			if (id) fields.delete(id);
		};
	}, [formRef, id]);
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/field/control/FieldControl.js
/**
* The form control to label and validate.
* Renders an `<input>` element.
*
* You can omit this part and use any Base UI input component instead. For example,
* [Input](https://base-ui.com/react/components/input), [Checkbox](https://base-ui.com/react/components/checkbox),
* or [Select](https://base-ui.com/react/components/select), among others, will work with Field out of the box.
*
* Documentation: [Base UI Field](https://base-ui.com/react/components/field)
*/
var FieldControl = /* @__PURE__ */ import_react.forwardRef(function FieldControl(componentProps, forwardedRef) {
	const { render, className, id: idProp, name: nameProp, value: valueProp, disabled: disabledProp = false, onValueChange, defaultValue, autoFocus = false, ...elementProps } = componentProps;
	const { state: fieldState, name: fieldName, disabled: fieldDisabled, setTouched, setDirty, validityData, setFocused, setFilled, validationMode, validation } = useFieldRootContext();
	const disabled = fieldDisabled || disabledProp;
	const name = fieldName ?? nameProp;
	const state = {
		...fieldState,
		disabled
	};
	const { labelId } = useLabelableContext();
	const id = useLabelableId({ id: idProp });
	useIsoLayoutEffect(() => {
		const hasExternalValue = valueProp != null;
		if (validation.inputRef.current?.value || hasExternalValue && valueProp !== "") setFilled(true);
		else if (hasExternalValue && valueProp === "") setFilled(false);
	}, [
		validation.inputRef,
		setFilled,
		valueProp
	]);
	const inputRef = import_react.useRef(null);
	useIsoLayoutEffect(() => {
		if (autoFocus && inputRef.current === activeElement(ownerDocument(inputRef.current))) setFocused(true);
	}, [autoFocus, setFocused]);
	const [valueUnwrapped] = useControlled({
		controlled: valueProp,
		default: defaultValue,
		name: "FieldControl",
		state: "value"
	});
	const isControlled = valueProp !== void 0;
	const value = isControlled ? valueUnwrapped : void 0;
	useField({
		id,
		name,
		commit: validation.commit,
		value,
		getValue: () => validation.inputRef.current?.value,
		controlRef: validation.inputRef
	});
	return useRenderElement("input", componentProps, {
		ref: [forwardedRef, inputRef],
		state,
		props: [
			{
				id,
				disabled,
				name,
				ref: validation.inputRef,
				"aria-labelledby": labelId,
				autoFocus,
				...isControlled ? { value } : { defaultValue },
				onChange(event) {
					const inputValue = event.currentTarget.value;
					onValueChange?.(inputValue, createChangeEventDetails(none, event.nativeEvent));
					setDirty(inputValue !== validityData.initialValue);
					setFilled(inputValue !== "");
				},
				onFocus() {
					setFocused(true);
				},
				onBlur(event) {
					setTouched(true);
					setFocused(false);
					if (validationMode === "onBlur") validation.commit(event.currentTarget.value);
				},
				onKeyDown(event) {
					if (event.currentTarget.tagName === "INPUT" && event.key === "Enter") {
						setTouched(true);
						validation.commit(event.currentTarget.value);
					}
				}
			},
			validation.getInputValidationProps(),
			elementProps
		],
		stateAttributesMapping: fieldValidityMapping
	});
});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/input/Input.js
var import_jsx_runtime = require_jsx_runtime();
var Input$1 = /* @__PURE__ */ import_react.forwardRef(function Input(props, forwardedRef) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(FieldControl, {
		ref: forwardedRef,
		...props
	});
});
//#endregion
//#region src/components/ui/input.tsx
function Input({ className, type, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Input$1, {
		type,
		"data-slot": "input",
		className: cn("h-8 w-full min-w-0 rounded-lg border border-input bg-transparent px-2.5 py-1 text-base transition-colors outline-none file:inline-flex file:h-6 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:cursor-not-allowed disabled:bg-input/50 disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 md:text-sm dark:bg-input/30 dark:disabled:bg-input/80 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40", className),
		...props
	});
}
//#endregion
//#region node_modules/.pnpm/@base-ui+utils@0.2.6_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/utils/esm/visuallyHidden.js
var visuallyHiddenBase = {
	clipPath: "inset(50%)",
	overflow: "hidden",
	whiteSpace: "nowrap",
	border: 0,
	padding: 0,
	width: 1,
	height: 1,
	margin: -1
};
var visuallyHidden = {
	...visuallyHiddenBase,
	position: "fixed",
	top: 0,
	left: 0
};
var visuallyHiddenInput = {
	...visuallyHiddenBase,
	position: "absolute"
};
//#endregion
//#region node_modules/.pnpm/@base-ui+utils@0.2.6_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/utils/esm/useOnFirstRender.js
function useOnFirstRender(fn) {
	const ref = import_react.useRef(true);
	if (ref.current) {
		ref.current = false;
		fn();
	}
}
//#endregion
//#region node_modules/.pnpm/@base-ui+utils@0.2.6_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/utils/esm/useValueAsRef.js
/**
* Untracks the provided value by turning it into a ref to remove its reactivity.
*
* Used to access the passed value inside `React.useEffect` without causing the effect to re-run when the value changes.
*/
function useValueAsRef(value) {
	const latest = useRefWithInit(createLatestRef, value).current;
	latest.next = value;
	useIsoLayoutEffect(latest.effect);
	return latest;
}
function createLatestRef(value) {
	const latest = {
		current: value,
		next: value,
		effect: () => {
			latest.current = latest.next;
		}
	};
	return latest;
}
//#endregion
//#region node_modules/.pnpm/reselect@5.1.1/node_modules/reselect/dist/reselect.mjs
var NOT_FOUND = /* @__PURE__ */ Symbol("NOT_FOUND");
function assertIsFunction(func, errorMessage = `expected a function, instead received ${typeof func}`) {
	if (typeof func !== "function") throw new TypeError(errorMessage);
}
function assertIsObject(object, errorMessage = `expected an object, instead received ${typeof object}`) {
	if (typeof object !== "object") throw new TypeError(errorMessage);
}
function assertIsArrayOfFunctions(array, errorMessage = `expected all items to be functions, instead received the following types: `) {
	if (!array.every((item) => typeof item === "function")) {
		const itemTypes = array.map((item) => typeof item === "function" ? `function ${item.name || "unnamed"}()` : typeof item).join(", ");
		throw new TypeError(`${errorMessage}[${itemTypes}]`);
	}
}
var ensureIsArray = (item) => {
	return Array.isArray(item) ? item : [item];
};
function getDependencies(createSelectorArgs) {
	const dependencies = Array.isArray(createSelectorArgs[0]) ? createSelectorArgs[0] : createSelectorArgs;
	assertIsArrayOfFunctions(dependencies, `createSelector expects all input-selectors to be functions, but received the following types: `);
	return dependencies;
}
function collectInputSelectorResults(dependencies, inputSelectorArgs) {
	const inputSelectorResults = [];
	const { length } = dependencies;
	for (let i = 0; i < length; i++) inputSelectorResults.push(dependencies[i].apply(null, inputSelectorArgs));
	return inputSelectorResults;
}
function createSingletonCache(equals) {
	let entry;
	return {
		get(key) {
			if (entry && equals(entry.key, key)) return entry.value;
			return NOT_FOUND;
		},
		put(key, value) {
			entry = {
				key,
				value
			};
		},
		getEntries() {
			return entry ? [entry] : [];
		},
		clear() {
			entry = void 0;
		}
	};
}
function createLruCache(maxSize, equals) {
	let entries = [];
	function get(key) {
		const cacheIndex = entries.findIndex((entry) => equals(key, entry.key));
		if (cacheIndex > -1) {
			const entry = entries[cacheIndex];
			if (cacheIndex > 0) {
				entries.splice(cacheIndex, 1);
				entries.unshift(entry);
			}
			return entry.value;
		}
		return NOT_FOUND;
	}
	function put(key, value) {
		if (get(key) === NOT_FOUND) {
			entries.unshift({
				key,
				value
			});
			if (entries.length > maxSize) entries.pop();
		}
	}
	function getEntries() {
		return entries;
	}
	function clear() {
		entries = [];
	}
	return {
		get,
		put,
		getEntries,
		clear
	};
}
var referenceEqualityCheck = (a, b) => a === b;
function createCacheKeyComparator(equalityCheck) {
	return function areArgumentsShallowlyEqual(prev, next) {
		if (prev === null || next === null || prev.length !== next.length) return false;
		const { length } = prev;
		for (let i = 0; i < length; i++) if (!equalityCheck(prev[i], next[i])) return false;
		return true;
	};
}
function lruMemoize(func, equalityCheckOrOptions) {
	const { equalityCheck = referenceEqualityCheck, maxSize = 1, resultEqualityCheck } = typeof equalityCheckOrOptions === "object" ? equalityCheckOrOptions : { equalityCheck: equalityCheckOrOptions };
	const comparator = createCacheKeyComparator(equalityCheck);
	let resultsCount = 0;
	const cache = maxSize <= 1 ? createSingletonCache(comparator) : createLruCache(maxSize, comparator);
	function memoized() {
		let value = cache.get(arguments);
		if (value === NOT_FOUND) {
			value = func.apply(null, arguments);
			resultsCount++;
			if (resultEqualityCheck) {
				const matchingEntry = cache.getEntries().find((entry) => resultEqualityCheck(entry.value, value));
				if (matchingEntry) {
					value = matchingEntry.value;
					resultsCount !== 0 && resultsCount--;
				}
			}
			cache.put(arguments, value);
		}
		return value;
	}
	memoized.clearCache = () => {
		cache.clear();
		memoized.resetResultsCount();
	};
	memoized.resultsCount = () => resultsCount;
	memoized.resetResultsCount = () => {
		resultsCount = 0;
	};
	return memoized;
}
var StrongRef = class {
	constructor(value) {
		this.value = value;
	}
	deref() {
		return this.value;
	}
};
var Ref = typeof WeakRef !== "undefined" ? WeakRef : StrongRef;
var UNTERMINATED = 0;
var TERMINATED = 1;
function createCacheNode() {
	return {
		s: UNTERMINATED,
		v: void 0,
		o: null,
		p: null
	};
}
function weakMapMemoize(func, options = {}) {
	let fnNode = createCacheNode();
	const { resultEqualityCheck } = options;
	let lastResult;
	let resultsCount = 0;
	function memoized() {
		let cacheNode = fnNode;
		const { length } = arguments;
		for (let i = 0, l = length; i < l; i++) {
			const arg = arguments[i];
			if (typeof arg === "function" || typeof arg === "object" && arg !== null) {
				let objectCache = cacheNode.o;
				if (objectCache === null) cacheNode.o = objectCache = /* @__PURE__ */ new WeakMap();
				const objectNode = objectCache.get(arg);
				if (objectNode === void 0) {
					cacheNode = createCacheNode();
					objectCache.set(arg, cacheNode);
				} else cacheNode = objectNode;
			} else {
				let primitiveCache = cacheNode.p;
				if (primitiveCache === null) cacheNode.p = primitiveCache = /* @__PURE__ */ new Map();
				const primitiveNode = primitiveCache.get(arg);
				if (primitiveNode === void 0) {
					cacheNode = createCacheNode();
					primitiveCache.set(arg, cacheNode);
				} else cacheNode = primitiveNode;
			}
		}
		const terminatedNode = cacheNode;
		let result;
		if (cacheNode.s === TERMINATED) result = cacheNode.v;
		else {
			result = func.apply(null, arguments);
			resultsCount++;
			if (resultEqualityCheck) {
				const lastResultValue = lastResult?.deref?.() ?? lastResult;
				if (lastResultValue != null && resultEqualityCheck(lastResultValue, result)) {
					result = lastResultValue;
					resultsCount !== 0 && resultsCount--;
				}
				lastResult = typeof result === "object" && result !== null || typeof result === "function" ? new Ref(result) : result;
			}
		}
		terminatedNode.s = TERMINATED;
		terminatedNode.v = result;
		return result;
	}
	memoized.clearCache = () => {
		fnNode = createCacheNode();
		memoized.resetResultsCount();
	};
	memoized.resultsCount = () => resultsCount;
	memoized.resetResultsCount = () => {
		resultsCount = 0;
	};
	return memoized;
}
function createSelectorCreator(memoizeOrOptions, ...memoizeOptionsFromArgs) {
	const createSelectorCreatorOptions = typeof memoizeOrOptions === "function" ? {
		memoize: memoizeOrOptions,
		memoizeOptions: memoizeOptionsFromArgs
	} : memoizeOrOptions;
	const createSelector2 = (...createSelectorArgs) => {
		let recomputations = 0;
		let dependencyRecomputations = 0;
		let lastResult;
		let directlyPassedOptions = {};
		let resultFunc = createSelectorArgs.pop();
		if (typeof resultFunc === "object") {
			directlyPassedOptions = resultFunc;
			resultFunc = createSelectorArgs.pop();
		}
		assertIsFunction(resultFunc, `createSelector expects an output function after the inputs, but received: [${typeof resultFunc}]`);
		const { memoize, memoizeOptions = [], argsMemoize = weakMapMemoize, argsMemoizeOptions = [], devModeChecks = {} } = {
			...createSelectorCreatorOptions,
			...directlyPassedOptions
		};
		const finalMemoizeOptions = ensureIsArray(memoizeOptions);
		const finalArgsMemoizeOptions = ensureIsArray(argsMemoizeOptions);
		const dependencies = getDependencies(createSelectorArgs);
		const memoizedResultFunc = memoize(function recomputationWrapper() {
			recomputations++;
			return resultFunc.apply(null, arguments);
		}, ...finalMemoizeOptions);
		const selector = argsMemoize(function dependenciesChecker() {
			dependencyRecomputations++;
			const inputSelectorResults = collectInputSelectorResults(dependencies, arguments);
			lastResult = memoizedResultFunc.apply(null, inputSelectorResults);
			return lastResult;
		}, ...finalArgsMemoizeOptions);
		return Object.assign(selector, {
			resultFunc,
			memoizedResultFunc,
			dependencies,
			dependencyRecomputations: () => dependencyRecomputations,
			resetDependencyRecomputations: () => {
				dependencyRecomputations = 0;
			},
			lastResult: () => lastResult,
			recomputations: () => recomputations,
			resetRecomputations: () => {
				recomputations = 0;
			},
			memoize,
			argsMemoize
		});
	};
	Object.assign(createSelector2, { withTypes: () => createSelector2 });
	return createSelector2;
}
var createSelector$1 = /* @__PURE__ */ createSelectorCreator(weakMapMemoize);
var createStructuredSelector = Object.assign((inputSelectorsObject, selectorCreator = createSelector$1) => {
	assertIsObject(inputSelectorsObject, `createStructuredSelector expects first argument to be an object where each property is a selector, instead received a ${typeof inputSelectorsObject}`);
	const inputSelectorKeys = Object.keys(inputSelectorsObject);
	return selectorCreator(inputSelectorKeys.map((key) => inputSelectorsObject[key]), (...inputSelectorResults) => {
		return inputSelectorResults.reduce((composition, value, index) => {
			composition[inputSelectorKeys[index]] = value;
			return composition;
		}, {});
	});
}, { withTypes: () => createStructuredSelector });
createSelectorCreator({
	memoize: lruMemoize,
	memoizeOptions: {
		maxSize: 1,
		equalityCheck: Object.is
	}
});
/**
* Creates a selector function that can be used to derive values from the store's state.
* The selector can take up to three additional arguments that can be used in the selector logic.
* This function accepts up to six functions and combines them into a single selector function.
* The last parameter is the combiner function that combines the results of the previous selectors.
*
* @example
* const selector = createSelector(
*  (state) => state.disabled
* );
*
* @example
* const selector = createSelector(
*   (state) => state.disabled,
*   (state) => state.open,
*   (disabled, open) => ({ disabled, open })
* );
*
*/
var createSelector = (a, b, c, d, e, f, ...other) => {
	if (other.length > 0) throw new Error(formatErrorMessage(1));
	let selector;
	if (a && b && c && d && e && f) selector = (state, a1, a2, a3) => {
		return f(a(state, a1, a2, a3), b(state, a1, a2, a3), c(state, a1, a2, a3), d(state, a1, a2, a3), e(state, a1, a2, a3), a1, a2, a3);
	};
	else if (a && b && c && d && e) selector = (state, a1, a2, a3) => {
		return e(a(state, a1, a2, a3), b(state, a1, a2, a3), c(state, a1, a2, a3), d(state, a1, a2, a3), a1, a2, a3);
	};
	else if (a && b && c && d) selector = (state, a1, a2, a3) => {
		return d(a(state, a1, a2, a3), b(state, a1, a2, a3), c(state, a1, a2, a3), a1, a2, a3);
	};
	else if (a && b && c) selector = (state, a1, a2, a3) => {
		return c(a(state, a1, a2, a3), b(state, a1, a2, a3), a1, a2, a3);
	};
	else if (a && b) selector = (state, a1, a2, a3) => {
		return b(a(state, a1, a2, a3), a1, a2, a3);
	};
	else if (a) selector = a;
	else throw new Error("Missing arguments");
	return selector;
};
//#endregion
//#region node_modules/.pnpm/@base-ui+utils@0.2.6_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/utils/esm/fastHooks.js
var hooks = [];
var currentInstance = void 0;
function getInstance() {
	return currentInstance;
}
function register(hook) {
	hooks.push(hook);
}
//#endregion
//#region node_modules/.pnpm/@base-ui+utils@0.2.6_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/utils/esm/store/useStore.js
var import_shim = require_shim();
var import_with_selector = require_with_selector();
var useStoreImplementation = isReactVersionAtLeast(19) ? useStoreFast : useStoreLegacy;
function useStore(store, selector, a1, a2, a3) {
	return useStoreImplementation(store, selector, a1, a2, a3);
}
function useStoreR19(store, selector, a1, a2, a3) {
	const getSelection = import_react.useCallback(() => selector(store.getSnapshot(), a1, a2, a3), [
		store,
		selector,
		a1,
		a2,
		a3
	]);
	return (0, import_shim.useSyncExternalStore)(store.subscribe, getSelection, getSelection);
}
register({
	before(instance) {
		instance.syncIndex = 0;
		if (!instance.didInitialize) {
			instance.syncTick = 1;
			instance.syncHooks = [];
			instance.didChangeStore = true;
			instance.getSnapshot = () => {
				let didChange = false;
				for (let i = 0; i < instance.syncHooks.length; i += 1) {
					const hook = instance.syncHooks[i];
					const value = hook.selector(hook.store.state, hook.a1, hook.a2, hook.a3);
					if (hook.didChange || !Object.is(hook.value, value)) {
						didChange = true;
						hook.value = value;
						hook.didChange = false;
					}
				}
				if (didChange) instance.syncTick += 1;
				return instance.syncTick;
			};
		}
	},
	after(instance) {
		if (instance.syncHooks.length > 0) {
			if (instance.didChangeStore) {
				instance.didChangeStore = false;
				instance.subscribe = (onStoreChange) => {
					const stores = /* @__PURE__ */ new Set();
					for (const hook of instance.syncHooks) stores.add(hook.store);
					const unsubscribes = [];
					for (const store of stores) unsubscribes.push(store.subscribe(onStoreChange));
					return () => {
						for (const unsubscribe of unsubscribes) unsubscribe();
					};
				};
			}
			(0, import_shim.useSyncExternalStore)(instance.subscribe, instance.getSnapshot, instance.getSnapshot);
		}
	}
});
function useStoreFast(store, selector, a1, a2, a3) {
	const instance = getInstance();
	if (!instance) return useStoreR19(store, selector, a1, a2, a3);
	const index = instance.syncIndex;
	instance.syncIndex += 1;
	let hook;
	if (!instance.didInitialize) {
		hook = {
			store,
			selector,
			a1,
			a2,
			a3,
			value: selector(store.getSnapshot(), a1, a2, a3),
			didChange: false
		};
		instance.syncHooks.push(hook);
	} else {
		hook = instance.syncHooks[index];
		if (hook.store !== store || hook.selector !== selector || !Object.is(hook.a1, a1) || !Object.is(hook.a2, a2) || !Object.is(hook.a3, a3)) {
			if (hook.store !== store) instance.didChangeStore = true;
			hook.store = store;
			hook.selector = selector;
			hook.a1 = a1;
			hook.a2 = a2;
			hook.a3 = a3;
			hook.didChange = true;
		}
	}
	return hook.value;
}
function useStoreLegacy(store, selector, a1, a2, a3) {
	return (0, import_with_selector.useSyncExternalStoreWithSelector)(store.subscribe, store.getSnapshot, store.getSnapshot, (state) => selector(state, a1, a2, a3));
}
//#endregion
//#region node_modules/.pnpm/@base-ui+utils@0.2.6_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/utils/esm/store/Store.js
/**
* A data store implementation that allows subscribing to state changes and updating the state.
* It uses an observer pattern to notify subscribers when the state changes.
*/
var Store = class {
	/**
	* The current state of the store.
	* This property is updated immediately when the state changes as a result of calling {@link setState}, {@link update}, or {@link set}.
	* To subscribe to state changes, use the {@link useState} method. The value returned by {@link useState} is updated after the component renders (similarly to React's useState).
	* The values can be used directly (to avoid subscribing to the store) in effects or event handlers.
	*
	* Do not modify properties in state directly. Instead, use the provided methods to ensure proper state management and listener notification.
	*/
	constructor(state) {
		this.state = state;
		this.listeners = /* @__PURE__ */ new Set();
		this.updateTick = 0;
	}
	/**
	* Registers a listener that will be called whenever the store's state changes.
	*
	* @param fn The listener function to be called on state changes.
	* @returns A function to unsubscribe the listener.
	*/
	subscribe = (fn) => {
		this.listeners.add(fn);
		return () => {
			this.listeners.delete(fn);
		};
	};
	/**
	* Returns the current state of the store.
	*/
	getSnapshot = () => {
		return this.state;
	};
	/**
	* Updates the entire store's state and notifies all registered listeners.
	*
	* @param newState The new state to set for the store.
	*/
	setState(newState) {
		if (this.state === newState) return;
		this.state = newState;
		this.updateTick += 1;
		const currentTick = this.updateTick;
		for (const listener of this.listeners) {
			if (currentTick !== this.updateTick) return;
			listener(newState);
		}
	}
	/**
	* Merges the provided changes into the current state and notifies listeners if there are changes.
	*
	* @param changes An object containing the changes to apply to the current state.
	*/
	update(changes) {
		for (const key in changes) if (!Object.is(this.state[key], changes[key])) {
			this.setState({
				...this.state,
				...changes
			});
			return;
		}
	}
	/**
	* Sets a specific key in the store's state to a new value and notifies listeners if the value has changed.
	*
	* @param key The key in the store's state to update.
	* @param value The new value to set for the specified key.
	*/
	set(key, value) {
		if (!Object.is(this.state[key], value)) this.setState({
			...this.state,
			[key]: value
		});
	}
	/**
	* Gives the state a new reference and updates all registered listeners.
	*/
	notifyAll() {
		const newState = { ...this.state };
		this.setState(newState);
	}
	use(selector, a1, a2, a3) {
		return useStore(this, selector, a1, a2, a3);
	}
};
//#endregion
//#region node_modules/.pnpm/@base-ui+utils@0.2.6_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/utils/esm/store/ReactStore.js
/**
* A Store that supports controlled state keys, non-reactive values and provides utility methods for React.
*/
var ReactStore = class extends Store {
	/**
	* Creates a new ReactStore instance.
	*
	* @param state Initial state of the store.
	* @param context Non-reactive context values.
	* @param selectors Optional selectors for use with `useState`.
	*/
	constructor(state, context = {}, selectors) {
		super(state);
		this.context = context;
		this.selectors = selectors;
	}
	/**
	* Non-reactive values such as refs, callbacks, etc.
	*/
	/**
	* Synchronizes a single external value into the store.
	*
	* Note that the while the value in `state` is updated immediately, the value returned
	* by `useState` is updated before the next render (similarly to React's `useState`).
	*/
	useSyncedValue(key, value) {
		import_react.useDebugValue(key);
		useIsoLayoutEffect(() => {
			if (this.state[key] !== value) this.set(key, value);
		}, [key, value]);
	}
	/**
	* Synchronizes a single external value into the store and
	* cleans it up (sets to `undefined`) on unmount.
	*
	* Note that the while the value in `state` is updated immediately, the value returned
	* by `useState` is updated before the next render (similarly to React's `useState`).
	*/
	useSyncedValueWithCleanup(key, value) {
		const store = this;
		useIsoLayoutEffect(() => {
			if (store.state[key] !== value) store.set(key, value);
			return () => {
				store.set(key, void 0);
			};
		}, [
			store,
			key,
			value
		]);
	}
	/**
	* Synchronizes multiple external values into the store.
	*
	* Note that the while the values in `state` are updated immediately, the values returned
	* by `useState` are updated before the next render (similarly to React's `useState`).
	*/
	useSyncedValues(statePart) {
		const store = this;
		useIsoLayoutEffect(() => {
			store.update(statePart);
		}, [store, ...Object.values(statePart)]);
	}
	/**
	* Registers a controllable prop pair (`controlled`, `defaultValue`) for a specific key. If `controlled`
	* is non-undefined, the store's state at `key` is updated to match `controlled`.
	*/
	useControlledProp(key, controlled) {
		import_react.useDebugValue(key);
		const isControlled = controlled !== void 0;
		useIsoLayoutEffect(() => {
			if (isControlled && !Object.is(this.state[key], controlled)) super.setState({
				...this.state,
				[key]: controlled
			});
		}, [
			key,
			controlled,
			isControlled
		]);
	}
	/** Gets the current value from the store using a selector with the provided key.
	*
	* @param key Key of the selector to use.
	*/
	select(key, a1, a2, a3) {
		const selector = this.selectors[key];
		return selector(this.state, a1, a2, a3);
	}
	/**
	* Returns a value from the store's state using a selector function.
	* Used to subscribe to specific parts of the state.
	* This methods causes a rerender whenever the selected state changes.
	*
	* @param key Key of the selector to use.
	*/
	useState(key, a1, a2, a3) {
		import_react.useDebugValue(key);
		return useStore(this, this.selectors[key], a1, a2, a3);
	}
	/**
	* Wraps a function with `useStableCallback` to ensure it has a stable reference
	* and assigns it to the context.
	*
	* @param key Key of the event callback. Must be a function in the context.
	* @param fn Function to assign.
	*/
	useContextCallback(key, fn) {
		import_react.useDebugValue(key);
		const stableFunction = useStableCallback(fn ?? NOOP);
		this.context[key] = stableFunction;
	}
	/**
	* Returns a stable setter function for a specific key in the store's state.
	* It's commonly used to pass as a ref callback to React elements.
	*
	* @param key Key of the state to set.
	*/
	useStateSetter(key) {
		const ref = import_react.useRef(void 0);
		if (ref.current === void 0) ref.current = (value) => {
			this.set(key, value);
		};
		return ref.current;
	}
	/**
	* Observes changes derived from the store's selectors and calls the listener when the selected value changes.
	*
	* @param key Key of the selector to observe.
	* @param listener Listener function called when the selector result changes.
	*/
	observe(selector, listener) {
		let selectFn;
		if (typeof selector === "function") selectFn = selector;
		else selectFn = this.selectors[selector];
		let prevValue = selectFn(this.state);
		listener(prevValue, prevValue, this);
		return this.subscribe((nextState) => {
			const nextValue = selectFn(nextState);
			if (!Object.is(prevValue, nextValue)) {
				const oldValue = prevValue;
				prevValue = nextValue;
				listener(nextValue, oldValue, this);
			}
		});
	}
};
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/FocusGuard.js
/**
* @internal
*/
var FocusGuard = /* @__PURE__ */ import_react.forwardRef(function FocusGuard(props, ref) {
	const [role, setRole] = import_react.useState();
	useIsoLayoutEffect(() => {
		if (isSafari) setRole("button");
	}, []);
	const restProps = {
		tabIndex: 0,
		role
	};
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("span", {
		...props,
		ref,
		style: visuallyHidden,
		"aria-hidden": role ? void 0 : true,
		...restProps,
		"data-base-ui-focus-guard": ""
	});
});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/utils/createAttribute.js
function createAttribute(name) {
	return `data-base-ui-${name}`;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/utils/enqueueFocus.js
var rafId = 0;
function enqueueFocus(el, options = {}) {
	const { preventScroll = false, cancelPrevious = true, sync = false } = options;
	if (cancelPrevious) cancelAnimationFrame(rafId);
	const exec = () => el?.focus({ preventScroll });
	if (sync) exec();
	else rafId = requestAnimationFrame(exec);
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/utils/markOthers.js
var counters = {
	inert: /* @__PURE__ */ new WeakMap(),
	"aria-hidden": /* @__PURE__ */ new WeakMap()
};
var markerName = "data-base-ui-inert";
var uncontrolledElementsSets = {
	inert: /* @__PURE__ */ new WeakSet(),
	"aria-hidden": /* @__PURE__ */ new WeakSet()
};
var markerCounterMap = /* @__PURE__ */ new WeakMap();
var lockCount = 0;
function getUncontrolledElementsSet(controlAttribute) {
	return uncontrolledElementsSets[controlAttribute];
}
function unwrapHost(node) {
	if (!node) return null;
	return isShadowRoot(node) ? node.host : unwrapHost(node.parentNode);
}
var correctElements = (parent, targets) => targets.map((target) => {
	if (parent.contains(target)) return target;
	const correctedTarget = unwrapHost(target);
	if (parent.contains(correctedTarget)) return correctedTarget;
	return null;
}).filter((x) => x != null);
var buildKeepSet = (targets) => {
	const keep = /* @__PURE__ */ new Set();
	targets.forEach((target) => {
		let node = target;
		while (node && !keep.has(node)) {
			keep.add(node);
			node = node.parentNode;
		}
	});
	return keep;
};
var collectOutsideElements = (root, keepElements, stopElements) => {
	const outside = [];
	const walk = (parent) => {
		if (!parent || stopElements.has(parent)) return;
		Array.from(parent.children).forEach((node) => {
			if (getNodeName(node) === "script") return;
			if (keepElements.has(node)) walk(node);
			else outside.push(node);
		});
	};
	walk(root);
	return outside;
};
function applyAttributeToOthers(uncorrectedAvoidElements, body, ariaHidden, inert, { mark = true, markerIgnoreElements = [] }) {
	const controlAttribute = inert ? "inert" : ariaHidden ? "aria-hidden" : null;
	let counterMap = null;
	let uncontrolledElementsSet = null;
	const avoidElements = correctElements(body, uncorrectedAvoidElements);
	const markerIgnoreTargets = mark ? correctElements(body, markerIgnoreElements) : [];
	const markerIgnoreSet = new Set(markerIgnoreTargets);
	const markerTargets = mark ? collectOutsideElements(body, buildKeepSet(avoidElements), new Set(avoidElements)).filter((target) => !markerIgnoreSet.has(target)) : [];
	const hiddenElements = [];
	const markedElements = [];
	if (controlAttribute) {
		const map = counters[controlAttribute];
		const currentUncontrolledElementsSet = getUncontrolledElementsSet(controlAttribute);
		uncontrolledElementsSet = currentUncontrolledElementsSet;
		counterMap = map;
		const ariaLiveElements = correctElements(body, Array.from(body.querySelectorAll("[aria-live]")));
		const controlElements = avoidElements.concat(ariaLiveElements);
		collectOutsideElements(body, buildKeepSet(controlElements), new Set(controlElements)).forEach((node) => {
			const attr = node.getAttribute(controlAttribute);
			const alreadyHidden = attr !== null && attr !== "false";
			const counterValue = (map.get(node) || 0) + 1;
			map.set(node, counterValue);
			hiddenElements.push(node);
			if (counterValue === 1 && alreadyHidden) currentUncontrolledElementsSet.add(node);
			if (!alreadyHidden) node.setAttribute(controlAttribute, controlAttribute === "inert" ? "" : "true");
		});
	}
	if (mark) markerTargets.forEach((node) => {
		const markerValue = (markerCounterMap.get(node) || 0) + 1;
		markerCounterMap.set(node, markerValue);
		markedElements.push(node);
		if (markerValue === 1) node.setAttribute(markerName, "");
	});
	lockCount += 1;
	return () => {
		if (counterMap) hiddenElements.forEach((element) => {
			const counterValue = (counterMap.get(element) || 0) - 1;
			counterMap.set(element, counterValue);
			if (!counterValue) {
				if (!uncontrolledElementsSet?.has(element) && controlAttribute) element.removeAttribute(controlAttribute);
				uncontrolledElementsSet?.delete(element);
			}
		});
		if (mark) markedElements.forEach((element) => {
			const markerValue = (markerCounterMap.get(element) || 0) - 1;
			markerCounterMap.set(element, markerValue);
			if (!markerValue) element.removeAttribute(markerName);
		});
		lockCount -= 1;
		if (!lockCount) {
			counters.inert = /* @__PURE__ */ new WeakMap();
			counters["aria-hidden"] = /* @__PURE__ */ new WeakMap();
			uncontrolledElementsSets.inert = /* @__PURE__ */ new WeakSet();
			uncontrolledElementsSets["aria-hidden"] = /* @__PURE__ */ new WeakSet();
			markerCounterMap = /* @__PURE__ */ new WeakMap();
		}
	};
}
function markOthers(avoidElements, options = {}) {
	const { ariaHidden = false, inert = false, mark = true, markerIgnoreElements = [] } = options;
	const body = ownerDocument(avoidElements[0]).body;
	return applyAttributeToOthers(avoidElements, body, ariaHidden, inert, {
		mark,
		markerIgnoreElements
	});
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/components/FloatingPortal.js
var PortalContext = /* @__PURE__ */ import_react.createContext(null);
var usePortalContext = () => import_react.useContext(PortalContext);
var attr = createAttribute("portal");
function useFloatingPortalNode(props = {}) {
	const { ref, container: containerProp, componentProps = EMPTY_OBJECT, elementProps } = props;
	const uniqueId = useId();
	const parentPortalNode = usePortalContext()?.portalNode;
	const [containerElement, setContainerElement] = import_react.useState(null);
	const [portalNode, setPortalNode] = import_react.useState(null);
	const setPortalNodeRef = useStableCallback((node) => {
		if (node !== null) setPortalNode(node);
	});
	const containerRef = import_react.useRef(null);
	useIsoLayoutEffect(() => {
		if (containerProp === null) {
			if (containerRef.current) {
				containerRef.current = null;
				setPortalNode(null);
				setContainerElement(null);
			}
			return;
		}
		if (uniqueId == null) return;
		const resolvedContainer = (containerProp && (isNode(containerProp) ? containerProp : containerProp.current)) ?? parentPortalNode ?? document.body;
		if (resolvedContainer == null) {
			if (containerRef.current) {
				containerRef.current = null;
				setPortalNode(null);
				setContainerElement(null);
			}
			return;
		}
		if (containerRef.current !== resolvedContainer) {
			containerRef.current = resolvedContainer;
			setPortalNode(null);
			setContainerElement(resolvedContainer);
		}
	}, [
		containerProp,
		parentPortalNode,
		uniqueId
	]);
	const portalElement = useRenderElement("div", componentProps, {
		ref: [ref, setPortalNodeRef],
		props: [{
			id: uniqueId,
			[attr]: ""
		}, elementProps]
	});
	return {
		portalNode,
		portalSubtree: containerElement && portalElement ? /* @__PURE__ */ import_react_dom.createPortal(portalElement, containerElement) : null
	};
}
/**
* Portals the floating element into a given container element — by default,
* outside of the app root and into the body.
* This is necessary to ensure the floating element can appear outside any
* potential parent containers that cause clipping (such as `overflow: hidden`),
* while retaining its location in the React tree.
* @see https://floating-ui.com/docs/FloatingPortal
* @internal
*/
var FloatingPortal = /* @__PURE__ */ import_react.forwardRef(function FloatingPortal(componentProps, forwardedRef) {
	const { children, container, className, render, renderGuards, ...elementProps } = componentProps;
	const { portalNode, portalSubtree } = useFloatingPortalNode({
		container,
		ref: forwardedRef,
		componentProps,
		elementProps
	});
	const beforeOutsideRef = import_react.useRef(null);
	const afterOutsideRef = import_react.useRef(null);
	const beforeInsideRef = import_react.useRef(null);
	const afterInsideRef = import_react.useRef(null);
	const [focusManagerState, setFocusManagerState] = import_react.useState(null);
	const modal = focusManagerState?.modal;
	const open = focusManagerState?.open;
	const shouldRenderGuards = typeof renderGuards === "boolean" ? renderGuards : !!focusManagerState && !focusManagerState.modal && focusManagerState.open && !!portalNode;
	import_react.useEffect(() => {
		if (!portalNode || modal) return;
		function onFocus(event) {
			if (portalNode && event.relatedTarget && isOutsideEvent(event)) (event.type === "focusin" ? enableFocusInside : disableFocusInside)(portalNode);
		}
		portalNode.addEventListener("focusin", onFocus, true);
		portalNode.addEventListener("focusout", onFocus, true);
		return () => {
			portalNode.removeEventListener("focusin", onFocus, true);
			portalNode.removeEventListener("focusout", onFocus, true);
		};
	}, [portalNode, modal]);
	import_react.useEffect(() => {
		if (!portalNode || open) return;
		enableFocusInside(portalNode);
	}, [open, portalNode]);
	const portalContextValue = import_react.useMemo(() => ({
		beforeOutsideRef,
		afterOutsideRef,
		beforeInsideRef,
		afterInsideRef,
		portalNode,
		setFocusManagerState
	}), [portalNode]);
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(import_react.Fragment, { children: [portalSubtree, /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(PortalContext.Provider, {
		value: portalContextValue,
		children: [
			shouldRenderGuards && portalNode && /* @__PURE__ */ (0, import_jsx_runtime.jsx)(FocusGuard, {
				"data-type": "outside",
				ref: beforeOutsideRef,
				onFocus: (event) => {
					if (isOutsideEvent(event, portalNode)) beforeInsideRef.current?.focus();
					else getPreviousTabbable(focusManagerState ? focusManagerState.domReference : null)?.focus();
				}
			}),
			shouldRenderGuards && portalNode && /* @__PURE__ */ (0, import_jsx_runtime.jsx)("span", {
				"aria-owns": portalNode.id,
				style: ownerVisuallyHidden
			}),
			portalNode && /* @__PURE__ */ import_react_dom.createPortal(children, portalNode),
			shouldRenderGuards && portalNode && /* @__PURE__ */ (0, import_jsx_runtime.jsx)(FocusGuard, {
				"data-type": "outside",
				ref: afterOutsideRef,
				onFocus: (event) => {
					if (isOutsideEvent(event, portalNode)) afterInsideRef.current?.focus();
					else {
						getNextTabbable(focusManagerState ? focusManagerState.domReference : null)?.focus();
						if (focusManagerState?.closeOnFocusOut) focusManagerState?.onOpenChange(false, createChangeEventDetails("focus-out", event.nativeEvent));
					}
				}
			})
		]
	})] });
});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/utils/createEventEmitter.js
function createEventEmitter() {
	const map = /* @__PURE__ */ new Map();
	return {
		emit(event, data) {
			map.get(event)?.forEach((listener) => listener(data));
		},
		on(event, listener) {
			if (!map.has(event)) map.set(event, /* @__PURE__ */ new Set());
			map.get(event).add(listener);
		},
		off(event, listener) {
			map.get(event)?.delete(listener);
		}
	};
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/components/FloatingTree.js
var FloatingNodeContext = /* @__PURE__ */ import_react.createContext(null);
var FloatingTreeContext = /* @__PURE__ */ import_react.createContext(null);
var useFloatingParentNodeId = () => import_react.useContext(FloatingNodeContext)?.id || null;
/**
* Returns the nearest floating tree context, if available.
*/
var useFloatingTree = (externalTree) => {
	const contextTree = import_react.useContext(FloatingTreeContext);
	return externalTree ?? contextTree;
};
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/components/FloatingFocusManager.js
function getEventType(event, lastInteractionType) {
	const win = getWindow(event.target);
	if (event instanceof win.KeyboardEvent) return "keyboard";
	if (event instanceof win.FocusEvent) return lastInteractionType || "keyboard";
	if ("pointerType" in event) return event.pointerType || "keyboard";
	if ("touches" in event) return "touch";
	if (event instanceof win.MouseEvent) return lastInteractionType || (event.detail === 0 ? "keyboard" : "mouse");
	return "";
}
var LIST_LIMIT = 20;
var previouslyFocusedElements = [];
function clearDisconnectedPreviouslyFocusedElements() {
	previouslyFocusedElements = previouslyFocusedElements.filter((entry) => {
		return entry.deref()?.isConnected;
	});
}
function addPreviouslyFocusedElement(element) {
	clearDisconnectedPreviouslyFocusedElements();
	if (element && getNodeName(element) !== "body") {
		previouslyFocusedElements.push(new WeakRef(element));
		if (previouslyFocusedElements.length > LIST_LIMIT) previouslyFocusedElements = previouslyFocusedElements.slice(-LIST_LIMIT);
	}
}
function getPreviouslyFocusedElement() {
	clearDisconnectedPreviouslyFocusedElements();
	return previouslyFocusedElements[previouslyFocusedElements.length - 1]?.deref();
}
function getFirstTabbableElement(container) {
	if (!container) return null;
	const tabbableOptions = getTabbableOptions();
	if (isTabbable(container, tabbableOptions)) return container;
	return tabbable(container, tabbableOptions)[0] || container;
}
function isFocusable(element) {
	if (!element || !element.isConnected) return false;
	if (typeof element.checkVisibility === "function") return element.checkVisibility();
	return isElementVisible(element);
}
function handleTabIndex(floatingFocusElement, orderRef) {
	if (!orderRef.current.includes("floating") && !floatingFocusElement.getAttribute("role")?.includes("dialog")) return;
	const options = getTabbableOptions();
	const tabbableContent = focusable(floatingFocusElement, options).filter((element) => {
		const dataTabIndex = element.getAttribute("data-tabindex") || "";
		return isTabbable(element, options) || element.hasAttribute("data-tabindex") && !dataTabIndex.startsWith("-");
	});
	const tabIndex = floatingFocusElement.getAttribute("tabindex");
	if (orderRef.current.includes("floating") || tabbableContent.length === 0) {
		if (tabIndex !== "0") floatingFocusElement.setAttribute("tabindex", "0");
	} else if (tabIndex !== "-1" || floatingFocusElement.hasAttribute("data-tabindex") && floatingFocusElement.getAttribute("data-tabindex") !== "-1") {
		floatingFocusElement.setAttribute("tabindex", "-1");
		floatingFocusElement.setAttribute("data-tabindex", "-1");
	}
}
/**
* Provides focus management for the floating element.
* @see https://floating-ui.com/docs/FloatingFocusManager
* @internal
*/
function FloatingFocusManager(props) {
	const { context, children, disabled = false, initialFocus = true, returnFocus = true, restoreFocus = false, modal = true, closeOnFocusOut = true, openInteractionType = "", nextFocusableElement, previousFocusableElement, beforeContentFocusGuardRef, externalTree, getInsideElements } = props;
	const store = "rootStore" in context ? context.rootStore : context;
	const open = store.useState("open");
	const domReference = store.useState("domReferenceElement");
	const floating = store.useState("floatingElement");
	const { events, dataRef } = store.context;
	const getNodeId = useStableCallback(() => dataRef.current.floatingContext?.nodeId);
	const ignoreInitialFocus = initialFocus === false;
	const isUntrappedTypeableCombobox = isTypeableCombobox(domReference) && ignoreInitialFocus;
	const orderRef = import_react.useRef(["content"]);
	const initialFocusRef = useValueAsRef(initialFocus);
	const returnFocusRef = useValueAsRef(returnFocus);
	const openInteractionTypeRef = useValueAsRef(openInteractionType);
	const tree = useFloatingTree(externalTree);
	const portalContext = usePortalContext();
	const preventReturnFocusRef = import_react.useRef(false);
	const isPointerDownRef = import_react.useRef(false);
	const pointerDownOutsideRef = import_react.useRef(false);
	const tabbableIndexRef = import_react.useRef(-1);
	const closeTypeRef = import_react.useRef("");
	const lastInteractionTypeRef = import_react.useRef("");
	const beforeGuardRef = import_react.useRef(null);
	const afterGuardRef = import_react.useRef(null);
	const mergedBeforeGuardRef = useMergedRefs(beforeGuardRef, beforeContentFocusGuardRef, portalContext?.beforeInsideRef);
	const mergedAfterGuardRef = useMergedRefs(afterGuardRef, portalContext?.afterInsideRef);
	const blurTimeout = useTimeout();
	const pointerDownTimeout = useTimeout();
	const restoreFocusFrame = useAnimationFrame();
	const isInsidePortal = portalContext != null;
	const floatingFocusElement = getFloatingFocusElement(floating);
	const getTabbableContent = useStableCallback((container = floatingFocusElement) => {
		return container ? tabbable(container, getTabbableOptions()) : [];
	});
	const getResolvedInsideElements = useStableCallback(() => getInsideElements?.().filter((element) => element != null) ?? []);
	const getTabbableElements = useStableCallback((container) => {
		const content = getTabbableContent(container);
		return orderRef.current.map(() => content).filter(Boolean).flat();
	});
	import_react.useEffect(() => {
		if (disabled || !modal) return;
		function onKeyDown(event) {
			if (event.key === "Tab") {
				if (contains(floatingFocusElement, activeElement(ownerDocument(floatingFocusElement))) && getTabbableContent().length === 0 && !isUntrappedTypeableCombobox) stopEvent(event);
			}
		}
		const doc = ownerDocument(floatingFocusElement);
		doc.addEventListener("keydown", onKeyDown);
		return () => {
			doc.removeEventListener("keydown", onKeyDown);
		};
	}, [
		disabled,
		domReference,
		floatingFocusElement,
		modal,
		orderRef,
		isUntrappedTypeableCombobox,
		getTabbableContent,
		getTabbableElements
	]);
	import_react.useEffect(() => {
		if (disabled || !open) return;
		const doc = ownerDocument(floatingFocusElement);
		function clearPointerDownOutside() {
			pointerDownOutsideRef.current = false;
		}
		function onPointerDown(event) {
			const target = getTarget(event);
			const insideElements = getResolvedInsideElements();
			pointerDownOutsideRef.current = !(contains(floating, target) || contains(domReference, target) || contains(portalContext?.portalNode, target) || insideElements.some((element) => element === target || contains(element, target)));
			lastInteractionTypeRef.current = event.pointerType || "keyboard";
			if (target?.closest(`[data-base-ui-click-trigger]`)) isPointerDownRef.current = true;
		}
		function onKeyDown() {
			lastInteractionTypeRef.current = "keyboard";
		}
		doc.addEventListener("pointerdown", onPointerDown, true);
		doc.addEventListener("pointerup", clearPointerDownOutside, true);
		doc.addEventListener("pointercancel", clearPointerDownOutside, true);
		doc.addEventListener("keydown", onKeyDown, true);
		return () => {
			doc.removeEventListener("pointerdown", onPointerDown, true);
			doc.removeEventListener("pointerup", clearPointerDownOutside, true);
			doc.removeEventListener("pointercancel", clearPointerDownOutside, true);
			doc.removeEventListener("keydown", onKeyDown, true);
		};
	}, [
		disabled,
		floating,
		domReference,
		floatingFocusElement,
		open,
		portalContext,
		getResolvedInsideElements
	]);
	import_react.useEffect(() => {
		if (disabled || !closeOnFocusOut) return;
		const doc = ownerDocument(floatingFocusElement);
		function handlePointerDown() {
			isPointerDownRef.current = true;
			pointerDownTimeout.start(0, () => {
				isPointerDownRef.current = false;
			});
		}
		function handleFocusIn(event) {
			const target = getTarget(event);
			const tabbableIndex = getTabbableContent().indexOf(target);
			if (tabbableIndex !== -1) tabbableIndexRef.current = tabbableIndex;
		}
		function handleFocusOutside(event) {
			const relatedTarget = event.relatedTarget;
			const currentTarget = event.currentTarget;
			const target = getTarget(event);
			queueMicrotask(() => {
				const nodeId = getNodeId();
				const triggers = store.context.triggerElements;
				const insideElements = getResolvedInsideElements();
				const isRelatedFocusGuard = relatedTarget?.hasAttribute(createAttribute("focus-guard")) && [
					beforeGuardRef.current,
					afterGuardRef.current,
					portalContext?.beforeInsideRef.current,
					portalContext?.afterInsideRef.current,
					portalContext?.beforeOutsideRef.current,
					portalContext?.afterOutsideRef.current,
					resolveRef(previousFocusableElement),
					resolveRef(nextFocusableElement)
				].includes(relatedTarget);
				const movedToUnrelatedNode = !(contains(domReference, relatedTarget) || contains(floating, relatedTarget) || contains(relatedTarget, floating) || contains(portalContext?.portalNode, relatedTarget) || insideElements.some((element) => element === relatedTarget || contains(element, relatedTarget)) || relatedTarget != null && triggers.hasElement(relatedTarget) || triggers.hasMatchingElement((trigger) => contains(trigger, relatedTarget)) || isRelatedFocusGuard || tree && (getNodeChildren(tree.nodesRef.current, nodeId).find((node) => contains(node.context?.elements.floating, relatedTarget) || contains(node.context?.elements.domReference, relatedTarget)) || getNodeAncestors(tree.nodesRef.current, nodeId).find((node) => [node.context?.elements.floating, getFloatingFocusElement(node.context?.elements.floating)].includes(relatedTarget) || node.context?.elements.domReference === relatedTarget)));
				if (currentTarget === domReference && floatingFocusElement) handleTabIndex(floatingFocusElement, orderRef);
				if (restoreFocus && currentTarget !== domReference && !isFocusable(target) && activeElement(doc) === doc.body) {
					if (isHTMLElement(floatingFocusElement)) {
						floatingFocusElement.focus();
						if (restoreFocus === "popup") {
							restoreFocusFrame.request(() => {
								floatingFocusElement.focus();
							});
							return;
						}
					}
					const prevTabbableIndex = tabbableIndexRef.current;
					const tabbableContent = getTabbableContent();
					const nodeToFocus = tabbableContent[prevTabbableIndex] || tabbableContent[tabbableContent.length - 1] || floatingFocusElement;
					if (isHTMLElement(nodeToFocus)) nodeToFocus.focus();
				}
				if (dataRef.current.insideReactTree) {
					dataRef.current.insideReactTree = false;
					return;
				}
				if ((isUntrappedTypeableCombobox ? true : !modal) && relatedTarget && movedToUnrelatedNode && !isPointerDownRef.current && (isUntrappedTypeableCombobox || relatedTarget !== getPreviouslyFocusedElement())) {
					preventReturnFocusRef.current = true;
					store.setOpen(false, createChangeEventDetails(focusOut, event));
				}
			});
		}
		function markInsideReactTree() {
			if (pointerDownOutsideRef.current) return;
			dataRef.current.insideReactTree = true;
			blurTimeout.start(0, () => {
				dataRef.current.insideReactTree = false;
			});
		}
		const domReferenceElement = isHTMLElement(domReference) ? domReference : null;
		const cleanups = [];
		if (!floating && !domReferenceElement) return;
		if (domReferenceElement) {
			domReferenceElement.addEventListener("focusout", handleFocusOutside);
			domReferenceElement.addEventListener("pointerdown", handlePointerDown);
			cleanups.push(() => {
				domReferenceElement.removeEventListener("focusout", handleFocusOutside);
				domReferenceElement.removeEventListener("pointerdown", handlePointerDown);
			});
		}
		if (floating) {
			floating.addEventListener("focusin", handleFocusIn);
			floating.addEventListener("focusout", handleFocusOutside);
			if (portalContext) {
				floating.addEventListener("focusout", markInsideReactTree, true);
				cleanups.push(() => {
					floating.removeEventListener("focusout", markInsideReactTree, true);
				});
			}
			cleanups.push(() => {
				floating.removeEventListener("focusin", handleFocusIn);
				floating.removeEventListener("focusout", handleFocusOutside);
			});
		}
		return () => {
			cleanups.forEach((cleanup) => {
				cleanup();
			});
		};
	}, [
		disabled,
		domReference,
		floating,
		floatingFocusElement,
		modal,
		tree,
		portalContext,
		store,
		closeOnFocusOut,
		restoreFocus,
		getTabbableContent,
		isUntrappedTypeableCombobox,
		getNodeId,
		orderRef,
		dataRef,
		blurTimeout,
		pointerDownTimeout,
		restoreFocusFrame,
		nextFocusableElement,
		previousFocusableElement,
		getResolvedInsideElements
	]);
	import_react.useEffect(() => {
		if (disabled || !floating || !open) return;
		const portalNodes = Array.from(portalContext?.portalNode?.querySelectorAll(`[${createAttribute("portal")}]`) || []);
		const rootAncestorComboboxDomReference = (tree ? getNodeAncestors(tree.nodesRef.current, getNodeId()) : []).find((node) => isTypeableCombobox(node.context?.elements.domReference || null))?.context?.elements.domReference;
		const ariaHiddenCleanup = markOthers([
			...[
				floating,
				...portalNodes,
				beforeGuardRef.current,
				afterGuardRef.current,
				portalContext?.beforeOutsideRef.current,
				portalContext?.afterOutsideRef.current,
				...getResolvedInsideElements()
			],
			rootAncestorComboboxDomReference,
			resolveRef(previousFocusableElement),
			resolveRef(nextFocusableElement),
			isUntrappedTypeableCombobox ? domReference : null
		].filter((x) => x != null), {
			ariaHidden: modal || isUntrappedTypeableCombobox,
			mark: false
		});
		const markerCleanup = markOthers([floating, ...portalNodes].filter((x) => x != null));
		return () => {
			markerCleanup();
			ariaHiddenCleanup();
		};
	}, [
		open,
		disabled,
		domReference,
		floating,
		modal,
		orderRef,
		portalContext,
		isUntrappedTypeableCombobox,
		tree,
		getNodeId,
		nextFocusableElement,
		previousFocusableElement,
		getResolvedInsideElements
	]);
	useIsoLayoutEffect(() => {
		if (!open || disabled || !isHTMLElement(floatingFocusElement)) return;
		const previouslyFocusedElement = activeElement(ownerDocument(floatingFocusElement));
		queueMicrotask(() => {
			const focusableElements = getTabbableElements(floatingFocusElement);
			const initialFocusValueOrFn = initialFocusRef.current;
			const resolvedInitialFocus = typeof initialFocusValueOrFn === "function" ? initialFocusValueOrFn(openInteractionTypeRef.current || "") : initialFocusValueOrFn;
			if (resolvedInitialFocus === void 0 || resolvedInitialFocus === false) return;
			let elToFocus;
			if (resolvedInitialFocus === true || resolvedInitialFocus === null) elToFocus = focusableElements[0] || floatingFocusElement;
			else elToFocus = resolveRef(resolvedInitialFocus);
			elToFocus = elToFocus || focusableElements[0] || floatingFocusElement;
			if (contains(floatingFocusElement, previouslyFocusedElement)) return;
			enqueueFocus(elToFocus, { preventScroll: elToFocus === floatingFocusElement });
		});
	}, [
		disabled,
		open,
		floatingFocusElement,
		ignoreInitialFocus,
		getTabbableElements,
		initialFocusRef,
		openInteractionTypeRef
	]);
	useIsoLayoutEffect(() => {
		if (disabled || !floatingFocusElement) return;
		const doc = ownerDocument(floatingFocusElement);
		addPreviouslyFocusedElement(activeElement(doc));
		function onOpenChangeLocal(details) {
			if (!details.open) closeTypeRef.current = getEventType(details.nativeEvent, lastInteractionTypeRef.current);
			if (details.reason === "trigger-hover" && details.nativeEvent.type === "mouseleave") preventReturnFocusRef.current = true;
			if (details.reason !== "outside-press") return;
			if (details.nested) preventReturnFocusRef.current = false;
			else if (isVirtualClick(details.nativeEvent) || isVirtualPointerEvent(details.nativeEvent)) preventReturnFocusRef.current = false;
			else {
				let isPreventScrollSupported = false;
				document.createElement("div").focus({ get preventScroll() {
					isPreventScrollSupported = true;
					return false;
				} });
				if (isPreventScrollSupported) preventReturnFocusRef.current = false;
				else preventReturnFocusRef.current = true;
			}
		}
		events.on("openchange", onOpenChangeLocal);
		function getReturnElement() {
			const returnFocusValueOrFn = returnFocusRef.current;
			let resolvedReturnFocusValue = typeof returnFocusValueOrFn === "function" ? returnFocusValueOrFn(closeTypeRef.current) : returnFocusValueOrFn;
			if (resolvedReturnFocusValue === void 0 || resolvedReturnFocusValue === false) return null;
			if (resolvedReturnFocusValue === null) resolvedReturnFocusValue = true;
			if (typeof resolvedReturnFocusValue === "boolean") {
				const el = domReference || getPreviouslyFocusedElement();
				return el && el.isConnected ? el : null;
			}
			const fallback = domReference || getPreviouslyFocusedElement();
			return resolveRef(resolvedReturnFocusValue) || fallback || null;
		}
		return () => {
			events.off("openchange", onOpenChangeLocal);
			const activeEl = activeElement(doc);
			const insideElements = getResolvedInsideElements();
			const isFocusInsideFloatingTree = contains(floating, activeEl) || insideElements.some((element) => element === activeEl || contains(element, activeEl)) || tree && getNodeChildren(tree.nodesRef.current, getNodeId(), false).some((node) => contains(node.context?.elements.floating, activeEl));
			const returnElement = getReturnElement();
			queueMicrotask(() => {
				const tabbableReturnElement = getFirstTabbableElement(returnElement);
				const hasExplicitReturnFocus = typeof returnFocusRef.current !== "boolean";
				if (returnFocusRef.current && !preventReturnFocusRef.current && isHTMLElement(tabbableReturnElement) && (!hasExplicitReturnFocus && tabbableReturnElement !== activeEl && activeEl !== doc.body ? isFocusInsideFloatingTree : true)) tabbableReturnElement.focus({ preventScroll: true });
				preventReturnFocusRef.current = false;
			});
		};
	}, [
		disabled,
		floating,
		floatingFocusElement,
		returnFocusRef,
		dataRef,
		events,
		tree,
		domReference,
		getNodeId,
		getResolvedInsideElements
	]);
	useIsoLayoutEffect(() => {
		if (!isWebKit$1 || open || !floating) return;
		const activeEl = activeElement(ownerDocument(floating));
		if (!isHTMLElement(activeEl) || !isTypeableElement(activeEl)) return;
		if (contains(floating, activeEl)) activeEl.blur();
	}, [open, floating]);
	useIsoLayoutEffect(() => {
		if (disabled || !portalContext) return;
		portalContext.setFocusManagerState({
			modal,
			closeOnFocusOut,
			open,
			onOpenChange: store.setOpen,
			domReference
		});
		return () => {
			portalContext.setFocusManagerState(null);
		};
	}, [
		disabled,
		portalContext,
		modal,
		open,
		store,
		closeOnFocusOut,
		domReference
	]);
	useIsoLayoutEffect(() => {
		if (disabled || !floatingFocusElement) return;
		handleTabIndex(floatingFocusElement, orderRef);
		return () => {
			queueMicrotask(clearDisconnectedPreviouslyFocusedElements);
		};
	}, [
		disabled,
		floatingFocusElement,
		orderRef
	]);
	const shouldRenderGuards = !disabled && (modal ? !isUntrappedTypeableCombobox : true) && (isInsidePortal || modal);
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(import_react.Fragment, { children: [
		shouldRenderGuards && /* @__PURE__ */ (0, import_jsx_runtime.jsx)(FocusGuard, {
			"data-type": "inside",
			ref: mergedBeforeGuardRef,
			onFocus: (event) => {
				if (modal) {
					const els = getTabbableElements();
					enqueueFocus(els[els.length - 1]);
				} else if (portalContext?.portalNode) {
					preventReturnFocusRef.current = false;
					if (isOutsideEvent(event, portalContext.portalNode)) getNextTabbable(domReference)?.focus();
					else resolveRef(previousFocusableElement ?? portalContext.beforeOutsideRef)?.focus();
				}
			}
		}),
		children,
		shouldRenderGuards && /* @__PURE__ */ (0, import_jsx_runtime.jsx)(FocusGuard, {
			"data-type": "inside",
			ref: mergedAfterGuardRef,
			onFocus: (event) => {
				if (modal) enqueueFocus(getTabbableElements()[0]);
				else if (portalContext?.portalNode) {
					if (closeOnFocusOut) preventReturnFocusRef.current = true;
					if (isOutsideEvent(event, portalContext.portalNode)) getPreviousTabbable(domReference)?.focus();
					else resolveRef(nextFocusableElement ?? portalContext.afterOutsideRef)?.focus();
				}
			}
		})
	] });
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/hooks/useClick.js
/**
* Opens or closes the floating element when clicking the reference element.
* @see https://floating-ui.com/docs/useClick
*/
function useClick(context, props = {}) {
	const store = "rootStore" in context ? context.rootStore : context;
	const dataRef = store.context.dataRef;
	const { enabled = true, event: eventOption = "click", toggle = true, ignoreMouse = false, stickIfOpen = true, touchOpenDelay = 0, reason = triggerPress } = props;
	const pointerTypeRef = import_react.useRef(void 0);
	const frame = useAnimationFrame();
	const touchOpenTimeout = useTimeout();
	const reference = import_react.useMemo(() => ({
		onPointerDown(event) {
			pointerTypeRef.current = event.pointerType;
		},
		onMouseDown(event) {
			const pointerType = pointerTypeRef.current;
			const nativeEvent = event.nativeEvent;
			const open = store.select("open");
			if (event.button !== 0 || eventOption === "click" || isMouseLikePointerType(pointerType, true) && ignoreMouse) return;
			const openEvent = dataRef.current.openEvent;
			const openEventType = openEvent?.type;
			const hasClickedOnInactiveTrigger = store.select("domReferenceElement") !== event.currentTarget;
			const nextOpen = open && hasClickedOnInactiveTrigger || !(open && toggle && (openEvent && stickIfOpen ? openEventType === "click" || openEventType === "mousedown" : true));
			if (isTypeableElement(nativeEvent.target)) {
				const details = createChangeEventDetails(reason, nativeEvent, nativeEvent.target);
				if (nextOpen && pointerType === "touch" && touchOpenDelay > 0) touchOpenTimeout.start(touchOpenDelay, () => {
					store.setOpen(true, details);
				});
				else store.setOpen(nextOpen, details);
				return;
			}
			const eventCurrentTarget = event.currentTarget;
			frame.request(() => {
				const details = createChangeEventDetails(reason, nativeEvent, eventCurrentTarget);
				if (nextOpen && pointerType === "touch" && touchOpenDelay > 0) touchOpenTimeout.start(touchOpenDelay, () => {
					store.setOpen(true, details);
				});
				else store.setOpen(nextOpen, details);
			});
		},
		onClick(event) {
			if (eventOption === "mousedown-only") return;
			const pointerType = pointerTypeRef.current;
			if (eventOption === "mousedown" && pointerType) {
				pointerTypeRef.current = void 0;
				return;
			}
			if (isMouseLikePointerType(pointerType, true) && ignoreMouse) return;
			const open = store.select("open");
			const openEvent = dataRef.current.openEvent;
			const hasClickedOnInactiveTrigger = store.select("domReferenceElement") !== event.currentTarget;
			const nextOpen = open && hasClickedOnInactiveTrigger || !(open && toggle && (openEvent && stickIfOpen ? isClickLikeEvent(openEvent) : true));
			const details = createChangeEventDetails(reason, event.nativeEvent, event.currentTarget);
			if (nextOpen && pointerType === "touch" && touchOpenDelay > 0) touchOpenTimeout.start(touchOpenDelay, () => {
				store.setOpen(true, details);
			});
			else store.setOpen(nextOpen, details);
		},
		onKeyDown() {
			pointerTypeRef.current = void 0;
		}
	}), [
		dataRef,
		eventOption,
		ignoreMouse,
		store,
		stickIfOpen,
		toggle,
		frame,
		touchOpenTimeout,
		touchOpenDelay,
		reason
	]);
	return import_react.useMemo(() => enabled ? { reference } : EMPTY_OBJECT, [enabled, reference]);
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/hooks/useDismiss.js
var bubbleHandlerKeys = {
	intentional: "onClick",
	sloppy: "onPointerDown"
};
function alwaysFalse() {
	return false;
}
function normalizeProp(normalizable) {
	return {
		escapeKey: typeof normalizable === "boolean" ? normalizable : normalizable?.escapeKey ?? false,
		outsidePress: typeof normalizable === "boolean" ? normalizable : normalizable?.outsidePress ?? true
	};
}
/**
* Closes the floating element when a dismissal is requested — by default, when
* the user presses the `escape` key or outside of the floating element.
* @see https://floating-ui.com/docs/useDismiss
*/
function useDismiss(context, props = {}) {
	const store = "rootStore" in context ? context.rootStore : context;
	const open = store.useState("open");
	const floatingElement = store.useState("floatingElement");
	const { dataRef } = store.context;
	const { enabled = true, escapeKey: escapeKey$1 = true, outsidePress: outsidePressProp = true, outsidePressEvent = "sloppy", referencePress = alwaysFalse, referencePressEvent = "sloppy", bubbles, externalTree } = props;
	const tree = useFloatingTree(externalTree);
	const outsidePressFn = useStableCallback(typeof outsidePressProp === "function" ? outsidePressProp : () => false);
	const outsidePress$1 = typeof outsidePressProp === "function" ? outsidePressFn : outsidePressProp;
	const outsidePressEnabled = outsidePress$1 !== false;
	const getOutsidePressEventProp = useStableCallback(() => outsidePressEvent);
	const pressStartedInsideRef = import_react.useRef(false);
	const pressStartPreventedRef = import_react.useRef(false);
	const suppressNextOutsideClickRef = import_react.useRef(false);
	const { escapeKey: escapeKeyBubbles, outsidePress: outsidePressBubbles } = normalizeProp(bubbles);
	const touchStateRef = import_react.useRef(null);
	const cancelDismissOnEndTimeout = useTimeout();
	const clearInsideReactTreeTimeout = useTimeout();
	const clearInsideReactTree = useStableCallback(() => {
		clearInsideReactTreeTimeout.clear();
		dataRef.current.insideReactTree = false;
	});
	const isComposingRef = import_react.useRef(false);
	const currentPointerTypeRef = import_react.useRef("");
	const isReferencePressEnabled = useStableCallback(referencePress);
	const closeOnEscapeKeyDown = useStableCallback((event) => {
		if (!open || !enabled || !escapeKey$1 || event.key !== "Escape") return;
		if (isComposingRef.current) return;
		const nodeId = dataRef.current.floatingContext?.nodeId;
		const children = tree ? getNodeChildren(tree.nodesRef.current, nodeId) : [];
		if (!escapeKeyBubbles) {
			if (children.length > 0) {
				let shouldDismiss = true;
				children.forEach((child) => {
					if (child.context?.open && !child.context.dataRef.current.__escapeKeyBubbles) shouldDismiss = false;
				});
				if (!shouldDismiss) return;
			}
		}
		const eventDetails = createChangeEventDetails(escapeKey, isReactEvent(event) ? event.nativeEvent : event);
		store.setOpen(false, eventDetails);
		if (!escapeKeyBubbles && !eventDetails.isPropagationAllowed) event.stopPropagation();
	});
	const markInsideReactTree = useStableCallback(() => {
		dataRef.current.insideReactTree = true;
		clearInsideReactTreeTimeout.start(0, clearInsideReactTree);
	});
	import_react.useEffect(() => {
		if (!open || !enabled) return;
		dataRef.current.__escapeKeyBubbles = escapeKeyBubbles;
		dataRef.current.__outsidePressBubbles = outsidePressBubbles;
		const compositionTimeout = new Timeout();
		const preventedPressSuppressionTimeout = new Timeout();
		function handleCompositionStart() {
			compositionTimeout.clear();
			isComposingRef.current = true;
		}
		function handleCompositionEnd() {
			compositionTimeout.start(isWebKit() ? 5 : 0, () => {
				isComposingRef.current = false;
			});
		}
		function suppressImmediateOutsideClickAfterPreventedStart() {
			suppressNextOutsideClickRef.current = true;
			preventedPressSuppressionTimeout.start(0, () => {
				suppressNextOutsideClickRef.current = false;
			});
		}
		function resetPressStartState() {
			pressStartedInsideRef.current = false;
			pressStartPreventedRef.current = false;
		}
		function getOutsidePressEvent() {
			const type = currentPointerTypeRef.current;
			const computedType = type === "pen" || !type ? "mouse" : type;
			const outsidePressEventValue = getOutsidePressEventProp();
			const resolved = typeof outsidePressEventValue === "function" ? outsidePressEventValue() : outsidePressEventValue;
			if (typeof resolved === "string") return resolved;
			return resolved[computedType];
		}
		function shouldIgnoreEvent(event) {
			const computedOutsidePressEvent = getOutsidePressEvent();
			return computedOutsidePressEvent === "intentional" && event.type !== "click" || computedOutsidePressEvent === "sloppy" && event.type === "click";
		}
		function isEventWithinFloatingTree(event) {
			const nodeId = dataRef.current.floatingContext?.nodeId;
			const targetIsInsideChildren = tree && getNodeChildren(tree.nodesRef.current, nodeId).some((node) => isEventTargetWithin(event, node.context?.elements.floating));
			return isEventTargetWithin(event, store.select("floatingElement")) || isEventTargetWithin(event, store.select("domReferenceElement")) || targetIsInsideChildren;
		}
		function closeOnPressOutside(event) {
			if (shouldIgnoreEvent(event)) {
				clearInsideReactTree();
				return;
			}
			if (dataRef.current.insideReactTree) {
				clearInsideReactTree();
				return;
			}
			const target = getTarget(event);
			const inertSelector = `[${createAttribute("inert")}]`;
			let markers = Array.from(ownerDocument(store.select("floatingElement")).querySelectorAll(inertSelector));
			const targetRoot = isElement(target) ? target.getRootNode() : null;
			if (isShadowRoot(targetRoot)) markers = markers.concat(Array.from(targetRoot.querySelectorAll(inertSelector)));
			const triggers = store.context.triggerElements;
			if (target && (triggers.hasElement(target) || triggers.hasMatchingElement((trigger) => contains(trigger, target)))) return;
			let targetRootAncestor = isElement(target) ? target : null;
			while (targetRootAncestor && !isLastTraversableNode(targetRootAncestor)) {
				const nextParent = getParentNode(targetRootAncestor);
				if (isLastTraversableNode(nextParent) || !isElement(nextParent)) break;
				targetRootAncestor = nextParent;
			}
			if (markers.length && isElement(target) && !isRootElement(target) && !contains(target, store.select("floatingElement")) && markers.every((marker) => !contains(targetRootAncestor, marker))) return;
			if (isHTMLElement(target) && !("touches" in event)) {
				const lastTraversableNode = isLastTraversableNode(target);
				const style = getComputedStyle$1(target);
				const scrollRe = /auto|scroll/;
				const isScrollableX = lastTraversableNode || scrollRe.test(style.overflowX);
				const isScrollableY = lastTraversableNode || scrollRe.test(style.overflowY);
				const canScrollX = isScrollableX && target.clientWidth > 0 && target.scrollWidth > target.clientWidth;
				const canScrollY = isScrollableY && target.clientHeight > 0 && target.scrollHeight > target.clientHeight;
				const isRTL = style.direction === "rtl";
				const pressedVerticalScrollbar = canScrollY && (isRTL ? event.offsetX <= target.offsetWidth - target.clientWidth : event.offsetX > target.clientWidth);
				const pressedHorizontalScrollbar = canScrollX && event.offsetY > target.clientHeight;
				if (pressedVerticalScrollbar || pressedHorizontalScrollbar) return;
			}
			if (isEventWithinFloatingTree(event)) return;
			if (getOutsidePressEvent() === "intentional" && suppressNextOutsideClickRef.current) {
				preventedPressSuppressionTimeout.clear();
				suppressNextOutsideClickRef.current = false;
				return;
			}
			if (typeof outsidePress$1 === "function" && !outsidePress$1(event)) return;
			const nodeId = dataRef.current.floatingContext?.nodeId;
			const children = tree ? getNodeChildren(tree.nodesRef.current, nodeId) : [];
			if (children.length > 0) {
				let shouldDismiss = true;
				children.forEach((child) => {
					if (child.context?.open && !child.context.dataRef.current.__outsidePressBubbles) shouldDismiss = false;
				});
				if (!shouldDismiss) return;
			}
			store.setOpen(false, createChangeEventDetails(outsidePress, event));
			clearInsideReactTree();
		}
		function handlePointerDown(event) {
			if (getOutsidePressEvent() !== "sloppy" || event.pointerType === "touch" || !store.select("open") || !enabled || isEventTargetWithin(event, store.select("floatingElement")) || isEventTargetWithin(event, store.select("domReferenceElement"))) return;
			closeOnPressOutside(event);
		}
		function handleTouchStart(event) {
			if (getOutsidePressEvent() !== "sloppy" || !store.select("open") || !enabled || isEventTargetWithin(event, store.select("floatingElement")) || isEventTargetWithin(event, store.select("domReferenceElement"))) return;
			const touch = event.touches[0];
			if (touch) {
				touchStateRef.current = {
					startTime: Date.now(),
					startX: touch.clientX,
					startY: touch.clientY,
					dismissOnTouchEnd: false,
					dismissOnMouseDown: true
				};
				cancelDismissOnEndTimeout.start(1e3, () => {
					if (touchStateRef.current) {
						touchStateRef.current.dismissOnTouchEnd = false;
						touchStateRef.current.dismissOnMouseDown = false;
					}
				});
			}
		}
		function handleTouchStartCapture(event) {
			currentPointerTypeRef.current = "touch";
			const target = getTarget(event);
			function callback() {
				handleTouchStart(event);
				target?.removeEventListener(event.type, callback);
			}
			target?.addEventListener(event.type, callback);
		}
		function closeOnPressOutsideCapture(event) {
			cancelDismissOnEndTimeout.clear();
			if (event.type === "pointerdown") currentPointerTypeRef.current = event.pointerType;
			if (event.type === "mousedown" && touchStateRef.current && !touchStateRef.current.dismissOnMouseDown) return;
			const target = getTarget(event);
			function callback() {
				if (event.type === "pointerdown") handlePointerDown(event);
				else closeOnPressOutside(event);
				target?.removeEventListener(event.type, callback);
			}
			target?.addEventListener(event.type, callback);
		}
		function handlePressEndCapture(event) {
			if (!pressStartedInsideRef.current) return;
			const pressStartedInsideDefaultPrevented = pressStartPreventedRef.current;
			resetPressStartState();
			if (getOutsidePressEvent() !== "intentional") return;
			if (event.type === "pointercancel") {
				if (pressStartedInsideDefaultPrevented) suppressImmediateOutsideClickAfterPreventedStart();
				return;
			}
			if (isEventWithinFloatingTree(event)) return;
			if (pressStartedInsideDefaultPrevented) {
				suppressImmediateOutsideClickAfterPreventedStart();
				return;
			}
			if (typeof outsidePress$1 === "function" && !outsidePress$1(event)) return;
			preventedPressSuppressionTimeout.clear();
			suppressNextOutsideClickRef.current = true;
			clearInsideReactTree();
		}
		function handleTouchMove(event) {
			if (getOutsidePressEvent() !== "sloppy" || !touchStateRef.current || isEventTargetWithin(event, store.select("floatingElement")) || isEventTargetWithin(event, store.select("domReferenceElement"))) return;
			const touch = event.touches[0];
			if (!touch) return;
			const deltaX = Math.abs(touch.clientX - touchStateRef.current.startX);
			const deltaY = Math.abs(touch.clientY - touchStateRef.current.startY);
			const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);
			if (distance > 5) touchStateRef.current.dismissOnTouchEnd = true;
			if (distance > 10) {
				closeOnPressOutside(event);
				cancelDismissOnEndTimeout.clear();
				touchStateRef.current = null;
			}
		}
		function handleTouchMoveCapture(event) {
			const target = getTarget(event);
			function callback() {
				handleTouchMove(event);
				target?.removeEventListener(event.type, callback);
			}
			target?.addEventListener(event.type, callback);
		}
		function handleTouchEnd(event) {
			if (getOutsidePressEvent() !== "sloppy" || !touchStateRef.current || isEventTargetWithin(event, store.select("floatingElement")) || isEventTargetWithin(event, store.select("domReferenceElement"))) return;
			if (touchStateRef.current.dismissOnTouchEnd) closeOnPressOutside(event);
			cancelDismissOnEndTimeout.clear();
			touchStateRef.current = null;
		}
		function handleTouchEndCapture(event) {
			const target = getTarget(event);
			function callback() {
				handleTouchEnd(event);
				target?.removeEventListener(event.type, callback);
			}
			target?.addEventListener(event.type, callback);
		}
		const doc = ownerDocument(floatingElement);
		if (escapeKey$1) {
			doc.addEventListener("keydown", closeOnEscapeKeyDown);
			doc.addEventListener("compositionstart", handleCompositionStart);
			doc.addEventListener("compositionend", handleCompositionEnd);
		}
		if (outsidePressEnabled) {
			doc.addEventListener("click", closeOnPressOutsideCapture, true);
			doc.addEventListener("pointerdown", closeOnPressOutsideCapture, true);
			doc.addEventListener("pointerup", handlePressEndCapture, true);
			doc.addEventListener("pointercancel", handlePressEndCapture, true);
			doc.addEventListener("mousedown", closeOnPressOutsideCapture, true);
			doc.addEventListener("mouseup", handlePressEndCapture, true);
			doc.addEventListener("touchstart", handleTouchStartCapture, true);
			doc.addEventListener("touchmove", handleTouchMoveCapture, true);
			doc.addEventListener("touchend", handleTouchEndCapture, true);
		}
		return () => {
			if (escapeKey$1) {
				doc.removeEventListener("keydown", closeOnEscapeKeyDown);
				doc.removeEventListener("compositionstart", handleCompositionStart);
				doc.removeEventListener("compositionend", handleCompositionEnd);
			}
			if (outsidePressEnabled) {
				doc.removeEventListener("click", closeOnPressOutsideCapture, true);
				doc.removeEventListener("pointerdown", closeOnPressOutsideCapture, true);
				doc.removeEventListener("pointerup", handlePressEndCapture, true);
				doc.removeEventListener("pointercancel", handlePressEndCapture, true);
				doc.removeEventListener("mousedown", closeOnPressOutsideCapture, true);
				doc.removeEventListener("mouseup", handlePressEndCapture, true);
				doc.removeEventListener("touchstart", handleTouchStartCapture, true);
				doc.removeEventListener("touchmove", handleTouchMoveCapture, true);
				doc.removeEventListener("touchend", handleTouchEndCapture, true);
			}
			compositionTimeout.clear();
			preventedPressSuppressionTimeout.clear();
			resetPressStartState();
			suppressNextOutsideClickRef.current = false;
		};
	}, [
		dataRef,
		floatingElement,
		escapeKey$1,
		outsidePressEnabled,
		outsidePress$1,
		open,
		enabled,
		escapeKeyBubbles,
		outsidePressBubbles,
		closeOnEscapeKeyDown,
		clearInsideReactTree,
		getOutsidePressEventProp,
		tree,
		store,
		cancelDismissOnEndTimeout
	]);
	import_react.useEffect(clearInsideReactTree, [outsidePress$1, clearInsideReactTree]);
	const reference = import_react.useMemo(() => ({
		onKeyDown: closeOnEscapeKeyDown,
		[bubbleHandlerKeys[referencePressEvent]]: (event) => {
			if (!isReferencePressEnabled()) return;
			store.setOpen(false, createChangeEventDetails(triggerPress, event.nativeEvent));
		},
		...referencePressEvent !== "intentional" && { onClick(event) {
			if (!isReferencePressEnabled()) return;
			store.setOpen(false, createChangeEventDetails("trigger-press", event.nativeEvent));
		} }
	}), [
		closeOnEscapeKeyDown,
		store,
		referencePressEvent,
		isReferencePressEnabled
	]);
	const markPressStartedInsideReactTree = useStableCallback((event) => {
		if (!open || !enabled || event.button !== 0) return;
		const target = getTarget(event.nativeEvent);
		if (!contains(store.select("floatingElement"), target)) return;
		if (!pressStartedInsideRef.current) {
			pressStartedInsideRef.current = true;
			pressStartPreventedRef.current = false;
		}
	});
	const markInsidePressStartPrevented = useStableCallback((event) => {
		if (!open || !enabled) return;
		if (!(event.defaultPrevented || event.nativeEvent.defaultPrevented)) return;
		if (pressStartedInsideRef.current) pressStartPreventedRef.current = true;
	});
	const floating = import_react.useMemo(() => ({
		onKeyDown: closeOnEscapeKeyDown,
		onPointerDown: markInsidePressStartPrevented,
		onMouseDown: markInsidePressStartPrevented,
		onClickCapture: markInsideReactTree,
		onMouseDownCapture(event) {
			markInsideReactTree();
			markPressStartedInsideReactTree(event);
		},
		onPointerDownCapture(event) {
			markInsideReactTree();
			markPressStartedInsideReactTree(event);
		},
		onMouseUpCapture: markInsideReactTree,
		onTouchEndCapture: markInsideReactTree,
		onTouchMoveCapture: markInsideReactTree
	}), [
		closeOnEscapeKeyDown,
		markInsideReactTree,
		markPressStartedInsideReactTree,
		markInsidePressStartPrevented
	]);
	return import_react.useMemo(() => enabled ? {
		reference,
		floating,
		trigger: reference
	} : {}, [
		enabled,
		reference,
		floating
	]);
}
//#endregion
//#region node_modules/.pnpm/@floating-ui+core@1.7.5/node_modules/@floating-ui/core/dist/floating-ui.core.mjs
function computeCoordsFromPlacement(_ref, placement, rtl) {
	let { reference, floating } = _ref;
	const sideAxis = getSideAxis(placement);
	const alignmentAxis = getAlignmentAxis(placement);
	const alignLength = getAxisLength(alignmentAxis);
	const side = getSide(placement);
	const isVertical = sideAxis === "y";
	const commonX = reference.x + reference.width / 2 - floating.width / 2;
	const commonY = reference.y + reference.height / 2 - floating.height / 2;
	const commonAlign = reference[alignLength] / 2 - floating[alignLength] / 2;
	let coords;
	switch (side) {
		case "top":
			coords = {
				x: commonX,
				y: reference.y - floating.height
			};
			break;
		case "bottom":
			coords = {
				x: commonX,
				y: reference.y + reference.height
			};
			break;
		case "right":
			coords = {
				x: reference.x + reference.width,
				y: commonY
			};
			break;
		case "left":
			coords = {
				x: reference.x - floating.width,
				y: commonY
			};
			break;
		default: coords = {
			x: reference.x,
			y: reference.y
		};
	}
	switch (getAlignment(placement)) {
		case "start":
			coords[alignmentAxis] -= commonAlign * (rtl && isVertical ? -1 : 1);
			break;
		case "end":
			coords[alignmentAxis] += commonAlign * (rtl && isVertical ? -1 : 1);
			break;
	}
	return coords;
}
/**
* Resolves with an object of overflow side offsets that determine how much the
* element is overflowing a given clipping boundary on each side.
* - positive = overflowing the boundary by that number of pixels
* - negative = how many pixels left before it will overflow
* - 0 = lies flush with the boundary
* @see https://floating-ui.com/docs/detectOverflow
*/
async function detectOverflow(state, options) {
	var _await$platform$isEle;
	if (options === void 0) options = {};
	const { x, y, platform, rects, elements, strategy } = state;
	const { boundary = "clippingAncestors", rootBoundary = "viewport", elementContext = "floating", altBoundary = false, padding = 0 } = evaluate(options, state);
	const paddingObject = getPaddingObject(padding);
	const element = elements[altBoundary ? elementContext === "floating" ? "reference" : "floating" : elementContext];
	const clippingClientRect = rectToClientRect(await platform.getClippingRect({
		element: ((_await$platform$isEle = await (platform.isElement == null ? void 0 : platform.isElement(element))) != null ? _await$platform$isEle : true) ? element : element.contextElement || await (platform.getDocumentElement == null ? void 0 : platform.getDocumentElement(elements.floating)),
		boundary,
		rootBoundary,
		strategy
	}));
	const rect = elementContext === "floating" ? {
		x,
		y,
		width: rects.floating.width,
		height: rects.floating.height
	} : rects.reference;
	const offsetParent = await (platform.getOffsetParent == null ? void 0 : platform.getOffsetParent(elements.floating));
	const offsetScale = await (platform.isElement == null ? void 0 : platform.isElement(offsetParent)) ? await (platform.getScale == null ? void 0 : platform.getScale(offsetParent)) || {
		x: 1,
		y: 1
	} : {
		x: 1,
		y: 1
	};
	const elementClientRect = rectToClientRect(platform.convertOffsetParentRelativeRectToViewportRelativeRect ? await platform.convertOffsetParentRelativeRectToViewportRelativeRect({
		elements,
		rect,
		offsetParent,
		strategy
	}) : rect);
	return {
		top: (clippingClientRect.top - elementClientRect.top + paddingObject.top) / offsetScale.y,
		bottom: (elementClientRect.bottom - clippingClientRect.bottom + paddingObject.bottom) / offsetScale.y,
		left: (clippingClientRect.left - elementClientRect.left + paddingObject.left) / offsetScale.x,
		right: (elementClientRect.right - clippingClientRect.right + paddingObject.right) / offsetScale.x
	};
}
var MAX_RESET_COUNT = 50;
/**
* Computes the `x` and `y` coordinates that will place the floating element
* next to a given reference element.
*
* This export does not have any `platform` interface logic. You will need to
* write one for the platform you are using Floating UI with.
*/
var computePosition$1 = async (reference, floating, config) => {
	const { placement = "bottom", strategy = "absolute", middleware = [], platform } = config;
	const platformWithDetectOverflow = platform.detectOverflow ? platform : {
		...platform,
		detectOverflow
	};
	const rtl = await (platform.isRTL == null ? void 0 : platform.isRTL(floating));
	let rects = await platform.getElementRects({
		reference,
		floating,
		strategy
	});
	let { x, y } = computeCoordsFromPlacement(rects, placement, rtl);
	let statefulPlacement = placement;
	let resetCount = 0;
	const middlewareData = {};
	for (let i = 0; i < middleware.length; i++) {
		const currentMiddleware = middleware[i];
		if (!currentMiddleware) continue;
		const { name, fn } = currentMiddleware;
		const { x: nextX, y: nextY, data, reset } = await fn({
			x,
			y,
			initialPlacement: placement,
			placement: statefulPlacement,
			strategy,
			middlewareData,
			rects,
			platform: platformWithDetectOverflow,
			elements: {
				reference,
				floating
			}
		});
		x = nextX != null ? nextX : x;
		y = nextY != null ? nextY : y;
		middlewareData[name] = {
			...middlewareData[name],
			...data
		};
		if (reset && resetCount < MAX_RESET_COUNT) {
			resetCount++;
			if (typeof reset === "object") {
				if (reset.placement) statefulPlacement = reset.placement;
				if (reset.rects) rects = reset.rects === true ? await platform.getElementRects({
					reference,
					floating,
					strategy
				}) : reset.rects;
				({x, y} = computeCoordsFromPlacement(rects, statefulPlacement, rtl));
			}
			i = -1;
		}
	}
	return {
		x,
		y,
		placement: statefulPlacement,
		strategy,
		middlewareData
	};
};
/**
* Optimizes the visibility of the floating element by flipping the `placement`
* in order to keep it in view when the preferred placement(s) will overflow the
* clipping boundary. Alternative to `autoPlacement`.
* @see https://floating-ui.com/docs/flip
*/
var flip$2 = function(options) {
	if (options === void 0) options = {};
	return {
		name: "flip",
		options,
		async fn(state) {
			var _middlewareData$arrow, _middlewareData$flip;
			const { placement, middlewareData, rects, initialPlacement, platform, elements } = state;
			const { mainAxis: checkMainAxis = true, crossAxis: checkCrossAxis = true, fallbackPlacements: specifiedFallbackPlacements, fallbackStrategy = "bestFit", fallbackAxisSideDirection = "none", flipAlignment = true, ...detectOverflowOptions } = evaluate(options, state);
			if ((_middlewareData$arrow = middlewareData.arrow) != null && _middlewareData$arrow.alignmentOffset) return {};
			const side = getSide(placement);
			const initialSideAxis = getSideAxis(initialPlacement);
			const isBasePlacement = getSide(initialPlacement) === initialPlacement;
			const rtl = await (platform.isRTL == null ? void 0 : platform.isRTL(elements.floating));
			const fallbackPlacements = specifiedFallbackPlacements || (isBasePlacement || !flipAlignment ? [getOppositePlacement(initialPlacement)] : getExpandedPlacements(initialPlacement));
			const hasFallbackAxisSideDirection = fallbackAxisSideDirection !== "none";
			if (!specifiedFallbackPlacements && hasFallbackAxisSideDirection) fallbackPlacements.push(...getOppositeAxisPlacements(initialPlacement, flipAlignment, fallbackAxisSideDirection, rtl));
			const placements = [initialPlacement, ...fallbackPlacements];
			const overflow = await platform.detectOverflow(state, detectOverflowOptions);
			const overflows = [];
			let overflowsData = ((_middlewareData$flip = middlewareData.flip) == null ? void 0 : _middlewareData$flip.overflows) || [];
			if (checkMainAxis) overflows.push(overflow[side]);
			if (checkCrossAxis) {
				const sides = getAlignmentSides(placement, rects, rtl);
				overflows.push(overflow[sides[0]], overflow[sides[1]]);
			}
			overflowsData = [...overflowsData, {
				placement,
				overflows
			}];
			if (!overflows.every((side) => side <= 0)) {
				var _middlewareData$flip2, _overflowsData$filter;
				const nextIndex = (((_middlewareData$flip2 = middlewareData.flip) == null ? void 0 : _middlewareData$flip2.index) || 0) + 1;
				const nextPlacement = placements[nextIndex];
				if (nextPlacement) {
					if (!(checkCrossAxis === "alignment" ? initialSideAxis !== getSideAxis(nextPlacement) : false) || overflowsData.every((d) => getSideAxis(d.placement) === initialSideAxis ? d.overflows[0] > 0 : true)) return {
						data: {
							index: nextIndex,
							overflows: overflowsData
						},
						reset: { placement: nextPlacement }
					};
				}
				let resetPlacement = (_overflowsData$filter = overflowsData.filter((d) => d.overflows[0] <= 0).sort((a, b) => a.overflows[1] - b.overflows[1])[0]) == null ? void 0 : _overflowsData$filter.placement;
				if (!resetPlacement) switch (fallbackStrategy) {
					case "bestFit": {
						var _overflowsData$filter2;
						const placement = (_overflowsData$filter2 = overflowsData.filter((d) => {
							if (hasFallbackAxisSideDirection) {
								const currentSideAxis = getSideAxis(d.placement);
								return currentSideAxis === initialSideAxis || currentSideAxis === "y";
							}
							return true;
						}).map((d) => [d.placement, d.overflows.filter((overflow) => overflow > 0).reduce((acc, overflow) => acc + overflow, 0)]).sort((a, b) => a[1] - b[1])[0]) == null ? void 0 : _overflowsData$filter2[0];
						if (placement) resetPlacement = placement;
						break;
					}
					case "initialPlacement":
						resetPlacement = initialPlacement;
						break;
				}
				if (placement !== resetPlacement) return { reset: { placement: resetPlacement } };
			}
			return {};
		}
	};
};
function getSideOffsets(overflow, rect) {
	return {
		top: overflow.top - rect.height,
		right: overflow.right - rect.width,
		bottom: overflow.bottom - rect.height,
		left: overflow.left - rect.width
	};
}
function isAnySideFullyClipped(overflow) {
	return sides.some((side) => overflow[side] >= 0);
}
/**
* Provides data to hide the floating element in applicable situations, such as
* when it is not in the same clipping context as the reference element.
* @see https://floating-ui.com/docs/hide
*/
var hide$3 = function(options) {
	if (options === void 0) options = {};
	return {
		name: "hide",
		options,
		async fn(state) {
			const { rects, platform } = state;
			const { strategy = "referenceHidden", ...detectOverflowOptions } = evaluate(options, state);
			switch (strategy) {
				case "referenceHidden": {
					const offsets = getSideOffsets(await platform.detectOverflow(state, {
						...detectOverflowOptions,
						elementContext: "reference"
					}), rects.reference);
					return { data: {
						referenceHiddenOffsets: offsets,
						referenceHidden: isAnySideFullyClipped(offsets)
					} };
				}
				case "escaped": {
					const offsets = getSideOffsets(await platform.detectOverflow(state, {
						...detectOverflowOptions,
						altBoundary: true
					}), rects.floating);
					return { data: {
						escapedOffsets: offsets,
						escaped: isAnySideFullyClipped(offsets)
					} };
				}
				default: return {};
			}
		}
	};
};
var originSides = /* @__PURE__ */ new Set(["left", "top"]);
async function convertValueToCoords(state, options) {
	const { placement, platform, elements } = state;
	const rtl = await (platform.isRTL == null ? void 0 : platform.isRTL(elements.floating));
	const side = getSide(placement);
	const alignment = getAlignment(placement);
	const isVertical = getSideAxis(placement) === "y";
	const mainAxisMulti = originSides.has(side) ? -1 : 1;
	const crossAxisMulti = rtl && isVertical ? -1 : 1;
	const rawValue = evaluate(options, state);
	let { mainAxis, crossAxis, alignmentAxis } = typeof rawValue === "number" ? {
		mainAxis: rawValue,
		crossAxis: 0,
		alignmentAxis: null
	} : {
		mainAxis: rawValue.mainAxis || 0,
		crossAxis: rawValue.crossAxis || 0,
		alignmentAxis: rawValue.alignmentAxis
	};
	if (alignment && typeof alignmentAxis === "number") crossAxis = alignment === "end" ? alignmentAxis * -1 : alignmentAxis;
	return isVertical ? {
		x: crossAxis * crossAxisMulti,
		y: mainAxis * mainAxisMulti
	} : {
		x: mainAxis * mainAxisMulti,
		y: crossAxis * crossAxisMulti
	};
}
/**
* Modifies the placement by translating the floating element along the
* specified axes.
* A number (shorthand for `mainAxis` or distance), or an axes configuration
* object may be passed.
* @see https://floating-ui.com/docs/offset
*/
var offset$2 = function(options) {
	if (options === void 0) options = 0;
	return {
		name: "offset",
		options,
		async fn(state) {
			var _middlewareData$offse, _middlewareData$arrow;
			const { x, y, placement, middlewareData } = state;
			const diffCoords = await convertValueToCoords(state, options);
			if (placement === ((_middlewareData$offse = middlewareData.offset) == null ? void 0 : _middlewareData$offse.placement) && (_middlewareData$arrow = middlewareData.arrow) != null && _middlewareData$arrow.alignmentOffset) return {};
			return {
				x: x + diffCoords.x,
				y: y + diffCoords.y,
				data: {
					...diffCoords,
					placement
				}
			};
		}
	};
};
/**
* Optimizes the visibility of the floating element by shifting it in order to
* keep it in view when it will overflow the clipping boundary.
* @see https://floating-ui.com/docs/shift
*/
var shift$2 = function(options) {
	if (options === void 0) options = {};
	return {
		name: "shift",
		options,
		async fn(state) {
			const { x, y, placement, platform } = state;
			const { mainAxis: checkMainAxis = true, crossAxis: checkCrossAxis = false, limiter = { fn: (_ref) => {
				let { x, y } = _ref;
				return {
					x,
					y
				};
			} }, ...detectOverflowOptions } = evaluate(options, state);
			const coords = {
				x,
				y
			};
			const overflow = await platform.detectOverflow(state, detectOverflowOptions);
			const crossAxis = getSideAxis(getSide(placement));
			const mainAxis = getOppositeAxis(crossAxis);
			let mainAxisCoord = coords[mainAxis];
			let crossAxisCoord = coords[crossAxis];
			if (checkMainAxis) {
				const minSide = mainAxis === "y" ? "top" : "left";
				const maxSide = mainAxis === "y" ? "bottom" : "right";
				const min = mainAxisCoord + overflow[minSide];
				const max = mainAxisCoord - overflow[maxSide];
				mainAxisCoord = clamp$1(min, mainAxisCoord, max);
			}
			if (checkCrossAxis) {
				const minSide = crossAxis === "y" ? "top" : "left";
				const maxSide = crossAxis === "y" ? "bottom" : "right";
				const min = crossAxisCoord + overflow[minSide];
				const max = crossAxisCoord - overflow[maxSide];
				crossAxisCoord = clamp$1(min, crossAxisCoord, max);
			}
			const limitedCoords = limiter.fn({
				...state,
				[mainAxis]: mainAxisCoord,
				[crossAxis]: crossAxisCoord
			});
			return {
				...limitedCoords,
				data: {
					x: limitedCoords.x - x,
					y: limitedCoords.y - y,
					enabled: {
						[mainAxis]: checkMainAxis,
						[crossAxis]: checkCrossAxis
					}
				}
			};
		}
	};
};
/**
* Built-in `limiter` that will stop `shift()` at a certain point.
*/
var limitShift$2 = function(options) {
	if (options === void 0) options = {};
	return {
		options,
		fn(state) {
			const { x, y, placement, rects, middlewareData } = state;
			const { offset = 0, mainAxis: checkMainAxis = true, crossAxis: checkCrossAxis = true } = evaluate(options, state);
			const coords = {
				x,
				y
			};
			const crossAxis = getSideAxis(placement);
			const mainAxis = getOppositeAxis(crossAxis);
			let mainAxisCoord = coords[mainAxis];
			let crossAxisCoord = coords[crossAxis];
			const rawOffset = evaluate(offset, state);
			const computedOffset = typeof rawOffset === "number" ? {
				mainAxis: rawOffset,
				crossAxis: 0
			} : {
				mainAxis: 0,
				crossAxis: 0,
				...rawOffset
			};
			if (checkMainAxis) {
				const len = mainAxis === "y" ? "height" : "width";
				const limitMin = rects.reference[mainAxis] - rects.floating[len] + computedOffset.mainAxis;
				const limitMax = rects.reference[mainAxis] + rects.reference[len] - computedOffset.mainAxis;
				if (mainAxisCoord < limitMin) mainAxisCoord = limitMin;
				else if (mainAxisCoord > limitMax) mainAxisCoord = limitMax;
			}
			if (checkCrossAxis) {
				var _middlewareData$offse, _middlewareData$offse2;
				const len = mainAxis === "y" ? "width" : "height";
				const isOriginSide = originSides.has(getSide(placement));
				const limitMin = rects.reference[crossAxis] - rects.floating[len] + (isOriginSide ? ((_middlewareData$offse = middlewareData.offset) == null ? void 0 : _middlewareData$offse[crossAxis]) || 0 : 0) + (isOriginSide ? 0 : computedOffset.crossAxis);
				const limitMax = rects.reference[crossAxis] + rects.reference[len] + (isOriginSide ? 0 : ((_middlewareData$offse2 = middlewareData.offset) == null ? void 0 : _middlewareData$offse2[crossAxis]) || 0) - (isOriginSide ? computedOffset.crossAxis : 0);
				if (crossAxisCoord < limitMin) crossAxisCoord = limitMin;
				else if (crossAxisCoord > limitMax) crossAxisCoord = limitMax;
			}
			return {
				[mainAxis]: mainAxisCoord,
				[crossAxis]: crossAxisCoord
			};
		}
	};
};
/**
* Provides data that allows you to change the size of the floating element —
* for instance, prevent it from overflowing the clipping boundary or match the
* width of the reference element.
* @see https://floating-ui.com/docs/size
*/
var size$2 = function(options) {
	if (options === void 0) options = {};
	return {
		name: "size",
		options,
		async fn(state) {
			var _state$middlewareData, _state$middlewareData2;
			const { placement, rects, platform, elements } = state;
			const { apply = () => {}, ...detectOverflowOptions } = evaluate(options, state);
			const overflow = await platform.detectOverflow(state, detectOverflowOptions);
			const side = getSide(placement);
			const alignment = getAlignment(placement);
			const isYAxis = getSideAxis(placement) === "y";
			const { width, height } = rects.floating;
			let heightSide;
			let widthSide;
			if (side === "top" || side === "bottom") {
				heightSide = side;
				widthSide = alignment === (await (platform.isRTL == null ? void 0 : platform.isRTL(elements.floating)) ? "start" : "end") ? "left" : "right";
			} else {
				widthSide = side;
				heightSide = alignment === "end" ? "top" : "bottom";
			}
			const maximumClippingHeight = height - overflow.top - overflow.bottom;
			const maximumClippingWidth = width - overflow.left - overflow.right;
			const overflowAvailableHeight = min(height - overflow[heightSide], maximumClippingHeight);
			const overflowAvailableWidth = min(width - overflow[widthSide], maximumClippingWidth);
			const noShift = !state.middlewareData.shift;
			let availableHeight = overflowAvailableHeight;
			let availableWidth = overflowAvailableWidth;
			if ((_state$middlewareData = state.middlewareData.shift) != null && _state$middlewareData.enabled.x) availableWidth = maximumClippingWidth;
			if ((_state$middlewareData2 = state.middlewareData.shift) != null && _state$middlewareData2.enabled.y) availableHeight = maximumClippingHeight;
			if (noShift && !alignment) {
				const xMin = max(overflow.left, 0);
				const xMax = max(overflow.right, 0);
				const yMin = max(overflow.top, 0);
				const yMax = max(overflow.bottom, 0);
				if (isYAxis) availableWidth = width - 2 * (xMin !== 0 || xMax !== 0 ? xMin + xMax : max(overflow.left, overflow.right));
				else availableHeight = height - 2 * (yMin !== 0 || yMax !== 0 ? yMin + yMax : max(overflow.top, overflow.bottom));
			}
			await apply({
				...state,
				availableWidth,
				availableHeight
			});
			const nextDimensions = await platform.getDimensions(elements.floating);
			if (width !== nextDimensions.width || height !== nextDimensions.height) return { reset: { rects: true } };
			return {};
		}
	};
};
//#endregion
//#region node_modules/.pnpm/@floating-ui+dom@1.7.6/node_modules/@floating-ui/dom/dist/floating-ui.dom.mjs
function getCssDimensions(element) {
	const css = getComputedStyle$1(element);
	let width = parseFloat(css.width) || 0;
	let height = parseFloat(css.height) || 0;
	const hasOffset = isHTMLElement(element);
	const offsetWidth = hasOffset ? element.offsetWidth : width;
	const offsetHeight = hasOffset ? element.offsetHeight : height;
	const shouldFallback = round(width) !== offsetWidth || round(height) !== offsetHeight;
	if (shouldFallback) {
		width = offsetWidth;
		height = offsetHeight;
	}
	return {
		width,
		height,
		$: shouldFallback
	};
}
function unwrapElement(element) {
	return !isElement(element) ? element.contextElement : element;
}
function getScale$1(element) {
	const domElement = unwrapElement(element);
	if (!isHTMLElement(domElement)) return createCoords(1);
	const rect = domElement.getBoundingClientRect();
	const { width, height, $ } = getCssDimensions(domElement);
	let x = ($ ? round(rect.width) : rect.width) / width;
	let y = ($ ? round(rect.height) : rect.height) / height;
	if (!x || !Number.isFinite(x)) x = 1;
	if (!y || !Number.isFinite(y)) y = 1;
	return {
		x,
		y
	};
}
var noOffsets = /* @__PURE__ */ createCoords(0);
function getVisualOffsets(element) {
	const win = getWindow(element);
	if (!isWebKit() || !win.visualViewport) return noOffsets;
	return {
		x: win.visualViewport.offsetLeft,
		y: win.visualViewport.offsetTop
	};
}
function shouldAddVisualOffsets(element, isFixed, floatingOffsetParent) {
	if (isFixed === void 0) isFixed = false;
	if (!floatingOffsetParent || isFixed && floatingOffsetParent !== getWindow(element)) return false;
	return isFixed;
}
function getBoundingClientRect(element, includeScale, isFixedStrategy, offsetParent) {
	if (includeScale === void 0) includeScale = false;
	if (isFixedStrategy === void 0) isFixedStrategy = false;
	const clientRect = element.getBoundingClientRect();
	const domElement = unwrapElement(element);
	let scale = createCoords(1);
	if (includeScale) if (offsetParent) {
		if (isElement(offsetParent)) scale = getScale$1(offsetParent);
	} else scale = getScale$1(element);
	const visualOffsets = shouldAddVisualOffsets(domElement, isFixedStrategy, offsetParent) ? getVisualOffsets(domElement) : createCoords(0);
	let x = (clientRect.left + visualOffsets.x) / scale.x;
	let y = (clientRect.top + visualOffsets.y) / scale.y;
	let width = clientRect.width / scale.x;
	let height = clientRect.height / scale.y;
	if (domElement) {
		const win = getWindow(domElement);
		const offsetWin = offsetParent && isElement(offsetParent) ? getWindow(offsetParent) : offsetParent;
		let currentWin = win;
		let currentIFrame = getFrameElement(currentWin);
		while (currentIFrame && offsetParent && offsetWin !== currentWin) {
			const iframeScale = getScale$1(currentIFrame);
			const iframeRect = currentIFrame.getBoundingClientRect();
			const css = getComputedStyle$1(currentIFrame);
			const left = iframeRect.left + (currentIFrame.clientLeft + parseFloat(css.paddingLeft)) * iframeScale.x;
			const top = iframeRect.top + (currentIFrame.clientTop + parseFloat(css.paddingTop)) * iframeScale.y;
			x *= iframeScale.x;
			y *= iframeScale.y;
			width *= iframeScale.x;
			height *= iframeScale.y;
			x += left;
			y += top;
			currentWin = getWindow(currentIFrame);
			currentIFrame = getFrameElement(currentWin);
		}
	}
	return rectToClientRect({
		width,
		height,
		x,
		y
	});
}
function getWindowScrollBarX(element, rect) {
	const leftScroll = getNodeScroll(element).scrollLeft;
	if (!rect) return getBoundingClientRect(getDocumentElement(element)).left + leftScroll;
	return rect.left + leftScroll;
}
function getHTMLOffset(documentElement, scroll) {
	const htmlRect = documentElement.getBoundingClientRect();
	return {
		x: htmlRect.left + scroll.scrollLeft - getWindowScrollBarX(documentElement, htmlRect),
		y: htmlRect.top + scroll.scrollTop
	};
}
function convertOffsetParentRelativeRectToViewportRelativeRect(_ref) {
	let { elements, rect, offsetParent, strategy } = _ref;
	const isFixed = strategy === "fixed";
	const documentElement = getDocumentElement(offsetParent);
	const topLayer = elements ? isTopLayer(elements.floating) : false;
	if (offsetParent === documentElement || topLayer && isFixed) return rect;
	let scroll = {
		scrollLeft: 0,
		scrollTop: 0
	};
	let scale = createCoords(1);
	const offsets = createCoords(0);
	const isOffsetParentAnElement = isHTMLElement(offsetParent);
	if (isOffsetParentAnElement || !isOffsetParentAnElement && !isFixed) {
		if (getNodeName(offsetParent) !== "body" || isOverflowElement(documentElement)) scroll = getNodeScroll(offsetParent);
		if (isOffsetParentAnElement) {
			const offsetRect = getBoundingClientRect(offsetParent);
			scale = getScale$1(offsetParent);
			offsets.x = offsetRect.x + offsetParent.clientLeft;
			offsets.y = offsetRect.y + offsetParent.clientTop;
		}
	}
	const htmlOffset = documentElement && !isOffsetParentAnElement && !isFixed ? getHTMLOffset(documentElement, scroll) : createCoords(0);
	return {
		width: rect.width * scale.x,
		height: rect.height * scale.y,
		x: rect.x * scale.x - scroll.scrollLeft * scale.x + offsets.x + htmlOffset.x,
		y: rect.y * scale.y - scroll.scrollTop * scale.y + offsets.y + htmlOffset.y
	};
}
function getClientRects(element) {
	return Array.from(element.getClientRects());
}
function getDocumentRect(element) {
	const html = getDocumentElement(element);
	const scroll = getNodeScroll(element);
	const body = element.ownerDocument.body;
	const width = max(html.scrollWidth, html.clientWidth, body.scrollWidth, body.clientWidth);
	const height = max(html.scrollHeight, html.clientHeight, body.scrollHeight, body.clientHeight);
	let x = -scroll.scrollLeft + getWindowScrollBarX(element);
	const y = -scroll.scrollTop;
	if (getComputedStyle$1(body).direction === "rtl") x += max(html.clientWidth, body.clientWidth) - width;
	return {
		width,
		height,
		x,
		y
	};
}
var SCROLLBAR_MAX = 25;
function getViewportRect(element, strategy) {
	const win = getWindow(element);
	const html = getDocumentElement(element);
	const visualViewport = win.visualViewport;
	let width = html.clientWidth;
	let height = html.clientHeight;
	let x = 0;
	let y = 0;
	if (visualViewport) {
		width = visualViewport.width;
		height = visualViewport.height;
		const visualViewportBased = isWebKit();
		if (!visualViewportBased || visualViewportBased && strategy === "fixed") {
			x = visualViewport.offsetLeft;
			y = visualViewport.offsetTop;
		}
	}
	const windowScrollbarX = getWindowScrollBarX(html);
	if (windowScrollbarX <= 0) {
		const doc = html.ownerDocument;
		const body = doc.body;
		const bodyStyles = getComputedStyle(body);
		const bodyMarginInline = doc.compatMode === "CSS1Compat" ? parseFloat(bodyStyles.marginLeft) + parseFloat(bodyStyles.marginRight) || 0 : 0;
		const clippingStableScrollbarWidth = Math.abs(html.clientWidth - body.clientWidth - bodyMarginInline);
		if (clippingStableScrollbarWidth <= SCROLLBAR_MAX) width -= clippingStableScrollbarWidth;
	} else if (windowScrollbarX <= SCROLLBAR_MAX) width += windowScrollbarX;
	return {
		width,
		height,
		x,
		y
	};
}
function getInnerBoundingClientRect(element, strategy) {
	const clientRect = getBoundingClientRect(element, true, strategy === "fixed");
	const top = clientRect.top + element.clientTop;
	const left = clientRect.left + element.clientLeft;
	const scale = isHTMLElement(element) ? getScale$1(element) : createCoords(1);
	return {
		width: element.clientWidth * scale.x,
		height: element.clientHeight * scale.y,
		x: left * scale.x,
		y: top * scale.y
	};
}
function getClientRectFromClippingAncestor(element, clippingAncestor, strategy) {
	let rect;
	if (clippingAncestor === "viewport") rect = getViewportRect(element, strategy);
	else if (clippingAncestor === "document") rect = getDocumentRect(getDocumentElement(element));
	else if (isElement(clippingAncestor)) rect = getInnerBoundingClientRect(clippingAncestor, strategy);
	else {
		const visualOffsets = getVisualOffsets(element);
		rect = {
			x: clippingAncestor.x - visualOffsets.x,
			y: clippingAncestor.y - visualOffsets.y,
			width: clippingAncestor.width,
			height: clippingAncestor.height
		};
	}
	return rectToClientRect(rect);
}
function hasFixedPositionAncestor(element, stopNode) {
	const parentNode = getParentNode(element);
	if (parentNode === stopNode || !isElement(parentNode) || isLastTraversableNode(parentNode)) return false;
	return getComputedStyle$1(parentNode).position === "fixed" || hasFixedPositionAncestor(parentNode, stopNode);
}
function getClippingElementAncestors(element, cache) {
	const cachedResult = cache.get(element);
	if (cachedResult) return cachedResult;
	let result = getOverflowAncestors(element, [], false).filter((el) => isElement(el) && getNodeName(el) !== "body");
	let currentContainingBlockComputedStyle = null;
	const elementIsFixed = getComputedStyle$1(element).position === "fixed";
	let currentNode = elementIsFixed ? getParentNode(element) : element;
	while (isElement(currentNode) && !isLastTraversableNode(currentNode)) {
		const computedStyle = getComputedStyle$1(currentNode);
		const currentNodeIsContaining = isContainingBlock(currentNode);
		if (!currentNodeIsContaining && computedStyle.position === "fixed") currentContainingBlockComputedStyle = null;
		if (elementIsFixed ? !currentNodeIsContaining && !currentContainingBlockComputedStyle : !currentNodeIsContaining && computedStyle.position === "static" && !!currentContainingBlockComputedStyle && (currentContainingBlockComputedStyle.position === "absolute" || currentContainingBlockComputedStyle.position === "fixed") || isOverflowElement(currentNode) && !currentNodeIsContaining && hasFixedPositionAncestor(element, currentNode)) result = result.filter((ancestor) => ancestor !== currentNode);
		else currentContainingBlockComputedStyle = computedStyle;
		currentNode = getParentNode(currentNode);
	}
	cache.set(element, result);
	return result;
}
function getClippingRect(_ref) {
	let { element, boundary, rootBoundary, strategy } = _ref;
	const clippingAncestors = [...boundary === "clippingAncestors" ? isTopLayer(element) ? [] : getClippingElementAncestors(element, this._c) : [].concat(boundary), rootBoundary];
	const firstRect = getClientRectFromClippingAncestor(element, clippingAncestors[0], strategy);
	let top = firstRect.top;
	let right = firstRect.right;
	let bottom = firstRect.bottom;
	let left = firstRect.left;
	for (let i = 1; i < clippingAncestors.length; i++) {
		const rect = getClientRectFromClippingAncestor(element, clippingAncestors[i], strategy);
		top = max(rect.top, top);
		right = min(rect.right, right);
		bottom = min(rect.bottom, bottom);
		left = max(rect.left, left);
	}
	return {
		width: right - left,
		height: bottom - top,
		x: left,
		y: top
	};
}
function getDimensions(element) {
	const { width, height } = getCssDimensions(element);
	return {
		width,
		height
	};
}
function getRectRelativeToOffsetParent(element, offsetParent, strategy) {
	const isOffsetParentAnElement = isHTMLElement(offsetParent);
	const documentElement = getDocumentElement(offsetParent);
	const isFixed = strategy === "fixed";
	const rect = getBoundingClientRect(element, true, isFixed, offsetParent);
	let scroll = {
		scrollLeft: 0,
		scrollTop: 0
	};
	const offsets = createCoords(0);
	function setLeftRTLScrollbarOffset() {
		offsets.x = getWindowScrollBarX(documentElement);
	}
	if (isOffsetParentAnElement || !isOffsetParentAnElement && !isFixed) {
		if (getNodeName(offsetParent) !== "body" || isOverflowElement(documentElement)) scroll = getNodeScroll(offsetParent);
		if (isOffsetParentAnElement) {
			const offsetRect = getBoundingClientRect(offsetParent, true, isFixed, offsetParent);
			offsets.x = offsetRect.x + offsetParent.clientLeft;
			offsets.y = offsetRect.y + offsetParent.clientTop;
		} else if (documentElement) setLeftRTLScrollbarOffset();
	}
	if (isFixed && !isOffsetParentAnElement && documentElement) setLeftRTLScrollbarOffset();
	const htmlOffset = documentElement && !isOffsetParentAnElement && !isFixed ? getHTMLOffset(documentElement, scroll) : createCoords(0);
	return {
		x: rect.left + scroll.scrollLeft - offsets.x - htmlOffset.x,
		y: rect.top + scroll.scrollTop - offsets.y - htmlOffset.y,
		width: rect.width,
		height: rect.height
	};
}
function isStaticPositioned(element) {
	return getComputedStyle$1(element).position === "static";
}
function getTrueOffsetParent(element, polyfill) {
	if (!isHTMLElement(element) || getComputedStyle$1(element).position === "fixed") return null;
	if (polyfill) return polyfill(element);
	let rawOffsetParent = element.offsetParent;
	if (getDocumentElement(element) === rawOffsetParent) rawOffsetParent = rawOffsetParent.ownerDocument.body;
	return rawOffsetParent;
}
function getOffsetParent(element, polyfill) {
	const win = getWindow(element);
	if (isTopLayer(element)) return win;
	if (!isHTMLElement(element)) {
		let svgOffsetParent = getParentNode(element);
		while (svgOffsetParent && !isLastTraversableNode(svgOffsetParent)) {
			if (isElement(svgOffsetParent) && !isStaticPositioned(svgOffsetParent)) return svgOffsetParent;
			svgOffsetParent = getParentNode(svgOffsetParent);
		}
		return win;
	}
	let offsetParent = getTrueOffsetParent(element, polyfill);
	while (offsetParent && isTableElement(offsetParent) && isStaticPositioned(offsetParent)) offsetParent = getTrueOffsetParent(offsetParent, polyfill);
	if (offsetParent && isLastTraversableNode(offsetParent) && isStaticPositioned(offsetParent) && !isContainingBlock(offsetParent)) return win;
	return offsetParent || getContainingBlock(element) || win;
}
var getElementRects = async function(data) {
	const getOffsetParentFn = this.getOffsetParent || getOffsetParent;
	const getDimensionsFn = this.getDimensions;
	const floatingDimensions = await getDimensionsFn(data.floating);
	return {
		reference: getRectRelativeToOffsetParent(data.reference, await getOffsetParentFn(data.floating), data.strategy),
		floating: {
			x: 0,
			y: 0,
			width: floatingDimensions.width,
			height: floatingDimensions.height
		}
	};
};
function isRTL(element) {
	return getComputedStyle$1(element).direction === "rtl";
}
var platform = {
	convertOffsetParentRelativeRectToViewportRelativeRect,
	getDocumentElement,
	getClippingRect,
	getOffsetParent,
	getElementRects,
	getClientRects,
	getDimensions,
	getScale: getScale$1,
	isElement,
	isRTL
};
function rectsAreEqual(a, b) {
	return a.x === b.x && a.y === b.y && a.width === b.width && a.height === b.height;
}
function observeMove(element, onMove) {
	let io = null;
	let timeoutId;
	const root = getDocumentElement(element);
	function cleanup() {
		var _io;
		clearTimeout(timeoutId);
		(_io = io) == null || _io.disconnect();
		io = null;
	}
	function refresh(skip, threshold) {
		if (skip === void 0) skip = false;
		if (threshold === void 0) threshold = 1;
		cleanup();
		const elementRectForRootMargin = element.getBoundingClientRect();
		const { left, top, width, height } = elementRectForRootMargin;
		if (!skip) onMove();
		if (!width || !height) return;
		const insetTop = floor(top);
		const insetRight = floor(root.clientWidth - (left + width));
		const insetBottom = floor(root.clientHeight - (top + height));
		const insetLeft = floor(left);
		const options = {
			rootMargin: -insetTop + "px " + -insetRight + "px " + -insetBottom + "px " + -insetLeft + "px",
			threshold: max(0, min(1, threshold)) || 1
		};
		let isFirstUpdate = true;
		function handleObserve(entries) {
			const ratio = entries[0].intersectionRatio;
			if (ratio !== threshold) {
				if (!isFirstUpdate) return refresh();
				if (!ratio) timeoutId = setTimeout(() => {
					refresh(false, 1e-7);
				}, 1e3);
				else refresh(false, ratio);
			}
			if (ratio === 1 && !rectsAreEqual(elementRectForRootMargin, element.getBoundingClientRect())) refresh();
			isFirstUpdate = false;
		}
		try {
			io = new IntersectionObserver(handleObserve, {
				...options,
				root: root.ownerDocument
			});
		} catch (_e) {
			io = new IntersectionObserver(handleObserve, options);
		}
		io.observe(element);
	}
	refresh(true);
	return cleanup;
}
/**
* Automatically updates the position of the floating element when necessary.
* Should only be called when the floating element is mounted on the DOM or
* visible on the screen.
* @returns cleanup function that should be invoked when the floating element is
* removed from the DOM or hidden from the screen.
* @see https://floating-ui.com/docs/autoUpdate
*/
function autoUpdate(reference, floating, update, options) {
	if (options === void 0) options = {};
	const { ancestorScroll = true, ancestorResize = true, elementResize = typeof ResizeObserver === "function", layoutShift = typeof IntersectionObserver === "function", animationFrame = false } = options;
	const referenceEl = unwrapElement(reference);
	const ancestors = ancestorScroll || ancestorResize ? [...referenceEl ? getOverflowAncestors(referenceEl) : [], ...floating ? getOverflowAncestors(floating) : []] : [];
	ancestors.forEach((ancestor) => {
		ancestorScroll && ancestor.addEventListener("scroll", update, { passive: true });
		ancestorResize && ancestor.addEventListener("resize", update);
	});
	const cleanupIo = referenceEl && layoutShift ? observeMove(referenceEl, update) : null;
	let reobserveFrame = -1;
	let resizeObserver = null;
	if (elementResize) {
		resizeObserver = new ResizeObserver((_ref) => {
			let [firstEntry] = _ref;
			if (firstEntry && firstEntry.target === referenceEl && resizeObserver && floating) {
				resizeObserver.unobserve(floating);
				cancelAnimationFrame(reobserveFrame);
				reobserveFrame = requestAnimationFrame(() => {
					var _resizeObserver;
					(_resizeObserver = resizeObserver) == null || _resizeObserver.observe(floating);
				});
			}
			update();
		});
		if (referenceEl && !animationFrame) resizeObserver.observe(referenceEl);
		if (floating) resizeObserver.observe(floating);
	}
	let frameId;
	let prevRefRect = animationFrame ? getBoundingClientRect(reference) : null;
	if (animationFrame) frameLoop();
	function frameLoop() {
		const nextRefRect = getBoundingClientRect(reference);
		if (prevRefRect && !rectsAreEqual(prevRefRect, nextRefRect)) update();
		prevRefRect = nextRefRect;
		frameId = requestAnimationFrame(frameLoop);
	}
	update();
	return () => {
		var _resizeObserver2;
		ancestors.forEach((ancestor) => {
			ancestorScroll && ancestor.removeEventListener("scroll", update);
			ancestorResize && ancestor.removeEventListener("resize", update);
		});
		cleanupIo?.();
		(_resizeObserver2 = resizeObserver) == null || _resizeObserver2.disconnect();
		resizeObserver = null;
		if (animationFrame) cancelAnimationFrame(frameId);
	};
}
/**
* Modifies the placement by translating the floating element along the
* specified axes.
* A number (shorthand for `mainAxis` or distance), or an axes configuration
* object may be passed.
* @see https://floating-ui.com/docs/offset
*/
var offset$1 = offset$2;
/**
* Optimizes the visibility of the floating element by shifting it in order to
* keep it in view when it will overflow the clipping boundary.
* @see https://floating-ui.com/docs/shift
*/
var shift$1 = shift$2;
/**
* Optimizes the visibility of the floating element by flipping the `placement`
* in order to keep it in view when the preferred placement(s) will overflow the
* clipping boundary. Alternative to `autoPlacement`.
* @see https://floating-ui.com/docs/flip
*/
var flip$1 = flip$2;
/**
* Provides data that allows you to change the size of the floating element —
* for instance, prevent it from overflowing the clipping boundary or match the
* width of the reference element.
* @see https://floating-ui.com/docs/size
*/
var size$1 = size$2;
/**
* Provides data to hide the floating element in applicable situations, such as
* when it is not in the same clipping context as the reference element.
* @see https://floating-ui.com/docs/hide
*/
var hide$2 = hide$3;
/**
* Built-in `limiter` that will stop `shift()` at a certain point.
*/
var limitShift$1 = limitShift$2;
/**
* Computes the `x` and `y` coordinates that will place the floating element
* next to a given reference element.
*/
var computePosition = (reference, floating, options) => {
	const cache = /* @__PURE__ */ new Map();
	const mergedOptions = {
		platform,
		...options
	};
	const platformWithCache = {
		...mergedOptions.platform,
		_c: cache
	};
	return computePosition$1(reference, floating, {
		...mergedOptions,
		platform: platformWithCache
	});
};
//#endregion
//#region node_modules/.pnpm/@floating-ui+react-dom@2.1.8_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@floating-ui/react-dom/dist/floating-ui.react-dom.mjs
var index = typeof document !== "undefined" ? import_react.useLayoutEffect : function noop() {};
function deepEqual(a, b) {
	if (a === b) return true;
	if (typeof a !== typeof b) return false;
	if (typeof a === "function" && a.toString() === b.toString()) return true;
	let length;
	let i;
	let keys;
	if (a && b && typeof a === "object") {
		if (Array.isArray(a)) {
			length = a.length;
			if (length !== b.length) return false;
			for (i = length; i-- !== 0;) if (!deepEqual(a[i], b[i])) return false;
			return true;
		}
		keys = Object.keys(a);
		length = keys.length;
		if (length !== Object.keys(b).length) return false;
		for (i = length; i-- !== 0;) if (!{}.hasOwnProperty.call(b, keys[i])) return false;
		for (i = length; i-- !== 0;) {
			const key = keys[i];
			if (key === "_owner" && a.$$typeof) continue;
			if (!deepEqual(a[key], b[key])) return false;
		}
		return true;
	}
	return a !== a && b !== b;
}
function getDPR(element) {
	if (typeof window === "undefined") return 1;
	return (element.ownerDocument.defaultView || window).devicePixelRatio || 1;
}
function roundByDPR(element, value) {
	const dpr = getDPR(element);
	return Math.round(value * dpr) / dpr;
}
function useLatestRef(value) {
	const ref = import_react.useRef(value);
	index(() => {
		ref.current = value;
	});
	return ref;
}
/**
* Provides data to position a floating element.
* @see https://floating-ui.com/docs/useFloating
*/
function useFloating$1(options) {
	if (options === void 0) options = {};
	const { placement = "bottom", strategy = "absolute", middleware = [], platform, elements: { reference: externalReference, floating: externalFloating } = {}, transform = true, whileElementsMounted, open } = options;
	const [data, setData] = import_react.useState({
		x: 0,
		y: 0,
		strategy,
		placement,
		middlewareData: {},
		isPositioned: false
	});
	const [latestMiddleware, setLatestMiddleware] = import_react.useState(middleware);
	if (!deepEqual(latestMiddleware, middleware)) setLatestMiddleware(middleware);
	const [_reference, _setReference] = import_react.useState(null);
	const [_floating, _setFloating] = import_react.useState(null);
	const setReference = import_react.useCallback((node) => {
		if (node !== referenceRef.current) {
			referenceRef.current = node;
			_setReference(node);
		}
	}, []);
	const setFloating = import_react.useCallback((node) => {
		if (node !== floatingRef.current) {
			floatingRef.current = node;
			_setFloating(node);
		}
	}, []);
	const referenceEl = externalReference || _reference;
	const floatingEl = externalFloating || _floating;
	const referenceRef = import_react.useRef(null);
	const floatingRef = import_react.useRef(null);
	const dataRef = import_react.useRef(data);
	const hasWhileElementsMounted = whileElementsMounted != null;
	const whileElementsMountedRef = useLatestRef(whileElementsMounted);
	const platformRef = useLatestRef(platform);
	const openRef = useLatestRef(open);
	const update = import_react.useCallback(() => {
		if (!referenceRef.current || !floatingRef.current) return;
		const config = {
			placement,
			strategy,
			middleware: latestMiddleware
		};
		if (platformRef.current) config.platform = platformRef.current;
		computePosition(referenceRef.current, floatingRef.current, config).then((data) => {
			const fullData = {
				...data,
				isPositioned: openRef.current !== false
			};
			if (isMountedRef.current && !deepEqual(dataRef.current, fullData)) {
				dataRef.current = fullData;
				import_react_dom.flushSync(() => {
					setData(fullData);
				});
			}
		});
	}, [
		latestMiddleware,
		placement,
		strategy,
		platformRef,
		openRef
	]);
	index(() => {
		if (open === false && dataRef.current.isPositioned) {
			dataRef.current.isPositioned = false;
			setData((data) => ({
				...data,
				isPositioned: false
			}));
		}
	}, [open]);
	const isMountedRef = import_react.useRef(false);
	index(() => {
		isMountedRef.current = true;
		return () => {
			isMountedRef.current = false;
		};
	}, []);
	index(() => {
		if (referenceEl) referenceRef.current = referenceEl;
		if (floatingEl) floatingRef.current = floatingEl;
		if (referenceEl && floatingEl) {
			if (whileElementsMountedRef.current) return whileElementsMountedRef.current(referenceEl, floatingEl, update);
			update();
		}
	}, [
		referenceEl,
		floatingEl,
		update,
		whileElementsMountedRef,
		hasWhileElementsMounted
	]);
	const refs = import_react.useMemo(() => ({
		reference: referenceRef,
		floating: floatingRef,
		setReference,
		setFloating
	}), [setReference, setFloating]);
	const elements = import_react.useMemo(() => ({
		reference: referenceEl,
		floating: floatingEl
	}), [referenceEl, floatingEl]);
	const floatingStyles = import_react.useMemo(() => {
		const initialStyles = {
			position: strategy,
			left: 0,
			top: 0
		};
		if (!elements.floating) return initialStyles;
		const x = roundByDPR(elements.floating, data.x);
		const y = roundByDPR(elements.floating, data.y);
		if (transform) return {
			...initialStyles,
			transform: "translate(" + x + "px, " + y + "px)",
			...getDPR(elements.floating) >= 1.5 && { willChange: "transform" }
		};
		return {
			position: strategy,
			left: x,
			top: y
		};
	}, [
		strategy,
		transform,
		elements.floating,
		data.x,
		data.y
	]);
	return import_react.useMemo(() => ({
		...data,
		update,
		refs,
		elements,
		floatingStyles
	}), [
		data,
		update,
		refs,
		elements,
		floatingStyles
	]);
}
/**
* Modifies the placement by translating the floating element along the
* specified axes.
* A number (shorthand for `mainAxis` or distance), or an axes configuration
* object may be passed.
* @see https://floating-ui.com/docs/offset
*/
var offset = (options, deps) => {
	const result = offset$1(options);
	return {
		name: result.name,
		fn: result.fn,
		options: [options, deps]
	};
};
/**
* Optimizes the visibility of the floating element by shifting it in order to
* keep it in view when it will overflow the clipping boundary.
* @see https://floating-ui.com/docs/shift
*/
var shift = (options, deps) => {
	const result = shift$1(options);
	return {
		name: result.name,
		fn: result.fn,
		options: [options, deps]
	};
};
/**
* Built-in `limiter` that will stop `shift()` at a certain point.
*/
var limitShift = (options, deps) => {
	return {
		fn: limitShift$1(options).fn,
		options: [options, deps]
	};
};
/**
* Optimizes the visibility of the floating element by flipping the `placement`
* in order to keep it in view when the preferred placement(s) will overflow the
* clipping boundary. Alternative to `autoPlacement`.
* @see https://floating-ui.com/docs/flip
*/
var flip = (options, deps) => {
	const result = flip$1(options);
	return {
		name: result.name,
		fn: result.fn,
		options: [options, deps]
	};
};
/**
* Provides data that allows you to change the size of the floating element —
* for instance, prevent it from overflowing the clipping boundary or match the
* width of the reference element.
* @see https://floating-ui.com/docs/size
*/
var size = (options, deps) => {
	const result = size$1(options);
	return {
		name: result.name,
		fn: result.fn,
		options: [options, deps]
	};
};
/**
* Provides data to hide the floating element in applicable situations, such as
* when it is not in the same clipping context as the reference element.
* @see https://floating-ui.com/docs/hide
*/
var hide$1 = (options, deps) => {
	const result = hide$2(options);
	return {
		name: result.name,
		fn: result.fn,
		options: [options, deps]
	};
};
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/components/FloatingRootStore.js
var selectors$1 = {
	open: createSelector((state) => state.open),
	domReferenceElement: createSelector((state) => state.domReferenceElement),
	referenceElement: createSelector((state) => state.positionReference ?? state.referenceElement),
	floatingElement: createSelector((state) => state.floatingElement),
	floatingId: createSelector((state) => state.floatingId)
};
var FloatingRootStore = class extends ReactStore {
	constructor(options) {
		const { nested, noEmit, onOpenChange, triggerElements, ...initialState } = options;
		super({
			...initialState,
			positionReference: initialState.referenceElement,
			domReferenceElement: initialState.referenceElement
		}, {
			onOpenChange,
			dataRef: { current: {} },
			events: createEventEmitter(),
			nested,
			noEmit,
			triggerElements
		}, selectors$1);
	}
	/**
	* Emits the `openchange` event through the internal event emitter and calls the `onOpenChange` handler with the provided arguments.
	*
	* @param newOpen The new open state.
	* @param eventDetails Details about the event that triggered the open state change.
	*/
	setOpen = (newOpen, eventDetails) => {
		if (!newOpen || !this.state.open || isClickLikeEvent(eventDetails.event)) this.context.dataRef.current.openEvent = newOpen ? eventDetails.event : void 0;
		if (!this.context.noEmit) {
			const details = {
				open: newOpen,
				reason: eventDetails.reason,
				nativeEvent: eventDetails.event,
				nested: this.context.nested,
				triggerElement: eventDetails.trigger
			};
			this.context.events.emit("openchange", details);
		}
		this.context.onOpenChange?.(newOpen, eventDetails);
	};
};
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/popups/popupTriggerMap.js
/**
* Data structure to keep track of popup trigger elements by their IDs.
* Uses both a set of Elements and a map of IDs to Elements for efficient lookups.
*/
var PopupTriggerMap = class {
	constructor() {
		this.elementsSet = /* @__PURE__ */ new Set();
		this.idMap = /* @__PURE__ */ new Map();
	}
	/**
	* Adds a trigger element with the given ID.
	*
	* Note: The provided element is assumed to not be registered under multiple IDs.
	*/
	add(id, element) {
		const existingElement = this.idMap.get(id);
		if (existingElement === element) return;
		if (existingElement !== void 0) this.elementsSet.delete(existingElement);
		this.elementsSet.add(element);
		this.idMap.set(id, element);
	}
	/**
	* Removes the trigger element with the given ID.
	*/
	delete(id) {
		const element = this.idMap.get(id);
		if (element) {
			this.elementsSet.delete(element);
			this.idMap.delete(id);
		}
	}
	/**
	* Whether the given element is registered as a trigger.
	*/
	hasElement(element) {
		return this.elementsSet.has(element);
	}
	/**
	* Whether there is a registered trigger element matching the given predicate.
	*/
	hasMatchingElement(predicate) {
		for (const element of this.elementsSet) if (predicate(element)) return true;
		return false;
	}
	/**
	* Returns the trigger element associated with the given ID, or undefined if no such element exists.
	*/
	getById(id) {
		return this.idMap.get(id);
	}
	/**
	* Returns an iterable of all registered trigger entries, where each entry is a tuple of [id, element].
	*/
	entries() {
		return this.idMap.entries();
	}
	/**
	* Returns an iterable of all registered trigger elements.
	*/
	elements() {
		return this.elementsSet.values();
	}
	/**
	* Returns the number of registered trigger elements.
	*/
	get size() {
		return this.idMap.size;
	}
};
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/hooks/useFloatingRootContext.js
function useFloatingRootContext(options) {
	const { open = false, onOpenChange, elements = {} } = options;
	const floatingId = useId();
	const nested = useFloatingParentNodeId() != null;
	const store = useRefWithInit(() => new FloatingRootStore({
		open,
		onOpenChange,
		referenceElement: elements.reference ?? null,
		floatingElement: elements.floating ?? null,
		triggerElements: new PopupTriggerMap(),
		floatingId,
		nested,
		noEmit: false
	})).current;
	useIsoLayoutEffect(() => {
		const valuesToSync = {
			open,
			floatingId
		};
		if (elements.reference !== void 0) {
			valuesToSync.referenceElement = elements.reference;
			valuesToSync.domReferenceElement = isElement(elements.reference) ? elements.reference : null;
		}
		if (elements.floating !== void 0) valuesToSync.floatingElement = elements.floating;
		store.update(valuesToSync);
	}, [
		open,
		floatingId,
		elements.reference,
		elements.floating,
		store
	]);
	store.context.onOpenChange = onOpenChange;
	store.context.nested = nested;
	store.context.noEmit = false;
	return store;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/hooks/useFloating.js
/**
* Provides data to position a floating element and context to add interactions.
* @see https://floating-ui.com/docs/useFloating
*/
function useFloating(options = {}) {
	const { nodeId, externalTree } = options;
	const internalRootStore = useFloatingRootContext(options);
	const rootContext = options.rootContext || internalRootStore;
	const rootContextElements = {
		reference: rootContext.useState("referenceElement"),
		floating: rootContext.useState("floatingElement"),
		domReference: rootContext.useState("domReferenceElement")
	};
	const [positionReference, setPositionReferenceRaw] = import_react.useState(null);
	const domReferenceRef = import_react.useRef(null);
	const tree = useFloatingTree(externalTree);
	useIsoLayoutEffect(() => {
		if (rootContextElements.domReference) domReferenceRef.current = rootContextElements.domReference;
	}, [rootContextElements.domReference]);
	const position = useFloating$1({
		...options,
		elements: {
			...rootContextElements,
			...positionReference && { reference: positionReference }
		}
	});
	const setPositionReference = import_react.useCallback((node) => {
		const computedPositionReference = isElement(node) ? {
			getBoundingClientRect: () => node.getBoundingClientRect(),
			getClientRects: () => node.getClientRects(),
			contextElement: node
		} : node;
		setPositionReferenceRaw(computedPositionReference);
		position.refs.setReference(computedPositionReference);
	}, [position.refs]);
	const [localDomReference, setLocalDomReference] = import_react.useState(null);
	const [localFloatingElement, setLocalFloatingElement] = import_react.useState(null);
	rootContext.useSyncedValue("referenceElement", localDomReference);
	rootContext.useSyncedValue("domReferenceElement", isElement(localDomReference) ? localDomReference : null);
	rootContext.useSyncedValue("floatingElement", localFloatingElement);
	const setReference = import_react.useCallback((node) => {
		if (isElement(node) || node === null) {
			domReferenceRef.current = node;
			setLocalDomReference(node);
		}
		if (isElement(position.refs.reference.current) || position.refs.reference.current === null || node !== null && !isElement(node)) position.refs.setReference(node);
	}, [position.refs, setLocalDomReference]);
	const setFloating = import_react.useCallback((node) => {
		setLocalFloatingElement(node);
		position.refs.setFloating(node);
	}, [position.refs]);
	const refs = import_react.useMemo(() => ({
		...position.refs,
		setReference,
		setFloating,
		setPositionReference,
		domReference: domReferenceRef
	}), [
		position.refs,
		setReference,
		setFloating,
		setPositionReference
	]);
	const elements = import_react.useMemo(() => ({
		...position.elements,
		domReference: rootContextElements.domReference
	}), [position.elements, rootContextElements.domReference]);
	const open = rootContext.useState("open");
	const floatingId = rootContext.useState("floatingId");
	const context = import_react.useMemo(() => ({
		...position,
		dataRef: rootContext.context.dataRef,
		open,
		onOpenChange: rootContext.setOpen,
		events: rootContext.context.events,
		floatingId,
		refs,
		elements,
		nodeId,
		rootStore: rootContext
	}), [
		position,
		refs,
		elements,
		nodeId,
		rootContext,
		open,
		floatingId
	]);
	useIsoLayoutEffect(() => {
		rootContext.context.dataRef.current.floatingContext = context;
		const node = tree?.nodesRef.current.find((n) => n.id === nodeId);
		if (node) node.context = context;
	});
	return import_react.useMemo(() => ({
		...position,
		context,
		refs,
		elements,
		rootStore: rootContext
	}), [
		position,
		refs,
		elements,
		context,
		rootContext
	]);
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/hooks/useInteractions.js
/**
* Merges an array of interaction hooks' props into prop getters, allowing
* event handler functions to be composed together without overwriting one
* another.
* @see https://floating-ui.com/docs/useInteractions
*/
function useInteractions(propsList = []) {
	const referenceDeps = propsList.map((key) => key?.reference);
	const floatingDeps = propsList.map((key) => key?.floating);
	const itemDeps = propsList.map((key) => key?.item);
	const triggerDeps = propsList.map((key) => key?.trigger);
	const getReferenceProps = import_react.useCallback((userProps) => mergeProps(userProps, propsList, "reference"), referenceDeps);
	const getFloatingProps = import_react.useCallback((userProps) => mergeProps(userProps, propsList, "floating"), floatingDeps);
	const getItemProps = import_react.useCallback((userProps) => mergeProps(userProps, propsList, "item"), itemDeps);
	const getTriggerProps = import_react.useCallback((userProps) => mergeProps(userProps, propsList, "trigger"), triggerDeps);
	return import_react.useMemo(() => ({
		getReferenceProps,
		getFloatingProps,
		getItemProps,
		getTriggerProps
	}), [
		getReferenceProps,
		getFloatingProps,
		getItemProps,
		getTriggerProps
	]);
}
function mergeProps(userProps, propsList, elementKey) {
	const eventHandlers = /* @__PURE__ */ new Map();
	const isItem = elementKey === "item";
	const outputProps = {};
	if (elementKey === "floating") {
		outputProps.tabIndex = -1;
		outputProps[FOCUSABLE_ATTRIBUTE] = "";
	}
	for (const key in userProps) {
		if (isItem && userProps) {
			if (key === "active" || key === "selected") continue;
		}
		outputProps[key] = userProps[key];
	}
	for (let i = 0; i < propsList.length; i += 1) {
		let props;
		const propsOrGetProps = propsList[i]?.[elementKey];
		if (typeof propsOrGetProps === "function") props = userProps ? propsOrGetProps(userProps) : null;
		else props = propsOrGetProps;
		if (!props) continue;
		mutablyMergeProps(outputProps, props, isItem, eventHandlers);
	}
	mutablyMergeProps(outputProps, userProps, isItem, eventHandlers);
	return outputProps;
}
function mutablyMergeProps(outputProps, props, isItem, eventHandlers) {
	for (const key in props) {
		const value = props[key];
		if (isItem && (key === "active" || key === "selected")) continue;
		if (!key.startsWith("on")) outputProps[key] = value;
		else {
			if (!eventHandlers.has(key)) eventHandlers.set(key, []);
			if (typeof value === "function") {
				eventHandlers.get(key)?.push(value);
				outputProps[key] = (...args) => {
					return eventHandlers.get(key)?.map((fn) => fn(...args)).find((val) => val !== void 0);
				};
			}
		}
	}
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/hooks/useListNavigation.js
var ESCAPE = "Escape";
function doSwitch(orientation, vertical, horizontal) {
	switch (orientation) {
		case "vertical": return vertical;
		case "horizontal": return horizontal;
		default: return vertical || horizontal;
	}
}
function isMainOrientationKey(key, orientation) {
	return doSwitch(orientation, key === "ArrowUp" || key === "ArrowDown", key === "ArrowLeft" || key === "ArrowRight");
}
function isMainOrientationToEndKey(key, orientation, rtl) {
	return doSwitch(orientation, key === "ArrowDown", rtl ? key === "ArrowLeft" : key === "ArrowRight") || key === "Enter" || key === " " || key === "";
}
function isCrossOrientationOpenKey(key, orientation, rtl) {
	return doSwitch(orientation, rtl ? key === ARROW_LEFT : key === ARROW_RIGHT, key === ARROW_DOWN);
}
function isCrossOrientationCloseKey(key, orientation, rtl, cols) {
	const vertical = rtl ? key === ARROW_RIGHT : key === ARROW_LEFT;
	const horizontal = key === ARROW_UP;
	if (orientation === "both" || orientation === "horizontal" && cols && cols > 1) return key === ESCAPE;
	return doSwitch(orientation, vertical, horizontal);
}
/**
* Adds arrow key-based navigation of a list of items, either using real DOM
* focus or virtual focus.
* @see https://floating-ui.com/docs/useListNavigation
*/
function useListNavigation(context, props) {
	const store = "rootStore" in context ? context.rootStore : context;
	const open = store.useState("open");
	const floatingElement = store.useState("floatingElement");
	const domReferenceElement = store.useState("domReferenceElement");
	const dataRef = store.context.dataRef;
	const { listRef, activeIndex, onNavigate: onNavigateProp = () => {}, enabled = true, selectedIndex = null, allowEscape = false, loopFocus = false, nested = false, rtl = false, virtual = false, focusItemOnOpen = "auto", focusItemOnHover = true, openOnArrowKeyDown = true, disabledIndices = void 0, orientation = "vertical", parentOrientation, cols = 1, id, resetOnPointerLeave = true, externalTree } = props;
	const floatingFocusElementRef = useValueAsRef(getFloatingFocusElement(floatingElement));
	const parentId = useFloatingParentNodeId();
	const tree = useFloatingTree(externalTree);
	useIsoLayoutEffect(() => {
		dataRef.current.orientation = orientation;
	}, [dataRef, orientation]);
	const typeableComboboxReference = isTypeableCombobox(domReferenceElement);
	const focusItemOnOpenRef = import_react.useRef(focusItemOnOpen);
	const indexRef = import_react.useRef(selectedIndex ?? -1);
	const keyRef = import_react.useRef(null);
	const isPointerModalityRef = import_react.useRef(true);
	const onNavigate = useStableCallback((event) => {
		onNavigateProp(indexRef.current === -1 ? null : indexRef.current, event);
	});
	const previousOnNavigateRef = import_react.useRef(onNavigate);
	const previousMountedRef = import_react.useRef(!!floatingElement);
	const previousOpenRef = import_react.useRef(open);
	const forceSyncFocusRef = import_react.useRef(false);
	const forceScrollIntoViewRef = import_react.useRef(false);
	const disabledIndicesRef = useValueAsRef(disabledIndices);
	const latestOpenRef = useValueAsRef(open);
	const selectedIndexRef = useValueAsRef(selectedIndex);
	const resetOnPointerLeaveRef = useValueAsRef(resetOnPointerLeave);
	const focusItem = useStableCallback(() => {
		function runFocus(item) {
			if (virtual) tree?.events.emit("virtualfocus", item);
			else enqueueFocus(item, {
				sync: forceSyncFocusRef.current,
				preventScroll: true
			});
		}
		const initialItem = listRef.current[indexRef.current];
		const forceScrollIntoView = forceScrollIntoViewRef.current;
		if (initialItem) runFocus(initialItem);
		(forceSyncFocusRef.current ? (v) => v() : requestAnimationFrame)(() => {
			const waitedItem = listRef.current[indexRef.current] || initialItem;
			if (!waitedItem) return;
			if (!initialItem) runFocus(waitedItem);
			if (item && (forceScrollIntoView || !isPointerModalityRef.current)) waitedItem.scrollIntoView?.({
				block: "nearest",
				inline: "nearest"
			});
		});
	});
	useIsoLayoutEffect(() => {
		if (!enabled) return;
		if (open && floatingElement) {
			indexRef.current = selectedIndex ?? -1;
			if (focusItemOnOpenRef.current && selectedIndex != null) {
				forceScrollIntoViewRef.current = true;
				onNavigate();
			}
		} else if (previousMountedRef.current) {
			indexRef.current = -1;
			previousOnNavigateRef.current();
		}
	}, [
		enabled,
		open,
		floatingElement,
		selectedIndex,
		onNavigate
	]);
	useIsoLayoutEffect(() => {
		if (!enabled) return;
		if (!open) {
			forceSyncFocusRef.current = false;
			return;
		}
		if (!floatingElement) return;
		if (activeIndex == null) {
			forceSyncFocusRef.current = false;
			if (selectedIndexRef.current != null) return;
			if (previousMountedRef.current) {
				indexRef.current = -1;
				focusItem();
			}
			if ((!previousOpenRef.current || !previousMountedRef.current) && focusItemOnOpenRef.current && (keyRef.current != null || focusItemOnOpenRef.current === true && keyRef.current == null)) {
				let runs = 0;
				const waitForListPopulated = () => {
					if (listRef.current[0] == null) {
						if (runs < 2) (runs ? requestAnimationFrame : queueMicrotask)(waitForListPopulated);
						runs += 1;
					} else {
						indexRef.current = keyRef.current == null || isMainOrientationToEndKey(keyRef.current, orientation, rtl) || nested ? getMinListIndex(listRef) : getMaxListIndex(listRef);
						keyRef.current = null;
						onNavigate();
					}
				};
				waitForListPopulated();
			}
		} else if (!isIndexOutOfListBounds(listRef, activeIndex)) {
			indexRef.current = activeIndex;
			focusItem();
			forceScrollIntoViewRef.current = false;
		}
	}, [
		enabled,
		open,
		floatingElement,
		activeIndex,
		selectedIndexRef,
		nested,
		listRef,
		orientation,
		rtl,
		onNavigate,
		focusItem,
		disabledIndicesRef
	]);
	useIsoLayoutEffect(() => {
		if (!enabled || floatingElement || !tree || virtual || !previousMountedRef.current) return;
		const nodes = tree.nodesRef.current;
		const parent = nodes.find((node) => node.id === parentId)?.context?.elements.floating;
		const activeEl = activeElement(ownerDocument(floatingElement));
		const treeContainsActiveEl = nodes.some((node) => node.context && contains(node.context.elements.floating, activeEl));
		if (parent && !treeContainsActiveEl && isPointerModalityRef.current) parent.focus({ preventScroll: true });
	}, [
		enabled,
		floatingElement,
		tree,
		parentId,
		virtual
	]);
	useIsoLayoutEffect(() => {
		previousOnNavigateRef.current = onNavigate;
		previousOpenRef.current = open;
		previousMountedRef.current = !!floatingElement;
	});
	useIsoLayoutEffect(() => {
		if (!open) {
			keyRef.current = null;
			focusItemOnOpenRef.current = focusItemOnOpen;
		}
	}, [open, focusItemOnOpen]);
	const hasActiveIndex = activeIndex != null;
	const item = import_react.useMemo(() => {
		function syncCurrentTarget(event) {
			if (!latestOpenRef.current) return;
			const index = listRef.current.indexOf(event.currentTarget);
			if (index !== -1 && indexRef.current !== index) {
				indexRef.current = index;
				onNavigate(event);
			}
		}
		return {
			onFocus(event) {
				forceSyncFocusRef.current = true;
				syncCurrentTarget(event);
			},
			onClick: ({ currentTarget }) => currentTarget.focus({ preventScroll: true }),
			onMouseMove(event) {
				forceSyncFocusRef.current = true;
				forceScrollIntoViewRef.current = false;
				if (focusItemOnHover) syncCurrentTarget(event);
			},
			onPointerLeave(event) {
				if (!latestOpenRef.current || !isPointerModalityRef.current || event.pointerType === "touch") return;
				forceSyncFocusRef.current = true;
				const relatedTarget = event.relatedTarget;
				if (!focusItemOnHover || listRef.current.includes(relatedTarget)) return;
				if (!resetOnPointerLeaveRef.current) return;
				enqueueFocus(null, { sync: true });
				indexRef.current = -1;
				onNavigate(event);
				if (!virtual) {
					const floatingFocusEl = floatingFocusElementRef.current;
					const activeEl = activeElement(ownerDocument(floatingFocusEl));
					if (floatingFocusEl && contains(floatingFocusEl, activeEl)) floatingFocusEl.focus({ preventScroll: true });
				}
			}
		};
	}, [
		latestOpenRef,
		floatingFocusElementRef,
		focusItemOnHover,
		listRef,
		onNavigate,
		resetOnPointerLeaveRef,
		virtual
	]);
	const getParentOrientation = import_react.useCallback(() => {
		return parentOrientation ?? tree?.nodesRef.current.find((node) => node.id === parentId)?.context?.dataRef?.current.orientation;
	}, [
		parentId,
		tree,
		parentOrientation
	]);
	const commonOnKeyDown = useStableCallback((event) => {
		isPointerModalityRef.current = false;
		forceSyncFocusRef.current = true;
		if (event.which === 229) return;
		if (!latestOpenRef.current && event.currentTarget === floatingFocusElementRef.current) return;
		if (nested && isCrossOrientationCloseKey(event.key, orientation, rtl, cols)) {
			if (!isMainOrientationKey(event.key, getParentOrientation())) stopEvent(event);
			store.setOpen(false, createChangeEventDetails(listNavigation, event.nativeEvent));
			if (isHTMLElement(domReferenceElement)) if (virtual) tree?.events.emit("virtualfocus", domReferenceElement);
			else domReferenceElement.focus();
			return;
		}
		const currentIndex = indexRef.current;
		const minIndex = getMinListIndex(listRef, disabledIndices);
		const maxIndex = getMaxListIndex(listRef, disabledIndices);
		if (!typeableComboboxReference) {
			if (event.key === "Home") {
				stopEvent(event);
				indexRef.current = minIndex;
				onNavigate(event);
			}
			if (event.key === "End") {
				stopEvent(event);
				indexRef.current = maxIndex;
				onNavigate(event);
			}
		}
		if (cols > 1) {
			const sizes = Array.from({ length: listRef.current.length }, () => ({
				width: 1,
				height: 1
			}));
			const cellMap = createGridCellMap(sizes, cols, false);
			const minGridIndex = cellMap.findIndex((index) => index != null && !isListIndexDisabled(listRef, index, disabledIndices));
			const maxGridIndex = cellMap.reduce((foundIndex, index, cellIndex) => index != null && !isListIndexDisabled(listRef, index, disabledIndices) ? cellIndex : foundIndex, -1);
			const index = cellMap[getGridNavigatedIndex({ current: cellMap.map((itemIndex) => itemIndex != null ? listRef.current[itemIndex] : null) }, {
				event,
				orientation,
				loopFocus,
				rtl,
				cols,
				disabledIndices: getGridCellIndices([...(typeof disabledIndices !== "function" ? disabledIndices : null) || listRef.current.map((_, listIndex) => isListIndexDisabled(listRef, listIndex, disabledIndices) ? listIndex : void 0), void 0], cellMap),
				minIndex: minGridIndex,
				maxIndex: maxGridIndex,
				prevIndex: getGridCellIndexOfCorner(indexRef.current > maxIndex ? minIndex : indexRef.current, sizes, cellMap, cols, event.key === "ArrowDown" ? "bl" : event.key === (rtl ? "ArrowLeft" : "ArrowRight") ? "tr" : "tl"),
				stopEvent: true
			})];
			if (index != null) {
				indexRef.current = index;
				onNavigate(event);
			}
			if (orientation === "both") return;
		}
		if (isMainOrientationKey(event.key, orientation)) {
			stopEvent(event);
			if (open && !virtual && activeElement(event.currentTarget.ownerDocument) === event.currentTarget) {
				indexRef.current = isMainOrientationToEndKey(event.key, orientation, rtl) ? minIndex : maxIndex;
				onNavigate(event);
				return;
			}
			if (isMainOrientationToEndKey(event.key, orientation, rtl)) if (loopFocus) if (currentIndex >= maxIndex) if (allowEscape && currentIndex !== listRef.current.length) indexRef.current = -1;
			else {
				forceSyncFocusRef.current = false;
				indexRef.current = minIndex;
			}
			else indexRef.current = findNonDisabledListIndex(listRef, {
				startingIndex: currentIndex,
				disabledIndices
			});
			else indexRef.current = Math.min(maxIndex, findNonDisabledListIndex(listRef, {
				startingIndex: currentIndex,
				disabledIndices
			}));
			else if (loopFocus) if (currentIndex <= minIndex) if (allowEscape && currentIndex !== -1) indexRef.current = listRef.current.length;
			else {
				forceSyncFocusRef.current = false;
				indexRef.current = maxIndex;
			}
			else indexRef.current = findNonDisabledListIndex(listRef, {
				startingIndex: currentIndex,
				decrement: true,
				disabledIndices
			});
			else indexRef.current = Math.max(minIndex, findNonDisabledListIndex(listRef, {
				startingIndex: currentIndex,
				decrement: true,
				disabledIndices
			}));
			if (isIndexOutOfListBounds(listRef, indexRef.current)) indexRef.current = -1;
			onNavigate(event);
		}
	});
	const ariaActiveDescendantProp = import_react.useMemo(() => {
		return virtual && open && hasActiveIndex && { "aria-activedescendant": `${id}-${activeIndex}` };
	}, [
		virtual,
		open,
		hasActiveIndex,
		id,
		activeIndex
	]);
	const floating = import_react.useMemo(() => {
		return {
			"aria-orientation": orientation === "both" ? void 0 : orientation,
			...!typeableComboboxReference ? ariaActiveDescendantProp : {},
			onKeyDown(event) {
				if (event.key === "Tab" && event.shiftKey && open && !virtual) {
					const target = getTarget(event.nativeEvent);
					if (target && !contains(floatingFocusElementRef.current, target)) return;
					stopEvent(event);
					store.setOpen(false, createChangeEventDetails(focusOut, event.nativeEvent));
					if (isHTMLElement(domReferenceElement)) domReferenceElement.focus();
					return;
				}
				commonOnKeyDown(event);
			},
			onPointerMove() {
				isPointerModalityRef.current = true;
			}
		};
	}, [
		ariaActiveDescendantProp,
		commonOnKeyDown,
		floatingFocusElementRef,
		orientation,
		typeableComboboxReference,
		store,
		open,
		virtual,
		domReferenceElement
	]);
	const trigger = import_react.useMemo(() => {
		function checkVirtualMouse(event) {
			if (focusItemOnOpen === "auto" && isVirtualClick(event.nativeEvent)) focusItemOnOpenRef.current = !virtual;
		}
		function checkVirtualPointer(event) {
			focusItemOnOpenRef.current = focusItemOnOpen;
			if (focusItemOnOpen === "auto" && isVirtualPointerEvent(event.nativeEvent)) focusItemOnOpenRef.current = true;
		}
		return {
			onKeyDown(event) {
				const currentOpen = store.select("open");
				isPointerModalityRef.current = false;
				const isArrowKey = event.key.startsWith("Arrow");
				const isParentCrossOpenKey = isCrossOrientationOpenKey(event.key, getParentOrientation(), rtl);
				const isMainKey = isMainOrientationKey(event.key, orientation);
				const isNavigationKey = (nested ? isParentCrossOpenKey : isMainKey) || event.key === "Enter" || event.key.trim() === "";
				if (virtual && currentOpen) return commonOnKeyDown(event);
				if (!currentOpen && !openOnArrowKeyDown && isArrowKey) return;
				if (isNavigationKey) {
					const isParentMainKey = isMainOrientationKey(event.key, getParentOrientation());
					keyRef.current = nested && isParentMainKey ? null : event.key;
				}
				if (nested) {
					if (isParentCrossOpenKey) {
						stopEvent(event);
						if (currentOpen) {
							indexRef.current = getMinListIndex(listRef, disabledIndicesRef.current);
							onNavigate(event);
						} else store.setOpen(true, createChangeEventDetails(listNavigation, event.nativeEvent, event.currentTarget));
					}
					return;
				}
				if (isMainKey) {
					if (selectedIndexRef.current != null) indexRef.current = selectedIndexRef.current;
					stopEvent(event);
					if (!currentOpen && openOnArrowKeyDown) store.setOpen(true, createChangeEventDetails(listNavigation, event.nativeEvent, event.currentTarget));
					else commonOnKeyDown(event);
					if (currentOpen) onNavigate(event);
				}
			},
			onFocus(event) {
				if (store.select("open") && !virtual) {
					indexRef.current = -1;
					onNavigate(event);
				}
			},
			onPointerDown: checkVirtualPointer,
			onPointerEnter: checkVirtualPointer,
			onMouseDown: checkVirtualMouse,
			onClick: checkVirtualMouse
		};
	}, [
		commonOnKeyDown,
		disabledIndicesRef,
		focusItemOnOpen,
		listRef,
		nested,
		onNavigate,
		store,
		openOnArrowKeyDown,
		orientation,
		getParentOrientation,
		rtl,
		selectedIndexRef,
		virtual
	]);
	const reference = import_react.useMemo(() => {
		return {
			...ariaActiveDescendantProp,
			...trigger
		};
	}, [ariaActiveDescendantProp, trigger]);
	return import_react.useMemo(() => enabled ? {
		reference,
		floating,
		item,
		trigger
	} : {}, [
		enabled,
		reference,
		floating,
		trigger,
		item
	]);
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/hooks/useTypeahead.js
/**
* Provides a matching callback that can be used to focus an item as the user
* types, often used in tandem with `useListNavigation()`.
* @see https://floating-ui.com/docs/useTypeahead
*/
function useTypeahead(context, props) {
	const store = "rootStore" in context ? context.rootStore : context;
	const dataRef = store.context.dataRef;
	const open = store.useState("open");
	const { listRef, elementsRef, activeIndex, onMatch: onMatchProp, onTypingChange, enabled = true, resetMs = 750, selectedIndex = null } = props;
	const timeout = useTimeout();
	const stringRef = import_react.useRef("");
	const prevIndexRef = import_react.useRef(selectedIndex ?? activeIndex ?? -1);
	const matchIndexRef = import_react.useRef(null);
	useIsoLayoutEffect(() => {
		if (!open && selectedIndex !== null) return;
		timeout.clear();
		matchIndexRef.current = null;
		if (stringRef.current !== "") stringRef.current = "";
	}, [
		open,
		selectedIndex,
		timeout
	]);
	useIsoLayoutEffect(() => {
		if (open && stringRef.current === "") prevIndexRef.current = selectedIndex ?? activeIndex ?? -1;
	}, [
		open,
		selectedIndex,
		activeIndex
	]);
	const setTypingChange = useStableCallback((value) => {
		if (value) {
			if (!dataRef.current.typing) {
				dataRef.current.typing = value;
				onTypingChange?.(value);
			}
		} else if (dataRef.current.typing) {
			dataRef.current.typing = value;
			onTypingChange?.(value);
		}
	});
	const onKeyDown = useStableCallback((event) => {
		function isVisible(index) {
			const element = elementsRef?.current[index];
			return !element || isElementVisible(element);
		}
		function getMatchingIndex(list, string, startIndex = 0) {
			if (list.length === 0) return -1;
			const normalizedStartIndex = (startIndex % list.length + list.length) % list.length;
			const lowerString = string.toLocaleLowerCase();
			for (let offset = 0; offset < list.length; offset += 1) {
				const index = (normalizedStartIndex + offset) % list.length;
				if (!list[index]?.toLocaleLowerCase().startsWith(lowerString) || !isVisible(index)) continue;
				return index;
			}
			return -1;
		}
		const listContent = listRef.current;
		if (stringRef.current.length > 0 && event.key === " ") {
			stopEvent(event);
			setTypingChange(true);
		}
		if (stringRef.current.length > 0 && stringRef.current[0] !== " ") {
			if (getMatchingIndex(listContent, stringRef.current) === -1 && event.key !== " ") setTypingChange(false);
		}
		if (listContent == null || event.key.length !== 1 || event.ctrlKey || event.metaKey || event.altKey) return;
		if (open && event.key !== " ") {
			stopEvent(event);
			setTypingChange(true);
		}
		const isNewSession = stringRef.current === "";
		if (isNewSession) prevIndexRef.current = selectedIndex ?? activeIndex ?? -1;
		if (listContent.every((text) => text ? text[0]?.toLocaleLowerCase() !== text[1]?.toLocaleLowerCase() : true) && stringRef.current === event.key) {
			stringRef.current = "";
			prevIndexRef.current = matchIndexRef.current;
		}
		stringRef.current += event.key;
		timeout.start(resetMs, () => {
			stringRef.current = "";
			prevIndexRef.current = matchIndexRef.current;
			setTypingChange(false);
		});
		const startIndex = ((isNewSession ? selectedIndex ?? activeIndex ?? -1 : prevIndexRef.current) ?? 0) + 1;
		const index = getMatchingIndex(listContent, stringRef.current, startIndex);
		if (index !== -1) {
			onMatchProp?.(index);
			matchIndexRef.current = index;
		} else if (event.key !== " ") {
			stringRef.current = "";
			setTypingChange(false);
		}
	});
	const onBlur = useStableCallback((event) => {
		const next = event.relatedTarget;
		const currentDomReferenceElement = store.select("domReferenceElement");
		const currentFloatingElement = store.select("floatingElement");
		const withinReference = contains(currentDomReferenceElement, next);
		const withinFloating = contains(currentFloatingElement, next);
		if (withinReference || withinFloating) return;
		timeout.clear();
		stringRef.current = "";
		prevIndexRef.current = matchIndexRef.current;
		setTypingChange(false);
	});
	const reference = import_react.useMemo(() => ({
		onKeyDown,
		onBlur
	}), [onKeyDown, onBlur]);
	const floating = import_react.useMemo(() => {
		return {
			onKeyDown,
			onBlur
		};
	}, [onKeyDown, onBlur]);
	return import_react.useMemo(() => enabled ? {
		reference,
		floating
	} : {}, [
		enabled,
		reference,
		floating
	]);
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/root/SelectRootContext.js
var SelectRootContext = /* @__PURE__ */ import_react.createContext(null);
var SelectFloatingContext = /* @__PURE__ */ import_react.createContext(null);
function useSelectRootContext() {
	const context = import_react.useContext(SelectRootContext);
	if (context === null) throw new Error(formatErrorMessage(60));
	return context;
}
function useSelectFloatingContext() {
	const context = import_react.useContext(SelectFloatingContext);
	if (context === null) throw new Error(formatErrorMessage(61));
	return context;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/itemEquality.js
var defaultItemEquality = (itemValue, selectedValue) => Object.is(itemValue, selectedValue);
function compareItemEquality(itemValue, selectedValue, comparer) {
	if (itemValue == null || selectedValue == null) return Object.is(itemValue, selectedValue);
	return comparer(itemValue, selectedValue);
}
function selectedValueIncludes(selectedValues, itemValue, comparer) {
	if (!selectedValues || selectedValues.length === 0) return false;
	return selectedValues.some((selectedValue) => {
		if (selectedValue === void 0) return false;
		return compareItemEquality(itemValue, selectedValue, comparer);
	});
}
function findItemIndex(itemValues, selectedValue, comparer) {
	if (!itemValues || itemValues.length === 0) return -1;
	return itemValues.findIndex((itemValue) => {
		if (itemValue === void 0) return false;
		return compareItemEquality(itemValue, selectedValue, comparer);
	});
}
function removeItem(selectedValues, itemValue, comparer) {
	return selectedValues.filter((selectedValue) => !compareItemEquality(itemValue, selectedValue, comparer));
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/serializeValue.js
function serializeValue(value) {
	if (value == null) return "";
	if (typeof value === "string") return value;
	try {
		return JSON.stringify(value);
	} catch {
		return String(value);
	}
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/resolveValueLabel.js
function isGroupedItems(items) {
	return items != null && items.length > 0 && typeof items[0] === "object" && items[0] != null && "items" in items[0];
}
/**
* Checks if the items array contains an item with a null value that has a non-null label.
*/
function hasNullItemLabel(items) {
	if (!Array.isArray(items)) return items != null && "null" in items;
	const arrayItems = items;
	if (isGroupedItems(arrayItems)) {
		for (const group of arrayItems) for (const item of group.items) if (item && item.value == null && item.label != null) return true;
		return false;
	}
	for (const item of arrayItems) if (item && item.value == null && item.label != null) return true;
	return false;
}
function stringifyAsLabel(item, itemToStringLabel) {
	if (itemToStringLabel && item != null) return itemToStringLabel(item) ?? "";
	if (item && typeof item === "object") {
		if ("label" in item && item.label != null) return String(item.label);
		if ("value" in item) return String(item.value);
	}
	return serializeValue(item);
}
function stringifyAsValue(item, itemToStringValue) {
	if (itemToStringValue && item != null) return itemToStringValue(item) ?? "";
	if (item && typeof item === "object" && "value" in item && "label" in item) return serializeValue(item.value);
	return serializeValue(item);
}
function resolveSelectedLabel(value, items, itemToStringLabel) {
	function fallback() {
		return stringifyAsLabel(value, itemToStringLabel);
	}
	if (itemToStringLabel && value != null) return itemToStringLabel(value);
	if (value && typeof value === "object" && "label" in value && value.label != null) return value.label;
	if (items && !Array.isArray(items)) return items[value] ?? fallback();
	if (Array.isArray(items)) {
		const arrayItems = items;
		const flatItems = isGroupedItems(arrayItems) ? arrayItems.flatMap((group) => group.items) : arrayItems;
		if (value == null || typeof value !== "object") {
			const match = flatItems.find((item) => item.value === value);
			if (match && match.label != null) return match.label;
			return fallback();
		}
		if ("value" in value) {
			const match = flatItems.find((item) => item && item.value === value.value);
			if (match && match.label != null) return match.label;
		}
	}
	return fallback();
}
function resolveMultipleLabels(values, items, itemToStringLabel) {
	return values.reduce((acc, value, index) => {
		if (index > 0) acc.push(", ");
		acc.push(/* @__PURE__ */ (0, import_jsx_runtime.jsx)(import_react.Fragment, { children: resolveSelectedLabel(value, items, itemToStringLabel) }, index));
		return acc;
	}, []);
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/store.js
var selectors = {
	id: createSelector((state) => state.id),
	labelId: createSelector((state) => state.labelId),
	modal: createSelector((state) => state.modal),
	multiple: createSelector((state) => state.multiple),
	items: createSelector((state) => state.items),
	itemToStringLabel: createSelector((state) => state.itemToStringLabel),
	itemToStringValue: createSelector((state) => state.itemToStringValue),
	isItemEqualToValue: createSelector((state) => state.isItemEqualToValue),
	value: createSelector((state) => state.value),
	hasSelectedValue: createSelector((state) => {
		const { value, multiple, itemToStringValue } = state;
		if (value == null) return false;
		if (multiple && Array.isArray(value)) return value.length > 0;
		return stringifyAsValue(value, itemToStringValue) !== "";
	}),
	hasNullItemLabel: createSelector((state, enabled) => {
		return enabled ? hasNullItemLabel(state.items) : false;
	}),
	open: createSelector((state) => state.open),
	mounted: createSelector((state) => state.mounted),
	forceMount: createSelector((state) => state.forceMount),
	transitionStatus: createSelector((state) => state.transitionStatus),
	openMethod: createSelector((state) => state.openMethod),
	activeIndex: createSelector((state) => state.activeIndex),
	selectedIndex: createSelector((state) => state.selectedIndex),
	isActive: createSelector((state, index) => state.activeIndex === index),
	isSelected: createSelector((state, index, itemValue) => {
		const comparer = state.isItemEqualToValue;
		const storeValue = state.value;
		if (state.multiple) return Array.isArray(storeValue) && storeValue.some((selectedItem) => compareItemEquality(itemValue, selectedItem, comparer));
		if (state.selectedIndex === index && state.selectedIndex !== null) return true;
		return compareItemEquality(itemValue, storeValue, comparer);
	}),
	isSelectedByFocus: createSelector((state, index) => {
		return state.selectedIndex === index;
	}),
	popupProps: createSelector((state) => state.popupProps),
	triggerProps: createSelector((state) => state.triggerProps),
	triggerElement: createSelector((state) => state.triggerElement),
	positionerElement: createSelector((state) => state.positionerElement),
	listElement: createSelector((state) => state.listElement),
	scrollUpArrowVisible: createSelector((state) => state.scrollUpArrowVisible),
	scrollDownArrowVisible: createSelector((state) => state.scrollDownArrowVisible),
	hasScrollArrows: createSelector((state) => state.hasScrollArrows)
};
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/useValueChanged.js
function useValueChanged(value, onChange) {
	const valueRef = import_react.useRef(value);
	const onChangeCallback = useStableCallback(onChange);
	useIsoLayoutEffect(() => {
		if (valueRef.current === value) return;
		onChangeCallback(valueRef.current);
	}, [value, onChangeCallback]);
	useIsoLayoutEffect(() => {
		valueRef.current = value;
	}, [value]);
}
//#endregion
//#region node_modules/.pnpm/@base-ui+utils@0.2.6_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/utils/esm/useEnhancedClickHandler.js
/**
* Provides a cross-browser way to determine the type of the pointer used to click.
* Safari and Firefox do not provide the PointerEvent to the click handler (they use MouseEvent) yet.
* Additionally, this implementation detects if the click was triggered by the keyboard.
*
* @param handler The function to be called when the button is clicked. The first parameter is the original event and the second parameter is the pointer type.
*/
function useEnhancedClickHandler(handler) {
	const lastClickInteractionTypeRef = import_react.useRef("");
	const handlePointerDown = import_react.useCallback((event) => {
		if (event.defaultPrevented) return;
		lastClickInteractionTypeRef.current = event.pointerType;
		handler(event, event.pointerType);
	}, [handler]);
	return {
		onClick: import_react.useCallback((event) => {
			if (event.detail === 0) {
				handler(event, "keyboard");
				return;
			}
			if ("pointerType" in event) handler(event, event.pointerType);
			else handler(event, lastClickInteractionTypeRef.current);
			lastClickInteractionTypeRef.current = "";
		}, [handler]),
		onPointerDown: handlePointerDown
	};
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/useOpenInteractionType.js
/**
* Determines the interaction type (keyboard, mouse, touch, etc.) that opened the component.
*
* @param open The open state of the component.
*/
function useOpenInteractionType(open) {
	const [openMethod, setOpenMethod] = import_react.useState(null);
	const handleTriggerClick = useStableCallback((_, interactionType) => {
		if (!open) setOpenMethod(interactionType || (isIOS ? "touch" : ""));
	});
	useValueChanged(open, (previousOpen) => {
		if (previousOpen && !open) setOpenMethod(null);
	});
	const { onClick, onPointerDown } = useEnhancedClickHandler(handleTriggerClick);
	return import_react.useMemo(() => ({
		openMethod,
		triggerProps: {
			onClick,
			onPointerDown
		}
	}), [
		openMethod,
		onClick,
		onPointerDown
	]);
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/root/SelectRoot.js
/**
* Groups all parts of the select.
* Doesn’t render its own HTML element.
*
* Documentation: [Base UI Select](https://base-ui.com/react/components/select)
*/
function SelectRoot(props) {
	const { id, value: valueProp, defaultValue = null, onValueChange, open: openProp, defaultOpen = false, onOpenChange, name: nameProp, autoComplete, disabled: disabledProp = false, readOnly = false, required = false, modal = true, actionsRef, inputRef, onOpenChangeComplete, items, multiple = false, itemToStringLabel, itemToStringValue, isItemEqualToValue = defaultItemEquality, highlightItemOnHover = true, children } = props;
	const { clearErrors } = useFormContext();
	const { setDirty, setTouched, setFocused, shouldValidateOnChange, validityData, setFilled, name: fieldName, disabled: fieldDisabled, validation, validationMode } = useFieldRootContext();
	const generatedId = useLabelableId({ id });
	const disabled = fieldDisabled || disabledProp;
	const name = fieldName ?? nameProp;
	const [value, setValueUnwrapped] = useControlled({
		controlled: valueProp,
		default: multiple ? defaultValue ?? EMPTY_ARRAY : defaultValue,
		name: "Select",
		state: "value"
	});
	const [open, setOpenUnwrapped] = useControlled({
		controlled: openProp,
		default: defaultOpen,
		name: "Select",
		state: "open"
	});
	const listRef = import_react.useRef([]);
	const labelsRef = import_react.useRef([]);
	const popupRef = import_react.useRef(null);
	const scrollHandlerRef = import_react.useRef(null);
	const scrollArrowsMountedCountRef = import_react.useRef(0);
	const valueRef = import_react.useRef(null);
	const valuesRef = import_react.useRef([]);
	const typingRef = import_react.useRef(false);
	const keyboardActiveRef = import_react.useRef(false);
	const selectedItemTextRef = import_react.useRef(null);
	const selectionRef = import_react.useRef({
		allowSelectedMouseUp: false,
		allowUnselectedMouseUp: false
	});
	const alignItemWithTriggerActiveRef = import_react.useRef(false);
	const { mounted, setMounted, transitionStatus } = useTransitionStatus(open);
	const { openMethod, triggerProps: interactionTypeProps } = useOpenInteractionType(open);
	const store = useRefWithInit(() => new Store({
		id: generatedId,
		labelId: void 0,
		modal,
		multiple,
		itemToStringLabel,
		itemToStringValue,
		isItemEqualToValue,
		value,
		open,
		mounted,
		transitionStatus,
		items,
		forceMount: false,
		openMethod: null,
		activeIndex: null,
		selectedIndex: null,
		popupProps: {},
		triggerProps: {},
		triggerElement: null,
		positionerElement: null,
		listElement: null,
		scrollUpArrowVisible: false,
		scrollDownArrowVisible: false,
		hasScrollArrows: false
	})).current;
	const activeIndex = useStore(store, selectors.activeIndex);
	const selectedIndex = useStore(store, selectors.selectedIndex);
	const triggerElement = useStore(store, selectors.triggerElement);
	const positionerElement = useStore(store, selectors.positionerElement);
	const serializedValue = import_react.useMemo(() => {
		if (multiple && Array.isArray(value) && value.length === 0) return "";
		return stringifyAsValue(value, itemToStringValue);
	}, [
		multiple,
		value,
		itemToStringValue
	]);
	const fieldStringValue = import_react.useMemo(() => {
		if (multiple && Array.isArray(value)) return value.map((currentValue) => stringifyAsValue(currentValue, itemToStringValue));
		return stringifyAsValue(value, itemToStringValue);
	}, [
		multiple,
		value,
		itemToStringValue
	]);
	const controlRef = useValueAsRef(store.state.triggerElement);
	useField({
		id: generatedId,
		commit: validation.commit,
		value,
		controlRef,
		name,
		getValue: () => fieldStringValue
	});
	const initialValueRef = import_react.useRef(value);
	useIsoLayoutEffect(() => {
		if (value !== initialValueRef.current) store.set("forceMount", true);
	}, [store, value]);
	useIsoLayoutEffect(() => {
		setFilled(multiple ? Array.isArray(value) && value.length > 0 : value != null);
	}, [
		multiple,
		value,
		setFilled
	]);
	useIsoLayoutEffect(function syncSelectedIndex() {
		if (open) return;
		const registry = valuesRef.current;
		if (multiple) {
			const currentValue = Array.isArray(value) ? value : [];
			if (currentValue.length === 0) {
				store.set("selectedIndex", null);
				return;
			}
			const lastValue = currentValue[currentValue.length - 1];
			const lastIndex = findItemIndex(registry, lastValue, isItemEqualToValue);
			store.set("selectedIndex", lastIndex === -1 ? null : lastIndex);
			return;
		}
		const index = findItemIndex(registry, value, isItemEqualToValue);
		store.set("selectedIndex", index === -1 ? null : index);
	}, [
		multiple,
		open,
		value,
		valuesRef,
		isItemEqualToValue,
		store
	]);
	useValueChanged(value, () => {
		clearErrors(name);
		setDirty(value !== validityData.initialValue);
		if (shouldValidateOnChange()) validation.commit(value);
		else validation.commit(value, true);
	});
	const setOpen = useStableCallback((nextOpen, eventDetails) => {
		onOpenChange?.(nextOpen, eventDetails);
		if (eventDetails.isCanceled) return;
		setOpenUnwrapped(nextOpen);
		if (!nextOpen && (eventDetails.reason === "focus-out" || eventDetails.reason === "outside-press")) {
			setTouched(true);
			setFocused(false);
			if (validationMode === "onBlur") validation.commit(value);
		}
		if (!nextOpen && store.state.activeIndex !== null) {
			const activeOption = listRef.current[store.state.activeIndex];
			queueMicrotask(() => {
				activeOption?.setAttribute("tabindex", "-1");
			});
		}
	});
	const handleUnmount = useStableCallback(() => {
		setMounted(false);
		store.set("activeIndex", null);
		onOpenChangeComplete?.(false);
	});
	useOpenChangeComplete({
		enabled: !actionsRef,
		open,
		ref: popupRef,
		onComplete() {
			if (!open) handleUnmount();
		}
	});
	import_react.useImperativeHandle(actionsRef, () => ({ unmount: handleUnmount }), [handleUnmount]);
	const setValue = useStableCallback((nextValue, eventDetails) => {
		onValueChange?.(nextValue, eventDetails);
		if (eventDetails.isCanceled) return;
		setValueUnwrapped(nextValue);
	});
	const handleScrollArrowVisibility = useStableCallback(() => {
		const scroller = store.state.listElement || popupRef.current;
		if (!scroller) return;
		const viewportTop = scroller.scrollTop;
		const viewportBottom = scroller.scrollTop + scroller.clientHeight;
		const shouldShowUp = viewportTop > 1;
		const shouldShowDown = viewportBottom < scroller.scrollHeight - 1;
		if (store.state.scrollUpArrowVisible !== shouldShowUp) store.set("scrollUpArrowVisible", shouldShowUp);
		if (store.state.scrollDownArrowVisible !== shouldShowDown) store.set("scrollDownArrowVisible", shouldShowDown);
	});
	const floatingContext = useFloatingRootContext({
		open,
		onOpenChange: setOpen,
		elements: {
			reference: triggerElement,
			floating: positionerElement
		}
	});
	const { getReferenceProps, getFloatingProps, getItemProps } = useInteractions([
		useClick(floatingContext, {
			enabled: !readOnly && !disabled,
			event: "mousedown"
		}),
		useDismiss(floatingContext, { bubbles: false }),
		useListNavigation(floatingContext, {
			enabled: !readOnly && !disabled,
			listRef,
			activeIndex,
			selectedIndex,
			disabledIndices: EMPTY_ARRAY,
			onNavigate(nextActiveIndex) {
				if (nextActiveIndex === null && !open) return;
				store.set("activeIndex", nextActiveIndex);
			},
			focusItemOnHover: false
		}),
		useTypeahead(floatingContext, {
			enabled: !readOnly && !disabled && (open || !multiple),
			listRef: labelsRef,
			activeIndex,
			selectedIndex,
			onMatch(index) {
				if (open) store.set("activeIndex", index);
				else setValue(valuesRef.current[index], createChangeEventDetails("none"));
			},
			onTypingChange(typing) {
				typingRef.current = typing;
			}
		})
	]);
	const mergedTriggerProps = import_react.useMemo(() => {
		return mergeProps$1(getReferenceProps(), interactionTypeProps, generatedId ? { id: generatedId } : EMPTY_OBJECT);
	}, [
		getReferenceProps,
		interactionTypeProps,
		generatedId
	]);
	useOnFirstRender(() => {
		store.update({
			popupProps: getFloatingProps(),
			triggerProps: mergedTriggerProps
		});
	});
	useIsoLayoutEffect(() => {
		store.update({
			id: generatedId,
			modal,
			multiple,
			value,
			open,
			mounted,
			transitionStatus,
			popupProps: getFloatingProps(),
			triggerProps: mergedTriggerProps,
			items,
			itemToStringLabel,
			itemToStringValue,
			isItemEqualToValue,
			openMethod
		});
	}, [
		store,
		generatedId,
		modal,
		multiple,
		value,
		open,
		mounted,
		transitionStatus,
		getFloatingProps,
		mergedTriggerProps,
		items,
		itemToStringLabel,
		itemToStringValue,
		isItemEqualToValue,
		openMethod
	]);
	const contextValue = import_react.useMemo(() => ({
		store,
		name,
		required,
		disabled,
		readOnly,
		multiple,
		itemToStringLabel,
		itemToStringValue,
		highlightItemOnHover,
		setValue,
		setOpen,
		listRef,
		popupRef,
		scrollHandlerRef,
		handleScrollArrowVisibility,
		scrollArrowsMountedCountRef,
		getItemProps,
		events: floatingContext.context.events,
		valueRef,
		valuesRef,
		labelsRef,
		typingRef,
		selectionRef,
		selectedItemTextRef,
		validation,
		onOpenChangeComplete,
		keyboardActiveRef,
		alignItemWithTriggerActiveRef,
		initialValueRef
	}), [
		store,
		name,
		required,
		disabled,
		readOnly,
		multiple,
		itemToStringLabel,
		itemToStringValue,
		highlightItemOnHover,
		setValue,
		setOpen,
		getItemProps,
		floatingContext.context.events,
		validation,
		onOpenChangeComplete,
		handleScrollArrowVisibility
	]);
	const ref = useMergedRefs(inputRef, validation.inputRef);
	const hasMultipleSelection = multiple && Array.isArray(value) && value.length > 0;
	const hiddenInputName = multiple ? void 0 : name;
	const hiddenInputs = import_react.useMemo(() => {
		if (!multiple || !Array.isArray(value) || !name) return null;
		return value.map((v) => {
			const currentSerializedValue = stringifyAsValue(v, itemToStringValue);
			return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("input", {
				type: "hidden",
				name,
				value: currentSerializedValue
			}, currentSerializedValue);
		});
	}, [
		multiple,
		value,
		name,
		itemToStringValue
	]);
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectRootContext.Provider, {
		value: contextValue,
		children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(SelectFloatingContext.Provider, {
			value: floatingContext,
			children: [
				children,
				/* @__PURE__ */ (0, import_jsx_runtime.jsx)("input", {
					...validation.getInputValidationProps({
						onFocus() {
							store.state.triggerElement?.focus({ focusVisible: true });
						},
						onChange(event) {
							if (event.nativeEvent.defaultPrevented) return;
							const nextValue = event.target.value;
							const details = createChangeEventDetails(none, event.nativeEvent);
							function handleChange() {
								if (multiple) return;
								const matchingValue = valuesRef.current.find((v) => {
									if (stringifyAsValue(v, itemToStringValue).toLowerCase() === nextValue.toLowerCase()) return true;
									return false;
								});
								if (matchingValue != null) {
									setDirty(matchingValue !== validityData.initialValue);
									setValue(matchingValue, details);
									if (shouldValidateOnChange()) validation.commit(matchingValue);
								}
							}
							store.set("forceMount", true);
							queueMicrotask(handleChange);
						}
					}),
					id: generatedId && hiddenInputName == null ? `${generatedId}-hidden-input` : void 0,
					name: hiddenInputName,
					autoComplete,
					value: serializedValue,
					disabled,
					required: required && !hasMultipleSelection,
					readOnly,
					ref,
					style: name ? visuallyHiddenInput : visuallyHidden,
					tabIndex: -1,
					"aria-hidden": true
				}),
				hiddenInputs
			]
		})
	});
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/resolveAriaLabelledBy.js
function resolveAriaLabelledBy(fieldLabelId, localLabelId) {
	return fieldLabelId ?? localLabelId;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/popupStateMapping.js
var CommonPopupDataAttributes = function(CommonPopupDataAttributes) {
	/**
	* Present when the popup is open.
	*/
	CommonPopupDataAttributes["open"] = "data-open";
	/**
	* Present when the popup is closed.
	*/
	CommonPopupDataAttributes["closed"] = "data-closed";
	/**
	* Present when the popup is animating in.
	*/
	CommonPopupDataAttributes[CommonPopupDataAttributes["startingStyle"] = TransitionStatusDataAttributes.startingStyle] = "startingStyle";
	/**
	* Present when the popup is animating out.
	*/
	CommonPopupDataAttributes[CommonPopupDataAttributes["endingStyle"] = TransitionStatusDataAttributes.endingStyle] = "endingStyle";
	/**
	* Present when the anchor is hidden.
	*/
	CommonPopupDataAttributes["anchorHidden"] = "data-anchor-hidden";
	/**
	* Indicates which side the popup is positioned relative to the trigger.
	* @type { 'top' | 'bottom' | 'left' | 'right' | 'inline-end' | 'inline-start'}
	*/
	CommonPopupDataAttributes["side"] = "data-side";
	/**
	* Indicates how the popup is aligned relative to specified side.
	* @type {'start' | 'center' | 'end'}
	*/
	CommonPopupDataAttributes["align"] = "data-align";
	return CommonPopupDataAttributes;
}({});
var CommonTriggerDataAttributes = /* @__PURE__ */ function(CommonTriggerDataAttributes) {
	/**
	* Present when the popup is open.
	*/
	CommonTriggerDataAttributes["popupOpen"] = "data-popup-open";
	/**
	* Present when a pressable trigger is pressed.
	*/
	CommonTriggerDataAttributes["pressed"] = "data-pressed";
	return CommonTriggerDataAttributes;
}({});
var TRIGGER_HOOK = { [CommonTriggerDataAttributes.popupOpen]: "" };
var PRESSABLE_TRIGGER_HOOK = {
	[CommonTriggerDataAttributes.popupOpen]: "",
	[CommonTriggerDataAttributes.pressed]: ""
};
var POPUP_OPEN_HOOK = { [CommonPopupDataAttributes.open]: "" };
var POPUP_CLOSED_HOOK = { [CommonPopupDataAttributes.closed]: "" };
var ANCHOR_HIDDEN_HOOK = { [CommonPopupDataAttributes.anchorHidden]: "" };
var triggerOpenStateMapping = { open(value) {
	if (value) return TRIGGER_HOOK;
	return null;
} };
var pressableTriggerOpenStateMapping = { open(value) {
	if (value) return PRESSABLE_TRIGGER_HOOK;
	return null;
} };
var popupStateMapping = {
	open(value) {
		if (value) return POPUP_OPEN_HOOK;
		return POPUP_CLOSED_HOOK;
	},
	anchorHidden(value) {
		if (value) return ANCHOR_HIDDEN_HOOK;
		return null;
	}
};
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/getPseudoElementBounds.js
function getPseudoElementBounds(element) {
	const elementRect = element.getBoundingClientRect();
	const beforeStyles = window.getComputedStyle(element, "::before");
	const afterStyles = window.getComputedStyle(element, "::after");
	if (!(beforeStyles.content !== "none" || afterStyles.content !== "none")) return elementRect;
	const beforeWidth = parseFloat(beforeStyles.width) || 0;
	const beforeHeight = parseFloat(beforeStyles.height) || 0;
	const afterWidth = parseFloat(afterStyles.width) || 0;
	const afterHeight = parseFloat(afterStyles.height) || 0;
	const totalWidth = Math.max(elementRect.width, beforeWidth, afterWidth);
	const totalHeight = Math.max(elementRect.height, beforeHeight, afterHeight);
	const widthDiff = totalWidth - elementRect.width;
	const heightDiff = totalHeight - elementRect.height;
	return {
		left: elementRect.left - widthDiff / 2,
		right: elementRect.right + widthDiff / 2,
		top: elementRect.top - heightDiff / 2,
		bottom: elementRect.bottom + heightDiff / 2
	};
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/trigger/SelectTrigger.js
var BOUNDARY_OFFSET = 2;
var SELECTED_DELAY = 400;
var UNSELECTED_DELAY = 200;
var stateAttributesMapping$2 = {
	...pressableTriggerOpenStateMapping,
	...fieldValidityMapping,
	value: () => null
};
/**
* A button that opens the select popup.
* Renders a `<button>` element.
*
* Documentation: [Base UI Select](https://base-ui.com/react/components/select)
*/
var SelectTrigger$1 = /* @__PURE__ */ import_react.forwardRef(function SelectTrigger(componentProps, forwardedRef) {
	const { render, className, id: idProp, disabled: disabledProp = false, nativeButton = true, ...elementProps } = componentProps;
	const { setTouched, setFocused, validationMode, state: fieldState, disabled: fieldDisabled } = useFieldRootContext();
	const { labelId: fieldLabelId } = useLabelableContext();
	const { store, setOpen, selectionRef, validation, readOnly, required, alignItemWithTriggerActiveRef, disabled: selectDisabled, keyboardActiveRef } = useSelectRootContext();
	const disabled = fieldDisabled || selectDisabled || disabledProp;
	const open = useStore(store, selectors.open);
	const value = useStore(store, selectors.value);
	const triggerProps = useStore(store, selectors.triggerProps);
	const positionerElement = useStore(store, selectors.positionerElement);
	const listElement = useStore(store, selectors.listElement);
	const rootId = useStore(store, selectors.id);
	const selectLabelId = useStore(store, selectors.labelId);
	const hasSelectedValue = useStore(store, selectors.hasSelectedValue);
	const shouldCheckNullItemLabel = !hasSelectedValue && open;
	const hasNullItemLabel = useStore(store, selectors.hasNullItemLabel, shouldCheckNullItemLabel);
	const id = idProp ?? rootId;
	const ariaLabelledBy = resolveAriaLabelledBy(fieldLabelId, selectLabelId);
	useLabelableId({ id });
	const positionerRef = useValueAsRef(positionerElement);
	const triggerRef = import_react.useRef(null);
	const { getButtonProps, buttonRef } = useButton({
		disabled,
		native: nativeButton
	});
	const mergedRef = useMergedRefs(forwardedRef, triggerRef, buttonRef, useStableCallback((element) => {
		store.set("triggerElement", element);
	}));
	const timeoutFocus = useTimeout();
	const timeoutMouseDown = useTimeout();
	const selectedDelayTimeout = useTimeout();
	const unselectedDelayTimeout = useTimeout();
	import_react.useEffect(() => {
		if (open) {
			if (!(hasSelectedValue || hasNullItemLabel)) selectedDelayTimeout.start(SELECTED_DELAY, () => {
				selectionRef.current.allowUnselectedMouseUp = true;
				selectionRef.current.allowSelectedMouseUp = true;
			});
			else unselectedDelayTimeout.start(UNSELECTED_DELAY, () => {
				selectionRef.current.allowUnselectedMouseUp = true;
				selectedDelayTimeout.start(UNSELECTED_DELAY, () => {
					selectionRef.current.allowSelectedMouseUp = true;
				});
			});
			return () => {
				selectedDelayTimeout.clear();
				unselectedDelayTimeout.clear();
			};
		}
		selectionRef.current = {
			allowSelectedMouseUp: false,
			allowUnselectedMouseUp: false
		};
		timeoutMouseDown.clear();
	}, [
		open,
		hasSelectedValue,
		hasNullItemLabel,
		selectionRef,
		timeoutMouseDown,
		selectedDelayTimeout,
		unselectedDelayTimeout
	]);
	const ariaControlsId = import_react.useMemo(() => {
		return listElement?.id ?? getFloatingFocusElement(positionerElement)?.id;
	}, [listElement, positionerElement]);
	const props = mergeProps$1(triggerProps, {
		id,
		role: "combobox",
		"aria-expanded": open ? "true" : "false",
		"aria-haspopup": "listbox",
		"aria-controls": open ? ariaControlsId : void 0,
		"aria-labelledby": ariaLabelledBy,
		"aria-readonly": readOnly || void 0,
		"aria-required": required || void 0,
		tabIndex: disabled ? -1 : 0,
		ref: mergedRef,
		onFocus(event) {
			setFocused(true);
			if (open && alignItemWithTriggerActiveRef.current) setOpen(false, createChangeEventDetails(none, event.nativeEvent));
			timeoutFocus.start(0, () => {
				store.set("forceMount", true);
			});
		},
		onBlur(event) {
			if (contains(positionerElement, event.relatedTarget)) return;
			setTouched(true);
			setFocused(false);
			if (validationMode === "onBlur") validation.commit(value);
		},
		onPointerMove() {
			keyboardActiveRef.current = false;
		},
		onKeyDown() {
			keyboardActiveRef.current = true;
		},
		onMouseDown(event) {
			if (open) return;
			const doc = ownerDocument(event.currentTarget);
			function handleMouseUp(mouseEvent) {
				if (!triggerRef.current) return;
				const mouseUpTarget = mouseEvent.target;
				if (contains(triggerRef.current, mouseUpTarget) || contains(positionerRef.current, mouseUpTarget) || mouseUpTarget === triggerRef.current) return;
				const bounds = getPseudoElementBounds(triggerRef.current);
				if (mouseEvent.clientX >= bounds.left - BOUNDARY_OFFSET && mouseEvent.clientX <= bounds.right + BOUNDARY_OFFSET && mouseEvent.clientY >= bounds.top - BOUNDARY_OFFSET && mouseEvent.clientY <= bounds.bottom + BOUNDARY_OFFSET) return;
				setOpen(false, createChangeEventDetails(cancelOpen, mouseEvent));
			}
			timeoutMouseDown.start(0, () => {
				doc.addEventListener("mouseup", handleMouseUp, { once: true });
			});
		}
	}, validation.getValidationProps, elementProps, getButtonProps);
	props.role = "combobox";
	const state = {
		...fieldState,
		open,
		disabled,
		value,
		readOnly,
		placeholder: !hasSelectedValue
	};
	return useRenderElement("button", componentProps, {
		ref: [forwardedRef, triggerRef],
		state,
		stateAttributesMapping: stateAttributesMapping$2,
		props
	});
});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/value/SelectValue.js
var stateAttributesMapping$1 = { value: () => null };
/**
* A text label of the currently selected item.
* Renders a `<span>` element.
*
* Documentation: [Base UI Select](https://base-ui.com/react/components/select)
*/
var SelectValue$1 = /* @__PURE__ */ import_react.forwardRef(function SelectValue(componentProps, forwardedRef) {
	const { className, render, children: childrenProp, placeholder, ...elementProps } = componentProps;
	const { store, valueRef } = useSelectRootContext();
	const value = useStore(store, selectors.value);
	const items = useStore(store, selectors.items);
	const itemToStringLabel = useStore(store, selectors.itemToStringLabel);
	const hasSelectedValue = useStore(store, selectors.hasSelectedValue);
	const shouldCheckNullItemLabel = !hasSelectedValue && placeholder != null && childrenProp == null;
	const hasNullLabel = useStore(store, selectors.hasNullItemLabel, shouldCheckNullItemLabel);
	const state = {
		value,
		placeholder: !hasSelectedValue
	};
	let children = null;
	if (typeof childrenProp === "function") children = childrenProp(value);
	else if (childrenProp != null) children = childrenProp;
	else if (!hasSelectedValue && placeholder != null && !hasNullLabel) children = placeholder;
	else if (Array.isArray(value)) children = resolveMultipleLabels(value, items, itemToStringLabel);
	else children = resolveSelectedLabel(value, items, itemToStringLabel);
	return useRenderElement("span", componentProps, {
		state,
		ref: [forwardedRef, valueRef],
		props: [{ children }, elementProps],
		stateAttributesMapping: stateAttributesMapping$1
	});
});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/icon/SelectIcon.js
/**
* An icon that indicates that the trigger button opens a select popup.
* Renders a `<span>` element.
*
* Documentation: [Base UI Select](https://base-ui.com/react/components/select)
*/
var SelectIcon = /* @__PURE__ */ import_react.forwardRef(function SelectIcon(componentProps, forwardedRef) {
	const { className, render, ...elementProps } = componentProps;
	const { store } = useSelectRootContext();
	return useRenderElement("span", componentProps, {
		state: { open: useStore(store, selectors.open) },
		ref: forwardedRef,
		props: [{
			"aria-hidden": true,
			children: "▼"
		}, elementProps],
		stateAttributesMapping: triggerOpenStateMapping
	});
});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/portal/SelectPortalContext.js
var SelectPortalContext = /* @__PURE__ */ import_react.createContext(void 0);
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/portal/SelectPortal.js
/**
* A portal element that moves the popup to a different part of the DOM.
* By default, the portal element is appended to `<body>`.
* Renders a `<div>` element.
*
* Documentation: [Base UI Select](https://base-ui.com/react/components/select)
*/
var SelectPortal = /* @__PURE__ */ import_react.forwardRef(function SelectPortal(portalProps, forwardedRef) {
	const { store } = useSelectRootContext();
	const mounted = useStore(store, selectors.mounted);
	const forceMount = useStore(store, selectors.forceMount);
	if (!(mounted || forceMount)) return null;
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectPortalContext.Provider, {
		value: true,
		children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(FloatingPortal, {
			ref: forwardedRef,
			...portalProps
		})
	});
});
//#endregion
//#region node_modules/.pnpm/@base-ui+utils@0.2.6_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/utils/esm/useScrollLock.js
var originalHtmlStyles = {};
var originalBodyStyles = {};
var originalHtmlScrollBehavior = "";
function hasInsetScrollbars(referenceElement) {
	if (typeof document === "undefined") return false;
	const doc = ownerDocument(referenceElement);
	return getWindow(doc).innerWidth - doc.documentElement.clientWidth > 0;
}
function supportsStableScrollbarGutter(referenceElement) {
	if (!(typeof CSS !== "undefined" && CSS.supports && CSS.supports("scrollbar-gutter", "stable")) || typeof document === "undefined") return false;
	const doc = ownerDocument(referenceElement);
	const html = doc.documentElement;
	const body = doc.body;
	const scrollContainer = isOverflowElement(html) ? html : body;
	const originalScrollContainerOverflowY = scrollContainer.style.overflowY;
	const originalHtmlStyleGutter = html.style.scrollbarGutter;
	html.style.scrollbarGutter = "stable";
	scrollContainer.style.overflowY = "scroll";
	const before = scrollContainer.offsetWidth;
	scrollContainer.style.overflowY = "hidden";
	const after = scrollContainer.offsetWidth;
	scrollContainer.style.overflowY = originalScrollContainerOverflowY;
	html.style.scrollbarGutter = originalHtmlStyleGutter;
	return before === after;
}
function preventScrollOverlayScrollbars(referenceElement) {
	const doc = ownerDocument(referenceElement);
	const html = doc.documentElement;
	const body = doc.body;
	const elementToLock = isOverflowElement(html) ? html : body;
	const originalElementToLockStyles = {
		overflowY: elementToLock.style.overflowY,
		overflowX: elementToLock.style.overflowX
	};
	Object.assign(elementToLock.style, {
		overflowY: "hidden",
		overflowX: "hidden"
	});
	return () => {
		Object.assign(elementToLock.style, originalElementToLockStyles);
	};
}
function preventScrollInsetScrollbars(referenceElement) {
	const doc = ownerDocument(referenceElement);
	const html = doc.documentElement;
	const body = doc.body;
	const win = getWindow(html);
	let scrollTop = 0;
	let scrollLeft = 0;
	let updateGutterOnly = false;
	const resizeFrame = AnimationFrame.create();
	if (isWebKit$1 && (win.visualViewport?.scale ?? 1) !== 1) return () => {};
	function lockScroll() {
		const htmlStyles = win.getComputedStyle(html);
		const bodyStyles = win.getComputedStyle(body);
		const scrollbarGutterValue = (htmlStyles.scrollbarGutter || "").includes("both-edges") ? "stable both-edges" : "stable";
		scrollTop = html.scrollTop;
		scrollLeft = html.scrollLeft;
		originalHtmlStyles = {
			scrollbarGutter: html.style.scrollbarGutter,
			overflowY: html.style.overflowY,
			overflowX: html.style.overflowX
		};
		originalHtmlScrollBehavior = html.style.scrollBehavior;
		originalBodyStyles = {
			position: body.style.position,
			height: body.style.height,
			width: body.style.width,
			boxSizing: body.style.boxSizing,
			overflowY: body.style.overflowY,
			overflowX: body.style.overflowX,
			scrollBehavior: body.style.scrollBehavior
		};
		const isScrollableY = html.scrollHeight > html.clientHeight;
		const isScrollableX = html.scrollWidth > html.clientWidth;
		const hasConstantOverflowY = htmlStyles.overflowY === "scroll" || bodyStyles.overflowY === "scroll";
		const hasConstantOverflowX = htmlStyles.overflowX === "scroll" || bodyStyles.overflowX === "scroll";
		const scrollbarWidth = Math.max(0, win.innerWidth - body.clientWidth);
		const scrollbarHeight = Math.max(0, win.innerHeight - body.clientHeight);
		const marginY = parseFloat(bodyStyles.marginTop) + parseFloat(bodyStyles.marginBottom);
		const marginX = parseFloat(bodyStyles.marginLeft) + parseFloat(bodyStyles.marginRight);
		const elementToLock = isOverflowElement(html) ? html : body;
		updateGutterOnly = supportsStableScrollbarGutter(referenceElement);
		if (updateGutterOnly) {
			html.style.scrollbarGutter = scrollbarGutterValue;
			elementToLock.style.overflowY = "hidden";
			elementToLock.style.overflowX = "hidden";
			return;
		}
		Object.assign(html.style, {
			scrollbarGutter: scrollbarGutterValue,
			overflowY: "hidden",
			overflowX: "hidden"
		});
		if (isScrollableY || hasConstantOverflowY) html.style.overflowY = "scroll";
		if (isScrollableX || hasConstantOverflowX) html.style.overflowX = "scroll";
		Object.assign(body.style, {
			position: "relative",
			height: marginY || scrollbarHeight ? `calc(100dvh - ${marginY + scrollbarHeight}px)` : "100dvh",
			width: marginX || scrollbarWidth ? `calc(100vw - ${marginX + scrollbarWidth}px)` : "100vw",
			boxSizing: "border-box",
			overflow: "hidden",
			scrollBehavior: "unset"
		});
		body.scrollTop = scrollTop;
		body.scrollLeft = scrollLeft;
		html.setAttribute("data-base-ui-scroll-locked", "");
		html.style.scrollBehavior = "unset";
	}
	function cleanup() {
		Object.assign(html.style, originalHtmlStyles);
		Object.assign(body.style, originalBodyStyles);
		if (!updateGutterOnly) {
			html.scrollTop = scrollTop;
			html.scrollLeft = scrollLeft;
			html.removeAttribute("data-base-ui-scroll-locked");
			html.style.scrollBehavior = originalHtmlScrollBehavior;
		}
	}
	function handleResize() {
		cleanup();
		resizeFrame.request(lockScroll);
	}
	lockScroll();
	win.addEventListener("resize", handleResize);
	return () => {
		resizeFrame.cancel();
		cleanup();
		if (typeof win.removeEventListener === "function") win.removeEventListener("resize", handleResize);
	};
}
var ScrollLocker = class {
	lockCount = 0;
	restore = null;
	timeoutLock = Timeout.create();
	timeoutUnlock = Timeout.create();
	acquire(referenceElement) {
		this.lockCount += 1;
		if (this.lockCount === 1 && this.restore === null) this.timeoutLock.start(0, () => this.lock(referenceElement));
		return this.release;
	}
	release = () => {
		this.lockCount -= 1;
		if (this.lockCount === 0 && this.restore) this.timeoutUnlock.start(0, this.unlock);
	};
	unlock = () => {
		if (this.lockCount === 0 && this.restore) {
			this.restore?.();
			this.restore = null;
		}
	};
	lock(referenceElement) {
		if (this.lockCount === 0 || this.restore !== null) return;
		const html = ownerDocument(referenceElement).documentElement;
		const htmlOverflowY = getWindow(html).getComputedStyle(html).overflowY;
		if (htmlOverflowY === "hidden" || htmlOverflowY === "clip") {
			this.restore = NOOP;
			return;
		}
		this.restore = isIOS || !hasInsetScrollbars(referenceElement) ? preventScrollOverlayScrollbars(referenceElement) : preventScrollInsetScrollbars(referenceElement);
	}
};
var SCROLL_LOCKER = new ScrollLocker();
/**
* Locks the scroll of the document when enabled.
*
* @param enabled - Whether to enable the scroll lock.
* @param referenceElement - Element to use as a reference for lock calculations.
*/
function useScrollLock(enabled = true, referenceElement = null) {
	useIsoLayoutEffect(() => {
		if (!enabled) return;
		return SCROLL_LOCKER.acquire(referenceElement);
	}, [enabled, referenceElement]);
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/floating-ui-react/middleware/arrow.js
/**
* Fork of the original `arrow` middleware from Floating UI that allows
* configuring the offset parent.
*/
var baseArrow = (options) => ({
	name: "arrow",
	options,
	async fn(state) {
		const { x, y, placement, rects, platform, elements, middlewareData } = state;
		const { element, padding = 0, offsetParent = "real" } = evaluate(options, state) || {};
		if (element == null) return {};
		const paddingObject = getPaddingObject(padding);
		const coords = {
			x,
			y
		};
		const axis = getAlignmentAxis(placement);
		const length = getAxisLength(axis);
		const arrowDimensions = await platform.getDimensions(element);
		const isYAxis = axis === "y";
		const minProp = isYAxis ? "top" : "left";
		const maxProp = isYAxis ? "bottom" : "right";
		const clientProp = isYAxis ? "clientHeight" : "clientWidth";
		const endDiff = rects.reference[length] + rects.reference[axis] - coords[axis] - rects.floating[length];
		const startDiff = coords[axis] - rects.reference[axis];
		const arrowOffsetParent = offsetParent === "real" ? await platform.getOffsetParent?.(element) : elements.floating;
		let clientSize = elements.floating[clientProp] || rects.floating[length];
		if (!clientSize || !await platform.isElement?.(arrowOffsetParent)) clientSize = elements.floating[clientProp] || rects.floating[length];
		const centerToReference = endDiff / 2 - startDiff / 2;
		const largestPossiblePadding = clientSize / 2 - arrowDimensions[length] / 2 - 1;
		const minPadding = Math.min(paddingObject[minProp], largestPossiblePadding);
		const maxPadding = Math.min(paddingObject[maxProp], largestPossiblePadding);
		const min = minPadding;
		const max = clientSize - arrowDimensions[length] - maxPadding;
		const center = clientSize / 2 - arrowDimensions[length] / 2 + centerToReference;
		const offset = clamp$1(min, center, max);
		const shouldAddOffset = !middlewareData.arrow && getAlignment(placement) != null && center !== offset && rects.reference[length] / 2 - (center < min ? minPadding : maxPadding) - arrowDimensions[length] / 2 < 0;
		const alignmentOffset = shouldAddOffset ? center < min ? center - min : center - max : 0;
		return {
			[axis]: coords[axis] + alignmentOffset,
			data: {
				[axis]: offset,
				centerOffset: center - offset - alignmentOffset,
				...shouldAddOffset && { alignmentOffset }
			},
			reset: shouldAddOffset
		};
	}
});
/**
* Provides data to position an inner element of the floating element so that it
* appears centered to the reference element.
* This wraps the core `arrow` middleware to allow React refs as the element.
* @see https://floating-ui.com/docs/arrow
*/
var arrow = (options, deps) => ({
	...baseArrow(options),
	options: [options, deps]
});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/hideMiddleware.js
var hide = {
	name: "hide",
	async fn(state) {
		const { width, height, x, y } = state.rects.reference;
		const anchorHidden = width === 0 && height === 0 && x === 0 && y === 0;
		return { data: { referenceHidden: (await hide$1().fn(state)).data?.referenceHidden || anchorHidden } };
	}
};
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/adaptiveOriginMiddleware.js
var DEFAULT_SIDES = {
	sideX: "left",
	sideY: "top"
};
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/useAnchorPositioning.js
function getLogicalSide(sideParam, renderedSide, isRtl) {
	const isLogicalSideParam = sideParam === "inline-start" || sideParam === "inline-end";
	return {
		top: "top",
		right: isLogicalSideParam ? isRtl ? "inline-start" : "inline-end" : "right",
		bottom: "bottom",
		left: isLogicalSideParam ? isRtl ? "inline-end" : "inline-start" : "left"
	}[renderedSide];
}
function getOffsetData(state, sideParam, isRtl) {
	const { rects, placement } = state;
	return {
		side: getLogicalSide(sideParam, getSide(placement), isRtl),
		align: getAlignment(placement) || "center",
		anchor: {
			width: rects.reference.width,
			height: rects.reference.height
		},
		positioner: {
			width: rects.floating.width,
			height: rects.floating.height
		}
	};
}
/**
* Provides standardized anchor positioning behavior for floating elements. Wraps Floating UI's
* `useFloating` hook.
*/
function useAnchorPositioning(params) {
	const { anchor, positionMethod = "absolute", side: sideParam = "bottom", sideOffset = 0, align = "center", alignOffset = 0, collisionBoundary, collisionPadding: collisionPaddingParam = 5, sticky = false, arrowPadding = 5, disableAnchorTracking = false, keepMounted = false, floatingRootContext, mounted, collisionAvoidance, shiftCrossAxis = false, nodeId, adaptiveOrigin, lazyFlip = false, externalTree } = params;
	const [mountSide, setMountSide] = import_react.useState(null);
	if (!mounted && mountSide !== null) setMountSide(null);
	const collisionAvoidanceSide = collisionAvoidance.side || "flip";
	const collisionAvoidanceAlign = collisionAvoidance.align || "flip";
	const collisionAvoidanceFallbackAxisSide = collisionAvoidance.fallbackAxisSide || "end";
	const anchorFn = typeof anchor === "function" ? anchor : void 0;
	const anchorFnCallback = useStableCallback(anchorFn);
	const anchorDep = anchorFn ? anchorFnCallback : anchor;
	const anchorValueRef = useValueAsRef(anchor);
	const isRtl = useDirection() === "rtl";
	const side = mountSide || {
		top: "top",
		right: "right",
		bottom: "bottom",
		left: "left",
		"inline-end": isRtl ? "left" : "right",
		"inline-start": isRtl ? "right" : "left"
	}[sideParam];
	const placement = align === "center" ? side : `${side}-${align}`;
	let collisionPadding = collisionPaddingParam;
	const bias = 1;
	const biasTop = sideParam === "bottom" ? bias : 0;
	const biasBottom = sideParam === "top" ? bias : 0;
	const biasLeft = sideParam === "right" ? bias : 0;
	const biasRight = sideParam === "left" ? bias : 0;
	if (typeof collisionPadding === "number") collisionPadding = {
		top: collisionPadding + biasTop,
		right: collisionPadding + biasRight,
		bottom: collisionPadding + biasBottom,
		left: collisionPadding + biasLeft
	};
	else if (collisionPadding) collisionPadding = {
		top: (collisionPadding.top || 0) + biasTop,
		right: (collisionPadding.right || 0) + biasRight,
		bottom: (collisionPadding.bottom || 0) + biasBottom,
		left: (collisionPadding.left || 0) + biasLeft
	};
	const commonCollisionProps = {
		boundary: collisionBoundary === "clipping-ancestors" ? "clippingAncestors" : collisionBoundary,
		padding: collisionPadding
	};
	const arrowRef = import_react.useRef(null);
	const sideOffsetRef = useValueAsRef(sideOffset);
	const alignOffsetRef = useValueAsRef(alignOffset);
	const middleware = [offset((state) => {
		const data = getOffsetData(state, sideParam, isRtl);
		const sideAxis = typeof sideOffsetRef.current === "function" ? sideOffsetRef.current(data) : sideOffsetRef.current;
		const alignAxis = typeof alignOffsetRef.current === "function" ? alignOffsetRef.current(data) : alignOffsetRef.current;
		return {
			mainAxis: sideAxis,
			crossAxis: alignAxis,
			alignmentAxis: alignAxis
		};
	}, [
		typeof sideOffset !== "function" ? sideOffset : 0,
		typeof alignOffset !== "function" ? alignOffset : 0,
		isRtl,
		sideParam
	])];
	const shiftDisabled = collisionAvoidanceAlign === "none" && collisionAvoidanceSide !== "shift";
	const crossAxisShiftEnabled = !shiftDisabled && (sticky || shiftCrossAxis || collisionAvoidanceSide === "shift");
	const flipMiddleware = collisionAvoidanceSide === "none" ? null : flip({
		...commonCollisionProps,
		padding: {
			top: collisionPadding.top + bias,
			right: collisionPadding.right + bias,
			bottom: collisionPadding.bottom + bias,
			left: collisionPadding.left + bias
		},
		mainAxis: !shiftCrossAxis && collisionAvoidanceSide === "flip",
		crossAxis: collisionAvoidanceAlign === "flip" ? "alignment" : false,
		fallbackAxisSideDirection: collisionAvoidanceFallbackAxisSide
	});
	const shiftMiddleware = shiftDisabled ? null : shift((data) => {
		const html = ownerDocument(data.elements.floating).documentElement;
		return {
			...commonCollisionProps,
			rootBoundary: shiftCrossAxis ? {
				x: 0,
				y: 0,
				width: html.clientWidth,
				height: html.clientHeight
			} : void 0,
			mainAxis: collisionAvoidanceAlign !== "none",
			crossAxis: crossAxisShiftEnabled,
			limiter: sticky || shiftCrossAxis ? void 0 : limitShift((limitData) => {
				if (!arrowRef.current) return {};
				const { width, height } = arrowRef.current.getBoundingClientRect();
				const sideAxis = getSideAxis(getSide(limitData.placement));
				const arrowSize = sideAxis === "y" ? width : height;
				const offsetAmount = sideAxis === "y" ? collisionPadding.left + collisionPadding.right : collisionPadding.top + collisionPadding.bottom;
				return { offset: arrowSize / 2 + offsetAmount / 2 };
			})
		};
	}, [
		commonCollisionProps,
		sticky,
		shiftCrossAxis,
		collisionPadding,
		collisionAvoidanceAlign
	]);
	if (collisionAvoidanceSide === "shift" || collisionAvoidanceAlign === "shift" || align === "center") middleware.push(shiftMiddleware, flipMiddleware);
	else middleware.push(flipMiddleware, shiftMiddleware);
	middleware.push(size({
		...commonCollisionProps,
		apply({ elements: { floating }, availableWidth, availableHeight, rects }) {
			const floatingStyle = floating.style;
			floatingStyle.setProperty("--available-width", `${availableWidth}px`);
			floatingStyle.setProperty("--available-height", `${availableHeight}px`);
			const dpr = window.devicePixelRatio || 1;
			const { x, y, width, height } = rects.reference;
			const anchorWidth = (Math.round((x + width) * dpr) - Math.round(x * dpr)) / dpr;
			const anchorHeight = (Math.round((y + height) * dpr) - Math.round(y * dpr)) / dpr;
			floatingStyle.setProperty("--anchor-width", `${anchorWidth}px`);
			floatingStyle.setProperty("--anchor-height", `${anchorHeight}px`);
		}
	}), arrow(() => ({
		element: arrowRef.current || document.createElement("div"),
		padding: arrowPadding,
		offsetParent: "floating"
	}), [arrowPadding]), {
		name: "transformOrigin",
		fn(state) {
			const { elements, middlewareData, placement: renderedPlacement, rects, y } = state;
			const currentRenderedSide = getSide(renderedPlacement);
			const currentRenderedAxis = getSideAxis(currentRenderedSide);
			const arrowEl = arrowRef.current;
			const arrowX = middlewareData.arrow?.x || 0;
			const arrowY = middlewareData.arrow?.y || 0;
			const arrowWidth = arrowEl?.clientWidth || 0;
			const arrowHeight = arrowEl?.clientHeight || 0;
			const transformX = arrowX + arrowWidth / 2;
			const transformY = arrowY + arrowHeight / 2;
			const shiftY = Math.abs(middlewareData.shift?.y || 0);
			const halfAnchorHeight = rects.reference.height / 2;
			const sideOffsetValue = typeof sideOffset === "function" ? sideOffset(getOffsetData(state, sideParam, isRtl)) : sideOffset;
			const isOverlappingAnchor = shiftY > sideOffsetValue;
			const adjacentTransformOrigin = {
				top: `${transformX}px calc(100% + ${sideOffsetValue}px)`,
				bottom: `${transformX}px ${-sideOffsetValue}px`,
				left: `calc(100% + ${sideOffsetValue}px) ${transformY}px`,
				right: `${-sideOffsetValue}px ${transformY}px`
			}[currentRenderedSide];
			const overlapTransformOrigin = `${transformX}px ${rects.reference.y + halfAnchorHeight - y}px`;
			elements.floating.style.setProperty("--transform-origin", crossAxisShiftEnabled && currentRenderedAxis === "y" && isOverlappingAnchor ? overlapTransformOrigin : adjacentTransformOrigin);
			return {};
		}
	}, hide, adaptiveOrigin);
	useIsoLayoutEffect(() => {
		if (!mounted && floatingRootContext) floatingRootContext.update({
			referenceElement: null,
			floatingElement: null,
			domReferenceElement: null
		});
	}, [mounted, floatingRootContext]);
	const autoUpdateOptions = import_react.useMemo(() => ({
		elementResize: !disableAnchorTracking && typeof ResizeObserver !== "undefined",
		layoutShift: !disableAnchorTracking && typeof IntersectionObserver !== "undefined"
	}), [disableAnchorTracking]);
	const { refs, elements, x, y, middlewareData, update, placement: renderedPlacement, context, isPositioned, floatingStyles: originalFloatingStyles } = useFloating({
		rootContext: floatingRootContext,
		placement,
		middleware,
		strategy: positionMethod,
		whileElementsMounted: keepMounted ? void 0 : (...args) => autoUpdate(...args, autoUpdateOptions),
		nodeId,
		externalTree
	});
	const { sideX, sideY } = middlewareData.adaptiveOrigin || DEFAULT_SIDES;
	const resolvedPosition = isPositioned ? positionMethod : "fixed";
	const floatingStyles = import_react.useMemo(() => {
		const base = adaptiveOrigin ? {
			position: resolvedPosition,
			[sideX]: x,
			[sideY]: y
		} : {
			position: resolvedPosition,
			...originalFloatingStyles
		};
		if (!isPositioned) base.opacity = 0;
		return base;
	}, [
		adaptiveOrigin,
		resolvedPosition,
		sideX,
		x,
		sideY,
		y,
		originalFloatingStyles,
		isPositioned
	]);
	const registeredPositionReferenceRef = import_react.useRef(null);
	useIsoLayoutEffect(() => {
		if (!mounted) return;
		const anchorValue = anchorValueRef.current;
		const resolvedAnchor = typeof anchorValue === "function" ? anchorValue() : anchorValue;
		const finalAnchor = (isRef(resolvedAnchor) ? resolvedAnchor.current : resolvedAnchor) || null;
		if (finalAnchor !== registeredPositionReferenceRef.current) {
			refs.setPositionReference(finalAnchor);
			registeredPositionReferenceRef.current = finalAnchor;
		}
	}, [
		mounted,
		refs,
		anchorDep,
		anchorValueRef
	]);
	import_react.useEffect(() => {
		if (!mounted) return;
		const anchorValue = anchorValueRef.current;
		if (typeof anchorValue === "function") return;
		if (isRef(anchorValue) && anchorValue.current !== registeredPositionReferenceRef.current) {
			refs.setPositionReference(anchorValue.current);
			registeredPositionReferenceRef.current = anchorValue.current;
		}
	}, [
		mounted,
		refs,
		anchorDep,
		anchorValueRef
	]);
	import_react.useEffect(() => {
		if (keepMounted && mounted && elements.domReference && elements.floating) return autoUpdate(elements.domReference, elements.floating, update, autoUpdateOptions);
	}, [
		keepMounted,
		mounted,
		elements,
		update,
		autoUpdateOptions
	]);
	const renderedSide = getSide(renderedPlacement);
	const logicalRenderedSide = getLogicalSide(sideParam, renderedSide, isRtl);
	const renderedAlign = getAlignment(renderedPlacement) || "center";
	const anchorHidden = Boolean(middlewareData.hide?.referenceHidden);
	/**
	* Locks the flip (makes it "sticky") so it doesn't prefer a given placement
	* and flips back lazily, not eagerly. Ideal for filtered lists that change
	* the size of the popup dynamically to avoid unwanted flipping when typing.
	*/
	useIsoLayoutEffect(() => {
		if (lazyFlip && mounted && isPositioned) setMountSide(renderedSide);
	}, [
		lazyFlip,
		mounted,
		isPositioned,
		renderedSide
	]);
	const arrowStyles = import_react.useMemo(() => ({
		position: "absolute",
		top: middlewareData.arrow?.y,
		left: middlewareData.arrow?.x
	}), [middlewareData.arrow]);
	const arrowUncentered = middlewareData.arrow?.centerOffset !== 0;
	return import_react.useMemo(() => ({
		positionerStyles: floatingStyles,
		arrowStyles,
		arrowRef,
		arrowUncentered,
		side: logicalRenderedSide,
		align: renderedAlign,
		physicalSide: renderedSide,
		anchorHidden,
		refs,
		context,
		isPositioned,
		update
	}), [
		floatingStyles,
		arrowStyles,
		arrowRef,
		arrowUncentered,
		logicalRenderedSide,
		renderedAlign,
		renderedSide,
		anchorHidden,
		refs,
		context,
		isPositioned,
		update
	]);
}
function isRef(param) {
	return param != null && "current" in param;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/positioner/SelectPositionerContext.js
var SelectPositionerContext = /* @__PURE__ */ import_react.createContext(void 0);
function useSelectPositionerContext() {
	const context = import_react.useContext(SelectPositionerContext);
	if (!context) throw new Error(formatErrorMessage(59));
	return context;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/InternalBackdrop.js
/**
* @internal
*/
var InternalBackdrop = /* @__PURE__ */ import_react.forwardRef(function InternalBackdrop(props, ref) {
	const { cutout, ...otherProps } = props;
	let clipPath;
	if (cutout) {
		const rect = cutout?.getBoundingClientRect();
		clipPath = `polygon(
      0% 0%,
      100% 0%,
      100% 100%,
      0% 100%,
      0% 0%,
      ${rect.left}px ${rect.top}px,
      ${rect.left}px ${rect.bottom}px,
      ${rect.right}px ${rect.bottom}px,
      ${rect.right}px ${rect.top}px,
      ${rect.left}px ${rect.top}px
    )`;
	}
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("div", {
		ref,
		role: "presentation",
		"data-base-ui-inert": "",
		...otherProps,
		style: {
			position: "fixed",
			inset: 0,
			userSelect: "none",
			WebkitUserSelect: "none",
			clipPath
		}
	});
});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/getDisabledMountTransitionStyles.js
function getDisabledMountTransitionStyles(transitionStatus) {
	return transitionStatus === "starting" ? DISABLED_TRANSITIONS_STYLE : EMPTY_OBJECT;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/popup/utils.js
function clearStyles(element, originalStyles) {
	if (element) Object.assign(element.style, originalStyles);
}
var LIST_FUNCTIONAL_STYLES = {
	position: "relative",
	maxHeight: "100%",
	overflowX: "hidden",
	overflowY: "auto"
};
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/positioner/SelectPositioner.js
var FIXED = { position: "fixed" };
/**
* Positions the select popup.
* Renders a `<div>` element.
*
* Documentation: [Base UI Select](https://base-ui.com/react/components/select)
*/
var SelectPositioner = /* @__PURE__ */ import_react.forwardRef(function SelectPositioner(componentProps, forwardedRef) {
	const { anchor, positionMethod = "absolute", className, render, side = "bottom", align = "center", sideOffset = 0, alignOffset = 0, collisionBoundary = "clipping-ancestors", collisionPadding, arrowPadding = 5, sticky = false, disableAnchorTracking, alignItemWithTrigger = true, collisionAvoidance = DROPDOWN_COLLISION_AVOIDANCE, ...elementProps } = componentProps;
	const { store, listRef, labelsRef, alignItemWithTriggerActiveRef, selectedItemTextRef, valuesRef, initialValueRef, popupRef, setValue } = useSelectRootContext();
	const floatingRootContext = useSelectFloatingContext();
	const open = useStore(store, selectors.open);
	const mounted = useStore(store, selectors.mounted);
	const modal = useStore(store, selectors.modal);
	const value = useStore(store, selectors.value);
	const openMethod = useStore(store, selectors.openMethod);
	const positionerElement = useStore(store, selectors.positionerElement);
	const triggerElement = useStore(store, selectors.triggerElement);
	const isItemEqualToValue = useStore(store, selectors.isItemEqualToValue);
	const transitionStatus = useStore(store, selectors.transitionStatus);
	const scrollUpArrowRef = import_react.useRef(null);
	const scrollDownArrowRef = import_react.useRef(null);
	const [controlledAlignItemWithTrigger, setControlledAlignItemWithTrigger] = import_react.useState(alignItemWithTrigger);
	const alignItemWithTriggerActive = mounted && controlledAlignItemWithTrigger && openMethod !== "touch";
	if (!mounted && controlledAlignItemWithTrigger !== alignItemWithTrigger) setControlledAlignItemWithTrigger(alignItemWithTrigger);
	useIsoLayoutEffect(() => {
		if (!mounted) {
			if (selectors.scrollUpArrowVisible(store.state)) store.set("scrollUpArrowVisible", false);
			if (selectors.scrollDownArrowVisible(store.state)) store.set("scrollDownArrowVisible", false);
		}
	}, [store, mounted]);
	import_react.useImperativeHandle(alignItemWithTriggerActiveRef, () => alignItemWithTriggerActive);
	useScrollLock((alignItemWithTriggerActive || modal) && open && openMethod !== "touch", triggerElement);
	const positioning = useAnchorPositioning({
		anchor,
		floatingRootContext,
		positionMethod,
		mounted,
		side,
		sideOffset,
		align,
		alignOffset,
		arrowPadding,
		collisionBoundary,
		collisionPadding,
		sticky,
		disableAnchorTracking: disableAnchorTracking ?? alignItemWithTriggerActive,
		collisionAvoidance,
		keepMounted: true
	});
	const renderedSide = alignItemWithTriggerActive ? "none" : positioning.side;
	const positionerStyles = alignItemWithTriggerActive ? FIXED : positioning.positionerStyles;
	const defaultProps = import_react.useMemo(() => {
		const hiddenStyles = {};
		if (!open) hiddenStyles.pointerEvents = "none";
		return {
			role: "presentation",
			hidden: !mounted,
			style: {
				...positionerStyles,
				...hiddenStyles
			}
		};
	}, [
		open,
		mounted,
		positionerStyles
	]);
	const state = {
		open,
		side: renderedSide,
		align: positioning.align,
		anchorHidden: positioning.anchorHidden
	};
	const element = useRenderElement("div", componentProps, {
		ref: [forwardedRef, useStableCallback((element) => {
			store.set("positionerElement", element);
		})],
		state,
		stateAttributesMapping: popupStateMapping,
		props: [
			defaultProps,
			getDisabledMountTransitionStyles(transitionStatus),
			elementProps
		]
	});
	const prevMapSizeRef = import_react.useRef(0);
	const onMapChange = useStableCallback((map) => {
		if (map.size === 0 && prevMapSizeRef.current === 0) return;
		if (valuesRef.current.length === 0) return;
		const prevSize = prevMapSizeRef.current;
		prevMapSizeRef.current = map.size;
		if (map.size === prevSize) return;
		const eventDetails = createChangeEventDetails(none);
		if (prevSize !== 0 && !store.state.multiple && value !== null) {
			if (findItemIndex(valuesRef.current, value, isItemEqualToValue) === -1) {
				const initialSelectedValue = initialValueRef.current;
				const nextValue = initialSelectedValue != null && findItemIndex(valuesRef.current, initialSelectedValue, isItemEqualToValue) !== -1 ? initialSelectedValue : null;
				setValue(nextValue, eventDetails);
				if (nextValue === null) {
					store.set("selectedIndex", null);
					selectedItemTextRef.current = null;
				}
			}
		}
		if (prevSize !== 0 && store.state.multiple && Array.isArray(value)) {
			const hasVisibleItem = (selectedItemValue) => findItemIndex(valuesRef.current, selectedItemValue, isItemEqualToValue) !== -1;
			const nextValue = value.filter((selectedItemValue) => hasVisibleItem(selectedItemValue));
			if (nextValue.length !== value.length || nextValue.some((selectedItemValue) => !selectedValueIncludes(value, selectedItemValue, isItemEqualToValue))) {
				setValue(nextValue, eventDetails);
				if (nextValue.length === 0) {
					store.set("selectedIndex", null);
					selectedItemTextRef.current = null;
				}
			}
		}
		if (open && alignItemWithTriggerActive) {
			store.update({
				scrollUpArrowVisible: false,
				scrollDownArrowVisible: false
			});
			const stylesToClear = { height: "" };
			clearStyles(positionerElement, stylesToClear);
			clearStyles(popupRef.current, stylesToClear);
		}
	});
	const contextValue = import_react.useMemo(() => ({
		...positioning,
		side: renderedSide,
		alignItemWithTriggerActive,
		setControlledAlignItemWithTrigger,
		scrollUpArrowRef,
		scrollDownArrowRef
	}), [
		positioning,
		renderedSide,
		alignItemWithTriggerActive,
		setControlledAlignItemWithTrigger
	]);
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(CompositeList, {
		elementsRef: listRef,
		labelsRef,
		onMapChange,
		children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(SelectPositionerContext.Provider, {
			value: contextValue,
			children: [mounted && modal && /* @__PURE__ */ (0, import_jsx_runtime.jsx)(InternalBackdrop, {
				inert: inertValue(!open),
				cutout: triggerElement
			}), element]
		})
	});
});
//#endregion
//#region node_modules/.pnpm/@base-ui+utils@0.2.6_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/utils/esm/isMouseWithinBounds.js
function isMouseWithinBounds(event) {
	const targetRect = event.currentTarget.getBoundingClientRect();
	return targetRect.top + 1 <= event.clientY && event.clientY <= targetRect.bottom - 1 && targetRect.left + 1 <= event.clientX && event.clientX <= targetRect.right - 1;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/styles.js
var DISABLE_SCROLLBAR_CLASS_NAME = "base-ui-disable-scrollbar";
var styleDisableScrollbar = {
	className: DISABLE_SCROLLBAR_CLASS_NAME,
	getElement(nonce) {
		return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("style", {
			nonce,
			href: DISABLE_SCROLLBAR_CLASS_NAME,
			precedence: "base-ui:low",
			children: `.${DISABLE_SCROLLBAR_CLASS_NAME}{scrollbar-width:none}.${DISABLE_SCROLLBAR_CLASS_NAME}::-webkit-scrollbar{display:none}`
		});
	}
};
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/toolbar/root/ToolbarRootContext.js
var ToolbarRootContext = /* @__PURE__ */ import_react.createContext(void 0);
function useToolbarRootContext(optional) {
	const context = import_react.useContext(ToolbarRootContext);
	if (context === void 0 && !optional) throw new Error(formatErrorMessage(69));
	return context;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/utils/clamp.js
function clamp(val, min = Number.MIN_SAFE_INTEGER, max = Number.MAX_SAFE_INTEGER) {
	return Math.max(min, Math.min(val, max));
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/popup/SelectPopup.js
var SCROLL_EPS_PX = 1;
var stateAttributesMapping = {
	...popupStateMapping,
	...transitionStatusMapping
};
/**
* A container for the select list.
* Renders a `<div>` element.
*
* Documentation: [Base UI Select](https://base-ui.com/react/components/select)
*/
var SelectPopup = /* @__PURE__ */ import_react.forwardRef(function SelectPopup(componentProps, forwardedRef) {
	const { render, className, finalFocus, ...elementProps } = componentProps;
	const { store, popupRef, onOpenChangeComplete, setOpen, valueRef, selectedItemTextRef, keyboardActiveRef, multiple, handleScrollArrowVisibility, scrollHandlerRef, highlightItemOnHover } = useSelectRootContext();
	const { side, align, alignItemWithTriggerActive, setControlledAlignItemWithTrigger, scrollDownArrowRef, scrollUpArrowRef } = useSelectPositionerContext();
	const insideToolbar = useToolbarRootContext(true) != null;
	const floatingRootContext = useSelectFloatingContext();
	const { nonce, disableStyleElements } = useCSPContext();
	const highlightTimeout = useTimeout();
	const id = useStore(store, selectors.id);
	const open = useStore(store, selectors.open);
	const mounted = useStore(store, selectors.mounted);
	const popupProps = useStore(store, selectors.popupProps);
	const transitionStatus = useStore(store, selectors.transitionStatus);
	const triggerElement = useStore(store, selectors.triggerElement);
	const positionerElement = useStore(store, selectors.positionerElement);
	const listElement = useStore(store, selectors.listElement);
	const initialHeightRef = import_react.useRef(0);
	const reachedMaxHeightRef = import_react.useRef(false);
	const maxHeightRef = import_react.useRef(0);
	const initialPlacedRef = import_react.useRef(false);
	const originalPositionerStylesRef = import_react.useRef({});
	const scrollArrowFrame = useAnimationFrame();
	const handleScroll = useStableCallback((scroller) => {
		if (!positionerElement || !popupRef.current || !initialPlacedRef.current) return;
		if (reachedMaxHeightRef.current || !alignItemWithTriggerActive) {
			handleScrollArrowVisibility();
			return;
		}
		const isTopPositioned = positionerElement.style.top === "0px";
		const isBottomPositioned = positionerElement.style.bottom === "0px";
		const currentHeight = positionerElement.getBoundingClientRect().height;
		const doc = ownerDocument(positionerElement);
		const positionerStyles = getComputedStyle(positionerElement);
		const marginTop = parseFloat(positionerStyles.marginTop);
		const marginBottom = parseFloat(positionerStyles.marginBottom);
		const maxPopupHeight = getMaxPopupHeight(getComputedStyle(popupRef.current));
		const maxAvailableHeight = Math.min(doc.documentElement.clientHeight - marginTop - marginBottom, maxPopupHeight);
		const scrollTop = scroller.scrollTop;
		const maxScrollTop = getMaxScrollTop(scroller);
		let nextPositionerHeight = 0;
		let nextScrollTop = null;
		let setReachedMax = false;
		let scrollToMax = false;
		const setHeight = (height) => {
			positionerElement.style.height = `${height}px`;
		};
		const handleSmallDiff = (diff, targetScrollTop) => {
			const heightDelta = clamp(diff, 0, maxAvailableHeight - currentHeight);
			if (heightDelta > 0) setHeight(currentHeight + heightDelta);
			scroller.scrollTop = targetScrollTop;
			if (maxAvailableHeight - (currentHeight + heightDelta) <= SCROLL_EPS_PX) reachedMaxHeightRef.current = true;
			handleScrollArrowVisibility();
		};
		if (isTopPositioned) {
			const diff = maxScrollTop - scrollTop;
			const idealHeight = currentHeight + diff;
			const nextHeight = Math.min(idealHeight, maxAvailableHeight);
			nextPositionerHeight = nextHeight;
			if (diff <= SCROLL_EPS_PX) {
				handleSmallDiff(diff, maxScrollTop);
				return;
			}
			if (maxAvailableHeight - nextHeight > SCROLL_EPS_PX) scrollToMax = true;
			else setReachedMax = true;
		} else if (isBottomPositioned) {
			const diff = scrollTop;
			const idealHeight = currentHeight + diff;
			const nextHeight = Math.min(idealHeight, maxAvailableHeight);
			const overshoot = idealHeight - maxAvailableHeight;
			nextPositionerHeight = nextHeight;
			if (diff <= SCROLL_EPS_PX) {
				handleSmallDiff(diff, 0);
				return;
			}
			if (maxAvailableHeight - nextHeight > SCROLL_EPS_PX) nextScrollTop = 0;
			else {
				setReachedMax = true;
				if (scrollTop < maxScrollTop) nextScrollTop = scrollTop - (diff - overshoot);
			}
		}
		nextPositionerHeight = Math.ceil(nextPositionerHeight);
		if (nextPositionerHeight !== 0) setHeight(nextPositionerHeight);
		if (scrollToMax || nextScrollTop != null) {
			const nextMaxScrollTop = getMaxScrollTop(scroller);
			const target = scrollToMax ? nextMaxScrollTop : clamp(nextScrollTop, 0, nextMaxScrollTop);
			if (Math.abs(scroller.scrollTop - target) > SCROLL_EPS_PX) scroller.scrollTop = target;
		}
		if (setReachedMax || nextPositionerHeight >= maxAvailableHeight - SCROLL_EPS_PX) reachedMaxHeightRef.current = true;
		handleScrollArrowVisibility();
	});
	import_react.useImperativeHandle(scrollHandlerRef, () => handleScroll, [handleScroll]);
	useOpenChangeComplete({
		open,
		ref: popupRef,
		onComplete() {
			if (open) onOpenChangeComplete?.(true);
		}
	});
	const state = {
		open,
		transitionStatus,
		side,
		align
	};
	useIsoLayoutEffect(() => {
		if (!positionerElement || !popupRef.current || Object.keys(originalPositionerStylesRef.current).length) return;
		originalPositionerStylesRef.current = {
			top: positionerElement.style.top || "0",
			left: positionerElement.style.left || "0",
			right: positionerElement.style.right,
			height: positionerElement.style.height,
			bottom: positionerElement.style.bottom,
			minHeight: positionerElement.style.minHeight,
			maxHeight: positionerElement.style.maxHeight,
			marginTop: positionerElement.style.marginTop,
			marginBottom: positionerElement.style.marginBottom
		};
	}, [popupRef, positionerElement]);
	useIsoLayoutEffect(() => {
		if (open || alignItemWithTriggerActive) return;
		initialPlacedRef.current = false;
		reachedMaxHeightRef.current = false;
		initialHeightRef.current = 0;
		maxHeightRef.current = 0;
		clearStyles(positionerElement, originalPositionerStylesRef.current);
	}, [
		open,
		alignItemWithTriggerActive,
		positionerElement,
		popupRef
	]);
	useIsoLayoutEffect(() => {
		const popupElement = popupRef.current;
		if (!open || !triggerElement || !positionerElement || !popupElement || store.state.transitionStatus === "ending") return;
		if (!alignItemWithTriggerActive) {
			initialPlacedRef.current = true;
			scrollArrowFrame.request(handleScrollArrowVisibility);
			popupElement.style.removeProperty("--transform-origin");
			return;
		}
		queueMicrotask(() => {
			const restoreTransformStyles = unsetTransformStyles(popupElement);
			popupElement.style.removeProperty("--transform-origin");
			try {
				const positionerStyles = getComputedStyle(positionerElement);
				const popupStyles = getComputedStyle(popupElement);
				const doc = ownerDocument(triggerElement);
				const win = getWindow(positionerElement);
				const scale = getScale(triggerElement);
				const triggerRect = normalizeRect(triggerElement.getBoundingClientRect(), scale);
				const positionerRect = normalizeRect(positionerElement.getBoundingClientRect(), scale);
				const triggerX = triggerRect.left;
				const triggerHeight = triggerRect.height;
				const scroller = listElement || popupElement;
				const scrollHeight = scroller.scrollHeight;
				const borderBottom = parseFloat(popupStyles.borderBottomWidth);
				const marginTop = parseFloat(positionerStyles.marginTop) || 10;
				const marginBottom = parseFloat(positionerStyles.marginBottom) || 10;
				const minHeight = parseFloat(positionerStyles.minHeight) || 100;
				const maxPopupHeight = getMaxPopupHeight(popupStyles);
				const paddingLeft = 5;
				const paddingRight = 5;
				const triggerCollisionThreshold = 20;
				const viewportHeight = doc.documentElement.clientHeight - marginTop - marginBottom;
				const viewportWidth = doc.documentElement.clientWidth;
				const availableSpaceBeneathTrigger = viewportHeight - triggerRect.bottom + triggerHeight;
				const textElement = selectedItemTextRef.current;
				const valueElement = valueRef.current;
				let textRect;
				let offsetX = 0;
				let offsetY = 0;
				if (textElement && valueElement) {
					const valueRect = normalizeRect(valueElement.getBoundingClientRect(), scale);
					textRect = normalizeRect(textElement.getBoundingClientRect(), scale);
					const valueLeftFromTriggerLeft = valueRect.left - triggerX;
					const textLeftFromPositionerLeft = textRect.left - positionerRect.left;
					const valueCenterFromPositionerTop = valueRect.top - triggerRect.top + valueRect.height / 2;
					const textCenterFromTriggerTop = textRect.top - positionerRect.top + textRect.height / 2;
					offsetX = valueLeftFromTriggerLeft - textLeftFromPositionerLeft;
					offsetY = textCenterFromTriggerTop - valueCenterFromPositionerTop;
				}
				const idealHeight = availableSpaceBeneathTrigger + offsetY + marginBottom + borderBottom;
				let height = Math.min(viewportHeight, idealHeight);
				const maxHeight = viewportHeight - marginTop - marginBottom;
				const scrollTop = idealHeight - height;
				const left = Math.max(paddingLeft, triggerX + offsetX);
				const maxRight = viewportWidth - paddingRight;
				const rightOverflow = Math.max(0, left + positionerRect.width - maxRight);
				positionerElement.style.left = `${left - rightOverflow}px`;
				positionerElement.style.height = `${height}px`;
				positionerElement.style.maxHeight = "auto";
				positionerElement.style.marginTop = `${marginTop}px`;
				positionerElement.style.marginBottom = `${marginBottom}px`;
				popupElement.style.height = "100%";
				const maxScrollTop = getMaxScrollTop(scroller);
				const isTopPositioned = scrollTop >= maxScrollTop - SCROLL_EPS_PX;
				if (isTopPositioned) height = Math.min(viewportHeight, positionerRect.height) - (scrollTop - maxScrollTop);
				const fallbackToAlignPopupToTrigger = triggerRect.top < triggerCollisionThreshold || triggerRect.bottom > viewportHeight - triggerCollisionThreshold || Math.ceil(height) + SCROLL_EPS_PX < Math.min(scrollHeight, minHeight);
				const isPinchZoomed = (win.visualViewport?.scale ?? 1) !== 1 && isWebKit$1;
				if (fallbackToAlignPopupToTrigger || isPinchZoomed) {
					initialPlacedRef.current = true;
					clearStyles(positionerElement, originalPositionerStylesRef.current);
					import_react_dom.flushSync(() => setControlledAlignItemWithTrigger(false));
					return;
				}
				if (isTopPositioned) {
					const topOffset = Math.max(0, viewportHeight - idealHeight);
					positionerElement.style.top = positionerRect.height >= maxHeight ? "0" : `${topOffset}px`;
					positionerElement.style.height = `${height}px`;
					scroller.scrollTop = getMaxScrollTop(scroller);
					initialHeightRef.current = Math.max(minHeight, height);
				} else {
					positionerElement.style.bottom = "0";
					initialHeightRef.current = Math.max(minHeight, height);
					scroller.scrollTop = scrollTop;
				}
				if (textRect) {
					const popupTop = positionerRect.top;
					const popupHeight = positionerRect.height;
					const textCenterY = textRect.top + textRect.height / 2;
					const clampedY = clamp(popupHeight > 0 ? (textCenterY - popupTop) / popupHeight * 100 : 50, 0, 100);
					popupElement.style.setProperty("--transform-origin", `50% ${clampedY}%`);
				}
				if (initialHeightRef.current === viewportHeight || height >= maxPopupHeight) reachedMaxHeightRef.current = true;
				handleScrollArrowVisibility();
				setTimeout(() => {
					initialPlacedRef.current = true;
				});
			} finally {
				restoreTransformStyles();
			}
		});
	}, [
		store,
		open,
		positionerElement,
		triggerElement,
		valueRef,
		selectedItemTextRef,
		popupRef,
		handleScrollArrowVisibility,
		alignItemWithTriggerActive,
		setControlledAlignItemWithTrigger,
		scrollArrowFrame,
		scrollDownArrowRef,
		scrollUpArrowRef,
		listElement
	]);
	import_react.useEffect(() => {
		if (!alignItemWithTriggerActive || !positionerElement || !open) return;
		const win = getWindow(positionerElement);
		function handleResize(event) {
			setOpen(false, createChangeEventDetails(windowResize, event));
		}
		win.addEventListener("resize", handleResize);
		return () => {
			win.removeEventListener("resize", handleResize);
		};
	}, [
		setOpen,
		alignItemWithTriggerActive,
		positionerElement,
		open
	]);
	const defaultProps = {
		...listElement ? {
			role: "presentation",
			"aria-orientation": void 0
		} : {
			role: "listbox",
			"aria-multiselectable": multiple || void 0,
			id: `${id}-list`
		},
		onKeyDown(event) {
			keyboardActiveRef.current = true;
			if (insideToolbar && COMPOSITE_KEYS.has(event.key)) event.stopPropagation();
		},
		onMouseMove() {
			keyboardActiveRef.current = false;
		},
		onPointerLeave(event) {
			if (!highlightItemOnHover || isMouseWithinBounds(event) || event.pointerType === "touch") return;
			const popup = event.currentTarget;
			highlightTimeout.start(0, () => {
				store.set("activeIndex", null);
				popup.focus({ preventScroll: true });
			});
		},
		onScroll(event) {
			if (listElement) return;
			handleScroll(event.currentTarget);
		},
		...alignItemWithTriggerActive && { style: listElement ? { height: "100%" } : LIST_FUNCTIONAL_STYLES }
	};
	const element = useRenderElement("div", componentProps, {
		ref: [forwardedRef, popupRef],
		state,
		stateAttributesMapping,
		props: [
			popupProps,
			defaultProps,
			getDisabledMountTransitionStyles(transitionStatus),
			{ className: !listElement && alignItemWithTriggerActive ? styleDisableScrollbar.className : void 0 },
			elementProps
		]
	});
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(import_react.Fragment, { children: [!disableStyleElements && styleDisableScrollbar.getElement(nonce), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(FloatingFocusManager, {
		context: floatingRootContext,
		modal: false,
		disabled: !mounted,
		returnFocus: finalFocus,
		restoreFocus: true,
		children: element
	})] });
});
function getMaxPopupHeight(popupStyles) {
	const maxHeightStyle = popupStyles.maxHeight || "";
	return maxHeightStyle.endsWith("px") ? parseFloat(maxHeightStyle) || Infinity : Infinity;
}
function getMaxScrollTop(scroller) {
	return Math.max(0, scroller.scrollHeight - scroller.clientHeight);
}
function getScale(element) {
	return platform.getScale(element);
}
function normalizeRect(rect, scale) {
	return rectToClientRect({
		x: rect.x / scale.x,
		y: rect.y / scale.y,
		width: rect.width / scale.x,
		height: rect.height / scale.y
	});
}
var TRANSFORM_STYLE_RESETS = [
	["transform", "none"],
	["scale", "1"],
	["translate", "0 0"]
];
function unsetTransformStyles(popupElement) {
	const { style } = popupElement;
	const originalStyles = {};
	for (const [property, value] of TRANSFORM_STYLE_RESETS) {
		originalStyles[property] = style.getPropertyValue(property);
		style.setProperty(property, value, "important");
	}
	return () => {
		for (const [property] of TRANSFORM_STYLE_RESETS) {
			const originalValue = originalStyles[property];
			if (originalValue) style.setProperty(property, originalValue);
			else style.removeProperty(property);
		}
	};
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/list/SelectList.js
/**
* A container for the select items.
* Renders a `<div>` element.
*
* Documentation: [Base UI Select](https://base-ui.com/react/components/select)
*/
var SelectList = /* @__PURE__ */ import_react.forwardRef(function SelectList(componentProps, forwardedRef) {
	const { className, render, ...elementProps } = componentProps;
	const { store, scrollHandlerRef } = useSelectRootContext();
	const { alignItemWithTriggerActive } = useSelectPositionerContext();
	const hasScrollArrows = useStore(store, selectors.hasScrollArrows);
	const openMethod = useStore(store, selectors.openMethod);
	const multiple = useStore(store, selectors.multiple);
	const defaultProps = {
		id: `${useStore(store, selectors.id)}-list`,
		role: "listbox",
		"aria-multiselectable": multiple || void 0,
		onScroll(event) {
			scrollHandlerRef.current?.(event.currentTarget);
		},
		...alignItemWithTriggerActive && { style: LIST_FUNCTIONAL_STYLES },
		className: hasScrollArrows && openMethod !== "touch" ? styleDisableScrollbar.className : void 0
	};
	return useRenderElement("div", componentProps, {
		ref: [forwardedRef, useStableCallback((element) => {
			store.set("listElement", element);
		})],
		props: [defaultProps, elementProps]
	});
});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/item/SelectItemContext.js
var SelectItemContext = /* @__PURE__ */ import_react.createContext(void 0);
function useSelectItemContext() {
	const context = import_react.useContext(SelectItemContext);
	if (!context) throw new Error(formatErrorMessage(57));
	return context;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/item/SelectItem.js
/**
* An individual option in the select popup.
* Renders a `<div>` element.
*
* Documentation: [Base UI Select](https://base-ui.com/react/components/select)
*/
var SelectItem$1 = /* @__PURE__ */ import_react.memo(/* @__PURE__ */ import_react.forwardRef(function SelectItem(componentProps, forwardedRef) {
	const { render, className, value: itemValue = null, label, disabled = false, nativeButton = false, ...elementProps } = componentProps;
	const textRef = import_react.useRef(null);
	const listItem = useCompositeListItem({
		label,
		textRef,
		indexGuessBehavior: IndexGuessBehavior.GuessFromOrder
	});
	const { store, getItemProps, setOpen, setValue, selectionRef, typingRef, valuesRef, keyboardActiveRef, multiple, highlightItemOnHover } = useSelectRootContext();
	const highlightTimeout = useTimeout();
	const highlighted = useStore(store, selectors.isActive, listItem.index);
	const selected = useStore(store, selectors.isSelected, listItem.index, itemValue);
	const selectedByFocus = useStore(store, selectors.isSelectedByFocus, listItem.index);
	const isItemEqualToValue = useStore(store, selectors.isItemEqualToValue);
	const index = listItem.index;
	const hasRegistered = index !== -1;
	const itemRef = import_react.useRef(null);
	const indexRef = useValueAsRef(index);
	useIsoLayoutEffect(() => {
		if (!hasRegistered) return;
		const values = valuesRef.current;
		values[index] = itemValue;
		return () => {
			delete values[index];
		};
	}, [
		hasRegistered,
		index,
		itemValue,
		valuesRef
	]);
	useIsoLayoutEffect(() => {
		if (!hasRegistered) return;
		const selectedValue = store.state.value;
		let selectedCandidate = selectedValue;
		if (multiple && Array.isArray(selectedValue) && selectedValue.length > 0) selectedCandidate = selectedValue[selectedValue.length - 1];
		if (selectedCandidate !== void 0 && compareItemEquality(itemValue, selectedCandidate, isItemEqualToValue)) store.set("selectedIndex", index);
	}, [
		hasRegistered,
		index,
		multiple,
		isItemEqualToValue,
		store,
		itemValue
	]);
	const state = {
		disabled,
		selected,
		highlighted
	};
	const rootProps = getItemProps({
		active: highlighted,
		selected
	});
	rootProps.onFocus = void 0;
	rootProps.id = void 0;
	const lastKeyRef = import_react.useRef(null);
	const pointerTypeRef = import_react.useRef("mouse");
	const didPointerDownRef = import_react.useRef(false);
	const { getButtonProps, buttonRef } = useButton({
		disabled,
		focusableWhenDisabled: true,
		native: nativeButton,
		composite: true
	});
	function commitSelection(event) {
		const selectedValue = store.state.value;
		if (multiple) {
			const currentValue = Array.isArray(selectedValue) ? selectedValue : [];
			setValue(selected ? removeItem(currentValue, itemValue, isItemEqualToValue) : [...currentValue, itemValue], createChangeEventDetails(itemPress, event));
		} else {
			setValue(itemValue, createChangeEventDetails(itemPress, event));
			setOpen(false, createChangeEventDetails(itemPress, event));
		}
	}
	const defaultProps = {
		role: "option",
		"aria-selected": selected,
		tabIndex: highlighted ? 0 : -1,
		onFocus() {
			store.set("activeIndex", index);
		},
		onMouseEnter() {
			if (!keyboardActiveRef.current && store.state.selectedIndex === null && highlightItemOnHover) store.set("activeIndex", index);
		},
		onMouseMove() {
			if (highlightItemOnHover) store.set("activeIndex", index);
		},
		onMouseLeave(event) {
			if (!highlightItemOnHover || keyboardActiveRef.current || isMouseWithinBounds(event)) return;
			highlightTimeout.start(0, () => {
				if (store.state.activeIndex === index) store.set("activeIndex", null);
			});
		},
		onTouchStart() {
			selectionRef.current = {
				allowSelectedMouseUp: false,
				allowUnselectedMouseUp: false
			};
		},
		onKeyDown(event) {
			lastKeyRef.current = event.key;
			store.set("activeIndex", index);
			if (event.key === " " && typingRef.current) event.preventDefault();
		},
		onClick(event) {
			didPointerDownRef.current = false;
			if (event.type === "keydown" && lastKeyRef.current === null) return;
			if (disabled || event.type === "keydown" && lastKeyRef.current === " " && typingRef.current || pointerTypeRef.current !== "touch" && !highlighted) return;
			lastKeyRef.current = null;
			commitSelection(event.nativeEvent);
		},
		onPointerEnter(event) {
			pointerTypeRef.current = event.pointerType;
		},
		onPointerDown(event) {
			pointerTypeRef.current = event.pointerType;
			didPointerDownRef.current = true;
		},
		onMouseUp() {
			if (disabled) return;
			if (didPointerDownRef.current) {
				didPointerDownRef.current = false;
				return;
			}
			const disallowSelectedMouseUp = !selectionRef.current.allowSelectedMouseUp && selected;
			const disallowUnselectedMouseUp = !selectionRef.current.allowUnselectedMouseUp && !selected;
			if (disallowSelectedMouseUp || disallowUnselectedMouseUp || pointerTypeRef.current !== "touch" && !highlighted) return;
			itemRef.current?.click();
		}
	};
	const element = useRenderElement("div", componentProps, {
		ref: [
			buttonRef,
			forwardedRef,
			listItem.ref,
			itemRef
		],
		state,
		props: [
			rootProps,
			defaultProps,
			elementProps,
			getButtonProps
		]
	});
	const contextValue = import_react.useMemo(() => ({
		selected,
		indexRef,
		textRef,
		selectedByFocus,
		hasRegistered
	}), [
		selected,
		indexRef,
		textRef,
		selectedByFocus,
		hasRegistered
	]);
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectItemContext.Provider, {
		value: contextValue,
		children: element
	});
}));
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/item-indicator/SelectItemIndicator.js
/**
* Indicates whether the select item is selected.
* Renders a `<span>` element.
*
* Documentation: [Base UI Select](https://base-ui.com/react/components/select)
*/
var SelectItemIndicator = /* @__PURE__ */ import_react.forwardRef(function SelectItemIndicator(componentProps, forwardedRef) {
	const keepMounted = componentProps.keepMounted ?? false;
	const { selected } = useSelectItemContext();
	if (!(keepMounted || selected)) return null;
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Inner, {
		...componentProps,
		ref: forwardedRef
	});
});
var Inner = /* @__PURE__ */ import_react.memo(/* @__PURE__ */ import_react.forwardRef((componentProps, forwardedRef) => {
	const { render, className, keepMounted, ...elementProps } = componentProps;
	const { selected } = useSelectItemContext();
	const indicatorRef = import_react.useRef(null);
	const { transitionStatus, setMounted } = useTransitionStatus(selected);
	const element = useRenderElement("span", componentProps, {
		ref: [forwardedRef, indicatorRef],
		state: {
			selected,
			transitionStatus
		},
		props: [{
			"aria-hidden": true,
			children: "✔️"
		}, elementProps],
		stateAttributesMapping: transitionStatusMapping
	});
	useOpenChangeComplete({
		open: selected,
		ref: indicatorRef,
		onComplete() {
			if (!selected) setMounted(false);
		}
	});
	return element;
}));
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/item-text/SelectItemText.js
/**
* A text label of the select item.
* Renders a `<div>` element.
*
* Documentation: [Base UI Select](https://base-ui.com/react/components/select)
*/
var SelectItemText = /* @__PURE__ */ import_react.memo(/* @__PURE__ */ import_react.forwardRef(function SelectItemText(componentProps, forwardedRef) {
	const { indexRef, textRef, selectedByFocus, hasRegistered } = useSelectItemContext();
	const { selectedItemTextRef } = useSelectRootContext();
	const { className, render, ...elementProps } = componentProps;
	return useRenderElement("div", componentProps, {
		ref: [
			import_react.useCallback((node) => {
				if (!node || !hasRegistered) return;
				const hasNoSelectedItemText = selectedItemTextRef.current === null || !selectedItemTextRef.current.isConnected;
				if (selectedByFocus || hasNoSelectedItemText && indexRef.current === 0) selectedItemTextRef.current = node;
			}, [
				selectedItemTextRef,
				indexRef,
				selectedByFocus,
				hasRegistered
			]),
			forwardedRef,
			textRef
		],
		props: elementProps
	});
}));
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/scroll-arrow/SelectScrollArrow.js
/**
* @internal
*/
var SelectScrollArrow = /* @__PURE__ */ import_react.forwardRef(function SelectScrollArrow(componentProps, forwardedRef) {
	const { render, className, direction, keepMounted = false, ...elementProps } = componentProps;
	const { store, popupRef, listRef, handleScrollArrowVisibility, scrollArrowsMountedCountRef } = useSelectRootContext();
	const { side, scrollDownArrowRef, scrollUpArrowRef } = useSelectPositionerContext();
	const stateVisible = useStore(store, direction === "up" ? selectors.scrollUpArrowVisible : selectors.scrollDownArrowVisible);
	const openMethod = useStore(store, selectors.openMethod);
	const visible = stateVisible && openMethod !== "touch";
	const timeout = useTimeout();
	const scrollArrowRef = direction === "up" ? scrollUpArrowRef : scrollDownArrowRef;
	const { transitionStatus, setMounted } = useTransitionStatus(visible);
	useIsoLayoutEffect(() => {
		scrollArrowsMountedCountRef.current += 1;
		if (!store.state.hasScrollArrows) store.set("hasScrollArrows", true);
		return () => {
			scrollArrowsMountedCountRef.current = Math.max(0, scrollArrowsMountedCountRef.current - 1);
			if (scrollArrowsMountedCountRef.current === 0 && store.state.hasScrollArrows) store.set("hasScrollArrows", false);
		};
	}, [store, scrollArrowsMountedCountRef]);
	useOpenChangeComplete({
		open: visible,
		ref: scrollArrowRef,
		onComplete() {
			if (!visible) setMounted(false);
		}
	});
	const element = useRenderElement("div", componentProps, {
		ref: [forwardedRef, scrollArrowRef],
		state: {
			direction,
			visible,
			side,
			transitionStatus
		},
		props: [{
			"aria-hidden": true,
			children: direction === "up" ? "▲" : "▼",
			style: { position: "absolute" },
			onMouseMove(event) {
				if (event.movementX === 0 && event.movementY === 0 || timeout.isStarted()) return;
				store.set("activeIndex", null);
				function scrollNextItem() {
					const scroller = store.state.listElement ?? popupRef.current;
					if (!scroller) return;
					store.set("activeIndex", null);
					handleScrollArrowVisibility();
					const isScrolledToTop = scroller.scrollTop === 0;
					const isScrolledToBottom = Math.round(scroller.scrollTop + scroller.clientHeight) >= scroller.scrollHeight;
					if (listRef.current.length === 0) if (direction === "up") store.set("scrollUpArrowVisible", !isScrolledToTop);
					else store.set("scrollDownArrowVisible", !isScrolledToBottom);
					if (direction === "up" && isScrolledToTop || direction === "down" && isScrolledToBottom) {
						timeout.clear();
						return;
					}
					if ((store.state.listElement || popupRef.current) && listRef.current && listRef.current.length > 0) {
						const items = listRef.current;
						const scrollArrowHeight = scrollArrowRef.current?.offsetHeight || 0;
						if (direction === "up") {
							let firstVisibleIndex = 0;
							const scrollTop = scroller.scrollTop + scrollArrowHeight;
							for (let i = 0; i < items.length; i += 1) {
								const item = items[i];
								if (item) {
									if (item.offsetTop >= scrollTop) {
										firstVisibleIndex = i;
										break;
									}
								}
							}
							const targetIndex = Math.max(0, firstVisibleIndex - 1);
							if (targetIndex < firstVisibleIndex) {
								const targetItem = items[targetIndex];
								if (targetItem) scroller.scrollTop = Math.max(0, targetItem.offsetTop - scrollArrowHeight);
							} else scroller.scrollTop = 0;
						} else {
							let lastVisibleIndex = items.length - 1;
							const scrollBottom = scroller.scrollTop + scroller.clientHeight - scrollArrowHeight;
							for (let i = 0; i < items.length; i += 1) {
								const item = items[i];
								if (item) {
									if (item.offsetTop + item.offsetHeight > scrollBottom) {
										lastVisibleIndex = Math.max(0, i - 1);
										break;
									}
								}
							}
							const targetIndex = Math.min(items.length - 1, lastVisibleIndex + 1);
							if (targetIndex > lastVisibleIndex) {
								const targetItem = items[targetIndex];
								if (targetItem) scroller.scrollTop = targetItem.offsetTop + targetItem.offsetHeight - scroller.clientHeight + scrollArrowHeight;
							} else scroller.scrollTop = scroller.scrollHeight - scroller.clientHeight;
						}
					}
					timeout.start(40, scrollNextItem);
				}
				timeout.start(40, scrollNextItem);
			},
			onMouseLeave() {
				timeout.clear();
			}
		}, elementProps]
	});
	if (!(visible || keepMounted)) return null;
	return element;
});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/scroll-down-arrow/SelectScrollDownArrow.js
/**
* An element that scrolls the select popup down when hovered. Does not render when using touch input.
* Renders a `<div>` element.
*
* Documentation: [Base UI Select](https://base-ui.com/react/components/select)
*/
var SelectScrollDownArrow = /* @__PURE__ */ import_react.forwardRef(function SelectScrollDownArrow(props, forwardedRef) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectScrollArrow, {
		...props,
		ref: forwardedRef,
		direction: "down"
	});
});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/select/scroll-up-arrow/SelectScrollUpArrow.js
/**
* An element that scrolls the select popup up when hovered. Does not render when using touch input.
* Renders a `<div>` element.
*
* Documentation: [Base UI Select](https://base-ui.com/react/components/select)
*/
var SelectScrollUpArrow = /* @__PURE__ */ import_react.forwardRef(function SelectScrollUpArrow(props, forwardedRef) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectScrollArrow, {
		...props,
		ref: forwardedRef,
		direction: "up"
	});
});
//#endregion
//#region src/components/ui/select.tsx
var Select = SelectRoot;
function SelectValue({ className, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectValue$1, {
		"data-slot": "select-value",
		className: cn("flex flex-1 text-left", className),
		...props
	});
}
function SelectTrigger({ className, size = "default", children, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(SelectTrigger$1, {
		"data-slot": "select-trigger",
		"data-size": size,
		className: cn("flex w-fit items-center justify-between gap-1.5 rounded-lg border border-input bg-transparent py-2 pr-2 pl-2.5 text-sm whitespace-nowrap transition-colors outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 data-placeholder:text-muted-foreground data-[size=default]:h-8 data-[size=sm]:h-7 data-[size=sm]:rounded-[min(var(--radius-md),10px)] *:data-[slot=select-value]:line-clamp-1 *:data-[slot=select-value]:flex *:data-[slot=select-value]:items-center *:data-[slot=select-value]:gap-1.5 dark:bg-input/30 dark:hover:bg-input/50 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4", className),
		...props,
		children: [children, /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectIcon, { render: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(ChevronDown, { className: "pointer-events-none size-4 text-muted-foreground" }) })]
	});
}
function SelectContent({ className, children, side = "bottom", sideOffset = 4, align = "center", alignOffset = 0, alignItemWithTrigger = true, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectPortal, { children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectPositioner, {
		side,
		sideOffset,
		align,
		alignOffset,
		alignItemWithTrigger,
		className: "isolate z-50",
		children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(SelectPopup, {
			"data-slot": "select-content",
			"data-align-trigger": alignItemWithTrigger,
			className: cn("relative isolate z-50 max-h-(--available-height) w-(--anchor-width) min-w-36 origin-(--transform-origin) overflow-x-hidden overflow-y-auto rounded-lg bg-popover text-popover-foreground shadow-md ring-1 ring-foreground/10 duration-100 data-[align-trigger=true]:animate-none data-[side=bottom]:slide-in-from-top-2 data-[side=inline-end]:slide-in-from-left-2 data-[side=inline-start]:slide-in-from-right-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95", className),
			...props,
			children: [
				/* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectScrollUpButton, {}),
				/* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectList, { children }),
				/* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectScrollDownButton, {})
			]
		})
	}) });
}
function SelectItem({ className, children, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(SelectItem$1, {
		"data-slot": "select-item",
		className: cn("relative flex w-full cursor-default items-center gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground not-data-[variant=destructive]:focus:**:text-accent-foreground data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 *:[span]:last:flex *:[span]:last:items-center *:[span]:last:gap-2", className),
		...props,
		children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectItemText, {
			className: "flex flex-1 shrink-0 gap-2 whitespace-nowrap",
			children
		}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectItemIndicator, {
			render: /* @__PURE__ */ (0, import_jsx_runtime.jsx)("span", { className: "pointer-events-none absolute right-2 flex size-4 items-center justify-center" }),
			children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Check, { className: "pointer-events-none" })
		})]
	});
}
function SelectScrollUpButton({ className, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectScrollUpArrow, {
		"data-slot": "select-scroll-up-button",
		className: cn("top-0 z-10 flex w-full cursor-default items-center justify-center bg-popover py-1 [&_svg:not([class*='size-'])]:size-4", className),
		...props,
		children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(ChevronUp, {})
	});
}
function SelectScrollDownButton({ className, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectScrollDownArrow, {
		"data-slot": "select-scroll-down-button",
		className: cn("bottom-0 z-10 flex w-full cursor-default items-center justify-center bg-popover py-1 [&_svg:not([class*='size-'])]:size-4", className),
		...props,
		children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(ChevronDown, {})
	});
}
//#endregion
//#region src/routes/scripts/index.tsx?tsr-split=component
function BrowsePage() {
	const { scripts, sites } = Route.useLoaderData();
	const search = Route.useSearch();
	const navigate = useNavigate();
	const [query, setQuery] = (0, import_react.useState)(search.q ?? "");
	(0, import_react.useEffect)(() => {
		setQuery(search.q ?? "");
	}, [search.q]);
	const updateSearch = (0, import_react.useCallback)((updates) => {
		navigate({
			to: "/scripts",
			search: (prev) => ({
				...prev,
				...updates
			}),
			replace: true
		});
	}, [navigate]);
	(0, import_react.useEffect)(() => {
		const timer = setTimeout(() => {
			const trimmed = query.trim();
			if (trimmed !== (search.q ?? "")) updateSearch({ q: trimmed || void 0 });
		}, 300);
		return () => clearTimeout(timer);
	}, [
		query,
		search.q,
		updateSearch
	]);
	const currentSort = search.sort || "popular";
	const currentSite = search.site || "";
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("main", {
		className: "mx-auto max-w-6xl px-4 py-8",
		children: [
			/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "mb-8",
				children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)("h1", {
					className: "text-2xl font-bold",
					children: "Browse Scripts"
				}), /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("p", {
					className: "text-muted-foreground",
					children: [
						scripts.length,
						" script",
						scripts.length !== 1 ? "s" : "",
						" found"
					]
				})]
			}),
			/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "mb-6 flex flex-wrap items-center gap-3",
				children: [
					/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
						className: "relative flex-1 min-w-[200px] max-w-sm",
						children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Search, { className: "absolute left-2.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" }), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Input, {
							type: "search",
							placeholder: "Search scripts…",
							value: query,
							onChange: (e) => setQuery(e.target.value),
							className: "pl-8"
						})]
					}),
					/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Select, {
						value: currentSort,
						onValueChange: (val) => updateSearch({ sort: val || void 0 }),
						children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectTrigger, {
							className: "w-[140px]",
							children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectValue, { placeholder: "Sort by" })
						}), /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(SelectContent, { children: [
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectItem, {
								value: "popular",
								children: "Popular"
							}),
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectItem, {
								value: "newest",
								children: "Newest"
							}),
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectItem, {
								value: "name",
								children: "Name"
							})
						] })]
					}),
					/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Select, {
						value: currentSite,
						onValueChange: (val) => updateSearch({ site: val || void 0 }),
						children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectTrigger, {
							className: "w-[160px]",
							children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectValue, { placeholder: "All sites" })
						}), /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(SelectContent, { children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectItem, {
							value: "",
							children: "All sites"
						}), sites.map((s) => /* @__PURE__ */ (0, import_jsx_runtime.jsx)(SelectItem, {
							value: s,
							children: s
						}, s))] })]
					}),
					(search.q || search.site) && /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Button, {
						variant: "ghost",
						size: "sm",
						onClick: () => navigate({
							to: "/scripts",
							search: { sort: search.sort }
						}),
						children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(X, { className: "size-3.5" }), "Clear filters"]
					})
				]
			}),
			scripts.length === 0 ? /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "py-16 text-center text-muted-foreground",
				children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)("p", {
					className: "text-lg",
					children: "No scripts found"
				}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("p", {
					className: "text-sm",
					children: "Try adjusting your search or filters."
				})]
			}) : /* @__PURE__ */ (0, import_jsx_runtime.jsx)("div", {
				className: "grid gap-4 sm:grid-cols-2 lg:grid-cols-3",
				children: scripts.map((script) => /* @__PURE__ */ (0, import_jsx_runtime.jsx)(ScriptCard, { script }, script.name))
			})
		]
	});
}
//#endregion
export { BrowsePage as component };
