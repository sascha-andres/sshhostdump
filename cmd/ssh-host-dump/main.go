package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sascha-andres/flag"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"livingit.de/code/sshhostdump"
)

const version = "develop"

var (
	sshMenuData sshhostdump.SSHMenu

	printJSON   bool
	printLines  bool
	flat        bool
	showVersion bool
)

func init() {
	flag.BoolVar(&printJSON, "json", false, "print host hierarchy as json")
	flag.BoolVar(&printLines, "lines", true, "print host hierarchy")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&flat, "flat", false, "print hosts with groups one at a line")
}

func main() {
	flag.Parse()
	if !flag.Parsed() {
		log.Fatal("flags have not been parsed")
	}

	if showVersion {
		fmt.Println("ssh-menu")
		fmt.Printf("version %s\n\n", version)
		os.Exit(0)
	}

	if printLines && printJSON {
		logrus.Error("lines and json are mutually exclusive")
		os.Exit(1)
	}

	stat, _ := os.Stdin.Stat()
	if !((stat.Mode() & os.ModeCharDevice) == 0) {
		_, _ = fmt.Fprintln(os.Stderr, "expecting list of config files to be piped")
		os.Exit(1)
	}

	sshMenuData, err := sshhostdump.NewSSHMenu("./")
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

}

// printDirectory prints out all hosts in a directory and
// for each directory calls printDirectory
func printDirectory(directory sshhostdump.DirectoryEntry) {
	for _, val := range directory.Hosts {
		fmt.Println(val)
	}

	for _, val := range directory.Directories {
		printDirectory(val)
	}
}
