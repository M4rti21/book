# 📚 book

Simple bookmark manager written in go

https://github.com/user-attachments/assets/150bcd63-2e35-4703-8e04-ea6dda61cbe9

## Index
- [Instalation](#instalation)
  * [Manual](#manual)
  * [Build](#build)
- [Configuration](#configuration)
  * [Settings](#settings)
  * [Bookmarks](#bookmarks)
- [Usage](#usage)

## Instalation
### Manual
Grab the binary from the [latest release](https://github.com/M4rti21/book/releases/latest) 
and put it somewhere detected by your `$PATH`

### Build
You can also build the project yourself:
```sh
git clone https://github.com/m4rti21/book.git
cd book
go build -o ..
cd ..
```

## Configuration
The default config directory is `$XDG_CONFIG_HOME/book`, if the variable
is not set then it will fallback to `~/.config/book`.

### Settings
The program will look for a file called `config.toml` in the config directory.
The allowed options are:
| name          | type      | default       | flag | description  |
|---------------|-----------|---------------|------|--------------|
| folderIcon    | string    | `""`           | -f   | icon used for folders |
| menu          | string    | `"dmenu"`       | -m   | program for displaying the options (has to accept stdin for the entires) |
| run           | string    | `"xdg-open"`    | -r   | program for opening the selected url |
| showUrl       | boolean   | `true`          | -u   | weather to show the url if a name is provided |

##### Example: 
```toml
# config.toml
menu = "tofi"
showUrl = false
```

> [!NOTE]
> The next options only exist as flags, they will be ignored if present in the
> config file

| name          | type      | default                               | flag  | description |
|---------------|-----------|---------------------------------------|-------|-------------|
| config        | string    | `"$XDG_CONFIG_HOME"/book/config.conf` | -c    | location of the base config file |
| directory     | string    | `"$XDG_CONFIG_HOME"/book`             | -d    | location of the directory where <collection-name>.conf will be searched for |

### Bookmarks
Inside the config directory create file called `<collection-name>.conf` where your
bookmarks will be stored, for example `personal.conf`. You can have as many 
collections as you want. The file is divided in two sections: `[config]` and `[bookmarks]`.

> [!IMPORTANT]
> Indentation in this file has to be 4 spaces and cannot be ommited

> [!TIP]
> Anything at the ritght side of the hashtag `#` symbol will be interpreted as a 
> comment and therefore it will be ignored

The `[config]` section can be used to override the global settings defined at the 
`config.toml` file (for example, there might be a collection of bookmarks you might 
want to show on a different menu or launch them with a different browser).

##### Example: 
```conf
# personal.conf
[config]
    menu = "tofi"
    run = "xdg-open"
    showUrl = true
```

The `[bookmarks]` section is where all the bookmarks will be defined, the syntax
goes as follows:

After the `[bookmarks]` tag we will define all the bookmarks that go on the base
level, meaning the bookmarks you will first see when running the program, each
bookmark must be on its own line. Said line can have 1 or 2 arguments, if only
1 is provided that must be the URL, if instead you have 2 the first one will be
the Label/Name of the url and the second one will then be the URL. Each argument
must start and end with double quotes `"` and cannot contain a double quote, they
must be separated by at least one space.


##### Example: 
```conf
# personal.conf
[bookmarks]
    "Example"                   "https://example.com"   # This is a comment
    "Another Example"           "https://someotherurl.com/i/dont/know"
    "https://unnamedurl.com"    # This url doesnt have a name
```

You can also define sub-folders, to define one the line must start with the star
`*` symbol, anything following will be the folder name


##### Example: 
```conf
# personal.conf
[config]
    menu = "tofi"
    run = "xdg-open"
    showUrl = true

[bookmarks]
    "Example"                   "https://example.com"   # This is a comment
    "Another Example"           "https://someotherurl.com/i/dont/know"
    "https://unnamedurl.com"    # This url doesnt have a name

    * This is a folder
        "Woah how cool"         "https://how.cool"
        "Woah how cool"         "https://how.cool"
        "Woah how cool"         "https://how.cool"

        * More???           # You can add as many sub-sub folders as you want
            "You get the point"     "https://i.am.runnign/out/of/url/ideas
            #...
            #...
            #...

    * This is another folder
        "I am in your walls"         "https://wake.up/wake/up/wake/up"
        "I am in your walls"         "https://wake.up/wake/up/wake/up"
        "I am in your walls"         "https://wake.up/wake/up/wake/up"
```

## Usage
Once configured its as simple as running `book <collection-name>` where 
`<collection-name>` is the name of one of the collections in your config directory, 
for example `book personal` will open `personal.conf`. You can pass flags to
the program after the `<collection-name>` is specified:

##### Example: 
```sh
book personal
book work -m "tofi" -d "~/.work/book"
```
