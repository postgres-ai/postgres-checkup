# check installed useful system diagnostics utilites

if [[ "${SSH_SUPPORT}" = "false" ]]; then
  echo "SSH is not supported, skipping..." >&2
  exit 1
fi

# define lists
cpu_list="ps,htop,top,mpstat,lscpu"
disk_usage_list="df,du"
io_list="pidstat,iostat,iotop,ftrace,blktrace"
memory_list="free,vmstat"
network_list="tcpdump,netstat,ss,iptraf,ethtool"
misc_list="dstat,strace,ltrace,perf,numastat"

#######################################
# Check utilites from list in remote server, make json
# Globals:
#   None
# Arguments:
#   json_obj_name, list
# Returns:
#   json
#######################################
check_list() {
  local json_obj_name="$1" # name of list
  local list="$2" # comma-separated list of utilites
  local json="\"${json_obj_name}\": {"
  SAVE_IFS="$IFS"
  IFS=","
  local cnt="0"
  local comma=""
  for util in $list ; do
    [[ "$cnt" -eq "0" ]] && comma="" || comma=","
    IFS="$SAVE_IFS" # non-standart IFS ruins ${CHECK_HOST_CMD}
    local res=$(${CHECK_HOST_CMD} "sudo which $util 2>&1")
    res=$(echo "$res" | grep -v "\[sudo\] password for ")
    if [[ "$res" != "" ]]; then
      json="${json}${comma} \"$util\": \"yes\""
    else
      json="${json}${comma} \"$util\": \"\""
    fi
    cnt=$(( cnt + 1 ))
  done
  json="${json} }"
  echo "$json"
}

res1=$(${CHECK_HOST_CMD} "which sudo")
res2=$(${CHECK_HOST_CMD} "sudo which sudo")
res2=$(echo "$res2" | grep -v "\[sudo\] password for ")
if [[ "$res1" != "$res2" ]]; then
  errmsg "ERROR: Can not execute 'which' on target server."
fi

# build json object to stdout
echo "{"
check_list "cpu" "$cpu_list" && echo -n ","
check_list "free_space" "$disk_usage_list" && echo -n ","
check_list "io" "$io_list" && echo -n ","
check_list "memory" "$memory_list" && echo -n ","
check_list "network" "$network_list" && echo -n ","
check_list "misc" "$misc_list"
echo "}"

