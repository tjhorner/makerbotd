package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (e apiResult) MarshalJSON() ([]byte, error) {
	var es []byte
	var err error

	if e.Error != nil {
		es, err = json.Marshal(e.Error.Error())
		if err != nil {
			return nil, err
		}
	} else {
		es = []byte("null")
	}

	rs, err := json.Marshal(e.Result)
	if err != nil {
		return nil, err
	}

	return []byte(`{"result":` + string(rs) + `,"error":` + string(es) + `}`), nil
}

// APIv1 is version 1 of the API
type APIv1 struct {
	context *mbContext
}

// Route implements API.Route
func (a *APIv1) Route(router *httprouter.Router) {
	prefix := "/api/v1/"

	router.GET(prefix+"printers", a.getPrinters)
	router.GET(prefix+"printers/:id", a.getPrinter)
	router.GET(prefix+"printers/:id/snapshot.jpg", a.getPrinterSnapshot)
	router.GET(prefix+"printers/:id/current_job", a.getPrinterCurrentJob)
	router.POST(prefix+"printers/:id/current_job/suspend", a.postPrinterCurrentJobSuspend)
	router.POST(prefix+"printers/:id/current_job/resume", a.postPrinterCurrentJobResume)
	router.DELETE(prefix+"printers/:id/current_job", a.deletePrinterCurrentJob)
	router.POST(prefix+"printers/:id/prints", a.postPrinterPrints)
	router.POST(prefix+"printers/:id/unload_filament/:tool_index", a.postPrinterUnloadFilament)
	router.POST(prefix+"printers/:id/load_filament/:tool_index", a.postPrinterLoadFilament)
}

func (a *APIv1) notFound(w http.ResponseWriter, r *http.Request) {
	nf, _ := json.Marshal(apiError(errors.New("not found")))
	http.Error(w, string(nf), http.StatusNotFound)
}

func (a *APIv1) badRequest(w http.ResponseWriter, r *http.Request) {
	nf, _ := json.Marshal(apiError(errors.New("bad request")))
	http.Error(w, string(nf), http.StatusBadRequest)
}

func (a *APIv1) internalError(w http.ResponseWriter, r *http.Request) {
	nf, _ := json.Marshal(apiError(errors.New("internal server error")))
	http.Error(w, string(nf), http.StatusInternalServerError)
}

func (a *APIv1) getPrinters(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	enc.Encode(apiSuccess(a.context.Printers.ConnectedPrinters()))
}

func (a *APIv1) getPrinter(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	printer, ok := a.context.Printers.Find(params.ByName("id"))
	if !ok {
		a.notFound(w, r)
		return
	}

	enc := json.NewEncoder(w)
	enc.Encode(apiSuccess(printer.connection.Printer))
}

func (a *APIv1) getPrinterSnapshot(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	printer, ok := a.context.Printers.Find(params.ByName("id"))
	if !ok {
		a.notFound(w, r)
		return
	}

	frame, err := printer.connection.GetCameraFrame()
	if err != nil {
		a.internalError(w, r)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(frame.Data)
}

func (a *APIv1) getPrinterCurrentJob(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	printer, ok := a.context.Printers.Find(params.ByName("id"))
	if !ok {
		a.notFound(w, r)
		return
	}

	if printer.connection.Printer.Metadata == nil {
		fmt.Fprintf(w, "null\n")
		return
	}

	enc := json.NewEncoder(w)
	enc.Encode(apiSuccess(printer.connection.Printer.Metadata.CurrentProcess))
}

func (a *APIv1) postPrinterCurrentJobSuspend(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	printer, ok := a.context.Printers.Find(params.ByName("id"))
	if !ok {
		a.notFound(w, r)
		return
	}

	enc := json.NewEncoder(w)

	_, err := printer.connection.Suspend()
	if err != nil {
		enc.Encode(apiError(err))
		return
	}

	enc.Encode(apiSuccess(true))
}

func (a *APIv1) postPrinterCurrentJobResume(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	printer, ok := a.context.Printers.Find(params.ByName("id"))
	if !ok {
		a.notFound(w, r)
		return
	}

	enc := json.NewEncoder(w)

	_, err := printer.connection.Resume()
	if err != nil {
		enc.Encode(apiError(err))
		return
	}

	enc.Encode(apiSuccess(true))
}

func (a *APIv1) deletePrinterCurrentJob(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	printer, ok := a.context.Printers.Find(params.ByName("id"))
	if !ok {
		a.notFound(w, r)
		return
	}

	enc := json.NewEncoder(w)

	_, err := printer.connection.Cancel()
	if err != nil {
		enc.Encode(apiError(err))
		return
	}

	enc.Encode(apiSuccess(true))
}

func (a *APIv1) postPrinterPrints(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	printer, ok := a.context.Printers.Find(params.ByName("id"))
	if !ok {
		a.notFound(w, r)
		return
	}

	r.ParseMultipartForm(52428800)

	file, meta, err := r.FormFile("printfile")
	if err != nil {
		a.badRequest(w, r)
		return
	}
	defer file.Close()

	enc := json.NewEncoder(w)

	err = printer.connection.Print(meta.Filename, file, int(meta.Size))
	if err != nil {
		enc.Encode(apiError(err))
		return
	}

	enc.Encode(apiSuccess(true))
}

func (a *APIv1) postPrinterUnloadFilament(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	printer, ok := a.context.Printers.Find(params.ByName("id"))
	if !ok {
		a.notFound(w, r)
		return
	}

	ti, err := strconv.Atoi(params.ByName("tool_index"))
	if err != nil {
		a.badRequest(w, r)
		return
	}

	enc := json.NewEncoder(w)

	_, err = printer.connection.UnloadFilament(ti)
	if err != nil {
		enc.Encode(apiError(err))
		return
	}

	enc.Encode(apiSuccess(true))
}

func (a *APIv1) postPrinterLoadFilament(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	printer, ok := a.context.Printers.Find(params.ByName("id"))
	if !ok {
		a.notFound(w, r)
		return
	}

	ti, err := strconv.Atoi(params.ByName("tool_index"))
	if err != nil {
		a.badRequest(w, r)
		return
	}

	enc := json.NewEncoder(w)

	_, err = printer.connection.LoadFilament(ti)
	if err != nil {
		enc.Encode(apiError(err))
		return
	}

	enc.Encode(apiSuccess(true))
}
