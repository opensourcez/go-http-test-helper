package testhelper

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Jeffail/gabs"
	"github.com/gorilla/handlers"
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

var logging = false

func TestMain(t *testing.T) {

	fmt.Println("Server started on port: " + os.Getenv("Port"))

	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	go http.ListenAndServe(":3333", handlers.CORS(originsOk, methodsOk)(setUpRoutes()))
	runTest(t)
}

func runTest(t *testing.T) {

	t.Run("BaseHelloTest1-parallel", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := TestHelper{
			ResponseBucket: make(map[string]map[string]*gabs.Container),
			ShouldLog:      logging,
			Cookies:        make(map[string]*http.Cookie),
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
				KeyValuesInBody: map[string]string{
					"hello": "hello back at you !",
				},
				KeysPresentInBody: []string{"hello"},
			}), t)

	})

	t.Run("BaseHelloTest2-parallel", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := TestHelper{
			ResponseBucket: make(map[string]map[string]*gabs.Container),
			ShouldLog:      logging,
			Cookies:        make(map[string]*http.Cookie),
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
				KeyValuesInBody: map[string]string{
					"hello": "hello back at you !",
				},
				KeysPresentInBody: []string{"hello"},
			}), t)

	})

	t.Run("HeaderTest", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := TestHelper{
			ResponseBucket: make(map[string]map[string]*gabs.Container),
			ShouldLog:      logging,
			Cookies:        make(map[string]*http.Cookie),
		}

		testInstance.TestThis(NewHTTPTest(
			HTTPTestIn{
				Label: "HeaderTest", TestCode: "HELPER-003",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/reflect-header",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValuesInBody: map[string]string{
					"content_type": "application/json",
				},
				KeysPresentInBody: []string{"content_type"},
			}), t)

	})

	t.Run("CookieTest", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := TestHelper{
			ResponseBucket: make(map[string]map[string]*gabs.Container),
			ShouldLog:      logging,
			Cookies:        make(map[string]*http.Cookie),
		}

		testInstance.TestThis(NewHTTPTest(
			HTTPTestIn{
				Label: "TestGetCookie", TestCode: "HELPER-004",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/get-cookie",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValuesInBody: map[string]string{
					"Name":  "cookiemonster",
					"Value": "cookiemonster",
				},
				KeysPresentInBody: []string{"Name", "Value", "Path", "MaxAge", "HttpOnly", "Domain", "Expires", "RawExpires", "Secure", "Raw", "Unparsed"},
			}), t)

		testInstance.TestThis(NewHTTPTest(
			HTTPTestIn{
				Label: "TestSendCookie", TestCode: "HELPER-005",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/send-cookie",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValuesInBody: map[string]string{
					"Name":  "cookiemonster",
					"Value": "cookiemonster",
				},
				KeysPresentInBody: []string{"Name", "Value", "Path", "MaxAge", "HttpOnly", "Domain", "Expires", "RawExpires", "Secure", "Raw", "Unparsed"},
			}), t)

	})

	t.Run("FailCases", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := TestHelper{
			ResponseBucket: make(map[string]map[string]*gabs.Container),
			ShouldLog:      logging,
			Cookies:        make(map[string]*http.Cookie),
		}

		testInstance.TestThis(NewHTTPTest(
			HTTPTestIn{
				Label: "TestEmptyBody", TestCode: "HELPER-006",
				Body:    nil,
				URL:     "http://localhost:3333/empty-body",
				Method:  "GET",
				Headers: headers,
			},
			HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValuesInBody: map[string]string{
					"Name":  "cookiemonster",
					"Value": "cookiemonster",
				},
				KeysPresentInBody: []string{"Name", "Value"},
			}), t)

		testInstance.TestThis(NewHTTPTest(
			HTTPTestIn{
				Label: "TestHelper", TestCode: "HELPER-007",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     "http://localhost:3333/test",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
				KeyValuesInBody: map[string]string{
					"hello": "hello back at you !",
				},
				KeysPresentInBody: []string{"goodbye"},
			}), t)

	})
}
