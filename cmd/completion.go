package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	return &cobra.Command{Use: "completion [bash|zsh|fish|powershell]", Short: "Generate shell completion", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, a []string) error {
		switch a[0] {
		case "bash":
			return c.Root().GenBashCompletion(c.OutOrStdout())
		case "zsh":
			return c.Root().GenZshCompletion(c.OutOrStdout())
		case "fish":
			return c.Root().GenFishCompletion(c.OutOrStdout(), true)
		case "powershell":
			return c.Root().GenPowerShellCompletion(c.OutOrStdout())
		}
		return fmt.Errorf("unsupported shell %q", a[0])
	}}
}
