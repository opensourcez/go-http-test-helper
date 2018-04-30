# Test Helper
## only works for JSON or non-body requests
## Only works for single layer json objects, not objects within object  ( this will be added soon )

## import
import th "github.com/zkynet/testhelper"

# INPUT
```
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
```

# OUTPUT
```
2018/04/30 12:45:58 ====== Sending method ( POST ) to =======
2018/04/30 12:45:58 http://localhost:1337/user
2018/04/30 12:45:58 ====== Received ( 201 Created ) =======
2018/04/30 12:45:58 {"id":"8dd5badd-96c1-413a-9f34-b2ac7a736e1a","name":"user1","phone":"5678","email":"user1@zkynet.io","disabled":false}
```


# INPUT WITH ERRORS
```
# IN
```
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
```
```

# OUTPUT WITH ERRORS
```
--- FAIL: TestUser (0.35s)
    --- FAIL: TestUser/USER-001:CreateNewUser (0.07s)
        helpers.go:170: Expected ( user@zkynet.io ) in key ( email ) but got ( user1@zkynet.io )
        helpers.go:145: Key (phone) is not suppose to be in response
FAIL
```