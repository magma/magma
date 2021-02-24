package servicers

import (
	"bytes"
	"fmt"
	"github.com/golang/glog"
	"golang.org/x/crypto/ssh"
	"time"
)

// restartIperfServer is a best efort function to restart iperf server, Will try to
// access the iperf server through ssh using the address and port defined in uesim.
// If ssh is not accessible it will just print error.

func restartIperfServer(address, port string) error {
	client, session, err := connectToHost("vagrant", "vagrant", address, port)
	if err != nil {
		glog.Errorf("Iperf server not restarted. Could not connect iperf server at %s: %s", trafficSrvIP, err)
		return err
	}
	defer client.Close()
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		glog.Errorf("Iperf server not restarted. Could not create StdinPipe at %s %s", trafficSrvIP, err)
		return err
	}

	// store output in variable
	var b bytes.Buffer
	session.Stdout = &b
	session.Stderr = &b

	// Enable system stdout
	// Uncomment to fwd to system stdout (comment then store output in variable)
	//session.Stdout = os.Stdout
	//session.Stderr = os.Stderr

	// Start remote shell
	err = session.Shell()
	if err != nil {
		glog.Errorf("Iperf server not restarted. Could not create shell at %s %s", trafficSrvIP, err)
		return err
	}

	// Send commands
	commands := []string{
		"pkill iperf3 ",
		fmt.Sprintf("nohup iperf3 -s --json -B %s > /dev/null &", trafficSrvIP),
		"exit",
	}
	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			glog.Errorf("Iperf server not restarted. Command  %s failed  not create shell at %s %s", cmd, trafficSrvIP, err)
			return err
		}
	}

	// DO NOT WAIT OR SSH WILL HANG UP
	//session.Wait()
	time.Sleep(500 * time.Millisecond)
	glog.V(5).Info("Restarted iperf server")
	return nil
}

func connectToHost(user, pass, host, port string) (*ssh.Client, *ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	hostAndPort := fmt.Sprintf("%s:%s", trafficSrvIP, trafficSrvSSHport)
	client, err := ssh.Dial("tcp", hostAndPort, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}
	return client, session, nil
}
