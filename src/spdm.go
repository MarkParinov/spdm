package main

import (
	"os"
)

func main() {

	// Get the arguments from the terminal
	args := os.Args
	args = args[1:]

	var getPackageWithGit bool

	// Parse the arguments for flags
	temp_args := args

	for _, v := range args {
		if string(v[0]) == "-" || string(v[1]) == "-" {
			switch v {
			case "--git":
				getPackageWithGit = true
			default:
				log_wrn("flag '", v, "' is not recognized, argument will be ignored")
			}
		} else {
			temp_args = append(temp_args, v)
		}
	}

	args = temp_args

	// If no commands were provided print avaible options for spdm
	if len(args) == 0 {
		spdmHelp()
		return
	}

	setupSPDMEnv()

	// Execute the command given as the first argument after calling the file
	switch args[0] {
	case "get":
		if len(args) == 2 {
			log_err("no package name was provided, unable to get")
			return
		}
		getPackage(args[1], getPackageWithGit)
	case "help":
		spdmHelp()
	default:
		log_err("command '", args[0], "' is not valid. Type 'spdm help' for avaible options")
	}
}
