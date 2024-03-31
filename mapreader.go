package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// The text-based representation of a map ingested from a file.
type TextMap struct {
	width, height int
	area          []string
}

// Prints a TextMap's string array into a square matching the TexMap's width
// and height.
func (tm *TextMap) printArea() {
	var cY int
	for i, symbol := range tm.area {
		coords := getCoords(i, tm.width, tm.height)
		if coords.y > cY {
			cY = coords.y
			fmt.Print("\n")
		}
		fmt.Print(symbol)
	}
	fmt.Print("\n")
}

// Loads a map in from a file located at fileDir
func readMap(fileDir string) (m TextMap, err error) {
	// First, open the indicated file, if there's an error, just crash out.
	if file, err := os.Open(fileDir); err == nil {
		defer file.Close() // Close the file when done.
		read := make([]byte, 1)
		rawMap := [][]string{}
	out:
		for { // Until interrupted...
			row := []string{} // ...make a new 'row' array...
			for {             // ...and until interrupted...
				_, err := io.ReadAtLeast(file, read, 1) // ...read 1 byte at a time (one character)
				// If the character is a newline (end of row)
				// OR the reader reaches the end of the file:
				if string(read) == "\n" || err != nil {
					rawMap = append(rawMap, row) // add to the row to the list of rows
					if len(row) > m.width {      // If the width of the new row is bigger than the current width
						m.width = len(row) // Then set it to be the current width
					}
					row = []string{} // create a new row
					if err != nil {  // If the reader has reached the end of the file:
						break out // Break out of the 'out:' loop
					}
				} else { // If the character is any other value
					row = append(row, string(read)) // Add it to the row
				}
			}
		}
		m.height = len(rawMap)                                       // The height of the map is how many rows it has.
		m.area = combineArrays(padArrayStringArray(m.width, rawMap)) // Flatten the array of arrays into a single []string after making each row the same length
		return m, nil
	} else {
		return m, err
	}
}

// Extends input array to length 'width', if array is already equal or larger
// than width, just returns array.
func padStringArray(width int, array []string) (output []string) {
	if len(array) < width {
		return append(array, make([]string, width-len(array))...)
	}
	return array
}

// Calls 'padStringArray' on all arrays in a 2-dimensional string array.
func padArrayStringArray(width int, arrays [][]string) (output [][]string) {
	for _, array := range arrays {
		output = append(output, padStringArray(width, array))
	}
	return
}

// Collapses down a 2-dimensional string array to a 1-dimensional array.
func combineArrays(arrays [][]string) (output []string) {
	for _, array := range arrays {
		output = append(output, array...)
	}
	return
}

// Converts the input index to an x,y touple from a 2-dimensional grid of
// height and width.
func getCoords(index int, width, height int) Coordinates {
	return Coordinates{index % width, index / width}
}

// Converts the input coordinate (x, y touple) from a 2-dimensional grid into
// a 1-dimensional index.
func getIndex(coords Coordinates, width int) (index int) {
	return coords.x + coords.y*width
}

// Returns true if input string is empty or a space.
// These symbols are used in TextMaps to indicate empty space on the map.
func isEmpty(symbol string) bool {
	return symbol == "" || symbol == " "
}

// Returns true if the input string is "|" or "-".
// These symbols are used in TextMaps to indicate vertical or horizontal
// room connetions respectively.
func isExit(symbol string) bool {
	return symbol == "|" || symbol == "-"
}

// Returns true if the input string is any other character than and empty
// string, a space, a "|" or a "-".
func isRoom(symbol string) bool {
	return !isEmpty(symbol) && !isExit(symbol)
}

// Creates a complex of Room objects from a TextMap and a collection of TextRooms.
func (tm *TextMap) buildMap(rooms map[string]TextRoom) []*Room {
	mapWorker := CreateMapWorker(tm, rooms) // Creates a 'worker' to manage mape creation.
	for i, symbol := range tm.area {        // For every character in the []string array of the text map
		if isRoom(symbol) {
			mapWorker.buildRoom(i)
		} else if isExit(symbol) {
			mapWorker.createLink(i)
		}
	}
	mapWorker.joinLinks()          // Connect together the rooms.
	output := mapWorker.getRooms() // Convert the map of indices to rooms to a flat room array.
	return output
}

