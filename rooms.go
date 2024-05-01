package main

import (
	"fmt"
	"strings"
)

type Room struct {
	description string
	name        string
	exits       []*Exit
	commands    []Command
}

func (r *Room) show() {
	padding := strings.Repeat("-", len(r.description))
	fmt.Printf("\n|%v|\n%v\n%v\n%v\n\n", r.name, padding, r.description, padding)
}

func (r *Room) getExitCommands() (commands []Command) {
	// We can filter out elements based on visibility, etc. later
	for _, exit := range r.exits {
		commands = append(commands, exit.generateCommands())
	}
	return
}

func (r *Room) listExits() {
	var exitString string
	for _, exit := range r.exits {
		exitString += exit.getPrimaryName() + ", "
	}
	if len(r.exits) > 0 {
		fmt.Printf("Visible Exits: %v\n", exitString)
	}
}

func newUnlinkedRoom(description string, name string) *Room {
	return &Room{description, name, []*Exit{}, []Command{}}
}

func newGenericRoom() *Room {
	return newUnlinkedRoom("This is a boring generic room.", "Generic Room")
}

func newLinkedRoom(description string, name string, exits []*Exit) *Room {
	room := newUnlinkedRoom(description, name)
	room.exits = append(room.exits, exits...)
	room.updateCommands()
	return room
}

func (r *Room) updateCommands() {
	r.commands = r.getExitCommands()
}

func (r *Room) enterRoom(p *Mob) bool {
	p.location = r
	r.show()
	r.listExits()
	// Later on, if something stops the movement, return false.
	return true
}
