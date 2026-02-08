package packets

const PACKET_ENCRYPTION_REQUEST int32 = 0x01
const PACKET_ENCRYPTION_RESPONSE int32 = 0x01

type EncryptionRequest struct {
	ServerID    string
	PublicKey   []byte
	VerifyToken []byte
	ShouldAuth  bool
}

func (e *EncryptionRequest) ID() int32 { return PACKET_ENCRYPTION_REQUEST }

func (e *EncryptionRequest) Encode(w *Writer) error {

	if err := w.WriteString(e.ServerID); err != nil {
		return err
	}

	if err := w.WritePrefixed(e.PublicKey); err != nil {
		return err
	}
	if err := w.WritePrefixed(e.VerifyToken); err != nil {
		return err
	}
	if err := w.WriteBoolean(e.ShouldAuth); err != nil {
		return err
	}

	return nil
}

func (e *EncryptionRequest) Decode(r *Reader) error { return nil }

type EncryptionResponse struct {
	SharedSecretCiphered []byte
	VerifyTokenCiphered  []byte
}

func (e *EncryptionResponse) ID() int32              { return PACKET_ENCRYPTION_RESPONSE }
func (e *EncryptionResponse) Encode(r *Writer) error { return nil }

func (e *EncryptionResponse) Decode(r *Reader) error {
	var err error
	e.SharedSecretCiphered, err = r.ReadPrefixed()
	if err != nil {
		return err
	}

	e.VerifyTokenCiphered, err = r.ReadPrefixed()
	if err != nil {
		return err
	}

	return nil
}