// Struct indicating a horizontal and vertical position on a 2-dimensional grid
type Coordinates struct {
	x, y int
}

// Returns bitflags for directions that are abutting the edge of a
// 2-dimensional grid.
func getEdgeFlags(tm *TextMap, index int, coords Coordinates) (dirFlags Direction) {
	if coords.y <= 0 {
		dirFlags += North
	}
	if coords.x >= tm.width-1 {
		dirFlags += East
	}
	if coords.y >= tm.height-1 {
		dirFlags += South
	}
	if coords.x <= 0 {
		dirFlags += West
	}
	return
}

// MapWorkers coordinate the steps involved in converting a text map into a
// collection of interconnected
type MapWorker struct {
	textMap        *TextMap            // The text map being converted
	completedRooms map[int]*Room       // A map of indices to Room structures
	links          []Link              // A list of connections between Room structures
	rooms          map[string]TextRoom // The 'content' to be loaded into Room structures
}

// Links represent connections between Rooms and are built in parallel to
// them.
// When all Links and Rooms have been instanced, Links are used to connect them
// together.
type Link struct {
	start     int       // Index of one side of the connection.
	end       int       // Index of the opposite side of the connection.
	direction Direction // Direction from 'start' to 'end'.
}

// Create a new MapWorker. Probably doesn't need its own function.
func CreateMapWorker(tm *TextMap, rooms map[string]TextRoom) (mw *MapWorker) {
	mw = &MapWorker{
		textMap:        tm,
		links:          []Link{},
		completedRooms: make(map[int]*Room),
		rooms:          rooms,
	}
	return mw
}

// Uses the array of Links and map of indices to Rooms to connect the map's
// Rooms together.
func (mw *MapWorker) joinLinks() {
	for _, link := range mw.links {
		startRoom, startOk := mw.completedRooms[link.start]
		endRoom, endOk := mw.completedRooms[link.end]
		if startOk && endOk && (link.direction != BadDir) {
			connectRoomsCardinally(startRoom, endRoom, link.direction)
		}
	}
}

// Builds a room at a given index.
func (mw *MapWorker) buildRoom(index int) {
	var output *Room
	if room, exists := mw.rooms[mw.textMap.area[index]]; exists {
		output = newUnlinkedRoom(room.Description, room.Title)
	} else {
		output = newGenericRoom()
	}
	mw.completedRooms[index] = output
}

// Creates a Link between either ends of a connection and determines its
// direction.
func (mw *MapWorker) createLink(index int) {
	dir := East + West
	if mw.textMap.area[index] == "|" {
		dir = North + South
	}
	link := []int{}
	for _, dirInt := range flagToIndicies(mw.textMap, index, dir) {
		link = append(link, dirInt)
	}
	if len(link) == 2 {
		start := link[0]
		end := link[1]
		mw.links = append(mw.links, Link{start, end, cardinalDirBetween(start, end, mw.textMap.width)})
	}
}

// Converts the MapWorker's map of text map indices to Rooms to a flat array
// of Room structs.
func (mw *MapWorker) getRooms() (output []*Room) {
	for _, room := range mw.completedRooms {
		output = append(output, room)
	}
	return
}

// TextRoom is the serialised format of Room descriptions, etc.
type TextRoom struct {
	Symbol      string `yaml:"symbol"` // Symbol is the single character on the Text Map that this TextRoom will be used for
	Title       string `yaml:"title"`
	Description string `yaml:"desc"`
}

// This collects a set of serialised room descriptions and symbols from a file
// at location 'dir' and creates a map of symbol to TextRoom.
func readRooms(dir string) (rooms map[string]TextRoom) {
	rooms = make(map[string]TextRoom)
	if rawRooms, err := os.Open(dir); err == nil {
		defer rawRooms.Close()
		decoder := yaml.NewDecoder(rawRooms)
		var rawRoom TextRoom
		for decoder.Decode(&rawRoom) == nil {
			rooms[rawRoom.Symbol] = rawRoom
		}
	}
	return
}
