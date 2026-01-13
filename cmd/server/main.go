package main

import (
	"MoonMS/cmd/server/config"
	"MoonMS/cmd/server/crypto"
	datatypes "MoonMS/internal/datatypes/varTypes"
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

const TRUE_VALUE byte = 0x01
const FALSE_VALUE byte = 0x00

const MOJANG_SESSION_CHECKER = `https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s`

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

func main() {

	logInfo("Reading main configuration file")
	var err error
	CFG, err = config.ReadConfigurationFile()
	if err != nil {
		logFatal(err)
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

	logInfo(fmt.Sprintf("Starting java tcp listener on port %d", CFG.Proprieties.ServerPort))
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

	startingData := make([]byte, 2)
	_, err := conn.Read(startingData)
	if err != nil {
		logError(err)
		return
	}

	if startingData[1] != packets.PACKET_HANDSHAKE {
		logError("The player don't started the connection with a handshake")
		return
	}

	buf := make([]byte, uint8(startingData[0])-1)

	_, err = conn.Read(buf)
	if err != nil {
		logError(err)
	}

	buff := bytes.NewBuffer(buf)
	n, readed, err := datatypes.ReadVarInt(buff)
	if err != nil {
		logError(err)
	}
	if uint16(n) != PROTOCOL_VERSION {
		logError("Mismatched version from the client")
		return
	}
	stringLenght := uint8(buf[readed])
	stringOffset := stringLenght + uint8(readed) + 1
	// serverAdress := string(buf[readed+1 : stringOffset])

	// port := binary.BigEndian.Uint16(buf[stringOffset : stringOffset+2])

	intention, _, err := datatypes.ReadVarInt(bytes.NewBuffer(buf[stringOffset+2:]))

	switch intention {
	case 1:
		initialPayload := make([]byte, 2)
		_, err := conn.Read(initialPayload)
		if err != nil {
			logError(err)
		}

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
		lenghtForResponsePayload := datatypes.WriteVarInt(int32(len(statusSerialized)))
		packageID := datatypes.WriteVarInt(packets.PACKET_HANDSHAKE)

		totalLenght := len(packageID) + lenghtForResponse + len(lenghtForResponsePayload)

		packageLenght := datatypes.WriteVarInt(int32(totalLenght))

		response := make([]byte, len(packageLenght)+totalLenght)

		offset := 0
		offset += copy(response[offset:], packageLenght)
		offset += copy(response[offset:], packageID)
		offset += copy(response[offset:], lenghtForResponsePayload)
		copy(response[offset:], statusSerialized)

		conn.Write(response)

	case 2:
		initialPayload := make([]byte, 512)
		n, err := conn.Read(initialPayload)
		if err != nil {
			logError(err)
			return
		}
		buf := make([]byte, n)

		copy(buf, initialPayload[:n])
		initialPayload = nil

		_, offset, err := datatypes.ReadVarInt(bytes.NewReader(buf))
		if err != nil {
			logError(err)
			return
		}

		protocolID, tmp, err := datatypes.ReadVarInt(bytes.NewReader(buf[offset:]))
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

		stringLenght, tmp, err := datatypes.ReadVarInt(bytes.NewReader(buf[offset:]))
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
			logInfo(fmt.Sprintf("Received connection from %s  %s (%s)", clientAddr, string(playerName), playerUUID.String()))

			packetID := datatypes.WriteVarInt(packets.PACKET_ENCRYPTION_REQUEST)

			serverID := []byte("")
			serverIDPrefix := datatypes.WriteVarInt(int32(len(serverID)))

			publicKeyMarshal, err := os.ReadFile("server-public.key")
			if err != nil {
				logError(err)
				return
			}
			publicKeyMarshalPrefix := datatypes.WriteVarInt(int32(len(publicKeyMarshal)))

			verifyToken := make([]byte, 4)
			rand.Read(verifyToken)
			verifyTokenPrefix := datatypes.WriteVarInt(4)

			var shouldAuth byte
			shouldAuth = TRUE_VALUE

			totalLenght := len(packetID) + len(serverIDPrefix) + len(serverID) + len(publicKeyMarshalPrefix) + len(publicKeyMarshal) + len(verifyTokenPrefix) + len(verifyToken) + 1 // +1 for shouldAuth
			lenght := datatypes.WriteVarInt(int32(totalLenght))

			response := make([]byte, totalLenght+len(lenght))

			offset := copy(response, lenght)
			offset += copy(response[offset:], packetID)

			offset += copy(response[offset:], serverIDPrefix)
			offset += copy(response[offset:], serverID)

			offset += copy(response[offset:], publicKeyMarshalPrefix)
			offset += copy(response[offset:], publicKeyMarshal)

			offset += copy(response[offset:], verifyTokenPrefix)
			copy(response[offset:], verifyToken)

			response[len(response)-1] = shouldAuth

			conn.Write(response)

			initialPayload := make([]byte, 4096)
			n, err := conn.Read(initialPayload)
			if err != nil {
				logError(err)
			}
			buf := make([]byte, n)

			copy(buf, initialPayload[:n])
			initialPayload = nil

			_, offset, err = datatypes.ReadVarInt(bytes.NewReader(buf))
			if err != nil {
				logError(err)
				return
			}

			protocolID, tmp, err := datatypes.ReadVarInt(bytes.NewReader(buf[offset:]))
			if err != nil {
				logError(err)
				return
			}
			offset += tmp

			if protocolID != packets.PACKET_ENCRYPTION_RESPONSE {
				logError("Player failed to answer encryption response")
				return
			}

			sharedSecretLenght, tmp, err := datatypes.ReadVarInt(bytes.NewBuffer(buf[offset:]))
			if err != nil {
				logError(err)
				return
			}
			offset += tmp
			tmp = 0

			sharedSecretCipher := make([]byte, sharedSecretLenght)

			copy(sharedSecretCipher, buf[offset:offset+int(sharedSecretLenght)])

			verifyTokenLenght, tmp, err := datatypes.ReadVarInt(bytes.NewBuffer(buf[offset+int(sharedSecretLenght):]))

			offset = offset + int(sharedSecretLenght) + tmp

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
					logError("Mojang didn't responded player check for the 10th time in 10 secods")
					return
				}
				if resp.StatusCode != 200 {
					resp.Body.Close()
				} else {
					break
				}
				time.Sleep(1 * time.Second)
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

			totalLenght = 0
			packetID = datatypes.WriteVarInt(packets.PACKET_LOGIN_SUCCESS)
			totalLenght += len(packetID)

			binUUUID, err := playerUUID.MarshalBinary()
			if err != nil {
				logError(err)
				return
			}
			totalLenght += len(binUUUID)

			playerNameLenght := datatypes.WriteVarInt(int32(len(playerName)))
			totalLenght += len(playerNameLenght) + len(playerName)

			arrayEntryPoint := datatypes.WriteVarInt(1)
			totalLenght += len(arrayEntryPoint)

			nameLenght := datatypes.WriteVarInt(int32(len([]byte(name))))
			totalLenght += len(nameLenght) + len([]byte(name))

			valueLenght := datatypes.WriteVarInt(int32(len([]byte(value))))
			totalLenght += len(valueLenght) + len([]byte(value))

			signatureLenght := datatypes.WriteVarInt(int32(len([]byte(signature))))
			totalLenght += len(signatureLenght) + len([]byte(signature))

			totalLenght += 1 // bool signature value
			lenght = datatypes.WriteVarInt(int32(totalLenght))

			response = make([]byte, totalLenght+len(lenght))

			offset = copy(response, lenght)
			offset += copy(response[offset:], packetID)
			offset += copy(response[offset:], binUUUID)
			offset += copy(response[offset:], playerNameLenght)
			offset += copy(response[offset:], playerName)
			offset += copy(response[offset:], arrayEntryPoint)
			offset += copy(response[offset:], nameLenght)
			offset += copy(response[offset:], []byte(name))
			offset += copy(response[offset:], valueLenght)
			offset += copy(response[offset:], []byte(value))
			offset += copy(response[offset:], []byte{TRUE_VALUE})
			offset += copy(response[offset:], signatureLenght)
			copy(response[offset:], []byte(signature))

			fmt.Println(totalLenght)
			block, err := aes.NewCipher(sharedSecret)
			if err != nil {
				logError(err)
				return
			}

			encrypter := cipher.NewCTR(block, sharedSecret)

			responseCipher := make([]byte, len(response))
			encrypter.XORKeyStream(responseCipher, response)
			conn.Write(responseCipher)

		} else {

			playerUUID := offline.NameToUUID(string(playerName))

			proprieties := datatypes.WriteVarInt(0)

			packetID := datatypes.WriteVarInt(packets.PACKET_LOGIN_SUCCESS)

			binaryUUID, err := playerUUID.MarshalBinary()
			if err != nil {
				logError(err)
				return
			}

			stringSize := datatypes.WriteVarInt(int32(len(playerName)))

			totalLenght := len(packetID) + len(binaryUUID) + len(stringSize) + len(playerName) + len(proprieties)
			lenght := datatypes.WriteVarInt(int32(totalLenght))

			response := make([]byte, totalLenght+len(lenght))

			offset := copy(response, lenght)
			offset += copy(response[offset:], packetID)
			offset += copy(response[offset:], binaryUUID)
			offset += copy(response[offset:], stringSize)
			offset += copy(response[offset:], playerName)
			copy(response[offset:], proprieties)

			conn.Write(response)
		}

	}

	for {
		innerBuffer := make([]byte, 4096)
		n, err := conn.Read(innerBuffer)
		if err != nil {
			if err == io.EOF {
				continue
			}
			logError(err)
		}
		if n == 0 {
			continue
		}

		buf := make([]byte, n)
		copy(buf, innerBuffer[:n])
		innerBuffer = nil

		packetLenght, offset, err := datatypes.ReadVarInt(bytes.NewReader(buf))
		if err != nil {
			logError(err)
		}

		protocolID, tmp, err := datatypes.ReadVarInt(bytes.NewReader(buf[offset:]))
		if err != nil {
			logError(err)
		}
		offset += tmp
		tmp = 0

		switch {
		case protocolID == packets.PACKET_PING && packetLenght == 10:
			conn.Write(packets.SerializePong(buf[offset:]))
		}
	}
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
