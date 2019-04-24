// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

import (
	"fmt"
	"net"
)

type socket struct {
	sink       ISocketSink
	addr       string
	listener   net.Listener
	serializer ISerializer
	peers      []*peer
}

func (s *socket) Listen() {
	var err error
	s.listener, err = net.Listen("tcp", s.addr)
	if nil != err {
		fmt.Println(err)
		return
	}

	for {
		conn, err := s.listener.Accept()
		if nil != err {
			fmt.Println(err)
			continue
		}

		peer := newPeer(conn, s.serializer, s.sink, false)
		s.peers = append(s.peers, peer)

		if nil != s.sink {
			s.sink.OnConnected(peer)
		}

		peer.run()
	}
}

func (s *socket) Dial() {
	if nil == s.sink {
		panic("Please call Init first!")
	}

	conn, err := net.Dial("tcp", s.addr)
	if nil != err {
		fmt.Println(err)
		return
	}

	peer := newPeer(conn, s.serializer, s.sink, true)
	s.peers = append(s.peers, peer)

	if nil != s.sink {
		s.sink.OnConnected(peer)
	}

	peer.run()
}

func (s *socket) Close() {
	for _, peer := range s.peers {
		peer.close()
	}
	s.peers = nil

	s.listener.Close()
	s.listener = nil
}

func (s *socket) Send(obj interface{}) {
	for _, peer := range s.peers {
		peer.Send(obj)
	}
}
