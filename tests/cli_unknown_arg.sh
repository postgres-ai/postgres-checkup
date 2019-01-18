#!/bin/bash

# don't allow invalid arguments

# Output styles (only BOLD is supported by default GNU screen)
BOLD=`tput md 2>/dev/null` || :
RESET=`tput me 2>/dev/null` || :

export PATH=$PATH:${BASH_SOURCE%/*}/..

output=$(./check --force 2>&1)

if [[ $output =~ "invalid argument" ]]; then
  echo -e "OK"
else
  >&2 echo -e "${BOLD}FAILED${RESET}"
  >&2 echo -e "Output: $output"
  exit 1
fi

