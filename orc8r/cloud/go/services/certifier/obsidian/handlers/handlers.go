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

	"github.com/labstack/echo"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/services/certifier/obsidian/models"
	"magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/storage"
)

const (
	UserParam        = ":username"
	User             = "user"
	Tokens           = "tokens"
	ListUser         = obsidian.V1Root + User
	ManageUser       = ListUser + obsidian.UrlSep + UserParam
	ManageUserTokens = ManageUser + obsidian.UrlSep + Tokens
	Login            = ListUser + obsidian.UrlSep + "login"
)

func GetHandlers(storage storage.CertifierStorage) []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: ListUser, Methods: obsidian.GET, HandlerFunc: listUsersHandler},
		{Path: ListUser, Methods: obsidian.POST, HandlerFunc: createUserHandler},
		{Path: ManageUser, Methods: obsidian.GET, HandlerFunc: getUserHandler},
		{Path: ManageUser, Methods: obsidian.PUT, HandlerFunc: updateUserHandler},
		{Path: ManageUser, Methods: obsidian.DELETE, HandlerFunc: deleteUserHandler},
		{Path: Login, Methods: obsidian.POST, HandlerFunc: loginHandler},
		{Path: ManageUserTokens, Methods: obsidian.GET, HandlerFunc: getUserTokensHandler},
		{Path: ManageUserTokens, Methods: obsidian.POST, HandlerFunc: addUserTokenHandler},
		{Path: ManageUserTokens, Methods: obsidian.DELETE, HandlerFunc: deleteUserTokenHandler},
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
	var data models.User
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error decoding request body for creating user: %v", err))
	}
	username := fmt.Sprintf("%v", *data.Username)
	password := []byte(fmt.Sprintf("%v", *data.Password))
	user := &protos.User{
		Username: username,
		Password: password,
	}
	err = certifier.CreateUser(c.Request().Context(), user)
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
	// username := c.Param("username")
	// var updatedPassword string
	// err := json.NewDecoder(c.Request().Body).Decode(updatedPassword)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error decoding request body for updating user"))
	// }
	return nil
}

func deleteUserHandler(c echo.Context) error {
	return nil
}

func loginHandler(c echo.Context) error {
	return nil
}
func getUserTokensHandler(c echo.Context) error {
	return nil
}
func addUserTokenHandler(c echo.Context) error {
	return nil
}
func deleteUserTokenHandler(c echo.Context) error {
	return nil
}

// func getUpdateUserHandler(storage storage.CertifierStorage) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		username := c.Param("username")
// 		data := make(map[string]interface{})
// 		err := json.NewDecoder(c.Request().Body).Decode(&data)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error decoding request for User: %v", err))
// 		}
//
// 		newPassword := []byte(fmt.Sprintf("%v", data["password"]))
// 		hashedPassword, err := bcrypt.GenerateFromPassword(newPassword, bcrypt.DefaultCost)
//
// 		user, err := storage.GetUser(username)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting existing user: %v", err))
// 		}
//
// 		newUser := &protos.User{
// 			Username: username,
// 			Password: hashedPassword,
// 			Tokens:   user.Tokens,
// 		}
// 		storage.PutUser(username, newUser)
//
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, err)
// 		}
// 		return nil
// 	}
// }
//
// func getCreateUserHandler(storage storage.CertifierStorage) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var data models.UserWithPolicy
// 		err := json.NewDecoder(c.Request().Body).Decode(&data)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error decoding request body for creating user: %v", err))
// 		}
//
// 		username := fmt.Sprintf("%v", *data.User.Username)
// 		password := []byte(fmt.Sprintf("%v", *data.User.Password))
// 		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error hashing password: %v", err))
// 		}
//
// 		token, err := certifier.GenerateToken(certifier.Personal)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error generating personal access token for user: %v", err))
// 		}
// 		user := &protos.User{
// 			Username: username,
// 			Password: hashedPassword,
// 			Tokens:   &protos.TokenList{Tokens: []string{token}},
// 		}
// 		if err = storage.PutUser(username, user); err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, err)
// 		}
// 		effect := matchEffect(data.Policy.Effect)
// 		action := matchAction(data.Policy.Action)
// 		resource := &protos.Resource{
// 			Effect:       effect,
// 			Action:       action,
// 			ResourceType: protos.ResourceType_URI,
// 			Resource:     "/**",
// 		}
// 		policy := &protos.Policy{
// 			Token: token,
// 			Resources: &protos.ResourceList{
// 				Resources: []*protos.Resource{resource},
// 			},
// 		}
// 		if err = storage.PutPolicy(token, policy); err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, err)
// 		}
//
// 		return nil
// 	}
// }
//
// func getDeleteUserHandler(storage storage.CertifierStorage) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		username := c.Param("username")
// 		err := storage.DeleteUser(username)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, err)
// 		}
// 		return nil
// 	}
// }
//
// func getLoginHandler(storage storage.CertifierStorage) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		data := make(map[string]interface{})
// 		err := json.NewDecoder(c.Request().Body).Decode(&data)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error decoding request for User: %v", err))
// 		}
//
// 		username := fmt.Sprintf("%v", data["username"])
// 		password := []byte(fmt.Sprintf("%v", data["password"]))
// 		user, err := storage.GetUser(username)
// 		if err != nil {
// 			return obsidian.MakeHTTPError(err, http.StatusInternalServerError)
// 		}
//
// 		hashedPassword := user.Password
// 		err = bcrypt.CompareHashAndPassword(hashedPassword, password)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusUnauthorized, "wrong password")
// 		}
//
// 		tokens := user.Tokens.Tokens
// 		return c.JSON(http.StatusOK, tokens)
// 	}
// }
//
// func matchEffect(rawEffect *string) protos.Effect {
// 	switch *rawEffect {
// 	case protos.Effect_DENY.String():
// 		return protos.Effect_DENY
// 	case protos.Effect_ALLOW.String():
// 		return protos.Effect_ALLOW
// 	default:
// 		return protos.Effect_UNKNOWN
// 	}
// }
//
// func matchAction(rawAction *string) protos.Action {
// 	switch *rawAction {
// 	case protos.Action_READ.String():
// 		return protos.Action_READ
// 	case protos.Action_WRITE.String():
// 		return protos.Action_WRITE
// 	default:
// 		return protos.Action_NONE
// 	}
// }
