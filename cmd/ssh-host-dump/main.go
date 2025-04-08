package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sascha-andres/reuse/flag"

	"livingit.de/code/sshhostdump"
)

const version = "develop"

var (
	printJSON   bool
	printLines  bool
	flat        bool
	showVersion bool
)

func init() {
	flag.SetEnvPrefix("SSH_HOST_DUMP")
	flag.BoolVar(&printJSON, "json", false, "print host hierarchy as json")
	flag.BoolVar(&printLines, "lines", true, "print host hierarchy")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&flat, "flat", false, "print hosts with groups one at a line")
}

func main() {
	log.SetPrefix("[ssh-host-dump] ")
	log.SetFlags(log.LstdFlags | log.LUTC)

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
		log.Fatal("lines and json are mutually exclusive")
	}

	stat, _ := os.Stdin.Stat()
	if !((stat.Mode() & os.ModeCharDevice) == 0) {
		log.Fatal("expecting list of config files to be piped")
	}

	sshMenuData, err := sshhostdump.NewSSHMenu("./")
	if err != nil {
		log.Fatalf("error creating ssh menu data handler: %s", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		configFile := scanner.Text()
		file, err := os.OpenFile(configFile, os.O_RDONLY, 0400)
		if err != nil {
			log.Fatalf("error reading config file %s:%s", configFile, err)
		}
		err = sshMenuData.FromReader(file)
		_ = file.Close()
		if err != nil {
			log.Fatalf("error getting host data: %s", err)
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
