package util

import "fmt"

func FormatBytes(size uint64, sp bool) string {
	//Determine the correct space character
	ws := ""
	if sp {
		ws = " "
	}

	//Get the size as a float64
	value := float64(size)

	//Determine the correct units
	unit := ""
	switch {
	case size >= 1<<40:
		unit = "TB"
		value /= 1 << 40
	case size >= 1<<30:
		unit = "GB"
		value /= 1 << 30
	case size >= 1<<20:
		unit = "MB"
		value /= 1 << 20
	case size >= 1<<10:
		unit = "KB"
		value /= 1 << 10
	default:
		return fmt.Sprintf("%d%sB", size, ws)
	}

	//Sprintf the output
	return fmt.Sprintf("%.2f%s%s", value, ws, unit)
}

func InRange(num int64, min int64, max int64) bool {
	return num >= min && num <= max
}

func EqualsAny[T comparable](targ T, items ...T) bool {
	for _, item := range items {
		if targ == item {
			return true
		}
	}
	return false
}
