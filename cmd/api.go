package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/theinventor/sunflower-tasks-cli/internal/client"
	"net/url"
	"os"
	"strings"
)

func newAPICommand() *cobra.Command {
	api := &cobra.Command{Use: "api", Short: "Inspect and call the Sunflower Tasks API"}
	api.AddCommand(newEndpointsCmd(), newRequestCmd())
	return api
}
func newEndpointsCmd() *cobra.Command {
	var jsonOut bool
	c := &cobra.Command{Use: "endpoints [search]", Args: cobra.MaximumNArgs(1), Short: "List every documented mobile API endpoint", Long: "Print the complete /api/mobile/v1 contract. Pass a search term to filter by path or summary. Use --json for agent-friendly output.", RunE: func(cmd *cobra.Command, args []string) error {
		needle := ""
		if len(args) > 0 {
			needle = strings.ToLower(args[0])
		}
		out := []Endpoint{}
		for _, e := range endpoints {
			if needle == "" || strings.Contains(strings.ToLower(e.Path), needle) || strings.Contains(strings.ToLower(e.Summary), needle) {
				out = append(out, e)
			}
		}
		if jsonOut {
			return printJSON(out)
		}
		for _, e := range out {
			fmt.Printf("%-6s %-62s %s\n", e.Method, e.Path, e.Summary)
		}
		fmt.Fprintf(os.Stderr, "%d endpoint operations\n", len(out))
		return nil
	}}
	c.Flags().BoolVar(&jsonOut, "json", false, "emit JSON")
	return c
}
func newRequestCmd() *cobra.Command {
	var data, query string
	c := &cobra.Command{Use: "request METHOD PATH", Short: "Call any documented or forward-compatible API path", Args: cobra.ExactArgs(2), Example: "sunflower-tasks api request GET /api/mobile/v1/properties\nsunflower-tasks api request PATCH /api/mobile/v1/profile --data '{\"preferred_language\":\"es\"}'", RunE: func(cmd *cobra.Command, args []string) error {
		cl, e := client.New(profile)
		if e != nil {
			return e
		}
		var body []byte
		if data != "" {
			if data == "@-" {
				body, e = os.ReadFile("/dev/stdin")
			} else if strings.HasPrefix(data, "@") {
				body, e = os.ReadFile(strings.TrimPrefix(data, "@"))
			} else {
				body = []byte(data)
			}
			if e != nil {
				return e
			}
		}
		q := url.Values{}
		if query != "" {
			q, e = url.ParseQuery(query)
			if e != nil {
				return fmt.Errorf("invalid --query: %w", e)
			}
		}
		status, b, e := cl.Do(strings.ToUpper(args[0]), args[1], body, q)
		if e != nil {
			if b != nil {
				fmt.Fprintln(os.Stderr, string(b))
			}
			return e
		}
		fmt.Fprintln(os.Stderr, fmt.Sprintf("HTTP %d", status))
		var v any
		if json.Unmarshal(b, &v) == nil {
			return printJSON(v)
		}
		fmt.Print(string(b))
		return nil
	}}
	c.Flags().StringVar(&data, "data", "", "JSON body, @file, or @- for stdin")
	c.Flags().StringVar(&query, "query", "", "URL query string (key=value&key2=value2)")
	return c
}
