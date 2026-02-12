package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSetClientConfig(t *testing.T) {
	// Coverage: verifies flags fallback, env override, and global assignment.
	tests := []struct {
		name      string
		env       map[string]string
		flags     ClientFlagConfig
		want      func(t *testing.T) ClientConfig
		expectErr bool
	}{
		{
			name: "flags are used when env is absent",
			flags: ClientFlagConfig{
				Address:         "127.0.0.1:9443",
				ClientCert:      "client.pem",
				PrivateKey:      "client.key",
				LocalStorage:    "./tmp/local",
				BinInputFolder:  "./tmp/in",
				BinOutputFolder: "./tmp/out",
				CheckTime:       15 * time.Second,
			},
			want: func(t *testing.T) ClientConfig {
				return ClientConfig{
					Address:         "127.0.0.1:9443",
					ClientCert:      expectedClientPath(t, "client.pem"),
					ServerCert:      "cert.pem",
					PrivateKey:      expectedClientPath(t, "client.key"),
					LocalStorage:    "./tmp/local",
					BinInputFolder:  "./tmp/in",
					UserDataFolder:  "./userdata",
					UserDataFile:    "./userdata/user.json",
					BinOutputFolder: "./tmp/out",
					CheckTime:       15 * time.Second,
				}
			},
		},
		{
			name: "env overrides flags",
			env: map[string]string{
				"GK_ADDRESS":         "10.0.0.2:443",
				"GK_PUBLICKEY":       "env-client.pem",
				"GK_CERT":            "env-server.pem",
				"GK_PRIVATEKEY":      "env-client.key",
				"GK_LOCALSTORAGE":    "./env/local",
				"GK_BINARIES_INPUT":  "./env/in",
				"GK_BINARIES_OUTPUT": "./env/out",
				"GK_CHECKINTERVAL":   "42s",
			},
			flags: ClientFlagConfig{
				Address:         "127.0.0.1:9443",
				ClientCert:      "client.pem",
				PrivateKey:      "client.key",
				LocalStorage:    "./tmp/local",
				BinInputFolder:  "./tmp/in",
				BinOutputFolder: "./tmp/out",
				CheckTime:       15 * time.Second,
			},
			want: func(t *testing.T) ClientConfig {
				return ClientConfig{
					Address:         "10.0.0.2:443",
					ClientCert:      "env-client.pem",
					ServerCert:      "env-server.pem",
					PrivateKey:      "env-client.key",
					LocalStorage:    "./env/local",
					BinInputFolder:  "./env/in",
					UserDataFolder:  "./userdata",
					UserDataFile:    "./userdata/user.json",
					BinOutputFolder: "./env/out",
					CheckTime:       42 * time.Second,
				}
			},
		},
		{
			name: "invalid env returns error",
			env: map[string]string{
				"GK_CHECKINTERVAL": "bad-duration",
			},
			flags:     DefaultClientFlags(),
			expectErr: true,
		},
	}

	envKeys := []string{
		"GK_ADDRESS",
		"GK_PUBLICKEY",
		"GK_CERT",
		"GK_PRIVATEKEY",
		"GK_LOCALSTORAGE",
		"GK_BINARIES_INPUT",
		"GK_BINARIES_OUTPUT",
		"GK_CHECKINTERVAL",
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			unsetEnv(t, envKeys...)
			setEnvMap(t, tc.env)

			oldCfg := ClientCfg
			t.Cleanup(func() { ClientCfg = oldCfg })

			got, err := SetClientConfig(tc.flags)
			if tc.expectErr {
				if err == nil {
					t.Fatal("SetClientConfig() error = nil, want non-nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("SetClientConfig() error = %v", err)
			}

			want := tc.want(t)
			if got != want {
				t.Fatalf("cfg = %+v, want %+v", got, want)
			}
			if ClientCfg != want {
				t.Fatalf("ClientCfg = %+v, want %+v", ClientCfg, want)
			}
		})
	}
}

func expectedClientPath(t *testing.T, relPath string) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	return filepath.Join(filepath.Dir(wd), relPath)
}
