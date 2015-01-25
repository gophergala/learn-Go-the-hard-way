package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Cheat struct {
	Name     string
	Contents []Content
}

type Content struct {
	Comment string
	Command string
}

func CheatSheet(command string) string {
	file, err := os.Open("./commands.json")
	if err != nil {
		panic(err)
	}
	var cheats []Cheat
	dec := json.NewDecoder(file)
	dec.Decode(&cheats)
	var out string
	for _, v := range cheats {
		if v.Name == command {
			for _, v := range v.Contents {
				out += fmt.Sprintf("%s\n", v.Comment)
				out += fmt.Sprintf("%s\n", v.Command)
			}
		}
	}
	return out
}

func main() {
	args := os.Args
	if len(args) != 1 {
		log.Fatal("want one  argument")
	}

	println(`Unix has a lot of commands to remenber.
To help us search the command quickly,we will create a small cheat sheet command.
We will store the commands as json.In this exercise you can play with Go IO and json encoding.`)
}
