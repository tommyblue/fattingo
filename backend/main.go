package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)

	flag, isPresent := os.LookupEnv("DEBUG")
	if isPresent && flag == "1" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	// time.Sleep(5 * time.Second)
	cfg, err := readConf()
	if err != nil {
		log.Fatal(err)
	}
	bk, err := newBackend(cfg)
	if err != nil {
		log.Fatal(err)
	}

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := bk.run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-quitCh
	log.Warn("CTRL+C caught, doing clean shutdown (use CTRL+\\ aka SIGQUIT to abort)")
	bk.stop()
}
