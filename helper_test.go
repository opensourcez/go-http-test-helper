package testhelper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Jeffail/gabs"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Test struct {
	Hello string `json:"hello"`
}

type Header struct {
	ContentType string `json:"content_type"`
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

		testInstance := TestHelper{
			TestResultMap: make(map[string]map[string]*gabs.Container),
			ShouldLog:     true,
			Cookies:       make(map[string]*http.Cookie),
		}

		testInstance.TestThis(NewHTTPTest(
			HTTPTestIn{
				Label: "TestHelper", TestCode: "HELPER-001",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/test",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
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

		testInstance := TestHelper{
			TestResultMap: make(map[string]map[string]*gabs.Container),
			ShouldLog:     true,
			Cookies:       make(map[string]*http.Cookie),
		}

		testInstance.TestThis(NewHTTPTest(
			HTTPTestIn{
				Label: "TestHelper", TestCode: "HELPER-002",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/test",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"hello": "hello back at you !",
				},
				KeyPresent: []string{"hello"},
			}), t)

	})

	// Could not find key ( goodbye ) in response body
	// Key (hello) is not suppose to be in response
	t.Run("HelloTestThatShouldFail", func(t *testing.T) {

		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := TestHelper{
			TestResultMap: make(map[string]map[string]*gabs.Container),
			ShouldLog:     true,
			Cookies:       make(map[string]*http.Cookie),
		}

		testInstance.TestThis(NewHTTPTest(
			HTTPTestIn{
				Label: "TestHelper", TestCode: "HELPER-003",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/test",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"hello": "hello back at you !",
				},
				KeyPresent: []string{"goodbye"},
			}), t)

	})

	t.Run("HeaderTest", func(t *testing.T) {

		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := TestHelper{
			TestResultMap: make(map[string]map[string]*gabs.Container),
			ShouldLog:     true,
			Cookies:       make(map[string]*http.Cookie),
		}

		testInstance.TestThis(NewHTTPTest(
			HTTPTestIn{
				Label: "TestHelper", TestCode: "HELPER-003",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/reflect-header",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"content_type": "application/json",
				},
				KeyPresent: []string{"content_type"},
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

	return r
}
