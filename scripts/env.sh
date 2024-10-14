#!/bin/bash

while read -r LINE; do
    if [ "$LINE" != "" ] && [ "$LINE" != '#'* ]; then
        ENV_VAR="$(echo $LINE | envsubst)"
        export $ENV_VAR
    fi
done < $(pwd)/.env.dev
