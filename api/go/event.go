package moonms

type TriggerID uint8
type ServerEvent func() any

var serverEvents map[TriggerID]ServerEvent = make(map[TriggerID]ServerEvent)

//go:wasmexport server_stopping_event
func server_stopping_event() {
	fn, exists := serverEvents[id_SERVER_STOP_EVENT]
	if !exists {
		return
	}
	fn()
}

func AssignServerStopEvent(fn func()) {
	serverEvents[id_SERVER_STOP_EVENT] = adaptNoReturn(fn)
}
