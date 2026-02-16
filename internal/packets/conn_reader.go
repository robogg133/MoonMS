package packets

import (
	"bytes"
	"fmt"
	"net"

	"github.com/robogg133/KernelCraft/internal/datatypes"
)

func ReadPackageFromConnecion(conn net.Conn) ([]byte, error) {

	startBuffer := make([]byte, 5)
	readedFromConn, err := conn.Read(startBuffer)
	if err != nil {
		return nil, err
	}

	if readedFromConn < 5 {
		return startBuffer[:readedFromConn], nil
	}

	lenght, readed, err := datatypes.ParseVarInt(bytes.NewReader(startBuffer))
	if err != nil {
		return nil, fmt.Errorf("failed to paser VarInt: %v", err)
	}

	needToReadAmmount := lenght - int32(5-readed)
	response := make([]byte, 5+needToReadAmmount)
	copy(response, startBuffer)

	n, err := conn.Read(response[5:])
	if err != nil {
		return nil, err
	}

	if int32(n) != needToReadAmmount {
		return nil, fmt.Errorf("reading %d but need to read %d", n, needToReadAmmount)
	}

	return response, nil
}
