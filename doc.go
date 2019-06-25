/*
Package trace provides simple helpers to trace the function calls, errors or logs reference

	package main

	import (
		"log"

		"github.com/vardius/trace"
	)

	func main() {
		log.Printf("%s\n\t%s", "Hello from:", trace.Here(trace.Lfile | trace.Lline | trace.Lfunction))
	}
*/
package trace
