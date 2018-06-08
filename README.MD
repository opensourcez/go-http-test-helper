
![alt text](./gopher.png)
# Test Helper
This package will save you some time when doing http testing. 

## installation
```
go get github.com/zkynet/testhelper
```

## What can this package do for you ?
 - Compare raw response bodies
 - Compare JSON response objects
    - the support for JSON only goes one level deep. If you want to go deeper use the raw body comparison
- Compare outgoing and incoming headers
- Compare key/value in response body (json-only)
- Check if key's exist in response body
    - sometimes you can't predict the return value of a field but you want to make sure it's there
- Compare response status and code

## Extras
 - Optional logging of http requests and responses
 - Optional error and success logging 
 - Custom colors for logging

## Conditions
If the HTTPTestOut.RawBody is present the test helper will not perform JSON object dependant checks.
```
if HTTPTest.HTTPTestOut.RawBody != "" {
    th.checkRawBody(string(body), HTTPTest.HTTPTestOut.RawBody, t)
} else {
    th.ResponseBucket[HTTPTest.HTTPTestIn.TestCode] = th.decodeBody(body, t)
    th.checkKeyValues(th.ResponseBucket[HTTPTest.HTTPTestIn.TestCode], HTTPTest.HTTPTestOut.KeyValues, t)
    th.checkFields(th.ResponseBucket[HTTPTest.HTTPTestIn.TestCode], HTTPTest.HTTPTestOut.Keys, t)
}
```



# How it works
The testhelper has an exported method called TestThis. TestThis will accept a HTTPTest struct wich contains HTTPTestIn and HTTPTestOut. 

1. Make a new <b>TestHelper</b> instance
```
package test 

import th "github.com/zkynet/testhelper"

func main(){
    errorColor := "\033[0m"
    successColor := "\033[0m"
    infoColor := "\033[0m"
    enableLogging := true
    testHelper := th.NewHTTPTestHelper(enableLogging, errorColor, successColor, infoColor)
}
```

2. Make an <b>HTTPTest</b> struct with a HTTPTestIn and a HTTPTestOut. 
```
type HTTPTest struct {
    HTTPTestIn
    HTTPTestOut
}

type HTTPTestIn struct {
    Note     string // A custom note that is printed if logging is enabled.
    Label    string // The label for your test.. testhelper will perform a t.run(label)
    TestCode string // The testhelper uses this code to store data from response bodies. (more on this in the examples)
    Body     []byte // The request body (optional)
    URL      string // The endpoint you are testing
    Method   string // Endpoint method
    Headers  map[string]string // Headers that you want to send (optional)
}

type HTTPTestOut struct {
    RawBody           string // A raw response body for comparison (optional)
    KeyValues         map[string]string // Key/values to check for in the response body (optional)
    Keys              []string // Keys to check for in the response body (optional)
    Status            string // Expected response status
    Code              int // Excpected response code
    Headers           map[string]string // Excpected response headers (optional)
    HeadersToIgnore   []string // Headers to ignore. If any. (optional)
}

```

2. Once the HTTPTest is constructed call the TestThis method:  <b>testhelper.TestThis(HTTPTest, t) </b>.. t is the *testing.T variable.

3. The test helper will now perform all the checks that it can depending on the input you gave it.  
 
4. The response Body and Cookies will be saved in the <b>ResponseBucket</b> and <b>CookieBucket</b> ( if you are working with a json object )

```
type TestHelper struct {
	ShouldLog      bool
	CookieBucket   map[string]*http.Cookie
	ResponseBucket map[string]map[string]*gabs.Container
	ErrorColor     string
	SuccessColor   string
	InfoColor      string
}
```

5. Access the response object using the <b>ResponseBucket</b>.
    - testhelper.ResponseBucket["TestCode"]["KEY"].Data().(type)
    - This bucket is a <b>Gabs Container: "github.com/Jeffail/gabs"</b>


# How to (as simple as possible)
## Initialize a new test helper
```
package test 

import th "github.com/zkynet/testhelper"

testHelper := th.NewHTTPTestHelper(true,"","","")
```
- To turn off any type of coloring just replace the color name with "none"

