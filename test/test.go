package test

import (
	"log"
	"os"
	"os/exec"
	"strings"
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
