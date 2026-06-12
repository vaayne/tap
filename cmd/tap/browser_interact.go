package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/chromedp/cdproto/input"
	"github.com/urfave/cli/v3"
)

func browserDblclickCmd() *cli.Command {
	return &cli.Command{
		Name:      "dblclick",
		Usage:     "Double-click an element by CSS selector or snapshot ref",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Description: `Dispatch a real double-click (clickCount=2) on the first visible element.

Examples:
  tap browser dblclick "td.editable"
  tap browser dblclick @e3`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("CSS selector or @eN ref required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.DblClickElement(ctx, cmd.String("session"), cmd.String("tab"), sel); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Double-clicked %s\n", sel)
			return nil
		},
	}
}

func browserFocusCmd() *cli.Command {
	return &cli.Command{
		Name:      "focus",
		Usage:     "Focus an element by CSS selector or snapshot ref",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Description: `Move browser focus to the element matching the selector.

Examples:
  tap browser focus "input[name=email]"
  tap browser focus @e1`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("CSS selector or @eN ref required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.FocusElement(ctx, cmd.String("session"), cmd.String("tab"), sel); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Focused %s\n", sel)
			return nil
		},
	}
}

func browserCheckCmd() *cli.Command {
	return &cli.Command{
		Name:      "check",
		Usage:     "Ensure a checkbox is checked",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Description: `Read the current checkbox state and click only if unchecked.
Dispatches React-compatible events on state change.

Examples:
  tap browser check "input[name=agree]"
  tap browser check @e4`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("CSS selector or @eN ref required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.CheckElement(ctx, cmd.String("session"), cmd.String("tab"), sel); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Checked %s\n", sel)
			return nil
		},
	}
}

func browserUncheckCmd() *cli.Command {
	return &cli.Command{
		Name:      "uncheck",
		Usage:     "Ensure a checkbox is unchecked",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Description: `Read the current checkbox state and click only if checked.
Dispatches React-compatible events on state change.

Examples:
  tap browser uncheck "input[name=newsletter]"
  tap browser uncheck @e5`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("CSS selector or @eN ref required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.UncheckElement(ctx, cmd.String("session"), cmd.String("tab"), sel); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Unchecked %s\n", sel)
			return nil
		},
	}
}

func browserScrollIntoViewCmd() *cli.Command {
	return &cli.Command{
		Name:      "scrollintoview",
		Usage:     "Scroll an element into the viewport",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Description: `Scroll the element matching the selector into view using DOM.scrollIntoViewIfNeeded.

Examples:
  tap browser scrollintoview "#footer"
  tap browser scrollintoview @e7`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("CSS selector or @eN ref required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.ScrollIntoViewElement(ctx, cmd.String("session"), cmd.String("tab"), sel); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Scrolled %s into view\n", sel)
			return nil
		},
	}
}

func browserUploadCmd() *cli.Command {
	return &cli.Command{
		Name:      "upload",
		Usage:     "Set files on a file input element",
		ArgsUsage: "<selector> <file> [file...]",
		Flags:     browserActionFlags(false),
		Description: `Set one or more local file paths on a <input type=file> element via
DOM.setFileInputFiles. The files must exist on the local filesystem.

Examples:
  tap browser upload "input[type=file]" /tmp/report.pdf
  tap browser upload "#avatar" photo.jpg`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser upload <selector> <file> [file...]")
			}
			sel := args[0]
			files := args[1:]
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.UploadFiles(ctx, cmd.String("session"), cmd.String("tab"), sel, files); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Uploaded %d file(s) to %s\n", len(files), sel)
			return nil
		},
	}
}

func browserDragCmd() *cli.Command {
	return &cli.Command{
		Name:      "drag",
		Usage:     "Mouse-based drag and drop from source to destination",
		ArgsUsage: "<src-selector> <dst-selector>",
		Flags:     browserActionFlags(false),
		Description: `Perform a real mouse drag: move→press→interpolate→release.

Both arguments are CSS selectors (or snapshot refs whose selector hints are used).

Examples:
  tap browser drag ".card" ".dropzone"
  tap browser drag @e2 @e8`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser drag <src-selector> <dst-selector>")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.Drag(ctx, cmd.String("session"), cmd.String("tab"), args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Dragged %s → %s\n", args[0], args[1])
			return nil
		},
	}
}

