// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// SMPP client for the command line.
//
// We bind to the SMSC as a transmitter, therefore can do SubmitSM
// (send Short Message) or QuerySM (query for message status). The
// latter may not be available depending on the SMSC.
package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/fiorix/go-smpp/smpp"
)

// Version of smppcli.
var Version = "tip"

// Author of smppcli.
var Author = "go-smpp authors"

func main() {
	app := cli.NewApp()
	app.Name = "smpp_client"
	app.Usage = "SMPP client for SMSC"
	app.Version = Version
	app.Author = Author

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Value: "localhost:6200",
			Usage: "Set SMPP server host:port",
		},
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "Set SMPP username",
		},
		cli.StringFlag{
			Name:  "passwd",
			Value: "",
			Usage: "Set SMPP password",
		},
		cli.BoolFlag{
			Name:  "tls",
			Usage: "Use client TLS connection",
		},
		cli.BoolFlag{
			Name:  "precaire",
			Usage: "Accept invalid TLS certificate",
		},
	}

	app.Commands = []cli.Command{
		runClient,
	}
	app.Run(os.Args)
}

var runClient = cli.Command{
	Name:  "runClient",
	Usage: "start SMPP Client",
	Action: func(c *cli.Context) {
		log.Println("Connecting...")
		tx := newTransmitter(c)
		defer tx.Close()
		log.Println("Connected to", tx.Addr)
		for {
			time.Sleep(100 * time.Millisecond)
			fmt.PrintLn("sleep... ")
		}
	},
}

func newTransmitter(c *cli.Context) *smpp.Transmitter {
	tx := &smpp.Transmitter{
		Addr:   c.GlobalString("addr"),
		User:   os.Getenv("SMPP_USER"),
		Passwd: os.Getenv("SMPP_PASSWD"),
	}
	if s := c.GlobalString("user"); s != "" {
		tx.User = s
	}
	if s := c.GlobalString("passwd"); s != "" {
		tx.Passwd = s
	}
	if c.GlobalBool("tls") {
		host, _, _ := net.SplitHostPort(tx.Addr)
		tx.TLS = &tls.Config{
			ServerName: host,
		}
		if c.GlobalBool("precaire") {
			tx.TLS.InsecureSkipVerify = true
		}
	}
	conn := <-tx.Bind()
	switch conn.Status() {
	case smpp.Connected:
	default:
		log.Fatalln("Connection failed:", conn.Error())
	}
	return tx
}
