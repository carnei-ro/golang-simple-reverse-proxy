package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Proxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	r.Host = p.target.Host
	w.Header().Set("X-Proxied-By", "go-reverse-proxy")
	p.proxy.ServeHTTP(w, r)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		value = fallback
	}
	return value
}

func logBanner(upstream string, port string) {
	log.Printf("Starting reverse proxy on port %s\n", port)
	log.Printf("Proxying requests to %s\n", upstream)
}

func main() {
	// Replace 'target' with the URL of the server you want to proxy to
	upstream := getEnv("UPSTREAM", "https://httpbin.org")
	target, err := url.Parse(upstream)
	if err != nil {
		panic(err)
	}

	// Create a new ReverseProxy instance
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Configure the reverse proxy to use HTTPS
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create a new Proxy instance
	p := &Proxy{target: target, proxy: proxy}

	// Start the HTTP server and register the Proxy instance as the handler
	port := getEnv("PORT", "8080")
	logBanner(upstream, port)
	err = http.ListenAndServe(":"+port, p)
	if err != nil {
		panic(err)
	}
}
