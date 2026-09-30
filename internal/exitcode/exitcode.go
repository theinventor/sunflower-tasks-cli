package exitcode

import "fmt"

const (
	Success  = 0
	Usage    = 2
	Auth     = 3
	NotFound = 4
	API      = 5
	Network  = 6
	Config   = 7
)

func Description(code int) string {
	switch code {
	case Success:
		return "success"
	case Usage:
		return "invalid command or arguments"
	case Auth:
		return "authentication required or rejected"
	case NotFound:
		return "resource not found"
	case API:
		return "API request failed"
	case Network:
		return "network or transport failure"
	case Config:
		return "local configuration failure"
	}
	return fmt.Sprintf("exit code %d", code)
}
func All() []int { return []int{Success, Usage, Auth, NotFound, API, Network, Config} }
