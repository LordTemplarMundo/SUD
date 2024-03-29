package main

type Player struct {
	location *Room
	name     string
	commands []Command
	world    *world
}

func newPlayer(name string, world *world) (p *Player) {
	start := world.getStartRoom()
	p = &Player{
		name:     name,
		location: start,
		commands: basicCommands(),
		world:    world,
	}
	start.enterRoom(p)
	return p
}
