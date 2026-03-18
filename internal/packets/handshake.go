package packets

import (
	"encoding/json"
)

const PACKET_HELLO int32 = 0x00
const PACKET_STATUS int32 = 0x00

type HelloPacket struct {
	ProtocolVersion int32
	ServerAdress    string
	ServerPort      uint16
	Intent          int32
}

func (h *HelloPacket) ID() int32 { return PACKET_HELLO }

func (h *HelloPacket) Encode(w *Writer) error {

	w.WriteVarInt(h.ProtocolVersion)

	w.WriteString(h.ServerAdress)

	w.WriteUnsignedShort(h.ServerPort)

	w.WriteVarInt(h.Intent)

	return nil
}

func (h *HelloPacket) Decode(r *Reader) error {

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

type PlayerListInfo struct {
	Username string `json:"name"`
	UUID     string `json:"id"`
}

type StatusPacket struct {
	Version struct {
		Name            string `json:"name"`
		ProtocolVersion int32  `json:"protocol"`
	} `json:"version"`

	Players struct {
		MaxPlayers    uint32           `json:"max"`
		OnlinePlayers uint32           `json:"online"`
		PlayerStatus  []PlayerListInfo `json:"sample"`
	} `json:"players"`

	Description struct {
		Text string `json:"text"`
	} `json:"description"`

	Favicon           string `json:"favicon"`
	EnforceSecureChat bool   `json:"enforcesSecureChat"`
}

func (h *StatusPacket) ID() int32 { return PACKET_STATUS }

func (h *StatusPacket) Encode(w *Writer) error {

	statusSerialized, err := json.Marshal(&h)
	if err != nil {
		return err
	}

	return w.WritePrefixed(statusSerialized)
}

func (h *StatusPacket) Decode(r *Reader) error {

	s, err := r.ReadPrefixed()
	if err != nil {
		return err
	}

	err = json.Unmarshal(s, h)

	return err
}
