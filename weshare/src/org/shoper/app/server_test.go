package main

import entry "org/shoper/app/message"
import proto "code.google.com/p/goprotobuf/proto"
import (
	"fmt"
	"log"
	"scan"
)

func main() {
	var requestTime = "2014-12-12"
	msg := &entry.Handshake{
		RequestTime: &requestTime,
	}
	encObj, err := proto.Marshal(msg)
	if nil == err {
		fmt.Println("length:", len(encObj))
		tobj := &entry.Handshake{}
		e := proto.Unmarshal(encObj, tobj)
		if nil == e {
			fmt.Println(tobj.GetRequestTime())
		} else {
			log.Fatalln("decode fail ", e)
		}
	} else {
		log.Fatalln("encode fail", err)
	}
}
