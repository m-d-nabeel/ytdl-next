package main

import (
	"os/exec"
)

type CmdInterface interface {
	GetYTMediaCmd() *exec.Cmd
}

func (a *API) GetYTMediaCmd() *exec.Cmd {
	cmd := exec.Command("yt-dlp")
	cmd.Args = append(cmd.Args, "-j")
	cmd.Args = append(cmd.Args, "--skip-download")
	return cmd
}
