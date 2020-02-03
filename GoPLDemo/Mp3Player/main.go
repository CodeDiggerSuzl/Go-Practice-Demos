package main

import (
	"bufio"
	"fmt"
	"mlib"
	"mp"
	"os"
	"strconv"
	"strings"
)

var lib *mlib.MusicManager
var id int = 1

var ctrl, signal chan int

func main() {
	p("start")
	lib = mlib.NewMusicManager()
	r := bufio.NewReader(os.Stdin)
	for {
		p("Enter command -> ")
		rawLine, _, _ := r.ReadLine()
		line := string(rawLine)
		if line == "q" || line == "e" {
			break
		}
		tokens := strings.Split(line, " ")
		if tokens[0] == "lib" {
			handleLibCommands(tokens)
		} else if tokens[0] == "play" {
			handlePlayCommand(tokens)
		} else {
			p("Unrecognized command:", tokens[0])
		}
	}
}

var p = fmt.Println

func handleLibCommands(tokens []string) {
	switch tokens[1] {
	case "list":
		for i := 0; i < lib.Len(); i++ {
			e, _ := lib.Get(i)
			p(i+1, ":", e.Name, e.Artist, e.Source, e.Type)
		}
	case "add":
		if len(tokens) == 6 {
			id++
			lib.Add(&mlib.MusicEntry{strconv.Itoa(id), tokens[2], tokens[3], tokens[4], tokens[5]})
		} else {
			p("USAGE: lib add <name> <artist> <source> <type>")
		}
	case "remove":
		if len(tokens) == 3 {
			lib.Remove(id)
		} else {
			p("USAGE: lib remove <id>")
		}
	default:
		p("Unrecognized commands")
	}

}

func handlePlayCommand(tokens []string) {
	if len(tokens) == 2 {
		p("USAGE: play <name>")
		return
	}
	a := lib.Find(tokens[1])
	if a == nil {
		p("Music", tokens[1], "does not exist")
		return
	}
	mp.PlayMusic(a.Source, a.Type)
}
