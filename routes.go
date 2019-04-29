package main

import (
	"github.com/julienschmidt/httprouter"
)

func getRouter(ctx *mbContext) *httprouter.Router {
	router := httprouter.New()

	v1 := APIv1{ctx}
	v1.Route(router)

	return router
}
