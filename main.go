package main

import (
	"fmt"
)

func main() {
	fmt.Println("Welcome to the cool game.")
	area, _ := readMap("./testmap.txt")
	rooms := readRooms("./testrooms.txt")
	world := newWorld(area.buildRooms(rooms))
	player := newPlayer("Mysterious Stranger", world)
	world.startWorld()
	var input string
	for world.running {
		fmt.Print(">")
		fmt.Scan(&input)
		availableActions := append(player.commands, player.location.getExitCommands()...)
		action := parseInput(input, availableActions)
		action(player)
	}
}
