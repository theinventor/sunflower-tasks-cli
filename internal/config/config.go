package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Profile struct {
	APIURL    string `json:"api_url"`
	APIKey    string `json:"api_key,omitempty"`
	AccountID string `json:"account_id,omitempty"`
}
type File struct {
	DefaultProfile string             `json:"default_profile,omitempty"`
	Profiles       map[string]Profile `json:"profiles"`
}

func path() string {
	if p := os.Getenv("SUNFLOWER_TASKS_CONFIG"); p != "" {
		return p
	}
	d, _ := os.UserConfigDir()
	return filepath.Join(d, "sunflower-tasks", "config.json")
}
func Load() (File, error) {
	b, e := os.ReadFile(path())
	if errors.Is(e, os.ErrNotExist) {
		return File{Profiles: map[string]Profile{}}, nil
	}
	if e != nil {
		return File{}, e
	}
	var f File
	if e = json.Unmarshal(b, &f); e != nil {
		return f, e
	}
	if f.Profiles == nil {
		f.Profiles = map[string]Profile{}
	}
	return f, nil
}
func Save(f File) error {
	p := path()
	if e := os.MkdirAll(filepath.Dir(p), 0700); e != nil {
		return e
	}
	b, e := json.MarshalIndent(f, "", "  ")
	if e != nil {
		return e
	}
	return os.WriteFile(p, append(b, '\n'), 0600)
}
func Names(f File) []string {
	a := make([]string, 0, len(f.Profiles))
	for n := range f.Profiles {
		a = append(a, n)
	}
	sort.Strings(a)
	return a
}
func Resolve(name string) (string, Profile, error) {
	if p := os.Getenv("SUNFLOWER_TASKS_API_KEY"); p != "" {
		u := os.Getenv("SUNFLOWER_TASKS_API_URL")
		if u == "" {
			u = "https://tasks.sunflower-vacations.com"
		}
		return "env", Profile{APIURL: u, APIKey: p, AccountID: os.Getenv("SUNFLOWER_TASKS_ACCOUNT_ID")}, nil
	}
	f, e := Load()
	if e != nil {
		return "", Profile{}, e
	}
	if name == "" {
		name = f.DefaultProfile
	}
	if name == "" && len(f.Profiles) == 1 {
		for n := range f.Profiles {
			name = n
		}
	}
	p, ok := f.Profiles[name]
	if !ok {
		return "", Profile{}, fmt.Errorf("no saved profile; run `sunflower-tasks auth login` or set SUNFLOWER_TASKS_API_KEY")
	}
	if p.APIURL == "" {
		p.APIURL = "https://tasks.sunflower-vacations.com"
	}
	return name, p, nil
}
