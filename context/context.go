package context

import (
	"fmt"
	"log"
	"net/http"
)

type Store interface {
	Fetch() string
	Cancel()
}

func Server(store Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		ctx := request.Context()
		data := make(chan string, 1)

		go func() {
			data <- store.Fetch()
		}()

		select {
		case d := <-data:
			_, err := fmt.Fprint(writer, d)
			if err != nil {
				log.Fatalf("error printing to writer: %+v", err)
			}
		case <-ctx.Done():
			store.Cancel()
		}
	}
}
