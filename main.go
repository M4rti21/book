package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var root string

type Config struct {
	FolderIcon string
	Menu       string
}

var config = Config{
	FolderIcon: "",
	Menu:       "dmenu",
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("At least 1 argument is required")
	}

	loadVariables()
	loadConfig()
	dir(os.Args[1])
}

func loadVariables() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		root = strings.Replace("~/.config/book/", "~", homeDir, 1)
	} else {
		root = configDir + "/book/"
	}

}

func loadConfig() {
	file, err := os.ReadFile(root + "config.toml")
	if err != nil {
		return
	}
	err = toml.Unmarshal(file, &config)
	if err != nil {
		return
	}
}

func dir(dirname string) {
	location := root + dirname

	var directories []string

	// Walk through the directory
	err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if it's a directory
		if info.IsDir() {
			directories = append(directories, path) // Append to the slice
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	links, err := readLines(location + "/index")
	if err != nil {
		log.Fatal(err)
	}

	var list []string
	var dirs []string

	var pipe string

	if len(directories) > 1 {
		for _, d := range directories[1:] {
			split := strings.Split(d, "/")
			dirs = append(dirs, split[len(split)-1])
			line := config.FolderIcon + " " + split[len(split)-1] + "\n"
			pipe += line
			list = append(list, line)
		}
	}

	for _, l := range links {
		pipe += l + "\n"
		list = append(list, l+"\n")
	}

	out := showPrompt(pipe)
	index := indexOf(list, out)

	if index == -1 {
		return
	} else if index < len(dirs) {
		dir(dirname + "/" + dirs[index])
	} else if index-len(dirs) < len(links) {
		selected := links[index-len(dirs)]
		link_split := strings.Split(selected, "#")
		if len(link_split) == 1 {
			cmd := exec.Command("xdg-open", selected)
			err := cmd.Start()
			if err != nil {
				return
			}
			return
		}
		cmd := exec.Command("xdg-open", strings.Trim(link_split[1], " "))
		err := cmd.Start()
		if err != nil {
			return
		}
	}

}

func showPrompt(pipe string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	menuCommand := strings.Replace(config.Menu, "~", homeDir, 1)
	menu := strings.Fields(menuCommand)
	cmd := exec.Command(menu[0], menu[1:]...)
	cmd.Stdin = strings.NewReader(pipe)
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(output)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func indexOf(arr []string, val string) int {
	for pos, v := range arr {
		if v == val {
			return pos
		}
	}
	return -1
}
