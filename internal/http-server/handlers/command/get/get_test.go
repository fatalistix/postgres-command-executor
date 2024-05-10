package get_test

import (
	"encoding/json"
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/database"
	"github.com/fatalistix/postgres-command-executor/internal/domain/models"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/get"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/get/mocks"
	"github.com/fatalistix/postgres-command-executor/internal/lib/log/slog"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetHandler(t *testing.T) {
	cases := []struct {
		name            string
		commandID       string
		responseError   string
		mockError       error
		mockCommand     models.Command
		mockCalledTimes int
		responseCode    int
	}{
		{
			name:            "invalid id as string",
			commandID:       "invalid",
			responseError:   "invalid id",
			mockError:       nil,
			mockCommand:     models.Command{},
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "invalid id as too big number",
			commandID:       "1000000000000000000000000000000000000000000000000000000000000000000",
			responseError:   "invalid id",
			mockError:       nil,
			mockCommand:     models.Command{},
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "empty id",
			commandID:       "",
			responseError:   "invalid id",
			mockError:       nil,
			mockCommand:     models.Command{},
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "not found",
			commandID:       "1",
			responseError:   "command not found",
			mockError:       database.ErrCommandNotFound,
			mockCommand:     models.Command{},
			mockCalledTimes: 1,
			responseCode:    http.StatusNotFound,
		},
		{
			name:            "internal error",
			commandID:       "1",
			responseError:   "error getting command",
			mockError:       database.ErrInternal,
			mockCommand:     models.Command{},
			mockCalledTimes: 1,
			responseCode:    http.StatusInternalServerError,
		},
		{
			name:          "success",
			commandID:     "1",
			responseError: "",
			mockError:     nil,
			mockCommand: models.Command{
				ID:      1,
				Command: "command",
			},
			mockCalledTimes: 1,
			responseCode:    http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			commandProviderMock := mocks.NewCommandProvider(t)

			if parseResult, err := strconv.ParseInt(tc.commandID, 10, 64); err == nil {
				commandProviderMock.On("Command", parseResult).Return(tc.mockCommand, tc.mockError)
			}

			handler := get.MakeGetHandlerFunc(slog.Discard(), commandProviderMock)

			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/command/%s", tc.commandID), nil)
			require.NoError(t, err)

			request.SetPathValue("id", tc.commandID)

			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			require.Equal(t, tc.responseCode, responseRecorder.Code)

			body := responseRecorder.Body.String()

			if tc.responseError != "" {
				require.Contains(t, body, tc.responseError)
			} else {
				require.NotEmpty(t, body)

				var response get.Response
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)

				require.Equal(t, tc.mockCommand.ID, response.ID)
				require.Equal(t, tc.mockCommand.Command, response.Command)
			}
		})
	}
}
