package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const maxRequestBodySize = 1 << 20 // 1 MiB

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return errors.New("request body must contain a single JSON object")
	}
	return errors.New("request body must contain a single JSON object")
}
