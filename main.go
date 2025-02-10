package main

import (
	"log"
	"github.com/biggi93/simple-file-server/certbot"
	"github.com/biggi93/simple-file-server/config"
	"github.com/biggi93/simple-file-server/fileserver"
	"github.com/biggi93/simple-file-server/hetzner"
)

func main() {
	config.Init()

	go func() {
		err := fileserver.GetCertbotHttpTestServer().ListenAndServe()
		log.Fatal(err)
	}()

	h := hetzner.GetHetznerObject()

	h.AddTcpPortServiceForSSL()

	h.WaitForTcp80LBService()

	h.TestLbService()

	certbot.Run()

	h.DeleteTcpPort80ServiceForSSL()

}
