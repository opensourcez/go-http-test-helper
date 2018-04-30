# Test Helper
#### only works for JSON or non-body requests
#### Only works for single layer json objects, not objects within object  ( this will be added soon )

# Import
import th "github.com/zkynet/testhelper"

# Test results map
All response bodies are saved in a map with their TestCode as an index
```
TestResultMap[HTTPTest.HTTPTestIn.TestCode] = decodedBody
```
This way you can access them later
```
	loginUser = th.TestResultMap["USER-005"]
	th.TestThis(th.NewHTTPTest(
		th.HTTPTestIn{
			Label: "LoginWith(USER-005)", TestCode: "USER-014",
			Body:   []byte(`{"email":"` + loginUser["email"].Data().(string) + `","password":"` + userPasswordPostCreate + `"}`),
			URL:    os.Getenv("Url") + ":" + os.Getenv("Port") + "/normal/login",
			Method: "POST",
			Headers: map[string]string{
				os.Getenv("RequestTokenName"): os.Getenv("RequestToken"),
				"Content-Type":                "application/json",
			},
		},
		th.HTTPTestOut{Body: "", Code: 200, Status: "200 OK",
			KeyValues: []*th.KeyValueInBody{
				&th.KeyValueInBody{Key: "name", Value: loginUser["name"].Data().(string)},
				&th.KeyValueInBody{Key: "email", Value: loginUser["email"].Data().(string)},
				&th.KeyValueInBody{Key: "phone", Value: loginUser["phone"].Data().(string)},
			},
			KeyPresent: []string{"id", "disabled", "name", "phone", "email", "disabled"},
		}), t)
```


# Example Test
```
func TestUser(t *testing.T) {
	userPasswordPostCreate := "meowmeow"
	th.TestThis(th.NewHTTPTest(
		th.HTTPTestIn{
			Label: "CreateNewUser", TestCode: "USER-001",
			Body:   []byte(`{"name":"user1","phone":"5678","email":"user1@zkynet.io","password":"` + userPasswordPostCreate + `","facebook_id":"012345678"}`),
			URL:    os.Getenv("Url") + ":" + os.Getenv("Port") + "/user",
			Method: "POST",
			Headers: map[string]string{
				os.Getenv("RequestTokenName"): os.Getenv("RequestToken"),
				"Content-Type":                "application/json",
			},
		},
		th.HTTPTestOut{Code: 201, Status: "201 Created",
			KeyValues: []*th.KeyValueInBody{
				&th.KeyValueInBody{Key: "name", Value: "user1"},
				&th.KeyValueInBody{Key: "email", Value: "user1@zkynet.io"},
				&th.KeyValueInBody{Key: "phone", Value: "5678"},
			},
			KeyPresent: []string{"id", "disabled", "name", "phone", "email"},
		}), t)
}
```

# Example error-free output
```
2018/04/30 12:45:58 ====== Sending method ( POST ) to =======
2018/04/30 12:45:58 http://localhost:1337/user
2018/04/30 12:45:58 ====== Received ( 201 Created ) =======
2018/04/30 12:45:58 {"id":"8dd5badd-96c1-413a-9f34-b2ac7a736e1a","name":"user1","phone":"5678","email":"user1@zkynet.io","disabled":false}
```


# Example test with errors
```
func TestUser(t *testing.T) {
	th.TestThis(th.NewHTTPTest(
		th.HTTPTestIn{
			Label: "CreateNewUser", TestCode: "USER-001",
			Body:   []byte(`{"name":"user1","phone":"5678","email":"user1@zkynet.io","password":"` + userPasswordPostCreate + `","facebook_id":"012345678"}`),
			URL:    os.Getenv("Url") + ":" + os.Getenv("Port") + "/user",
			Method: "POST",
			Headers: map[string]string{
				os.Getenv("RequestTokenName"): os.Getenv("RequestToken"),
				"Content-Type":                "application/json",
			},
		},
		th.HTTPTestOut{Code: 201, Status: "201 Created",
			KeyValues: []*th.KeyValueInBody{
				&th.KeyValueInBody{Key: "name", Value: "user1"},
				&th.KeyValueInBody{Key: "email", Value: "user@zkynet.io"},
				&th.KeyValueInBody{Key: "phone", Value: "5678"},
			},
			KeyPresent: []string{"id", "disabled", "name", "email"},
		}), t)
}
```

# Example output with errors
```
--- FAIL: TestUser (0.35s)
    --- FAIL: TestUser/USER-001:CreateNewUser (0.07s)
        helpers.go:170: Expected ( user@zkynet.io ) in key ( email ) but got ( user1@zkynet.io )
        helpers.go:145: Key (phone) is not suppose to be in response
FAIL
```