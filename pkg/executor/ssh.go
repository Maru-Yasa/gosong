package executor

import "golang.org/x/crypto/ssh"

type SSHExecutor struct {
	Client *ssh.Client
}

func (s *SSHExecutor) Run(cmd string) (string, error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return "", err
	}

	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	return string(output), err
}
