package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"
    "strings"
)

type FileWatcher struct {
	ts   int64
	path string
}

func NewFileWatcher(path string) *FileWatcher {
	return &FileWatcher{
		ts:   0,
		path: path,
	}
}

func (fw *FileWatcher) updated() bool {
	stat, err := os.Stat(fw.path)
	if err != nil {
		return false
	}
	stamp := stat.ModTime().UnixNano()
	if stamp != fw.ts {
		fw.ts = stamp
		return true
	}
	return false
}

func (fw *FileWatcher) watch(fn func()) {
	for {
		if fw.updated() {
			fmt.Printf("%s changed\n", fw.path)
			fn()
		}
		time.Sleep(time.Second)
	}
}

type SliceAndDicePrefs struct {
	XMLName xml.Name `xml:"properties"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",innerxml"`
}

type CursesData struct {
	D Curses `json:"d"`
}

type Curses struct {
	M []string `json:"m"`
}

func writeCurses(file string, out_file string, gamemode string) {
	xmlFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
    if err != nil {
		fmt.Println(err)
		return
    }

	var prefs SliceAndDicePrefs
	err = xml.Unmarshal(byteValue, &prefs)
    if err != nil {
		fmt.Println(err)
		return
    }

	var curses []string
	for _, entry := range prefs.Entries {
		if entry.Key == gamemode {
			var data CursesData
			json.Unmarshal([]byte(entry.Value), &data)
			for _, curse := range data.D.M {
				curses = append(curses, regexp.MustCompile(`\[\w+\]`).ReplaceAllString(curse, ""))
			}
			break
		}
	}

    out := strings.Join(curses, "\n")
    err = ioutil.WriteFile(out_file, []byte(out), os.ModePerm)
    if err != nil {
        fmt.Println(err)
        return
    }
}

func main() {
	input := flag.String("input", filepath.Join(os.Getenv("HOME"), ".prefs", "slice-and-dice-2"), "Path to the slice and dice save file")
	output := flag.String("output", "", "Path to output the current curse text file")
	gamemode := flag.String("gamemode", "classic", "The gamemode to use")
	flag.Parse()

	if *output == "" {
		fmt.Println("Output path must be specified")
		return
	}

	watcher := NewFileWatcher(*input)
	watcher.watch(func() { writeCurses(*input, *output, *gamemode) })
}
