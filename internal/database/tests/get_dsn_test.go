package database_tests

import (
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/database"
	"os"
	"path/filepath"
	"testing"
)

// в тестах мы не должны ссылаться на глобальный .env
// а также при изменении переменных окружения должны возвращать их в
// значение до тестов
func TestGetDsn(t *testing.T) {
	envVar := "TEST_GET_DSN"
	badEnvVar := "BAD_NAME"

	tests := []struct {
		name        string
		envVar      string
		loadGlobal  bool
		badPath     bool
		globalEmpty bool
		wantDsn     string
		wantErr     error
	}{
		{
			name:       "success global env",
			envVar:     envVar,
			loadGlobal: true,
			wantDsn:    "testdsn_global",
			wantErr:    nil,
		},
		{
			name:    "success local env",
			envVar:  envVar,
			wantDsn: "testdsn_local",
			wantErr: nil,
		},
		{
			name:       "success global env, bad path",
			envVar:     envVar,
			badPath:    true,
			loadGlobal: true,
			wantDsn:    "testdsn_global",
			wantErr:    nil,
		},
		{
			name:    "err, bad path",
			badPath: true,
			envVar:  envVar,
			wantDsn: "",
			wantErr: apperrors.ErrInvalidFileLoad,
		},
		{
			name:    "err, bad name",
			envVar:  badEnvVar,
			wantDsn: "",
			wantErr: apperrors.ErrBadEnvVariableName,
		},
		{
			name:        "empty global",
			envVar:      envVar,
			loadGlobal:  true,
			globalEmpty: true,
			wantDsn:     "",
			wantErr:     apperrors.ErrBadEnvironmentConfiguration,
		},
		{
			name:    "err, empty global, empty .env",
			envVar:  badEnvVar,
			wantDsn: "",
			wantErr: apperrors.ErrBadEnvVariableName,
		},
	}

	val, exists := os.LookupEnv(envVar)
	badVal, badExists := os.LookupEnv(badEnvVar)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				if exists {
					os.Setenv(envVar, val)
				} else {
					os.Unsetenv(envVar)
				}

				if badExists {
					os.Setenv(badEnvVar, badVal)
				} else {
					os.Unsetenv(badEnvVar)
				}
			})

			if tt.loadGlobal {
				if tt.globalEmpty {
					t.Setenv(envVar, "")
				} else {
					t.Setenv(envVar, "testdsn_global")
				}
			} else {
				os.Unsetenv(envVar)
			}

			if badExists {
				os.Unsetenv(badEnvVar)
			}

			path := t.TempDir()
			envPath := filepath.Join(path, "config.env")
			err := os.WriteFile(
				envPath,
				[]byte("TEST_GET_DSN=testdsn_local"), 0644,
			)
			if err != nil {
				t.Fatalf("failed to write in config.env %v", err)
			}

			if tt.badPath {
				envPath = ""
			}

			dsn, err := database.GetDsn(envPath, tt.envVar)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected err to be %v but got %v", tt.wantErr, err)
			}

			if dsn != tt.wantDsn {
				t.Errorf("expected dsn to be %v but got %v", tt.wantDsn, dsn)
			}
		})
	}
}
