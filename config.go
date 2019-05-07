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
	ConnectionType string // ConnectionType should be either "local" or "remote". "local" = direct connect via IP, "remote" = remotely connect via MakerBot Reflector service.
	ID             string // ID should be provided if the connection type is "remote". This is the ID of the printer as returned by MakerBot Reflector. It is usually the serial number.
	IP             string // IP should be provided if the connection type is "local"
	Port           string // Port should be provided if the connection type is "port"
}

type config struct {
	Debug               bool            // Debug makes output more verbose
	ThingiverseUsername string          // ThingiverseUsername defines the username of the authenticated Thingiverse account
	ThingiverseToken    string          // ThingiverseToken defines the auth token of the authenticated Thingiverse account
	ListenSocket        bool            // ListenSocket defines whether or not makerbotd will listen on a unix domain socket
	ListenSocketPath    string          // ListenSocketPath defines the unix domain socket to listen on if ListenSocket is true
	ListenTCP           bool            // ListenTCP defines whether or not makerbotd will listen on a TCP port
	ListenTCPAddress    string          // ListenTCPPort defines the TCP port to listen on if ListenTCP is true
	AutoAddPrinters     bool            // AutoAddPrinters defines whether or not printers should automatically be added from the authenticated Thingiverse account (DOES NOTHING RIGHT NOW)
	ReadOnly            bool            // ReadOnly makes the API exposed by makerbotd read-only, e.g. print jobs cannot be sent, cancelled, etc. This is useful if you are publicly exposing the makerbotd API.
	Printers            []printerConfig // Printers is the list of MakerBot printers that will automatically be connected when makerbotd starts
}

func writeDefaultConfig(path string) (*config, error) {
	dc := config{
		Debug:            false,
		AutoAddPrinters:  false,
		ListenSocket:     true,
		ListenSocketPath: "/var/run/makerbot.socket",
		ListenTCP:        false,
		ListenTCPAddress: ":6969", // nice
		ReadOnly:         false,
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
