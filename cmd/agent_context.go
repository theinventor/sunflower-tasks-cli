package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theinventor/sunflower-tasks-cli/internal/config"
	"github.com/theinventor/sunflower-tasks-cli/internal/exitcode"
)

const AgentContextSchemaVersion = "1"

func newAgentContextCmd() *cobra.Command {
	return &cobra.Command{Use: "agent-context", Short: "Emit a versioned machine-readable CLI and API description", Long: "AI agents should call this once at startup. It describes commands, flags, endpoint methods, authentication, exit codes, and saved profiles without scraping terminal help.", RunE: func(cmd *cobra.Command, _ []string) error {
		f, _ := config.Load()
		profiles := []string{}
		if f.Profiles != nil {
			profiles = config.Names(f)
		}
		return printJSON(map[string]any{"schema_version": AgentContextSchemaVersion, "cli_version": Version, "commands": commandContext(cmd.Root()), "global_flags": map[string]any{"--profile": map[string]string{"type": "string", "usage": "saved auth profile for this invocation"}}, "profiles": profiles, "authentication": map[string]any{"environment": "SUNFLOWER_TASKS_API_KEY", "header": "Authorization: Bearer <token>", "account_header": "X-Sunflower-Account-ID"}, "endpoints": endpoints, "exit_codes": exitContext()})
	}}
}
func commandContext(root *cobra.Command) map[string]any {
	out := map[string]any{}
	for _, c := range root.Commands() {
		if c.Name() == "help" || c.Hidden {
			continue
		}
		entry := map[string]any{"summary": c.Short}
		if len(c.Commands()) > 0 {
			entry["subcommands"] = commandContext(c)
		}
		out[c.Name()] = entry
	}
	return out
}
func exitContext() map[string]string {
	m := map[string]string{}
	for _, c := range exitcode.All() {
		m[itoa(c)] = exitcode.Description(c)
	}
	return m
}
func itoa(n int) string { return string(rune('0' + n)) }
