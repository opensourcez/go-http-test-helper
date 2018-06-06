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

type TestHelper struct {
	ShouldLog      bool
	Cookies        map[string]*http.Cookie
	ResponseBucket map[string]map[string]*gabs.Container
}
type HTTPTestIn struct {
	Label    string
	TestCode string
	Body     []byte
	URL      string
	Method   string
	Headers  map[string]string
}

type HTTPTestOut struct {
	Body              string
	KeyValuesInBody   map[string]string
	KeysPresentInBody []string
	Status            string
	Code              int
	Headers           map[string]string
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

func (th *TestHelper) sendRequest(HTTPTest *HTTPTest, t *testing.T) (*http.Response, []byte) {

	if th.ShouldLog {
		t.Log("\033[35m==============================================================\033[0m")
		t.Log(HTTPTest.HTTPTestIn.Method, "(", HTTPTest.HTTPTestIn.URL, ")")
	}
	req, err := http.NewRequest(HTTPTest.HTTPTestIn.Method, HTTPTest.HTTPTestIn.URL, bytes.NewBuffer(HTTPTest.HTTPTestIn.Body))
	if err != nil {
		t.Error("\033[31mCould not make request:\033[0m ", err)
		t.Skip()
	}
	for _, v := range th.Cookies {
		req.AddCookie(v)
	}

	for i, v := range HTTPTest.HTTPTestIn.Headers {
		req.Header.Set(i, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error("\033[31mCould not send request:\033[0m ", err)
		t.Skip()
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if th.ShouldLog {
		t.Log("CODE (", strconv.Itoa(resp.StatusCode), ") STATUS (", resp.Status, ")")
		if len(body) > 0 {
			t.Log(strings.TrimSuffix(string(body), "\n"))
		} else {
			t.Log("NO RESPONSE BODY")
		}

		if len(resp.Cookies()) < 1 {
			t.Log("\033[35m==============================================================\033[0m")
		}
	}

	if len(resp.Cookies()) > 0 {
		if th.ShouldLog {
			t.Log("\033[35m==================== RECEIVED NEW COOKIES ====================\033[0m")
		}

		for _, v := range resp.Cookies() {
			if th.ShouldLog {
				t.Log(v)
			}
			th.Cookies[v.Name] = v
		}
		if th.ShouldLog {
			t.Log("\033[35m==============================================================\033[0m")
		}

	}

	return resp, body
}

func (th *TestHelper) checkHTTPStatus(response *http.Response, expectedStatus string, t *testing.T) {
	if response.Status != expectedStatus {
		t.Error("\033[31mExpected Status (", expectedStatus, ") but got (", response.Status, ")\033[0m")
	}
	if th.ShouldLog {
		t.Log("Wanted status (", expectedStatus, ") and got (", response.Status, ")")
	}

}

func (th *TestHelper) checkHTTPCode(response *http.Response, expectedCode int, t *testing.T) {
	if response.StatusCode != expectedCode {
		t.Error("\033[31mExpected Code (", strconv.Itoa(expectedCode), ") but got (", strconv.Itoa(response.StatusCode), ")\033[0m")
	}
	if th.ShouldLog {
		t.Log("Wanted code (", strconv.Itoa(expectedCode), ") and got( ", strconv.Itoa(response.StatusCode), ")")
	}
}

func (th *TestHelper) checkFields(decodedBody map[string]*gabs.Container, Fields []string, t *testing.T) {
	if len(decodedBody) < 1 && len(Fields) < 1 {
		return
	} else if len(decodedBody) < 1 && len(Fields) > 0 {
		if th.ShouldLog {
			t.Error("\033[31mExpecting (", strconv.Itoa(len(Fields)), ") fields in response but got (", strconv.Itoa(len(decodedBody)), ")\033[0m")
		}
		return
	}

	for _, key := range Fields {
		continueInOuterLoop := false
		for decodedBodyKey := range decodedBody {
			if decodedBodyKey == key {
				if th.ShouldLog {
					t.Log("Key (", key, ") found in response")
				}
				continueInOuterLoop = true
				continue

			}
		}
		if continueInOuterLoop {
			continue
		}
		t.Error("\033[31mKey (", key, ") not found in response\033[0m")
	}

}
func (th *TestHelper) checkKeyValues(decodedBody map[string]*gabs.Container, KeyValues map[string]string, t *testing.T) {
	for key, value := range KeyValues {

		var valueToCheck string
		decodedBodyValue := decodedBody[key].Data()
		if decodedBodyValue == nil {
			t.Error("\033[31mKey ( " + key + " ) with value (" + value + ") not found in request\033[0m")
			continue
		}

		switch reflect.TypeOf(decodedBodyValue).Kind() {
		case reflect.Bool:
			if th.ShouldLog {
				t.Log("Key ", key, "is of type ( bool )")
			}
			valueToCheck = strconv.FormatBool(decodedBodyValue.(bool))
		case reflect.Int:
			if th.ShouldLog {
				t.Log("Key ", key, "is of type ( int )")
			}
			valueToCheck = strconv.Itoa(decodedBodyValue.(int))
		case reflect.Float64:
			if th.ShouldLog {
				t.Log("Key ", key, "is of type ( float64 )")
			}
			valueToCheck = strconv.FormatFloat(decodedBodyValue.(float64), 'f', -1, 64)
		default:
			if th.ShouldLog {
				t.Log("Key ", key, "is of type ( string )")
			}
			valueToCheck = decodedBodyValue.(string)
		}

		if valueToCheck != value {
			t.Error("\033[31mExpected ( " + value + " ) in key ( " + key + " ) but got ( " + valueToCheck + " )\033[0m")
			return
		}

		if th.ShouldLog {
			t.Log("Wanted value (", value, ") in key (", key, ") and got (", valueToCheck, ")")
		}
	}
}

func (th *TestHelper) decodeBody(body []byte, t *testing.T) map[string]*gabs.Container {

	if len(body) < 1 {
		return nil
	}

	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		t.Error("\033[31mRequest body could not be converted to JSON or XML:\033[30m", err)
		return nil
	}

	children, err := jsonParsed.S().ChildrenMap()
	if err != nil {
		t.Error("\033[31mJSON coult not be converted to GABS container:\033[30m", err)
		return nil
	}
	return children
}

func (th *TestHelper) TestThis(
	HTTPTest *HTTPTest,
	t *testing.T) {
	t.Run(HTTPTest.HTTPTestIn.TestCode+":"+HTTPTest.HTTPTestIn.Label, func(t *testing.T) {

		response, body := th.sendRequest(HTTPTest, t)

		th.ResponseBucket[HTTPTest.HTTPTestIn.TestCode] = th.decodeBody(body, t)

		th.checkHTTPStatus(response, HTTPTest.HTTPTestOut.Status, t)
		th.checkHTTPCode(response, HTTPTest.HTTPTestOut.Code, t)
		th.checkKeyValues(th.ResponseBucket[HTTPTest.HTTPTestIn.TestCode], HTTPTest.HTTPTestOut.KeyValuesInBody, t)
		th.checkFields(th.ResponseBucket[HTTPTest.HTTPTestIn.TestCode], HTTPTest.HTTPTestOut.KeysPresentInBody, t)

	})
}
