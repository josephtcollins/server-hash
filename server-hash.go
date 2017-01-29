package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"encoding/json"
	"crypto/sha512"
	"encoding/base64"
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

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.w.Add(1)
	defer h.w.Done()
	h.ServeMux.ServeHTTP(w, r)
	w.(http.Flusher).Flush()
}

// '/'
func hashPassword(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	json.NewEncoder(w).Encode(getSHA512(r.FormValue("password")))
}

func getSHA512(pw string) string {
	h := sha512.New()
	h.Write([]byte(pw))
	sha1_hash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return string(sha1_hash)
}

// '/shutdown'
func close(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Finishing requests...")
	lis.Close()
}
