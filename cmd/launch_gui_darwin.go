package cmd

func buildLaunchGuiCmd(cmd string, args ...string) (string, []string) {
	newArgs := make([]string, len(args)+3)
	newArgs[0] = "-a"
	newArgs[1] = cmd
	newArgs[2] = "--args"
	copy(newArgs[3:], args)
	return "open", newArgs
}
