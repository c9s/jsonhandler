package jsonhandler

import "net/http"
import "fmt"
import "encoding/json"
import "jsondata"

// General Function to write json response.
func writeJson(w http.ResponseWriter, val interface{}) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	fmt.Fprint(w, string(b))
	return nil
}

func New(handler func(http.ResponseWriter, *http.Request) interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				writeJson(w, jsondata.Map{"error": true, "message": err})
			}
		}()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var resp interface{} = handler(w, r)
		if resp != nil {
			if jsonmap, ok := resp.(*jsondata.Map); ok {
				fmt.Fprint(w, jsonmap)
			} else {
				// for arbitrary data
				writeJson(w, resp)
			}
		}
	}
}
