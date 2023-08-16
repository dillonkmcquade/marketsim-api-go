package handlers

import (
	"net/http"
)

func HealthCheck(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("OK!\n"))
}
