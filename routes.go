package main

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/julienschmidt/httprouter"
)

type statsResponse struct {
	NumGoroutine     int
	NumCPU           int
	AllocMemory      uint64
	TotalAllocMemory uint64
}

func stats(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	res := statsResponse{
		NumGoroutine:     runtime.NumGoroutine(),
		NumCPU:           runtime.NumCPU(),
		AllocMemory:      ms.Alloc,
		TotalAllocMemory: ms.TotalAlloc,
	}

	enc := json.NewEncoder(w)
	enc.Encode(apiSuccess(res))
}

func getRouter(ctx *mbContext) *httprouter.Router {
	router := httprouter.New()

	router.GET("/_/stats", stats)

	v1 := APIv1{ctx}
	v1.Route(router)

	return router
}
