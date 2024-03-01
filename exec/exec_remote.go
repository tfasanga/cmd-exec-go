package exec

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

type sshExecutionContext struct {
	host      string
	ipAddr    string
	port      int
	sshConfig *ssh.ClientConfig
}

func (rc *sshExecutionContext) String() string {
	return fmt.Sprintf("%s@%s -p %d", rc.sshConfig.User, rc.host, rc.port)
}

// ExecuteCmd implements Machine
func (rc *sshExecutionContext) ExecuteCmd(io CommandInOut, dir, command string, arg ...string) (string, error) {
	if rc.host == "" {
		return "", fmt.Errorf("cannot execute ssh command, the remote hostname is not set")
	}
	serverAddress := joinHostPort(rc.host, rc.port)
	return remoteExec(io, serverAddress, rc.sshConfig, dir, command, arg...)
}

// RunCmd implements Machine
func (rc *sshExecutionContext) RunCmd(io CommandInOut, dir, command string, arg ...string) error {
	serverAddress := joinHostPort(rc.host, rc.port)
	return remoteRun(serverAddress, rc.sshConfig, io, dir, command, arg...)
}

// User implements Machine
func (rc *sshExecutionContext) User() string {
	return rc.sshConfig.User
}

// Host implements Machine
func (rc *sshExecutionContext) Host() string {
	return rc.host
}

// IpAddr implements Machine
func (rc *sshExecutionContext) IpAddr() string {
	// resolve IP lazily
	if rc.ipAddr == "" {
		rc.ipAddr = resolveIpAddr(rc.host)
	}
	return rc.ipAddr
}

// Port implements Machine
func (rc *sshExecutionContext) Port() int {
	return rc.port
}

func remoteExec(io CommandInOut, serverAddress string, sshConfig *ssh.ClientConfig, dir, command string, arg ...string) (string, error) {
	// Establish an SSH connection
	sshClient, err := ssh.Dial("tcp", serverAddress, sshConfig)
	if err != nil {
		return "", fmt.Errorf("%w: failed to establish SSH connection", err)
	}

	defer sshClient.Close()

	// Create a session on the SSH connection
	session, err := sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("%w: failed to create SSH session", err)
	}

	if io.In() != nil {
		session.Stdin = io.In()
	}

	defer session.Close()

	// Run the remote command
	var actualCmd string
	if dir == "" {
		actualCmd = command
	} else {
		actualCmd = fmt.Sprintf("cd %s && %s", dir, command)
	}

	logCommand(io, serverAddress, "", actualCmd, arg...)

	runCmd := actualCmd
	if len(arg) > 0 {
		runCmd = fmt.Sprintf("%s %s", actualCmd, strings.Join(arg, " "))
	}

	output, err := session.CombinedOutput(runCmd)
	if err != nil {
		logCommand(io, serverAddress, "ERR", actualCmd, arg...)

		return "", fmt.Errorf("%w: failed to run remote command. Output: %s", err, string(output))
	}

	logCommand(io, serverAddress, "OK", actualCmd, arg...)

	return string(output), nil
}

func remoteRun(serverAddress string, sshConfig *ssh.ClientConfig, io CommandInOut, dir, command string, arg ...string) error {
	// Establish an SSH connection
	sshClient, err := ssh.Dial("tcp", serverAddress, sshConfig)
	if err != nil {
		return fmt.Errorf("%w: failed to establish SSH connection", err)
	}

	defer sshClient.Close()

	// Create a session on the SSH connection
	session, err := sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("%w: failed to create SSH session", err)
	}

	defer session.Close()

	// Run the remote command
	var actualCmd string
	if dir == "" {
		actualCmd = command
	} else {
		actualCmd = fmt.Sprintf("cd %s && %s", dir, command)
	}

	logCommand(io, serverAddress, "", actualCmd, arg...)

	runCmd := actualCmd
	if len(arg) > 0 {
		runCmd = fmt.Sprintf("%s %s", actualCmd, strings.Join(arg, " "))
	}

	if io.Err() != nil {
		session.Stderr = io.Err()
	}
	if io.Out() != nil {
		session.Stdout = io.Out()
	}
	if io.In() != nil {
		session.Stdin = io.In()
	}

	if err := session.Run(runCmd); err != nil {
		logCommand(io, serverAddress, "ERR", actualCmd, arg...)
		return fmt.Errorf("%w when executing command: [%s]", err, actualCmd)
	}

	logCommand(io, serverAddress, "OK", actualCmd, arg...)

	return nil
}
