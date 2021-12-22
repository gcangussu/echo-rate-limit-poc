go build ./client_src/client.go
chmod +x ./client

pids=()

DURATION=15 # seconds

# ./client <id> <req/s> <duration (s)>

./client A 4 $DURATION > out_a_1 &
pids+=($!)

./client B 20 $DURATION > out_b &
pids+=($!)

./client A 4 $DURATION > out_a_2 &
pids+=($!)

./client A 3 $DURATION > out_a_3 &
pids+=($!)

# wait for all pids
for pid in ${pids[*]}; do
  wait $pid
done
