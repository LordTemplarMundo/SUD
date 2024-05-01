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

func newTestWorld() *World {
	return &World{
		rooms:   makeTestDungeon(),
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

func makeTestDungeon() []*Room {
	sideRoom := newUnlinkedRoom("This room is unimportant. Go away.", "Side Room")
	endRoom := newUnlinkedRoom("This room is marginally better. Maybe.", "End Room")
	startRoom := newUnlinkedRoom("This room sucks. Really bad. Wow!", "Start Room")
	connectRooms(startRoom, endRoom, []string{"north", "n"}, []string{"south", "s"})
	connectRooms(startRoom, sideRoom, []string{"west", "w"}, []string{"east", "e"})
	return []*Room{startRoom, endRoom, sideRoom}
}

func (w *World) beat() {
	for {
		time.Sleep(time.Second)
		w.Mutex.Lock()
		for _, thing := range w.things {
			thing <- true
		}
		w.Mutex.Unlock()
	}
}
