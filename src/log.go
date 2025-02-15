// Three basic log functions and a special log_dev function which
// was designed for easier development, feel free to borrow it, i'm
// really proud of that one for some reason :D

package main

import "fmt"

func log_err(args ...string) {
	fmt.Print(RED, "spdm", RESET, ": ")
	for _, v := range args {
		fmt.Print(v)
	}
	fmt.Println()
}

func log_inf(args ...string) {
	fmt.Print(BLUE, "spdm", RESET, ": ")
	for _, v := range args {
		fmt.Print(v)
	}
	fmt.Println()
}

func log_wrn(args ...string) {
	fmt.Print(YELLOW, "spdm", RESET, ": ")
	for _, v := range args {
		fmt.Print(v)
	}
	fmt.Println()
}

func log_dev(description string, args ...string) {
	fmt.Print(MAGENTA, "spdm_dev", RESET, "[", GREEN, description, RESET, "]", ": ")
	for _, v := range args {
		fmt.Print(v)
	}
	fmt.Println()
}
