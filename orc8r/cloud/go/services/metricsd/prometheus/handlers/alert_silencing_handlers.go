package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/lib/go/metrics"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/prometheus/alertmanager/api/v2/client/silence"
	"github.com/prometheus/alertmanager/api/v2/models"
	"github.com/prometheus/alertmanager/pkg/parse"
	"github.com/prometheus/prometheus/pkg/labels"
)

const (
	silenceIDParam     = "silence_id"
	filterParam        = "filter"
	activeStatusParam  = "active"
	pendingStatusParam = "pending"
	expiredStatusParam = "expired"

	getSilencesPath    = "/silences"
	postSilencesPath   = "/silences"
	deleteSilencesPath = "/silence"
)

func GetPostSilencerHandler(alertmanagerURL string, client HttpClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		silencerURL := alertmanagerURL + getSilencesPath
		return postSilencer(networkID, silencerURL, c, client)
	}
}

func GetGetSilencersHandler(alertmanagerURL string, client HttpClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		silencerURL := alertmanagerURL + postSilencesPath
		return getSilencers(networkID, silencerURL, c, client)
	}
}

func GetDeleteSilencerHandler(alertmanagerURL string, client HttpClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		silencerURL := alertmanagerURL + deleteSilencesPath
		silenceID := c.QueryParam(silenceIDParam)
		return deleteSilencer(silenceID, silencerURL, c, client)
	}
}

func postSilencer(networkID, silencerURL string, c echo.Context, client HttpClient) error {
	silencer, err := buildSilencerFromContext(c)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	isRegex := false
	labelName := metrics.NetworkLabelName
	networkMatcher := models.Matcher{
		IsRegex: &isRegex,
		Name:    &labelName,
		Value:   &networkID,
	}
	silencer.Matchers = append(silencer.Matchers, &networkMatcher)

	newSilencerBytes, err := json.Marshal(silencer)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "make silencer"), http.StatusInternalServerError)
	}

	resp, err := client.Post(silencerURL, "application/json", bytes.NewBuffer(newSilencerBytes))
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return obsidian.HttpError(fmt.Errorf("error posting silencer: %s", respBody), resp.StatusCode)
	}
	var silenceResponse silence.PostSilencesOKBody
	err = json.NewDecoder(resp.Body).Decode(&silenceResponse)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("error decoding alertmanager response: %v", err), http.StatusInternalServerError)
	}
	return c.String(http.StatusOK, silenceResponse.SilenceID)
}

func getSilencers(networkID, silencerURL string, c echo.Context, client HttpClient) error {
	filters, err := parse.Matchers(c.QueryParam(filterParam))
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	filters = append(filters, &labels.Matcher{Type: labels.MatchEqual, Name: metrics.NetworkLabelName, Value: networkID})

	filteredURL, err := url.Parse(silencerURL)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	q := filteredURL.Query()
	q.Set(filterParam, filtersToString(filters))
	filteredURL.RawQuery = q.Encode()

	resp, err := client.Get(filteredURL.String())
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("error getting silences: %v", err), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return obsidian.HttpError(fmt.Errorf("error getting silences: %s", respBody), resp.StatusCode)
	}

	var silencers []models.GettableSilence
	err = json.NewDecoder(resp.Body).Decode(&silencers)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("error decoding server response: %v", err), http.StatusInternalServerError)
	}

	// Alertmanager API doesn't implement filtering on status so we have to do
	// it here
	returnSilencers := make([]models.GettableSilence, 0)
	getActiveSilences := c.QueryParam(activeStatusParam) != "false"
	getPendingSilences := c.QueryParam(pendingStatusParam) != "false"
	getExpiredSilences := c.QueryParam(expiredStatusParam) != "false"
	for _, sil := range silencers {
		switch *sil.Status.State {
		case models.SilenceStatusStateActive:
			if getActiveSilences {
				returnSilencers = append(returnSilencers, sil)
			}
		case models.SilenceStatusStatePending:
			if getPendingSilences {
				returnSilencers = append(returnSilencers, sil)
			}
		case models.SilenceStatusStateExpired:
			if getExpiredSilences {
				returnSilencers = append(returnSilencers, sil)
			}
		}
	}
	return c.JSON(http.StatusOK, returnSilencers)
}

func deleteSilencer(silenceID, silencerURL string, c echo.Context, client HttpClient) error {
	req, err := http.NewRequest(http.MethodDelete, silencerURL+"/"+silenceID, nil)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	resp, err := client.Do(req)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("error deleting silence: %v", err), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return obsidian.HttpError(fmt.Errorf("error deleting silence: %s", respBody), resp.StatusCode)
	}
	return c.NoContent(http.StatusOK)
}

func buildSilencerFromContext(c echo.Context) (models.Silence, error) {
	var silencer models.Silence
	err := json.NewDecoder(c.Request().Body).Decode(&silencer)
	if err != nil {
		return models.Silence{}, err
	}
	return silencer, nil
}

func filtersToString(filters []*labels.Matcher) string {
	ret := strings.Builder{}
	ret.WriteString("{")
	for i, filter := range filters {
		if i > 0 {
			ret.WriteString(",")
		}
		ret.WriteString(filterToString(*filter))
	}
	ret.WriteString("}")
	return ret.String()
}

func filterToString(filter labels.Matcher) string {
	ret := strings.Builder{}
	ret.WriteString(filter.Name)
	switch filter.Type {
	case labels.MatchEqual:
		ret.WriteString("=")
	case labels.MatchNotEqual:
		ret.WriteString("!=")
	case labels.MatchRegexp:
		ret.WriteString("=~")
	case labels.MatchNotRegexp:
		ret.WriteString("!~")
	}
	ret.WriteString(fmt.Sprintf(`"%s"`, filter.Value))
	return ret.String()
}
