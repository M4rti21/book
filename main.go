package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/k0kubun/pp"
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
	// prev_indent := 0
	// prev_folder := urls
	current_folder := &urls
	current_folder.folders = make(map[string]*Folder)
	for scanner.Scan() {
		raw_line := scanner.Text()
		// indent := indentLevel(raw_line)
		line := cleanLine(raw_line)
		if strings.HasPrefix(line, FOLDER_DELIMITER) {
			split := strings.Split(line, FOLDER_DELIMITER)
			folder_name := strings.Trim(split[1], " ")
			current_folder := &createFolder()
		} else if strings.HasPrefix(line, BOOKMARK_DELIMITER) {
			split := strings.Split(line, BOOKMARK_DELIMITER)
			switch len(split) {
			case 3:
				current_folder.index = append(current_folder.index, Bookmark{
					Name: "",
					URL:  split[1],
				})
				break
			case 5:
				current_folder.index = append(current_folder.index, Bookmark{
					Name: split[1],
					URL:  split[3],
				})
				break
			}
		}
	}
	pp.Print(folders)
	// pp.Print(current_folder)
}

func createFolder() Folder {
	index := len(folders)
	folders = append(folders, Folder{})
	folders[index].folders = make(map[string]*Folder)
	return folders[index]
}
