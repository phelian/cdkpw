package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var execCommand = exec.Command

type CDKCommand struct {
	Action    string   // diff, deploy, etc.
	StackName string   // the first non-flag positional arg
	Profile   string   // value from --profile if present
	RawArgs   []string // full CLI args
	Context   []string // all `-c` and `--context` switches
	Flags     []string // any other flags (e.g. --exclusively)
}

func (c *CDKCommand) SetProfile(profile string) {
	if c.Profile != "" {
		return
	}
	c.Profile = profile
	c.RawArgs = append(c.RawArgs, "--profile", profile)
}

func (c *CDKCommand) Execute(cdk string) {
	cmd := execCommand(os.ExpandEnv(cdk), c.RawArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running cdk command:", err)
		os.Exit(cmd.ProcessState.ExitCode())
	}
}

func (c *CDKCommand) IsProfiled() bool {
	return c.Profile != ""
}

func parseArgs(args []string) *CDKCommand {
	cmd := CDKCommand{
		RawArgs: args,
	}

	if len(args) == 0 {
		return &cmd
	}

	cmd.Action = args[0]

	for i := 1; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "--profile" && i+1 < len(args):
			cmd.Profile = args[i+1]
			i++
		case arg == "-c" || strings.HasPrefix(arg, "-c") || strings.HasPrefix(arg, "--context"):
			cmd.Context = append(cmd.Context, arg)
			if (arg == "-c" || arg == "--context") && i+1 < len(args) {
				cmd.Context = append(cmd.Context, args[i+1])
				i++
			}
		case strings.HasPrefix(arg, "-"):
			cmd.Flags = append(cmd.Flags, arg)
		case cmd.StackName == "":
			cmd.StackName = arg
		}
	}

	return &cmd
}
