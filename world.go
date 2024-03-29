package main

type world struct {
	rooms   []*Room
	running bool
}

func newWorld(rooms []*Room) *world {
	return &world{
		rooms:   rooms,
		running: false,
	}
}

func newTestWorld() *world {
	return &world{
		rooms:   makeTestDungeon(),
		running: false,
	}
}

func (w *world) getStartRoom() *Room {
	return w.rooms[0]
}

func (w *world) startWorld() {
	w.running = true
}

func (w *world) stopWorld() {
	w.running = false
}

func makeTestDungeon() []*Room {
	sideRoom := newUnlinkedRoom("This room is unimportant. Go away.", "Side Room")
	endRoom := newUnlinkedRoom("This room is marginally better. Maybe.", "End Room")
	startRoom := newUnlinkedRoom("This room sucks. Really bad. Wow!", "Start Room")
	connectRooms(startRoom, endRoom, []string{"north", "n"}, []string{"south", "s"})
	connectRooms(startRoom, sideRoom, []string{"west", "w"}, []string{"east", "e"})
	return []*Room{startRoom, endRoom, sideRoom}
}
