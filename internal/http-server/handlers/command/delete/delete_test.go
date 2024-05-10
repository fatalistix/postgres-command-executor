package delete_test

import (
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/database"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/delete"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/delete/mocks"
	"github.com/fatalistix/postgres-command-executor/internal/lib/log/slog"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name            string
		commandID       string
		responseError   string
		mockError       error
		mockCalledTimes int
		responseCode    int
	}{
		{
			name:            "invalid id as string",
			commandID:       "invalid",
			responseError:   "invalid id",
			mockError:       nil,
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "invalid id as too big number",
			commandID:       "1000000000000000000000000000000000000000000000000000000000000000000",
			responseError:   "invalid id",
			mockError:       nil,
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "empty id",
			commandID:       "",
			responseError:   "invalid id",
			mockError:       nil,
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "not found",
			commandID:       "1",
			responseError:   "command not found",
			mockError:       database.ErrCommandNotFound,
			mockCalledTimes: 1,
			responseCode:    http.StatusNotFound,
		},
		{
			name:            "internal error",
			commandID:       "1",
			responseError:   "error deleting command",
			mockError:       database.ErrInternal,
			mockCalledTimes: 1,
			responseCode:    http.StatusInternalServerError,
		},
		{
			name:            "success",
			commandID:       "1",
			responseError:   "",
			mockError:       nil,
			mockCalledTimes: 1,
			responseCode:    http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			commandDeleterMock := mocks.NewCommandDeleter(t)

			if parseResult, err := strconv.ParseInt(tc.commandID, 10, 64); err == nil {
				commandDeleterMock.
					On("DeleteCommand", parseResult).
					Return(tc.mockError).
					Times(tc.mockCalledTimes)
			}

			handler := delete.MakeDeleteHandlerFunc(slog.Discard(), commandDeleterMock)

			request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/command/%s", tc.commandID), nil)
			require.NoError(t, err)

			request.SetPathValue("id", tc.commandID)

			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			require.Equal(t, tc.responseCode, responseRecorder.Code)

			body := responseRecorder.Body.String()

			require.Contains(t, body, tc.responseError)
		})
	}
}
