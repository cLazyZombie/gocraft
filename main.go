package main

import (
	"flag"
	"log"
	"time"

	_ "image/png"

	"net/http"
	_ "net/http/pprof"

	. "github.com/cLazyZombie/gocraft/internal"
	"github.com/faiface/mainthread"
)

var (
	pprofPort = flag.String("pprof", "", "http pprof port")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	flag.Parse()
	go func() {
		if *pprofPort != "" {
			log.Fatal(http.ListenAndServe(*pprofPort, nil))
		}
	}()
	mainthread.Run(run)
}

func run() {
	err := LoadTextureDesc()
	if err != nil {
		log.Fatal(err)
	}

	err = InitStore()
	if err != nil {
		log.Fatal(err)
	}
	defer GlobalStore.Close()

	game, err := NewGame(800, 600)
	if err != nil {
		log.Fatal(err)
	}
	game.Camera().Restore(GlobalStore.GetCamera())
	tick := time.Tick(time.Second / 60)
	for !game.ShouldClose() {
		<-tick
		game.Update()
	}
	GlobalStore.UpdateCamera(game.Camera().State())
}
