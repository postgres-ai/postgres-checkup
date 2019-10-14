#!/bin/bash
# don't allow invalid arguments

export PATH=$PATH:${BASH_SOURCE%/*}/..

# put invalid '--force' argument
output=$(./checkup --hostname postgres --ssh-hostname postgres 2>&1)

if [[ $output =~ "only one of options '--hostname', '--ssh-hostname' or '--pg-hostname' must be set" ]]; then
  echo -e "\e[36mOK\e[39m"
else
  >&2 echo -e "\e[31mFAILED\e[39m"
  >&2 echo -e "Output: $output"
  exit 1
fi

