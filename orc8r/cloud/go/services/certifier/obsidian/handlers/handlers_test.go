package handlers_test

import (
	"testing"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/services/certifier/obsidian/handlers"
	"magma/orc8r/cloud/go/services/certifier/test_utils"

	"github.com/labstack/echo"
)

const ROOT_USERNAME = "root"

func TestHTTPBasicAuth(t *testing.T) {
	store := test_utils.GetCertifierBlobstore(t)
	handlers := handlers.GetHandlers(store)
	testURLRoot := "/magma/v1/http_basic_auth"
	listHTTPBasicAuth := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc
	createHTTPBasicAuth := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.POST).HandlerFunc
	// TODO(christinewang5): test these later
	// testUpdateDeleteRoot := "/magma/v1/http_basic_auth/:username"
	// deleteHTTPBasicAuth := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.DELETE).HandlerFunc
	// updateHTTPBasicAuth := tests.GetHandlerByPathAndMethod(t, handlers, testUpdateDeleteRoot, obsidian.PUT).HandlerFunc

	bob, _ := test_utils.CreateTestUser(t, "bob", "password")
	root, _ := test_utils.CreateTestUser(t, ROOT_USERNAME, "password")

	e := echo.New()

	// create bob user
	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(bob),
		Handler:        createHTTPBasicAuth,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// create root user
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(root),
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

}
