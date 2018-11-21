#Collect system information
HOSTS=(localhost)
CPU_INFO=""
MEM_INFO=""
OS_INFO=""
DISK_INFO=""
CTL_INFO=""

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
  res_obj="$res_obj, \"raw\": \"$cpu_info\""
  res_obj="${res_obj} }"
  CPU_INFO=$res_obj #$(jq -n "$res_obj")
}

function get_mem_info() {
  local host=$1
  local res=""
  #local mem_info="$(ssh "$host" "cat /proc/meminfo" | sed 's/"/\\"/g')"
  local mem_info="$(cat /proc/meminfo | sed 's/"/\\"/g')"

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
  #local sys_info="$(ssh "$host" "uname -a" | sed 's/"/\\"/g')"
  local sys_info="$(uname -a | sed 's/"/\\"/g')"
  res_obj="{\"cmd2check\": \"uname -a\", \"raw\": \"$sys_info\"}"
  OS_INFO=$res_obj #$(jq -n "$res_obj")
}

function get_ctl_info() {
  local host=$1
  #local sys_info="$(ssh "$host" "hostnamectl status" | sed 's/"/\\"/g')"
  local ctl_info="$(hostnamectl status | sed 's/"/\\"/g')"
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
  #local sys_info="$(ssh "$host" "dt -T" | sed 's/"/\\"/g')"
  local disk_info="$(df -T | sed 's/"/\\"/g')"
  res_obj="{\"cmd2check\": \"df -T\", \"raw\": \"$disk_info\"}"
  DISK_INFO=$res_obj #$(jq -n "$res_obj")
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
    OS_INFO=""
    DISK_INFO=""
    CTL_INFO=""
    get_cpu_info $host
    get_mem_info $host
    get_system_info $host
    get_disk_info $host
    get_ctl_info $host
    host_obj="{\"cpu\": $CPU_INFO, \"ram\": $MEM_INFO, \"system\": $OS_INFO, \"disk\": $DISK_INFO}"
    host_obj="{\"cpu\": $CPU_INFO, \"ram\": $MEM_INFO, \"system\": $OS_INFO, \"disk\": $DISK_INFO, \"virtualization\": $CTL_INFO }"
    result="${result}{\"$host\": ${host_obj}}"
done
result="${result}]"

result=$(jq -n "$result")
echo "$result"
