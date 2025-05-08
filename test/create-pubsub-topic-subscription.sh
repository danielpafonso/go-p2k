# create topic
curl -X PUT http://localhost:8085/v1/projects/acme/topics/events

# create subscription
curl -H 'content-type: application/json' -X PUT -d '{"topic": "projects/acme/topics/events"}' \
  "http://localhost:8085/v1/projects/acme/subscriptions/sub-events"
