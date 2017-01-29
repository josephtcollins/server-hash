package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"encoding/json"
	"crypto/sha512"
	"time"
)

var lis net.Listener

type Handler struct {
	w sync.WaitGroup
	*http.ServeMux
}

func main() {
	var err error
	lis, err = net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	h := &Handler{ServeMux: http.NewServeMux()}
	h.ServeMux.HandleFunc("/", hashPassword)
	h.ServeMux.HandleFunc("/shutdown", close)

	fmt.Println("Listening on:", lis.Addr())
	http.Serve(lis, h)

	h.w.Wait()
}

// at a new request increment wait group, decrement once scope closes
// call default ServeHttp, call flush to make sure responses are sent
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.w.Add(1)
	defer h.w.Done()
	h.ServeMux.ServeHTTP(w, r)
	w.(http.Flusher).Flush()
}

// '/'
func hashPassword(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)

	pw := r.FormValue("password")
	if len(pw) == 0 {
		panic("No password received")
	}

	h := sha512.New()
	h.Write([]byte(pw))

	// byte[] encodes as base64 string
	json.NewEncoder(w).Encode(h.Sum(nil))
}

// '/shutdown'
func close(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Finishing requests...")
	lis.Close()
}
