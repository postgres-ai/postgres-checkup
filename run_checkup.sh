#!/bin/bash

echo "Collect the first set of snapshots (only for K*** reports)..."

HOSTS=$CHECKUP_HOSTS
if [[ ! -z ${PG_CHECKUP_HOSTS+x} ]]; then
  HOSTS=$PG_CHECKUP_HOSTS
fi
if [[ ! -z ${SSH_CHECKUP_HOSTS+x} ]]; then
  HOSTS=$SSH_CHECKUP_HOSTS
fi

for host in $CHECKUP_HOSTS; do
  ./checkup collect \
    --config "${CHECKUP_CONFIG_PATH}" \
    --hostname "${host}" \
    --file resources/checks/K000_query_analysis.sh
done

echo "The first set of snapshots has been created. Wait ${CHECKUP_SNAPSHOT_DISTANCE_SECONDS} seconds..."
sleep "${CHECKUP_SNAPSHOT_DISTANCE_SECONDS}"
# the distance ^^^ recommended: large enough to get good data, at least 10 minutes

echo "Collect the second set of snapshots and build reports..."
for host in $HOSTS; do
  ./checkup collect \
    --config "${CHECKUP_CONFIG_PATH}" \
    --hostname "${host}"
done

echo "Generate human-readable reports..."
./checkup process --config "${CHECKUP_CONFIG_PATH}"

echo "Upload the report to Postgres.ai platform..."
./checkup upload --config "${CHECKUP_CONFIG_PATH}"
