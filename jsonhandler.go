package jsonhandler

import (
	"encoding/json"
	"fmt"
	"github.com/c9s/jsondata"
	"log"
	"net/http"
	"runtime/debug"
)

const Padding = "  "

func WriteHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func DecodeBody(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(val)
}

// General Function to write json response.
func WriteJson(w http.ResponseWriter, val interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	b, err := json.Marshal(val)
	if err != nil {
		WriteErrorJson(w, err)
	}
	fmt.Fprint(w, string(b))
}

func WriteError(w http.ResponseWriter, e interface{}) {
	WriteErrorJson(w, e)
}

func WriteErrorJson(w http.ResponseWriter, e interface{}) {
	if err, ok := e.(error); ok {
		log.Println("ERROR: ", err)
		debug.PrintStack()
		WriteJson(w, jsondata.Map{"error": true, "message": err.Error()})
	} else {
		log.Println("RESPONSE ERROR: ", e)
		WriteJson(w, jsondata.Map{"error": true, "message": e})
	}
	WriteJson(w, e)
}

var ErrorHandler = func(w http.ResponseWriter, r *http.Request) {
	if e := recover(); e != nil {
		if err, ok := e.(error); ok {
			log.Println("ERROR: ", r.RequestURI, err)
			debug.PrintStack()
			WriteJson(w, jsondata.Map{"error": true, "message": err.Error()})
		} else {
			log.Println("RESPONSE ERROR: ", r.RequestURI, e)
			WriteJson(w, jsondata.Map{"error": true, "message": e})
		}
	}
}

func New(handler func(http.ResponseWriter, *http.Request) interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer ErrorHandler(w, r)
		if resp := handler(w, r); resp != nil {
			// if resp != nil {
			if jsonmap, ok := resp.(*jsondata.Map); ok {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				fmt.Fprint(w, jsonmap)
				resp = nil
			} else if err, ok := resp.(error); ok {
				WriteJson(w, jsondata.Map{"error": true, "message": err.Error()})
			} else {
				// for arbitrary data
				WriteJson(w, resp)
				resp = nil // free
			}
		}
	}
}
