package main

import (
	"os"

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
	bk, err := NewBackend(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Add goroutines for recurring jobs
	log.Info("Serving...")
	bk.Run()

}
