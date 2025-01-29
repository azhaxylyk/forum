package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"
)

func StartServer(mux http.Handler) {
	certPath := os.Getenv("SSL_CERT_PATH")
	keyPath := os.Getenv("SSL_KEY_PATH")
	if certPath == "" || keyPath == "" {
		log.Fatal("Certificate paths not set in environment variables")
	}
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
			},
		},
	}

	log.Println("Server started on https://localhost:8080")
	err := server.ListenAndServeTLS(certPath, keyPath)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
