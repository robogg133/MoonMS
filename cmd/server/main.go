package main

import (
	"MoonMS/cmd/server/config"
	datatypes "MoonMS/internal/datatypes/varTypes"
	"MoonMS/internal/packets"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	_ "embed"
)

var PROTOCOL_VERSION uint16 = 774
var CURRENT_VERSION string = "1.21.11"

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

	clientAddr := conn.RemoteAddr().String()

	logInfo(fmt.Sprintf("Received connection from %s", clientAddr))

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
		}
		buf := make([]byte, n)

		copy(buf, initialPayload[:n])

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
