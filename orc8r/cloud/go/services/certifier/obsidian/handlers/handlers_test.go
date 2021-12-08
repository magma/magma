package handlers_test

import (
	"testing"

	"github.com/labstack/echo"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/services/certifier/obsidian/handlers"
	"magma/orc8r/cloud/go/services/certifier/obsidian/models"
	"magma/orc8r/cloud/go/services/certifier/test_utils"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
)

func TestUserEndpoints(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	handlers := handlers.GetHandlers()
	testUser := "/magma/v1/user"
	testManageUser := testUser + "/:username"
	testManageUserTokens := testManageUser + "/tokens"
	testLogin := testUser + "/login"
	listUser := tests.GetHandlerByPathAndMethod(t, handlers, testUser, obsidian.GET).HandlerFunc
	createUser := tests.GetHandlerByPathAndMethod(t, handlers, testUser, obsidian.POST).HandlerFunc
	deleteUser := tests.GetHandlerByPathAndMethod(t, handlers, testManageUser, obsidian.DELETE).HandlerFunc
	updateUser := tests.GetHandlerByPathAndMethod(t, handlers, testManageUser, obsidian.PUT).HandlerFunc
	getUserTokens := tests.GetHandlerByPathAndMethod(t, handlers, testManageUserTokens, obsidian.GET).HandlerFunc
	addUserTokens := tests.GetHandlerByPathAndMethod(t, handlers, testManageUserTokens, obsidian.POST).HandlerFunc
	login := tests.GetHandlerByPathAndMethod(t, handlers, testLogin, obsidian.POST).HandlerFunc

	// TODO(christinewang5): is it possible to get the response from RunUnitTest?
	// deleteUserTokens := tests.GetHandlerByPathAndMethod(t, handlers, testManageUserTokens, obsidian.DELETE).HandlerFunc

	e := echo.New()

	// Create user endpoints
	username := test_utils.TestUsername
	password := test_utils.TestPassword
	createBobRequest := &models.User{
		Username: &username,
		Password: &password,
	}
	tc := tests.Test{
		Method:         "POST",
		URL:            testUser,
		Payload:        tests.JSONMarshaler(createBobRequest),
		Handler:        createUser,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	username = test_utils.TestRootUsername
	createRootRequest := &models.User{
		Username: &username,
		Password: &password,
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            testUser,
		Payload:        tests.JSONMarshaler(createRootRequest),
		Handler:        createUser,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            testUser,
		Handler:        listUser,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "PUT",
		URL:            testManageUser,
		Handler:        updateUser,
		ParamNames:     []string{"username"},
		ParamValues:    []string{test_utils.TestUsername},
		Payload:        tests.JSONMarshaler("newPassword"),
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "DELETE",
		URL:            testManageUser,
		Handler:        deleteUser,
		ParamNames:     []string{"username"},
		ParamValues:    []string{test_utils.TestUsername},
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// Test token endpoints
	writeAllResource := []*models.Resource{
		{
			Effect:       models.ResourceEffectALLOW,
			Action:       models.ResourceActionWRITE,
			ResourceType: models.ResourceResourceTypeURI,
			Resource:     "**",
		},
		{
			Effect:       models.ResourceEffectALLOW,
			Action:       models.ResourceActionWRITE,
			ResourceType: models.ResourceResourceTypeURI,
			Resource:     "**",
		},
		{
			Effect:       models.ResourceEffectALLOW,
			Action:       models.ResourceActionWRITE,
			ResourceType: models.ResourceResourceTypeURI,
			Resource:     "**",
		},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            testManageUserTokens,
		Handler:        addUserTokens,
		ParamNames:     []string{"username"},
		ParamValues:    []string{test_utils.TestRootUsername},
		Payload:        tests.JSONMarshaler(writeAllResource),
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            testManageUserTokens,
		ParamNames:     []string{"username"},
		ParamValues:    []string{test_utils.TestRootUsername},
		Handler:        getUserTokens,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// Test login endpoints
	tc = tests.Test{
		Method:         "POST",
		URL:            testLogin,
		Payload:        tests.JSONMarshaler(createBobRequest),
		Handler:        login,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)
	badPassword := "blah"
	badPayload := &models.User{
		Username: &username,
		Password: &badPassword,
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            testLogin,
		Payload:        tests.JSONMarshaler(badPayload),
		Handler:        login,
		ExpectedStatus: 500,
		ExpectedError:  "wrong password",
	}
	tests.RunUnitTest(t, e, tc)
}
