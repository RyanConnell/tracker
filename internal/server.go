package server

import (
	"fmt"
	"net/http"
)

func Serve(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
