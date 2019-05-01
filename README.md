# makerbotd

makerbotd is a daemon that manages connections to MakerBot 3D printers. With it, you can control your printers with simple HTTP API calls.

**This software is mid-development. It is not stable yet.** That said, feel free to play around with it!

## Examples

Here are some examples of what makerbotd can do.

Send `GET /api/v1/printers`, get a list of connected printers back:

```json
{
    "result": [
        {
            "machine_type": "platypus",
            "vid": 9153,
            "ip": "10.65.1.99",
            "pid": 5,
            "api_version": "1.9.0",
            "iserial": "23C100053C7059018291",
            "ssl_port": "12309",
            "machine_name": "Kevin",
            "motor_driver_version": "4.6",
            "bot_type": "replicator_5",
            "port": "9999",
            "firmware_version": {
                "major": 2,
                "minor": 6,
                "bugfix": 1,
                "build": 724
            }
        }
    ],
    "error": null
}
```

Send a print file to the printer with `POST /api/v1/printers/23C100053C7059018291/prints`.

Grab a snapshot from the printer's camera with `GET /api/v1/printers/23C100053C7059018291/snapshot.jpg`:

![](https://user-images.githubusercontent.com/2646487/57029732-71b08a00-6bf7-11e9-90ad-3f3339c0d181.png)

## Setup

Since makerbotd is a daemon, it expects to be run as a background process. If you use systemd, you can install the `makerbotd` service with `make install-systemd` (may need sudo). This command will:

- Build `makerbotd`
- Move it to `/usr/local/bin/makerbotd`
- Copy the `makerbotd.service` file to the `/etc/systemd/system` directory

It will not enable or start makerbotd after you run this command. You should run `systemctl enable makerbotd` or `systemctl start makerbotd` after this if you wish.

After `makerbotd` starts for the first time, it will create a config file in the directory specified. With the sample `makerbotd.service` in this directory, it will create it at `/etc/makerbotd/config.json`. You should edit this config file to suit your needs and then `systemctl restart makerbotd`. (Reloading is not supported yet.)

## Configuration

Since this project is in a pretty early state, the schema of the config file may change from time to time. Here it is as of right now:

```golang
type config struct {
	Debug               bool            // Debug makes output more verbose
	ThingiverseUsername string          // ThingiverseUsername defines the username of the authenticated Thingiverse account
	ThingiverseToken    string          // ThingiverseToken defines the auth token of the authenticated Thingiverse account
	ListenSocket        bool            // ListenSocket defines whether or not makerbotd will listen on a unix domain socket
	ListenSocketPath    string          // ListenSocketPath defines the unix domain socket to listen on if ListenSocket is true
	ListenTCP           bool            // ListenTCP defines whether or not makerbotd will listen on a TCP port
	ListenTCPAddress    string          // ListenTCPPort defines the TCP port to listen on if ListenTCP is true
	AutoAddPrinters     bool            // AutoAddPrinters defines whether or not printers should automatically be added from the authenticated Thingiverse account (DOES NOTHING RIGHT NOW)
	Printers            []printerConfig // Printers is the list of MakerBot printers that will automatically be connected when makerbotd starts
}

type printerConfig struct {
	ConnectionType string // ConnectionType should be either "local" or "remote". "local" = direct connect via IP, "remote" = remotely connect via MakerBot Reflector service.
	ID             string // ID should be provided if the connection type is "remote". This is the ID of the printer as returned by MakerBot Reflector. It is usually the serial number.
	IP             string // IP should be provided if the connection type is "local"
	Port           string // Port should be provided if the connection type is "port"
}
```

A sane default config is written on first start that connects to no printers and listens at `/var/run/makerbot.socket`.

## API

Check out `api_v1.go` to see what API routes are available. I don't want to write proper documentation just yet since it will likely change pretty often while this is in development. I recommend using Postman for testing the API out -- it has pretty good UNIX domain socket support. For example: `unix:///var/run/makerbot.socket:/api/v1/printers`

## License

TBD