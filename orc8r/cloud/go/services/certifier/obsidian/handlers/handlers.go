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
	certProtos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/storage"
)

const (
	UserParam           = ":username"
	HTTPBasicAuth       = "http_basic_auth"
	ListHTTPBasicAuth   = obsidian.V1Root + HTTPBasicAuth
	ManageHTTPBasicAuth = ListHTTPBasicAuth + obsidian.UrlSep + UserParam
	Login               = ListHTTPBasicAuth + obsidian.UrlSep + "login"
)

func GetHandlers(storage storage.CertifierStorage) []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: ListHTTPBasicAuth, Methods: obsidian.GET, HandlerFunc: getListHTTPBasicAuthHandler(storage)},
		{Path: ListHTTPBasicAuth, Methods: obsidian.POST, HandlerFunc: getCreateHTTPBasicAuthHandler(storage)},
		{Path: ManageHTTPBasicAuth, Methods: obsidian.PUT, HandlerFunc: getUpdateHTTPBasicAuthHandler(storage)},
		{Path: ManageHTTPBasicAuth, Methods: obsidian.DELETE, HandlerFunc: getDeleteHTTPBasicAuthHandler(storage)},
		{Path: Login, Methods: obsidian.POST, HandlerFunc: getLoginHandler(storage)},
	}
	return ret
}

func getListHTTPBasicAuthHandler(storage storage.CertifierStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := storage.ListHTTPBasicAuth()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, users)
	}
}

// TODO(christinewang5): should not be able to create users that already exist
func getCreateHTTPBasicAuthHandler(storage storage.CertifierStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := make(map[string]interface{})
		err := json.NewDecoder(c.Request().Body).Decode(&data)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error decoding request for HTTP basic auth: %v", err))
		}

		username := fmt.Sprintf("%v", data["username"])
		password := []byte(fmt.Sprintf("%v", data["password"]))
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error hashing password: %v", err))
		}

		token, err := certifier.GenerateToken(certifier.Personal)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error generating personal access token for operator: %v", err))
		}

		operator := &certProtos.Operator{
			Username: username,
			Password: hashedPassword,
			Tokens:   &certProtos.Operator_TokenList{Token: []string{token}},
		}
		if err = storage.PutHTTPBasicAuth(username, operator); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		// TODO(christinewang5): remove this once bootstrapping is finished...
		policy := &certProtos.Policy{
			Token:     token,
			Effect:    certProtos.Effect_ALLOW,
			Action:    certProtos.Action_WRITE,
			Resources: &certProtos.Policy_ResourceList{Resource: []string{"*"}},
		}
		if err = storage.PutPolicy(token, policy); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return nil
	}
}

func getUpdateHTTPBasicAuthHandler(storage storage.CertifierStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param("username")
		data := make(map[string]interface{})
		err := json.NewDecoder(c.Request().Body).Decode(&data)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error decoding request for HTTP basic auth: %v", err))
		}

		newPassword := []byte(fmt.Sprintf("%v", data["password"]))
		hashedPassword, err := bcrypt.GenerateFromPassword(newPassword, bcrypt.DefaultCost)

		// get old operator
		operator, err := storage.GetHTTPBasicAuth(username)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting existing user: %v", err))
		}

		// update new operator blob
		newOperator := &certProtos.Operator{
			Username: username,
			Password: hashedPassword,
			Tokens:   operator.Tokens,
		}
		storage.PutHTTPBasicAuth(username, newOperator)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return nil
	}
}

func getDeleteHTTPBasicAuthHandler(storage storage.CertifierStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param("username")
		err := storage.DeleteHTTPBasicAuth(username)
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
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error decoding request for HTTP basic auth: %v", err))
		}

		username := fmt.Sprintf("%v", data["username"])
		password := []byte(fmt.Sprintf("%v", data["password"]))

		operator, err := storage.GetHTTPBasicAuth(username)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		// check password hash
		hashedPassword := operator.GetPassword()
		err = bcrypt.CompareHashAndPassword(hashedPassword, password)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "wrong password")
		}

		// return tokens for access if correct password
		tokens := operator.GetTokens().GetToken()
		return c.JSON(http.StatusOK, tokens)
	}
}
