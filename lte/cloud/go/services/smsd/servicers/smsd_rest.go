package servicers

import (
	"net/http"

	lteHandlers "magma/lte/cloud/go/services/lte/obsidian/handlers"
	"magma/lte/cloud/go/services/smsd/obsidian/models"
	"magma/lte/cloud/go/services/smsd/storage"
	"magma/orc8r/cloud/go/obsidian"

	"github.com/labstack/echo"
	"github.com/thoas/go-funk"
)

const (
	SmsRootPath   = lteHandlers.ManageNetworkPath + obsidian.UrlSep + "sms"
	SmsManagePath = SmsRootPath + obsidian.UrlSep + ":sms_pk"
)

func NewRESTServicer(store storage.SMSStorage) *SMSDRestServicer {
	return &SMSDRestServicer{store: store}
}

type SMSDRestServicer struct {
	store storage.SMSStorage
}

func (s *SMSDRestServicer) GetHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		{Path: SmsRootPath, Methods: obsidian.GET, HandlerFunc: s.listMessages},
		{Path: SmsRootPath, Methods: obsidian.POST, HandlerFunc: s.createMessage},
		{Path: SmsManagePath, Methods: obsidian.GET, HandlerFunc: s.getMessage},
		{Path: SmsManagePath, Methods: obsidian.DELETE, HandlerFunc: s.deleteMessage},
	}
}

func (s *SMSDRestServicer) listMessages(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	messages, err := s.store.GetSMSs(networkID, nil, nil, false, nil, nil)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	out := make([]*models.SmsMessage, 0, len(messages))
	for _, msg := range messages {
		out = append(out, (&models.SmsMessage{}).FromProto(msg))
	}
	return c.JSON(http.StatusOK, out)
}

func (s *SMSDRestServicer) getMessage(c echo.Context) error {
	networkID, pk, nerr := getNetworkAndSMSID(c)
	if nerr != nil {
		return nerr
	}

	msgs, err := s.store.GetSMSs(networkID, []string{pk}, nil, false, nil, nil)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if funk.IsEmpty(msgs) {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, (&models.SmsMessage{}).FromProto(msgs[0]))
}

func (s *SMSDRestServicer) createMessage(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload := &models.MutableSmsMessage{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	pk, err := s.store.CreateSMS(networkID, payload.ToProto())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, pk)
}

func (s *SMSDRestServicer) deleteMessage(c echo.Context) error {
	networkID, pk, nerr := getNetworkAndSMSID(c)
	if nerr != nil {
		return nerr
	}

	err := s.store.DeleteSMSs(networkID, []string{pk})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)

}

func getNetworkAndSMSID(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", "sms_pk")
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}
