package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dop251/goja"
	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/agentbrowser"
)

type workflowBrowser interface {
	Run(context.Context, ...string) (json.RawMessage, error)
	Eval(context.Context, string) (any, error)
}

func runCmd() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "Run a host-side JavaScript browser workflow",
		ArgsUsage: "[script.js]",
		Description: `Execute JavaScript that drives the active agent-browser session.
Read the workflow from a file, or from stdin when no file (or -) is given.

The runtime exposes browser.cmd(...args), browser.eval(script), console.log(),
and shortcuts for open and snapshot.

Examples:
  tap run workflow.js
  tap run <<'JS'
  await browser.open("https://example.com")
  const page = await browser.snapshot("-i")
  console.log(page.snapshot)
  JS`,
		Action: runWorkflow,
	}
}

func runWorkflow(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() > 1 {
		return fmt.Errorf("tap run accepts at most one script file")
	}

	reader := cmd.Root().Reader
	if reader == nil {
		reader = os.Stdin
	}
	source, filename, err := readWorkflow(cmd.Args().First(), reader)
	if err != nil {
		return err
	}

	if timeout := cmd.Duration("timeout"); timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	stdout := cmd.Root().Writer
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := cmd.Root().ErrWriter
	if stderr == nil {
		stderr = os.Stderr
	}
	browser := agentbrowser.New(cmd.String("agent-browser"))
	return executeWorkflow(ctx, browser, source, filename, stdout, stderr)
}

func readWorkflow(path string, stdin io.Reader) (source, filename string, err error) {
	if path == "" || path == "-" {
		data, err := io.ReadAll(stdin)
		if err != nil {
			return "", "", fmt.Errorf("read workflow from stdin: %w", err)
		}
		return string(data), "<stdin>", nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("read workflow %s: %w", path, err)
	}
	return string(data), path, nil
}

func executeWorkflow(
	ctx context.Context,
	browser workflowBrowser,
	source string,
	filename string,
	stdout io.Writer,
	stderr io.Writer,
) error {
	vm := goja.New()
	stopInterrupt := context.AfterFunc(ctx, func() {
		vm.Interrupt(ctx.Err())
	})
	defer stopInterrupt()

	if err := bindConsole(vm, stdout, stderr); err != nil {
		return err
	}
	if err := bindBrowser(vm, ctx, browser); err != nil {
		return err
	}
	if _, err := vm.RunString(`
for (const command of ["open", "snapshot"]) {
  Object.defineProperty(browser, command, {
    value: (...args) => browser.cmd(command, ...args)
  });
}
`); err != nil {
		return fmt.Errorf("initialize workflow runtime: %w", err)
	}

	program := "(async () => {\n" + source + "\n})()"
	value, err := vm.RunScript(filename, program)
	if err != nil {
		return fmt.Errorf("run workflow: %w", err)
	}
	promise, ok := value.Export().(*goja.Promise)
	if !ok {
		return fmt.Errorf("run workflow: async program did not return a promise")
	}
	switch promise.State() {
	case goja.PromiseStateFulfilled:
		return nil
	case goja.PromiseStateRejected:
		return fmt.Errorf("run workflow: %w", promiseError(promise.Result()))
	default:
		return fmt.Errorf("run workflow: pending promises require unsupported asynchronous callbacks")
	}
}

func bindConsole(vm *goja.Runtime, stdout, stderr io.Writer) error {
	console := vm.NewObject()
	for name, writer := range map[string]io.Writer{
		"log":   stdout,
		"info":  stdout,
		"warn":  stderr,
		"error": stderr,
	} {
		if err := console.Set(name, consoleWriter(vm, writer)); err != nil {
			return fmt.Errorf("initialize console.%s: %w", name, err)
		}
	}
	if err := vm.Set("console", console); err != nil {
		return fmt.Errorf("initialize console: %w", err)
	}
	return nil
}

func consoleWriter(vm *goja.Runtime, writer io.Writer) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		parts := make([]string, len(call.Arguments))
		for index, value := range call.Arguments {
			parts[index] = value.String()
		}
		if _, err := fmt.Fprintln(writer, strings.Join(parts, " ")); err != nil {
			panic(vm.NewGoError(err))
		}
		return goja.Undefined()
	}
}

func bindBrowser(vm *goja.Runtime, ctx context.Context, browser workflowBrowser) error {
	object := vm.NewObject()
	if err := object.Set("cmd", func(call goja.FunctionCall) goja.Value {
		args := make([]string, len(call.Arguments))
		for index, value := range call.Arguments {
			args[index] = value.String()
		}
		data, err := browser.Run(ctx, args...)
		if err != nil {
			panic(vm.NewGoError(err))
		}
		var result any
		if len(data) > 0 {
			if err := json.Unmarshal(data, &result); err != nil {
				panic(vm.NewGoError(fmt.Errorf("decode agent-browser command result: %w", err)))
			}
		}
		return vm.ToValue(result)
	}); err != nil {
		return fmt.Errorf("initialize browser.cmd: %w", err)
	}
	if err := object.Set("eval", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			panic(vm.NewGoError(fmt.Errorf("browser.eval requires JavaScript")))
		}
		result, err := browser.Eval(ctx, call.Argument(0).String())
		if err != nil {
			panic(vm.NewGoError(err))
		}
		return vm.ToValue(result)
	}); err != nil {
		return fmt.Errorf("initialize browser.eval: %w", err)
	}
	if err := vm.Set("browser", object); err != nil {
		return fmt.Errorf("initialize browser: %w", err)
	}
	return nil
}

func promiseError(value goja.Value) error {
	if object, ok := value.(*goja.Object); ok {
		if stack := object.Get("stack"); stack != nil && !goja.IsUndefined(stack) {
			return errors.New(stack.String())
		}
	}
	return errors.New(value.String())
}
