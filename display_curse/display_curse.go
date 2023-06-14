package main

import (
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//go:embed curse.csv
var curseCSV string

//go:embed tweak.csv
var tweakCSV string

//go:embed blessing.csv
var blessingCSV string

type Modifier struct {
	Name string
	Tier int
}

func loadModifier(reader io.Reader) ([]Modifier, error) {
	var mods []Modifier
	r := csv.NewReader(reader)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		tier, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}
		mods = append(mods, Modifier{Name: normalizeName(record[0]), Tier: tier})
	}
	return mods, nil
}

func LoadModifiers() ([]Modifier, error) {
	var modifiers []Modifier
	curses, err := loadModifier(strings.NewReader(curseCSV))
	if err != nil {
		return nil, err
	}
	modifiers = append(modifiers, curses...)

	tweaks, err := loadModifier(strings.NewReader(tweakCSV))
	if err != nil {
		return nil, err
	}
	modifiers = append(modifiers, tweaks...)

	blessings, err := loadModifier(strings.NewReader(blessingCSV))
	if err != nil {
		return nil, err
	}
	modifiers = append(modifiers, blessings...)

	return modifiers, nil
}

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

func normalizeName(name string) string {
	return regexp.MustCompile(`\[\w+]`).ReplaceAllString(name, "")
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

type Data struct {
	D Modifiers `json:"d"`
}

type Modifiers struct {
	M []string `json:"m"`
}

func writeCurses(file string, out_file string, gamemode string, modifierTiers map[string]int) {
	xmlFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer xmlFile.Close()

	byteValue, err := io.ReadAll(xmlFile)
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

	var modifiers []string
	for _, entry := range prefs.Entries {
		if entry.Key == gamemode {
			var data Data
			json.Unmarshal([]byte(entry.Value), &data)
			for _, modifier := range data.D.M {
				modifier = normalizeName(modifier)
				tier, ok := modifierTiers[modifier]
				if ok {
					modifiers = append(modifiers, fmt.Sprintf("[%d] %s", tier, modifier))
				} else {
					modifiers = append(modifiers, modifier)
				}
			}
			break
		}
	}

	out := strings.Join(modifiers, "\n")
	err = os.WriteFile(out_file, []byte(out), os.ModePerm)
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
		os.Exit(1)
	}

	modifiers, err := LoadModifiers()
	if err != nil {
		log.Fatal(err)
	}

	modifierTiers := make(map[string]int)
	for _, mod := range modifiers {
		modifierTiers[mod.Name] = mod.Tier
	}

	watcher := NewFileWatcher(*input)
	watcher.watch(func() { writeCurses(*input, *output, *gamemode, modifierTiers) })
}
