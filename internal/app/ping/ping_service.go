package ping

import (
	"fmt"
	"log"
	"net/http"
)

func Serve() {
	http.HandleFunc("/ping", ping)
	go func() {
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			log.Panic(err)
		}
	}()
}

func ping(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprint(writer, "OK")
	if err != nil {
		log.Print(err)
	}
}
