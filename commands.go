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
			p.print <- p.location.displayRoom()
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
			p.despawn()
			disconnectUserFromMob(p)
			return true
		},
	}
}

func sayCommand() Command {
	return Command{
		names: []string{"say", "'"},
		action: func(p *Mob) bool {
			world.roomEmit(fmt.Sprintf("%v makes a noise!\n", p.getName()), p.location)
			return true
		},
	}
}

func noCommandAction(p *Mob) bool {
	p.print <- "I don't know how to do that!\n"
	return true
}

func generateExitAction(exit *Exit) Cmd {
	return func(m *Mob) bool {
		roomLeft := exit.room.leaveRoom(m)
		if !roomLeft {
			m.print <- "You can't get out of here!"
		}
		world.roomEmit(fmt.Sprintf("%v leaves to the %v.\n", m.name, exit.getPrimaryName()), exit.room)
		roomEntered := exit.destination.room.enterRoom(m)
		if !roomEntered {
			m.print <- "You can't get in there!"
		}
		world.roomEmit(fmt.Sprintf("%v enters from the %v.\n", m.name, exit.destination.getPrimaryName()), exit.destination.room)
		m.print <- exit.destination.room.displayRoom()
		return roomEntered
	}
}

func basicCommands() (output []Command) {
	output = append(output, []Command{lookCommand(), exitCommand(), quitCommand(), sayCommand()}...)
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
