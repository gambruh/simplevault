package main

import (
	"os"
	"testing"

	"github.com/gambruh/simplevault/internal/config"
)

func TestNewRootCmd(t *testing.T) {
	// Coverage: verifies Cobra flag parsing and run wiring.
	tests := []struct {
		name   string
		env    map[string]string
		args   []string
		expect config.Config
	}{
		{
			name: "flags configure server",
			args: []string{
				"-a", "127.0.0.1:9090",
				"--cert", "server.crt",
				"--privatekey", "server.key",
				"-d", "postgres://flag-db",
				"-k", "flag-hash",
			},
			expect: config.Config{
				Address:     "127.0.0.1:9090",
				Certificate: "server.crt",
				PrivateKey:  "server.key",
				Database:    "postgres://flag-db",
				Key:         "flag-hash",
			},
		},
		{
			name: "env overrides flags",
			env: map[string]string{
				"GK_ADDRESS":     "10.0.0.1:443",
				"GK_CERT":        "env.crt",
				"GK_PRIVATE_KEY": "env.key",
				"GK_DATABASE":    "postgres://env-db",
				"GK_HASHKEY":     "env-hash",
			},
			args: []string{
				"-a", "127.0.0.1:9090",
				"--cert", "server.crt",
				"--privatekey", "server.key",
				"-d", "postgres://flag-db",
				"-k", "flag-hash",
			},
			expect: config.Config{
				Address:     "10.0.0.1:443",
				Certificate: "env.crt",
				PrivateKey:  "env.key",
				Database:    "postgres://env-db",
				Key:         "env-hash",
			},
		},
	}

	envKeys := []string{"GK_ADDRESS", "GK_CERT", "GK_PRIVATE_KEY", "GK_DATABASE", "GK_HASHKEY"}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			configUnsetEnv(t, envKeys...)
			configSetEnvMap(t, tc.env)

			oldCfg := config.Cfg
			t.Cleanup(func() { config.Cfg = oldCfg })

			called := false
			cmd := newRootCmd(func() error {
				called = true
				return nil
			})
			cmd.SetArgs(tc.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if !called {
				t.Fatal("run func was not called")
			}
			if config.Cfg != tc.expect {
				t.Fatalf("config.Cfg = %+v, want %+v", config.Cfg, tc.expect)
			}
		})
	}
}

func configSetEnvMap(t *testing.T, env map[string]string) {
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

func configUnsetEnv(t *testing.T, keys ...string) {
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
