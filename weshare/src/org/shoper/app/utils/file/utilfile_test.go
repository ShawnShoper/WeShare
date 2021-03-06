package file

import (
	//"bytes"
	"crypto/md5"
	"encoding/hex"
	//f "org/shoper/app/utils/file"
	"fmt"
	"os"
	"testing"
)

func TestToMD5(t *testing.T) {
	file, _ := os.Open("/Users/ShawnShoper/Downloads/cn_windows_7_ultimate_with_sp1_x64_dvd_u_677408.iso")
	b := make([]byte, 1024*1024)
	//info := bytes.Buffer{}
	m := md5.New()
	for {
		n, _ := file.Read(b)
		if n == 0 {
			break
		}

		//	info.Write(b[:n])
		m.Write(b[:n])
	}
	//[]byte(info.String())
	cipherStr := m.Sum(nil)
	fmt.Println("Generate md5 for :" + file.Name() + ":" + hex.EncodeToString(cipherStr))
}
