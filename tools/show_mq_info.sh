PORT=6600
HOST="localhost"
DB=3
EXECUTE_CMD="redis-cli -p $PORT -n $DB"

cmd="smembers rmq::connections"
conns=`$EXECUTE_CMD $cmd`
echo "Connections: $conns"
echo

cmd="smembers rmq::queues"
queues=`$EXECUTE_CMD $cmd`
echo "Queues: $queues"
echo

for queue_name in $queues; do
  echo "Ready queue name: $queue_name"

  cmd="llen rmq::queue::[$queue_name]::ready"
  num_elems=`$EXECUTE_CMD $cmd`
  echo "  Number of elements: $num_elems"

  cmd="lrange rmq::queue::[$queue_name]::ready 0 -1"
  elems=`$EXECUTE_CMD $cmd`
  echo "  Elements: $elems"
done

echo

for queue_name in $queues; do
  echo "Rejected queue name: $queue_name"

  cmd="llen rmq::queue::[$queue_name]::rejected"
  num_elems=`$EXECUTE_CMD $cmd`
  echo "  Number of elements: $num_elems"

  cmd="lrange rmq::queue::[$queue_name]::rejected 0 -1"
  elems=`$EXECUTE_CMD $cmd`
  echo "  Elements: $elems"
done
