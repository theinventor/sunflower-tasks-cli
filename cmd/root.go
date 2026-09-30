package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"runtime/debug"
)

var Version = "dev"
var Commit = ""
var BuildDate = ""
var profile string

func init() {
	if info, ok := debug.ReadBuildInfo(); ok && Version == "dev" && info.Main.Version != "" && info.Main.Version != "(devel)" {
		Version = info.Main.Version
	}
}
func NewRootCmd() *cobra.Command {
	r := &cobra.Command{Use: "sunflower-tasks", Aliases: []string{"stasks", "st"}, Short: "Sunflower Tasks API CLI for humans and AI agents", Long: "Operate Sunflower Tasks through the documented mobile API. Commands emit JSON on stdout for reliable piping; diagnostics go to stderr.", Version: Version, SilenceUsage: true, SilenceErrors: true}
	r.PersistentFlags().StringVar(&profile, "profile", "", "saved auth profile for this invocation")
	r.AddCommand(newAgentContextCmd(), newAuthCmd(), newAPICommand(), newWhoamiCmd(), newUpdateCmd())
	r.AddCommand(newCompletionCmd())
	return r
}
func PrintError(e error) { fmt.Fprintln(os.Stderr, "sunflower-tasks:", e) }
func printJSON(v any) error {
	b, e := json.MarshalIndent(v, "", "  ")
	if e != nil {
		return e
	}
	fmt.Println(string(b))
	return nil
}
