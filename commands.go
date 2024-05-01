package main

import (
	"fmt"
	"strings"
)

type Cmd = func(*Mob) bool

type Command struct {
	names  []string
	action Cmd
}

func lookCommand() Command {
	return Command{
		names: []string{"look", "l"},
		action: func(p *Mob) bool {
			p.print <- p.location.show()
			return true
		},
	}
}

func exitCommand() Command {
	return Command{
		names: []string{"exits", "doors", "dirs"},
		action: func(p *Mob) bool {
			p.print <- p.location.listExits()
			return true
		},
	}
}

func quitCommand() Command {
	return Command{
		names: []string{"quit", "q"},
		action: func(p *Mob) bool {
			fmt.Println("Goodbye!")
			p.world.stopWorld()
			return true
		},
	}
}

func noCommandAction(p *Mob) bool {
	p.print <- "I don't know how to do that!"
	return true
}

func generateExitAction(exit *Exit) Cmd {
	return func(m *Mob) bool {
		roomEntered := exit.destination.room.enterRoom(m)
		if !roomEntered {
			m.print <- "You can't go that way!"
		}
		exit.room.leaveRoom(m)
		return roomEntered
	}
}

func basicCommands() (output []Command) {
	output = append(output, []Command{lookCommand(), exitCommand(), quitCommand()}...)
	return
}

func parseInput(input string, commands []Command) Cmd {
	for _, cmd := range commands {
		for _, alias := range cmd.names {
			if strings.ToLower(alias) == strings.ToLower(input) {
				return cmd.action
			}
		}
	}
	return noCommandAction
}
