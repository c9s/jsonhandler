package jsonhandler

import "net/http"
import "fmt"
import "encoding/json"
import "jsondata"

var Debug = true

// General Function to write json response.
func writeJson(w http.ResponseWriter, val interface{}) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	fmt.Fprint(w, string(b))
	return nil
}

var ErrorHandler = func(w http.ResponseWriter) {
	if e := recover(); e != nil {
		if err, ok := e.(error); ok {
			writeJson(w, jsondata.Map{"error": true, "message": err})
		} else {
			writeJson(w, jsondata.Map{"error": true, "message": e})
		}
	}
}

func New(handler func(http.ResponseWriter, *http.Request) interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer ErrorHandler(w)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var resp interface{} = handler(w, r)
		if resp != nil {
			if jsonmap, ok := resp.(*jsondata.Map); ok {
				if Debug {
					fmt.Println(jsonmap)
				}
				fmt.Fprint(w, jsonmap)
			} else if err, ok := resp.(error); ok {
				if Debug {
					fmt.Println(err)
				}
				writeJson(w, jsondata.Map{"error": true, "message": err.Error()})
			} else {
				// for arbitrary data
				if Debug {
					fmt.Println(resp)
				}
				writeJson(w, resp)
			}
		}
	}
}
