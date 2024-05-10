package delete_test

import (
	"github.com/fatalistix/postgres-command-executor/internal/database"
	"github.com/fatalistix/postgres-command-executor/internal/domain/wrapper"
	"github.com/fatalistix/postgres-command-executor/internal/lib/syncmap"
	"github.com/fatalistix/postgres-command-executor/internal/services/process/delete"
	"github.com/fatalistix/postgres-command-executor/internal/services/process/delete/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"os/exec"
	"testing"
)

type syncMapEntries struct {
	Key   uuid.UUID
	Value *wrapper.CmdWrapper
}

func TestDeleteProcess(t *testing.T) {
	cases := []struct {
		name            string
		processID       uuid.UUID
		mockError       error
		mockCalledTimes int
		entries         []syncMapEntries
	}{
		{
			name:            "process not found",
			processID:       uuid.New(),
			mockError:       database.ErrProcessNotFound,
			mockCalledTimes: 1,
			entries:         []syncMapEntries{},
		},
		{
			name:            "success with only one process",
			processID:       uuid.New(),
			mockError:       nil,
			mockCalledTimes: 1,
			entries: []syncMapEntries{
				// uuid.Nil will be replaced by processID
				{uuid.Nil, &wrapper.CmdWrapper{Cmd: exec.Command("ls")}},
			},
		},
		{
			name:            "success with multiple processes",
			processID:       uuid.New(),
			mockError:       nil,
			mockCalledTimes: 1,
			entries: []syncMapEntries{
				{uuid.Nil, &wrapper.CmdWrapper{Cmd: exec.Command("/bin/bash", "-c", "ls; sleep 10; ls")}},
				{uuid.New(), &wrapper.CmdWrapper{Cmd: exec.Command("ls")}},
				{uuid.New(), &wrapper.CmdWrapper{Cmd: exec.Command("ls")}},
				{uuid.New(), &wrapper.CmdWrapper{Cmd: exec.Command("ls")}},
			},
		},
		{
			name:            "success with no processes",
			processID:       uuid.New(),
			mockError:       nil,
			mockCalledTimes: 1,
			entries:         []syncMapEntries{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			processDeleterMock := mocks.NewProcessDeleter(t)

			processDeleterMock.
				On("DeleteProcess", tc.processID).
				Return(tc.mockError).
				Times(tc.mockCalledTimes)

			// fill sync map
			sm := syncmap.NewSyncMap[uuid.UUID, *wrapper.CmdWrapper]()
			syncMapSize := len(tc.entries)
			for _, entry := range tc.entries {
				stdout, err := entry.Value.Cmd.StdoutPipe()
				require.NoError(t, err)
				stderr, err := entry.Value.Cmd.StderrPipe()
				require.NoError(t, err)
				entry.Value.Stdout = stdout
				entry.Value.Stderr = stderr
				if entry.Key == uuid.Nil {
					entry.Key = tc.processID
					err := entry.Value.Cmd.Start()
					require.NoError(t, err)
					syncMapSize--
				}
				sm.Store(entry.Key, entry.Value)
			}

			processDeleter := delete.NewService(processDeleterMock, sm)

			err := processDeleter.DeleteProcess(tc.processID)
			require.ErrorIs(t, err, tc.mockError)

			// check sync map
			require.Equal(t, syncMapSize, sm.Size())

			for _, entry := range tc.entries {
				_, ok := sm.Load(entry.Key)
				if entry.Key == uuid.Nil {
					require.False(t, ok)
				} else {
					require.True(t, ok)
				}
			}
		})
	}
}
