package main

import (
	"sort"
	"strings"

	"github.com/vaayne/tap"
	"github.com/vaayne/tap/script"
)

func parseArgs(raw []string) map[string]string {
	args := make(map[string]string)
	for _, s := range raw {
		if k, v, ok := strings.Cut(s, "="); ok {
			args[k] = v
		}
	}
	return args
}

func formatArgHints(s *script.Script, color bool) string {
	if len(s.Meta.Args) == 0 {
		return ""
	}
	var parts []string
	for name, def := range s.Meta.Args {
		if def.Required {
			parts = append(parts, yellow(color, name+"*"))
		} else {
			parts = append(parts, dim(color, name))
		}
	}
	sort.Strings(parts) // consistent ordering
	return " [" + strings.Join(parts, ", ") + "]"
}

// groupScripts groups scripts by their site prefix (part before /).
func groupScripts(scripts []*script.Script) map[string][]*script.Script {
	groups := make(map[string][]*script.Script)
	for _, s := range scripts {
		group := s.Meta.Name
		if idx := strings.Index(group, "/"); idx != -1 {
			group = group[:idx]
		}
		groups[group] = append(groups[group], s)
	}
	return groups
}

func scriptNames(client *tap.Client) []string {
	scripts := client.ListScripts()
	names := make([]string, len(scripts))
	for i, s := range scripts {
		names[i] = s.Meta.Name
	}
	return names
}
