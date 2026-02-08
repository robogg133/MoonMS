package main

import (
	"MoonMS/cmd/server/compress"
	"MoonMS/cmd/server/config"
	"MoonMS/cmd/server/crypto"
	"MoonMS/internal/datatypes"
	"MoonMS/internal/offline"
	"MoonMS/internal/packets"
	"MoonMS/internal/plugins"
	"MoonMS/internal/server"
	"bytes"
	"compress/gzip"
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
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	_ "embed"

	"github.com/google/uuid"
)

const MOJANG_SESSION_CHECKER = `https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s`
const DEADLINE = time.Second * 30

type MojangAnswer struct {
	Properties []map[string]string
}

var AnonymousPlayer = &packets.PlayerMinimunInfo{Username: "Anonymous Player", UUID: "00000000-0000-0000-0000-000000000000"}

var ServerData *server.ServerData

func CheckFilesToStart() error {

	_ = os.Mkdir("plugins", 0755)
	_ = os.Mkdir("logs", 0755)

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
			server.LogInfo("Generating server key pair")
			if err := crypto.GenerateKeyPair(config.DefaultValuesForServerConfig().Proprieties.RSAKeyBits); err != nil {
				server.LogFatal(err)
			}
		}
	} else {
		if _, err = os.Stat("server-public.key"); err != nil {
			if os.IsNotExist(err) {
				server.LogError("Public key don't exist, but private exist, generating another public key from the private key")
				pemKey, err := os.ReadFile("server-private-key.pem")
				if err != nil {
					server.LogFatal(err)
				}
				p, _ := pem.Decode(pemKey)

				privateKey, err := x509.ParsePKCS1PrivateKey(p.Bytes)
				if err != nil {
					server.LogFatal(err)
				}
				publicKey := privateKey.Public()
				publicKeyMarshalized, err := x509.MarshalPKIXPublicKey(publicKey)
				if err != nil {
					server.LogFatal(err)
				}

				if err := os.WriteFile("server-public.key", publicKeyMarshalized, 0644); err != nil {
					log.Fatal(err)
				}
			}
		}

	}

	return nil
}

func CompressLog() error {

	oldLog, err := os.Open("logs/latest.log")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer oldLog.Close()

	f, err := os.Create(fmt.Sprintf("logs/%s.log.gz", time.Now().Format("2006-01-02 15:04:05")))
	if err != nil {
		return err
	}
	defer f.Close()

	writer := gzip.NewWriter(f)

	_, err = io.Copy(writer, oldLog)
	if err != nil {
		return err
	}

	return os.Remove("logs/latest.log")
}

