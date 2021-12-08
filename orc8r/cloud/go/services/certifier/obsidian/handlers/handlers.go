/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/services/certifier/obsidian/models"
	"magma/orc8r/cloud/go/services/certifier/protos"
)

const (
	User             = "user"
	UserParam        = ":username"
	Tokens           = "tokens"
	TokenParam       = ":token"
	ListUser         = obsidian.V1Root + User
	ManageUser       = ListUser + obsidian.UrlSep + UserParam
	ListUserTokens   = ManageUser + obsidian.UrlSep + Tokens
	ManageUserTokens = ListUserTokens + obsidian.UrlSep + TokenParam
	Login            = ListUser + obsidian.UrlSep + "login"
)

func GetHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: ListUser, Methods: obsidian.GET, HandlerFunc: listUsersHandler},
		{Path: ListUser, Methods: obsidian.POST, HandlerFunc: createUserHandler},
		{Path: ManageUser, Methods: obsidian.GET, HandlerFunc: getUserHandler},
		{Path: ManageUser, Methods: obsidian.PUT, HandlerFunc: updateUserHandler},
		{Path: ManageUser, Methods: obsidian.DELETE, HandlerFunc: deleteUserHandler},
		{Path: ListUserTokens, Methods: obsidian.GET, HandlerFunc: getUserTokensHandler},
		{Path: ListUserTokens, Methods: obsidian.POST, HandlerFunc: addUserTokenHandler},
		{Path: ManageUserTokens, Methods: obsidian.DELETE, HandlerFunc: deleteUserTokenHandler},
		{Path: Login, Methods: obsidian.POST, HandlerFunc: loginHandler},
	}
	return ret
}

func listUsersHandler(c echo.Context) error {
	users, err := certifier.ListUsers(c.Request().Context())
	if err != nil {
		return obsidian.MakeHTTPError(err)
	}
	return c.JSON(http.StatusOK, users)
}

func createUserHandler(c echo.Context) error {
	data := &models.User{}
	if err := c.Bind(data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := data.Validate(strfmt.Default); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	username := fmt.Sprintf("%v", *data.Username)
	password := []byte(fmt.Sprintf("%v", *data.Password))
	user := &protos.User{
		Username: username,
		Password: password,
	}
	err := certifier.CreateUser(c.Request().Context(), user)
	return err
}

func getUserHandler(c echo.Context) error {
	username := c.Param("username")
	user, err := certifier.GetUser(c.Request().Context(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, user)
}

func updateUserHandler(c echo.Context) error {
	username := c.Param("username")
	var updatePassword string
	err := json.NewDecoder(c.Request().Body).Decode(&updatePassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error decoding request body for updating user: %v", err))
	}
	newUser := &protos.User{Username: username, Password: []byte(updatePassword)}
	certifier.UpdateUser(c.Request().Context(), newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return nil
}

func deleteUserHandler(c echo.Context) error {
	username := c.Param("username")
	deleteUser := &protos.User{Username: username}
	err := certifier.DeleteUser(c.Request().Context(), deleteUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error deleting user: %v", err))
	}
	return nil
}

func getUserTokensHandler(c echo.Context) error {
	username := c.Param("username")
	res, err := certifier.ListUserTokens(c.Request().Context(), &protos.User{Username: username})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to list user tokens: %v", err))
	}
	return c.JSON(http.StatusOK, res)
}

func addUserTokenHandler(c echo.Context) error {
	username := c.Param("username")

	resources := &models.Resources{}
	if err := c.Bind(resources); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := resources.Validate(strfmt.Default); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	resourceList := resourcesModelToProto(resources)
	req := &protos.AddUserTokenRequest{
		Username:  username,
		Resources: resourceList,
	}
	err := certifier.AddUserToken(c.Request().Context(), req)
	return err
}

func deleteUserTokenHandler(c echo.Context) error {
	username := c.Param("username")
	token := c.Param("token")
	req := &protos.DeleteUserTokenRequest{
		Username: username,
		Token:    token,
	}
	err := certifier.DeleteUserToken(c.Request().Context(), req)
	return err
}

func loginHandler(c echo.Context) error {
	data := &models.User{}
	if err := c.Bind(data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := data.Validate(strfmt.Default); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	username := fmt.Sprintf("%v", *data.Username)
	password := []byte(fmt.Sprintf("%v", *data.Password))
	user := &protos.User{
		Username: username,
		Password: password,
	}
	res, err := certifier.Login(c.Request().Context(), &protos.LoginRequest{User: user})
	if err != nil {
		return obsidian.MakeHTTPError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, res.Policies)
}

func resourcesModelToProto(resources *models.Resources) *protos.ResourceList {
	resourceList := make([]*protos.Resource, len(*resources))
	for i, resource := range *resources {
		resourceProto := &protos.Resource{
			Effect:       matchEffect(resource.Effect),
			Action:       matchAction(resource.Action),
			ResourceType: matchResourceType(resource.ResourceType),
			Resource:     resource.Resource,
		}
		resourceList[i] = resourceProto
	}
	return &protos.ResourceList{Resources: resourceList}
}

func matchResourceType(rawType string) protos.ResourceType {
	switch rawType {
	case protos.ResourceType_NETWORK_ID.String():
		return protos.ResourceType_NETWORK_ID
	case protos.ResourceType_TENANT_ID.String():
		return protos.ResourceType_TENANT_ID
	default:
		return protos.ResourceType_URI
	}
}

func matchEffect(rawEffect string) protos.Effect {
	switch rawEffect {
	case protos.Effect_DENY.String():
		return protos.Effect_DENY
	case protos.Effect_ALLOW.String():
		return protos.Effect_ALLOW
	default:
		return protos.Effect_UNKNOWN
	}
}

func matchAction(rawAction string) protos.Action {
	switch rawAction {
	case protos.Action_READ.String():
		return protos.Action_READ
	case protos.Action_WRITE.String():
		return protos.Action_WRITE
	default:
		return protos.Action_NONE
	}
}
