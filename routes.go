package main

import (
	"encoding/json"
	"net/http"
	"net/http/pprof"
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

	if ctx.Config.Debug {
		router.GET("/_/stats", stats)

		router.HandlerFunc("GET", "/debug/pprof/", pprof.Index)
		router.HandlerFunc("GET", "/debug/pprof/cmdline", pprof.Cmdline)
		router.HandlerFunc("GET", "/debug/pprof/profile", pprof.Profile)
		router.HandlerFunc("GET", "/debug/pprof/symbol", pprof.Symbol)

		router.Handler("GET", "/debug/pprof/allocs", pprof.Handler("allocs"))
		router.Handler("GET", "/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.Handler("GET", "/debug/pprof/heap", pprof.Handler("heap"))
		router.Handler("GET", "/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		router.Handler("GET", "/debug/pprof/block", pprof.Handler("block"))
	}

	v1 := APIv1{ctx}
	v1.Route(router)

	return router
}
