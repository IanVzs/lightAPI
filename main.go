package main

import (
	"fmt"
	"net/http"

	"github.com/IanVzs/lightAPI/chat"
	"github.com/IanVzs/lightAPI/flag_parse"
	"github.com/IanVzs/lightAPI/log"
	"github.com/IanVzs/lightAPI/rds"
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
}
func main() {
	defer rds.Client.Close()
	http.HandleFunc("/", hello)
	http.HandleFunc("/chat/polling", chat.ChatAlert)
	log.Logger.Info("server run: " + *flag_parse.Addr)
	err := http.ListenAndServe(*flag_parse.Addr, nil)
	if err != nil {
		log.Logger.Fatal("ListenAndServe: ", err)
	}
}
