package utilities
import (
    "crypto/md5"
    "encoding/hex"
)

func Md5Sum(enc string) string{
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(enc))
    cipherStr := md5Ctx.Sum(nil)
    return hex.EncodeToString(cipherStr)
}
