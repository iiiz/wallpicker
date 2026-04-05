# WallPicker

Simple wallpaper picker utility.

![wallpicker](./preview.gif)

```
Usage: wallpicker [--grid] [--persist] [--inc-extension] [--command COMMAND] [DIR]

Positional arguments:
  DIR

Options:
  --grid, -g             Display wallpapers in a grid.
  --persist, -p          Persist, remain after choosing a wallpaper.
  --inc-extension, -e    Include file extension as second argument to the target command or script. ie ($2)
  --command COMMAND, -c COMMAND
                         Settings command, expects a command to run [command] [image path] [~extension~]. [default: feh --bg-fill]
  --help, -h             display this help and exit
```

## Table of Contents

- [Features](#features)
- [Setup](#setup)
- [Install](#install)
- [Usage](#usage)
- [Hyprpaper Example](#hyprpaper-config-example)
- [Why Another Wallpaper picker?](#why-another-wallpaper-picker)
- [Contributing](#contributing)

## Features

- Persist option, remains open until closed. For when you can't make up your mind. :)
- Configurable command. Defaults to `feh --bg-fill` but can be configured to any command of your choice.
- Display wallpaper previews in a grid or wide list layout.

## Setup

- Install golang, you can download it from the official [Go website](https://go.dev/doc/install).
- Set up [fyne](https://docs.fyne.io/started/).

Clone the repository:

```sh
git clone https://github.com/iiiz/wallpicker.git
cd wallpicker
go mod tidy
```

To build the project:

```sh
make build
```

## Install

Default install location is `/usr/bin/wallpicker`

```sh
go mod download

# clean & build
make
make install
```

## Usage

To run the project:

```sh
# development
go run . /some/dir/with/wallpapers

# build
make

# run
./bin/wallpicker /some/dir/with/wallpapers
```

## Hyprpaper config example

For single monitor or mirrored. `wallpicker -c ~/.wallpaper/script/set.sh ~/Pictures/wallpapers`

```sh
#!/bin/bash

hyprctl hyprpaper wallpaper "DP-1,$1,cover"
```

For multi monitor setup. `wallpicker -e -c ~/.wallpaper/script/set.sh ~/Pictures/wallpapers`

```sh
#!/bin/bash

IMAGE_PATH=$1
IMAGE_EXT=$2

# Clean up dir of type and convert to match monitors
rm ~/.wallpaper/*$IMAGE_EXT
magick "$IMAGE_PATH" -crop 33.33%x100% ~/.wallpaper/wall$IMAGE_EXT

# set wallpaper
hyprctl hyprpaper wallpaper "DP-2,/home/iiiz/.wallpaper/wall-0$IMAGE_EXT,cover"
hyprctl hyprpaper wallpaper "DP-3,/home/iiiz/.wallpaper/wall-1$IMAGE_EXT,cover"
hyprctl hyprpaper wallpaper "DP-1,/home/iiiz/.wallpaper/wall-2$IMAGE_EXT,cover"
```

## Why another wallpaper picker

I use a very minimal ~XMonad~ hyprland desktop config and wanted a visual wallpaper picker to bind to a shortcut.
Other options didn't really fit what I was looking for so here we are. ¯\_(ツ)\_/¯

## Contributing

Contributions and forks are welcome!
