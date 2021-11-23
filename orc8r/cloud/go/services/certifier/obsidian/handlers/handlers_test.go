package handlers_test

import (
	"testing"

	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"

	"github.com/labstack/echo"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/services/certifier/obsidian/handlers"
	"magma/orc8r/cloud/go/services/certifier/test_utils"
)

const RootUsername = "root"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Policy struct {
	Effect    string   `json:"effect"`
	Action    string   `json:"action"`
	Resources []string `json:"resource"`
}

type CreateUserRequest struct {
	User   `json:"user"`
	Policy `json:"policy"`
}

func TestHTTPBasicAuth(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	store := test_utils.GetCertifierBlobstore(t)
	handlers := handlers.GetHandlers(store)
	testURLRoot := "/magma/v1/http_basic_auth"
	listHTTPBasicAuth := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc
	createHTTPBasicAuth := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.POST).HandlerFunc
	// TODO(christinewang5): test these later
	testUpdateDeleteRoot := "/magma/v1/http_basic_auth/:username"
	deleteHTTPBasicAuth := tests.GetHandlerByPathAndMethod(t, handlers, testUpdateDeleteRoot, obsidian.DELETE).HandlerFunc
	// updateHTTPBasicAuth := tests.GetHandlerByPathAndMethod(t, handlers, testUpdateDeleteRoot, obsidian.PUT).HandlerFunc

	e := echo.New()

	// create bob user
	createBobRequest := CreateUserRequest{
		User: User{
			Username: "bob",
			Password: "password",
		},
		Policy: Policy{
			Effect:    "ALLOW",
			Action:    "READ",
			Resources: []string{"*"},
		},
	}
	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(createBobRequest),
		Handler:        createHTTPBasicAuth,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// create root user
	createRootRequest := CreateUserRequest{
		User: User{
			Username: RootUsername,
			Password: "password",
		},
		Policy: Policy{
			Effect:    "ALLOW",
			Action:    "WRITE",
			Resources: []string{"*"},
		},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(createRootRequest),
		Handler:        createHTTPBasicAuth,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listHTTPBasicAuth,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"bob", "root"}),
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "DELETE",
		URL:            testUpdateDeleteRoot,
		Handler:        deleteHTTPBasicAuth,
		ParamNames:     []string{"username"},
		ParamValues:    []string{"bob"},
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)
}
