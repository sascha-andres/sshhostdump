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

func main() {
	printJSON := false
	flaggy.Bool(&printJSON, "j", "json", "print host hierarchy as json")
	flaggy.Parse()

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

	fmt.Println("ssh-menu")
	fmt.Printf("version %s\n\n", version)
}
