#Collect system information
[[ -z ${HOST+x} ]] && HOST=localhost
CPU_INFO=""
MEM_INFO=""
OS_INFO=""
DISK_INFO=""
CTL_INFO=""

function get_cpu_info() {
  local host=$1
  local res=""
  local cpu_info="$(ssh "$host" "lscpu")"
  #local cpu_info="$("lscpu")"
  cpu_info="${cpu_info/\"/\\\"}"
  local res_obj="{\"cmd2check\": \"lscpu\""
  while read -r line; do
    arg=$(echo "$line" | sed 's/:.*$//g' )
    value=$(echo "$line" | sed 's/^.*: *//g' )
    res_obj="$res_obj, \"$arg\": \"$value\""
  done <<< "$cpu_info"
  res_obj="$res_obj, \"raw\": \"$cpu_info\""
  res_obj="${res_obj} }"
  CPU_INFO=$res_obj #$(jq -n "$res_obj")
}

function get_mem_info() {
  local host=$1
  local res=""
  local mem_info="$(ssh "$host" "cat /proc/meminfo")"
  #local mem_info="$(cat /proc/meminfo)"
  mem_info="${mem_info/\"/\\\"}"
  local res_obj="{\"cmd2check\": \"cat /proc/meminfo\""
  while read -r line; do
    arg=$(echo "$line" | sed 's/:.*$//g' )
    value=$(echo "$line" | sed 's/^.*: *//g' )
    res_obj="$res_obj, \"$arg\": \"$value\""
  done <<< "$mem_info"
  res_obj="$res_obj, \"raw\": \"$mem_info\""
  res_obj="${res_obj} }"
  MEM_INFO=$res_obj #$(jq -n "$res_obj")
}

function get_system_info() {
  local host=$1
  local sys_info="$(ssh "$host" "uname -a")"
  #local sys_info="$(uname -a | sed 's/"/\\"/g')"
  res_obj="{\"cmd2check\": \"uname -a\", \"raw\": \"$sys_info\"}"
  OS_INFO=$res_obj #$(jq -n "$res_obj")
}

function get_ctl_info() {
  local host=$1
  local sys_info="$(ssh "$host" "hostnamectl status")"
  #local ctl_info="$(hostnamectl status)"
  ctl_info="${ctl_info/\"/\\\"}"
  local res_obj="{\"cmd2check\": \"hostnamectl status\""
  while read -r line; do
    arg=$(echo "$line" | sed 's/:.*$//g' )
    value=$(echo "$line" | sed 's/^.*: *//g' )
    res_obj="$res_obj, \"$arg\": \"$value\""
  done <<< "$ctl_info"
  res_obj="$res_obj, \"raw\": \"$ctl_info\""
  res_obj="${res_obj} }"
  CTL_INFO=$res_obj #$(jq -n "$res_obj")
}

function get_disk_info() {
  local host=$1
  local sys_info="$(ssh "$host" "df -T")"
  #local disk_info="$(df -T | sed 's/"/\\"/g')"
  res_obj="{\"cmd2check\": \"df -T\", \"raw\": \"$disk_info\"}"
  DISK_INFO=$res_obj #$(jq -n "$res_obj")
}

not_first=false
get_cpu_info $HOST
get_mem_info $HOST
get_system_info $HOST
get_disk_info $HOST
get_ctl_info $HOST
host_obj="{\"cpu\": $CPU_INFO, \"ram\": $MEM_INFO, \"system\": $OS_INFO, \"disk\": $DISK_INFO}"
host_obj="{\"cpu\": $CPU_INFO, \"ram\": $MEM_INFO, \"system\": $OS_INFO, \"disk\": $DISK_INFO, \"virtualization\": $CTL_INFO }"
result="{\"$HOST\": ${host_obj}}"
result="[${result}]"

result=$(jq -n "$result")
echo "$result"
