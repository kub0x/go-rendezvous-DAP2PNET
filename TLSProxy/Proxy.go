package TLSProxy

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

var (
	TLSProxyErrReadCA     = errors.New("cannot read ca certificate")
	TLSProxyErrPoolAppend = errors.New("cannot add cert to pool")
	TLSProxyErrLoadChain  = errors.New("cannot load certificate chain")
)

type TLSProxy struct {
	CACertPath      string
	ServerChainPath string
	ServerKeyPath   string
	HostRedirectURL string
}

func NewTLSProxy() *TLSProxy {
	tlsProxy := &TLSProxy{
		CACertPath:      "./certs/ca.pem",
		ServerChainPath: "./certs/rendezvous.dap2p.net.pem",
		ServerKeyPath:   "./certs/rendezvous.dap2p.net.key",
		HostRedirectURL: "https://rendezvous.dap2p.net:6668", // internal gin https server
	}

	return tlsProxy
}

func (tlsProxy *TLSProxy) Listen() error {

	caBytes, err := ioutil.ReadFile(tlsProxy.CACertPath)
	if err != nil {
		return TLSProxyErrReadCA
	}

	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(caBytes); !ok {
		return TLSProxyErrPoolAppend
	}

	tlsCertChain, err := tls.LoadX509KeyPair(tlsProxy.ServerChainPath, tlsProxy.ServerKeyPath)
	if err != nil {
		return TLSProxyErrLoadChain
	}

	tlsConfig := &tls.Config{
		// Only accept client certificates signed by our PKI
		ClientAuth: tls.RequireAndVerifyClientCert,
		// Must validate client cert chain against our CA
		ClientCAs: clientCertPool,
		// Supported suites
		CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256, tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA, tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
		MinVersion:   tls.VersionTLS11,
		// Avoid tls3 for debugging purposes
		MaxVersion: tls.VersionTLS12,
		// Send the certificate chain
		Certificates: []tls.Certificate{tlsCertChain},
	}

	httpServer := &http.Server{
		Addr:      ":6667",
		TLSConfig: tlsConfig,
	}

	// func that acts as a gateway

	http.HandleFunc("/", tlsProxy.gateWay)

	return httpServer.ListenAndServeTLS(tlsProxy.ServerChainPath, tlsProxy.ServerKeyPath)
}

func (tlsProxy *TLSProxy) gateWay(w http.ResponseWriter, req *http.Request) {
	// trust dap2pnet CA as gin uses rendezvous cert too
	certPool := x509.NewCertPool()

	caBytes, err := ioutil.ReadFile("./certs/ca.pem")
	if err != nil {
		log.Println(err.Error())
	}

	certPool.AppendCertsFromPEM(caBytes)
	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// build new redirect url with same url path

	httpClient := &http.Client{Transport: tr}
	req.RequestURI = ""
	req.URL, _ = url.Parse(tlsProxy.HostRedirectURL + req.URL.Path)

	req.Header.Add("Authorization", req.TLS.PeerCertificates[0].Subject.CommonName) // identify the peer that requests a resource

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	req.Header.Add("X-Forwarded-For", req.RemoteAddr)
	log.Printf("Redirecting CN=%v with IP=%v to %v\n", req.TLS.PeerCertificates[0].Subject.CommonName, ip, req.URL.Path)

	resp, err := httpClient.Do(req)
	if err != nil { // 500 when we cannot connect to internal gin https server
		log.Println("failed to initiate internal http request via tlsproxy: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil { // 500 when we cannot retrieve a valid response body from internal gin
		log.Println("failed to retrieve the http response via tlsproxy: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	// assing internal status code and response body to proxy response

	w.WriteHeader(resp.StatusCode)
	fmt.Fprintf(w, "%v", string(body))

}
