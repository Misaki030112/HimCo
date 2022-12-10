package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"web.misaki.world/FinalExam/handler"
)

var (
	storageParentDir string
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime)
	argsAnalyze()
	ctx, cancelCtx := context.WithCancel(context.Background())
	mux := http.NewServeMux()
	mux.HandleFunc("/album", handler.CrawlAlbum)
	mux.HandleFunc("/convert", handler.ConvertAudioToJson)
	mux.HandleFunc("/analyzeJson", handler.AnalyzeOutJson)
	server := &http.Server{
		Addr:    ":80",
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			return context.WithValue(ctx, "StorageParentDir", storageParentDir)
		},
	}

	err := server.ListenAndServe()
	cancelCtx()
	log.Panicf("can not start server,Here is the reason:\n%v\n", err)

}

func argsAnalyze() {
	flag.StringVar(&storageParentDir, "s", "./", "the storage Parent Path")
	flag.Parse()
}
