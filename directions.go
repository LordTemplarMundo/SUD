package main

import (
	"fmt"
)

// Bitflags for cardinal directions (North, South, East or West)
type Direction byte

const (
	BadDir  Direction = 0  //0000 0000 Used when not recognised.
	North   Direction = 1  //0000 0001
	East    Direction = 2  //0000 0010
	South   Direction = 4  //0000 0100
	West    Direction = 8  //0000 1000
	DirMask Direction = 15 //0000 1111
)

func dirToString(dir Direction) string {
	switch dir {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	}
	return "NaD"
}

func dirToCommandStrings(dir Direction) []string {
	full := dirToString(dir)
	return []string{full, string(full[0])}
}

func invertDir(dir Direction) Direction {
	switch dir {
	case North:
		return South
	case East:
		return West
	case South:
		return North
	case West:
		return East
	}
	return BadDir
}

// Get the cardinal direction to i2 from i1
func cardinalDirBetween(i1, i2, width int) Direction {
	if (i1-i2)%width == 0 {
		if i1 > i2 {
			return North
		}
		return South
	}
	diff := i1 - i2
	if diff < 0 {
		diff *= -1
	}
	if i1-i2 < width {
		if i1 > i2 {
			return West
		}
		return East
	}
	return BadDir
}

// Returns a map of direction bitflags to the indices in that direction from
// the input directions.
func flagToIndicies(tm *TextMap, index int, dirFlag Direction) map[Direction]int {
	output := make(map[Direction]int)
	if n := index - tm.width; dirFlag&North == North {
		output[North] = n
	}
	if e := index + 1; dirFlag&East == East {
		output[East] = e
	}
	if s := index + tm.width; dirFlag&South == South {
		output[South] = s
	}
	if w := index - 1; dirFlag&West == West {
		output[West] = w
	}
	return output
}

func flagToIndex(tm *TextMap, index int, dir Direction) (int, error) {
	if output, ok := flagToIndicies(tm, index, dir)[dir]; ok {
		return output, nil
	} else {
		return 0, fmt.Errorf("No index in that direction.")
	}
}
