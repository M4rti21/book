package main

import (
	"flag"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func loadArguments(config *Config) {
	var args_config Config
	flag.StringVar(&args_config.Root, "d", "", "Specify bookmark file {/path/to/directory/} where <bookmark-name>.conf is at")
	flag.StringVar(&args_config.ConfigFile, "c", "", "Specify a base config file")
	flag.StringVar(&args_config.FolderIcon, "f", "", "Specify folder icon")
	flag.StringVar(&args_config.Menu, "m", "", "Specify menu command")
	flag.BoolVar(&args_config.ShowUrl, "u", false, "Show URL in menu")
	flag.StringVar(&args_config.Run, "r", "", "Specify run command")

	flag.Parse()

	if len(flag.Args()) == 0 {
		panic("Error: Missing required positional argument.")
	}

	if args_config.Root != "" {
		config.Root = fixString(config.Root)
	}
	if args_config.ConfigFile != "" {
		config.ConfigFile = fixString(args_config.ConfigFile)
	}
	if args_config.Menu != "" {
		config.Menu = fixString(args_config.Menu)
	}
	if args_config.Run != "" {
		config.Run = fixString(args_config.Run)
	}
	if args_config.FolderIcon != "" {
		config.FolderIcon = fixString(args_config.FolderIcon)
	}
	if args_config.ShowUrl {
		config.ShowUrl = true
	}

	config.Bookmarks = flag.Args()[0]
}

func loadConfig(config *Config) {
	base_config := Config{}

	file, _ := os.ReadFile(config.ConfigFile)
	_ = toml.Unmarshal(file, &base_config)

	if base_config.Root != "" {
		config.Root = fixString(config.Root)
	}
	if base_config.ConfigFile != "" {
		config.ConfigFile = fixString(base_config.ConfigFile)
	}
	if base_config.Menu != "" {
		config.Menu = fixString(base_config.Menu)
	}
	if base_config.Run != "" {
		config.Run = fixString(base_config.Run)
	}
	if base_config.FolderIcon != "" {
		config.FolderIcon = fixString(base_config.FolderIcon)
	}
	if base_config.ShowUrl {
		config.ShowUrl = true
	}
}
