#!/usr/bin/env bash

#RESULT=${PWD##*/}
#RESULT=${RESULT:-/}
#
#if [ "$RESULT" != "postgres-command-executor" ]; then
#  echo "==> Possibly you are in a wrong directory. You must run '.sh' files only from project root directory"
#  exit 1
#fi

ENV_FILE=".env"

if [ -f $ENV_FILE ]; then
  echo "==> Clearing '$ENV_FILE' file"
  echo "" > $ENV_FILE
  echo "==> '$ENV_FILE' was cleared"
else
	echo "==> Creating '$ENV_FILE' file"
	touch $ENV_FILE
	echo "==> '$ENV_FILE' created successfully"
fi