package main

import (
	"github.com/hansels/sense_backend/common/log"
	"github.com/hansels/sense_backend/src/firebase"
	"github.com/hansels/sense_backend/src/ml"
	"github.com/hansels/sense_backend/src/sense"
	"github.com/hansels/sense_backend/src/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	os.Exit(Main())
}

func Main() int {
	firestore := firebase.InitFirestore()
	storage := firebase.InitStorage()

	model := ml.NewCoco()
	err := model.Load()
	if err != nil {
		log.Errorf("Error loading model: %v", err)
		panic(err)
	}

	opts := &sense.Opts{Firestore: firestore, Storage: storage, Model: model}
	modules := sense.New(opts)

	api := server.New(&server.Opts{ListenAddress: ":3001", Modules: modules})

	go api.Run()

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-term:
		_ = firestore.Close()
		log.Println("Exiting gracefully...", s)
	}
	log.Info("👋")

	return 0
}
