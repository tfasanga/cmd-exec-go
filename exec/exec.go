package exec

import (
	"fmt"
)

type CommandExecutor interface {
	ExecuteCmd(io CommandInOut, dir, command string, arg ...string) (string, error)
	RunCmd(io CommandInOut, dir, command string, arg ...string) error
}

type Machine interface {
	CommandExecutor
	User() string
	Host() string
	IpAddr() string
	Port() int
}

func Scp(io CommandInOut, sourceMachine Machine, sourceFile string, destinationMachine Machine, destinationFile string) error {
	if destinationMachine.Host() == sourceMachine.Host() {
		_, err := fmt.Fprintf(io.Out(), "Skipping, source and destination are the same: %s\n", destinationMachine.Host())
		if err != nil {
			return err
		}
		return nil
	}

	var execMachine Machine
	if IsLocal(destinationMachine) {
		execMachine = destinationMachine
	} else {
		execMachine = sourceMachine
	}

	cmd := "scp"

	var source, destination string
	if IsLocal(sourceMachine) || sourceMachine == execMachine {
		source = sourceFile
	} else {
		source = fmt.Sprintf("%s@%s:%s", sourceMachine.User(), sourceMachine.IpAddr(), sourceFile)
	}
	if IsLocal(destinationMachine) || destinationMachine == execMachine {
		destination = destinationFile
	} else {
		destination = fmt.Sprintf("%s@%s:%s", destinationMachine.User(), destinationMachine.IpAddr(), destinationFile)
	}
	args := []string{source, destination}

	return execMachine.RunCmd(io, "", cmd, args...)
}

func Rsync(io CommandInOut, sourceMachine Machine, sourceRootDir, sourceRelativeDir string, destinationMachine Machine, destinationRootDir string, excluded []string) error {
	if destinationMachine.Host() == sourceMachine.Host() {
		_, err := fmt.Fprintf(io.Out(), "Skipping, source and destination are the same: %s\n", destinationMachine.Host())
		if err != nil {
			return err
		}
		return nil
	}
	if IsLocal(destinationMachine) {
		return fmt.Errorf("remote machine cannot be %s", destinationMachine.Host())
	}
	cmd, args := buildRsyncCmdAndArgs(sourceRootDir, sourceRelativeDir, destinationMachine, destinationRootDir, excluded)
	return sourceMachine.RunCmd(io, "", cmd, args...)
}

func Mkdirs(machine Machine, io CommandInOut, dirName string) error {
	return machine.RunCmd(io, "", "mkdir", "-p", dirName)
}

func FileExists(machine Machine, io CommandInOut, fileName string) (bool, error) {
	return fileTest(machine, io, fileName, "-f")
}

func DirectoryExists(machine Machine, io CommandInOut, fileName string) (bool, error) {
	return fileTest(machine, io, fileName, "-d")
}
