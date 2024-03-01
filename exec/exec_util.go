package exec

import (
	"golang.org/x/crypto/ssh"
	"net"
	"os"
)

//goland:noinspection GoUnusedExportedFunction
func PublicKey(path string) (ssh.AuthMethod, error) {
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func fileTest(machine Machine, io CommandInOut, fileName string, option string) (bool, error) {
	err := machine.RunCmd(io, "", "test", option, fileName)
	if err != nil {
		switch e := err.(type) {
		case *ssh.ExitError:
			if e.ExitStatus() != 0 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func resolveIpAddr(host string) string {
	if host == "localhost" {
		return host
	}

	ipAddrs, err := net.LookupHost(host)
	if err != nil {
		return host
	}

	if len(ipAddrs) == 0 {
		return host
	}

	return ipAddrs[0]

}
