package respond

import (
	"encoding/json"
	"net/http"
)

func Json(w http.ResponseWriter, statusCode int, v interface{}) {
	var err error
	b := []byte(``)
	if v != nil {
		b, err = json.Marshal(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
