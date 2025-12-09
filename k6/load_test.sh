#!/bin/bash

# URL of the endpoint
URL="http://localhost:9000/test/mail"

# JSON payload
PAYLOAD='{"to":["tharunkumarra@gmail.com"],"subject":"Load Test","type":"welcome"}'

# Number of requests to send
NUM_REQUESTS=5

# Concurrency level (how many requests to send in parallel)
CONCURRENCY=5

# Function to send a request
send_request() {
  echo "Sending request..."
  curl -s -o /dev/null -w "%{http_code}\n" -X POST -H "Content-Type: application/json" -d "$PAYLOAD" "$URL"
}

# Send requests in parallel
for i in $(seq 1 $NUM_REQUESTS)
do
  # Limit the number of concurrent jobs
  if [[ $(jobs -r -p | wc -l) -ge $CONCURRENCY ]]; then
    # Wait for any job to finish
    wait -n
  fi
  send_request &
done

# Wait for all remaining background jobs to finish
wait
echo "All requests sent."
