//
//package http_proxy
//
//import (
//	"bufio"
//	"context"
//	"fmt"
//	"github.com/libp2p/go-libp2p"
//	"github.com/libp2p/go-libp2p-core/host"
//	"github.com/libp2p/go-libp2p-core/network"
//	"github.com/libp2p/go-libp2p-core/peer"
//	"github.com/libp2p/go-libp2p-core/peerstore"
//	ma "github.com/multiformats/go-multiaddr"
//	manet "github.com/multiformats/go-multiaddr-net"
//	"log"
//	"net/http"
//	"strings"
//)
//
//const Protocol = "/proxy-example/0.0.1"
//
//func makeRandomHOst(port int) host.Host {
//	host, err := libp2p.New(context.Background(), libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)))
//	if err != nil {
//		log.Fatalln(err)
//	}
//	return host
//}
//
//type ProxyService struct {
//	host      host.Host
//	dest      peer.ID
//	proxyAddr ma.Multiaddr
//}
//
//func NewProxyService(h host.Host, proxyAddr ma.Multiaddr, dest peer.ID) *ProxyService {
//	h.SetStreamHandler(Protocol, streamHandler)
//	fmt.Println("Proxy server is ready")
//	fmt.Println("libp2p-peer addresses:")
//	for _, a := range h.Addrs() {
//		fmt.Printf("%s/ipfs/%s\n", a, peer.IDB58Encode(h.ID()))
//	}
//
//	return &ProxyService{
//		host:      h,
//		dest:      dest,
//		proxyAddr: proxyAddr,
//	}
//}
//
//func streamHandler(stream network.Stream) {
//	defer stream.Close()
//
//	buf := bufio.NewReader(stream)
//	req, err := http.ReadRequest(buf)
//	if err != nil {
//		stream.Reset()
//		log.Println(err)
//		return
//	}
//	defer req.Body.Close()
//
//	req.URL.Scheme = "http"
//	hp := strings.Split(req.Host, ":")
//	if len(hp) > 1 && hp[1] == "443" {
//		req.URL.Scheme = "https"
//	} else {
//		req.URL.Scheme = "http"
//	}
//	req.URL.Host = req.Host
//	outreq := new(http.Request)
//	*outreq = *req
//
//	fmt.Printf("Making request to %s\n", req.URL)
//	resp, err := http.DefaultTransport.RoundTrip(outreq)
//	if err != nil {
//		stream.Reset()
//		log.Println(err)
//		return
//	}
//
//	resp.Write(stream)
//}
//
//func (p *ProxyService) Serve() {
//	_, serveArgs, _ := manet.DialArgs(p.proxyAddr)
//	fmt.Println("proxy listening on ", serveArgs)
//	if p.dest != "" {
//		http.ListenAndServe(serveArgs, p)
//	}
//}
//
