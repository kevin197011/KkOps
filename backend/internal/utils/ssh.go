// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package utils

import (
	"context"
	"io"
	"net"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHClient wraps an SSH client connection
type SSHClient struct {
	client *ssh.Client
}

// NewSSHClient creates a new SSH client connection with private key
func NewSSHClient(host string, port int, user string, privateKey []byte, timeout time.Duration) (*SSHClient, error) {
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use proper host key verification
		Timeout:         timeout,
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return &SSHClient{client: client}, nil
}

// NewSSHClientWithPassphrase creates a new SSH client connection with passphrase-protected private key
func NewSSHClientWithPassphrase(host string, port int, user string, privateKey []byte, passphrase []byte, timeout time.Duration) (*SSHClient, error) {
	var signer ssh.Signer
	var err error

	if len(passphrase) > 0 {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(privateKey, passphrase)
	} else {
		signer, err = ssh.ParsePrivateKey(privateKey)
	}
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use proper host key verification
		Timeout:         timeout,
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return &SSHClient{client: client}, nil
}

// NewSSHClientWithPassword creates a new SSH client connection with password authentication
func NewSSHClientWithPassword(host string, port int, user string, password string, timeout time.Duration) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use proper host key verification
		Timeout:         timeout,
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return &SSHClient{client: client}, nil
}

// ExecuteCommand executes a command on the remote host
func (c *SSHClient) ExecuteCommand(command string) (string, int, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", -1, err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*ssh.ExitError); ok {
			exitCode = exitError.ExitStatus()
		} else {
			return "", -1, err
		}
	}

	return string(output), exitCode, nil
}

// ExecuteCommandWithTimeout executes a command with a timeout
func (c *SSHClient) ExecuteCommandWithTimeout(ctx context.Context, command string) (string, int, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", -1, err
	}
	defer session.Close()

	// Create a channel to receive the result
	type result struct {
		output   []byte
		exitCode int
		err      error
	}
	resultCh := make(chan result, 1)

	go func() {
		output, err := session.CombinedOutput(command)
		exitCode := 0
		if err != nil {
			if exitError, ok := err.(*ssh.ExitError); ok {
				exitCode = exitError.ExitStatus()
				err = nil // Not a fatal error, just non-zero exit
			}
		}
		resultCh <- result{output: output, exitCode: exitCode, err: err}
	}()

	select {
	case <-ctx.Done():
		// Timeout or cancellation - try to close session to kill the command
		session.Signal(ssh.SIGKILL)
		return "", -1, ctx.Err()
	case res := <-resultCh:
		if res.err != nil {
			return "", -1, res.err
		}
		return string(res.output), res.exitCode, nil
	}
}

// ExecuteCommandWithStream executes a command and streams output
func (c *SSHClient) ExecuteCommandWithStream(command string, stdout, stderr io.Writer) (int, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return -1, err
	}
	defer session.Close()

	session.Stdout = stdout
	session.Stderr = stderr

	err = session.Run(command)
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*ssh.ExitError); ok {
			exitCode = exitError.ExitStatus()
		} else {
			return -1, err
		}
	}

	return exitCode, nil
}

// Close closes the SSH connection
func (c *SSHClient) Close() error {
	return c.client.Close()
}

// Client returns the underlying SSH client
func (c *SSHClient) Client() *ssh.Client {
	return c.client
}
