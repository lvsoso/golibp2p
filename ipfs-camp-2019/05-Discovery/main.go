package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	mplex "github.com/libp2p/go-libp2p-mplex"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)

	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)

	security := libp2p.Security(secio.ID, secio.New)

	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0",
		"/ip4/0.0.0.0/tcp/0/ws",
	)

	host, err := libp2p.New(
		ctx,
		transports,
		listenAddrs,
		muxers,
		security,
	)
	if err != nil {
		panic(err)
	}

	host.SetStreamHandler(chatProtocol, chatHandler)

	for _, addr := range host.Addrs() {
		fmt.Println("Listening on", addr)
	}
	fmt.Println("ID: ", host.ID().Pretty())

	//input := bufio.NewScanner(os.Stdin)
	//inputArray := []string{}
	//fmt.Printf("Please type in something:\n")
	//for input.Scan() {
	//	line := input.Text()
	//	if line == "end" {
	//		break
	//	}
	//	inputArray = append(inputArray, line)
	//}
	//targetAddrInput := inputArray[0]
	////targetAddr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/63785/p2p/QmWjz6xb8v9K4KnYEwP5Yk75k5mMBCehzWFLCvvQpYxF3d")
	//targetAddr, err := multiaddr.NewMultiaddr(targetAddrInput)
	//if err != nil {
	//	panic(err)
	//}
	//
	//targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = host.Connect(ctx, *targetInfo)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("Connected to", targetInfo.ID)

	donec := make(chan struct{}, 1)
	go chatInputLoop(ctx, host, donec)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	select {
	case <-stop:
		host.Close()
		os.Exit(0)
	case <-donec:
		host.Close()
	}
}
