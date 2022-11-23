PORT=6600
HOST="localhost"
DB=3
EXECUTE_CMD="redis-cli -p $PORT -n $DB"

cmd="smembers rmq::connections"
conns=`$EXECUTE_CMD $cmd`
echo "Connections: $conns"
echo

