package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/integrii/flaggy"
	"github.com/sirupsen/logrus"
	"livingit.de/code/sshmenu"
	"os"
)

const version = "develop"

var (
	sshMenuData sshmenu.SSHMenu
)

func main() {
	printJSON := false
	printLines := true
	flat := true
	flaggy.Bool(&printJSON, "j", "json", "print host hierarchy as json")
	flaggy.Bool(&printLines, "l", "lines", "print host hierarchy")
	flaggy.Bool(&flat, "f", "flat", "print hosts with groups one at a line")
	flaggy.Parse()

	if printLines && printJSON {
		logrus.Error("lines and json are mutually exclusive")
		os.Exit(1)
	}

	stat, _ := os.Stdin.Stat()
	if !((stat.Mode() & os.ModeCharDevice) == 0) {
		_, _ = fmt.Fprintln(os.Stderr, "expecting list of config files to be piped")
		os.Exit(1)
	}

	sshMenuData, err := sshmenu.NewSSHMenu("./")
	if err != nil {
		logrus.Errorf("error creating ssh menu data handler: %s", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		configFile := scanner.Text()
		file, err := os.OpenFile(configFile, os.O_RDONLY, 0400)
		if err != nil {
			logrus.Errorf("error reading config file %s:%s", configFile, err)
			os.Exit(1)
		}
		err = sshMenuData.FromReader(file)
		_ = file.Close()
		if err != nil {
			logrus.Errorf("error getting host data: %s", err)
			os.Exit(1)
		}
	}

	if printJSON {
		data, err := json.MarshalIndent(sshMenuData, "", " ")
		if err != nil {
			os.Exit(1)
		}
		fmt.Println(string(data))
		os.Exit(0)
	}

	if printLines {
		printDirectory(sshMenuData.Data)
	}

	if false {
		fmt.Println("ssh-menu")
		fmt.Printf("version %s\n\n", version)
	}

	//err = syscall.Exec("/usr/local/bin/ssh", []string{"halcon"}, nil)
	//logrus.Errorf("error: syscall failed: %s", err)
}

func printDirectory(directory sshmenu.DirectoryEntry) {
	for _, val := range directory.Hosts {
		fmt.Println(val)
	}

	for _, val := range directory.Directories {
		printDirectory(val)
	}
}
