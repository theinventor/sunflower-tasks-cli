package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/theinventor/sunflower-tasks-cli/internal/client"
)

func newWhoamiCmd() *cobra.Command {
	return &cobra.Command{Use: "whoami", Short: "Show the authenticated Sunflower Tasks identity", RunE: func(*cobra.Command, []string) error {
		c, e := client.New(profile)
		if e != nil {
			return e
		}
		_, b, e := c.Do("GET", "/api/mobile/v1/profile", nil, nil)
		if e != nil {
			return e
		}
		var v any
		if json.Unmarshal(b, &v) == nil {
			return printJSON(v)
		}
		fmt.Print(string(b))
		return nil
	}}
}
