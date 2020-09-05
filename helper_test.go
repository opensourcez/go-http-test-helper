package testhelper

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gorilla/handlers"
)

var logging = true
var port = "31337"
var baseURL = "http://localhost:" + port

func TestMain(t *testing.T) {

	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	go func() {
		err := http.ListenAndServe(":"+port, handlers.CORS(originsOk, methodsOk)(setUpRoutes()))
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		fmt.Println("Test server started on port: " + port)
	}()
	// sleeping for 2 seconds to make sure webserver is launched
	time.Sleep(2 * time.Second)
	runTests(t)
}

func runTests(t *testing.T) {

	t.Run("BaseTestParallel1", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := NewHTTPTestHelper(true, "", "", "")

		testInstance.TestThis(&HTTPTest{
			HTTPTestIn{
				Note:  "This is a test!",
				Label: "TestHelper", TestCode: "HELPER-001",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     baseURL + "/test",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"hello": "hello back at you !",
				},
				Keys: []string{"hello"},
				Headers: map[string]string{
					"Content-Type": "text/plain; charset=utf-8",
				},
				IgnoredHeaders: []string{"Host"},
			}}, t)

	})

	t.Run("BaseTestParallel2", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := NewHTTPTestHelper(true, "", "", "")

		testInstance.TestThis(&HTTPTest{
			HTTPTestIn{
				Label: "TestHelper", TestCode: "HELPER-002",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     baseURL + "/test",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"hello": "hello back at you !",
				},
				Keys: []string{"hello"},
			}}, t)

	})

	t.Run("HeaderTest", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := NewHTTPTestHelper(true, "", "", "")

		testInstance.TestThis(&HTTPTest{
			HTTPTestIn{
				Label: "HeaderTest", TestCode: "HELPER-003",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     baseURL + "/reflect-header",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"content_type": "application/json",
				},
				Keys: []string{"content_type"},
			}}, t)

	})

	t.Run("RawBodyTest", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}
		testInstance := NewHTTPTestHelper(true, "", "", "")

		testInstance.TestThis(&HTTPTest{
			HTTPTestIn{
				Label: "TestGetCookie", TestCode: "HELPER-004",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     baseURL + "/raw-body",
				Method:  "GET",
				Headers: headers,
			},
			HTTPTestOut{RawBody: []byte(`{"raw-body":"test value"}`), Code: 200, Status: "200 OK"}}, t)

	})

	t.Run("CookieTest", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}
		testInstance := NewHTTPTestHelper(true, "", "", "")

		testInstance.TestThis(&HTTPTest{
			HTTPTestIn{
				Label: "TestGetCookie", TestCode: "HELPER-004",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     baseURL + "/get-cookie",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"Name":  "cookiemonster",
					"Value": "cookiemonster",
				},
				Keys: []string{"Name", "Value", "Path", "MaxAge", "HttpOnly", "Domain", "Expires", "RawExpires", "Secure", "Raw", "Unparsed"},
			}}, t)

		testInstance.TestThis(&HTTPTest{
			HTTPTestIn{
				Label: "TestSendCookie", TestCode: "HELPER-005",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     baseURL + "/send-cookie",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"Name":  "cookiemonster",
					"Value": "cookiemonster",
				},
				Keys: []string{"Name", "Value", "Path", "MaxAge", "HttpOnly", "Domain", "Expires", "RawExpires", "Secure", "Raw", "Unparsed"},
			}}, t)

	})

	t.Run("FailCases", func(t *testing.T) {
		t.Parallel()
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		testInstance := NewHTTPTestHelper(true, "", "", "")

		testInstance.TestThis(&HTTPTest{
			HTTPTestIn{
				Note:  "This test should fail because of two missing keys in body: name and value",
				Label: "TestEmptyBody", TestCode: "HELPER-006",
				Body:    nil,
				URL:     baseURL + "/empty-body",
				Method:  "GET",
				Headers: headers,
			},
			HTTPTestOut{Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"Name":  "cookiemonster",
					"Value": "cookiemonster",
				},
				Keys: []string{"Name", "Value"},
			}}, t)

		testInstance.TestThis(&HTTPTest{
			HTTPTestIn{
				Note:  "This test should fail with the message: Key ( goodbye ) not found in response",
				Label: "TestHelper", TestCode: "HELPER-007",
				Body:    []byte(`{"hello":"hello back at you !"}`),
				URL:     baseURL + "/test",
				Method:  "POST",
				Headers: headers,
			},
			HTTPTestOut{Code: 200, Status: "200 OK",
				KeyValues: map[string]string{
					"hello": "hello back at you !",
				},
				Keys: []string{"goodbye"},
			}}, t)

	})
}
