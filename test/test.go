package test

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// SetVars for testing
func init() {
	cmd := exec.Command("docker", "inspect", "--format", "'{{ .NetworkSettings.IPAddress }}'", "trash-hunt-redis")
	redis, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command("docker", "inspect", "--format", "'{{ .NetworkSettings.IPAddress }}'", "trash-hunt-pg")
	pg, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}

	// set persistent storage
	os.Setenv("DB_USER", "dev")
	os.Setenv("DB_PASS", "dev_secret")
	os.Setenv("DB_HOST", strings.Trim(string(pg), "\n'"))
	os.Setenv("DB_TABLE", "dev_secret")
	os.Setenv("DB_SSL_MODE", "disable")

	// set temporary storage
	os.Setenv("REDIS_HOST", strings.Trim(string(redis), "\n'"))
	os.Setenv("REDIS_PORT", "6379")
}

// FatalOnError to be reused by testing package
func FatalOnError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Unexpected err: %s", err.Error())
	}
}

// ExpectEqualString expects  values to be equal
func ExpectEqualString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("Expected equal string values: (%s) (%s)", expected, actual)
	}
}

// ExpectEqualInt64s expects int64 values to be equal
func ExpectEqualInt64s(t *testing.T, expected, actual int64) {
	if expected != actual {
		t.Fatalf("Expected equal int64 values: (%d) (%d)", expected, actual)
	}
}
