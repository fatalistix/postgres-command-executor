package list_test

import (
	"encoding/json"
	"github.com/fatalistix/postgres-command-executor/internal/database"
	"github.com/fatalistix/postgres-command-executor/internal/domain/models"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/list"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/list/mocks"
	"github.com/fatalistix/postgres-command-executor/internal/lib/log/slog"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListHandler(t *testing.T) {
	cases := []struct {
		name            string
		responseError   string
		mockError       error
		mockCommands    []models.Command
		mockCalledTimes int
		responseCode    int
	}{
		{
			name:            "internal error",
			responseError:   "error getting commands",
			mockError:       database.ErrInternal,
			mockCommands:    nil,
			mockCalledTimes: 1,
			responseCode:    http.StatusInternalServerError,
		},
		{
			name:            "success",
			responseError:   "",
			mockError:       nil,
			mockCommands:    []models.Command{{ID: 1, Command: "command1"}, {ID: 2, Command: "command2"}},
			mockCalledTimes: 1,
			responseCode:    http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			commandProviderMock := mocks.NewCommandProvider(t)

			commandProviderMock.
				On("Commands").
				Return(tc.mockCommands, tc.mockError).
				Times(tc.mockCalledTimes)

			handler := list.MakeListHandlerFunc(slog.Discard(), commandProviderMock)

			request, err := http.NewRequest(http.MethodGet, "/commands", nil)
			require.NoError(t, err)

			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			require.Equal(t, tc.responseCode, responseRecorder.Code)

			body := responseRecorder.Body.String()

			if tc.responseError != "" {
				require.Contains(t, responseRecorder.Body.String(), tc.responseError)
			} else {
				require.NotEmpty(t, body)

				var response list.Response

				err = json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)

				require.Equal(t, len(tc.mockCommands), len(response.Commands))

				for i := range len(tc.mockCommands) {
					require.Equal(t, tc.mockCommands[i].ID, response.Commands[i].ID)
					require.Equal(t, tc.mockCommands[i].Command, response.Commands[i].Command)
				}
			}
		})
	}
}
