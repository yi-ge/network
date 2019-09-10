// +build windows

package exec

import "strings"

// Command is part of the Interface interface.
func parseCommand(cmd string, args ...string) (string, []string) {
	return "cmd", []string{"/d/s/c", "chcp 65001 >nul & " + cmd + " " + strings.Join(args, " ")}
}
