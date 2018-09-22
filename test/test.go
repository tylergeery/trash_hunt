package test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// SetVars for testing
func SetVars() {
	cmd := exec.Command("docker", "inspect", "--format", "'{{ .NetworkSettings.IPAddress }}'", "trash-hunt-redis")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		log.Fatal(err)
	}

	// set persistent storage
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASS", "test")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_TABLE", "test")
	os.Setenv("DB_SSL_MODE", "verify")

	// set temporary storage
	os.Setenv("REDIS_HOST", strings.Trim(string(output), "\n'"))
	os.Setenv("REDIS_PORT", "6379")
}
