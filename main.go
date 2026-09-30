package main

import (
	"errors"
	"github.com/theinventor/sunflower-tasks-cli/cmd"
	"github.com/theinventor/sunflower-tasks-cli/internal/client"
	"os"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		cmd.PrintError(err)
		code := 1
		var he *client.HTTPError
		if errors.As(err, &he) {
			if he.Status == 401 || he.Status == 403 {
				code = 3
			} else if he.Status == 404 {
				code = 4
			} else {
				code = 5
			}
		}
		os.Exit(code)
	}
}
