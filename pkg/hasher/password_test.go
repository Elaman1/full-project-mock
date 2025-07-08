package hasher

import "testing"

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"simple password", "hello123"},
		{"unicode password", "Пароль123"},
		{"symbols", "!@#$%^&*()"},
		{"long password", string(make([]byte, 1024))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := HashPassword(tt.password)
			if err != nil {
				t.Errorf("HashPassword() error = %v", err)
			}

			if h == tt.password {
				t.Errorf("HashPassword() = %v, want %v", h, tt.password)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	const MyPassword = "hello123"
	hash, err := HashPassword(MyPassword)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	tests := []struct {
		name     string
		password string
		mustErr  bool
	}{
		{"correct password", MyPassword, false},
		{"wrong password", "wrong", true},
		{"empty password", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = Verify(hash, tt.password)
			if (err != nil) != tt.mustErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.mustErr)
			}
		})
	}
}

func TestVerify_InvalidHashFormat(t *testing.T) {
	err := Verify("this-is-not-a-valid-format", "pass")
	if err == nil {
		t.Errorf("expected error on invalid hash format")
	}
}

func TestPasswordHashingFlow(t *testing.T) {
	password := "SecureP@ssword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}

	if err = Verify(hash, password); err != nil {
		t.Errorf("expected match, got error: %v", err)
	}

	if err = Verify(hash, "wrong"); err == nil {
		t.Errorf("expected mismatch error")
	}
}
