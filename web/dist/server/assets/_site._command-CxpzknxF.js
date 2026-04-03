import { Y as __toESM, l as require_jsx_runtime, q as require_react } from "./react-dom-BccjEFPc.js";
import { H as Link, _ as formatErrorMessage, g as useCompositeRootContext, h as CompositeRootContext, i as useRenderElement, l as EMPTY_ARRAY$1, m as useButton, n as cn, p as useMergedRefs, r as cva, t as Button, u as EMPTY_OBJECT, x as useStableCallback, y as useIsoLayoutEffect } from "./button-BAIw3jLv.js";
import { i as createLucideIcon } from "./server-fns-BWiKIn9S.js";
import { A as getMaxListIndex, Bt as ownerDocument, C as inertValue, D as getGridCellIndexOfCorner, E as findNonDisabledListIndex, Ht as CompositeList, Lt as useCompositeListItem, N as isIndexOutOfListBounds, O as getGridCellIndices, P as isListIndexDisabled, Rt as useBaseUiId, S as transitionStatusMapping, T as createGridCellMap, Tt as createChangeEventDetails, Ut as useControlled, Vt as useDirection, Wt as Check, _ as useOpenChangeComplete, a as ARROW_RIGHT, ct as contains, d as MODIFIER_KEYS, f as VERTICAL_KEYS, g as useTransitionStatus, h as scrollIntoViewIfNeeded, i as ARROW_LEFT, j as getMinListIndex, jt as none, k as getGridNavigatedIndex, l as HORIZONTAL_KEYS, m as isNativeInput, n as ARROW_DOWN, o as ARROW_UP, p as VERTICAL_KEYS_WITH_EXTRA_KEYS, r as ARROW_KEYS, st as activeElement, t as ALL_KEYS, u as HORIZONTAL_KEYS_WITH_EXTRA_KEYS, x as TransitionStatusDataAttributes } from "./composite-CSdsaCvZ.js";
import { c as TrendingUp, n as Card, o as CardHeader, r as CardContent, s as CardTitle, t as Badge } from "./badge-B7SmroTr.js";
import { t as Route } from "./_site._command-wEfuCnhP.js";
var ArrowLeft = createLucideIcon("arrow-left", [["path", {
	d: "m12 19-7-7 7-7",
	key: "1l729n"
}], ["path", {
	d: "M19 12H5",
	key: "x3x0zl"
}]]);
var Calendar = createLucideIcon("calendar", [
	["path", {
		d: "M8 2v4",
		key: "1cmpym"
	}],
	["path", {
		d: "M16 2v4",
		key: "4m81vk"
	}],
	["rect", {
		width: "18",
		height: "18",
		x: "3",
		y: "4",
		rx: "2",
		key: "1hopcy"
	}],
	["path", {
		d: "M3 10h18",
		key: "8toen8"
	}]
]);
var Copy = createLucideIcon("copy", [["rect", {
	width: "14",
	height: "14",
	x: "8",
	y: "8",
	rx: "2",
	ry: "2",
	key: "17jyea"
}], ["path", {
	d: "M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2",
	key: "zix9uf"
}]]);
var Globe = createLucideIcon("globe", [
	["circle", {
		cx: "12",
		cy: "12",
		r: "10",
		key: "1mglay"
	}],
	["path", {
		d: "M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20",
		key: "13o1zl"
	}],
	["path", {
		d: "M2 12h20",
		key: "9i4pu4"
	}]
]);
var LockOpen = createLucideIcon("lock-open", [["rect", {
	width: "18",
	height: "11",
	x: "3",
	y: "11",
	rx: "2",
	ry: "2",
	key: "1w4ew1"
}], ["path", {
	d: "M7 11V7a5 5 0 0 1 9.9-1",
	key: "1mm8w8"
}]]);
var Lock = createLucideIcon("lock", [["rect", {
	width: "18",
	height: "11",
	x: "3",
	y: "11",
	rx: "2",
	ry: "2",
	key: "1w4ew1"
}], ["path", {
	d: "M7 11V7a5 5 0 0 1 10 0v4",
	key: "fwvmzm"
}]]);
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/tabs/root/TabsRootContext.js
var import_react = /* @__PURE__ */ __toESM(require_react());
/**
* @internal
*/
var TabsRootContext = /* @__PURE__ */ import_react.createContext(void 0);
function useTabsRootContext() {
	const context = import_react.useContext(TabsRootContext);
	if (context === void 0) throw new Error(formatErrorMessage(64));
	return context;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/tabs/root/TabsRootDataAttributes.js
var TabsRootDataAttributes = /* @__PURE__ */ function(TabsRootDataAttributes) {
	/**
	* Indicates the direction of the activation (based on the previous active tab).
	* @type {'left' | 'right' | 'up' | 'down' | 'none'}
	*/
	TabsRootDataAttributes["activationDirection"] = "data-activation-direction";
	/**
	* Indicates the orientation of the tabs.
	* @type {'horizontal' | 'vertical'}
	*/
	TabsRootDataAttributes["orientation"] = "data-orientation";
	return TabsRootDataAttributes;
}({});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/tabs/root/stateAttributesMapping.js
var tabsStateAttributesMapping = { tabActivationDirection: (dir) => ({ [TabsRootDataAttributes.activationDirection]: dir }) };
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/tabs/root/TabsRoot.js
var import_jsx_runtime = require_jsx_runtime();
/**
* Groups the tabs and the corresponding panels.
* Renders a `<div>` element.
*
* Documentation: [Base UI Tabs](https://base-ui.com/react/components/tabs)
*/
var TabsRoot = /* @__PURE__ */ import_react.forwardRef(function TabsRoot(componentProps, forwardedRef) {
	const { className, defaultValue: defaultValueProp = 0, onValueChange: onValueChangeProp, orientation = "horizontal", render, value: valueProp, ...elementProps } = componentProps;
	const direction = useDirection();
	const hasExplicitDefaultValueProp = Object.hasOwn(componentProps, "defaultValue");
	const tabPanelRefs = import_react.useRef([]);
	const [mountedTabPanels, setMountedTabPanels] = import_react.useState(() => /* @__PURE__ */ new Map());
	const [value, setValue] = useControlled({
		controlled: valueProp,
		default: defaultValueProp,
		name: "Tabs",
		state: "value"
	});
	const isControlled = valueProp !== void 0;
	const [tabMap, setTabMap] = import_react.useState(() => /* @__PURE__ */ new Map());
	const [tabActivationDirection, setTabActivationDirection] = import_react.useState("none");
	const onValueChange = useStableCallback((newValue, eventDetails) => {
		onValueChangeProp?.(newValue, eventDetails);
		if (eventDetails.isCanceled) return;
		setValue(newValue);
		setTabActivationDirection(eventDetails.activationDirection);
	});
	const registerMountedTabPanel = useStableCallback((panelValue, panelId) => {
		setMountedTabPanels((prev) => {
			if (prev.get(panelValue) === panelId) return prev;
			const next = new Map(prev);
			next.set(panelValue, panelId);
			return next;
		});
	});
	const unregisterMountedTabPanel = useStableCallback((panelValue, panelId) => {
		setMountedTabPanels((prev) => {
			if (!prev.has(panelValue) || prev.get(panelValue) !== panelId) return prev;
			const next = new Map(prev);
			next.delete(panelValue);
			return next;
		});
	});
	const getTabPanelIdByValue = import_react.useCallback((tabValue) => {
		return mountedTabPanels.get(tabValue);
	}, [mountedTabPanels]);
	const getTabIdByPanelValue = import_react.useCallback((tabPanelValue) => {
		for (const tabMetadata of tabMap.values()) if (tabPanelValue === tabMetadata?.value) return tabMetadata?.id;
	}, [tabMap]);
	const getTabElementBySelectedValue = import_react.useCallback((selectedValue) => {
		if (selectedValue === void 0) return null;
		for (const [tabElement, tabMetadata] of tabMap.entries()) if (tabMetadata != null && selectedValue === (tabMetadata.value ?? tabMetadata.index)) return tabElement;
		return null;
	}, [tabMap]);
	const tabsContextValue = import_react.useMemo(() => ({
		direction,
		getTabElementBySelectedValue,
		getTabIdByPanelValue,
		getTabPanelIdByValue,
		onValueChange,
		orientation,
		registerMountedTabPanel,
		setTabMap,
		unregisterMountedTabPanel,
		tabActivationDirection,
		value
	}), [
		direction,
		getTabElementBySelectedValue,
		getTabIdByPanelValue,
		getTabPanelIdByValue,
		onValueChange,
		orientation,
		registerMountedTabPanel,
		setTabMap,
		unregisterMountedTabPanel,
		tabActivationDirection,
		value
	]);
	const selectedTabMetadata = import_react.useMemo(() => {
		for (const tabMetadata of tabMap.values()) if (tabMetadata != null && tabMetadata.value === value) return tabMetadata;
	}, [tabMap, value]);
	const firstEnabledTabValue = import_react.useMemo(() => {
		for (const tabMetadata of tabMap.values()) if (tabMetadata != null && !tabMetadata.disabled) return tabMetadata.value;
	}, [tabMap]);
	useIsoLayoutEffect(() => {
		if (isControlled || tabMap.size === 0) return;
		const selectionIsDisabled = selectedTabMetadata?.disabled;
		const selectionIsMissing = selectedTabMetadata == null && value !== null;
		if (hasExplicitDefaultValueProp && selectionIsDisabled && value === defaultValueProp) return;
		if (!selectionIsDisabled && !selectionIsMissing) return;
		const fallbackValue = firstEnabledTabValue ?? null;
		if (value === fallbackValue) return;
		setValue(fallbackValue);
		setTabActivationDirection("none");
	}, [
		defaultValueProp,
		firstEnabledTabValue,
		hasExplicitDefaultValueProp,
		isControlled,
		selectedTabMetadata,
		setTabActivationDirection,
		setValue,
		tabMap,
		value
	]);
	const element = useRenderElement("div", componentProps, {
		state: {
			orientation,
			tabActivationDirection
		},
		ref: forwardedRef,
		props: elementProps,
		stateAttributesMapping: tabsStateAttributesMapping
	});
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(TabsRootContext.Provider, {
		value: tabsContextValue,
		children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(CompositeList, {
			elementsRef: tabPanelRefs,
			children: element
		})
	});
});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/composite/constants.js
var ACTIVE_COMPOSITE_ITEM = "data-composite-item-active";
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/composite/item/useCompositeItem.js
function useCompositeItem(params = {}) {
	const { highlightItemOnHover, highlightedIndex, onHighlightedIndexChange } = useCompositeRootContext();
	const { ref, index } = useCompositeListItem(params);
	const isHighlighted = highlightedIndex === index;
	const itemRef = import_react.useRef(null);
	const mergedRef = useMergedRefs(ref, itemRef);
	return {
		compositeProps: import_react.useMemo(() => ({
			tabIndex: isHighlighted ? 0 : -1,
			onFocus() {
				onHighlightedIndexChange(index);
			},
			onMouseMove() {
				const item = itemRef.current;
				if (!highlightItemOnHover || !item) return;
				const disabled = item.hasAttribute("disabled") || item.ariaDisabled === "true";
				if (!isHighlighted && !disabled) item.focus();
			}
		}), [
			isHighlighted,
			onHighlightedIndexChange,
			index,
			highlightItemOnHover
		]),
		compositeRef: mergedRef,
		index
	};
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/tabs/list/TabsListContext.js
var TabsListContext = /* @__PURE__ */ import_react.createContext(void 0);
function useTabsListContext() {
	const context = import_react.useContext(TabsListContext);
	if (context === void 0) throw new Error(formatErrorMessage(65));
	return context;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/tabs/tab/TabsTab.js
/**
* An individual interactive tab button that toggles the corresponding panel.
* Renders a `<button>` element.
*
* Documentation: [Base UI Tabs](https://base-ui.com/react/components/tabs)
*/
var TabsTab = /* @__PURE__ */ import_react.forwardRef(function TabsTab(componentProps, forwardedRef) {
	const { className, disabled = false, render, value, id: idProp, nativeButton = true, ...elementProps } = componentProps;
	const { value: activeTabValue, getTabPanelIdByValue, orientation } = useTabsRootContext();
	const { activateOnFocus, highlightedTabIndex, onTabActivation, registerTabResizeObserverElement, setHighlightedTabIndex, tabsListElement } = useTabsListContext();
	const id = useBaseUiId(idProp);
	const { compositeProps, compositeRef, index } = useCompositeItem({ metadata: import_react.useMemo(() => ({
		disabled,
		id,
		value
	}), [
		disabled,
		id,
		value
	]) });
	const active = value === activeTabValue;
	const isNavigatingRef = import_react.useRef(false);
	const tabElementRef = import_react.useRef(null);
	import_react.useEffect(() => {
		const tabElement = tabElementRef.current;
		if (!tabElement) return;
		return registerTabResizeObserverElement(tabElement);
	}, [registerTabResizeObserverElement]);
	useIsoLayoutEffect(() => {
		if (isNavigatingRef.current) {
			isNavigatingRef.current = false;
			return;
		}
		if (!(active && index > -1 && highlightedTabIndex !== index)) return;
		const listElement = tabsListElement;
		if (listElement != null) {
			const activeEl = activeElement(ownerDocument(listElement));
			if (activeEl && contains(listElement, activeEl)) return;
		}
		if (!disabled) setHighlightedTabIndex(index);
	}, [
		active,
		index,
		highlightedTabIndex,
		setHighlightedTabIndex,
		disabled,
		tabsListElement
	]);
	const { getButtonProps, buttonRef } = useButton({
		disabled,
		native: nativeButton,
		focusableWhenDisabled: true
	});
	const tabPanelId = getTabPanelIdByValue(value);
	const isPressingRef = import_react.useRef(false);
	const isMainButtonRef = import_react.useRef(false);
	function onClick(event) {
		if (active || disabled) return;
		onTabActivation(value, createChangeEventDetails(none, event.nativeEvent, void 0, { activationDirection: "none" }));
	}
	function onFocus(event) {
		if (active) return;
		if (index > -1 && !disabled) setHighlightedTabIndex(index);
		if (disabled) return;
		if (activateOnFocus && (!isPressingRef.current || isPressingRef.current && isMainButtonRef.current)) onTabActivation(value, createChangeEventDetails(none, event.nativeEvent, void 0, { activationDirection: "none" }));
	}
	function onPointerDown(event) {
		if (active || disabled) return;
		isPressingRef.current = true;
		function handlePointerUp() {
			isPressingRef.current = false;
			isMainButtonRef.current = false;
		}
		if (!event.button || event.button === 0) {
			isMainButtonRef.current = true;
			ownerDocument(event.currentTarget).addEventListener("pointerup", handlePointerUp, { once: true });
		}
	}
	return useRenderElement("button", componentProps, {
		state: {
			disabled,
			active,
			orientation
		},
		ref: [
			forwardedRef,
			buttonRef,
			compositeRef,
			tabElementRef
		],
		props: [
			compositeProps,
			{
				role: "tab",
				"aria-controls": tabPanelId,
				"aria-selected": active,
				id,
				onClick,
				onFocus,
				onPointerDown,
				[ACTIVE_COMPOSITE_ITEM]: active ? "" : void 0,
				onKeyDownCapture() {
					isNavigatingRef.current = true;
				}
			},
			elementProps,
			getButtonProps
		]
	});
});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/tabs/panel/TabsPanelDataAttributes.js
var TabsPanelDataAttributes = function(TabsPanelDataAttributes) {
	/**
	* Indicates the index of the tab panel.
	*/
	TabsPanelDataAttributes["index"] = "data-index";
	/**
	* Indicates the direction of the activation (based on the previous active tab).
	* @type {'left' | 'right' | 'up' | 'down' | 'none'}
	*/
	TabsPanelDataAttributes["activationDirection"] = "data-activation-direction";
	/**
	* Indicates the orientation of the tabs.
	* @type {'horizontal' | 'vertical'}
	*/
	TabsPanelDataAttributes["orientation"] = "data-orientation";
	/**
	* Present when the panel is hidden.
	*/
	TabsPanelDataAttributes["hidden"] = "data-hidden";
	/**
	* Present when the panel is animating in.
	*/
	TabsPanelDataAttributes[TabsPanelDataAttributes["startingStyle"] = TransitionStatusDataAttributes.startingStyle] = "startingStyle";
	/**
	* Present when the panel is animating out.
	*/
	TabsPanelDataAttributes[TabsPanelDataAttributes["endingStyle"] = TransitionStatusDataAttributes.endingStyle] = "endingStyle";
	return TabsPanelDataAttributes;
}({});
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/tabs/panel/TabsPanel.js
var stateAttributesMapping = {
	...tabsStateAttributesMapping,
	...transitionStatusMapping
};
/**
* A panel displayed when the corresponding tab is active.
* Renders a `<div>` element.
*
* Documentation: [Base UI Tabs](https://base-ui.com/react/components/tabs)
*/
var TabsPanel = /* @__PURE__ */ import_react.forwardRef(function TabPanel(componentProps, forwardedRef) {
	const { className, value, render, keepMounted = false, ...elementProps } = componentProps;
	const { value: selectedValue, getTabIdByPanelValue, orientation, tabActivationDirection, registerMountedTabPanel, unregisterMountedTabPanel } = useTabsRootContext();
	const id = useBaseUiId();
	const { ref: listItemRef, index } = useCompositeListItem({ metadata: import_react.useMemo(() => ({
		id,
		value
	}), [id, value]) });
	const open = value === selectedValue;
	const { mounted, transitionStatus, setMounted } = useTransitionStatus(open);
	const hidden = !mounted;
	const correspondingTabId = getTabIdByPanelValue(value);
	const state = {
		hidden,
		orientation,
		tabActivationDirection,
		transitionStatus
	};
	const panelRef = import_react.useRef(null);
	const element = useRenderElement("div", componentProps, {
		state,
		ref: [
			forwardedRef,
			listItemRef,
			panelRef
		],
		props: [{
			"aria-labelledby": correspondingTabId,
			hidden,
			id,
			role: "tabpanel",
			tabIndex: open ? 0 : -1,
			inert: inertValue(!open),
			[TabsPanelDataAttributes.index]: index
		}, elementProps],
		stateAttributesMapping
	});
	useOpenChangeComplete({
		open,
		ref: panelRef,
		onComplete() {
			if (!open) setMounted(false);
		}
	});
	useIsoLayoutEffect(() => {
		if (hidden && !keepMounted) return;
		if (id == null) return;
		registerMountedTabPanel(value, id);
		return () => {
			unregisterMountedTabPanel(value, id);
		};
	}, [
		hidden,
		keepMounted,
		value,
		id,
		registerMountedTabPanel,
		unregisterMountedTabPanel
	]);
	if (!(keepMounted || mounted)) return null;
	return element;
});
//#endregion
//#region node_modules/.pnpm/@base-ui+utils@0.2.6_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/utils/esm/isElementDisabled.js
function isElementDisabled(element) {
	return element == null || element.hasAttribute("disabled") || element.getAttribute("aria-disabled") === "true";
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/composite/root/useCompositeRoot.js
var EMPTY_ARRAY = [];
function useCompositeRoot(params) {
	const { itemSizes, cols = 1, loopFocus = true, dense = false, orientation = "both", direction, highlightedIndex: externalHighlightedIndex, onHighlightedIndexChange: externalSetHighlightedIndex, rootRef: externalRef, enableHomeAndEndKeys = false, stopEventPropagation = false, disabledIndices, modifierKeys = EMPTY_ARRAY } = params;
	const [internalHighlightedIndex, internalSetHighlightedIndex] = import_react.useState(0);
	const isGrid = cols > 1;
	const rootRef = import_react.useRef(null);
	const mergedRef = useMergedRefs(rootRef, externalRef);
	const elementsRef = import_react.useRef([]);
	const hasSetDefaultIndexRef = import_react.useRef(false);
	const highlightedIndex = externalHighlightedIndex ?? internalHighlightedIndex;
	const onHighlightedIndexChange = useStableCallback((index, shouldScrollIntoView = false) => {
		(externalSetHighlightedIndex ?? internalSetHighlightedIndex)(index);
		if (shouldScrollIntoView) {
			const newActiveItem = elementsRef.current[index];
			scrollIntoViewIfNeeded(rootRef.current, newActiveItem, direction, orientation);
		}
	});
	const onMapChange = useStableCallback((map) => {
		if (map.size === 0 || hasSetDefaultIndexRef.current) return;
		hasSetDefaultIndexRef.current = true;
		const sortedElements = Array.from(map.keys());
		const activeItem = sortedElements.find((compositeElement) => compositeElement?.hasAttribute("data-composite-item-active")) ?? null;
		const activeIndex = activeItem ? sortedElements.indexOf(activeItem) : -1;
		if (activeIndex !== -1) onHighlightedIndexChange(activeIndex);
		scrollIntoViewIfNeeded(rootRef.current, activeItem, direction, orientation);
	});
	const props = import_react.useMemo(() => ({
		"aria-orientation": orientation === "both" ? void 0 : orientation,
		ref: mergedRef,
		onFocus(event) {
			if (!rootRef.current || !isNativeInput(event.target)) return;
			event.target.setSelectionRange(0, event.target.value.length ?? 0);
		},
		onKeyDown(event) {
			const RELEVANT_KEYS = enableHomeAndEndKeys ? ALL_KEYS : ARROW_KEYS;
			if (!RELEVANT_KEYS.has(event.key)) return;
			if (isModifierKeySet(event, modifierKeys)) return;
			if (!rootRef.current) return;
			const isRtl = direction === "rtl";
			const horizontalForwardKey = isRtl ? ARROW_LEFT : ARROW_RIGHT;
			const forwardKey = {
				horizontal: horizontalForwardKey,
				vertical: ARROW_DOWN,
				both: horizontalForwardKey
			}[orientation];
			const horizontalBackwardKey = isRtl ? ARROW_RIGHT : ARROW_LEFT;
			const backwardKey = {
				horizontal: horizontalBackwardKey,
				vertical: ARROW_UP,
				both: horizontalBackwardKey
			}[orientation];
			if (isNativeInput(event.target) && !isElementDisabled(event.target)) {
				const selectionStart = event.target.selectionStart;
				const selectionEnd = event.target.selectionEnd;
				const textContent = event.target.value ?? "";
				if (selectionStart == null || event.shiftKey || selectionStart !== selectionEnd) return;
				if (event.key !== backwardKey && selectionStart < textContent.length) return;
				if (event.key !== forwardKey && selectionStart > 0) return;
			}
			let nextIndex = highlightedIndex;
			const minIndex = getMinListIndex(elementsRef, disabledIndices);
			const maxIndex = getMaxListIndex(elementsRef, disabledIndices);
			if (isGrid) {
				const sizes = itemSizes || Array.from({ length: elementsRef.current.length }, () => ({
					width: 1,
					height: 1
				}));
				const cellMap = createGridCellMap(sizes, cols, dense);
				const minGridIndex = cellMap.findIndex((index) => index != null && !isListIndexDisabled(elementsRef, index, disabledIndices));
				const maxGridIndex = cellMap.reduce((foundIndex, index, cellIndex) => index != null && !isListIndexDisabled(elementsRef, index, disabledIndices) ? cellIndex : foundIndex, -1);
				nextIndex = cellMap[getGridNavigatedIndex({ current: cellMap.map((itemIndex) => itemIndex ? elementsRef.current[itemIndex] : null) }, {
					event,
					orientation,
					loopFocus,
					cols,
					disabledIndices: getGridCellIndices([...disabledIndices || elementsRef.current.map((_, index) => isListIndexDisabled(elementsRef, index) ? index : void 0), void 0], cellMap),
					minIndex: minGridIndex,
					maxIndex: maxGridIndex,
					prevIndex: getGridCellIndexOfCorner(highlightedIndex > maxIndex ? minIndex : highlightedIndex, sizes, cellMap, cols, event.key === "ArrowDown" ? "bl" : event.key === "ArrowRight" ? "tr" : "tl"),
					rtl: isRtl
				})];
			}
			const forwardKeys = {
				horizontal: [horizontalForwardKey],
				vertical: [ARROW_DOWN],
				both: [horizontalForwardKey, ARROW_DOWN]
			}[orientation];
			const backwardKeys = {
				horizontal: [horizontalBackwardKey],
				vertical: [ARROW_UP],
				both: [horizontalBackwardKey, ARROW_UP]
			}[orientation];
			const preventedKeys = isGrid ? RELEVANT_KEYS : {
				horizontal: enableHomeAndEndKeys ? HORIZONTAL_KEYS_WITH_EXTRA_KEYS : HORIZONTAL_KEYS,
				vertical: enableHomeAndEndKeys ? VERTICAL_KEYS_WITH_EXTRA_KEYS : VERTICAL_KEYS,
				both: RELEVANT_KEYS
			}[orientation];
			if (enableHomeAndEndKeys) {
				if (event.key === "Home") nextIndex = minIndex;
				else if (event.key === "End") nextIndex = maxIndex;
			}
			if (nextIndex === highlightedIndex && (forwardKeys.includes(event.key) || backwardKeys.includes(event.key))) if (loopFocus && nextIndex === maxIndex && forwardKeys.includes(event.key)) nextIndex = minIndex;
			else if (loopFocus && nextIndex === minIndex && backwardKeys.includes(event.key)) nextIndex = maxIndex;
			else nextIndex = findNonDisabledListIndex(elementsRef, {
				startingIndex: nextIndex,
				decrement: backwardKeys.includes(event.key),
				disabledIndices
			});
			if (nextIndex !== highlightedIndex && !isIndexOutOfListBounds(elementsRef, nextIndex)) {
				if (stopEventPropagation) event.stopPropagation();
				if (preventedKeys.has(event.key)) event.preventDefault();
				onHighlightedIndexChange(nextIndex, true);
				queueMicrotask(() => {
					elementsRef.current[nextIndex]?.focus();
				});
			}
		}
	}), [
		cols,
		dense,
		direction,
		disabledIndices,
		elementsRef,
		enableHomeAndEndKeys,
		highlightedIndex,
		isGrid,
		itemSizes,
		loopFocus,
		mergedRef,
		modifierKeys,
		onHighlightedIndexChange,
		orientation,
		stopEventPropagation
	]);
	return import_react.useMemo(() => ({
		props,
		highlightedIndex,
		onHighlightedIndexChange,
		elementsRef,
		disabledIndices,
		onMapChange,
		relayKeyboardEvent: props.onKeyDown
	}), [
		props,
		highlightedIndex,
		onHighlightedIndexChange,
		elementsRef,
		disabledIndices,
		onMapChange
	]);
}
function isModifierKeySet(event, ignoredModifierKeys) {
	for (const key of MODIFIER_KEYS.values()) {
		if (ignoredModifierKeys.includes(key)) continue;
		if (event.getModifierState(key)) return true;
	}
	return false;
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/composite/root/CompositeRoot.js
/**
* @internal
*/
function CompositeRoot(componentProps) {
	const { render, className, refs = EMPTY_ARRAY$1, props = EMPTY_ARRAY$1, state = EMPTY_OBJECT, stateAttributesMapping, highlightedIndex: highlightedIndexProp, onHighlightedIndexChange: onHighlightedIndexChangeProp, orientation, dense, itemSizes, loopFocus, cols, enableHomeAndEndKeys, onMapChange: onMapChangeProp, stopEventPropagation = true, rootRef, disabledIndices, modifierKeys, highlightItemOnHover = false, tag = "div", ...elementProps } = componentProps;
	const { props: defaultProps, highlightedIndex, onHighlightedIndexChange, elementsRef, onMapChange: onMapChangeUnwrapped, relayKeyboardEvent } = useCompositeRoot({
		itemSizes,
		cols,
		loopFocus,
		dense,
		orientation,
		highlightedIndex: highlightedIndexProp,
		onHighlightedIndexChange: onHighlightedIndexChangeProp,
		rootRef,
		stopEventPropagation,
		enableHomeAndEndKeys,
		direction: useDirection(),
		disabledIndices,
		modifierKeys
	});
	const element = useRenderElement(tag, componentProps, {
		state,
		ref: refs,
		props: [
			defaultProps,
			...props,
			elementProps
		],
		stateAttributesMapping
	});
	const contextValue = import_react.useMemo(() => ({
		highlightedIndex,
		onHighlightedIndexChange,
		highlightItemOnHover,
		relayKeyboardEvent
	}), [
		highlightedIndex,
		onHighlightedIndexChange,
		highlightItemOnHover,
		relayKeyboardEvent
	]);
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(CompositeRootContext.Provider, {
		value: contextValue,
		children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(CompositeList, {
			elementsRef,
			onMapChange: (newMap) => {
				onMapChangeProp?.(newMap);
				onMapChangeUnwrapped(newMap);
			},
			children: element
		})
	});
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/tabs/list/TabsList.js
/**
* Groups the individual tab buttons.
* Renders a `<div>` element.
*
* Documentation: [Base UI Tabs](https://base-ui.com/react/components/tabs)
*/
var TabsList$1 = /* @__PURE__ */ import_react.forwardRef(function TabsList(componentProps, forwardedRef) {
	const { activateOnFocus = false, className, loopFocus = true, render, ...elementProps } = componentProps;
	const { getTabElementBySelectedValue, onValueChange, orientation, value, setTabMap, tabActivationDirection } = useTabsRootContext();
	const [highlightedTabIndex, setHighlightedTabIndex] = import_react.useState(0);
	const [tabsListElement, setTabsListElement] = import_react.useState(null);
	const indicatorUpdateListenersRef = import_react.useRef(/* @__PURE__ */ new Set());
	const tabResizeObserverElementsRef = import_react.useRef(/* @__PURE__ */ new Set());
	const resizeObserverRef = import_react.useRef(null);
	const notifyIndicatorUpdateListeners = useStableCallback(() => {
		indicatorUpdateListenersRef.current.forEach((listener) => {
			listener();
		});
	});
	import_react.useEffect(() => {
		if (typeof ResizeObserver === "undefined") return;
		const resizeObserver = new ResizeObserver(() => {
			if (!indicatorUpdateListenersRef.current.size) return;
			notifyIndicatorUpdateListeners();
		});
		resizeObserverRef.current = resizeObserver;
		if (tabsListElement) resizeObserver.observe(tabsListElement);
		tabResizeObserverElementsRef.current.forEach((element) => {
			resizeObserver.observe(element);
		});
		return () => {
			resizeObserver.disconnect();
			resizeObserverRef.current = null;
		};
	}, [tabsListElement, notifyIndicatorUpdateListeners]);
	const registerIndicatorUpdateListener = useStableCallback((listener) => {
		indicatorUpdateListenersRef.current.add(listener);
		return () => {
			indicatorUpdateListenersRef.current.delete(listener);
		};
	});
	const registerTabResizeObserverElement = useStableCallback((element) => {
		tabResizeObserverElementsRef.current.add(element);
		resizeObserverRef.current?.observe(element);
		return () => {
			tabResizeObserverElementsRef.current.delete(element);
			resizeObserverRef.current?.unobserve(element);
		};
	});
	const detectActivationDirection = useActivationDirectionDetector(value, orientation, tabsListElement, getTabElementBySelectedValue);
	const onTabActivation = useStableCallback((newValue, eventDetails) => {
		if (newValue !== value) {
			eventDetails.activationDirection = detectActivationDirection(newValue);
			onValueChange(newValue, eventDetails);
		}
	});
	const state = {
		orientation,
		tabActivationDirection
	};
	const defaultProps = {
		"aria-orientation": orientation === "vertical" ? "vertical" : void 0,
		role: "tablist"
	};
	const tabsListContextValue = import_react.useMemo(() => ({
		activateOnFocus,
		highlightedTabIndex,
		registerIndicatorUpdateListener,
		registerTabResizeObserverElement,
		onTabActivation,
		setHighlightedTabIndex,
		tabsListElement
	}), [
		activateOnFocus,
		highlightedTabIndex,
		registerIndicatorUpdateListener,
		registerTabResizeObserverElement,
		onTabActivation,
		setHighlightedTabIndex,
		tabsListElement
	]);
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(TabsListContext.Provider, {
		value: tabsListContextValue,
		children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(CompositeRoot, {
			render,
			className,
			state,
			refs: [forwardedRef, setTabsListElement],
			props: [defaultProps, elementProps],
			stateAttributesMapping: tabsStateAttributesMapping,
			highlightedIndex: highlightedTabIndex,
			enableHomeAndEndKeys: true,
			loopFocus,
			orientation,
			onHighlightedIndexChange: setHighlightedTabIndex,
			onMapChange: setTabMap,
			disabledIndices: EMPTY_ARRAY$1
		})
	});
});
function getInset(tab, tabsList) {
	const { left: tabLeft, top: tabTop } = tab.getBoundingClientRect();
	const { left: listLeft, top: listTop } = tabsList.getBoundingClientRect();
	return {
		left: tabLeft - listLeft,
		top: tabTop - listTop
	};
}
function useActivationDirectionDetector(activeTabValue, orientation, tabsListElement, getTabElement) {
	const [previousTabEdge, setPreviousTabEdge] = import_react.useState(null);
	useIsoLayoutEffect(() => {
		if (activeTabValue == null || tabsListElement == null) {
			setPreviousTabEdge(null);
			return;
		}
		const activeTab = getTabElement(activeTabValue);
		if (activeTab == null) {
			setPreviousTabEdge(null);
			return;
		}
		const { left, top } = getInset(activeTab, tabsListElement);
		setPreviousTabEdge(orientation === "horizontal" ? left : top);
	}, [
		orientation,
		getTabElement,
		tabsListElement,
		activeTabValue
	]);
	return import_react.useCallback((newValue) => {
		if (newValue === activeTabValue) return "none";
		if (newValue == null) {
			setPreviousTabEdge(null);
			return "none";
		}
		if (newValue != null && tabsListElement != null) {
			const activeTabElement = getTabElement(newValue);
			if (activeTabElement != null) {
				const { left, top } = getInset(activeTabElement, tabsListElement);
				if (previousTabEdge == null) {
					setPreviousTabEdge(orientation === "horizontal" ? left : top);
					return "none";
				}
				if (orientation === "horizontal") {
					if (left < previousTabEdge) {
						setPreviousTabEdge(left);
						return "left";
					}
					if (left > previousTabEdge) {
						setPreviousTabEdge(left);
						return "right";
					}
				} else if (top < previousTabEdge) {
					setPreviousTabEdge(top);
					return "up";
				} else if (top > previousTabEdge) {
					setPreviousTabEdge(top);
					return "down";
				}
			}
		}
		return "none";
	}, [
		getTabElement,
		orientation,
		previousTabEdge,
		tabsListElement,
		activeTabValue
	]);
}
//#endregion
//#region src/components/ui/tabs.tsx
function Tabs({ className, orientation = "horizontal", ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(TabsRoot, {
		"data-slot": "tabs",
		"data-orientation": orientation,
		className: cn("group/tabs flex gap-2 data-horizontal:flex-col", className),
		...props
	});
}
var tabsListVariants = cva("group/tabs-list inline-flex w-fit items-center justify-center rounded-lg p-[3px] text-muted-foreground group-data-horizontal/tabs:h-8 group-data-vertical/tabs:h-fit group-data-vertical/tabs:flex-col data-[variant=line]:rounded-none", {
	variants: { variant: {
		default: "bg-muted",
		line: "gap-1 bg-transparent"
	} },
	defaultVariants: { variant: "default" }
});
function TabsList({ className, variant = "default", ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(TabsList$1, {
		"data-slot": "tabs-list",
		"data-variant": variant,
		className: cn(tabsListVariants({ variant }), className),
		...props
	});
}
function TabsTrigger({ className, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(TabsTab, {
		"data-slot": "tabs-trigger",
		className: cn("relative inline-flex h-[calc(100%-1px)] flex-1 items-center justify-center gap-1.5 rounded-md border border-transparent px-1.5 py-0.5 text-sm font-medium whitespace-nowrap text-foreground/60 transition-all group-data-vertical/tabs:w-full group-data-vertical/tabs:justify-start hover:text-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 focus-visible:outline-ring disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 dark:text-muted-foreground dark:hover:text-foreground group-data-[variant=default]/tabs-list:data-active:shadow-sm group-data-[variant=line]/tabs-list:data-active:shadow-none [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4", "group-data-[variant=line]/tabs-list:bg-transparent group-data-[variant=line]/tabs-list:data-active:bg-transparent dark:group-data-[variant=line]/tabs-list:data-active:border-transparent dark:group-data-[variant=line]/tabs-list:data-active:bg-transparent", "data-active:bg-background data-active:text-foreground dark:data-active:border-input dark:data-active:bg-input/30 dark:data-active:text-foreground", "after:absolute after:bg-foreground after:opacity-0 after:transition-opacity group-data-horizontal/tabs:after:inset-x-0 group-data-horizontal/tabs:after:bottom-[-5px] group-data-horizontal/tabs:after:h-0.5 group-data-vertical/tabs:after:inset-y-0 group-data-vertical/tabs:after:-right-1 group-data-vertical/tabs:after:w-0.5 group-data-[variant=line]/tabs-list:data-active:after:opacity-100", className),
		...props
	});
}
function TabsContent({ className, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(TabsPanel, {
		"data-slot": "tabs-content",
		className: cn("flex-1 text-sm outline-none", className),
		...props
	});
}
//#endregion
//#region src/components/ui/table.tsx
function Table({ className, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("div", {
		"data-slot": "table-container",
		className: "relative w-full overflow-x-auto",
		children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)("table", {
			"data-slot": "table",
			className: cn("w-full caption-bottom text-sm", className),
			...props
		})
	});
}
function TableHeader({ className, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("thead", {
		"data-slot": "table-header",
		className: cn("[&_tr]:border-b", className),
		...props
	});
}
function TableBody({ className, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("tbody", {
		"data-slot": "table-body",
		className: cn("[&_tr:last-child]:border-0", className),
		...props
	});
}
function TableRow({ className, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("tr", {
		"data-slot": "table-row",
		className: cn("border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted", className),
		...props
	});
}
function TableHead({ className, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("th", {
		"data-slot": "table-head",
		className: cn("h-10 px-2 text-left align-middle font-medium whitespace-nowrap text-foreground [&:has([role=checkbox])]:pr-0", className),
		...props
	});
}
function TableCell({ className, ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)("td", {
		"data-slot": "table-cell",
		className: cn("p-2 align-middle whitespace-nowrap [&:has([role=checkbox])]:pr-0", className),
		...props
	});
}
//#endregion
//#region node_modules/.pnpm/@base-ui+react@1.3.0_@types+react@19.2.14_react-dom@19.2.4_react@19.2.4__react@19.2.4/node_modules/@base-ui/react/esm/separator/Separator.js
/**
* A separator element accessible to screen readers.
* Renders a `<div>` element.
*
* Documentation: [Base UI Separator](https://base-ui.com/react/components/separator)
*/
var Separator$1 = /* @__PURE__ */ import_react.forwardRef(function SeparatorComponent(componentProps, forwardedRef) {
	const { className, render, orientation = "horizontal", ...elementProps } = componentProps;
	return useRenderElement("div", componentProps, {
		state: { orientation },
		ref: forwardedRef,
		props: [{
			role: "separator",
			"aria-orientation": orientation
		}, elementProps]
	});
});
//#endregion
//#region src/components/ui/separator.tsx
function Separator({ className, orientation = "horizontal", ...props }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Separator$1, {
		"data-slot": "separator",
		orientation,
		className: cn("shrink-0 bg-border data-horizontal:h-px data-horizontal:w-full data-vertical:w-px data-vertical:self-stretch", className),
		...props
	});
}
//#endregion
//#region src/routes/scripts/$site.$command.tsx?tsr-split=component
function DetailPage() {
	const script = Route.useLoaderData();
	if (!script) return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("main", {
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
	const args = Object.entries(script.args);
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("main", {
		className: "mx-auto max-w-4xl px-4 py-8",
		children: [
			/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Button, {
				variant: "ghost",
				size: "sm",
				className: "mb-4",
				render: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Link, { to: "/scripts" }),
				children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(ArrowLeft, { className: "size-3.5" }), "Back"]
			}),
			/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "mb-6",
				children: [
					/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
						className: "flex flex-wrap items-center gap-2 mb-2",
						children: [
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)("h1", {
								className: "text-2xl font-bold font-mono",
								children: script.name
							}),
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Badge, {
								variant: "secondary",
								children: script.site
							}),
							script.readOnly ? /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Badge, {
								variant: "outline",
								children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Lock, { className: "size-3" }), "Read-only"]
							}) : /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Badge, {
								variant: "outline",
								children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(LockOpen, { className: "size-3" }), "Read-write"]
							})
						]
					}),
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)("p", {
						className: "text-muted-foreground",
						children: script.description || "No description"
					}),
					script.domain && /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
						className: "mt-1 flex items-center gap-1 text-sm text-muted-foreground/70",
						children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Globe, { className: "size-3.5" }), script.domain]
					})
				]
			}),
			/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "mb-6 grid grid-cols-3 gap-3",
				children: [
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(UsageCard, {
						label: "Last 7 days",
						count: script.usage.last7d
					}),
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(UsageCard, {
						label: "Last 30 days",
						count: script.usage.last30d
					}),
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(UsageCard, {
						label: "All time",
						count: script.usage.total
					})
				]
			}),
			/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Separator, { className: "mb-6" }),
			/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Tabs, {
				defaultValue: "code",
				children: [
					/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(TabsList, { children: [
						/* @__PURE__ */ (0, import_jsx_runtime.jsx)(TabsTrigger, {
							value: "code",
							children: "Code"
						}),
						args.length > 0 && /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(TabsTrigger, {
							value: "args",
							children: [
								"Arguments (",
								args.length,
								")"
							]
						}),
						script.example && /* @__PURE__ */ (0, import_jsx_runtime.jsx)(TabsTrigger, {
							value: "example",
							children: "Example"
						})
					] }),
					/* @__PURE__ */ (0, import_jsx_runtime.jsx)(TabsContent, {
						value: "code",
						className: "mt-4",
						children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(CodeBlock, {
							content: script.content,
							name: script.name
						})
					}),
					args.length > 0 && /* @__PURE__ */ (0, import_jsx_runtime.jsx)(TabsContent, {
						value: "args",
						className: "mt-4",
						children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Card, { children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Table, { children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(TableHeader, { children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(TableRow, { children: [
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(TableHead, { children: "Name" }),
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(TableHead, { children: "Required" }),
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(TableHead, { children: "Description" })
						] }) }), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(TableBody, { children: args.map(([name, arg]) => /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(TableRow, { children: [
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(TableCell, {
								className: "font-mono text-sm",
								children: name
							}),
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(TableCell, { children: arg.required ? /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Badge, {
								variant: "default",
								className: "text-xs",
								children: "Required"
							}) : /* @__PURE__ */ (0, import_jsx_runtime.jsx)("span", {
								className: "text-muted-foreground text-xs",
								children: "Optional"
							}) }),
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(TableCell, {
								className: "text-muted-foreground whitespace-normal",
								children: arg.description
							})
						] }, name)) })] }) })
					}),
					script.example && /* @__PURE__ */ (0, import_jsx_runtime.jsx)(TabsContent, {
						value: "example",
						className: "mt-4",
						children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Card, { children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(CardContent, { children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)("pre", {
							className: "overflow-x-auto text-sm font-mono whitespace-pre-wrap",
							children: script.example
						}) }) })
					})
				]
			}),
			/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Separator, { className: "my-6" }),
			/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "flex flex-wrap items-center gap-4 text-xs text-muted-foreground",
				children: [
					/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("span", {
						className: "flex items-center gap-1",
						children: [
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Calendar, { className: "size-3" }),
							"Updated ",
							formatDate(script.updatedAt)
						]
					}),
					/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("span", {
						className: "flex items-center gap-1",
						children: [
							/* @__PURE__ */ (0, import_jsx_runtime.jsx)(Calendar, { className: "size-3" }),
							"Created ",
							formatDate(script.createdAt)
						]
					}),
					/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("span", {
						className: "font-mono",
						children: [
							"SHA-256: ",
							script.hash.slice(0, 12),
							"…"
						]
					})
				]
			})
		]
	});
}
function UsageCard({ label, count }) {
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Card, {
		size: "sm",
		children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(CardHeader, {
			className: "pb-0",
			children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)(CardTitle, {
				className: "text-xs font-normal text-muted-foreground",
				children: label
			})
		}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(CardContent, {
			className: "pt-0",
			children: /* @__PURE__ */ (0, import_jsx_runtime.jsxs)("div", {
				className: "flex items-center gap-1",
				children: [/* @__PURE__ */ (0, import_jsx_runtime.jsx)(TrendingUp, { className: "size-3.5 text-muted-foreground" }), /* @__PURE__ */ (0, import_jsx_runtime.jsx)("span", {
					className: "text-2xl font-bold",
					children: count
				})]
			})
		})]
	});
}
function CodeBlock({ content, name }) {
	const [copied, setCopied] = (0, import_react.useState)(false);
	function handleCopy() {
		navigator.clipboard.writeText(content);
		setCopied(true);
		setTimeout(() => setCopied(false), 2e3);
	}
	return /* @__PURE__ */ (0, import_jsx_runtime.jsxs)(Card, { children: [/* @__PURE__ */ (0, import_jsx_runtime.jsxs)(CardHeader, {
		className: "flex-row items-center justify-between border-b py-2",
		children: [/* @__PURE__ */ (0, import_jsx_runtime.jsxs)("span", {
			className: "font-mono text-xs text-muted-foreground",
			children: [name, ".js"]
		}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Button, {
			variant: "ghost",
			size: "icon-xs",
			onClick: handleCopy,
			children: copied ? /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Check, { className: "size-3" }) : /* @__PURE__ */ (0, import_jsx_runtime.jsx)(Copy, { className: "size-3" })
		})]
	}), /* @__PURE__ */ (0, import_jsx_runtime.jsx)(CardContent, {
		className: "p-0",
		children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)("pre", {
			className: "overflow-x-auto p-4 text-sm font-mono leading-relaxed",
			children: /* @__PURE__ */ (0, import_jsx_runtime.jsx)("code", { children: content })
		})
	})] });
}
function formatDate(dateStr) {
	return (/* @__PURE__ */ new Date(dateStr.replace(" ", "T") + "Z")).toLocaleDateString("en-US", {
		year: "numeric",
		month: "short",
		day: "numeric"
	});
}
//#endregion
export { DetailPage as component };