// browserMouseCmd returns the "mouse" parent command with move/down/up/wheel sub-commands.
func browserMouseCmd() *cli.Command {
	return &cli.Command{
		Name:  "mouse",
		Usage: "Low-level mouse event dispatch",
		Description: `Dispatch individual mouse events.

Sub-commands:
  move <x> <y>              Move cursor to absolute position
  down [button]             Press button (left|right|middle, default left)
  up   [button]             Release button
  wheel <dy> [dx]           Scroll vertically (and optionally horizontally)`,
		Commands: []*cli.Command{
			browserMouseMoveCmd(),
			browserMouseDownCmd(),
			browserMouseUpCmd(),
			browserMouseWheelCmd(),
		},
	}
}

func browserMouseMoveCmd() *cli.Command {
	return &cli.Command{
		Name:      "move",
		Usage:     "Move the mouse cursor to absolute coordinates",
		ArgsUsage: "<x> <y>",
		Flags:     browserActionFlags(false),
		Description: `Dispatch a mouseMoved event to the given absolute pixel position.
Coordinates are relative to the viewport top-left corner.

Examples:
  tap browser mouse move 640 480
  tap browser mouse move 0 0`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser mouse move <x> <y>")
			}
			x, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return fmt.Errorf("invalid x: %w", err)
			}
			y, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return fmt.Errorf("invalid y: %w", err)
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.MouseMove(ctx, cmd.String("session"), cmd.String("tab"), x, y); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Mouse moved to (%s, %s)\n", args[0], args[1])
			return nil
		},
	}
}

func resolveMouseButton(name string) (input.MouseButton, error) {
	switch name {
	case "", "left":
		return input.Left, nil
	case "right":
		return input.Right, nil
	case "middle":
		return input.Middle, nil
	default:
		return input.None, fmt.Errorf("unknown button %q: use left, right, or middle", name)
	}
}

func browserMouseDownCmd() *cli.Command {
	return &cli.Command{
		Name:      "down",
		Usage:     "Press a mouse button",
		ArgsUsage: "[button]",
		Flags:     browserActionFlags(false),
		Description: `Dispatch a mousePressed event. button defaults to left.

Examples:
  tap browser mouse down
  tap browser mouse down right`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			btn, err := resolveMouseButton(cmd.Args().First())
			if err != nil {
				return err
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.MouseDown(ctx, cmd.String("session"), cmd.String("tab"), btn); err != nil {
				return err
			}
			btnName := string(btn)
			if btnName == "" {
				btnName = "left"
			}
			fmt.Fprintf(os.Stderr, "Mouse %s button down\n", btnName)
			return nil
		},
	}
}

func browserMouseUpCmd() *cli.Command {
	return &cli.Command{
		Name:      "up",
		Usage:     "Release a mouse button",
		ArgsUsage: "[button]",
		Flags:     browserActionFlags(false),
		Description: `Dispatch a mouseReleased event. button defaults to left.

Examples:
  tap browser mouse up
  tap browser mouse up right`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			btn, err := resolveMouseButton(cmd.Args().First())
			if err != nil {
				return err
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.MouseUp(ctx, cmd.String("session"), cmd.String("tab"), btn); err != nil {
				return err
			}
			btnName := string(btn)
			if btnName == "" {
				btnName = "left"
			}
			fmt.Fprintf(os.Stderr, "Mouse %s button up\n", btnName)
			return nil
		},
	}
}

