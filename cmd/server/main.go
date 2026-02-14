package main

import (
	"MoonMS/app"
	"MoonMS/internal/packets"
	"crypto/rand"
	"crypto/rsa"
	"os"
	"os/signal"
	"time"

	_ "embed"
)

const MOJANG_SESSION_CHECKER = `https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s`
const DEADLINE = time.Second * 30

type MojangAnswer struct {
	Properties []map[string]string
}

var AnonymousPlayer = &packets.PlayerMinimunInfo{Username: "Anonymous Player", UUID: "00000000-0000-0000-0000-000000000000"}

func main() {

	cfg := app.MinecraftServerConfig{}
	if err := cfg.ConfigFile(); err != nil {
		panic(err)
	}

	scfg := app.Config{
		LatestLogFile: "logs/latest.log",
		StartName:     "java",
		DebugEnabled:  false,
		PluginsFolder: "plugins",
	}

	if os.Getenv("DEBUG") == "true" {
		scfg.DebugEnabled = true
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, int(cfg.Advanced.RSAKeyBitAmmount))
	if err != nil {
		panic(err)
	}

	server := app.New(cfg, scfg, privateKey)
	if err := server.StartLogger(); err != nil {
		panic(err)
	}
	server.InitPlugins()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig

		server.Stop()

		os.Exit(0)
	}()

	server.Start()
}

