package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/pin/tftp/v3"
)

var nextBootOption = "0"

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprint(w, nextBootOption)
		return
	}

	if req.Method == http.MethodPost {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("[HTTP] Failed to set next boot option: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(b) == 0 {
			log.Println("[HTTP] Failed to set next boot option: empty body")
			http.Error(w, "Empty body", http.StatusBadRequest)
			return
		}
		nextBootOption = string(b)
		log.Printf("[HTTP] Set next boot option to '%s'\n", nextBootOption)
		fmt.Fprint(w, "OK")
		return
	}
}

func handleTFTP(filename string, rf io.ReaderFrom) error {
	ip := rf.(tftp.OutgoingTransfer).RemoteAddr().IP.String()
	log.Printf("[TFTP] %s is booting '%s'\n", ip, nextBootOption)
	rf.ReadFrom(strings.NewReader(fmt.Sprintf("set default=\"%s\"\n", nextBootOption)))
	return nil
}

func main() {
	log.Println("Starting...")
	log.Printf("Next boot option is '%s'\n", nextBootOption)
	log.Println("[HTTP] Binding to :80")
	log.Println("[TFTP] Binding to :69")

	http.HandleFunc("/{$}", handleHTTP)
	go func() {
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			log.Fatalf("[HTTP] Failed to setup server: %v\n", err)
		}
	}()

	s := tftp.NewServer(handleTFTP, nil)
	err := s.ListenAndServe(":69")
	if err != nil {
		log.Fatalf("[TFTP] Failed to setup server: %v\n", err)
	}
}
