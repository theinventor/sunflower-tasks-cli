package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type releaseAsset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}
type releaseInfo struct {
	Tag    string         `json:"tag_name"`
	Assets []releaseAsset `json:"assets"`
}

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{Use: "update", Short: "Check for and install a newer release", Long: "Downloads the matching release from GitHub. Development builds never replace a source checkout.", RunE: func(*cobra.Command, []string) error {
		if Version == "dev" {
			return fmt.Errorf("development build cannot self-update; install a release or use `go install github.com/theinventor/sunflower-tasks-cli@latest`")
		}
		resp, e := http.Get("https://api.github.com/repos/theinventor/sunflower-tasks-cli/releases/latest")
		if e != nil {
			return e
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			return fmt.Errorf("GitHub release lookup: %s", resp.Status)
		}
		var rel releaseInfo
		if e = json.NewDecoder(resp.Body).Decode(&rel); e != nil {
			return e
		}
		if rel.Tag == Version || rel.Tag == "v"+strings.TrimPrefix(Version, "v") {
			fmt.Printf("already up to date (%s)\n", Version)
			return nil
		}
		want := fmt.Sprintf("sunflower-tasks_%s_%s_%s.tar.gz", strings.TrimPrefix(rel.Tag, "v"), runtime.GOOS, runtime.GOARCH)
		var u string
		for _, a := range rel.Assets {
			if a.Name == want {
				u = a.URL
			}
		}
		if u == "" {
			return fmt.Errorf("release %s has no asset %s", rel.Tag, want)
		}
		tmp, e := os.CreateTemp("", "sunflower-tasks-update-*.tar.gz")
		if e != nil {
			return e
		}
		defer os.Remove(tmp.Name())
		rr, e := http.Get(u)
		if e != nil {
			return e
		}
		defer rr.Body.Close()
		if rr.StatusCode >= 300 {
			return fmt.Errorf("download: %s", rr.Status)
		}
		if _, e = io.Copy(tmp, rr.Body); e != nil {
			return e
		}
		tmp.Close()
		dir := filepath.Dir(os.Args[0])
		if e = exec.Command("tar", "-xzf", tmp.Name(), "-C", dir).Run(); e != nil {
			return fmt.Errorf("extract release: %w", e)
		}
		fmt.Printf("updated to %s\n", rel.Tag)
		return nil
	}}
}