func main() {

	packets.Init()

	var err error

	if err = CompressLog(); err != nil {
		server.LogError(err)
	}

	if err = CheckFilesToStart(); err != nil {
		server.LogError(err)
		return
	}
	ServerData, err = server.InitServerData()
	if err != nil {
		server.LogFatal(err)
	}

	allDirFiles, err := os.ReadDir("plugins")
	if err != nil {
		server.LogFatal(err)
	}

	server.LogInfo("Intializing plugins")
	for _, d := range allDirFiles {
		if d.IsDir() {
			continue
		}

		path := filepath.Join("plugins", d.Name())

		f, err := os.Open(path)
		if err != nil {
			server.LogError(fmt.Sprintf("Error opening file: %v", err))
			server.LogInfo(fmt.Sprintf("SKIPPING %v", d.Name()))
			continue
		}

		stat, err := f.Stat()
		if err != nil {
			server.LogError(fmt.Sprintf("Error getting file status: %v", err))
			server.LogInfo(fmt.Sprintf("SKIPPING %v", d.Name()))
			continue
		}

		plugin, err := plugins.ReadPluginFile(f, stat.Size(), path)
		if err != nil {
			server.LogError(fmt.Sprintf("Error parsing plugin: %v", err))
			server.LogInfo(fmt.Sprintf("SKIPPING %v", d.Name()))
			continue
		}

		server.LogInfo(fmt.Sprintf("Starting %s", plugin.Identifier))
		go plugin.LoadPlugin()
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig

		var wg sync.WaitGroup

		for _, pl := range *plugins.GetAllPlugins() {
			wg.Add(1)
			go pl.RunEventServerStopping(&wg)
		}

		wg.Wait()

		os.Exit(0)
	}()

	server.LogInfo(fmt.Sprintf("Starting minecraft java edittion tcp listener on port %d", ServerData.ServerPort))
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ServerData.ServerPort))
	if err != nil {
		server.LogFatal(err)
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

	reader := packets.NewReaderFromReader(conn)

	handshake, err := packets.UnmarshalPacket(reader)
	if err != nil {
		server.LogError(err)
		return
	}

	if handshake.ID() != packets.PACKET_HANDSHAKE {
		server.LogError("The player didn't started the connection with a handshake")
		return
	}

	switch handshake.(*packets.Handshake).Intent {
	case 1:
		statuspkg, err := packets.ReadPackageFromConnecion(conn)
		if err != nil {
			server.LogError(err)
			return
		}

		if int32(statuspkg[1]) != packets.PACKET_HANDSHAKE {
			server.LogError("mismatched packet id gaved by client")
			return
		}
		statuspkg = nil

		var status packets.HandShakeResponseStatus

		status.Version.Name = server.CURRENT_VERSION
		status.Version.ProtocolVersion = server.PROTOCOL_VERSION

		status.Players.MaxPlayers = ServerData.MaxPlayers
		status.Players.OnlinePlayers = 0
		status.Players.PlayerStatus = []packets.PlayerMinimunInfo{}

		if ServerData.ServerIcon == "" {
			status.Favicon = ""
		} else {
			status.Favicon, err = GetBase64Image(ServerData.ServerIcon)
			if err != nil {
				server.LogError(err)
				status.Favicon = ""
			} else {
				status.Favicon = fmt.Sprintf("data:image/png;base64,%s", status.Favicon)
			}
		}

		status.Description.Text = ServerData.Motd

		data, err := packets.MarshalPacket(&status, nil)
		if err != nil {
			server.LogError(err)
			return
		}
		conn.Write(data)

		ping, err := packets.UnmarshalPacket(reader)

		if ping.ID() == packets.PACKET_PING_PONG {
			data, err := packets.MarshalPacket(ping, nil)
			if err != nil {
				server.LogError(err)
				return
			}
			conn.Write(data)
		}

		return

	case 2:

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
			data, err := packets.MarshalPacket(&response, nil)
			if err != nil {
				server.LogError(err)
				return
			}
			conn.Write(data)

			compressionStartPkg := &packets.CompressionStart{
				Threshould: ServerData.Threshold,
			}

			marshalCompressStart, err := packets.MarshalPacket(compressionStartPkg, nil)
			if err != nil {
				server.LogError(err)
				return
			}
			conn.Write(marshalCompressStart)

			pkg, err := packets.UnmarshalPacket(reader)
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

			if !bytes.Equal(plainToken, verifyToken) {
				server.LogError("Invalid Session")
				return
			}

			sharedSecret, err = privateKey.Decrypt(nil, encR.SharedSecretCiphered, nil)
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
			var namebuff bytes.Buffer

			namebuff.Write(datatypes.NewVarInt(int32(len(playerName))))
			namebuff.Write(playerName)
			loginSuccesspkg.Profile.Username = datatypes.String(namebuff.Bytes())
			namebuff.Reset()
			loginSuccesspkg.Profile.UUID, err = playerUUID.MarshalBinary()
			if err != nil {
				server.LogError(err)
				return
			}
			loginSuccesspkg.Profile.Name = name
			loginSuccesspkg.Profile.Value = value
			loginSuccesspkg.Profile.HaveSignature = true
			loginSuccesspkg.Profile.Signature = signature
			block, err := aes.NewCipher(sharedSecret)
			if err != nil {
				server.LogError(err)
				return
			}

			encrypter := cipher.NewCTR(block, sharedSecret)
			tmpBuff, err := compress.Compress(loginSuccesspkg.Serialize(), ServerData.Threshold)
			if err != nil {
				server.LogError(err)
				return
			}

			packets.UnmarshalPacket(packets.NewReaderFromReader(packets.NewCipherReader(conn, encrypter)))

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
				Threshould: ServerData.Threshold,
			}

			conn.Write(compressionStartPkg.Serialize())

			playerUUID := offline.NameToUUID(string(playerName))

			server.LogInfo(fmt.Sprintf("Received connection from %s as %s (%s) [OFFLINE]", clientAddr, string(playerName), playerUUID.String()))
			var loginSuccesspkg packets.LoginSuccessPacket
			var namebuff bytes.Buffer

			namebuff.Write(datatypes.NewVarInt(int32(len(playerName))))
			namebuff.Write(playerName)
			loginSuccesspkg.Profile.Username = datatypes.String(namebuff.Bytes())
			namebuff.Reset()
			loginSuccesspkg.Profile.UUID, err = playerUUID.MarshalBinary()
			if err != nil {
				server.LogError(err)
				return
			}

			response := loginSuccesspkg.Serialize()

			compressedResponse, err := compress.Compress(response, ServerData.Threshold)
			if err != nil {
				server.LogError(err)
				return
			}

			conn.Write(compressedResponse)

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
