package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	"sync"
	"time"
)

type mbContext struct {
	Printers *printerConnections
	Config   *config
}

func (ctx *mbContext) Debugln(v ...interface{}) {
	if ctx.Config.Debug {
		log.Println(v...)
	}
}

func (ctx *mbContext) Debugf(fmt string, v ...interface{}) {
	if ctx.Config.Debug {
		log.Printf(fmt, v...)
	}
}

func main() {
	var defaultConfig string

	usr, err := user.Current()
	if err != nil {
		defaultConfig = "/etc/makerbotd/config.json"
	} else {
		defaultConfig = usr.HomeDir + "/.makerbotd/config.json"
	}

	confPath := flag.String("config", defaultConfig, "the path to the makerbotd config file")
	forceListen := flag.Bool("force-listen", false, "force listen on unix socket if it is in use")
	flag.Parse()

	conf, err := loadConfig(*confPath)
	if err != nil {
		panic(err)
	}

	ctx := mbContext{Config: conf}

	printers := printerConnections{}

	// Set up printer connections
	for _, pc := range conf.Printers {
		conn := newPrinterConnection(&ctx, pc)
		go conn.Connect()
		printers = append(printers, conn)
	}

	ctx.Printers = &printers

	router := getRouter(&ctx)

	server := http.Server{
		Handler:     router,
		ReadTimeout: 5 * time.Minute,
	}

	var wg sync.WaitGroup

	if conf.ListenSocket {
		if *forceListen {
			os.Remove(conf.ListenSocketPath)
		}

		wg.Add(1)
		go func() {
			sock, err := net.Listen("unix", conf.ListenSocketPath)
			if err != nil {
				panic(err)
			}

			defer sock.Close()
			defer os.Remove(conf.ListenSocketPath)

			log.Printf("HTTP server listening on UNIX domain socket (force=%v): %s", *forceListen, conf.ListenSocketPath)

			err = server.Serve(sock)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}()
	}

	if conf.ListenTCP {
		wg.Add(1)
		go func() {
			conn, err := net.Listen("tcp", conf.ListenTCPAddress)
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			log.Printf("HTTP server listening on TCP address: %s", conf.ListenTCPAddress)

			err = server.Serve(conn)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
