#!/usr/bin/env bash

ENV_FILE=".env"

echo "==> Configuring PostgreSQL credentials"
echo -n "Enter database host: "
read -r host

echo -n "Enter database port: "
read -r port

echo -n "Enter username: "
read -r username

echo -n "Enter password: "
read -r password

echo -n "Enter database name: "
read -r dbname

echo -n "Enter database ssl mode (disable, require, verify-full, verify-ca): "
read -r sslmode

{
  echo "POSTGRES_HOST=$host";
  echo "POSTGRES_PORT=$port"
  echo "POSTGRES_USER=$username";
  echo "POSTGRES_PASSWORD=$password";
  echo "POSTGRES_DBNAME=$dbname";
  echo "POSTGRES_SSLMODE=$sslmode";
} >> $ENV_FILE

echo "==> Credentials stored successfully"
