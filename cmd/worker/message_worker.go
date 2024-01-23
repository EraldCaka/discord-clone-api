package main

import (
	"github.com/EraldCaka/discord-clone-api/pkg/workers/message/pkg"
	"log"
)

func main() {
	cfg := &pkg.Config{
		ListenAddr: ":3003",
		StoreProducerFunc: func() pkg.Storer {
			return pkg.NewMemoryStore()
		},
	}
	s, err := pkg.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
}
