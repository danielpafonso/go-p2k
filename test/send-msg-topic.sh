encode=$(echo -n $1 | base64)
printf '{"messages": [{"data": "%s"}]}' $encode

curl -H "Content-Type: application/json" \
  "localhost:8085/v1/projects/acme/topics/events:publish" \
  -d  "{\"messages\": [{\"data\": \"$encode\"}]}"
