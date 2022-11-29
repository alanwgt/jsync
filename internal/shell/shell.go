package shell

import "os/exec"

func Exec(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	return c.Run()
}
