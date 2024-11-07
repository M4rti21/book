package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var root string
var startLocation string

func main() {

	configDir := os.Getenv("XDG_CONFIG_HOME")
	var root string
	if configDir == "" {
		root = fixString("~/.config/book/")
	} else {
		root = configDir + "/book/"
	}

	var config Config = Config{
		Root:       root,
		ConfigFile: root + "config.toml",
		Run:        "xdg-open",
		Menu:       "dmenu",
		FolderIcon: "",
		ShowUrl:    true,
	}

	var folder Folder = Folder{
		Name:     "index",
		Indent:   0,
		Children: make(map[string]*Folder),
	}

	loadArguments(&config)
	loadConfig(&config)

	err := parseFile(&config, &folder)

	if err != nil {
		panic(err)
	}

	run(&folder, &config)

}

func run(folder *Folder, config *Config) {
	entries := stringifyFolder(*folder, *config)
	selected := showPrompt(entries, config.Menu)
	fmt.Println(selected)
	index := indexOf(entries, selected)
	if index == -1 {
		return
	}
	if folder.Parent != nil {
		if ".." == strings.TrimPrefix(selected, config.FolderIcon+" ") {
			run(folder.Parent, config)
			return
		} else {
			index--
		}
	}
	if index < len(folder.Children) {
		selected = strings.TrimPrefix(selected, config.FolderIcon+" ")
		run(folder.Children[selected], config)
	} else if index-len(folder.Children) < len(folder.Index) {
		bookmark := folder.Index[index-len(folder.Children)]
		fmt.Println(index, bookmark.Name, bookmark.Url)
		run_split := strings.Fields(config.Run)
		run_split = append(run_split, bookmark.Url)
		cmd := exec.Command(run_split[0], run_split[1:]...)
		cmd.Start()
	}
}

func indexOf(arr []string, val string) int {
	for pos, v := range arr {
		if v == val {
			return pos
		}
	}
	return -1
}

func stringifyFolder(folder Folder, config Config) []string {
	var entries []string

	if folder.Parent != nil {
		entries = append(entries, config.FolderIcon+" ..")
	}

	for k := range folder.Children {
		entries = append(entries, config.FolderIcon+" "+k)
	}

	for _, b := range folder.Index {
		var line string
		if !config.ShowUrl || b.Name == "" {
			line = b.Url
		} else {
			line = b.Name + " -> " + b.Url
		}
		entries = append(entries, line)
	}

	return entries
}

func showPrompt(entries []string, menuCommand string) string {
	var entriesTxt string

	for _, e := range entries {
		entriesTxt += e + "\n"
	}

	menu := strings.Fields(menuCommand)
	cmd := exec.Command(menu[0], menu[1:]...)
	cmd.Stdin = strings.NewReader(entriesTxt)
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSuffix(string(output), "\n")
}
