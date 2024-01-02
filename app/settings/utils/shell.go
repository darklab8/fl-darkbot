package utils

import (
	"darkbot/app/settings/darkbot_logus"
	"fmt"
	"os"
	"os/exec"
)

func ShellRunArgs(program string, args ...string) {
	darkbot_logus.Log.Debug(fmt.Sprintf("OK attempting to run: %s", program), darkbot_logus.Args(args))
	executable, _ := exec.LookPath(program)

	args = append([]string{""}, args...)
	command := exec.Cmd{
		Path:   executable,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	err := command.Run()

	darkbot_logus.Log.CheckFatal(err, "failed to run shell command")
}
