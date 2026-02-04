package seed

import (
	"crypto/rand"
	"encoding/binary"
)



func javaStringHash(s string) int32 {
	var hash int32
	for _, r := range s {
		hash = 31*hash + int32(r)
	}
	return hash
}

func GenerateSeed() int64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return int64(binary.BigEndian.Uint64(b[:]))
}

func GenerateByString(s string) int64 {
	return int64(javaStringHash(s))
}


