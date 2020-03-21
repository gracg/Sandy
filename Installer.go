package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	User, e := user.Current()
	if e != nil {
		panic(e)
	}

	if User.Uid != "0" {
		panic("Installer must be run as root")
	}

	fmt.Println("User is root")
	fmt.Println("Writing Systemd service file")
	ioutil.WriteFile("/etc/systemd/system/Sandy.service", []byte(systemdServiceFile()), 0644)

	DirExists, err := exists("/usr/local/bin/Sandy")
	if err != nil {
		panic(err)
	}

	if DirExists == false {
		fmt.Println("Creating software directory")
		err := os.Mkdir("/usr/local/bin/Sandy", 0755)
		if err != nil {
			panic(err)
		}
		FileSlice := []string{"config.json", "Sandy"}
		BasePath := BinaryPath()

		fmt.Println("Copying software to directory")
		for _, file := range FileSlice {
			copy(BasePath+"/"+file, "/usr/local/bin/Sandy/"+file)
		}

		bashcommand := "systemctl"
		args := []string{"enable", "Sandy"}
		// Create arbitrary command.
		terminalCommand(bashcommand, args)

		bashcommand = "systemctl"
		args = []string{"start", "Sandy"}
		// Create arbitrary command.
		terminalCommand(bashcommand, args)

		bashcommand = "systemctl"
		args = []string{"daemon-reload"}
		// Create arbitrary command.
		terminalCommand(bashcommand, args)

	}
}

func systemdServiceFile() string {
	var ServiceFile string

	User, e := user.Current()
	if e != nil {
		panic(e)
	}

	ServiceFile = `[Unit]
Description=Sync public SSH keys to Authorized file

[Service]
Type=simple
User=`
	ServiceFile += User.Username + "\n"
	ServiceFile += `WorkingDirectory=/usr/local/bin/Sandy
ExecStart=/usr/local/bin/Sandy/Sandy -client
Restart=always
StandardOutput=journal


[Install]
WantedBy=multi-user.target`

	return ServiceFile
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func copy(src string, dst string) {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}
	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0755)
	if err != nil {
		panic(err)
	}
}

func BinaryPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func terminalCommand(command string, args []string) error {

	// Create arbitrary command.
	c := exec.Command(command, args...)

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	return nil
}
