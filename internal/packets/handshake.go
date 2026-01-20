package packets

import "MoonMS/internal/datatypes"

type Handshake struct {
	ProtocolVersion uint16
	ServerAdress    string
	ServerPort      datatypes.UnsignedShort
	Intent          uint8
}

type PlayerMinimunInfo struct {
	Username string `json:"name"`
	UUID     string `json:"id"`
}

type HandShakeResponseStatus struct {
	Version struct {
		Name            string `json:"name"`
		ProtocolVersion uint16 `json:"protocol"`
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
