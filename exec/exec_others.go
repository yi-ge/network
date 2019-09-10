// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package exec

// Command is part of the Interface interface.
func parseCommand(cmd string, args ...string) (string, []string) {
	return cmd, args
}
