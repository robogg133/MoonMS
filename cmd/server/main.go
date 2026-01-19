package main

import (
	"MoonMS/cmd/server/compress"
	"MoonMS/cmd/server/config"
	"MoonMS/cmd/server/crypto"
	"MoonMS/internal/datatypes"
	"MoonMS/internal/offline"
	"MoonMS/internal/packets"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "embed"

	"github.com/google/uuid"
)

var PROTOCOL_VERSION uint16 = 774
var CURRENT_VERSION string = "1.21.11"

const MOJANG_SESSION_CHECKER = `https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s`
const DEADLINE = time.Second * 30

type MojangAnswer struct {
	Properties []map[string]string
}

func logInfo(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[34m[INFO]:\033[0m %v\n", timestamp, message)
}

func logError(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[31m[ERROR]:\033[0m %v\n", timestamp, message)
}

func logFatal(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[31m[FATAL]:\033[0m %v\n", timestamp, message)
	os.Exit(1)
}

var CFG config.Configs

var AnonymousPlayer = &packets.PlayerMinimunInfo{Username: "Anonymous Player", UUID: "00000000-0000-0000-0000-000000000000"}

func CheckFilesToStart() error {
	_, err := os.Stat("banned-ips.txt")
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create("banned-ips.txt")
			if err != nil {

				return err
			}
			f.Close()

		}
	}
	_, err = os.Stat("banned-accounts.txt")
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create("banned-accounts.txt")
			if err != nil {
				return err
			}
			f.Close()

		}
	}

	if _, err = os.Stat("server-private-key.pem"); err != nil {
		if os.IsNotExist(err) {
			logInfo("Generating server key pair")
			if err := crypto.GenerateKeyPair(CFG.Proprieties.RSAKeyBits); err != nil {
				logFatal(err)
			}
		}
	} else {
		if _, err = os.Stat("server-public.key"); err != nil {
			if os.IsNotExist(err) {
				logError("Public key don't exist, but private exist, generating another public key from the private key")
				pemKey, err := os.ReadFile("server-private-key.pem")
				if err != nil {
					logFatal(err)
				}
				p, _ := pem.Decode(pemKey)

				privateKey, err := x509.ParsePKCS1PrivateKey(p.Bytes)
				if err != nil {
					logFatal(err)
				}
				publicKey := privateKey.Public()

				publicKeyMarshalized, err := x509.MarshalPKIXPublicKey(publicKey)
				if err != nil {
					logFatal(err)
				}

				if err := os.WriteFile("server-public.key", publicKeyMarshalized, 0644); err != nil {
					log.Fatal(err)
				}
			}
		}

	}

	return nil
}

func main() {

	logInfo("Reading main configuration file")
	var err error
	CFG, err = config.ReadConfigurationFile()
	if err != nil {
		logFatal(err)
	}

	if err := CheckFilesToStart(); err != nil {
		logError(err)
		return
	}

	logInfo(fmt.Sprintf("Starting minecraft java edittion tcp listener on port %d", CFG.Proprieties.ServerPort))
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", CFG.Proprieties.ServerPort))
	if err != nil {
		logFatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var sharedSecret []byte

	buff, err := ReadPackageFromConnecion(conn)
	if err != nil {
		logError(err)
		return
	}

	_, err = conn.Read(buff)
	if err != nil {
		logError(err)
		return
	}

	n, readed, err := datatypes.ParseVarInt(bytes.NewBuffer(buff))
	if err != nil {
		logError(err)
	}
	if uint16(n) != PROTOCOL_VERSION {
		logError("Mismatched version from the client")
		return
	}
	stringLenght := uint8(buff[readed])
	stringOffset := stringLenght + uint8(readed) + 1
	// serverAdress := string(buf[readed+1 : stringOffset])

	// port := binary.BigEndian.Uint16(buf[stringOffset : stringOffset+2])

	intention, _, err := datatypes.ParseVarInt(bytes.NewBuffer(buff[stringOffset+2:]))

	switch intention {
	case 1:
		statuspkg, err := ReadPackageFromConnecion(conn)
		if err != nil {
			logError(err)
			return
		}

		if statuspkg[1] != packets.PACKET_HANDSHAKE {
			logError("mismatched packet id gaved by client")
			return
		}
		statuspkg = nil

		var status packets.HandShakeResponseStatus

		status.Version.Name = CURRENT_VERSION
		status.Version.ProtocolVersion = uint16(PROTOCOL_VERSION)

		status.Players.MaxPlayers = CFG.Proprieties.MaxPlayers
		status.Players.OnlinePlayers = 0
		status.Players.PlayerStatus = []packets.PlayerMinimunInfo{}

		if CFG.Proprieties.ServerIcon == "" {
			status.Favicon = ""
		} else {
			status.Favicon, err = GetBase64Image(CFG.Proprieties.ServerIcon)
			if err != nil {
				logError(err)
				status.Favicon = ""
			} else {
				status.Favicon = fmt.Sprintf("data:image/png;base64,%s", status.Favicon)
			}
		}

		status.Description.Text = CFG.Proprieties.Motd

		statusSerialized, err := json.Marshal(&status)
		if err != nil {
			logError(err)
			return
		}
		lenghtForResponse := len(statusSerialized)
		lenghtForResponsePayload := datatypes.NewVarInt(int32(len(statusSerialized)))
		packageID := datatypes.NewVarInt(packets.PACKET_HANDSHAKE)

		totalLenght := len(packageID) + lenghtForResponse + len(lenghtForResponsePayload)

		packageLenght := datatypes.NewVarInt(int32(totalLenght))

		var response bytes.Buffer

		response.Write(packageLenght)
		response.Write(packageID)
		response.Write(lenghtForResponsePayload)
		response.Write(statusSerialized)

		conn.Write(response.Bytes())

		pingPkg := make([]byte, 10)

		conn.Read(pingPkg)

		offset := 0
		_, n, err := datatypes.ParseVarInt(bytes.NewReader(pingPkg))
		if err != nil {
			logError(err)
			return
		}

		offset += n

		pkgID, n, err := datatypes.ParseVarInt(bytes.NewReader(pingPkg[n:]))
		if err != nil {
			logError(err)
			return
		}

		offset += n
		n = 0

		if pkgID == packets.PACKET_PING {
			conn.Write(packets.SerializePong(pingPkg[offset:]))
		}

		return

	case 2:

		Deadline := time.Now().Add(DEADLINE)
		conn.SetReadDeadline(Deadline)

		buf, err := ReadPackageFromConnecion(conn)
		if err != nil {
			logError(err)
			return
		}

		_, offset, err := datatypes.ParseVarInt(bytes.NewReader(buf))
		if err != nil {
			logError(err)
			return
		}

		protocolID, tmp, err := datatypes.ParseVarInt(bytes.NewReader(buf[offset:]))
		if err != nil {
			logError(err)
			return
		}
		offset += tmp
		tmp = 0

		if protocolID != packets.PACKET_HANDSHAKE {
			logError("Invalid package received from player")
			return
		}

		stringLenght, tmp, err := datatypes.ParseVarInt(bytes.NewReader(buf[offset:]))
		if err != nil {
			logError(err)
			return
		}

		offset += tmp
		tmp = 0

		playerName := buf[offset : offset+int(stringLenght)]
		clientAddr := conn.RemoteAddr().String()

		if CFG.Proprieties.OnlineMode {
			playerUUID, err := uuid.FromBytes(buf[offset+int(stringLenght):])
			if err != nil {
				logError(err)
				return
			}
			logInfo(fmt.Sprintf("Received connection from %s as %s (%s)", clientAddr, string(playerName), playerUUID.String()))

			packetID := datatypes.NewVarInt(int32(packets.PACKET_ENCRYPTION_REQUEST))

			serverID := []byte("")
			serverIDPrefix := datatypes.NewVarInt(int32(len(serverID)))

			publicKeyMarshal, err := os.ReadFile("server-public.key")
			if err != nil {
				logError(err)
				return
			}
			publicKeyMarshalPrefix := datatypes.NewVarInt(int32(len(publicKeyMarshal)))

			verifyToken := make([]byte, 4)
			rand.Read(verifyToken)
			verifyTokenPrefix := datatypes.NewVarInt(4)

			shouldAuth := datatypes.NewBoolean(true)

			totalLenght := len(packetID) + len(serverIDPrefix) + len(serverID) + len(publicKeyMarshalPrefix) + len(publicKeyMarshal) + len(verifyTokenPrefix) + len(verifyToken) + 1 // +1 for shouldAuth
			lenght := datatypes.NewVarInt(int32(totalLenght))

			var responseBuffer bytes.Buffer

			responseBuffer.Write(lenght)
			responseBuffer.Write(packetID)

			responseBuffer.Write(serverIDPrefix)
			responseBuffer.Write(serverID)

			responseBuffer.Write(publicKeyMarshalPrefix)
			responseBuffer.Write(publicKeyMarshal)

			responseBuffer.Write(verifyTokenPrefix)
			responseBuffer.Write(verifyToken)

			responseBuffer.WriteByte(byte(shouldAuth))

			//Sending encryption request
			conn.Write(responseBuffer.Bytes())
			responseBuffer.Reset()

			buf, err := ReadPackageFromConnecion(conn)
			if err != nil {
				logError(err)
				return
			}

			compressionStartPkg := &packets.CompressionStart{
				Threshould: CFG.Proprieties.ServerThreshold,
			}

			conn.Write(compressionStartPkg.Serialize())

			_, offset, err = datatypes.ParseVarInt(bytes.NewReader(buf))
			if err != nil {
				logError(err)
				return
			}

			protocolID, tmp, err := datatypes.ParseVarInt(bytes.NewReader(buf[offset:]))
			if err != nil {
				logError(err)
				return
			}
			offset += tmp

			if protocolID != packets.PACKET_ENCRYPTION_RESPONSE {
				logError("Player failed to answer encryption response")
				return
			}

			sharedSecretLenght, tmp, err := datatypes.ParseVarInt(bytes.NewBuffer(buf[offset:]))
			if err != nil {
				logError(err)
				return
			}
			offset += tmp
			tmp = 0

			sharedSecretCipher := make([]byte, sharedSecretLenght)

			copy(sharedSecretCipher, buf[offset:offset+int(sharedSecretLenght)])

			verifyTokenLenght, tmp, err := datatypes.ParseVarInt(bytes.NewBuffer(buf[offset+int(sharedSecretLenght):]))
			if err != nil {
				logError(err)
				return
			}
			offset += int(sharedSecretLenght) + tmp

			verifyTokenClientCiphered := make([]byte, verifyTokenLenght)
			copy(verifyTokenClientCiphered, buf[offset:offset+int(verifyTokenLenght)])

			privKeyPem, err := os.ReadFile("server-private-key.pem")
			if err != nil {
				logError(err)
				return
			}

			p, _ := pem.Decode(privKeyPem)

			privateKey, err := x509.ParsePKCS1PrivateKey(p.Bytes)
			if err != nil {
				logError(fmt.Sprintf("error parsing private key: %v", err))
				return
			}

			plainVerifyToken, err := privateKey.Decrypt(nil, verifyTokenClientCiphered, nil)
			if err != nil {
				logError(err)
				return
			}

			if !bytes.Equal(plainVerifyToken, verifyToken) {
				payload, err := packets.DisconnectLogin("Invalid Session")
				if err != nil {
					logError(err)
					return
				}
				conn.Write(payload)
				return
			}

			sharedSecret, err = privateKey.Decrypt(nil, sharedSecretCipher, nil)
			if err != nil {
				logError(err)
				return
			}

			hash := sha1.New()
			hash.Write(serverID)
			hash.Write(sharedSecret)
			hash.Write(publicKeyMarshal)

			sum := hash.Sum(nil)

			var resp *http.Response
			for i := range 10 {
				resp, err = http.Get(fmt.Sprintf(MOJANG_SESSION_CHECKER, string(playerName), hex.EncodeToString(sum)))
				if err != nil {
					logError(err)
					return
				}
				if resp.StatusCode != 200 && i == 9 {
					defer resp.Body.Close()
					payload, err := packets.DisconnectLogin("Invalid session")
					if err != nil {
						logError(err)
						return
					}
					conn.Write(payload)
				}
				if resp.StatusCode != 200 {
					resp.Body.Close()
				} else {
					break
				}
			}
			defer resp.Body.Close()

			mojangPayload, err := io.ReadAll(resp.Body)
			if err != nil {
				logError(err)
				return
			}

			var mojangResponse MojangAnswer
			if err := json.Unmarshal(mojangPayload, &mojangResponse); err != nil {
				logError(err)
				return
			}

			var name string
			var value string
			var signature string

			for _, v := range mojangResponse.Properties {
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

			var loginSuccesspkg packets.LoginSuccessPacket
			var namebuff bytes.Buffer

			namebuff.Write(datatypes.NewVarInt(int32(len(playerName))))
			namebuff.Write(playerName)
			loginSuccesspkg.Profile.Username = datatypes.String(namebuff.Bytes())
			namebuff.Reset()
			loginSuccesspkg.Profile.UUID, err = playerUUID.MarshalBinary()
			if err != nil {
				logError(err)
				return
			}
			loginSuccesspkg.Profile.Name = name
			loginSuccesspkg.Profile.Value = value
			loginSuccesspkg.Profile.HaveSignature = true
			loginSuccesspkg.Profile.Signature = signature
			block, err := aes.NewCipher(sharedSecret)
			if err != nil {
				logError(err)
				return
			}

			encrypter := cipher.NewCTR(block, sharedSecret)
			tmpBuff, err := compress.Compress(loginSuccesspkg.Serialize(), CFG.Proprieties.ServerThreshold)
			if err != nil {
				logError(err)
				return
			}

			responseCipher := make([]byte, len(tmpBuff))
			encrypter.XORKeyStream(responseCipher, tmpBuff)
			conn.Write(responseCipher)
			tmpBuff = nil
			fmt.Println("sent everything")

			aknowledge := make([]byte, 3)
			_, err = conn.Read(aknowledge)
			if err != nil {
				return
			}
			fmt.Println(aknowledge)

		} else {

			compressionStartPkg := &packets.CompressionStart{
				Threshould: CFG.Proprieties.ServerThreshold,
			}

			conn.Write(compressionStartPkg.Serialize())

			playerUUID := offline.NameToUUID(string(playerName))

			logInfo(fmt.Sprintf("Received connection from %s as %s (%s) [OFFLINE]", clientAddr, string(playerName), playerUUID.String()))
			var loginSuccesspkg packets.LoginSuccessPacket
			var namebuff bytes.Buffer

			namebuff.Write(datatypes.NewVarInt(int32(len(playerName))))
			namebuff.Write(playerName)
			loginSuccesspkg.Profile.Username = datatypes.String(namebuff.Bytes())
			namebuff.Reset()
			loginSuccesspkg.Profile.UUID, err = playerUUID.MarshalBinary()
			if err != nil {
				logError(err)
				return
			}

			response := loginSuccesspkg.Serialize()

			compressedResponse, err := compress.Compress(response, CFG.Proprieties.ServerThreshold)
			if err != nil {
				logError(err)
				return
			}

			conn.Write(compressedResponse)

			aknowledge := make([]byte, 3)
			_, err = conn.Read(aknowledge)
			if err != nil {
				return
			}

			if aknowledge[2] != byte(packets.PACKET_LOGIN_AKNOWLEDGED) {
				logError("Invalid packet")
				return
			}
		}

	}

	conn.SetReadDeadline(time.Time{})
	for {
		time.Sleep(5 * time.Second)
	}
}

func isConnAlive(conn net.Conn, makeDeadline time.Time) bool {
	conn.SetReadDeadline(time.Now())
	defer conn.SetReadDeadline(makeDeadline)

	zero := make([]byte, 0)
	_, err := conn.Read(zero)

	if err == nil {
		return true
	}

	return false

}

func GetBase64Image(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("Can't find specified server image using none")
		}
		return "", err
	}

	if len(content) > 5120 {
		return "", fmt.Errorf("The file size is too big!! needs to be lower than 5KB, using none")
	}

	return base64.StdEncoding.EncodeToString(content), nil
}