func browserMouseWheelCmd() *cli.Command {
	return &cli.Command{
		Name:      "wheel",
		Usage:     "Dispatch a mouse wheel scroll event",
		ArgsUsage: "<dy> [dx]",
		Flags:     browserActionFlags(false),
		Description: `Dispatch a mouseWheel event. dy scrolls vertically (positive = down),
dx scrolls horizontally (default 0).

Examples:
  tap browser mouse wheel 300
  tap browser mouse wheel -200 50`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 1 {
				return fmt.Errorf("usage: tap browser mouse wheel <dy> [dx]")
			}
			dy, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return fmt.Errorf("invalid dy: %w", err)
			}
			var dx float64
			if len(args) >= 2 {
				dx, err = strconv.ParseFloat(args[1], 64)
				if err != nil {
					return fmt.Errorf("invalid dx: %w", err)
				}
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.MouseWheel(ctx, cmd.String("session"), cmd.String("tab"), dy, dx); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Mouse wheel dy=%s dx=%s\n", args[0], func() string {
				if len(args) >= 2 {
					return args[1]
				}
				return "0"
			}())
			return nil
		},
	}
}

// browserKeyboardCmd returns the "keyboard" parent command with type/insert sub-commands.
func browserKeyboardCmd() *cli.Command {
	return &cli.Command{
		Name:  "keyboard",
		Usage: "Low-level keyboard event dispatch",
		Description: `Dispatch keyboard events at the current focus point.

Sub-commands:
  type <text>     Send per-character key events (real typing)
  insert <text>   Paste text instantly via Input.insertText (no key events)`,
		Commands: []*cli.Command{
			browserKeyboardTypeCmd(),
			browserKeyboardInsertCmd(),
		},
	}
}

func browserKeyboardTypeCmd() *cli.Command {
	return &cli.Command{
		Name:      "type",
		Usage:     "Type text at current focus with real key events",
		ArgsUsage: "<text>",
		Flags:     browserActionFlags(false),
		Description: `Send per-character keyDown/char/keyUp events for the given text.
Use this when the site validates per-keystroke input.

Examples:
  tap browser keyboard type "hello world"`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			text := cmd.Args().First()
			if text == "" {
				return fmt.Errorf("text required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.KeyboardType(ctx, cmd.String("session"), cmd.String("tab"), text); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Typed %q\n", text)
			return nil
		},
	}
}

func browserKeyboardInsertCmd() *cli.Command {
	return &cli.Command{
		Name:      "insert",
		Usage:     "Insert text instantly without key events",
		ArgsUsage: "<text>",
		Flags:     browserActionFlags(false),
		Description: `Insert text via Input.insertText — no keyDown/keyUp events are dispatched.
Faster than 'type' but bypasses keystroke listeners.

Examples:
  tap browser keyboard insert "paste this text"`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			text := cmd.Args().First()
			if text == "" {
				return fmt.Errorf("text required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.KeyboardInsert(ctx, cmd.String("session"), cmd.String("tab"), text); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Inserted %q\n", text)
			return nil
		},
	}
}

func browserKeydownCmd() *cli.Command {
	return &cli.Command{
		Name:      "keydown",
		Usage:     "Hold a key down",
		ArgsUsage: "<key>",
		Flags:     browserActionFlags(false),
		Description: `Dispatch a rawKeyDown event for the given key name.
Use 'keyup' to release. Standard key names: Enter, Tab, Escape, Shift, Control, etc.

Examples:
  tap browser keydown Shift
  tap browser keydown Control`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			key := cmd.Args().First()
			if key == "" {
				return fmt.Errorf("key required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.Keydown(ctx, cmd.String("session"), cmd.String("tab"), key); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Key down: %s\n", key)
			return nil
		},
	}
}

func browserKeyupCmd() *cli.Command {
	return &cli.Command{
		Name:      "keyup",
		Usage:     "Release a held key",
		ArgsUsage: "<key>",
		Flags:     browserActionFlags(false),
		Description: `Dispatch a keyUp event for the given key name.

Examples:
  tap browser keyup Shift
  tap browser keyup Control`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			key := cmd.Args().First()
			if key == "" {
				return fmt.Errorf("key required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.Keyup(ctx, cmd.String("session"), cmd.String("tab"), key); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Key up: %s\n", key)
			return nil
		},
	}
}
