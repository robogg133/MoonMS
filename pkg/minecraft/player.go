package minecraft

import (
	"MoonMS/internal/packets"

	"github.com/google/uuid"
)

type Language uint8

const (
	EN_US Language = iota
	PT_BR
	EN_UK
)

type PlayerVisualFlags byte

const (
	N_RIGHT_HAND uint8 = iota
	N_CAPE
	N_LEFT_SLEEVE
	N_LEFT_PANT_LEG
	N_HAT
	N_JACKET
	N_RIGHT_SLEEVE
	N_RIGHT_PANT_LEG
)

type PlayerInfo struct {
	Language    Language
	ChatEnabled bool
	ChatColors  bool

	VisualFlags PlayerVisualFlags

	TextFilterEnabled bool
	AllowServerList   bool

	ParticleStatus uint8
}

type Player struct {
	Username       string
	UUID           uuid.UUID
	isOnlinePlayer bool

	GameProfile packets.GameProfile

	PlayerData string

	ConnectionAdress string
}

func (s *PlayerVisualFlags) CheckIsTrue(n *uint8) bool {
	b := byte(*s)

	if b&(1<<*n) != 0 {
		return true
	}

	return false
}

func (s *PlayerVisualFlags) SetValue(n *uint8, value bool) {
	if value {
		*s |= 1 << *n
	} else {
		*s &^= 1 << *n
	}
}
