#!/usr/bin/env bash
if [[ -z "$1" ]]; then echo "\$1 - issuer"; exit 1;  fi
if [[ -z "$2" ]]; then echo "\$2 - paper number"; exit 1;  fi
source _common.sh

GET /cpaper/$1/$2