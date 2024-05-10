package get_test

import (
	"encoding/json"
	"github.com/fatalistix/postgres-command-executor/internal/database"
	"github.com/fatalistix/postgres-command-executor/internal/domain/models"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/process/get"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/process/get/mocks"
	"github.com/fatalistix/postgres-command-executor/internal/lib/log/slog"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetHandler(t *testing.T) {
	cases := []struct {
		name            string
		processID       string
		responseError   string
		mockError       error
		mockProcess     models.Process
		mockCalledTimes int
		responseCode    int
	}{
		{
			name:            "invalid id as non-uuid string",
			processID:       "some-non-uuid-string",
			responseError:   "invalid id",
			mockError:       nil,
			mockProcess:     models.Process{},
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "invalid id as too big number",
			processID:       "1000000000000000000000000000000000000000000000000000000000000000000",
			responseError:   "invalid id",
			mockError:       nil,
			mockProcess:     models.Process{},
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "invalid id as small number",
			processID:       "0",
			responseError:   "invalid id",
			mockError:       nil,
			mockProcess:     models.Process{},
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "empty id",
			processID:       "",
			responseError:   "invalid id",
			mockError:       nil,
			mockProcess:     models.Process{},
			mockCalledTimes: 0,
			responseCode:    http.StatusBadRequest,
		},
		{
			name:            "not found",
			processID:       uuid.New().String(),
			responseError:   "process not found",
			mockError:       database.ErrProcessNotFound,
			mockProcess:     models.Process{},
			mockCalledTimes: 1,
			responseCode:    http.StatusNotFound,
		},
		{
			name:            "internal error",
			processID:       uuid.New().String(),
			responseError:   "error getting process",
			mockError:       database.ErrInternal,
			mockProcess:     models.Process{},
			mockCalledTimes: 1,
			responseCode:    http.StatusInternalServerError,
		},
		{
			name:          "success",
			processID:     uuid.New().String(),
			responseError: "",
			mockError:     nil,
			mockProcess: models.Process{
				ID:       uuid.New(), // won't be synchronized with request's ID
				Output:   "some output",
				Error:    "some error",
				Status:   models.StatusFinished,
				ExitCode: 0,
			},
			mockCalledTimes: 1,
			responseCode:    http.StatusOK,
		},
		{
			name:          "success with raw uuid (without '-' symbol)",
			processID:     strings.ReplaceAll(uuid.New().String(), "-", ""),
			responseError: "",
			mockError:     nil,
			mockProcess: models.Process{
				ID:       uuid.New(), // won't be synchronized with request's ID
				Output:   "some output",
				Error:    "some error",
				Status:   models.StatusFinished,
				ExitCode: 0,
			},
			mockCalledTimes: 1,
			responseCode:    http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			processProviderMock := mocks.NewProcessProvider(t)

			if parseResult, err := uuid.Parse(tc.processID); err == nil {
				processProviderMock.
					On("Process", parseResult).
					Return(&tc.mockProcess, tc.mockError).
					Times(tc.mockCalledTimes)
			}

			handler := get.MakeGetHandlerFunc(slog.Discard(), processProviderMock)

			request, err := http.NewRequest(http.MethodGet, "/processes/"+tc.processID, nil)
			require.NoError(t, err)

			request.SetPathValue("id", tc.processID)

			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			require.Equal(t, tc.responseCode, responseRecorder.Code)

			body := responseRecorder.Body.String()

			if tc.responseError != "" {
				require.Contains(t, responseRecorder.Body.String(), tc.responseError)
			} else {
				require.NotEmpty(t, body)

				var response get.Response
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)

				require.Equal(t, tc.mockProcess.ID, response.ID)
				require.Equal(t, tc.mockProcess.Output, response.Output)
				require.Equal(t, tc.mockProcess.Error, response.Error)
				require.Equal(t, tc.mockProcess.Status, models.ProcessStatus(response.Status))
				require.Equal(t, tc.mockProcess.ExitCode, response.ExitCode)
			}
		})
	}
}
