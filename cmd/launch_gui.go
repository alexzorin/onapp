// +build !darwin

package cmd

func buildLaunchGuiCmd(cmd string, args ...string) (string, []string) {
	return cmd, args
}
