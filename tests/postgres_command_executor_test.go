package tests

import (
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/save"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"net/url"
	"testing"
)

const (
	host = "localhost:8089"
)

func TestPostgresCommandExecutor_SaveAndDeleteCommand(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	var id int64
	e.POST("/commands").
		WithJSON(save.Request{
			Command: "ls; ls",
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().
		Object().
		ContainsKey("id").
		Value("id").
		Decode(&id)

	e.DELETE(fmt.Sprintf("/command/%d", id)).
		Expect().
		Status(http.StatusOK)
}
