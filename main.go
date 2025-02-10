package main

import (
	"log"
	"simple-file-server/certbot"
	"simple-file-server/config"
	"simple-file-server/fileserver"
	"simple-file-server/hetzner"
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
