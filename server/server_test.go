package server

import (
	"log"
	"net"
	"testing"

	"github.com/gwuhaolin/lightsocks"
)

func TestCloseListen(t *testing.T) {
	password := lightsocks.RandPassword()
	localS, err := NewLsServer(password, "127.0.0.1:1234")
	if err != nil {
		log.Fatalln(err)
	}
	closeListen, err := localS.Listen(func(listenAddr *net.TCPAddr) {
		log.Println(listenAddr)
	})
	if err != nil {
		log.Fatalln(err)
	}
	localS, _ = NewLsServer(password, "127.0.0.1:1234")
	_, err = localS.Listen(func(listenAddr *net.TCPAddr) {
	})
	if err == nil {
		t.Error("closeListen 失败")
	}
	closeListen()
	localS, _ = NewLsServer(password, "127.0.0.1:1234")
	_, err = localS.Listen(func(listenAddr *net.TCPAddr) {
	})
	if err != nil {
		t.Error(err)
	}
}


