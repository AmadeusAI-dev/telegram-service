package server

import "net/http"

func Run(ch chan<- error, url string, handler http.Handler) {
	go func() {
		err := http.ListenAndServe(url, handler)
		if err != nil {
			ch <- err
		}
	}()
}
