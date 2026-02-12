package config

import (
	"os"
	"testing"
)

func TestSetConfig(t *testing.T) {
	// Coverage: verifies flag fallback and env override precedence.
	tests := []struct {
		name   string
		env    map[string]string
		flags  FlagConfig
		expect Config
	}{
		{
			name: "flags are used when env is absent",
			flags: FlagConfig{
				Address:     "127.0.0.1:9090",
				Certificate: "server.crt",
				PrivateKey:  "server.key",
				Key:         "hash-key",
				Database:    "postgres://from-flag",
			},
			expect: Config{
				Address:     "127.0.0.1:9090",
				Certificate: "server.crt",
				PrivateKey:  "server.key",
				Key:         "hash-key",
				Database:    "postgres://from-flag",
			},
		},
		{
			name: "env overrides flags",
			env: map[string]string{
				"GK_ADDRESS":     "10.0.0.1:443",
				"GK_CERT":        "env.crt",
				"GK_PRIVATE_KEY": "env.key",
				"GK_HASHKEY":     "env-hash",
				"GK_DATABASE":    "postgres://from-env",
			},
			flags: FlagConfig{
				Address:     "127.0.0.1:9090",
				Certificate: "server.crt",
				PrivateKey:  "server.key",
				Key:         "hash-key",
				Database:    "postgres://from-flag",
			},
			expect: Config{
				Address:     "10.0.0.1:443",
				Certificate: "env.crt",
				PrivateKey:  "env.key",
				Key:         "env-hash",
				Database:    "postgres://from-env",
			},
		},
	}

	envKeys := []string{"GK_ADDRESS", "GK_CERT", "GK_PRIVATE_KEY", "GK_HASHKEY", "GK_DATABASE"}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			unsetEnv(t, envKeys...)
			setEnvMap(t, tc.env)

			oldCfg := Cfg
			t.Cleanup(func() { Cfg = oldCfg })

			if err := SetConfig(tc.flags); err != nil {
				t.Fatalf("SetConfig() error = %v", err)
			}

			if Cfg != tc.expect {
				t.Fatalf("Cfg = %+v, want %+v", Cfg, tc.expect)
			}
		})
	}
}

func setEnvMap(t *testing.T, env map[string]string) {
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

func unsetEnv(t *testing.T, keys ...string) {
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
