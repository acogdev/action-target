package monitor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `hosts = ["localhost", "example.com", "192.168.1.1"]
port = 8080
interval = 30`

	tempFile := createTempConfigFile(t, configContent)
	defer os.Remove(tempFile)

	// Test reading the config
	config := ReadConfig(tempFile)

	// Verify the parsed values
	expectedHosts := []string{"localhost", "example.com", "192.168.1.1"}
	if len(config.Hosts) != len(expectedHosts) {
		t.Errorf("Expected %d hosts, got %d", len(expectedHosts), len(config.Hosts))
	}
	for i, host := range expectedHosts {
		if config.Hosts[i] != host {
			t.Errorf("Expected host[%d] to be %s, got %s", i, host, config.Hosts[i])
		}
	}

	if config.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", config.Port)
	}

	if config.Interval != 30 {
		t.Errorf("Expected interval 30, got %d", config.Interval)
	}
}

func TestReadConfigEmptyHosts(t *testing.T) {
	configContent := `hosts = []
port = 3000
interval = 60`

	tempFile := createTempConfigFile(t, configContent)
	defer os.Remove(tempFile)

	config := ReadConfig(tempFile)

	if len(config.Hosts) != 0 {
		t.Errorf("Expected empty hosts array, got %v", config.Hosts)
	}
	if config.Port != 3000 {
		t.Errorf("Expected port 3000, got %d", config.Port)
	}
	if config.Interval != 60 {
		t.Errorf("Expected interval 60, got %d", config.Interval)
	}
}

func TestReadConfigMinimalConfig(t *testing.T) {
	configContent := `hosts = ["test.com"]
port = 80
interval = 5`

	tempFile := createTempConfigFile(t, configContent)
	defer os.Remove(tempFile)

	config := ReadConfig(tempFile)

	if len(config.Hosts) != 1 || config.Hosts[0] != "test.com" {
		t.Errorf("Expected hosts ['test.com'], got %v", config.Hosts)
	}
	if config.Port != 80 {
		t.Errorf("Expected port 80, got %d", config.Port)
	}
	if config.Interval != 5 {
		t.Errorf("Expected interval 5, got %d", config.Interval)
	}
}

func TestReadConfigNonExistentFile(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when reading non-existent file")
		}
	}()

	ReadConfig("non_existent_file.toml")
}

func TestReadConfigInvalidTOML(t *testing.T) {
	configContent := `hosts = ["localhost"
port = invalid
interval = 30`

	tempFile := createTempConfigFile(t, configContent)
	defer os.Remove(tempFile)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when parsing invalid TOML")
		}
	}()

	ReadConfig(tempFile)
}

func TestReadConfigMissingFields(t *testing.T) {
	configContent := `hosts = ["localhost"]`

	tempFile := createTempConfigFile(t, configContent)
	defer os.Remove(tempFile)

	config := ReadConfig(tempFile)

	// Should have default zero values for missing fields
	if config.Port != 0 {
		t.Errorf("Expected default port 0, got %d", config.Port)
	}
	if config.Interval != 0 {
		t.Errorf("Expected default interval 0, got %d", config.Interval)
	}
}

func TestReadConfigEmptyFile(t *testing.T) {
	tempFile := createTempConfigFile(t, "")
	defer os.Remove(tempFile)

	config := ReadConfig(tempFile)

	// Should have default zero values
	if len(config.Hosts) != 0 {
		t.Errorf("Expected empty hosts, got %v", config.Hosts)
	}
	if config.Port != 0 {
		t.Errorf("Expected default port 0, got %d", config.Port)
	}
	if config.Interval != 0 {
		t.Errorf("Expected default interval 0, got %d", config.Interval)
	}
}

// Helper function to create temporary config files for testing
func createTempConfigFile(t *testing.T, content string) string {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_config.toml")
	
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	
	return tempFile
}

