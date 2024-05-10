package delete_test

import (
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/database"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/process/delete"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/process/delete/mocks"
	"github.com/fatalistix/postgres-command-executor/internal/lib/log/slog"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name            string
		processID       string
		responseError   string
		mockError       error
		mockCalledTimes int
		responseCode    int
	}{
		{
			name:            "invalid id as non-uuid string",
			processID:       "some-non-uuid-string",
			responseError:   "invalid id",
			mockError:       nil,
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "invalid id as too big number",
			processID:       "1000000000000000000000000000000000000000000000000000000000000000000",
			responseError:   "invalid id",
			mockError:       nil,
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "invalid id as small number",
			processID:       "0",
			responseError:   "invalid id",
			mockError:       nil,
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "empty id",
			processID:       "",
			responseError:   "invalid id",
			mockError:       nil,
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "not found",
			processID:       uuid.New().String(),
			responseError:   "process not found",
			mockError:       database.ErrProcessNotFound,
			mockCalledTimes: 1,
			responseCode:    http.StatusNotFound,
		},
		{
			name:            "internal error",
			processID:       uuid.New().String(),
			responseError:   "error deleting process",
			mockError:       database.ErrInternal,
			mockCalledTimes: 1,
			responseCode:    http.StatusInternalServerError,
		},
		{
			name:            "success",
			processID:       uuid.New().String(),
			responseError:   "",
			mockError:       nil,
			mockCalledTimes: 1,
			responseCode:    http.StatusOK,
		},
		{
			name:            "success with raw uuid (without '-' symbol)",
			processID:       strings.ReplaceAll(uuid.New().String(), "-", ""),
			responseError:   "",
			mockError:       nil,
			mockCalledTimes: 1,
			responseCode:    http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			processDeleterMock := mocks.NewProcessDeleter(t)

			if parseResult, err := uuid.Parse(tc.processID); err == nil {
				processDeleterMock.
					On("DeleteProcess", parseResult).
					Return(tc.mockError).
					Times(tc.mockCalledTimes)
			}

			handler := delete.MakeDeleteHandlerFunc(slog.Discard(), processDeleterMock)

			request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/process/%s", tc.processID), nil)
			require.NoError(t, err)

			request.SetPathValue("id", tc.processID)

			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			require.Equal(t, tc.responseCode, responseRecorder.Code)

			body := responseRecorder.Body.String()

			require.Contains(t, body, tc.responseError)
		})
	}
}
