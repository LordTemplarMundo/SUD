package main

import (
	"fmt"
	"strings"
	"sync"
)

type Room struct {
	sync.RWMutex
	description string
	name        string
	exits       []*Exit
	commands    []Command
	contents    []Thing
}

func (r *Room) show() string {
	padding := strings.Repeat("-", len(r.description))
	return fmt.Sprintf("\n|%v|\n%v\n%v\n%v\n\n", r.name, padding, r.description, padding)
}

func (r *Room) getExitCommands() (commands []Command) {
	// We can filter out elements based on visibility, etc. later
	for _, exit := range r.exits {
		commands = append(commands, exit.generateCommands())
	}
	return
}

func (r *Room) listExits() string {
	var exitString string
	for _, exit := range r.exits {
		exitString += exit.getPrimaryName() + ", "
	}
	if len(r.exits) > 0 {
		return fmt.Sprintf("Visible Exits: %v\n", exitString)
	}
	return "You can't see any exits."
}

func newUnlinkedRoom(description string, name string) *Room {
	return &Room{sync.RWMutex{}, description, name, []*Exit{}, []Command{}, []Thing{}}
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
	p.print <- fmt.Sprintf("%v\n%v\n", r.show(), r.listExits())
	r.RWMutex.Lock()
	r.contents = append(r.contents, p)
	r.RWMutex.Unlock()
	// Later on, if something stops the movement, return false.
	return true
}

func (r *Room) leaveRoom(p *Mob) bool {

	for pos, thing := range r.contents {
		if thing.getName() == p.getName() {
			r.RWMutex.Lock()
			r.contents = append(r.contents[0:pos], r.contents[pos:len(r.contents)-1]...)
			r.RWMutex.Unlock()
		}
	}
	return true
}
