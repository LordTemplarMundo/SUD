package main

import (
	"fmt"
)

func main() {
	fmt.Println("Welcome to the cool game.")
	rooms := CreateMap("testmap")
	world := newWorld(rooms)
	world.startWorld()
	player := newMob("Mysterious Stranger", world)
	var input string
	for world.running {
		fmt.Print(">")
		fmt.Scan(&input)
		availableActions := append(player.commands, player.location.getExitCommands()...)
		action := parseInput(input, availableActions)
		player.cmdQueue = append(player.cmdQueue, action)
	}
}
