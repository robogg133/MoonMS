package packets

const PACKET_CLIENT_INFORMATION int32 = 0x00

const (
	CHAT_ENUM_MODE_ON       int32 = 0
	CHAT_ENUM_MODE_ONLY_CMD int32 = 1
	CHAT_ENUM_MODE_OFF      int32 = 2
)

type PlayerSkinFlags byte

const (
	SKIN_ENUM_CAPE uint8 = iota
	SKIN_ENUM_JACKET
	SKIN_ENUM_LEFT_SLEEVE
	SKIN_ENUM_RIGHT_SLEEVE
	SKIN_ENUM_LEFT_PANT_LEG
	SKIN_ENUM_RIGHT_PANT_LEG
	SKIN_ENUM_HAT
)

const (
	PARTICLE_ENUM_ALL     int32 = 0
	PARTICLE_ENUM_REDUCED int32 = 1
	PARTICLE_ENUM_MINIMUN int32 = 2
)

type ClientInformationPacket struct {
	Locale       string
	ViewDistance uint8

	ChatMode   int32
	ChatColors bool

	SkinProprieties PlayerSkinFlags

	IsRightHand bool

	EnableTextFiltering bool

	AllowServerList bool

	ParticleStatus int32
}

func (c *ClientInformationPacket) ID() int32 { return PACKET_CLIENT_INFORMATION }

func (c *ClientInformationPacket) Encode(w *Writer) error { return nil }

func (c *ClientInformationPacket) Decode(r *Reader) error {
	var err error
	c.Locale, err = r.ReadString()
	if err != nil {
		return err
	}

	c.ViewDistance, err = r.ReadByte()
	if err != nil {
		return err
	}

	c.ChatMode, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	c.ChatColors, err = r.ReadBoolean()
	if err != nil {
		return err
	}

	sb, err := r.ReadByte()
	if err != nil {
		return err
	}
	c.SkinProprieties = PlayerSkinFlags(sb)

	c.IsRightHand, err = r.ReadBoolean()
	if err != nil {
		return err
	}

	c.EnableTextFiltering, err = r.ReadBoolean()
	if err != nil {
		return err
	}

	c.AllowServerList, err = r.ReadBoolean()
	if err != nil {
		return err
	}

	c.ParticleStatus, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	return nil
}

func (s *PlayerSkinFlags) CheckIsTrue(n *uint8) bool {
	b := byte(*s)

	if b&(1<<*n) != 0 {
		return true
	}

	return false
}

func (s *PlayerSkinFlags) SetValue(n *uint8, value bool) {
	if value {
		*s |= 1 << *n
	} else {
		*s &^= 1 << *n
	}
}
