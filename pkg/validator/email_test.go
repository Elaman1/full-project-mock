package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"", false},
		{"   ", false},
		{"plainaddress", false},
		{"@missingusername.com", false},
		{"user@", false},
		{"user@site", false},
		{"user@site.com", true},
		{"USER@SITE.COM", true},   // проверим приведение к lower
		{" user@site.com ", true}, // проверим trim
		{"user.name+tag@sub.domain.co.uk", true},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			assert.Equal(t, tt.valid, IsValidEmail(tt.email))
		})
	}
}
