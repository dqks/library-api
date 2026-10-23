package database_tests

import (
	"library-api/internal/apperrors"
	"library-api/internal/database"
	"os"
	"testing"
)

func TestGetDsn(t *testing.T) {
	tests := []struct {
		name       string
		envPath    string
		envVar     string
		loadGlobal bool
		wantDsn    string
		wantErr    error
	}{
		{
			name:    "success local env",
			envPath: "../../../.env",
			envVar:  "TEST_GET_DSN",
			wantDsn: "testdsn",
			wantErr: nil,
		},
		{
			name:       "success global env",
			envPath:    "../../../.env",
			envVar:     "TEST_GET_DSN",
			loadGlobal: true,
			wantDsn:    "testdsn",
			wantErr:    nil,
		},
		{
			name:       "success global env, bad path",
			envPath:    "../.env",
			envVar:     "TEST_GET_DSN",
			loadGlobal: true,
			wantDsn:    "testdsn",
			wantErr:    nil,
		},
		{
			name:    "err, bad path",
			envPath: "",
			envVar:  "TEST_GET_DSN",
			wantDsn: "",
			wantErr: apperrors.ErrBadEnvVariablePath,
		},
		{
			name:    "err, bad name",
			envPath: "../../../.env",
			envVar:  "BAD_NAME",
			wantDsn: "",
			wantErr: apperrors.ErrBadEnvVariableName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				os.Unsetenv("TEST_GET_DSN")
			})

			if tt.loadGlobal {
				os.Setenv("TEST_GET_DSN", "testdsn")
			}

			dsn, err := database.GetDsn(tt.envPath, tt.envVar)

			if err != tt.wantErr {
				t.Errorf("expected err to be %v but got %v", tt.wantErr, err)
			}

			if dsn != tt.wantDsn {
				t.Errorf("expected dsn to be %v but bot %v", tt.wantDsn, dsn)
			}
		})
	}
}
