package file

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

func CalcMD5(file *os.File) string {
	b := make([]byte, 1024*1024)
	//info := bytes.Buffer{}
	m := md5.New()
	for {
		n, _ := file.Read(b)
		if n == 0 {
			break
		}
		m.Write(b[:n])
	}
	cipherStr := m.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
