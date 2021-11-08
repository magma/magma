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

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/certifier/storage"
)

const (
	HTTPBasicAuth       = "http_basic_auth"
	ListHTTPBasicAuth   = obsidian.V1Root + HTTPBasicAuth
	ManageHTTPBasicAuth = ListHTTPBasicAuth + obsidian.UrlSep + ":username"
)

func GetHandlers(storage storage.CertifierStorage) []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: ListHTTPBasicAuth, Methods: obsidian.GET, HandlerFunc: getListHTTPBasicAuthHandler(storage)},
		{Path: ListHTTPBasicAuth, Methods: obsidian.POST, HandlerFunc: getCreateHTTPBasicAuthHandler(storage)},
		{Path: ManageHTTPBasicAuth, Methods: obsidian.PUT, HandlerFunc: getUpdateHTTPBasicAuthHandler(storage)},
		{Path: ManageHTTPBasicAuth, Methods: obsidian.DELETE, HandlerFunc: getDeleteHTTPBasicAuthHandler(storage)},
	}
	return ret
}

func getListHTTPBasicAuthHandler(storage storage.CertifierStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := storage.ListHTTPBasicAuth()
		glog.Errorf("christine wtf are these users %s", users)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, users)
	}
}

func getCreateHTTPBasicAuthHandler(storage storage.CertifierStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := make(map[string]interface{})
		err := json.NewDecoder(c.Request().Body).Decode(&data)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error decoding request for HTTP basic auth: %v", err))
		}
		username := fmt.Sprintf("%v", data["username"])
		strPassword := fmt.Sprintf("%v", data["password"])
		password := []byte(strPassword)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error hashing password: %v", err))
		}

		glog.Errorf("christine username %s", username)
		glog.Errorf("christine WHY IS THE PASSWORD ALWAYS THE SAME?? %s		 %s", strPassword, hashedPassword)

		err = storage.UpdateHTTPBasicAuth(username, hashedPassword)
		if err != nil {
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
		strPassword := fmt.Sprintf("%v", data["password"])
		password := []byte(strPassword)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		err = storage.UpdateHTTPBasicAuth(username, hashedPassword)
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
