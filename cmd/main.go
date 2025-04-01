package main

import (
	"fmt"
	"os"
)

func main() {
	cdkCommand := parseArgs(os.Args[1:])

	if !cdkCommand.IsProfiled() {
		switch cdkCommand.Action {
		case "diff", "deploy", "destroy", "bootstrap":
			config, err := loadConfig()
			if err != nil {
				fmt.Println("Error loading config:", err)
				os.Exit(1)
			}

			if profile, found := findProfile(cdkCommand.StackName, config); found {
				cdkCommand.SetProfile(profile)
			}
		default:
			//  do nothing
		}
	}

	cdkCommand.Execute()
}
