#!/bin/bash
# Check that with remote SSH calls, when --psql-binary "sudo psql" is used,
# statement_timeout and appname are set properly.

export PATH=$PATH:${BASH_SOURCE%/*}/..

# First, let's allow remote ssh connection to ourselves, via '-h localhost123'

sudo sh -c "echo \"127.0.0.1 localhost123\" > /etc/hosts"

rm -rf /tmp/_checkup_check_loc.key
ssh-keygen -t rsa -N "" -f /tmp/_checkup_check_loc.key
cat /tmp/_checkup_check_loc.key.pub >> ~/.ssh/authorized_keys

# put invalid '--force' argument
output=$(./checkup \
  -h localhost123 \
  --ssh-identity-file /tmp/_checkup_check_loc.key \
  -U postgres \
  --project test \
  -d test \
  -e 1 \
  2>&1 \
)

if [[ $output =~ "invalid argument" ]]; then
  echo -e "\e[36mOK\e[39m"
else
  >&2 echo -e "\e[31mFAILED\e[39m"
  >&2 echo -e "Output: $output"
  exit 1
fi

