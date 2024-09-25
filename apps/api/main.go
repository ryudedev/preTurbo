package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var secretKey string = "hardcoded_secret_key"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Golang API!")
	})

	http.HandleFunc("/user", handleUser)
	http.HandleFunc("/file", handleFile)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		return
	}

	// Insecure password hashing
	password := r.URL.Query().Get("password")
	hasher := md5.New()
	hasher.Write([]byte(password))
	hashedPassword := fmt.Sprintf("%x", hasher.Sum(nil))

	fmt.Fprintf(w, "User %s registered with hashed password: %s", username, hashedPassword)
}

func handleFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Missing filename", http.StatusBadRequest)
		return
	}

	// Potential directory traversal vulnerability
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	w.Write(content)
}

func init() {
	// Potential resource leak
	file, _ := os.Open("config.txt")
	defer file.Close()

	// Ignoring errors
	data, _ := ioutil.ReadAll(file)
	fmt.Println(string(data))
}
