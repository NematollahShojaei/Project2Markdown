// Copyright 2026 R3D HILLS. All Rights Reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package core

import (
	"testing"

	"github.com/nematollahshojaei/project2markdown/v2/internal/utils"
)

// BenchmarkEstimateTokens ensures the custom tokenizer remains lightning-fast.
// Performance Pillar: Must process millions of characters without memory allocation overhead.
func BenchmarkEstimateTokens(b *testing.B) {
	sampleText := "func main() { fmt.Println(\"Hello, Enterprise AI Context Engine!\") }"
	b.ResetTimer()
	// Zero-Bug Policy: Modernized benchmark loop for Go 1.24+
	for b.Loop() {
		utils.EstimateTokens(sampleText)
	}
}

// BenchmarkIgnoreEngine ensures the new Security Blacklist and Explicit Includes
// do not introduce performance bottlenecks during massive directory walks.
func BenchmarkIgnoreEngine(b *testing.B) {
	engine := NewIgnoreEngine(".")
	engine.SetExplicitIncludes("src/, *.go")

	b.ResetTimer()
	// Zero-Bug Policy: Modernized benchmark loop for Go 1.24+
	for b.Loop() {
		// Test a safe file
		engine.IsIgnored("src/main.go", false)
		// Test a blocked sensitive file
		engine.IsIgnored(".env", false)
	}
}
