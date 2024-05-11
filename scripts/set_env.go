package main

import (
	"fmt"
	"os"
)

func main() {
	envVars := []struct {
		name        string
		description string
	}{
		{
			name:        "CONFIG_PATH",
			description: "path to config file",
		},
		{
			name:        "POSTGRES_HOST",
			description: "Postgres database host",
		},
		{
			name:        "POSTGRES_PORT",
			description: "Postgres database port",
		},
		{
			name:        "POSTGRES_USER",
			description: "Postgres database user",
		},
		{
			name:        "POSTGRES_PASSWORD",
			description: "Postgres database password",
		},
		{
			name:        "POSTGRES_SSL_MODE",
			description: "Postgres database SSL mode",
		},
		{
			name:        "POSTGRES_DB",
			description: "Postgres database name",
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

		fmt.Printf("Insert %s: ", envVar.description)
		_, err := fmt.Scan(&value)
		if err != nil {
			panic(err)
		}

		_, err = file.WriteString(fmt.Sprintf("%s=%s\n", envVar.name, value))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Environment variables set")
}
