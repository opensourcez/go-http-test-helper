# Test Helper
### Now works for parallel execution (yay!)
#### only works for JSON or non-body requests
#### Only works for single layer json objects, not objects within object  ( this will be added soon )


# Import
import th "github.com/zkynet/testhelper"

# Variables to define
```
testInstance := th.TestHelper{
        // Map of responses from endpoints
    TestResultMap: make(map[string]map[string]*gabs.Container),
        // if this particular instance should log requests
    ShouldLog:     true,
        // cookies from http responses 
    Cookies:       make(map[string]*http.Cookie),
}
```

# Test results map
All response bodies are saved in a map with their TestCode as an index
```
th.TestResultMap[HTTPTest.HTTPTestIn.TestCode] = decodedBody
```
This way you can access them later
```
	loginUser := testInstance.TestResultMap["SomeTestCode"]
	testInstance.TestThis(th.NewHTTPTest(
		th.HTTPTestIn{
			Label: "LoginWith(USER-001)", TestCode: "USER-002",
			Body:   []byte(`{"email":"` + loginUser["email"].Data().(string) + `","password":"` + loginUser["password"].Data().(string) + `"}`),
			URL:    "https://localhost:3333/login",
			Method: "POST",
			Headers: map[string]string{
				"Content-Type":                "application/json",
			},
		},
		th.HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
			KeyValues: map[string]string{
				"name":  loginUser["name"].Data().(string),
				"email": loginUser["email"].Data().(string),
				"phone": loginUser["phone"].Data().(string),
			},
			KeyPresent: []string{"id", "disabled", "name", "phone", "email", "disabled"},
		}), t)
```


# Example Test
```
func TestUser(t *testing.T) {
    t.Parallel()
	
	testInstance.TestThis(th.NewHTTPTest(
		th.HTTPTestIn{
			Label: "CreateNewUser", TestCode: "USER-001",
			Body:    []byte(`{"name":"user1","phone":"5678","email":"user1@zkynet.io","password":"mowmeowm","facebook_id":"012345678"}`),
			URL:     "https://localhost:3333/user",
			Method:  "POST",
			Headers: headers,
		},
		th.HTTPTestOut{Code: 201, Status: "201 Created",
			KeyValues: map[string]string{
				"name":  "user1",
				"email": "user1@zkynet.io",
				"phone": "5678",
			},
			KeyPresent: []string{"id", "disabled", "name", "phone", "email"},
		}), t)


}
```

# Example error-free output
```
2018/04/30 16:42:13 ====== Sending method ( POST ) to =======
2018/04/30 16:42:13 http://localhost:3333/user
2018/04/30 16:42:13 ====== Received ( 201 Created ) =======
2018/04/30 16:42:13 {"id":"e10eddb7-4e02-4ed5-9105-df1ec903a8b7","name":"user1","phone":"5678","email":"user1@zkynet.io","disabled":false}
2018/04/30 16:42:13 ===============================================
```


# Example test with errors
```
func TestUser(t *testing.T) {
    t.Parallel()

	testInstance.TestThis(th.NewHTTPTest(
		th.HTTPTestIn{
			Label: "CreateNewUser", TestCode: "USER-001",
			Body:    []byte(`{"name":"user1","phone":"5678","email":"user1@zkynet.io","password":"mowmeowe","facebook_id":"012345678"}`),
			URL:     "https://localhost:3333/user",
			Method:  "POST",
			Headers: headers,
		},
		th.HTTPTestOut{Code: 201, Status: "201 Created",
			KeyValues: map[string]string{
				"name":  "user1",
				"email": "user@zkynet.io",
				"phone": "5678",
			},
			KeyPresent: []string{"id", "disabled", "name", "phone", "email"},
		}), t)

}
```

# Example output with errors
```
--- FAIL: TestUser (0.33s)
    --- FAIL: TestUser/USER-001:CreateNewUser (0.08s)
        helper.go:169: Expected ( user@zkynet.io ) in key ( email ) but got ( user1@zkynet.io )


```


# Tessting
The test unit will spawn a web sever on port 3333
```
go test -v
```