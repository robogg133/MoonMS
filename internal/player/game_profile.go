package player

import "github.com/google/uuid"

type GameProfile struct {
	UUID        uuid.UUID
	Username    string
	Proprieties struct {
		Name      string
		Value     string
		Signature string
	}
}
