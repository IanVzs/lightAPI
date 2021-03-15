package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/IanVzs/lightAPI/chat"
	"github.com/IanVzs/lightAPI/log"
)

func hello(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// http.ServeFile(w, r, "home.html")
	fmt.Fprintf(w, "hello\n")
}
func init() {
	flag.Parse()
}
func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/chat/polling", chat.ChatAlert)
	log.Logger.Info("server run: " + *log.Addr)
	err := http.ListenAndServe(*log.Addr, nil)
	if err != nil {
		log.Logger.Fatal("ListenAndServe: ", err)
	}
}
