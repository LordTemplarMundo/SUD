package main

import (
	"fmt"
	"sync"
	"time"
)

type World struct {
	sync.Mutex
	rooms   []*Room
	running bool
	things  []chan interface{}
}

func newWorld(rooms []*Room) *World {
	return &World{
		rooms:   rooms,
		running: false,
	}
}

func (w *World) getStartRoom() *Room {
	return w.rooms[0]
}

func (w *World) registerThing(p Pulsable) (pulse <-chan interface{}, err error) {
	if w.running {
		w.Lock()
		defer w.Unlock()
		pulse := make(chan interface{}, 10)
		p.Register(pulse)
		w.things = append(w.things, pulse)
		return pulse, err
	} else {
		return nil, fmt.Errorf("Unable to register - world is not running.")
	}
}

func (w *World) startWorld() {
	w.running = true
	go w.beat()
}

func (w *World) stopWorld() {
	w.running = false
	w.Mutex.Lock()
	for _, thing := range w.things {
		close(thing)
	}
	w.Mutex.Unlock()
}

func (w *World) beat() {
	for {
		time.Sleep(time.Millisecond)
		w.Mutex.Lock()
		for _, thing := range w.things {
			thing <- true
		}
		w.Mutex.Unlock()
	}
}

func (w *World) roomEmit(sound string, location *Room) {
	usersLock.Lock()
	for _, user := range users {
		if user.Mob.location.name == location.name {
			user.Mob.print <- sound
		}
	}
	usersLock.Unlock()
}
