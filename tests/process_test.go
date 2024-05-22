package tests

import (
	"encoding/json"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/process/execute"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/process/get"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestExecutionStartAndResult(t *testing.T) {
	commandCreateResp := mustSaveCommand(t, `echo \"1\"; sleep 3; echo \"2\"`)

	processExecuteResp, err := http.Post(
		url+"/processes",
		"application/json",
		strings.NewReader(
			`{ "command_id": `+strconv.FormatInt(commandCreateResp.ID, 10)+` }`,
		),
	)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, processExecuteResp.StatusCode)

	body, err := io.ReadAll(processExecuteResp.Body)
	require.NoError(t, err)

	var executeResp execute.Response

	err = json.Unmarshal(body, &executeResp)
	require.NoError(t, err)

	t.Logf("process_id: %s, waiting for result", executeResp.ProcessID)
	time.Sleep(time.Second)

	processGetResp := makeGetProcess(t, executeResp.ProcessID.String())

	require.Equal(t, "1\n", processGetResp.Output)
	require.Equal(t, "", processGetResp.Error)
	require.Equal(t, -1, processGetResp.ExitCode)
	require.Equal(t, "executing", processGetResp.Status)

	time.Sleep(time.Second * 3)

	processGetResp = makeGetProcess(t, executeResp.ProcessID.String())

	require.Equal(t, "1\n2\n", processGetResp.Output)
	require.Equal(t, "", processGetResp.Error)
	require.Equal(t, 0, processGetResp.ExitCode)
	require.Equal(t, "finished", processGetResp.Status)
}

func makeGetProcess(t *testing.T, processID string) get.Response {
	processGetResp, err := http.Get(url + "/process/" + processID)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, processGetResp.StatusCode)

	body, err := io.ReadAll(processGetResp.Body)
	require.NoError(t, err)

	var getRespBody get.Response

	err = json.Unmarshal(body, &getRespBody)
	require.NoError(t, err)

	return getRespBody
}
