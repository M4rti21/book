package main

type Config struct {
	Root       string
	ConfigFile string
	Bookmarks  string
	Menu       string
	Run        string
	FolderIcon string
	ShowUrl    bool
}

type Bookmark struct {
	Name string
	Url  string
}

type Folder struct {
	Parent   *Folder
	Name     string
	Indent   int
	Index    []Bookmark
	Children map[string]*Folder
}
