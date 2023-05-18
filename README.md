# Slice and dice tools

A set of slice and dice utilities

## Display curse

This tool will read the curses for your current run and output it to a text file.
The text file can be read by programs such as OBS to display your current
curses to viewers.

### Usage

Download the program from the [releases page](https://github.com/MaT1g3R/slice-and-dice-tools/releases/) based on your
OS and CPU architecture. For example, if you are using an x86 Windows PC, download `slice-and-dice-tools_Windows_x86_64.zip`

Unpack the downloaded archive and put the `display-curse` program somewhere.

The options to the program are:

```bash
Usage display-curse
  -gamemode string
        The gamemode to use (default "classic")
  -input string
        Path to the slice and dice save file (default "$HOME/.prefs/slice-and-dice-2")
  -output string
        Path to output the current curse text file
```

You will want to set the `-output` option to some file location that OBS will use as a text source. On windows the default
`-input` location might not work, the save file might be located at: `C/Users/YOUR_NAME/.prefs/slice-and-dice-2`

As an example, running this program on Windows via either powershell or cmd:

```bash
./display-curse -input C/Users/YOUR_NAME/.prefs/slice-and-dice-2 -output C/Users/YOUR_NAME/some-file.txt
```

Then you can have OBS read `C/Users/YOUR_NAME/some-file.txt` to display curses on stream.
