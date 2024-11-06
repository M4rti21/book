# 📚 book

Simple bookmark manager written in go

https://github.com/user-attachments/assets/150bcd63-2e35-4703-8e04-ea6dda61cbe9

## Index
- [Instalation](#instalation)
  * [Manual](#manual)
  * [Build](#build)
- [Configuration](#configuration)
  * [Bookmarks](#bookmarks)
    + [Without description](#without-description)
    + [With description](#with-description)
  * [Config](#config)
- [Usage](#usage)

## Instalation
### Manual
Grab the binary from the [latest release](https://github.com/M4rti21/book/releases/latest) 
and put it somewhere detected by your `$PATH`

### Build
You can also build the project yourself:
```sh
git clone https://github.com/m4rti21/book.git
cd book/src
go build -o ..
cd ..
```

## Configuration
The config directory its gonna be at `$XDG_CONFIG_HOME/book`, if the variable
is not set then it'll search at `~/.config/book`.

### Bookmarks
Inside the config directory create a folder for your bookmarks, for example 
`Personal`, inside of that folder you have fill an `index` file with the urls
you want as bookmarks, you can create subdirectories at any level and put an
`index` file on each.

An example structure could be like this:

```
~/.config/book/
--------------------
├── config.toml
├── Personal
|  index
│  ├── Shopping
│  │  └── index
│  ├── index
│  └── Entertainment
│     └── index
└── Work
      └── index
```

The contents of the index file can be in one of two ways:

#### Without description
```
index
------
https://archlinux.org
https://github.com/m4rti21
```

#### With description
```
index
------
Arch Linux Website          # https://archlinux.org
Personal GitHub profile     # https://github.com/m4rti21
```
It is important to note that in this case the urls and the names must be 
separated by a `#` with the names at the left and the urls at the right

### Config
The program will look for a file called `config.toml`.
The default values are:

```toml
folderIcon = ""
menu = "dmenu"
```
| name       | type   | description  |
|------------|--------|--------------|
| folderIcon | string | icon used for folders |
| menu       | string |command for displaying the options (has to accept stdin for the entires)|

## Usage
Once configured its as simple as running `book <dirname>` where `<dirname>` is
one of the folders inside the config directory, for example `book Personal`.
