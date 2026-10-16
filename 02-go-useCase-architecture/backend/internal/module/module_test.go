package module

import "testing"

func TestBuildProvidesRouteHandlers(t *testing.T) {
	modules := Build(nil)
	if modules.Auth == nil {
		t.Fatal("Build returned no auth handler")
	}
	if modules.Lists == nil {
		t.Fatal("Build returned no list handler")
	}
	if modules.Tasks == nil {
		t.Fatal("Build returned no task handler")
	}
}
