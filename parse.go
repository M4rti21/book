package main

import (
	"bufio"
	"os"
	"strings"
)

const FOLDER_DELIMITER = "*"
const BOOKMARK_DELIMITER = "\""
const INDENT_LEVEL = 4

func parseFile(config *Config, folder *Folder) error {
	override_config := Config{}
	file, err := os.Open(fixString(config.Root + config.Bookmarks + ".conf"))

	defer file.Close()

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := cleanLine(scanner.Text())
		if strings.Contains(line, "[config]") {
			parseConfig(*scanner, config)
		} else if strings.Contains(line, "[bookmarks]") {
			parseBookmarks(*scanner, folder)
		}
	}

	if override_config.Run != "" {
		config.Run = override_config.Run
	}
	if override_config.FolderIcon != "" {
		config.FolderIcon = override_config.FolderIcon
	}
	if override_config.Menu != "" {
		config.Menu = override_config.Menu
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func cleanLine(line string) string {
	return strings.Trim(line, " ")
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

func parseConfig(scanner bufio.Scanner, config *Config) {
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
			parseVariable(config, split[0], split[1])
		}
	}
}

func parseVariable(config *Config, key string, val string) {
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
	case "showAll":
		config.ShowAll = val == "true"
		break
	}
}

func parseBookmarks(scanner bufio.Scanner, folder *Folder) {
	current_folder := folder

	for scanner.Scan() {
		raw_line := scanner.Text()
		indent := indentLevel(raw_line)
		line := cleanLine(raw_line)

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, FOLDER_DELIMITER) {
			split := strings.Split(line, FOLDER_DELIMITER)
			folder_name := strings.TrimSpace(split[1])
			new_folder := createFolder(folder_name)
			new_folder.Indent = indent

			if indent > current_folder.Indent {
				new_folder.Parent = current_folder
				current_folder.Children[folder_name] = new_folder
			} else if indent == current_folder.Indent {
				new_folder.Parent = current_folder.Parent
				if current_folder.Parent != nil {
					current_folder.Parent.Children[folder_name] = new_folder
				}
			} else if indent < current_folder.Indent {
				for current_folder.Parent != nil && indent < current_folder.Indent {
					current_folder = current_folder.Parent
				}
				new_folder.Parent = current_folder.Parent
				if current_folder.Parent != nil {
					current_folder.Parent.Children[folder_name] = new_folder
				}
			}

			current_folder = new_folder

		} else if strings.HasPrefix(line, BOOKMARK_DELIMITER) {
			split := strings.Split(line, BOOKMARK_DELIMITER)
			switch len(split) {
			case 3:
				current_folder.Index = append(current_folder.Index, Bookmark{
					Name: "",
					Url:  strings.TrimSpace(split[1]),
				})
			case 5:
				current_folder.Index = append(current_folder.Index, Bookmark{
					Name: strings.TrimSpace(split[1]),
					Url:  strings.TrimSpace(split[3]),
				})
			}
		}
	}
}

func createFolder(name string) *Folder {
	new_folder := Folder{
		Name:     name,
		Children: make(map[string]*Folder),
	}
	return &new_folder
}
