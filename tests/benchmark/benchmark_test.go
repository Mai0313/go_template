package benchmark

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"go-template/core/config"
	"go-template/internal/testutil"
)

// BenchmarkConfigLoad benchmarks config loading performance
func BenchmarkConfigLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := config.Load()
		if err != nil {
			b.Fatalf("Config load failed: %v", err)
		}
	}
}

// BenchmarkConfigLoadFromFile benchmarks loading config from file
func BenchmarkConfigLoadFromFile(b *testing.B) {
	// Create a config file for benchmarking
	configContent := `{
		"version": "1.0.0",
		"environment": "benchmark",
		"log_level": "info",
		"debug": false
	}`

	configPath, cleanup := testutil.CreateTempConfig(b, configContent)
	defer cleanup()

	// Set environment variable
	cleanup2 := testutil.SetEnv(b, "CONFIG_PATH", configPath)
	defer cleanup2()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := config.Load()
		if err != nil {
			b.Fatalf("Config load failed: %v", err)
		}
	}
}

// BenchmarkConfigParsing benchmarks JSON parsing performance
func BenchmarkConfigParsing(b *testing.B) {
	configJSON := `{
		"version": "1.0.0",
		"environment": "benchmark",
		"log_level": "info",
		"debug": false
	}`

	data := []byte(configJSON)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var cfg config.Config
		err := json.Unmarshal(data, &cfg)
		if err != nil {
			b.Fatalf("JSON unmarshal failed: %v", err)
		}
	}
}

// BenchmarkConfigMarshaling benchmarks JSON marshaling performance
func BenchmarkConfigMarshaling(b *testing.B) {
	cfg := config.Config{
		Version:     "1.0.0",
		Environment: "benchmark",
		LogLevel:    "info",
		Debug:       false,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(cfg)
		if err != nil {
			b.Fatalf("JSON marshal failed: %v", err)
		}
	}
}

// BenchmarkFileOperations benchmarks file I/O operations
func BenchmarkFileOperations(b *testing.B) {
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	configContent := `{
		"version": "1.0.0",
		"environment": "benchmark",
		"log_level": "info",
		"debug": false
	}`

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Write file
		err := os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			b.Fatalf("Write file failed: %v", err)
		}
		b.StartTimer()

		// Read file (this is what we're benchmarking)
		_, err = os.ReadFile(configPath)
		if err != nil {
			b.Fatalf("Read file failed: %v", err)
		}

		b.StopTimer()
		// Clean up
		os.Remove(configPath)
		b.StartTimer()
	}
}

// BenchmarkConfigLoadWithDifferentSizes benchmarks config loading with different file sizes
func BenchmarkConfigLoadWithDifferentSizes(b *testing.B) {
	sizes := []struct {
		name    string
		content string
	}{
		{
			name: "Small",
			content: `{
				"version": "1.0.0",
				"environment": "production",
				"log_level": "info",
				"debug": false
			}`,
		},
		{
			name: "Medium",
			content: `{
				"version": "1.0.0",
				"environment": "production",
				"log_level": "info",
				"debug": false,
				"extra_field_1": "some value",
				"extra_field_2": "another value",
				"extra_field_3": "yet another value",
				"extra_field_4": "more data",
				"extra_field_5": "even more data"
			}`,
		},
		{
			name: "Large",
			content: func() string {
				cfg := map[string]interface{}{
					"version":     "1.0.0",
					"environment": "production",
					"log_level":   "info",
					"debug":       false,
				}
				// Add many extra fields
				for i := 0; i < 100; i++ {
					cfg[fmt.Sprintf("extra_field_%d", i)] = fmt.Sprintf("value_%d", i)
				}
				data, _ := json.Marshal(cfg)
				return string(data)
			}(),
		},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			configPath, cleanup := testutil.CreateTempConfig(b, size.content)
			defer cleanup()

			cleanup2 := testutil.SetEnv(b, "CONFIG_PATH", configPath)
			defer cleanup2()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := config.Load()
				if err != nil {
					b.Fatalf("Config load failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkConcurrentConfigLoad benchmarks concurrent config loading
func BenchmarkConcurrentConfigLoad(b *testing.B) {
	configContent := `{
		"version": "1.0.0",
		"environment": "concurrent",
		"log_level": "info",
		"debug": false
	}`

	configPath, cleanup := testutil.CreateTempConfig(b, configContent)
	defer cleanup()

	cleanup2 := testutil.SetEnv(b, "CONFIG_PATH", configPath)
	defer cleanup2()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := config.Load()
			if err != nil {
				b.Fatalf("Config load failed: %v", err)
			}
		}
	})
}

// BenchmarkConfigMemoryUsage measures memory allocation during config operations
func BenchmarkConfigMemoryUsage(b *testing.B) {
	configContent := `{
		"version": "1.0.0",
		"environment": "memory-test",
		"log_level": "info",
		"debug": false
	}`

	configPath, cleanup := testutil.CreateTempConfig(b, configContent)
	defer cleanup()

	cleanup2 := testutil.SetEnv(b, "CONFIG_PATH", configPath)
	defer cleanup2()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := config.Load()
		if err != nil {
			b.Fatalf("Config load failed: %v", err)
		}
	}
}
