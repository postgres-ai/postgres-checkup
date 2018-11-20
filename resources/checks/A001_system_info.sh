#Collect system information
HOSTS=(localhost)
CPU_INFO=""
MEM_INFO=""

function get_cpu_info() {
  local host=$1
  local res=""
  #local cpu_info="$(ssh "$host" "lscpu" | sed 's/"/\\"/g')"
  local cpu_info="$("lscpu" | sed 's/"/\\"/g')"
  local res_obj="{\"cmd2check\": \"lscpu\""
  while read -r line; do
    arg=$(echo "$line" | sed 's/:.*$//g' )
    value=$(echo "$line" | sed 's/^.*: *//g' )
    res_obj="$res_obj, \"$arg\": \"$value\""
  done <<< "$cpu_info"
  res_obj="${res_obj} }"
  CPU_INFO=$(jq -n "$res_obj")
}

function get_mem_info() {
  local host=$1
  local res=""
  #local mem_info="$(ssh "$host" "cat /proc/meminfo" | sed 's/"/\\"/g')"
  local mem_info="$(cat /proc/meminfo | sed 's/"/\\"/g')"
  local res_obj="{\"cmd2check\": \"meminfo\""
  while read -r line; do
    arg=$(echo "$line" | sed 's/:.*$//g' )
    value=$(echo "$line" | sed 's/^.*: *//g' )
    res_obj="$res_obj, \"$arg\": \"$value\""
  done <<< "$mem_info"
  res_obj="${res_obj} }"
  MEM_INFO=$(jq -n "$res_obj")
}

result="["
not_first=false
for host in ${HOSTS[@]}; do
    if [[ "$not_first" = true ]]; then
      result="${result}, "
    else
      not_first=true
    fi
    CPU_INFO=""
    MEM_INFO=""
    get_cpu_info $host
    get_mem_info $host
    host_obj="[$CPU_INFO,$MEM_INFO]"
    result="${result}{\"$host\": ${host_obj}}"
done
result="${result}]"

result=$(jq -n "$result")
echo "$result"
