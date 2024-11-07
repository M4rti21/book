package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	// "github.com/k0kubun/pp"
)

type Config struct {
	FolderIcon string
	Menu       string
	Run        string
	ShowUrl    bool
}

type Bookmark struct {
	Name string
	URL  string
}

type Folder struct {
	parent  *Folder
	indent  int
	index   []Bookmark
	folders map[string]*Folder
}

type Data struct {
	Config    Config
	Bookmarks Folder
}

var config = Config{
	FolderIcon: "",
	Menu:       "dmenu",
	Run:        "xdg-open",
	ShowUrl:    false,
}

var folders = []Folder{}

var urls = Folder{}

const FOLDER_DELIMITER = "*"
const BOOKMARK_DELIMITER = "\""
const INDENT_LEVEL = 4

func main() {
	file, err := os.Open("test.conf")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := cleanLine(scanner.Text())
		if strings.Contains(line, "[config]") {
			parseConfig(*scanner)
		} else if strings.Contains(line, "[bookmarks]") {
			parseBookmarks(*scanner)
		}
		// fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	// fmt.Println(config)
}

func cleanLine(line string) string {
	return strings.Trim(strings.Split(line, "#")[0], " ")
}

func indentLevel(s string) int {
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			count++
		} else {
			break
		}
	}
	return count
}

func parseConfig(scanner bufio.Scanner) {
	prev_indent := 0
	for scanner.Scan() {
		raw_line := scanner.Text()
		indent := indentLevel(raw_line)
		line := cleanLine(raw_line)
		if indent < prev_indent {
			return
		}
		if line == "" {
			return
		}
		split := strings.Split(line, "=")
		if len(split) > 1 {
			parseVariable(split[0], split[1])
		}
	}
}

func parseVariable(key string, val string) {
	key = strings.Trim(key, " ")
	val = strings.Trim(val, " ")
	switch key {
	case "folderIcon":
		config.FolderIcon = strings.Trim(val, "\"")
		break
	case "menu":
		config.Menu = strings.Trim(val, "\"")
		break
	case "run":
		config.Run = strings.Trim(val, "\"")
		break
	case "showUrl":
		config.ShowUrl = val == "true"
		break
	}
}

func parseBookmarks(scanner bufio.Scanner) {
	current_folder := &urls
	current_folder.indent = 0 // Initialize the indent level for the root folder
	current_folder.folders = make(map[string]*Folder)

	for scanner.Scan() {
		raw_line := scanner.Text()
		indent := indentLevel(raw_line)
		line := cleanLine(raw_line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Debug output for current line and indent level
		fmt.Printf("%d %s\n", indent, line)

		if strings.HasPrefix(line, FOLDER_DELIMITER) {
			// Handle folder creation
			split := strings.Split(line, FOLDER_DELIMITER)
			folder_name := strings.TrimSpace(split[1])
			new_folder := createFolder()
			new_folder.parent = current_folder // Set the parent of the new folder
			new_folder.indent = indent         // Set the indent level for the new folder

			if indent > current_folder.indent {
				// If the new folder is indented more, it is a child of the current folder
				fmt.Printf("Creating child folder: '%s' under '%s'\n", folder_name, current_folder.index)
				current_folder.folders[folder_name] = new_folder
				current_folder = new_folder // Move into the new folder
			} else if indent == current_folder.indent {
				// If at the same level, add it as a sibling in the parent folder
				if current_folder.parent != nil {
					fmt.Printf("Creating sibling folder: '%s' under parent '%s'\n", folder_name, current_folder.parent.index)
					current_folder.parent.folders[folder_name] = new_folder
				}
			} else if indent < current_folder.indent {
				// If indent level decreases, go back to the parent folder
				fmt.Printf("Going back to parent folder from '%s'\n", current_folder.index)
				for current_folder.parent != nil && indent < current_folder.indent {
					current_folder = current_folder.parent
				}
				// Now add the new folder as a sibling if we're back at the parent level
				if current_folder.parent != nil {
					fmt.Printf("Creating sibling folder: '%s' under parent '%s'\n", folder_name, current_folder.parent.index)
					current_folder.parent.folders[folder_name] = new_folder
				}
				current_folder = new_folder // Move into the new folder
			}

		} else {
			// Handle bookmark parsing
			if strings.HasPrefix(line, BOOKMARK_DELIMITER) {
				split := strings.Split(line, BOOKMARK_DELIMITER)
				switch len(split) {
				case 3:
					fmt.Printf("Adding bookmark with no name: '%s'\n", split[1])
					current_folder.index = append(current_folder.index, Bookmark{
						Name: "", // No name provided
						URL:  strings.TrimSpace(split[1]),
					})
				case 5:
					fmt.Printf("Adding bookmark: Name: '%s', URL: '%s'\n", split[1], split[3])
					current_folder.index = append(current_folder.index, Bookmark{
						Name: strings.TrimSpace(split[1]),
						URL:  strings.TrimSpace(split[3]),
					})
				}
			}
		}
	}
	// pp.Print(urls) // Print the populated URLs and folder structure
}

func createFolder() *Folder {
	index := len(folders)
	folders = append(folders, Folder{})
	folders[index].folders = make(map[string]*Folder)
	return &folders[index]
}
