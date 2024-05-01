package main

type Exit struct {
	names       []string
	room        *Room
	destination *Exit
}

func connectRooms(start *Room, end *Room, startDir []string, endDir []string) {
	startSide := Exit{names: startDir, room: start}
	endSide := Exit{names: endDir, room: end}
	startSide.destination = &endSide
	endSide.destination = &startSide
	start.exits = append(start.exits, &startSide)
	end.exits = append(end.exits, &endSide)
}

func connectRoomsCardinally(start *Room, end *Room, dir Direction) {
	connectRooms(start, end, dirToCommandStrings(dir), dirToCommandStrings(invertDir(dir)))
}

func (e *Exit) getDestination() *Room {
	return e.destination.room
}

func (e *Exit) getPrimaryName() string {
	if len(e.names) >= 1 {
		return e.names[0]
	}
	return "Void"
}

func (e *Exit) generateCommands() Command {
	return Command{
		names:  e.names,
		action: generateExitAction(e),
	}
}
