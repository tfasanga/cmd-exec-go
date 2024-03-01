package exec

import (
	"fmt"
	"strings"
)

type testExecutionContext struct {
	user string
	host string
	port int
}

func (rc *testExecutionContext) IsLocal() bool {
	return rc.host == "localhost"
}

// ExecuteCmd implements Machine
func (rc *testExecutionContext) ExecuteCmd(_io CommandInOut, _dir string, command string, arg ...string) (string, error) {
	cmd := rc.buildCmd(command, arg...)
	return cmd, nil
}

// RunCmd implements Machine
func (rc *testExecutionContext) RunCmd(io CommandInOut, _dir string, command string, arg ...string) error {
	cmd := rc.buildCmd(command, arg...)
	_, err := fmt.Fprintf(io.Log(), "%s", cmd)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(io.Out(), "%s", cmd)
	if err != nil {
		return err
	}
	return nil
}

// User implements Machine
func (rc *testExecutionContext) User() string {
	return rc.user
}

// Host implements Machine
func (rc *testExecutionContext) Host() string {
	return rc.host
}

// IpAddr implements Machine
func (rc *testExecutionContext) IpAddr() string {
	return rc.host
}

// Port implements Machine
func (rc *testExecutionContext) Port() int {
	return rc.port
}

func (rc *testExecutionContext) buildCmd(command string, arg ...string) string {
	cmd := buildCmd(command, arg...)
	if rc.IsLocal() {
		return cmd
	} else {
		return wrapSsh(rc, cmd)
	}
}

func buildCmd(command string, arg ...string) string {
	cmd := append([]string{command}, arg...)
	return strings.Join(cmd, " ")
}

func wrapSsh(rc *testExecutionContext, cmd string) string {
	if rc.port == 0 || rc.port == 22 {
		return fmt.Sprintf("ssh %s@%s -- %s",
			rc.user, rc.host, cmd)
	} else {
		return fmt.Sprintf("ssh %s@%s -p %d -- %s",
			rc.user, rc.host, rc.port, cmd)
	}
}
