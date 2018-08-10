/********************************************************************************
 * Copyright (c) 2018 Edgeworx, Inc.
 *
 * This program and the accompanying materials are made available under the
 * terms of the Eclipse Public License v. 2.0 which is available at
 * http://www.eclipse.org/legal/epl-2.0
 *
 * SPDX-License-Identifier: EPL-2.0
 ********************************************************************************/

package cn

import (
	"crypto/tls"
	"time"
)

type PublicConnection struct {
	*ComSatConn
	containerConn *ContainerConn
}

func newPublicConnection(id int,
	address, remoteAddress, passcode string,
	hbInterval, hbThreshold time.Duration,
	tlsConfig *tls.Config) *PublicConnection {
	return &PublicConnection{
		ComSatConn:    newConn(id, address, passcode, hbInterval, hbThreshold, tlsConfig),
		containerConn: newContainerConn(id, remoteAddress),
	}
}

func (p *PublicConnection) Connect() {
	reconnect := make(chan byte)
	defer func() {
		close(reconnect)
	}()
	go p.ComSatConn.Connect()
	go p.readConnection(reconnect)

	for {
		select {
		case <-reconnect:
			go p.readConnection(reconnect)
		}
	}
}

func (p *PublicConnection) Disconnect() {
	p.ComSatConn.Disconnect()
	p.containerConn.Disconnect()
}

func (p *PublicConnection) proxy(done chan byte) {
	for {
		select {
		case data, ok := <-p.containerConn.out:
			if !ok {
				done <- 0
				return
			}
			p.in <- data
		}
	}
}

func (p *PublicConnection) readConnection(done chan byte) {
	for {
		select {
		case <-done:
			return
		case data := <-p.out:
			if !p.containerConn.isConnected {
				if err := p.containerConn.Connect(); err != nil {
					logger.Printf("[ PublicConnection #%d ] Error when connecting to container: %s\n",
						p.id, err.Error())
					continue
				} else {
					go p.containerConn.Start()
					go p.proxy(done)
				}
			}
			p.containerConn.in <- data
		}
	}
}
