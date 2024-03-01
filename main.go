package main

import (
	"github.com/tfasanga/cmd-exec/exec"
	"os"
)

func main() {
	machine := exec.NewLocalMachine("test")
	inOut := exec.NewCommandInOut(os.Stdout, os.Stderr, nil, os.Stdin)
	err := machine.RunCmd(inOut, "", "ls", "-l")
	if err != nil {
		println(err)
		os.Exit(1)
	}
}
