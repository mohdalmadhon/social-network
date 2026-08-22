package validation

import "testing"

func TestValidateName(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
		label string
	}{
		{name: "Noa", valid: true, label: "normal mixed-case name"},
		{name: "A", valid: false, label: "too short"},
		{name: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", valid: false, label: "too long"},
		{name: "Noa2", valid: false, label: "contains a number"},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			err := validateName(test.name)
			if test.valid && err != nil {
				t.Fatalf("expected valid name, got %v", err)
			}
			if !test.valid && err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
