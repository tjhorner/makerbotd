package main

import (
	"errors"

	"github.com/tjhorner/makerbot-rpc"
)

type printerConnection struct {
	Connected  bool
	context    *mbContext
	config     printerConfig
	connection *makerbot.Client
}

type printerConnections []*printerConnection

func (pcs *printerConnections) ConnectedPrinters() *[]makerbot.Printer {
	printers := []makerbot.Printer{}
	for _, c := range *pcs {
		if !c.Connected {
			continue
		}

		printers = append(printers, *c.connection.Printer)
	}

	return &printers
}

func (pcs *printerConnections) BySerial(serial string) (conn *printerConnection, ok bool) {
	ok = false

	for _, c := range *pcs {
		if !c.Connected || c.connection.Printer.Serial != serial {
			continue
		}

		conn = c
		ok = true
	}

	return conn, ok
}

func newPrinterConnection(context *mbContext, conf printerConfig) *printerConnection {
	return &printerConnection{Connected: false, context: context, config: conf}
}

func (pc *printerConnection) handleDisconnect() {
	pc.context.Debugln("printerConnection: disconnected! attempting reconnect...")
	pc.Connected = false
	pc.connection = nil

	pc.Connect()
}

func (pc *printerConnection) connectLocal() error {
	pc.context.Debugln("printerConnection: connecting local...")

	err := pc.connection.ConnectLocal(pc.config.IP, pc.config.Port)
	if err != nil {
		return err
	}

	err = pc.connection.AuthenticateWithThingiverse(pc.context.Config.ThingiverseToken, pc.context.Config.ThingiverseUsername)
	if err != nil {
		return err
	}

	pc.Connected = true
	pc.context.Debugln("printerConnection: connected!")
	return nil
}

func (pc *printerConnection) connectRemote() error {
	pc.context.Debugln("printerConnection: connecting remote...")

	err := pc.connection.ConnectRemote(pc.config.ID, pc.context.Config.ThingiverseToken)
	if err != nil {
		return nil
	}

	pc.Connected = true
	pc.context.Debugln("printerConnection: connected!")
	return nil
}

func (pc *printerConnection) Connect() error {
	pc.context.Debugln("printerConnection: Connect() called...")
	cl := makerbot.NewClient()
	cl.HandleDisconnect(pc.handleDisconnect)

	pc.connection = &cl

	if pc.config.ConnectionType == connectionTypeLocal {
		return pc.connectLocal()
	}

	if pc.config.ConnectionType == connectionTypeRemote {
		return pc.connectRemote()
	}

	return errors.New("connection type is wrong")
}
