for i in {1..100}; do
    curl -sS -X GET "http://localhost:9000/test"
    sleep 3
done
