package main

import (
	"fmt"

	"github.com/YumeMichi/yohane/events"
)

func main() {
	ch := make(chan string)
	go events.Round(ch)

	for {
		out := <-ch
		fmt.Println(out)
	}
}
