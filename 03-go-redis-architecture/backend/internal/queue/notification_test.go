package queue

import "testing"

func TestNotificationActionAndMessage(t *testing.T) {
	for _, test := range []struct {
		job, resource, action, message string
	}{
		{"LIST_CREATE", "list", "created", "Lista criada com sucesso."},
		{"TASK_UPDATE", "task", "updated", "Tarefa atualizada com sucesso."},
		{"TAG_DELETE", "tag", "deleted", "Tag excluída com sucesso."},
	} {
		resource, action := notificationAction(test.job)
		if resource != test.resource || action != test.action {
			t.Fatalf("notificationAction(%q) = (%q, %q)", test.job, resource, action)
		}
		if message := notificationMessage(resource, action); message != test.message {
			t.Fatalf("notificationMessage = %q, want %q", message, test.message)
		}
	}
}
