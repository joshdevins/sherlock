package main

import (
	"flag"
	"github.com/gorilla/pat"
	"log"
	"net/http"
)

// A simple, stand-alone, fingerprint index. This is meant as a small-scale
// evaluation tool for testing various search and indexing parameters and
// algorithms of the Philips fingerprinter [1].
//
// [1] J. Haitsma and A. Kalker, “A Highly Robust Audio Fingerprinting System,”
// in _Proc. International Symposium on Music Information Retrieval (ISMIR)_,
// 2002.
//
func main() {
	serverAddr := flag.String("server.addr", ":8080", "HTTP server listen address")
	flag.Parse()

	// routes
	r := pat.New()
	r.Get("/-/stats", statsHandler())

	// serve
	log.Printf("Listening on: %s", *serverAddr)
	log.Fatal(http.ListenAndServe(*serverAddr, r))
}
