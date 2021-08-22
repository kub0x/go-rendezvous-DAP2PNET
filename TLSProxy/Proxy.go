package TLSProxy

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
}

func NewTLSProxy() *TLSProxy {
	tlsProxy := &TLSProxy{
		CACertPath:      "./certs/ca.pem",
		ServerChainPath: "./certs/rendezvous.dap2p.net.pem",
		ServerKeyPath:   "./certs/rendezvous.dap2p.net.key",
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
		// Only accept client certificate signed by our PKI
		ClientAuth: tls.RequireAndVerifyClientCert,
		// Must validate client cert chain against our CA
		ClientCAs: clientCertPool,
		// Supported suites
		CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256, tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA, tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
		// Force it server side
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS11,
		// Avoid tls3 for debugging purposes
		MaxVersion: tls.VersionTLS12,
		// Send the certificate chain
		Certificates: []tls.Certificate{tlsCertChain},
	}

	httpServer := &http.Server{
		Addr:      ":6667",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/", HelloUser)

	return httpServer.ListenAndServeTLS(tlsProxy.ServerChainPath, tlsProxy.ServerKeyPath)
}

func HelloUser(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello %v! you've been signed by %v \n", req.TLS.PeerCertificates[0].Subject.CommonName, req.TLS.PeerCertificates[1].Subject.CommonName)
}
