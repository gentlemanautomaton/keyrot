package keyrot

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func generate(bits int) string {
	// Convert bits to bytes (rounding up)
	blen := bits / 8
	if bits%8 != 0 {
		blen++
	}

	// Read the requested number of bytes from crypto/rand
	b := make([]byte, blen)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Errorf("unable to generate auth key: %v", err))
	}

	// Convert the bytes to a hexidecimal string
	value := make([]byte, hex.EncodedLen(blen))
	hex.Encode(value, b[:])
	s := string(value)

	return s
}
