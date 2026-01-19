package compress_test

import (
	"MoonMS/cmd/server/compress"
	"bytes"
	"fmt"
	"testing"
)

var payloadTest = []byte{121, 0, 119, 123, 34, 100, 101, 115, 99, 114, 105, 112, 116, 105, 111, 110, 34, 58, 34, 65, 32, 77, 105, 110, 101, 99, 114, 97, 102, 116, 32, 83, 101, 114, 118, 101, 114, 34, 44, 34, 112, 108, 97, 121, 101, 114, 115, 34, 58, 123, 34, 109, 97, 120, 34, 58, 50, 48, 44, 34, 111, 110, 108, 105, 110, 101, 34, 58, 48, 125, 44, 34, 118, 101, 114, 115, 105, 111, 110, 34, 58, 123, 34, 110, 97, 109, 101, 34, 58, 34, 80, 117, 114, 112, 117, 114, 32, 49, 46, 50, 49, 46, 49, 49, 34, 44, 34, 112, 114, 111, 116, 111, 99, 111, 108, 34, 58, 55, 55, 52, 125, 125}

func TestCompress(T *testing.T) {

	var uncompressed []byte
	T.Run("Testing Compress less than Threshold", func(t *testing.T) {
		var err error
		uncompressed, err = compress.Compress(payloadTest, 256)
		if err != nil {
			t.Errorf("Error compressing payload: %v", err)
		}
		fmt.Println(uncompressed)
	})

	var compressed []byte
	T.Run("Testing Compress higher than Threshold", func(t *testing.T) {
		var err error
		compressed, err = compress.Compress(payloadTest, 121)
		if err != nil {
			t.Errorf("Error compressing payload: %v", err)
		}
		fmt.Println(compressed)
	})

	T.Run("Testing Uncompress greater than Threshold", func(t *testing.T) {
		answer, err := compress.Uncompress(compressed)
		if err != nil {
			t.Errorf("Error uncompressing payload: %v", err)
		}
		if !bytes.Equal(answer, payloadTest) {
			t.Errorf("Uncompressed payload does not match original")
			fmt.Println(answer)
		}
	})

	T.Run("Testing Uncompress less than Threshold", func(t *testing.T) {
		answer, err := compress.Uncompress(uncompressed)
		if err != nil {
			t.Errorf("Error uncompressing payload: %v", err)
		}

		if !bytes.Equal(answer, payloadTest) {
			t.Errorf("Uncompressed payload does not match original")
			fmt.Println(answer)
		}
	})
}
