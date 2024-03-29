package main

import (
	"fmt"
	"io"
	"os"
	"sync"
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
	if file, err := os.Open(fileDir); err == nil {
		read := make([]byte, 1)
		rawMap := [][]string{}
	out:
		for {
			row := []string{}
			for {
				_, err := io.ReadAtLeast(file, read, 1)
				if string(read) == "\n" || err != nil {
					rawMap = append(rawMap, row)
					row = []string{}
					if err != nil {
						break out
					}
				} else {
					row = append(row, string(read))
					if len(row) > m.width {
						m.width = len(row)
					}
					continue
				}
			}
		}
		m.height = len(rawMap)
		m.area = combineArrays(padArrayStringArray(m.width, rawMap))
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

func (tm *TextMap) buildRooms() []*Room {
	tm.printArea()
	mapWorker := CreateMapWorker(tm)
	for i, symbol := range tm.area {
		if isRoom(symbol) {
			mapWorker.buildRoom(i)
		} else if isExit(symbol) {
			mapWorker.buildExit(i)
		}
	}
	mapWorker.joinLinks()
	output := mapWorker.getRooms()
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

type MapWorker struct {
	workerCount    int
	waitGroup      *sync.WaitGroup
	textMap        *TextMap
	completedRooms map[int]*Room
	links          []Link
}

type Link struct {
	start     int
	end       int
	direction Direction
}

func CreateMapWorker(tm *TextMap) (mw *MapWorker) {
	mw = &MapWorker{
		workerCount:    0,
		textMap:        tm,
		links:          []Link{},
		completedRooms: make(map[int]*Room),
	}
	return mw
}

func (mw *MapWorker) joinLinks() {
	for _, link := range mw.links {
		startRoom, startOk := mw.completedRooms[link.start]
		endRoom, endOk := mw.completedRooms[link.end]
		if startOk && endOk && (link.direction != BadDir) {
			connectRoomsCardinally(startRoom, endRoom, link.direction)
		}
	}
}

func (mw *MapWorker) buildRoom(index int) {
	output := newGenericRoom()
	output.name = mw.textMap.area[index]
	mw.completedRooms[index] = output
}

func (mw *MapWorker) buildExit(index int) {
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

func (mw *MapWorker) getRooms() (output []*Room) {
	for _, room := range mw.completedRooms {
		output = append(output, room)
	}
	return
}
