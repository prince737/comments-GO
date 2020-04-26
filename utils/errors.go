package utils

import (
	"fmt"
	"net/http"
)

//InternalServerError returns internal server error
func InternalServerError(w http.ResponseWriter) {
	fmt.Println("here")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}
