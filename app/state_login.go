package app

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/robogg133/KernelCraft/internal/offline"

	"github.com/robogg133/KernelCraft/internal/packets"

	"github.com/robogg133/KernelCraft/internal/datatypes"

	"github.com/Tnze/go-mc/net/CFB8"
	"github.com/google/uuid"
)

const STATE_NAME_LOGIN = "handshake"
const MOJANG_SESSION_CHECKER string = "https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s"

var ErrInvalidSession = errors.New("invalid session")

type LoginState struct{}

type MojangAnswer struct {
	Properties []map[string]string
}

func (s *LoginState) Name() string { return STATE_NAME_LOGIN }

func (s *LoginState) Handle(sess *Session) error {
	sess.Server.LogDebug("START STATE LOGIN")

	// Registering Encryption response
	new := make(packets.KnownPackets)
	new.RegisterPacket(packets.PACKET_ENCRYPTION_RESPONSE, func() packets.Packet {
		return &packets.EncryptionResponse{}
	})
	new.RegisterPacket(packets.PACKET_LOGIN_ACKNOWLEDGED, func() packets.Packet {
		return &packets.LoginAcknowledgedPacket{}
	})
	sess.KnownPkgs = new
	// Starting to read basic info like uuid and username
	buff, err := packets.ReadPackageFromConnecion(sess.Conn)
	if err != nil {
		return err
	}

	_, offset, err := datatypes.ParseVarInt(bytes.NewReader(buff))
	if err != nil {
		return err
	}

	protocolID, tmp, err := datatypes.ParseVarInt(bytes.NewReader(buff[offset:]))
	if err != nil {
		return err
	}
	offset += tmp
	tmp = 0

	if protocolID != packets.PACKET_HANDSHAKE {
		return err
	}

	stringLenght, tmp, err := datatypes.ParseVarInt(bytes.NewReader(buff[offset:]))
	if err != nil {
		return err
	}

	offset += tmp
	tmp = 0

	playerName := buff[offset : offset+int(stringLenght)]
	playerUUID, err := uuid.FromBytes(buff[offset+int(stringLenght):])
	if err != nil {
		return err
	}

	if !sess.Server.MinecraftConfig.Proprieties.OnlineMode {
		playerUUID = offline.NameToUUID(string(playerName))
	}

	sess.Server.LogInfo(fmt.Sprintf("Received connection from %s as %s (%s)", sess.Conn.RemoteAddr().String(), string(playerName), playerUUID.String()))
	var usernameCode string
	if sess.Server.MinecraftConfig.Proprieties.OnlineMode || sess.Server.MinecraftConfig.Advanced.OfflineEncryption {

		sess.EncryptCipher, sess.DecryptCipher, usernameCode, err = setupEncryption(sess)
		sess.Server.LogDebug("Setup encryption")
		sess.PkgReader = packets.NewReaderFromReader(packets.NewCipherReader(sess.Conn, sess.DecryptCipher))
	}

	if sess.Server.MinecraftConfig.Advanced.Threshold > -1 {
		compressStart := packets.CompressionStart{Threshould: sess.Server.MinecraftConfig.Advanced.Threshold}

		if err := sess.WritePacket(&compressStart); err != nil {
			return err
		}
		sess.Server.LogDebug("Sent compress Start")
		sess.Threshold = sess.Server.MinecraftConfig.Advanced.Threshold
	}
	var resp MojangAnswer

	if sess.Server.MinecraftConfig.Proprieties.OnlineMode {
		resp, err = mojangCheck(string(playerName), usernameCode)
		if err != nil {
			if err == ErrInvalidSession {
				// Need to kick msg
				return ErrNoReason
			} else {
				sess.Server.LogDebug("what happned here?")
				return err
			}
		}
	} else {
		resp.Properties = make([]map[string]string, 0)
	}

	b, err := playerUUID.MarshalBinary()
	if err != nil {
		panic(err)
	}

	var name string
	var value string
	var signature string

	for _, v := range resp.Properties {
		if _, exists := v["name"]; exists {
			name = v["name"]
		}
		if _, exists := v["value"]; exists {
			value = v["value"]
		}
		if _, exists := v["signature"]; exists {
			signature = v["signature"]
		}
	}

	loginSucc := packets.LoginSuccessPacket{
		Profile: packets.GameProfile{
			UUID:     b,
			Username: string(playerName),
			Name:     name,
			Value:    value,
		},
	}

	if signature != "" {
		loginSucc.Profile.HaveSignature = true
		loginSucc.Profile.Signature = signature
	}

	if err = sess.WritePacket(&loginSucc); err != nil {
		return err
	}
	sess.Server.LogDebug("Sent Login success")

	ak, err := sess.ReadPacket()
	if err != nil {
		return err
	}

	if ak.ID() != packets.PACKET_LOGIN_ACKNOWLEDGED {
		sess.Server.LogDebug("Protocol violation")
		return ErrNoReason
	}
	sess.Server.LogDebug("Login acknowledged")

	sess.State = &ConfigState{}

	return err
}

func setupEncryption(sess *Session) (encKey, decKey *CFB8.CFB8, usernameHex string, err error) {

	var serverid string = ""

	verifyToken := make([]byte, 4)
	_, err = io.ReadFull(rand.Reader, verifyToken)
	if err != nil {
		return nil, nil, "", err
	}

	pk, err := x509.MarshalPKIXPublicKey(sess.Server.ServerPrivateKey.Public())
	if err != nil {
		return nil, nil, "", err
	}

	var response packets.EncryptionRequest
	response.ServerID = serverid
	response.PublicKey = pk
	response.VerifyToken = verifyToken
	response.ShouldAuth = true

	if err := sess.WritePacket(&response); err != nil {
		return nil, nil, "", err
	}
	sess.Server.LogDebug("Sent encryption request")

	pkg, err := sess.ReadPacket()
	if err != nil {
		sess.Server.LogDebug("err here1")
		return nil, nil, "", err
	}

	encR := pkg.(*packets.EncryptionResponse)

	plainToken, err := sess.Server.ServerPrivateKey.Decrypt(nil, encR.VerifyTokenCiphered, nil)

	if !bytes.Equal(plainToken, verifyToken) {
		sess.Server.LogWarn("Invalid Session")
		sess.Conn.Close()
		return nil, nil, "", nil
	}

	sharedSecret, err := sess.Server.ServerPrivateKey.Decrypt(nil, encR.SharedSecretCiphered, nil)
	if err != nil {
		return nil, nil, "", err
	}

	hash := sha1.New()
	hash.Write([]byte(serverid))
	hash.Write(sharedSecret)
	hash.Write(pk)

	sum := hash.Sum(nil)

	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return nil, nil, "", err
	}

	return CFB8.NewCFB8Encrypt(block, sharedSecret), CFB8.NewCFB8Decrypt(block, sharedSecret), hex.EncodeToString(sum), nil
}

func mojangCheck(playerName, userCode string) (MojangAnswer, error) {
	resp, err := http.Get(fmt.Sprintf(MOJANG_SESSION_CHECKER, string(playerName), userCode))
	if err != nil {
		return MojangAnswer{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return MojangAnswer{}, ErrInvalidSession
	}

	var mojangResponse MojangAnswer

	mojangPayload, err := io.ReadAll(resp.Body)
	if err != nil {
		return MojangAnswer{}, err
	}
	if err := json.Unmarshal(mojangPayload, &mojangResponse); err != nil {
		return MojangAnswer{}, err
	}

	return mojangResponse, nil
}
