package utils

import "os/exec"

func HasExeOnPath(exe string) bool {
	path, err := exec.LookPath(exe)
	return err == nil && len(path) > 0
}
