/*
Package trace provides simple helpers to trace the function calls, errors or logs `File:Line Function` reference

	package main

	import (
		"log"

		"github.com/vardius/trace"
	)

	func main() {
		log.Printf("%s\n\t%s", "Hello from:", trace.Here())
	}
*/
package trace
