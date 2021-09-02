package main

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/krkhan/go-std-server/router"
	"github.com/krkhan/go-std-server/store"
)

const (
	PasswordKey        = "password"
	Sha512DelaySeconds = 5
)

func handleError(msg string, w http.ResponseWriter) {
	log.Printf(msg)
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, msg)
}

func postHash(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling POST /hash")

	err := r.ParseForm()
	if err != nil {
		msg := fmt.Sprintf("Error parsing request body: %s", err)
		handleError(msg, w)
		return
	}

	passwords, ok := r.PostForm[PasswordKey]
	if !ok {
		msg := fmt.Sprintf("Error parsing field: %s", PasswordKey)
		handleError(msg, w)
		return
	}

	if len(passwords) > 1 {
		msg := fmt.Sprintf("Multiple (%d) passwords provided", len(passwords))
		handleError(msg, w)
		return
	}

	digest := sha512.Sum512([]byte(passwords[0]))
	reqStore := r.Context().Value(store.Sha512DigestStoreContextKey{}).(*store.Sha512DigestStore)
	newKey := reqStore.AddDigest(digest)

	io.WriteString(w, fmt.Sprintf("%d", newKey))
}

func getHash(w http.ResponseWriter, r *http.Request) {
	param := router.GetParam(r, 0)
	key, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("Invalid key: %s", param)
		handleError(msg, w)
		return
	}

	log.Printf("Handling GET /hash key=%d", key)

	reqStore := r.Context().Value(store.Sha512DigestStoreContextKey{}).(*store.Sha512DigestStore)
	digest, ok := reqStore.GetDigest(key)
	if !ok {
		msg := fmt.Sprintf("Key not found: %d", key)
		handleError(msg, w)
		return
	}
	b64Encoded := base64.StdEncoding.EncodeToString(digest[:])

	io.WriteString(w, b64Encoded)
}

func getStats(w http.ResponseWriter, r *http.Request) {
	log.Print("Handling GET /stats")
	routes := r.Context().Value(RoutesContextKey{}).([]router.Route)
	for _, route := range routes {
		if route.Name == "postHash" {
			route.Stats.StatsLock.RLock()
			defer route.Stats.StatsLock.RUnlock()
			totalRequests := route.Stats.TotalRequests
			averageTime := route.Stats.TotalTime / totalRequests
			io.WriteString(w, fmt.Sprintf(`{"total": "%d", "average": "%d"}`, totalRequests, averageTime))
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
}

type ServerHandler struct {
	Routes []router.Route
	Store  *store.Sha512DigestStore
}

type RoutesContextKey struct{}

func (h ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), store.Sha512DigestStoreContextKey{}, h.Store)
	ctx = context.WithValue(ctx, RoutesContextKey{}, h.Routes)
	router.Serve(h.Routes, w, r.WithContext(ctx))
}

func startHttpServer(wg *sync.WaitGroup, addr string) *http.Server {
	digestsMap := make(map[uint64][sha512.Size]byte)
	handler := ServerHandler{
		Routes: []router.Route{
			router.NewRoute("postHash", http.MethodPost, "/hash", postHash),
			router.NewRoute("getHash", http.MethodGet, "/hash/([0-9]+)", getHash),
			router.NewRoute("getStats", http.MethodGet, "/stats", getStats),
		},
		Store: &store.Sha512DigestStore{
			Counter:      0,
			Digests:      &digestsMap,
			DigestsLock:  &sync.RWMutex{},
			DelaySeconds: Sha512DelaySeconds,
		},
	}
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		defer wg.Done()

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed with: %v", err)
		}
	}()

	return server
}

func main() {
	serverAddr := ":8080"
	if len(os.Args) > 1 {
		serverAddr = os.Args[1]
	}

	log.Printf("Launching HTTP server on %s", serverAddr)

	serverExited := &sync.WaitGroup{}
	serverExited.Add(1)
	server := startHttpServer(serverExited, serverAddr)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("Received signal '%s', shutting down HTTP server", sig)

		if err := server.Shutdown(context.TODO()); err != nil {
			panic(err)
		}
	}()

	serverExited.Wait()

	log.Printf("HTTP server terminated successfully")
}
