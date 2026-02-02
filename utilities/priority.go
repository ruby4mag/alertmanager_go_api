package utilities

import (
	"strings"
	"strconv"
)

// PriorityToInt converts a priority string (e.g., "P1", "P2") to an integer (1, 2).
// Lower number means higher priority.
// Defaults to 5 (lowest) if invalid.
func PriorityToInt(p string) int {
	p = strings.ToUpper(strings.TrimSpace(p))
	if strings.HasPrefix(p, "P") {
		val, err := strconv.Atoi(p[1:])
		if err == nil {
			return val
		}
	}
	// Handle cases where might be just "0", "1", "2" etc
    if val, err := strconv.Atoi(p); err == nil {
        return val
    }
	// Default to lowest if unknown. Assuming P4 is lowest for now based on user request "P0 to P4".
	return 4
}

// IntToPriority converts an integer to a priority string (e.g., 1 -> "P1").
func IntToPriority(i int) string {
	return "P" + strconv.Itoa(i)
}
