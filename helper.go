package testhelper

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/Jeffail/gabs"
)

var ShouldLog bool

var cookies = map[string]*http.Cookie{}

type HTTPTestIn struct {
	Label    string
	TestCode string
	Body     []byte
	URL      string
	Method   string
	Headers  map[string]string
}

type KeyValueInBody struct {
	Key   string
	Value string
}

type HTTPTestOut struct {
	Body       string
	KeyValues  []*KeyValueInBody
	KeyPresent []string
	Status     string
	Code       int
}

type HTTPTest struct {
	HTTPTestIn
	HTTPTestOut
}

func NewHTTPTest(
	HTTPTestIn HTTPTestIn,
	HTTPTestOut HTTPTestOut,
) *HTTPTest {
	return &HTTPTest{
		HTTPTestIn:  HTTPTestIn,
		HTTPTestOut: HTTPTestOut,
	}
}

var TestResultMap = make(map[string]map[string]*gabs.Container)

func sendRequest(HTTPTest *HTTPTest) (*http.Response, []byte) {
	if shouldLog {
		log.Println("====== Sending method ( " + HTTPTest.HTTPTestIn.Method + " ) to =======")
		log.Println(HTTPTest.HTTPTestIn.URL)
	}

	req, err := http.NewRequest(HTTPTest.HTTPTestIn.Method, HTTPTest.HTTPTestIn.URL, bytes.NewBuffer(HTTPTest.HTTPTestIn.Body))
	for _, v := range cookies {
		req.AddCookie(v)
	}

	for i, v := range HTTPTest.HTTPTestIn.Headers {
		req.Header.Set(i, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if shouldLog {
		log.Println("====== Received ( " + resp.Status + " ) =======")
		if len(body) > 0 {
			log.Println(strings.TrimSuffix(string(body), "\n"))
		}
		log.Println("===============================================")
	}

	if len(resp.Cookies()) > 0 {
		if shouldLog {
			log.Println("====== Received new cookies =======")
		}

		for _, v := range resp.Cookies() {
			if shouldLog {
				log.Println(v)
			}
			cookies[v.Name] = v
		}
		if shouldLog {
			log.Println("===============================================")
		}

	}

	return resp, body
}

func checkForStatusAndCode(response *http.Response, expectedCode int, expectedStatus string, t *testing.T) bool {
	isOk := true
	if response.Status != expectedStatus {
		isOk = false
		t.Error("Expected Status ( " + expectedStatus + " ) but got: " + response.Status)
	}
	if response.StatusCode != expectedCode {
		isOk = false
		t.Error("Expected Code ( " + strconv.Itoa(expectedCode) + " ) but got: " + strconv.Itoa(response.StatusCode))
	}
	return isOk
}

func checkFields(decodedBody map[string]*gabs.Container, Fields []string, t *testing.T) {
	if len(decodedBody) < 1 && len(Fields) < 1 {
		return
	} else if len(decodedBody) < 1 && len(Fields) > 0 {
		t.Error("No fields in response body but should have (" + strconv.Itoa(len(Fields)) + " )")
	}

	for _, key := range Fields {
		if decodedBody[key].Data() == nil {
			t.Error("Could not find key ( " + key + " ) in response body")
		}
	}
	for decodedBodyKey := range decodedBody {
		shouldContinue := false
		for _, key := range Fields {
			if decodedBodyKey == key {
				shouldContinue = true
			}
		}
		if shouldContinue {
			continue
		}
		t.Error("Key (" + decodedBodyKey + ") is not suppose to be in response")
	}

}
func checkKeyValues(decodedBody map[string]*gabs.Container, KeyValues []*KeyValueInBody, t *testing.T) {
	if len(decodedBody) < 1 {
		return
	}

	for _, KeyValue := range KeyValues {
		var valueToCheck string

		if decodedBody[KeyValue.Key].Data() == nil {
			t.Error("Key ( " + KeyValue.Key + " ) with value (" + KeyValue.Value + ") not found in request")
			continue
		}

		if reflect.TypeOf(decodedBody[KeyValue.Key].Data()).Kind() == reflect.Bool {
			valueToCheck = strconv.FormatBool(decodedBody[KeyValue.Key].Data().(bool))
		} else if reflect.TypeOf(decodedBody[KeyValue.Key].Data()).Kind() == reflect.Float64 {
			valueToCheck = strconv.FormatFloat(decodedBody[KeyValue.Key].Data().(float64), 'f', -1, 64)
		} else {
			valueToCheck = decodedBody[KeyValue.Key].Data().(string)
		}
		if valueToCheck != KeyValue.Value {
			t.Error("Expected ( " + KeyValue.Value + " ) in key ( " + KeyValue.Key + " ) but got ( " + valueToCheck + " )")
		}
	}
}

func decodeBody(body []byte) map[string]*gabs.Container {

	if len(body) < 1 {
		return nil
	}

	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		log.Println(err)
	}

	children, err := jsonParsed.S().ChildrenMap()
	if err != nil {
		log.Println(err)
	}

	return children
}

func TestThis(HTTPTest *HTTPTest, t *testing.T) {
	t.Run(HTTPTest.HTTPTestIn.TestCode+":"+HTTPTest.HTTPTestIn.Label, func(t *testing.T) {

		response, body := sendRequest(HTTPTest)

		decodedBody := decodeBody(body)
		TestResultMap[HTTPTest.HTTPTestIn.TestCode] = decodedBody

		if checkForStatusAndCode(response, HTTPTest.HTTPTestOut.Code, HTTPTest.HTTPTestOut.Status, t) {
			checkKeyValues(decodedBody, HTTPTest.HTTPTestOut.KeyValues, t)
			checkFields(decodedBody, HTTPTest.HTTPTestOut.KeyPresent, t)
		}
	})
}
