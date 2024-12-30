package tool

import "fmt"

func HumanDuration(nano int64) string {
	duration := float64(nano)
	unit := "ns"
	if duration >= 1000 {
		duration /= 1000
		unit = "us"
	}

	if duration >= 1000 {
		duration /= 1000
		unit = "ms"
	}

	if duration >= 1000 {
		duration /= 1000
		unit = " s"
	}

	return fmt.Sprintf("%6.2f%s", duration, unit)
}
