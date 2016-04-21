package utilities

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5Sum calculate and return the input string's md5
// Used for users' password
func Md5Sum(enc string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(enc))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
