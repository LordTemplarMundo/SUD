package main

import (
	log "github.com/sirupsen/logrus"
)

type Mob struct {
	location *Room
	name     string
	commands []Command
	world    *World
	pulse    <-chan interface{}
	cmdQueue []func(*Mob) bool
}

type Pulsable interface {
	Register(<-chan interface{})
}

func (m *Mob) Register(pulse <-chan interface{}) {
	m.pulse = pulse
}

func newMob(name string, world *World) (p *Mob) {
	start := world.getStartRoom()
	p = &Mob{
		name:     name,
		location: start,
		commands: basicCommands(),
		world:    world,
	}
	pulse, err := world.registerThing(p)
	if err != nil {
		log.WithError(err).Errorf("Failed to create mob: %s", name)
		return p
	}
	p.pulse = pulse
	start.enterRoom(p)
	log.WithFields(log.Fields{
		"mob_name":      name,
		"starting_room": start.name,
	}).Info("Mob instanced.")
	go p.beat()
	return p
}

func (m *Mob) beat() {
	for {
		select {
		case _ = <-m.pulse:
			if len(m.cmdQueue) >= 1 {
				var nextCommand Cmd
				if len(m.cmdQueue) >= 2 {
					nextCommand, m.cmdQueue = m.cmdQueue[0], m.cmdQueue[1:]
				} else {
					nextCommand, m.cmdQueue = m.cmdQueue[0], []Cmd{}
				}
				nextCommand(m)
			}
		}
	}
}
