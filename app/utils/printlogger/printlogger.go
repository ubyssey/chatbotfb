package printlogger

import (
	"fmt"
	"time"
)

// Generic logging with time prefix
func PrintLogger(format string, args ...interface{}) {
	// Include stack traces maybe? See errgo package
	fmt.Printf("[LOG] "+time.Now().Format("2017-05-27 00:00:00")+" "+format, args...)
}