/*
func handleConnection(conn net.Conn) {
	defer conn.Close()
	server.Debug("CONN -> ", conn.LocalAddr())

	var EncryptKey cipher.Stream

	reader := packets.NewReaderFromReader(conn)

	handshake, err := packets.UnmarshalPacket(reader, -1)
	if err != nil {
		server.LogError(err)
		return
	}

	if handshake.ID() != packets.PACKET_HANDSHAKE {
		server.LogError("The player didn't started the connection with a handshake")
		return
	}

	switch handshake.(*packets.Handshake).Intent {

	case 2:
		packets.RegisterPacket(packets.PACKET_ENCRYPTION_RESPONSE, func() packets.Packet {
			return &packets.EncryptionResponse{}
		})

		Deadline := time.Now().Add(DEADLINE)
		conn.SetReadDeadline(Deadline)

		buff, err := packets.ReadPackageFromConnecion(conn)
		if err != nil {
			server.LogError(err)
			return
		}

		_, offset, err := datatypes.ParseVarInt(bytes.NewReader(buff))
		if err != nil {
			server.LogError(err)
			return
		}

		protocolID, tmp, err := datatypes.ParseVarInt(bytes.NewReader(buff[offset:]))
		if err != nil {
			server.LogError(err)
			return
		}
		offset += tmp
		tmp = 0

		if protocolID != packets.PACKET_HANDSHAKE {
			server.LogError("Invalid package received from player")
			return
		}

		stringLenght, tmp, err := datatypes.ParseVarInt(bytes.NewReader(buff[offset:]))
		if err != nil {
			server.LogError(err)
			return
		}

		offset += tmp
		tmp = 0

		playerName := buff[offset : offset+int(stringLenght)]
		clientAddr := conn.RemoteAddr().String()

		if ServerData.OnlineMode {
			playerUUID, err := uuid.FromBytes(buff[offset+int(stringLenght):])
			if err != nil {
				server.LogError(err)
				return
			}
			server.LogInfo(fmt.Sprintf("Received connection from %s as %s (%s)", clientAddr, string(playerName), playerUUID.String()))

			publicKeyMarshal, err := os.ReadFile("server-public.key")
			if err != nil {
				server.LogError(err)
				return
			}

			verifyToken := make([]byte, 4)
			rand.Read(verifyToken)

			var response packets.EncryptionRequest
			response.ServerID = ""
			response.PublicKey = publicKeyMarshal
			response.VerifyToken = verifyToken
			response.ShouldAuth = true
			data, err := packets.MarshalPacket(&response, nil, -1)
			if err != nil {
				server.LogError(err)
				return
			}
			conn.Write(data)

			server.Debug("Write Encryption request")

			pkg, err := packets.UnmarshalPacket(reader, -1)
			if err != nil {
				return
			}

			encR := pkg.(*packets.EncryptionResponse)

			privKeyPem, err := os.ReadFile("server-private-key.pem")
			if err != nil {
				server.LogError(err)
				return
			}

			p, _ := pem.Decode(privKeyPem)

			privateKey, err := x509.ParsePKCS1PrivateKey(p.Bytes)
			if err != nil {
				server.LogError(fmt.Sprintf("error parsing private key: %v", err))
				return
			}

			plainToken, err := privateKey.Decrypt(nil, encR.VerifyTokenCiphered, nil)
			if err != nil {
				server.LogError(err)
				return
			}

			server.Debug("Im still here")

			if !bytes.Equal(plainToken, verifyToken) {
				server.LogError("Invalid Session")
				return
			}

			sharedSecret, err := privateKey.Decrypt(nil, encR.SharedSecretCiphered, nil)
			if err != nil {
				server.LogError(err)
				return
			}

			hash := sha1.New()
			hash.Write([]byte(""))
			hash.Write(sharedSecret)
			hash.Write(publicKeyMarshal)

			sum := hash.Sum(nil)

			var resp *http.Response
			for i := range 10 {
				resp, err = http.Get(fmt.Sprintf(MOJANG_SESSION_CHECKER, string(playerName), hex.EncodeToString(sum)))
				if err != nil {
					server.LogError(err)
					return
				}
				if resp.StatusCode != 200 && i == 9 {
					defer resp.Body.Close()
					payload, err := packets.DisconnectLogin("Invalid session")
					if err != nil {
						server.LogError(err)
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
				server.LogError(err)
				return
			}

			var mojangResponse MojangAnswer
			if err := json.Unmarshal(mojangPayload, &mojangResponse); err != nil {
				server.LogError(err)
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

			loginSuccesspkg.Profile.UUID, err = playerUUID.MarshalBinary()
			if err != nil {
				server.LogError(err)
				return
			}
			loginSuccesspkg.Profile.Username = string(playerName)
			loginSuccesspkg.Profile.Name = name
			loginSuccesspkg.Profile.Value = value
			loginSuccesspkg.Profile.Signature = signature

			block, err := aes.NewCipher(sharedSecret)
			if err != nil {
				server.LogError(err)
				return
			}

			EncryptKey = cipher.NewCTR(block, sharedSecret)
			// Compression
			compressionStartPkg := &packets.CompressionStart{
				Threshould: ServerData.Threshold,
			}

			marshalCompressStart, err := packets.MarshalPacket(compressionStartPkg, EncryptKey, -1)
			if err != nil {
				server.LogError(err)
				return
			}
			conn.Write(marshalCompressStart)

			server.Debug("Write Compress Start")

			// Compression

			reader = packets.NewReaderFromReader(packets.NewCipherReader(conn, EncryptKey))

			marshalized, err := packets.MarshalPacket(&loginSuccesspkg, EncryptKey, int(ServerData.Threshold))
			if err != nil {
				server.LogError(err)
				return
			}
			conn.Write(marshalized)

			server.Debug("Sent everything")

			aknowledge := make([]byte, 3)
			_, err = conn.Read(aknowledge)
			if err != nil {
				return
			}

		} else {

			//		compressionStartPkg := &packets.CompressionStart{
			//			Threshould: ServerData.Threshold,
			//		}

			//	conn.Write(compressionStartPkg.Serialize())

			playerUUID := offline.NameToUUID(string(playerName))

			server.LogInfo(fmt.Sprintf("Received connection from %s as %s (%s) [OFFLINE]", clientAddr, string(playerName), playerUUID.String()))
			var loginSuccesspkg packets.LoginSuccessPacket
			var namebuff bytes.Buffer

			namebuff.Write(datatypes.NewVarInt(int32(len(playerName))))
			namebuff.Write(playerName)
			//		loginSuccesspkg.Profile.Username = datatypes.String(namebuff.Bytes())
			namebuff.Reset()
			loginSuccesspkg.Profile.UUID, err = playerUUID.MarshalBinary()
			if err != nil {
				server.LogError(err)
				return
			}

			//		response := loginSuccesspkg.Serialize()

			//		compressedResponse, err := compress.Compress(response, ServerData.Threshold)
			if err != nil {
				server.LogError(err)
				return
			}

			//		conn.Write(compressedResponse)

			aknowledge, err := packets.ReadPackageFromConnecion(conn)
			if err != nil {
				server.LogError(err)
				return
			}

			if aknowledge[2] != byte(packets.PACKET_LOGIN_AKNOWLEDGED) {
				server.LogError("Invalid packet")
				return
			}
		}

	}

	conn.SetReadDeadline(time.Time{})

	// Configuration Process

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
*/
