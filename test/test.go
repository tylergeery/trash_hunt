package test

import (
	"fmt"
	"os"
)

// SetVars for testing
func SetVars() {
	// set persistent storage
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASS", "test")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_TABLE", "test")
	os.Setenv("DB_SSL_MODE", "verify")

	fmt.Println("ENV VARIABLES BEING SET")

	// set temporary storage
	os.Setenv("REDIS_USER", "test")
	os.Setenv("REDIS_SECRET", "test")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("REDIS_DB_NUMBER", "0")
}
