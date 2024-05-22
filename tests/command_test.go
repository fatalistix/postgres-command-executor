package tests

import (
	"encoding/json"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/save"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

const url = "http://localhost:8089"

// TestCommandCreate assumes, that there is no `ls -la` command in database
func TestCommandCreate(t *testing.T) {
	saveResp := mustCreateCommand(t, "ls -la")

	require.True(t, saveResp.ID > 0)
}

func mustCreateCommand(t *testing.T, command string) save.Response {
	resp, err := http.Post(
		url+"/commands",
		"application/json",
		strings.NewReader(
			`{ "command": "`+command+`"}`,
		),
	)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var saveResp save.Response

	err = json.Unmarshal(body, &saveResp)
	require.NoError(t, err)

	return saveResp
}

// TestCommandDuplicate assumes, that there is no `nvim .` command in database. It creates `nvim .` command and then
// tries to create it again. On first request it waits for 201 and on second request it waits for 409
func TestCommandDuplicate(t *testing.T) {
	const command = "nvim ."

	saveResp := mustCreateCommand(t, command)

	require.True(t, saveResp.ID > 0)

	resp, err := http.Post(
		url+"/commands",
		"application/json",
		strings.NewReader(
			`{ "command": "`+command+`"}`,
		),
	)
	require.NoError(t, err)

	require.Equal(t, http.StatusConflict, resp.StatusCode)
}

// TestCommandGet creates command and then tries to get it
func TestCommandGet(t *testing.T) {
	const command = "touch some_file.txt"
	saveResp := mustCreateCommand(t, command)

	resp, err := http.Get(url + "/command/" + strconv.FormatInt(saveResp.ID, 10))
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	sb := string(body)

	require.Contains(t, sb, command)
}

// TestCommandList creates a list of commands and then tries to get a list. It checks that all created commands are in
// the list
func TestCommandList(t *testing.T) {
	var commands = []string{
		"echo 1",
		"echo 2",
		"echo 3",
		"touch 1.c",
		"touch 2.c",
		"touch 3.c",
	}

	for _, command := range commands {
		mustCreateCommand(t, command)
	}

	resp, err := http.Get(url + "/commands")
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	sb := string(body)

	for _, command := range commands {
		require.Contains(t, sb, command)
	}
}

// TestCommandDelete creates a command and then tries to delete it
func TestCommandDelete(t *testing.T) {
	const command = "cd /"

	saveResp := mustCreateCommand(t, command)

	req, err := http.NewRequest(http.MethodDelete, url+"/command/"+strconv.FormatInt(saveResp.ID, 10), nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
