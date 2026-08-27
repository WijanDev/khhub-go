package config

import "testing"

func TestEnvFlags(t *testing.T) {
	t.Parallel()
	cases := []struct {
		env        string
		production bool
		reset      bool
		demo       bool
		strict     bool
	}{
		{env: "development", reset: true, demo: true},
		{env: "staging", demo: true, strict: true},
		{env: "production", production: true, strict: true},
	}
	for _, tc := range cases {
		c := Config{AppEnv: tc.env}
		if c.Production() != tc.production {
			t.Fatalf("%s Production: got %v", tc.env, c.Production())
		}
		if c.AllowsSeedReset() != tc.reset {
			t.Fatalf("%s AllowsSeedReset: got %v", tc.env, c.AllowsSeedReset())
		}
		if c.AutoDemoSeed() != tc.demo {
			t.Fatalf("%s AutoDemoSeed: got %v", tc.env, c.AutoDemoSeed())
		}
		if c.StrictSecrets() != tc.strict {
			t.Fatalf("%s StrictSecrets: got %v", tc.env, c.StrictSecrets())
		}
	}
}
