package handlers_test

import (
	"testing"

	"github.com/labstack/echo"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/services/certifier/obsidian/handlers"
	"magma/orc8r/cloud/go/services/certifier/test_utils"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
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

func TestUserEndpoints(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	store := test_utils.GetCertifierBlobstore(t)
	handlers := handlers.GetHandlers(store)
	testURLRoot := "/magma/v1/user"
	listUser := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc
	createUser := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.POST).HandlerFunc
	testUpdateDeleteRoot := "/magma/v1/user/:username"
	deleteUser := tests.GetHandlerByPathAndMethod(t, handlers, testUpdateDeleteRoot, obsidian.DELETE).HandlerFunc
	updateUser := tests.GetHandlerByPathAndMethod(t, handlers, testUpdateDeleteRoot, obsidian.PUT).HandlerFunc

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
		Handler:        createUser,
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
		Handler:        createUser,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listUser,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"bob", "root"}),
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:      "PUT",
		URL:         testUpdateDeleteRoot,
		Handler:     updateUser,
		ParamNames:  []string{"username"},
		ParamValues: []string{"bob"},
		Payload: tests.JSONMarshaler(struct {
			Password string
		}{Password: "newPassword"}),
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "DELETE",
		URL:            testUpdateDeleteRoot,
		Handler:        deleteUser,
		ParamNames:     []string{"username"},
		ParamValues:    []string{"bob"},
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)
}
