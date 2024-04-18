package exec

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func buildRsyncCmdAndArgs(sourceRootDir string, sourceRelativeDir string, to Machine, destinationRootDir string, options []string) (string, []string) {
	cmd := "rsync"
	args := buildRsyncArgs(sourceRootDir, sourceRelativeDir, to, destinationRootDir, options)
	return cmd, args
}

func buildRsyncArgs(sourceRootDir string, sourceRelativeDir string, to Machine, destinationRootDir string, options []string) []string {
	// must append "/./" in order to copy relative paths, see man for "rsync -R"
	source := fmt.Sprintf("%s/./%s", sourceRootDir, sourceRelativeDir)
	destination := fmt.Sprintf("%s@%s:%s", to.User(), to.IpAddr(), destinationRootDir)
	return append(options, source, destination)
}

func joinHostPort(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}

func logCommand(io CommandInOut, host string, flag string, command string, args ...string) {
	log := io.Log()
	if log == nil {
		return
	}

	joined := strings.Join(args, " ")

	if len(flag) == 0 {
		fmt.Fprintf(log, "[%s] %s %s\n", host, command, joined)

	} else {
		fmt.Fprintf(log, "[%s] %s %s [%s]\n", host, command, joined, flag)
	}

}
