package updatelistener

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type requestBody struct {
	Data []string `json:"data"`
}

// StartListener is
func StartListener(updateHandler func(values []string)) {
	log.Println("Starting Update Listener Service...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
		w.WriteHeader(http.StatusOK)

		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		log.Printf("Received Data: %v\n", string(body))

		log.Println("Attempting to Parse received data...")
		var j requestBody
		err := json.Unmarshal(body, &j)
		if err != nil {
			log.Panicln("Failed to parse received data: ", err)
			return
		}
		log.Printf("Parsed Data: %v\n", j)

		fmt.Fprintf(w, "%v", len(body))
		defer updateHandler(j.Data)
	})

	http.HandleFunc("/healthz/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("."))
	})

	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), nil))
}
