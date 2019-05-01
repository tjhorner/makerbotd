package main

import (
	"github.com/julienschmidt/httprouter"
)

// API represents an API (yep)
type API interface {
	Route(router *httprouter.Router)
}

type apiResult struct {
	Error  error       `json:"error"`
	Result interface{} `json:"result"`
}

func apiSuccess(result interface{}) apiResult {
	return apiResult{Result: result}
}

func apiError(err error) apiResult {
	return apiResult{Error: err}
}
