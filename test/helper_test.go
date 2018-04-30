package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Jeffail/gabs"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	th "github.com/zkynet/testhelper"
)

type Test struct {
	Hello string `json:"hello"`
}

type Header struct {
	ContentType string `json:"content_type"`
}
type Cookie struct {
	Value string `json:"Value"`
}

func TestMain(t *testing.T) {

	fmt.Println("Server started on port: " + os.Getenv("Port"))

	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	go http.ListenAndServe(":3333", handlers.CORS(originsOk, methodsOk)(setUpRoutes()))
	runTests(t)
}

func runTests(t *testing.T) {

	t.Run("BaseHelloTest1-parallel", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := th.TestHelper{
			TestResultMap: make(map[string]map[string]*gabs.Container),
			ShouldLog:     true,
			Cookies:       make(map[string]*http.Cookie),
		}

		testInstance.TestThis(th.NewHTTPTest(
			th.HTTPTestIn{
				Label: "TestHelper", TestCode: "HELPER-001",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/test",
				Method:  "POST",
				Headers: headers,
			},
			th.HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"hello": "hello back at you !",
				},
				KeyPresent: []string{"hello"},
			}), t)

	})

	t.Run("BaseHelloTest2-parallel", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := th.TestHelper{
			TestResultMap: make(map[string]map[string]*gabs.Container),
			ShouldLog:     true,
			Cookies:       make(map[string]*http.Cookie),
		}

		testInstance.TestThis(th.NewHTTPTest(
			th.HTTPTestIn{
				Label: "TestHelper", TestCode: "HELPER-002",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/test",
				Method:  "POST",
				Headers: headers,
			},
			th.HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"hello": "hello back at you !",
				},
				KeyPresent: []string{"hello"},
			}), t)

	})

	// Could not find key ( goodbye ) in response body
	// Key (hello) is not suppose to be in response
	t.Run("HelloTestThatShouldFail", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := th.TestHelper{
			TestResultMap: make(map[string]map[string]*gabs.Container),
			ShouldLog:     true,
			Cookies:       make(map[string]*http.Cookie),
		}

		testInstance.TestThis(th.NewHTTPTest(
			th.HTTPTestIn{
				Label: "TestHelper", TestCode: "HELPER-003",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/test",
				Method:  "POST",
				Headers: headers,
			},
			th.HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"hello": "hello back at you !",
				},
				KeyPresent: []string{"goodbye"},
			}), t)

	})

	t.Run("HeaderTest", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := th.TestHelper{
			TestResultMap: make(map[string]map[string]*gabs.Container),
			ShouldLog:     true,
			Cookies:       make(map[string]*http.Cookie),
		}

		testInstance.TestThis(th.NewHTTPTest(
			th.HTTPTestIn{
				Label: "HeaderTest", TestCode: "HELPER-004",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/reflect-header",
				Method:  "POST",
				Headers: headers,
			},
			th.HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"content_type": "application/json",
				},
				KeyPresent: []string{"content_type"},
			}), t)

	})

	t.Run("CookieTest", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := th.TestHelper{
			TestResultMap: make(map[string]map[string]*gabs.Container),
			ShouldLog:     true,
			Cookies:       make(map[string]*http.Cookie),
		}

		testInstance.TestThis(th.NewHTTPTest(
			th.HTTPTestIn{
				Label: "TestGetCookie", TestCode: "HELPER-005",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/get-cookie",
				Method:  "POST",
				Headers: headers,
			},
			th.HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"Name":  "cookiemonster",
					"Value": "cookiemonster",
				},
				KeyPresent: []string{"Name", "Value", "Path", "MaxAge", "HttpOnly", "Domain", "Expires", "RawExpires", "Secure", "Raw", "Unparsed"},
			}), t)
		fmt.Println(testInstance.TestResultMap)
		testInstance.TestThis(th.NewHTTPTest(
			th.HTTPTestIn{
				Label: "TestSendCookie", TestCode: "HELPER-006",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/send-cookie",
				Method:  "POST",
				Headers: headers,
			},
			th.HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"Name":  "cookiemonster",
					"Value": "cookiemonster",
				},
				KeyPresent: []string{"Name", "Value", "Path", "MaxAge", "HttpOnly", "Domain", "Expires", "RawExpires", "Secure", "Raw", "Unparsed"},
			}), t)

	})
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

	return r
}
