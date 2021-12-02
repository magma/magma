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
	"golang.org/x/crypto/bcrypt"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/storage"
)

const (
	UserParam  = ":username"
	User       = "user"
	ListUser   = obsidian.V1Root + User
	ManageUser = ListUser + obsidian.UrlSep + UserParam
	Login      = ListUser + obsidian.UrlSep + "login"
)

func GetHandlers(storage storage.CertifierStorage) []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: ListUser, Methods: obsidian.GET, HandlerFunc: getListUserHandler(storage)},
		{Path: ListUser, Methods: obsidian.POST, HandlerFunc: getCreateUserHandler(storage)},
		{Path: ManageUser, Methods: obsidian.PUT, HandlerFunc: getUpdateUserHandler(storage)},
		{Path: ManageUser, Methods: obsidian.DELETE, HandlerFunc: getDeleteUserHandler(storage)},
		{Path: Login, Methods: obsidian.POST, HandlerFunc: getLoginHandler(storage)},
	}
	return ret
}

func getListUserHandler(storage storage.CertifierStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := storage.ListUser()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, users)
	}
}

type CreateUserRequest struct {
	User struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"user"`
	Policy struct {
		Effect    string   `json:"effect"`
		Action    string   `json:"action"`
		Resources []string `json:"resource"`
	} `json:"policy"`
}

func getCreateUserHandler(storage storage.CertifierStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		var data CreateUserRequest
		err := json.NewDecoder(c.Request().Body).Decode(&data)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error decoding request body for creating user: %v", err))
		}

		username := fmt.Sprintf("%v", data.User.Username)
		password := []byte(fmt.Sprintf("%v", data.User.Password))
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error hashing password: %v", err))
		}

		token, err := certifier.GenerateToken(certifier.Personal)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error generating personal access token for user: %v", err))
		}

		user := &protos.User{
			Username: username,
			Password: hashedPassword,
			Tokens:   &protos.TokenList{Token: []string{token}},
		}
		if err = storage.PutUser(username, user); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var effect protos.Effect
		switch data.Policy.Effect {
		case protos.Effect_DENY.String():
			effect = protos.Effect_DENY
		case protos.Effect_ALLOW.String():
			effect = protos.Effect_ALLOW
		default:
			effect = protos.Effect_UNKNOWN
		}
		var action protos.Action
		switch data.Policy.Action {
		case protos.Action_READ.String():
			action = protos.Action_READ
		case protos.Action_WRITE.String():
			action = protos.Action_WRITE
		default:
			action = protos.Action_NONE
		}
		policy := &protos.Policy{
			Token:  token,
			Effect: effect,
			Action: action,
			Resources: &protos.ResourceList{
				Resource: data.Policy.Resources,
			},
		}
		if err = storage.PutPolicy(token, policy); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return nil
	}
}

func getUpdateUserHandler(storage storage.CertifierStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param("username")
		data := make(map[string]interface{})
		err := json.NewDecoder(c.Request().Body).Decode(&data)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error decoding request for User: %v", err))
		}

		newPassword := []byte(fmt.Sprintf("%v", data["password"]))
		hashedPassword, err := bcrypt.GenerateFromPassword(newPassword, bcrypt.DefaultCost)

		user, err := storage.GetUser(username)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting existing user: %v", err))
		}

		newUser := &protos.User{
			Username: username,
			Password: hashedPassword,
			Tokens:   user.Tokens,
		}
		storage.PutUser(username, newUser)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return nil
	}
}

func getDeleteUserHandler(storage storage.CertifierStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param("username")
		err := storage.DeleteUser(username)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return nil
	}
}

func getLoginHandler(storage storage.CertifierStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := make(map[string]interface{})
		err := json.NewDecoder(c.Request().Body).Decode(&data)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error decoding request for User: %v", err))
		}

		username := fmt.Sprintf("%v", data["username"])
		password := []byte(fmt.Sprintf("%v", data["password"]))

		user, err := storage.GetUser(username)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		hashedPassword := user.GetPassword()
		err = bcrypt.CompareHashAndPassword(hashedPassword, password)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "wrong password")
		}

		tokens := user.GetTokens().GetToken()
		return c.JSON(http.StatusOK, tokens)
	}
}
