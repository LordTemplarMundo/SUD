package main

import (
	log "github.com/sirupsen/logrus"
)

type Mob struct {
	location    *Room
	name        string
	commands    []Command
	world       *World
	pulse       <-chan interface{}
	cmdQueue    []func(*Mob) bool
	print       chan<- string
	description string
}

type Pulsable interface {
	Register(<-chan interface{})
}

func (m *Mob) Register(pulse <-chan interface{}) {
	m.pulse = pulse
}

func newMob() (p *Mob) {
	return &Mob{
		name:        "",
		commands:    basicCommands(),
		description: "A generic looking person.",
	}
}

func (m *Mob) spawn(name string, world *World) error {
	m.name = name
	start := world.getStartRoom()
	pulse, err := world.registerThing(m)
	if err != nil {
		log.WithError(err).Errorf("Failed to create mob: %s", name)
		return err
	}
	m.pulse = pulse
	go m.beat()
	start.enterRoom(m)
	log.WithFields(log.Fields{
		"mob_name":      m.name,
		"starting_room": m.location.name,
	}).Info("Mob spawned.")
	return nil
}

func (m *Mob) connect(print chan<- string) {
	m.print = print
}

func (m *Mob) getDescription() string {
	return m.description
}

func (m *Mob) getName() string {
	return m.name
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
