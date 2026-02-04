//Needs tinygo compiler
//go:build wasip1

package moonms_api

import "fmt"

type events struct {
	ServerStopping []func()
}

var allEvents events

func RegisterOnServerStop(fn func()) {
	allEvents.ServerStopping = append(allEvents.ServerStopping, fn)
}

//export on_server_stop
func on_server_stop() {
	fmt.Println("got called")
	for _, fn := range allEvents.ServerStopping {
		go fn()
	}
}
