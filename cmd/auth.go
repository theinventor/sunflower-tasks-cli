package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/theinventor/sunflower-tasks-cli/internal/client"
	"github.com/theinventor/sunflower-tasks-cli/internal/config"
	"net/http"
	"os"
	"strings"
	"time"
)

func newAuthCmd() *cobra.Command {
	a := &cobra.Command{Use: "auth", Short: "Manage saved API profiles"}
	a.AddCommand(newAuthLogin(), newAuthSave(), newAuthList(), newAuthStatus(), newAuthUse(), newAuthLogout())
	return a
}
func newAuthLogin() *cobra.Command {
	var name, url, email, password, otp, device, account string
	c := &cobra.Command{Use: "login", Short: "Sign in and save a bearer token", Long: "Calls POST /api/mobile/v1/session. The password is read from --password or SUNFLOWER_TASKS_PASSWORD; it is never printed.", RunE: func(*cobra.Command, []string) error {
		if name == "" {
			name = email
		}
		if name == "" || email == "" {
			return fmt.Errorf("--profile and --email are required")
		}
		if password == "" {
			password = os.Getenv("SUNFLOWER_TASKS_PASSWORD")
		}
		if password == "" {
			return fmt.Errorf("provide --password or set SUNFLOWER_TASKS_PASSWORD")
		}
		if url == "" {
			url = "https://tasks.sunflower-vacations.com"
		}
		payload := map[string]any{"email": email, "password": password}
		if otp != "" {
			payload["otp_attempt"] = otp
		}
		if device != "" {
			payload["device_name"] = device
		}
		if account != "" {
			payload["account_id"] = account
		}
		b, _ := json.Marshal(payload)
		cl := &client.Client{BaseURL: strings.TrimRight(url, "/"), HTTP: &http.Client{Timeout: 45 * time.Second}}
		_, body, e := cl.Do("POST", "/api/mobile/v1/session", b, nil)
		if e != nil {
			return e
		}
		var reply struct {
			Token string `json:"token"`
		}
		if e = json.Unmarshal(body, &reply); e != nil || reply.Token == "" {
			return fmt.Errorf("login response did not include a token")
		}
		f, e := config.Load()
		if e != nil {
			return e
		}
		f.Profiles[name] = config.Profile{APIURL: strings.TrimRight(url, "/"), APIKey: reply.Token, AccountID: account}
		if f.DefaultProfile == "" {
			f.DefaultProfile = name
		}
		if e = config.Save(f); e != nil {
			return e
		}
		fmt.Printf("logged in as %s (token hidden)\n", email)
		return nil
	}}
	c.Flags().StringVar(&name, "profile", "", "saved profile name")
	c.Flags().StringVar(&url, "api-url", "", "API base URL")
	c.Flags().StringVar(&email, "email", "", "account email")
	c.Flags().StringVar(&password, "password", "", "password (prefer env)")
	c.Flags().StringVar(&otp, "otp-attempt", "", "two-factor code")
	c.Flags().StringVar(&device, "device-name", "sunflower-tasks-cli", "device label")
	c.Flags().StringVar(&account, "account-id", "", "account to select")
	_ = c.MarkFlagRequired("profile")
	_ = c.MarkFlagRequired("email")
	return c
}
func newAuthSave() *cobra.Command {
	var name, url, key, account string
	c := &cobra.Command{Use: "save", Short: "Save an API token without echoing it", RunE: func(*cobra.Command, []string) error {
		if name == "" {
			return fmt.Errorf("--profile is required")
		}
		if key == "" {
			key = os.Getenv("SUNFLOWER_TASKS_API_KEY")
		}
		if key == "" {
			return fmt.Errorf("provide --api-key or set SUNFLOWER_TASKS_API_KEY")
		}
		if url == "" {
			url = "https://tasks.sunflower-vacations.com"
		}
		f, e := config.Load()
		if e != nil {
			return e
		}
		f.Profiles[name] = config.Profile{APIURL: strings.TrimRight(url, "/"), APIKey: key, AccountID: account}
		if f.DefaultProfile == "" {
			f.DefaultProfile = name
		}
		if e = config.Save(f); e != nil {
			return e
		}
		fmt.Printf("saved profile %q (token hidden)\n", name)
		return nil
	}}
	c.Flags().StringVar(&name, "profile", "", "profile name")
	c.Flags().StringVar(&url, "api-url", "", "API base URL")
	c.Flags().StringVar(&key, "api-key", "", "bearer token (prefer env)")
	c.Flags().StringVar(&account, "account-id", "", "default account ID")
	_ = c.MarkFlagRequired("profile")
	return c
}
func newAuthList() *cobra.Command {
	return &cobra.Command{Use: "list", Short: "List saved profile names", RunE: func(*cobra.Command, []string) error {
		f, e := config.Load()
		if e != nil {
			return e
		}
		for _, n := range config.Names(f) {
			mark := " "
			if n == f.DefaultProfile {
				mark = "*"
			}
			fmt.Printf("%s %s\n", mark, n)
		}
		return nil
	}}
}
func newAuthStatus() *cobra.Command {
	var name string
	c := &cobra.Command{Use: "status", Short: "Show profile metadata without its token", RunE: func(*cobra.Command, []string) error {
		f, e := config.Load()
		if e != nil {
			return e
		}
		if name == "" {
			name = f.DefaultProfile
		}
		p, ok := f.Profiles[name]
		if !ok {
			return fmt.Errorf("profile %q not found", name)
		}
		return printJSON(map[string]any{"profile": name, "api_url": p.APIURL, "account_id": p.AccountID, "token_configured": p.APIKey != "", "default": name == f.DefaultProfile})
	}}
	c.Flags().StringVar(&name, "profile", "", "profile name")
	return c
}
func newAuthUse() *cobra.Command {
	return &cobra.Command{Use: "use PROFILE", Short: "Set the default profile", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, a []string) error {
		f, e := config.Load()
		if e != nil {
			return e
		}
		if _, ok := f.Profiles[a[0]]; !ok {
			return fmt.Errorf("profile %q not found", a[0])
		}
		f.DefaultProfile = a[0]
		if e = config.Save(f); e != nil {
			return e
		}
		fmt.Printf("default profile: %s\n", a[0])
		return nil
	}}
}
func newAuthLogout() *cobra.Command {
	return &cobra.Command{Use: "logout PROFILE", Short: "Delete a saved profile", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, a []string) error {
		f, e := config.Load()
		if e != nil {
			return e
		}
		if _, ok := f.Profiles[a[0]]; !ok {
			return fmt.Errorf("profile %q not found", a[0])
		}
		delete(f.Profiles, a[0])
		if f.DefaultProfile == a[0] {
			f.DefaultProfile = ""
		}
		if e = config.Save(f); e != nil {
			return e
		}
		fmt.Printf("removed profile %s\n", a[0])
		return nil
	}}
}
