package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gambruh/simplevault/internal/config"
)

func TestNewRootCmd(t *testing.T) {
	// Coverage: verifies Cobra flag parsing, env precedence, and run wiring.
	tests := []struct {
		name   string
		env    map[string]string
		args   []string
		expect func(t *testing.T) config.ClientConfig
	}{
		{
			name: "flags configure client",
			args: []string{
				"-a", "127.0.0.1:9443",
				"-s", "client.pem",
				"-p", "client.key",
				"--localstorage", "./tmp/local",
				"-t", "15s",
				"--bininputfolder", "./tmp/in",
				"--binoutputfolder", "./tmp/out",
			},
			expect: func(t *testing.T) config.ClientConfig {
				return config.ClientConfig{
					Address:         "127.0.0.1:9443",
					ClientCert:      expectedPath(t, "client.pem"),
					ServerCert:      "cert.pem",
					PrivateKey:      expectedPath(t, "client.key"),
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
			args: []string{
				"-a", "127.0.0.1:9443",
				"-s", "client.pem",
				"-p", "client.key",
				"--localstorage", "./tmp/local",
				"-t", "15s",
				"--bininputfolder", "./tmp/in",
				"--binoutputfolder", "./tmp/out",
			},
			expect: func(t *testing.T) config.ClientConfig {
				return config.ClientConfig{
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
			unsetCliEnv(t, envKeys...)
			setCliEnvMap(t, tc.env)

			oldCfg := config.ClientCfg
			t.Cleanup(func() { config.ClientCfg = oldCfg })

			var got config.ClientConfig
			cmd := newRootCmd(func(cfg config.ClientConfig) error {
				got = cfg
				return nil
			})
			cmd.SetArgs(tc.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			want := tc.expect(t)
			if got != want {
				t.Fatalf("runner cfg = %+v, want %+v", got, want)
			}
			if config.ClientCfg != want {
				t.Fatalf("config.ClientCfg = %+v, want %+v", config.ClientCfg, want)
			}
		})
	}
}

func expectedPath(t *testing.T, relPath string) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	return filepath.Join(filepath.Dir(wd), relPath)
}

func setCliEnvMap(t *testing.T, env map[string]string) {
	t.Helper()

	for key, value := range env {
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("Setenv(%q) error = %v", key, err)
		}

		keyCopy := key
		t.Cleanup(func() {
			_ = os.Unsetenv(keyCopy)
		})
	}
}

func unsetCliEnv(t *testing.T, keys ...string) {
	t.Helper()

	for _, key := range keys {
		keyCopy := key
		original, ok := os.LookupEnv(keyCopy)
		if err := os.Unsetenv(keyCopy); err != nil {
			t.Fatalf("Unsetenv(%q) error = %v", keyCopy, err)
		}

		t.Cleanup(func() {
			if !ok {
				_ = os.Unsetenv(keyCopy)
				return
			}
			_ = os.Setenv(keyCopy, original)
		})
	}
}
