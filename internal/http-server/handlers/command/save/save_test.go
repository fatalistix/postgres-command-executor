package save_test

import (
	"encoding/json"
	"github.com/fatalistix/postgres-command-executor/internal/database"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/save"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/save/mocks"
	"github.com/fatalistix/postgres-command-executor/internal/lib/log/slog"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name            string
		requestJSON     string
		httpHeader      string
		responseError   string
		mockError       error
		mockID          int64
		mockCommand     string
		mockCalledTimes int
		responseCode    int
	}{
		{
			name:            "'Content-Type' header is 'application/xml'",
			requestJSON:     "",
			httpHeader:      "Application/xml",
			responseError:   "no 'application/json' header found",
			mockError:       nil,
			mockID:          0,
			mockCommand:     "",
			mockCalledTimes: 0,
			responseCode:    http.StatusUnsupportedMediaType,
		},
		{
			name:            "no 'Content-Type' header",
			requestJSON:     "",
			httpHeader:      "",
			responseError:   "no 'application/json' header found",
			mockError:       nil,
			mockID:          0,
			mockCommand:     "",
			mockCalledTimes: 0,
			responseCode:    http.StatusUnsupportedMediaType,
		},
		{
			name:            "empty body",
			requestJSON:     "",
			httpHeader:      "Application/json",
			responseError:   "unable to decode request's body",
			mockError:       nil,
			mockID:          0,
			mockCommand:     "",
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "empty json",
			requestJSON:     `{}`,
			httpHeader:      "Application/json",
			responseError:   "empty command",
			mockError:       nil,
			mockID:          0,
			mockCommand:     "",
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "not a json",
			requestJSON:     `some that is not a json <xmL>`,
			httpHeader:      "Application/json",
			responseError:   "unable to decode request's body",
			mockError:       nil,
			mockID:          0,
			mockCommand:     "",
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "empty command",
			requestJSON:     `{"command": ""}`,
			httpHeader:      "Application/json",
			responseError:   "empty command",
			mockError:       nil,
			mockID:          0,
			mockCommand:     "",
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "no command json field",
			requestJSON:     `{"invalid_field": "invalid"}`,
			httpHeader:      "Application/json",
			responseError:   "unable to decode request's body",
			mockError:       nil,
			mockID:          0,
			mockCommand:     "",
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "unknown json field",
			requestJSON:     `{"command": "some command", "invalid_field": "invalid"}`,
			httpHeader:      "Application/json",
			responseError:   "unable to decode request's body",
			mockError:       nil,
			mockID:          0,
			mockCommand:     "",
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "already exists",
			requestJSON:     `{"command": "some command"}`,
			httpHeader:      "Application/json",
			responseError:   "command already exists",
			mockError:       database.ErrCommandExists,
			mockID:          0,
			mockCommand:     "some command",
			mockCalledTimes: 1,
			responseCode:    http.StatusConflict,
		},
		{
			name:            "success",
			requestJSON:     `{"command": "some command"}`,
			httpHeader:      "Application/json",
			responseError:   "",
			mockError:       nil,
			mockID:          1,
			mockCommand:     "some command",
			mockCalledTimes: 1,
			responseCode:    http.StatusCreated,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			commandSaverMock := mocks.NewCommandSaver(t)

			if tc.mockCalledTimes != 0 {
				commandSaverMock.
					On("SaveCommand", tc.mockCommand).
					Return(tc.mockID, tc.mockError).
					Times(tc.mockCalledTimes)
			}

			handler := save.MakeSaveHandlerFunc(slog.Discard(), commandSaverMock)

			request, err := http.NewRequest(http.MethodPost, "/commands", strings.NewReader(tc.requestJSON))
			require.NoError(t, err)

			request.Header.Set("Content-Type", tc.httpHeader)

			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			require.Equal(t, tc.responseCode, responseRecorder.Code)

			body := responseRecorder.Body.String()

			if tc.responseError != "" {
				require.Contains(t, body, tc.responseError)
			} else {
				require.NotEmpty(t, body)

				var response save.Response
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)

				require.Equal(t, tc.mockID, response.ID)
			}
		})
	}
}
