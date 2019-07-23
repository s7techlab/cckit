#!/usr/bin/env bash
if [[ -z "$1" ]]; then echo "\$1 - path to cpaper issue payload"; exit 1;  fi
source _common.sh

POST /cpaper/issue @$1