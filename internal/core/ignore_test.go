// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package core

import (
	"testing"
)

func TestIgnoreEngine_SecurityBlacklist(t *testing.T) {
	engine := NewIgnoreEngine(".")

	// Test default blacklist behavior (Secure by Default)
	sensitiveFiles := []string{".env", "config/.env.local", "id_rsa", "cert.pem", "secrets.json"}
	for _, file := range sensitiveFiles {
		if !engine.IsIgnored(file, false) {
			t.Errorf("Expected sensitive file %s to be ignored by default", file)
		}
	}

	// Test allowing secrets (Flexible by Choice)
	engine.SetAllowSecrets(true)
	for _, file := range sensitiveFiles {
		if engine.IsIgnored(file, false) {
			t.Errorf("Expected sensitive file %s to be allowed when allowSecrets is true", file)
		}
	}
}

func TestIgnoreEngine_ExplicitIncludes(t *testing.T) {
	engine := NewIgnoreEngine(".")

	// .env is blocked by default
	if !engine.IsIgnored(".env", false) {
		t.Errorf("Expected .env to be ignored by default")
	}

	// Explicitly include .env
	engine.SetExplicitIncludes(".env, src/")

	// Now .env should bypass the blacklist
	if engine.IsIgnored(".env", false) {
		t.Errorf("Expected .env to bypass blacklist due to explicit include")
	}
}
