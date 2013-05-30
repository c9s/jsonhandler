package jsonhandler

import (
	"encoding/json"
	"fmt"
	"jsondata"
	"log"
	"net/http"
	"runtime/debug"
)

const Padding = "  "

// General Function to write json response.
func writeJson(w http.ResponseWriter, val interface{}) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	fmt.Fprint(w, string(b))
	return nil
}

var ErrorHandler = func(w http.ResponseWriter, r *http.Request) {
	if e := recover(); e != nil {
		if err, ok := e.(error); ok {
			log.Println("ERROR: ", r.RequestURI, err)
			debug.PrintStack()
			writeJson(w, jsondata.Map{"error": true, "message": err.Error()})
		} else {
			log.Println("RESPONSE ERROR: ", r.RequestURI, e)
			writeJson(w, jsondata.Map{"error": true, "message": e})
		}
	}
}

func New(handler func(http.ResponseWriter, *http.Request) interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer ErrorHandler(w, r)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var resp interface{} = handler(w, r)
		if resp != nil {
			if jsonmap, ok := resp.(*jsondata.Map); ok {
				fmt.Fprint(w, jsonmap)
			} else if err, ok := resp.(error); ok {
				writeJson(w, jsondata.Map{"error": true, "message": err.Error()})
			} else {
				// for arbitrary data
				writeJson(w, resp)
			}
		}
	}
}
