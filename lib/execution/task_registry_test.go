package execution

import "testing"

type mockTask struct{}

func (m mockTask) Run(reportPayload string) error {
	return nil
}

func TestTaskRegistryRegisterAndGetHandler(t *testing.T) {
	registry := newRegistry()
	handler := mockTask{}

	err := registry.Register("sample-task", handler)
	if err != nil {
		t.Fatalf("unexpected register error: %v", err)
	}

	h, err := registry.GetHandler("sample-task")
	if err != nil {
		t.Fatalf("unexpected get handler error: %v", err)
	}
	if h == nil {
		t.Fatalf("expected non-nil handler")
	}

	if len(registry.List()) != 1 {
		t.Fatalf("expected one registered task, got: %d", len(registry.List()))
	}
}

func TestTaskRegistryRegisterValidations(t *testing.T) {
	registry := newRegistry()
	handler := mockTask{}

	if err := registry.Register("", handler); err == nil {
		t.Fatalf("expected error for empty task name")
	}

	if err := registry.Register("dup-task", handler); err != nil {
		t.Fatalf("unexpected error on first register: %v", err)
	}
	if err := registry.Register("dup-task", handler); err == nil {
		t.Fatalf("expected error for duplicate task name")
	}
}

func TestTaskRegistryGetHandlerMissing(t *testing.T) {
	registry := GetRegistry()

	h, err := registry.GetHandler("missing")
	if err == nil {
		t.Fatalf("expected error for missing task")
	}
	if h != nil {
		t.Fatalf("expected nil handler for missing task")
	}
}