## Make a base test (Key/Value)

```
t.Run("SeriesOfBaseTests", func(t *testing.T) {
    t.Parallel()

testHelper := th.NewHTTPTestHelper(true, "", "", "")


baseHTTPTest := &HTTPTest{
    HTTPTestIn{
        Note:  "This is a custom note for this test, it will show up in the test output if logging is enabled..",
        Label: "BaseTestForKeysAndValues", TestCode: "BASE-001",
        Body:    []byte(`{"KEY1":"VALUE1","KEY2":"VALUE2"}`),
        URL:      "http://your-domain.com/test",
        Method:  "POST",
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    },
    HTTPTestOut{
        Code: 200, 
        Status: "200 OK",
        KeyValues: map[string]string{
            "KEY1": "VALUE1",
            "KEY2": "VALUE2",
        },
    }
}

// you can also just construct the HTTPTest inside the method.
testHelper.TestThis(baseHTTPTest, t)

})
```

## Access data from the previous test
The first index in the ResponseBucket is the TestCode, the second index is the Key in your response. 
```
testHelper.ResponseBucket["BASE-001"]["KEY1"].Data().(string)
// output: VALUE1

testHelper.ResponseBucket["BASE-001"]["KEY2"].Data().(string)
// output: VALUE2
```


## Make a base test ( Key check only )

```
&HTTPTest{
    HTTPTestIn{
        Note:  "This is a custom note for this test, it will show up in the test output if logging is enabled..",
        Label: "BaseTestKeysOnly", TestCode: "BASE-002",
        Body:    []byte(`{"KEY1":"VALUE1","KEY2":"VALUE2"}`),
        URL:      "http://your-domain.com/test",
        Method:  "POST",
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    },
    HTTPTestOut{
        Code: 200, 
        Status: "200 OK",
        Keys: []string{"KEY1", "KEY2"},
    }
}
```




## Full JSON test

```
&HTTPTest{
    HTTPTestIn{
        Note:  "This is a custom note for this test, it will show up in the test output if logging is enabled..",
        Label: "BaseTestFullJSON", TestCode: "BASE-003",
        Body:    []byte(`{"KEY1":"VALUE1","KEY2":"VALUE2"}`),
        URL:      "http://your-domain.com/test",
        Method:  "POST",
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    },
    HTTPTestOut{
        Code: 200, 
        Status: "200 OK",
        KeyValues: map[string]string{
            "KEY1":  "VALUE1",
            "KEY2":  "VALUE2",
        },
        Keys: []string{"KEY1", "KEY2"},
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }
}

```

## Full Raw Body test

```
&HTTPTest{
    HTTPTestIn{
        Note:  "This is a custom note for this test, it will show up in the test output if logging is enabled..",
        Label: "BaseTestRawBody", TestCode: "BASE-004",
        Body:    []byte(`{"KEY1":"VALUE1","KEY2":"VALUE2"}`),
        URL:      "http://your-domain.com/test",
        Method:  "POST",
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    },
    HTTPTestOut{
        Code: 200, 
        Status: "200 OK",
        RawBody:    []byte(`{"KEY1":"VALUE1","KEY2":"VALUE2"}`),
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }
}

```

## Ignoring headers

```
&HTTPTest{
    HTTPTestIn{
        Note:  "This is a custom note for this test, it will show up in the test output if logging is enabled..",
        Label: "BaseTestIgnoreHeaders", TestCode: "BASE-005",
        Body:    []byte(`{"KEY1":"VALUE1","KEY2":"VALUE2"}`),
        URL:      "http://your-domain.com/test",
        Method:  "POST",
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    },
    HTTPTestOut{
        Code: 200, 
        Status: "200 OK",
        RawBody:    []byte(`{"KEY1":"VALUE1","KEY2":"VALUE2"}`),
        Headers: map[string]string{
            "Content-Type": "application/json",
            "X-CUSTOM-HEADER": "super-secret-header-value",
        },
        IgnoredHeaders: []string{"Header-To-Ignore"},
    }
}

```