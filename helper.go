package testhelper

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/Jeffail/gabs"
)

var endColor = "\033[0m"

type TestHelper struct {
	ShouldLog      bool
	CookieBucket   map[string]*http.Cookie
	ResponseBucket map[string]map[string]*gabs.Container
	ErrorColor     string
	SuccessColor   string
	InfoColor      string
}

func NewHTTPTestHelper(
	logging bool,
	errorColor string,
	successColor string,
	infoColor string,
) *TestHelper {
	if errorColor == "" {
		errorColor = "\033[31m"
	}
	if errorColor == "none" {
		errorColor = ""
	}
	if successColor == "" {
		successColor = "\033[32m"
	}
	if successColor == "none" {
		successColor = ""
	}
	if infoColor == "" {
		infoColor = "\033[0m"
	}
	if infoColor == "none" {
		infoColor = ""
	}
	return &TestHelper{
		ErrorColor:     errorColor,
		SuccessColor:   successColor,
		InfoColor:      infoColor,
		ShouldLog:      logging,
		ResponseBucket: make(map[string]map[string]*gabs.Container),
		CookieBucket:   make(map[string]*http.Cookie),
	}

}

type HTTPTestIn struct {
	Note     string
	Label    string
	TestCode string
	Body     []byte
	URL      string
	Method   string
	Headers  map[string]string
}

type HTTPTestOut struct {
	RawBody        []byte
	KeyValues      map[string]string
	Keys           []string
	Status         string
	Code           int
	Headers        map[string]string
	IgnoredHeaders []string
}

type HTTPTest struct {
	HTTPTestIn
	HTTPTestOut
}

