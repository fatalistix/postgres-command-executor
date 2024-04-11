#!/usr/bin/env bash

ENV_FILE=".env"

echo "==> Configuring config path env variable"
echo -n "Enter path to config file: "
read -r CONFIG_PATH
echo "CONFIG_PATH=$CONFIG_PATH" >> $ENV_FILE
echo "==> Path configured successfully"
