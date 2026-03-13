package common

import "testing"

func TestExecutorType_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		et    ExecutorType
		valid bool
	}{
		{"ssh is valid", ExecutorSSH, true},
		{"local is valid", ExecutorLocal, true},
		{"unknown is invalid", ExecutorType("unknown"), false},
		{"empty is invalid", ExecutorType(""), false},
		{"mixed case is invalid", ExecutorType("SSH"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.et.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}
