package main

import (
	"gitlab.com/kendellfab/mtgstudio/internal"
	"log"
	"net/http"
)
import _ "net/http/pprof"

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	mtgStudio := internal.NewMtgStudio()
	mtgStudio.Start()
}
