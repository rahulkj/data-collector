package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"code.google.com/p/log4go"
	"net/http"
	"io/ioutil"
	"os"
	"github.com/gorilla/mux"
	"github.com/rahulkj/domain"
)

const (
	HostVar = "VCAP_APP_HOST"
	PortVar = "VCAP_APP_PORT"
)

// error response contains everything we need to use http.Error
type handlerError struct {
	Error   error
	Message string
	Code    int
}

// a custom type that we can use for handling errors and formatting responses
type handler func(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError)

// attach the standard ServeHTTP method to our handler so the http library can call it
func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := make(log4go.Logger)
	log.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())
	// here we could do some prep work before calling the handler if we wanted to

	// call the actual handler
	response, err := fn(w, r)

	// check for errors
	if err != nil {
		log.Error("ERROR: %v\n", err.Error)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Message), err.Code)
		return
	}
	if response == nil {
		log.Debug("ERROR: response from method is nil\n")
		http.Error(w, "Internal server error. Check the logs.", http.StatusInternalServerError)
		return
	}

	// turn the response into JSON
	bytes, e := json.Marshal(response)
	if e != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	// send the response and log
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
	log.Debug("%s %s %s %d", r.RemoteAddr, r.Method, r.URL, 200)
}

func parseDataRequest(r *http.Request) (domain.Data, *handlerError) {
	log := make(log4go.Logger)
	log.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())
	// the book payload is in the request body
	requestData, e := ioutil.ReadAll(r.Body)

	if e != nil {
		return domain.Data{}, &handlerError{e, "Could not read request", http.StatusBadRequest}
	}

	// turn the request body (JSON) into a book object
	var payload domain.Data
	e = json.Unmarshal(requestData, &payload)

	if e != nil {
		return domain.Data{}, &handlerError{e, "Could not parse JSON", http.StatusBadRequest}
	}

	return payload, nil
}

func saveData(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	log := make(log4go.Logger)
	log.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())

	payload, e := parseDataRequest(r)
	if e != nil {
		return nil, e
	}

	payload = domain.SaveData(payload)

	// it's our job to assign IDs, ignore what (if anything) the client sent
	log.Debug("the paylod id is %v\n", payload.Id)
	log.Debug("the environment is %v\n", payload.Environment)

	// we return the book we just made so the client can see the ID if they want
	return payload, nil
}

func main() {
	log := make(log4go.Logger)
	log.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())

	dir := flag.String("directory", "web/", "directory of web files")
	flag.Parse()

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)

	// setup routes
	router := mux.NewRouter()
	router.Handle("/", http.RedirectHandler("/static/", 302))
	router.Handle("/data", handler(saveData)).Methods("POST")
	router.Handle("/data/{id}", handler(saveData)).Methods("POST")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileHandler))
	http.Handle("/", router)

	var port string
	if port = os.Getenv(PortVar); port == "" {
		port = "8080"
	}
	log.Debug("Listening at port %v\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
