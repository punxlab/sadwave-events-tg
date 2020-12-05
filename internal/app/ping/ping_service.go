package ping

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Serve() {
	startDate := time.Now().Format("2006-01-02 15:04:05")
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := fmt.Fprint(writer, startDate)
		if err != nil {
			log.Printf("ping: %v", err)
		}
	})
	go func() {
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			log.Printf("serve ping: %v", err)
		}
	}()
}
