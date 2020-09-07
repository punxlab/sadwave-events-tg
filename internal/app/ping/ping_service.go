package ping

import (
	"fmt"
	"log"
	"net/http"
)

func Serve() {
	http.HandleFunc("/", ping)
	go func() {
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			log.Printf("serve ping: %v", err)
		}
	}()
}

func ping(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprint(writer, "OK")
	if err != nil {
		log.Printf("ping: %v", err)
	}
}
