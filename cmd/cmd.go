package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// ExecuteP execute and show progress
func ExecuteP(command string, args ...string) (string, error) {
	if runtime.GOOS == "windows" {
		args = prepend(args, command)
		command = "powershell"
	}

	cmd := exec.Command(command, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return "", errors.New(string(stderr.Bytes()))
	}

	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				fmt.Print(".")
				time.Sleep(3 * time.Second)
			}
		}
	}()

	err = cmd.Wait()
	quit <- true

	if err != nil {
		return "", errors.New(string(stderr.Bytes()))
	}
	fmt.Println()

	return string(stdout.Bytes()), nil
}

// Execute execute commands and return output as string or error
func Execute(command string, args ...string) (string, error) {
	if runtime.GOOS == "windows" {
		args = prepend(args, command)
		command = "powershell"
	}

	cmd := exec.Command(command, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return "", errors.New(string(stderr.Bytes()))
	}

	return string(stdout.Bytes()), nil
}

// Execute execute commands in interactive mode
func ExecuteI(command string, args ...string) error {
	if runtime.GOOS == "windows" {
		args = prepend(args, command)
		command = "powershell"
	}

	cmd := exec.Command(command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

// Run runs command and return nothing
func Run(command string, args ...string) error {
	if runtime.GOOS == "windows" {
		args = prepend(args, command)
		command = "powershell"
	}
	return exec.Command(command, args...).Run()
}

func prepend(x []string, y string) []string {
	x = append(x, "")
	copy(x[1:], x)
	x[0] = y
	return x
}
