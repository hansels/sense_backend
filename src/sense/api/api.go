package api

import (
	myRouter "github.com/hansels/sense_backend/common/router"
	"github.com/hansels/sense_backend/src/sense"
	"github.com/julienschmidt/httprouter"
)

func (a *API) Register(router *httprouter.Router) {
	router.GET("/ping", myRouter.HandleNow("/ping", a.Ping))

	router.POST("/check-user", myRouter.HandleNow("/check-user", a.CheckUser))
	router.POST("/login", myRouter.HandleNow("/login", a.Login))
	router.POST("/register", myRouter.HandleNow("/register", a.RegisterUser))
	router.POST("/predict", myRouter.HandleNow("/predict", a.Predict))
}

type API struct {
	Module *sense.Module
}

func New(module *sense.Module) *API {
	return &API{Module: module}
}
