package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/tjhorner/makerbot-rpc"
)

type apiResult struct {
	Error  *string         `json:"error"`
	Result json.RawMessage `json:"result"`
}

// Client is a client that talks to makerbotd
type Client struct {
	http    *http.Client
	baseURL string
}

func (c *Client) url(endpoint string) string {
	return fmt.Sprintf("%s%s", c.baseURL, endpoint)
}

func (c *Client) request(req *http.Request, result interface{}) error {
	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	dec := json.NewDecoder(res.Body)

	var ar apiResult

	dec.Decode(&ar)
	if ar.Error != nil {
		return errors.New(*ar.Error)
	}

	return json.Unmarshal(ar.Result, &result)
}

func (c *Client) httpGet(endpoint string, result interface{}) error {
	req, err := http.NewRequest("GET", c.url(endpoint), nil)
	if err != nil {
		return err
	}

	return c.request(req, result)
}

func (c *Client) httpPost(endpoint string, result interface{}) error {
	req, err := http.NewRequest("POST", c.url(endpoint), nil)
	if err != nil {
		return err
	}

	return c.request(req, result)
}

func (c *Client) httpPostFile(endpoint string, path string, result interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("printfile", stat.Name())
	if err != nil {
		return err
	}
	part.Write(data)

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.url(endpoint), body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	return c.request(req, result)
}

func (c *Client) httpDelete(endpoint string, result interface{}) error {
	req, err := http.NewRequest("DELETE", c.url(endpoint), nil)
	if err != nil {
		return err
	}

	return c.request(req, result)
}

// GetPrinters gets a list of connected printers from makerbotd
func (c *Client) GetPrinters() (*[]makerbot.Printer, error) {
	var printers []makerbot.Printer

	err := c.httpGet("/api/v1/printers", &printers)
	if err != nil {
		return nil, err
	}

	return &printers, nil
}

// GetPrinter gets a printer with `id`
func (c *Client) GetPrinter(id string) (*makerbot.Printer, error) {
	var printer makerbot.Printer

	err := c.httpGet("/api/v1/printers/"+id, &printer)
	if err != nil {
		return nil, err
	}

	return &printer, nil
}

// GetPrinterSnapshot gets a single frame from the printer's camera
func (c *Client) GetPrinterSnapshot(id string) (*[]byte, error) {
	req, err := http.NewRequest("GET", c.url(fmt.Sprintf("/api/v1/printers/%s/snapshot.jpg?%d", id, time.Now().Unix())), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// Print tells makerbotd to print on a specified printer
func (c *Client) Print(id, path string) (*bool, error) {
	var result bool

	err := c.httpPostFile("/api/v1/printers/"+id+"/prints", path, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCurrentJob gets a current job yeet
func (c *Client) GetCurrentJob(id string) (*makerbot.PrinterProcess, error) {
	var job makerbot.PrinterProcess

	err := c.httpGet("/api/v1/printers/"+id+"/current_job", &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

// SuspendCurrentJob tells makerbotd to suspend the current job on a specified printer
func (c *Client) SuspendCurrentJob(printerID string) (*bool, error) {
	var result bool

	err := c.httpPost("/api/v1/printers/"+printerID+"/current_job/suspend", &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ResumeCurrentJob tells makerbotd to resume the current job on a specified printer
func (c *Client) ResumeCurrentJob(printerID string) (*bool, error) {
	var result bool

	err := c.httpPost("/api/v1/printers/"+printerID+"/current_job/resume", &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CancelCurrentJob tells makerbotd to cancel the current job on a specified printer
func (c *Client) CancelCurrentJob(printerID string) (*bool, error) {
	var result bool

	err := c.httpDelete("/api/v1/printers/"+printerID+"/current_job", &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// LoadFilament tells makerbotd to tell the specified printer to start loading filament
func (c *Client) LoadFilament(printerID string, toolIndex string) (*bool, error) {
	var result bool

	err := c.httpPost("/api/v1/printers/"+printerID+"/load_filament/"+toolIndex, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// UnloadFilament tells makerbotd to tell the specified printer to start unloading filament
func (c *Client) UnloadFilament(printerID string, toolIndex string) (*bool, error) {
	var result bool

	err := c.httpPost("/api/v1/printers/"+printerID+"/unload_filament/"+toolIndex, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
