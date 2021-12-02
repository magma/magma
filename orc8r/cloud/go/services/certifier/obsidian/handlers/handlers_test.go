package handlers_test

//
// import (
// 	"testing"
//
// 	"github.com/labstack/echo"
//
// 	"magma/orc8r/cloud/go/obsidian"
// 	"magma/orc8r/cloud/go/obsidian/tests"
// 	"magma/orc8r/cloud/go/services/certifier/obsidian/handlers"
// 	"magma/orc8r/cloud/go/services/certifier/obsidian/models"
// 	"magma/orc8r/cloud/go/services/certifier/test_utils"
// 	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
// )
//
// // TODO(christinewang5): fix all these darn tests later ugh
// func TestUserEndpoints(t *testing.T) {
// 	configuratorTestInit.StartTestService(t)
// 	store := test_utils.GetCertifierBlobstore(t)
// 	handlers := handlers.GetHandlers(store)
// 	testURLRoot := "/magma/v1/user"
// 	listUser := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc
// 	createUser := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.POST).HandlerFunc
// 	testUpdateDeleteRoot := "/magma/v1/user/:username"
// 	deleteUser := tests.GetHandlerByPathAndMethod(t, handlers, testUpdateDeleteRoot, obsidian.DELETE).HandlerFunc
// 	updateUser := tests.GetHandlerByPathAndMethod(t, handlers, testUpdateDeleteRoot, obsidian.PUT).HandlerFunc
//
// 	e := echo.New()
//
// 	// Create non-admin user
// 	username := test_utils.TestUsername
// 	password := test_utils.TestPassword
// 	resource := &models.Resource{
// 		Action:       models.ResourceActionREAD,
// 		Effect:       models.ResourceEffectALLOW,
// 		ResourceType: models.ResourceResourceTypeURI,
// 		Resources:    "/**",
// 	}
// 	action := models.ResourceActionREAD
// 	createBobRequest := models.UserWithPolicy{
// 		User: &models.User{
// 			Username: &username,
// 			Password: &password,
// 		},
// 		Policy: &models.Policy{
// 			Effect:    &effect,
// 			Action:    &action,
// 			Resources: []*models.Resources{resource},
// 		},
// 	}
// 	tc := tests.Test{
// 		Method:         "POST",
// 		URL:            testURLRoot,
// 		Payload:        tests.JSONMarshaler(createBobRequest),
// 		Handler:        createUser,
// 		ExpectedStatus: 200,
// 	}
// 	tests.RunUnitTest(t, e, tc)
//
// 	// Create root user request
// 	username = test_utils.TestRootUsername
// 	effect = models.ResourceEffectALLOW
// 	action = models.ResourceActionWRITE
// 	createRootRequest := models.UserWithPolicy{
// 		User: &models.User{
// 			Username: &username,
// 			Password: &password,
// 		},
// 		Policy: &models.Policy{
// 			Effect:    &effect,
// 			Action:    &action,
// 			Resources: []string{"/**"},
// 		},
// 	}
// 	tc = tests.Test{
// 		Method:         "POST",
// 		URL:            testURLRoot,
// 		Payload:        tests.JSONMarshaler(createRootRequest),
// 		Handler:        createUser,
// 		ExpectedStatus: 200,
// 	}
// 	tests.RunUnitTest(t, e, tc)
//
// 	tc = tests.Test{
// 		Method:         "GET",
// 		URL:            testURLRoot,
// 		Handler:        listUser,
// 		ExpectedStatus: 200,
// 	}
// 	tests.RunUnitTest(t, e, tc)
//
// 	tc = tests.Test{
// 		Method:      "PUT",
// 		URL:         testUpdateDeleteRoot,
// 		Handler:     updateUser,
// 		ParamNames:  []string{"username"},
// 		ParamValues: []string{"bob"},
// 		Payload: tests.JSONMarshaler(struct {
// 			Password string
// 		}{Password: "newPassword"}),
// 		ExpectedStatus: 200,
// 	}
// 	tests.RunUnitTest(t, e, tc)
//
// 	tc = tests.Test{
// 		Method:         "DELETE",
// 		URL:            testUpdateDeleteRoot,
// 		Handler:        deleteUser,
// 		ParamNames:     []string{"username"},
// 		ParamValues:    []string{"bob"},
// 		ExpectedStatus: 200,
// 	}
// 	tests.RunUnitTest(t, e, tc)
// }
