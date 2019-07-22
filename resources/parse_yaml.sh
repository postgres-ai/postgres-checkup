#!/bin/bash

#######################################
# Parse simple YAML file
# Globals:
#   None
# Arguments:
#   (text) path to yaml file
# Returns:
#   None
#######################################
function parse_yaml() {
  local yaml_file=$1
  local prefix=$2
  local s
  local w
  local fs

  s='[[:space:]]*'
  w='[a-zA-Z0-9_.-]*'
  fs="$(echo @|tr @ '\034')"

  (
    sed -ne '/^--/s|--||g; s|\"|\\\"|g; s/\s*$//g;' \
      -e "/#.*[\"\']/!s| #.*||g; /^#/s|#.*||g;" \
      -e  "s|^\($s\)\($w\)$s:$s\"\(.*\)\"$s\$|\1$fs\2$fs\3|p" \
      -e "s|^\($s\)\($w\)$s[:-]$s\(.*\)$s\$|\1$fs\2$fs\3|p" |
    awk -F"$fs" '{
      indent = length($1)/2;
      if (length($2) == 0) { conj[indent]="+";} else {conj[indent]="";}
      vname[indent] = $2;
      for (i in vname) {if (i > indent) {delete vname[i]}}
        if (length($3) > 0) {
          vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
          printf("%s%s%s%s=(\"%s\")\n", "'"$prefix"'",vn, $2, conj[indent-1],$3);
        }
      }' |
    sed -e 's/_=/+=/g' |
    awk 'BEGIN {
         FS="=";
         OFS="="
       }
       /(-|\.).*=/ {
         gsub("-|\\.", "_", $1)
       }
       { print }'
  ) < "$yaml_file"
}

eval $(parse_yaml $1 "$2") 2>&1 > /dev/null