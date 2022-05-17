package main

import (
	"flag"
	"log"
)

func main() {
	var target string
	var password string
	var downloadPath string
	var localPath string
	var reverse bool
	var tcp bool
	flag.StringVar(&target, "target", "127.0.0.1:58993", "Target IP:Port")
	flag.StringVar(&password, "password", "", "Password for specified target")
	flag.StringVar(&downloadPath, "downloadPath", "", "Path to be downloaded")
	flag.StringVar(&localPath, "localPath", "", "Path where file will be stored")
	flag.BoolVar(
		&reverse,
		"reverse",
		false,
		"Whether connection should be reversed to upload file from client, should match on both ends",
	)
	flag.BoolVar(&tcp, "tcp", false, "Whether to use TCP")
	flag.Parse()

	if downloadPath != "" && localPath != "" {
		log.Println("Starting client")
		runClient(target, password, downloadPath, localPath, reverse, tcp)
	} else {
		log.Println("Starting server")
		runServer(target, password, reverse, tcp)
	}
}
