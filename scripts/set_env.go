package main

import (
	"fmt"
	"os"
)

func main() {
	envVars := []struct {
		name         string
		description  string
		defaultValue string
	}{
		{
			name:         "CONFIG_PATH",
			description:  "path to config file",
			defaultValue: "./config/prod.json",
		},
		{
			name:         "POSTGRES_HOST",
			description:  "Postgres database host (local: localhost, docker: db)",
			defaultValue: "localhost",
		},
		{
			name:         "POSTGRES_PORT",
			description:  "Postgres database port",
			defaultValue: "5432",
		},
		{
			name:         "POSTGRES_USER",
			description:  "Postgres database user",
			defaultValue: "command-executor-owner",
		},
		{
			name:         "POSTGRES_PASSWORD",
			description:  "Postgres database password",
			defaultValue: "command-executor-password",
		},
		{
			name:         "POSTGRES_SSL_MODE",
			description:  "Postgres database SSL mode",
			defaultValue: "disable",
		},
		{
			name:         "POSTGRES_DB",
			description:  "Postgres database name",
			defaultValue: "command-executor",
		},
	}
	file, err := os.OpenFile(".env", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = file.Close()
	}()

	for _, envVar := range envVars {
		var value string

		fmt.Printf("Insert %s (default: %s): ", envVar.description, envVar.defaultValue)
		_, err := fmt.Scan(&value)
		if err != nil {
			panic(err)
		}

		if value == "" {
			value = envVar.defaultValue
		}

		_, err = file.WriteString(fmt.Sprintf("%s=%s\n", envVar.name, value))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Environment variables set")
}
