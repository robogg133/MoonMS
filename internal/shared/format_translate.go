package shared

import "encoding/json"

func TranslationMessage(s string) json.RawMessage {
	return []byte(`{"translate":"` + s + `"}`)
}
