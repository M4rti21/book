package main

type Config struct {
	Root       string
	ConfigFile string
	Bookmarks  string
	Menu       string
	Run        string
	FolderIcon string
	ShowUrl    bool
	ShowAll    bool
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

type Entry struct {
	Name     string
	bookmark *Bookmark
	folder   *Folder
	isFolder bool
}
