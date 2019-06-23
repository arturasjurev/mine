package api

import (
	"encoding/json"
	"net/http"
)

func SegmentString(mux map[string]string, seg string) string {
	v, _ := mux[seg]
	return string(v)
}

func BindJSON(r *http.Request, data interface{}) (bool, error) {
	err := json.NewDecoder(r.Body).Decode(&data)
	return err == nil, err
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	output, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(output)
	return
}
