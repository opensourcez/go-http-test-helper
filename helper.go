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
var Cookies map[string]*http.Cookie
var TestResultMap map[string]map[string]*gabs.Container

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

func sendRequest(HTTPTest *HTTPTest) (*http.Response, []byte) {
	if ShouldLog {
		log.Println("====== Sending method ( " + HTTPTest.HTTPTestIn.Method + " ) to =======")
		log.Println(HTTPTest.HTTPTestIn.URL)
	}

	req, err := http.NewRequest(HTTPTest.HTTPTestIn.Method, HTTPTest.HTTPTestIn.URL, bytes.NewBuffer(HTTPTest.HTTPTestIn.Body))
	for _, v := range Cookies {
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

	if ShouldLog {
		log.Println("====== Received ( " + resp.Status + " ) =======")
		if len(body) > 0 {
			log.Println(strings.TrimSuffix(string(body), "\n"))
		}
		log.Println("===============================================")
	}

	if len(resp.Cookies()) > 0 {
		if ShouldLog {
			log.Println("====== Received new cookies =======")
		}

		for _, v := range resp.Cookies() {
			if ShouldLog {
				log.Println(v)
			}
			Cookies[v.Name] = v
		}
		if ShouldLog {
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
		decodedBodyValue := decodedBody[KeyValue.Key].Data()

		if decodedBodyValue == nil {
			t.Error("Key ( " + KeyValue.Key + " ) with value (" + KeyValue.Value + ") not found in request")
			continue
		}

		switch reflect.TypeOf(decodedBodyValue).Kind() {
		case reflect.Bool:
			valueToCheck = strconv.FormatBool(decodedBodyValue.(bool))
		case reflect.Float64:
			valueToCheck = strconv.FormatFloat(decodedBodyValue.(float64), 'f', -1, 64)
		default:
			valueToCheck = decodedBodyValue.(string)
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
