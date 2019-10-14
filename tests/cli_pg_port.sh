#!/bin/bash
# don't allow invalid arguments

export PATH=$PATH:${BASH_SOURCE%/*}/..

# put invalid '--force' argument
output=$(./checkup --ssh-hostname postgres --pg-port 5432 2>&1)

if [[ $output =~ "'--pg-port' must be set only with '--pg-hostname'" ]]; then
  echo -e "\e[36mOK\e[39m"
else
  >&2 echo -e "\e[31mFAILED\e[39m"
  >&2 echo -e "Output: $output"
  exit 1
fi

