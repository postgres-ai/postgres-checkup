#!/bin/bash
# don't allow invalid arguments

export PATH=$PATH:${BASH_SOURCE%/*}/..

# put invalid '--force' argument
output=$(./checkup -h postgres --username ${POSTGRES_USER} --force --project test --dbname ${POSTGRES_DB} 2>&1)

if [[ $output =~ "invalid argument" ]]; then
  echo -e "\e[36mOK\e[39m"
else
  >&2 echo -e "\e[31mFAILED\e[39m"
  >&2 echo -e "Output: $output"
  exit 1
fi

