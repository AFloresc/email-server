package security

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
)

func HashIP(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		ip = remoteAddr
	}

	hash := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(hash[:])
}
