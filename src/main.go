package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var root string
var startLocation string

type Config struct {
	FolderIcon string
	Menu       string
	Run        string
}

var config = Config{
	FolderIcon: "",
	Menu:       "dmenu",
	Run:        "xdg-open",
}

type Arguments struct {
	Bookmarks  string
	Config     string
	Directory  string
	FolderIcon string
	Menu       string
	Run        string
}

var arguments = Arguments{}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("At least 1 argument is required")
	}
	loadVariables()
}

func loadArguments() {

	flag.StringVar(&arguments.Config, "c", "", "Specify config file")
	flag.StringVar(&arguments.Directory, "d", "", "Specify bookmark directory")
	flag.StringVar(&arguments.FolderIcon, "f", "", "Specify folder icon")
	flag.StringVar(&arguments.Menu, "m", "", "Specify menu command")
	flag.StringVar(&arguments.Run, "r", "", "Specify run command")

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Error: Missing required positional argument.")
		os.Exit(1)
	}

	arguments.Bookmarks = flag.Args()[0]
}

func loadConfig() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	if arguments.Directory != "" {
		root = strings.Replace(arguments.Directory, "~", homeDir, 1)
	} else {
		configDir := os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			root = strings.Replace("~/.config/book/", "~", homeDir, 1)
		} else {
			root = configDir + "/book/"
		}
		startLocation = root + arguments.Bookmarks
	}

	var configFile string

	if arguments.Config != "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		configFile = strings.Replace(arguments.Config, "~", homeDir, 1)
	} else {
		configFile = root + "config.toml"
	}

	file, err := os.ReadFile(configFile)
	if err != nil {
	}
	err = toml.Unmarshal(file, &config)
	if err != nil {
	}
	if arguments.Run != "" {
		config.Run = arguments.Run
	}
	if arguments.FolderIcon != "" {
		config.FolderIcon = arguments.FolderIcon
	}
	if arguments.Menu != "" {
		config.Menu = arguments.Menu
	}
}

func dir(dirname string) {
	location := root + dirname
	var directories []string

	entries, err := os.ReadDir(location)
	if err != nil {
		panic(err)
	}

	// Loop through the entries
	for _, entry := range entries {
		// Check if it's a directory
		if entry.IsDir() {
			directories = append(directories, location+"/"+entry.Name()) // Append full path
		}
	}

	links, err := readLines(location + "/index")
	if err != nil {
		log.Fatal(err)
	}

	var list []string
	var dirs []string
	var pipe string

	if filepath.Clean(location) != startLocation {
		directories = append([]string{".."}, directories...)
	}

	if len(directories) > 0 {
		for _, d := range directories {
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
		run := strings.Fields(config.Run)
		if len(link_split) == 1 {
			run = append(run, selected)
		} else {
			run = append(run, strings.Trim(link_split[1], " "))
		}
		cmd := exec.Command(run[0], run[1:]...)
		err := cmd.Start()
		if err != nil {
			return
		}
		return
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
