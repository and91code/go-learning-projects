package dto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/seu-usuario/taskflow-backend/internal/types"
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
			value:    "2026-08-01 00:00",
			expected: "2026-08-01T00:00:00Z",
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
	for _, value := range []string{"2026-08-01", "2026-08-01 00:00:00", "2026-08-01T00:00:00Z"} {
		var request CreateTaskRequest
		err := json.Unmarshal([]byte(`{"due_date":"`+value+`"}`), &request)
		require.ErrorContains(t, err, "YYYY-MM-DD HH:mm")
	}
}

func TestCreateTaskRequest_ToDomainMapsTagIDs(t *testing.T) {
	var request CreateTaskRequest
	require.NoError(t, json.Unmarshal([]byte(`{"tag_ids":[7,9]}`), &request))

	task := request.ToDomain()

	require.Equal(t, []types.Tag{{ID: 7}, {ID: 9}}, task.Tags)
	require.NotNil(t, task.Tags)
}

func TestUpdateTaskRequest_ToDomainDistinguishesOmittedAndEmptyTagIDs(t *testing.T) {
	var omitted UpdateTaskRequest
	require.NoError(t, json.Unmarshal([]byte(`{}`), &omitted))
	require.Nil(t, omitted.ToDomain().Tags)

	var empty UpdateTaskRequest
	require.NoError(t, json.Unmarshal([]byte(`{"tag_ids":[]}`), &empty))
	require.NotNil(t, empty.ToDomain().Tags)
	require.Empty(t, empty.ToDomain().Tags)
}

func TestTaskResponseFromDomain_IncludesTags(t *testing.T) {
	task := types.Task{
		ID:        42,
		DueDate:   timePointer(mustParseTime(t, "2026-08-01T12:34:00Z")),
		CreatedAt: mustParseTime(t, "2026-08-01T10:11:00Z"),
		UpdatedAt: mustParseTime(t, "2026-08-01T10:12:00Z"),
		Tags:      []types.Tag{{ID: 7, Name: "urgente", Color: "#ff0000"}},
	}

	response := TaskResponseFromDomain(task)

	require.Equal(t, []TagResponse{{ID: 7, Name: "urgente", Color: "#ff0000"}}, response.Tags)

	encoded, err := json.Marshal(response)
	require.NoError(t, err)
	require.JSONEq(t, `{"id":42,"list_id":0,"title":"","status":"","priority":"","due_date":"2026-08-01 12:34","created_at":"2026-08-01 10:11","updated_at":"2026-08-01 10:12","tags":[{"id":7,"name":"urgente","color":"#ff0000"}]}`, string(encoded))
}

func TestTaskResponseFromDomain_UsesEmptyTagsArray(t *testing.T) {
	response := TaskResponseFromDomain(types.Task{})

	require.NotNil(t, response.Tags)
	require.Empty(t, response.Tags)
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339, value)
	require.NoError(t, err)
	return parsed
}

func timePointer(value time.Time) *time.Time {
	return &value
}
