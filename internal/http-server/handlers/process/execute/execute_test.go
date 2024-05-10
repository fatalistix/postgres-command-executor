package execute_test

import (
	"encoding/json"
	"github.com/fatalistix/postgres-command-executor/internal/database"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/process/execute"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/process/execute/mocks"
	"github.com/fatalistix/postgres-command-executor/internal/lib/log/slog"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExecuteHandler(t *testing.T) {
	cases := []struct {
		name            string
		requestJSON     string
		httpHeader      string
		responseError   string
		mockError       error
		mockCommandID   int64
		mockProcessID   uuid.UUID
		mockCalledTimes int
		responseCode    int
	}{
		{
			name:            "'Content-Type' header is 'application/xml'",
			requestJSON:     "",
			httpHeader:      "Application/xml",
			responseError:   "no 'application/json' header found",
			mockError:       nil,
			mockCommandID:   0,
			mockProcessID:   uuid.Nil,
			mockCalledTimes: 0,
			responseCode:    http.StatusUnsupportedMediaType,
		},
		{
			name:            "no 'Content-Type' header",
			requestJSON:     "",
			httpHeader:      "",
			responseError:   "no 'application/json' header found",
			mockError:       nil,
			mockCommandID:   0,
			mockProcessID:   uuid.Nil,
			mockCalledTimes: 0,
			responseCode:    http.StatusUnsupportedMediaType,
		},
		{
			name:            "empty body",
			requestJSON:     "",
			httpHeader:      "Application/json",
			responseError:   "unable to decode request's body",
			mockError:       nil,
			mockCommandID:   0,
			mockProcessID:   uuid.Nil,
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "not a json",
			requestJSON:     `I am not a json <XMl>`,
			httpHeader:      "Application/json",
			responseError:   "unable to decode request's body",
			mockError:       nil,
			mockCommandID:   0,
			mockProcessID:   uuid.Nil,
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "no command id json field",
			requestJSON:     `{"output": "some output"}`,
			httpHeader:      "Application/json",
			responseError:   "unable to decode request's body",
			mockError:       nil,
			mockCommandID:   0,
			mockProcessID:   uuid.Nil,
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "unknown json field",
			requestJSON:     `{"command_id": 1, "output": "some output"}`,
			httpHeader:      "Application/json",
			responseError:   "unable to decode request's body",
			mockError:       nil,
			mockCommandID:   0,
			mockProcessID:   uuid.Nil,
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "invalid command id",
			requestJSON:     `{"command_id": 1}`,
			httpHeader:      "Application/json",
			responseError:   "command not found",
			mockError:       database.ErrCommandNotFound,
			mockCommandID:   1,
			mockProcessID:   uuid.Nil,
			mockCalledTimes: 1,
			responseCode:    http.StatusNotFound,
		},
		{
			name:            "internal error",
			requestJSON:     `{"command_id": 1}`,
			httpHeader:      "Application/json",
			responseError:   "error starting command execution",
			mockError:       database.ErrInternal,
			mockCommandID:   1,
			mockProcessID:   uuid.Nil,
			mockCalledTimes: 1,
			responseCode:    http.StatusInternalServerError,
		},
		{
			name:            "success",
			requestJSON:     `{"command_id": 1}`,
			httpHeader:      "Application/json",
			responseError:   "",
			mockError:       nil,
			mockCommandID:   1,
			mockProcessID:   uuid.New(),
			mockCalledTimes: 1,
			responseCode:    http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			commandExecutionStarterMock := mocks.NewCommandExecutionStarter(t)

			if tc.mockCalledTimes != 0 {
				commandExecutionStarterMock.
					On("StartCommandExecution", tc.mockCommandID).
					Return(tc.mockProcessID, tc.mockError).
					Times(tc.mockCalledTimes)
			}

			handler := execute.MakeExecuteHandlerFunc(slog.Discard(), commandExecutionStarterMock)

			request, err := http.NewRequest(http.MethodPost, "/processes", strings.NewReader(tc.requestJSON))
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

				var response execute.Response
				err = json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)

				require.Equal(t, tc.mockProcessID, response.ProcessID)
			}
		})
	}
}
