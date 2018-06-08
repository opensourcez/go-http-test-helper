package testhelper

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Test struct {
	Hello string `json:"hello"`
}

type RawBody struct {
	RawBody string `json:"raw-body"`
}

type Header struct {
	ContentType string `json:"content_type"`
}

func setUpRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {

		var test Test
		if err := json.NewDecoder(r.Body).Decode(&test); err != nil {
			panic(err)
		}

		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(test); err != nil {
			panic(err)
		}

	}).Methods("POST")

	r.HandleFunc("/reflect-header", func(w http.ResponseWriter, r *http.Request) {

		header := Header{
			ContentType: r.Header.Get("Content-Type"),
		}

		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(header); err != nil {
			panic(err)
		}

	}).Methods("POST")

	r.HandleFunc("/send-cookie", func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("cookiemonster")
		if err != nil {
			panic(err)
		}

		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(cookie); err != nil {
			panic(err)
		}

	}).Methods("POST")

	r.HandleFunc("/raw-body", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(&RawBody{RawBody: "test value"}); err != nil {
			panic(err)
		}

	}).Methods("GET")

	r.HandleFunc("/get-cookie", func(w http.ResponseWriter, r *http.Request) {

		cookie := &http.Cookie{
			Name:     "cookiemonster",
			Value:    "cookiemonster",
			Path:     "/",
			Secure:   false,
			HttpOnly: false,
			Domain:   "localhost",
		}
		cookie.Unparsed = nil
		http.SetCookie(w, cookie)

		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(cookie); err != nil {
			panic(err)
		}

	}).Methods("POST")

	r.HandleFunc("/empty-body", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}).Methods("GET")

	return r
}
