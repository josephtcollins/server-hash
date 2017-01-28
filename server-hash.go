package main

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "time"
)

type Password struct {
	Password string
}

func handlePosts(w http.ResponseWriter, r *http.Request) {
  postPassword(w, r)
}

func postPassword(w http.ResponseWriter, r *http.Request) {
  fmt.Println(r.FormValue("password"))
	json.NewEncoder(w).Encode(getSHA256(r.FormValue("password")))
}

func handleRequests() {
	http.HandleFunc("/", postPassword)
	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func getSHA256(pw string) string {
	h := sha512.New()
	h.Write([]byte(pw))
	sha1_hash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return string(sha1_hash)
}

func main() {
	handleRequests()
}
