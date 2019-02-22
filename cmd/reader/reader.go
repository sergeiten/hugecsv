package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/sergeiten/hugecsv"
	"github.com/sergeiten/hugecsv/reader"
)

var file string

func init() {
	flag.StringVar(&file, "file", "", "CSV file name")
	flag.Parse()
}

func main() {
	if file == "" {
		log.Fatal("CSV file name empty")
	}
	// sleep to wait consumer service
	ctx, cancel := context.WithCancel(context.Background())
	gracefulStop := make(chan os.Signal)

	signal.Notify(gracefulStop, os.Interrupt)

	reader := reader.New(file)

	go func() {
		fmt.Println("waiting for termination")
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v\n", sig)
		cancel()
		fmt.Println("wait for 1 second to finish processing")
		time.Sleep(time.Second)
		os.Exit(0)
	}()

	err := reader.Serve(ctx)
	hugecsv.LogFatal(err, "failed to start reader serve")
}
