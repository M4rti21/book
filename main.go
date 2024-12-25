package main

import (
	"errors"
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
		ShowUrl:    false,
		ShowAll:    false,
	}

	var folder Folder = Folder{
		Name:     "index",
		Indent:   0,
		Children: make(map[string]*Folder),
	}

	loadConfig(&config)
	loadArguments(&config)

	err := parseFile(&config, &folder)

	if err != nil {
		panic(err)
	}

	run(&folder, &config)
}

func run(folder *Folder, config *Config) {
	entries := generateEntries([]Entry{}, folder, "", *config)
	selected, err := showPrompt(entries, config.Menu)
	if err != nil {
		panic(err)
	}

	if selected.isFolder {
		run(selected.folder, config)
	} else {
		run_split := strings.Fields(config.Run)
		run_split = append(run_split, selected.bookmark.Url)
		cmd := exec.Command(run_split[0], run_split[1:]...)
		cmd.Start()
	}
}

func indexOf(arr []Entry, val string) int {
	for pos, v := range arr {
		if v.Name == val {
			return pos
		}
	}
	return -1
}

func generateEntries(entries []Entry, folder *Folder, basename string, config Config) []Entry {

	createBookmarkEntry := func(b Bookmark) Entry {
		var name string
		if b.Name == "" {
			name = b.Url
		} else if config.ShowUrl {
			name = b.Name + " -> " + b.Url
		} else {
			name = b.Name
		}
		return Entry{
			Name:     basename + name,
			isFolder: false,
			bookmark: &b,
		}
	}

	if !config.ShowAll && folder.Parent != nil {
		entries = append(entries, Entry{
			Name:     config.FolderIcon + " ..",
			isFolder: true,
			folder:   folder.Parent,
		})
	}

	if !config.ShowAll {
		for k, v := range folder.Children {
			entries = append(entries, Entry{
				Name:     config.FolderIcon + " " + k,
				isFolder: true,
				folder:   v,
			})
		}
	}

	for _, b := range folder.Index {
		entries = append(entries, createBookmarkEntry(b))
	}

	if config.ShowAll {
		for _, f := range folder.Children {
			entries = generateEntries(entries, f, basename+f.Name+"/", config)
		}
	}

	return entries
}

func showPrompt(entries []Entry, menuCommand string) (Entry, error) {
	var entriesTxt string

	for _, e := range entries {
		entriesTxt += e.Name + "\n"
	}

	menu := strings.Fields(menuCommand)
	cmd := exec.Command(menu[0], menu[1:]...)
	cmd.Stdin = strings.NewReader(entriesTxt)
	output, err := cmd.Output()
	if err != nil {
		return Entry{}, errors.New("something went wrong")
	}

	res := strings.TrimSuffix(string(output), "\n")
	index := indexOf(entries, res)
	if index == -1 {
		return Entry{}, errors.New("no entry selected")
	}

	return entries[index], nil
}
