# HouseMon 0.9.x

Real-time home monitoring and automation.

> **This is the DEVELOPMENT branch of HouseMon. See README-0.8 for older info.**

The server side of HouseMon is undergoing some radical changes, caused by  
replacing Node.js with an MQTT + LevelDB + Lua server, written in Go.  
This server core architecture is also used in Tosqa and is called "JeeBus".

![HM-0.9-diagram.png](HM-0.9-diagram.png)

# Installation

Right now, you still need all the pieces of 0.8 for development (i.e. Node.js).

In addition, you need to install Go, see <http://golang.org/doc/install>.

Once Go is installed, you need JeeBus, see <https://github.com/jcw/jeebus>.

To launch this early HouseMon 0.9 setup, perform the following steps:

* make sure you're in the HouseMon root directory, next to this README
* launch JeeBus as follows, and keep it running: `jb run :3334`
* in a separate terminal window, launch Node.js as follows: `node .`
* now go to the HouseMon pages by browsing to `localhost:3333` (not a typo)

The definitive way to launch HouseMon will be simplified once the dust settles.

# Details

In this initial phase, JeeBus has been tied into HouseMon as follows:

* there's a `main.go` file which defines additional HM-specific sub-commands
* the `app/jeebus` module connects a _second_ websocket to port 3334, i.e. JB
* the `./logger` and `./storage` folders used from Node.js have been renamed

This allows keeping the Primus Live update mechanism and automatic compilation
of CoffeeScript, Jade, and Stylus. As things get moved over, and all Primus-
based RPC calls are replaced by calls over the JeeBus websocket, the websocket
on port 3333 will end up doing less and less, until it only handles reload and
CSS update triggers whenever source file changes are detected by Primus Live.

HouseMon is available at `localhost:3333` as long as both the JeeBus server
(`jb run :3334`) and the Primus Live server (`node .`) are kept running. For
additional services specific to HouseMon, more processes need to be started,
such as `housemon logger` (i.e. `go run main.go logger` in developer mode).

Note that the old data is now in `./storage-old`. There's no data sharing with
the new JeeBus data (in `./storage`). The *-old stuff will soon be dropped.

# Documentation

There is a _little_ documentation about HouseMon in the `README-0.8.md` file.

# License

MIT
