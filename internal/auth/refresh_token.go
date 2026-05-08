package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	key := make([]byte, 32)
	// fills key with crypotgraphically safe random bytes
	rand.Read(key)
	// returns a hexidecimal code as the refresh token
	return hex.EncodeToString(key)
}
