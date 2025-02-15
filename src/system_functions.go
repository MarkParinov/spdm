// Basically a kernel of the whole app, stores functions for getting packages,
// generating config file and directories, command executing

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

func check(e error) {
	if e != nil {
		log_err(e.Error())
	}
}

func executeShellCommand(name string, args ...string) (err bool) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		log_err("could not execute shell command, ", err.Error())
		return true
	}
	return false
}

func parseGitName(url string) string {
	if len(url) < 14 {
		// log_dev("url type", "url type = invalid")
		return url
	} else {
		if url[:8] == "https://" {
			// log_dev("url type", "url type = https")
			return url[19:]
		} else if url[:15] == "git@github.com:" {
			// log_dev("url type", "url type = git@")
			return url[15:]
		} else {
			// log_dev("url type", "undefined url type")
			return url
		}
	}
}

func cloneGitRepoToPackages(git_url string, destination_path string) {
	// Clone the git repo to packages directory, the long command below basically does that
	if executeShellCommand("/bin/bash", "-c", fmt.Sprintf("git clone %v %v", git_url, fmt.Sprintf("%v/%v", destination_path, parseGitName(git_url)))) {
		log_err("error while getting package, make sure that the package is valid and check for spelling mistakes in the url")
	} else {
		log_inf("done getting '", git_url, "'")
	}
}

func getUser() *user.User {
	USERNAME, err := user.Current()
	check(err)
	return USERNAME
}

func spdmHelp() {
	fmt.Print("\nspdm - Simple Package Distribution Manager\n\n")
	// ABOUT PAGE
	fmt.Print("ABOUT\n\n")
	fmt.Println("  A command line application, mainly inspired by npm. It's main goal is to")
	fmt.Println("  quickly install packages on the machine it's used on, whether it's a")
	fmt.Println("  GitHub repository or a spdm package. Here, any set of files can be")
	fmt.Println("  referenced as a package, spdm can be used to quickly install any")
	fmt.Println("  publicly shared software, library or project.")
	fmt.Println()
	// COMMANDS PAGE
	fmt.Print("COMMANDS\n\n")
	fmt.Println("  get  - installs the entered package into the 'packages' subdirectory,")
	fmt.Println(" 		  located in 'spdm' root directory. Related flags: --git")
	fmt.Println("  help - displays information and usage of the application. No related flags")
	fmt.Println()
	// FLAGS PAGE
	fmt.Print("FLAGS\n\n")
	// still under construction
}

// Generates a config for spdm. Default path: /home/{user}/spdm/packages/
func generateSPDMConfig() {
	// CONFIG STRUCTURE

	// LINE  EXAMPLE VALUE				COMMENT
	// 1. 	 /home/user/spdm			user root directory
	// 2.	 /home/user/spdm/packages	packages directory

	// Write following data to a config file

	path := fmt.Sprintf("/home/%v/.spdmconfig", getUser().Name)
	data_to_write := fmt.Sprintf("/home/%v/spdm\n", getUser().Name)
	data_to_write += fmt.Sprintf("/home/%v/spdm/packages\n", getUser().Name)

	err := os.WriteFile(path, []byte(data_to_write), 0644)
	check(err)
}

func getConfigData() []string {

	// Parses the config file and returns the decoded result
	out := []string{}
	file, err := os.Open(fmt.Sprintf("/home/%v/.spdmconfig", getUser().Name))
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		out = append(out, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log_err(err.Error())
	}

	return out
}

func setupSPDMEnv() {

	USERNAME := getUser()

	if _, err := os.Stat(fmt.Sprintf("/home/%v/.spdmconfig", USERNAME.Name)); errors.Is(err, os.ErrNotExist) {
		log_inf("~/.spdmconfig wasn't located, generating default config file")
		generateSPDMConfig()
	}

	necessaryDirectories := []string{"", "packages"}
	configData := getConfigData()

	default_dir := configData[0]

	// First, check if config file is there, if not - generate one

	// Create all necessary directories
	for _, v := range necessaryDirectories {
		if _, err := os.Stat(fmt.Sprintf("%v/%v", default_dir, v)); os.IsNotExist(err) {
			log_inf("spdm default directory not found, creating ~/spdm/", v)
			err := os.Mkdir(fmt.Sprintf("%v/%v", default_dir, v), 0755)
			check(err)
		}
	}
}

func getPackage(pckgName string, git bool) {
	// USERNAME := getUser()
	if git {
		log_inf("getting package ", pckgName, " with git")
	} else {
		log_inf("getting package ", pckgName, " without git")
	}

	// Check if the package is already installed

	default_dir := getConfigData()[1]

	if _, err := os.Stat(fmt.Sprintf("%v/%v", default_dir, parseGitName(pckgName))); !os.IsNotExist(err) {
		log_err("package already installed at ", fmt.Sprintf("%v/%v", default_dir, parseGitName(pckgName)))
		return
	}

	// cmd := exec.Command("bash", "cd", fmt.Sprintf("/home/%v/spdm/packages", USERNAME.Name))
	if git {
		cloneGitRepoToPackages(pckgName, getConfigData()[1])
	} else {
		log_err("spdm doesn't support installation other than from GitHub yet")
	}
}
