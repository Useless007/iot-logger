package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	port := "5000"
	logDir := "dht_logs"
	os.MkdirAll(logDir, 0755)

	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST", http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		temp := r.FormValue("temp")
		hum := r.FormValue("hum")
		if temp == "" || hum == "" {
			http.Error(w, "Missing temp/hum", http.StatusBadRequest)
			return
		}

		ts := time.Now().Format("2006-01-02 15:04:05")
		filename := filepath.Join(logDir, time.Now().Format("2006-01-02")+".csv")
		f, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString(fmt.Sprintf("%s,%s,%s\n", ts, temp, hum))
		log.Printf("[%s] %s°C, %s%%", ts, temp, hum)
		fmt.Fprintf(w, "Logged\n")
	})

	log.Println("Logger running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
