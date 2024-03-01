package exec

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"os/exec"
)

func IsLocal(m Machine) bool {
	return IsLocalHost(m.Host())
}

func IsLocalHost(host string) bool {
	return host == "localhost"
}

type localExecutionContext struct {
	localUser string
}

func (rc *localExecutionContext) String() string {
	return fmt.Sprintf("local")
}

// ExecuteCmd implements Machine
func (rc *localExecutionContext) ExecuteCmd(io CommandInOut, dir, command string, arg ...string) (string, error) {
	return localExec(io, dir, command, arg...)
}

// RunCmd implements Machine
func (rc *localExecutionContext) RunCmd(io CommandInOut, dir, command string, arg ...string) error {
	return localRun(io, dir, command, arg...)
}

// User implements Machine
func (rc *localExecutionContext) User() string {
	return rc.localUser
}

// Host implements Machine
func (rc *localExecutionContext) Host() string {
	return "localhost"
}

// IpAddr implements Machine
func (rc *localExecutionContext) IpAddr() string {
	return "localhost"
}

// Port implements Machine
func (rc *localExecutionContext) Port() int {
	return 22
}

func NewLocalMachine(user string) Machine {
	localMachine := &localExecutionContext{
		localUser: user,
	}
	return localMachine
}

func NewSshMachine(hostname string, port int, sshConfig *ssh.ClientConfig) Machine {
	sshMachine := &sshExecutionContext{
		host:      hostname,
		port:      port,
		sshConfig: sshConfig,
	}
	return sshMachine
}

func localExec(io CommandInOut, dir, command string, arg ...string) (string, error) {
	cmd := exec.Command(command, arg...)
	cmd.Dir = dir
	if io.In() != nil {
		cmd.Stdin = io.In()
	} else {
		cmd.Stdin = os.Stdin
	}

	logCommand(io, "localhost", "", command, arg...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: failed to run local command", err)
	}

	return string(output), nil
}

func localRun(io CommandInOut, dir, command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	cmd.Dir = dir
	if io.Out() != nil {
		cmd.Stdout = io.Out()
	}
	if io.Err() != nil {
		cmd.Stderr = io.Err()
	}
	if io.In() != nil {
		cmd.Stdin = io.In()
	}

	logCommand(io, "localhost", "", command, arg...)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%w: failed to run local command", err)
	}

	return nil
}
