package main

import (
	"os/exec"
)

type CmdInterface interface {
	GetYTMediaInfoCmd() *exec.Cmd
}

func (a *API) GetYTMediaInfoCmd() *exec.Cmd {
	cmd := exec.Command("yt-dlp")
	cmd.Args = append(cmd.Args, "-j")
	cmd.Args = append(cmd.Args, "--skip-download")
	return cmd
}