func (th *TestHelper) sendRequest(HTTPTest *HTTPTest, t *testing.T) (*http.Response, []byte) {

	if th.ShouldLog {
		t.Log(th.InfoColor, "===============================================================", endColor)
		t.Log(th.InfoColor, HTTPTest.HTTPTestIn.Method, "(", HTTPTest.HTTPTestIn.URL, ")", endColor)
	}
	req, err := http.NewRequest(HTTPTest.HTTPTestIn.Method, HTTPTest.HTTPTestIn.URL, bytes.NewBuffer(HTTPTest.HTTPTestIn.Body))
	if err != nil {
		t.Error(th.ErrorColor, "Could not make request:", endColor, err)
		t.Skip()
	}
	for _, v := range th.CookieBucket {
		req.AddCookie(v)
	}

	for i, v := range HTTPTest.HTTPTestIn.Headers {
		req.Header.Set(i, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(th.ErrorColor, "Could not send request:", endColor, err)
		t.Skip()
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if th.ShouldLog {
		t.Log(th.InfoColor, "CODE (", strconv.Itoa(resp.StatusCode), ") STATUS (", resp.Status, ")", endColor)
		if len(body) > 0 {
			t.Log(th.InfoColor, strings.TrimSuffix(string(body), "\n"), endColor)
		} else {
			t.Log(th.InfoColor, "NO RESPONSE BODY", endColor)
		}

		if len(resp.Cookies()) < 1 {
			t.Log(th.InfoColor, "===============================================================", endColor)
		}
	}

	if len(resp.Cookies()) > 0 {
		if th.ShouldLog {
			t.Log(th.InfoColor, "==================== RECEIVED NEW CookieBucket ===============", endColor)
		}

		for _, v := range resp.Cookies() {
			if th.ShouldLog {
				t.Log(v)
			}
			th.CookieBucket[v.Name] = v
		}
		if th.ShouldLog {
			t.Log(th.InfoColor, "==============================================================", endColor)
		}

	}

	return resp, body
}

func (th *TestHelper) checkHTTPStatus(response *http.Response, expectedStatus string, t *testing.T) {
	if response.Status != expectedStatus {
		t.Error(th.ErrorColor, "Expected Status (", expectedStatus, ") but got (", response.Status, ")", endColor)
	}
	if th.ShouldLog {
		t.Log(th.SuccessColor, "Wanted status (", expectedStatus, ") and got (", response.Status, ")", endColor)
	}

}

func (th *TestHelper) checkHTTPCode(response *http.Response, expectedCode int, t *testing.T) {
	if response.StatusCode != expectedCode {
		t.Error(th.ErrorColor, "Expected Code (", strconv.Itoa(expectedCode), ") but got (", strconv.Itoa(response.StatusCode), ")", endColor)
	}
	if th.ShouldLog {
		t.Log(th.SuccessColor, "Wanted code (", strconv.Itoa(expectedCode), ") and got( ", strconv.Itoa(response.StatusCode), ")", endColor)
	}
}

func (th *TestHelper) checkFields(decodedBody map[string]*gabs.Container, Fields []string, t *testing.T) {
	if len(decodedBody) < 1 && len(Fields) < 1 {
		return
	} else if len(decodedBody) < 1 && len(Fields) > 0 {
		if th.ShouldLog {
			t.Error(th.ErrorColor, "Expecting (", strconv.Itoa(len(Fields)), ") fields in response but got (", strconv.Itoa(len(decodedBody)), ")", endColor)
		}
		return
	}

	for _, key := range Fields {
		continueInOuterLoop := false
		for decodedBodyKey := range decodedBody {
			if decodedBodyKey == key {
				if th.ShouldLog {
					t.Log(th.SuccessColor, "Key (", key, ") found in response", endColor)
				}
				continueInOuterLoop = true
				continue

			}
		}
		if continueInOuterLoop {
			continue
		}
		t.Error(th.ErrorColor, "Key (", key, ") not found in response", endColor)
	}

}
func (th *TestHelper) checkKeyValues(decodedBody map[string]*gabs.Container, KeyValues map[string]string, t *testing.T) {
	for key, value := range KeyValues {

		var valueToCheck string
		decodedBodyValue := decodedBody[key].Data()
		if decodedBodyValue == nil {
			t.Error(th.ErrorColor, "Key ( "+key+" ) with value ("+value+") not found in request", endColor)
			continue
		}

		switch reflect.TypeOf(decodedBodyValue).Kind() {
		case reflect.Bool:
			if th.ShouldLog {
				t.Log(th.InfoColor, "Key ", key, "is of type ( bool )", endColor)
			}
			valueToCheck = strconv.FormatBool(decodedBodyValue.(bool))
		case reflect.Int:
			if th.ShouldLog {
				t.Log(th.InfoColor, "Key ", key, "is of type ( int )", endColor)
			}
			valueToCheck = strconv.Itoa(decodedBodyValue.(int))
		case reflect.Float64:
			if th.ShouldLog {
				t.Log(th.InfoColor, "Key ", key, "is of type ( float64 )", endColor)
			}
			valueToCheck = strconv.FormatFloat(decodedBodyValue.(float64), 'f', -1, 64)
		default:
			if th.ShouldLog {
				t.Log(th.InfoColor, "Key ", key, "is of type ( string )", endColor)
			}
			valueToCheck = decodedBodyValue.(string)
		}

		if valueToCheck != value {
			t.Error(th.ErrorColor, "Expected ( "+value+" ) in key ( "+key+" ) but got ( "+valueToCheck+" )", endColor)
			return
		}

		if th.ShouldLog {
			t.Log(th.SuccessColor, "Wanted value (", value, ") in key (", key, ") and got (", valueToCheck, ")", endColor)
		}
	}
}

func (th *TestHelper) decodeBody(body []byte, t *testing.T) map[string]*gabs.Container {

	if len(body) < 1 {
		return nil
	}

	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		t.Error(th.ErrorColor, "Request body could not be converted to JSON or XML:\033[30m", endColor, err)
		return nil
	}

	children, err := jsonParsed.S().ChildrenMap()
	if err != nil {
		t.Error(th.ErrorColor, "JSON coult not be converted to GABS container:", endColor, err)
		return nil
	}
	return children
}

func (th *TestHelper) checkHeaders(response *http.Response, out *HTTPTestOut, t *testing.T) {

	for header, expectedHeaderValue := range out.Headers {
		for _, headerToBeIgnored := range out.IgnoredHeaders {
			if header == headerToBeIgnored {
				continue
			}
		}

		actualHeaderValue := response.Header.Get(header)
		if actualHeaderValue != expectedHeaderValue {
			t.Error(th.ErrorColor, "Excpected header (", header, ") with value (", expectedHeaderValue, ") but got (", actualHeaderValue, ")", endColor)
		}

		if th.ShouldLog {
			t.Log(th.SuccessColor, "Found header (", header, ") with value (", actualHeaderValue, ") in response", endColor)
		}
	}
}

func (th *TestHelper) checkRawBody(responseBody string, expectedBody string, t *testing.T) {
	if strings.TrimRight(responseBody, "\n") != strings.TrimRight(expectedBody, "\n") {
		t.Error(th.ErrorColor, "Excpected body: ", strings.TrimRight(responseBody, "\n"), "\n but got: ", strings.TrimRight(expectedBody, "\n"), endColor)
	}
}

func (th *TestHelper) TestThis(
	HTTPTest *HTTPTest,
	t *testing.T) {
	t.Run(HTTPTest.HTTPTestIn.TestCode+":"+HTTPTest.HTTPTestIn.Label, func(t *testing.T) {

		if HTTPTest.HTTPTestIn.Note != "" {
			t.Log(th.InfoColor, "==================== NOTE =====================================", endColor)
			t.Log(th.InfoColor, HTTPTest.HTTPTestIn.Note, endColor)
		}

		response, body := th.sendRequest(HTTPTest, t)

		th.checkHTTPStatus(response, HTTPTest.HTTPTestOut.Status, t)
		th.checkHTTPCode(response, HTTPTest.HTTPTestOut.Code, t)
		th.checkHeaders(response, &HTTPTest.HTTPTestOut, t)

		if HTTPTest.HTTPTestOut.RawBody != nil {
			th.checkRawBody(string(body), string(HTTPTest.HTTPTestOut.RawBody), t)
		} else {
			th.ResponseBucket[HTTPTest.HTTPTestIn.TestCode] = th.decodeBody(body, t)
			th.checkKeyValues(th.ResponseBucket[HTTPTest.HTTPTestIn.TestCode], HTTPTest.HTTPTestOut.KeyValues, t)
			th.checkFields(th.ResponseBucket[HTTPTest.HTTPTestIn.TestCode], HTTPTest.HTTPTestOut.Keys, t)
		}

	})
}
