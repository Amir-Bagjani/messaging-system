package environmentvariable

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func NewEnv() {
	if err := LoadEnv(".env"); err != nil {
		log.Fatalf("Error loading environment file %v", err)
	}

	requiredVars := []string{"DATABASE_URL"}
	if err := ValidateEnv(requiredVars); err != nil {
		fmt.Println(err)
	}
}

func LoadEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening .env file %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line in .env file %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if e := os.Setenv(key, value); e != nil {
			return fmt.Errorf("error setting environment variable %s, %s", key, value)
		}
	}

	if e := scanner.Err(); e != nil {
		return fmt.Errorf("error reading .env file %v", e)
	}

	return nil
}

func ValidateEnv(vars []string) error {
	for _, k := range vars {
		if _, exists := os.LookupEnv(k); !exists {
			return fmt.Errorf("environment variable %s not set", k)
		}
	}

	return nil
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
