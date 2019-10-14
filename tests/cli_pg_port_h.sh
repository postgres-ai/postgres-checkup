#!/bin/bash
# don't allow invalid arguments

export PATH=$PATH:${BASH_SOURCE%/*}/..

# put invalid '--force' argument
output=$(./checkup --pg-hostname postgres --ssh-port 22 2>&1)

if [[ $output =~ "'--ssh-port' must be set only with '--ssh-hostname'" ]]; then
  echo -e "\e[36mOK\e[39m"
else
  >&2 echo -e "\e[31mFAILED\e[39m"
  >&2 echo -e "Output: $output"
  exit 1
fi

