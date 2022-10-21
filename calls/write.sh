#!/bin/sh

URL="http://localhost:8891/api/v1/logs"

max=10
for (( i=0; i < max; i++ ))
do

payload=$(
  cat <<EOF
{
  "raw": "raw #${i}"
}
EOF
)

curl -X POST ${URL} \
    -H 'Content-Type: application/json' \
    -d "${payload}"

done

