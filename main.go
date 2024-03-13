package main

import (
	"errors"
	"fmt"
	"strings"
)

type room struct {
	description string
	name        string
	exits       map[string]*room
}

func (r *room) show() {
	padding := strings.Repeat("-", len(r.description))
	fmt.Printf("\n|%v|\n%v\n%v\n%v\n\n", r.name, padding, r.description, padding)
}

func (r *room) listExits() {
	var exitString string
	for exit := range r.exits {
		exitString += exit + ", "
	}
	if len(r.exits) > 0 {
		fmt.Printf("Visible Exits: %v\n", exitString)
	}
}

func linkRooms(origin *room, dest *room, dir string) error {
	var antiDir string
	switch dir {
	case "north":
		antiDir = "south"
	case "south":
		antiDir = "north"
	case "east":
		antiDir = "west"
	case "west":
		antiDir = "east"
	case "up":
		antiDir = "down"
	case "down":
		antiDir = "up"
	default:
		return errors.New("Direction not recognised!")
	}
	origin.exits[dir] = dest
	dest.exits[antiDir] = origin
	return nil
}

func newRoom(description string, name string) *room {
	return &room{description, name, make(map[string]*room)}
}

func (r *room) enterRoom() {
	r.show()
	r.listExits()
}

type player struct {
	location *room
	name     string
}

func (p *player) move(dir string) {
	if dest, ok := p.location.exits[dir]; ok {
		fmt.Printf("You move %v.\n", dir)
		p.location = dest
		p.location.enterRoom()
	} else {
		fmt.Printf("You can't go that way!\n")
	}
}

func main() {
	fmt.Println("Welcome to the cool game.")
	sideRoom := newRoom("This room is unimportant. Go away.", "Side Room")
	endRoom := newRoom("This room is marginally better. Maybe.", "End Room")
	startRoom := newRoom("This room sucks. Really bad. Wow!", "Start Room")
	linkRooms(startRoom, endRoom, "north")
	linkRooms(startRoom, sideRoom, "west")
	player := player{location: startRoom, name: "Mysterious Stranger"}
	var i string
	player.location.enterRoom()
	for i != "quit" && i != "q" {
		fmt.Println("What you want to do lol???")
		fmt.Print(">")
		fmt.Scan(&i)
		switch i {
		case "look":
			player.location.show()
		case "exits":
			player.location.listExits()
		default:
			player.move(i)
		}
	}
}
