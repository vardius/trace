package trace

import (
	"bytes"
	"log"
)

// Here returns string representation of a reference
// to a File:Line calling this function
func Here() string {
	var buf bytes.Buffer
	var logger = log.New(&buf, "", log.Llongfile)

	logger.Output(2, "")

	return buf.String()
}
