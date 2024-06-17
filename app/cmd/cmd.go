package cmd

import "gptsovits/app/server"

func RunServer() {
	server.Init()
	server.Run()
}
