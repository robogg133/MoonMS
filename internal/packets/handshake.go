package packets

import (
	"encoding/json"
)

const PACKET_HANDSHAKE int32 = 0x00

type Handshake struct {
	ProtocolVersion int32
	ServerAdress    string
	ServerPort      uint16
	Intent          int32
}

func (h *Handshake) ID() int32 { return PACKET_HANDSHAKE }

func (h *Handshake) Encode(w *Writer) error {
	return nil
}

func (h *Handshake) Decode(r *Reader) error {

	n, err := r.ReadVarInt()
	if err != nil {
		return err
	}
	h.ProtocolVersion = n

	h.ServerAdress, err = r.ReadString()
	if err != nil {
		return err
	}

	h.ServerPort = r.ReadUnsignedShort()

	h.Intent, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	return nil
}

type PlayerMinimunInfo struct {
	Username string `json:"name"`
	UUID     string `json:"id"`
}

type HandShakeResponseStatus struct {
	Version struct {
		Name            string `json:"name"`
		ProtocolVersion int32  `json:"protocol"`
	} `json:"version"`

	Players struct {
		MaxPlayers    uint                `json:"max"`
		OnlinePlayers uint32              `json:"online"`
		PlayerStatus  []PlayerMinimunInfo `json:"sample"`
	} `json:"players"`

	Description struct {
		Text string `json:"text"`
	} `json:"description"`

	Favicon           string `json:"favicon"`
	EnforceSecureChat bool   `json:"enforcesSecureChat"`
}

func (h *HandShakeResponseStatus) ID() int32 { return PACKET_HANDSHAKE }

func (h *HandShakeResponseStatus) Encode(w *Writer) error {

	statusSerialized, err := json.Marshal(&h)
	if err != nil {
		return err
	}

	return w.WritePrefixed(statusSerialized)
}

func (h *HandShakeResponseStatus) Decode(r *Reader) error { return nil }
