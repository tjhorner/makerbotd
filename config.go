package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	connectionTypeLocal  = "local"
	connectionTypeRemote = "remote"
)

type printerConfig struct {
	ConnectionType string
	ID             string
	IP             string
	Port           string
}

type config struct {
	Debug               bool            // Debug makes output more verbose
	ThingiverseUsername string          // ThingiverseUsername defines the username of the authenticated Thingiverse account
	ThingiverseToken    string          // ThingiverseToken defines the auth token of the authenticated Thingiverse account
	ListenSocket        bool            // ListenSocket defines whether or not makerbotd will listen on a unix domain socket
	ListenSocketPath    string          // ListenSocketPath defines the unix domain socket to listen on if ListenSocket is true
	ListenTCP           bool            // ListenTCP defines whether or not makerbotd will listen on a TCP port
	ListenTCPAddress    string          // ListenTCPPort defines the TCP port to listen on if ListenTCP is true
	AutoAddPrinters     bool            // AutoAddPrinters defines whether or not printers should automatically be added from the authenticated Thingiverse account
	Printers            []printerConfig // Printers is the list of MakerBot printers that will automatically be connected when makerbotd starts
}

func writeDefaultConfig(path string) (*config, error) {
	dc := config{
		Debug:            false,
		AutoAddPrinters:  false,
		ListenSocket:     true,
		ListenSocketPath: "/var/run/makerbot.socket",
		ListenTCP:        false,
		ListenTCPAddress: ":6969", // har har
		Printers:         []printerConfig{},
	}

	conf, err := json.MarshalIndent(dc, "", "  ")
	if err != nil {
		return nil, err
	}

	return &dc, ioutil.WriteFile(path, conf, 0600)
}

func getConfig(path string) (*config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf config
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func loadConfig(path string) (*config, error) {
	conf, err := getConfig(path)
	if err != nil {
		err = os.MkdirAll(filepath.Dir(path), 0766)
		if err != nil {
			return nil, err
		}

		_, err = os.Create(path)
		if err != nil {
			return nil, err
		}

		conf, err = writeDefaultConfig(path)
		if err != nil {
			return nil, err
		}
	}

	return conf, err
}
