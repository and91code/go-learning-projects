package dto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateTaskRequest_UnmarshalDueDate(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "mysql datetime",
			value:    "2026-08-01 00:00:00",
			expected: "2026-08-01T00:00:00Z",
		},
		{
			name:     "rfc3339",
			value:    "2026-08-01T00:00:00-03:00",
			expected: "2026-08-01T00:00:00-03:00",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var request CreateTaskRequest
			err := json.Unmarshal([]byte(`{"due_date":"`+test.value+`"}`), &request)
			require.NoError(t, err)
			require.NotNil(t, request.DueDate)
			require.True(t, time.Time(*request.DueDate).Equal(mustParseTime(t, test.expected)))

			domainTask := request.ToDomain()
			require.NotNil(t, domainTask.DueDate)
			require.True(t, domainTask.DueDate.Equal(mustParseTime(t, test.expected)))
		})
	}
}

func TestCreateTaskRequest_UnmarshalInvalidDueDate(t *testing.T) {
	var request CreateTaskRequest
	err := json.Unmarshal([]byte(`{"due_date":"2026-08-01"}`), &request)

	require.ErrorContains(t, err, "RFC3339")
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339, value)
	require.NoError(t, err)
	return parsed
}
