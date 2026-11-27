#!/usr/bin/env bash
set -euo pipefail
reports="$(realpath ../reports/runtime)"
export REPORTS="$reports"
script_dir="$(cd ../scripts && pwd)"
go version > "$reports/go.txt"

loopback_compose() {
  python3 - "$1" <<'PY'
import pathlib, re, sys
p = pathlib.Path(sys.argv[1])
text = p.read_text()
text = re.sub(r'(?m)^(\s*-\s*")(?=(?:\d+|\$\{[A-Z_]+:-\d+\}):\d+")', r'\g<1>127.0.0.1:', text)
p.write_text(text)
PY
}

case "$1" in
  wukong)
    go test -tags=e2e -count=1 -timeout=20m -json ./test/e2e/message/cross_node_delivery \
      -run '^TestThreeNodeClusterCrossNodeUsersExchangeMessages$' | tee "$reports/exchange.jsonl"
    jq -se 'any(.[]; .Test == "TestThreeNodeClusterCrossNodeUsersExchangeMessages" and .Action == "pass")' "$reports/exchange.jsonl"
    ;;
  juggle)
    cleanup() {
      docker compose ps -a > "$reports/container-status.txt" || true
      docker compose logs --no-color --tail=120 > "$reports/compose.log" 2>&1 || true
      docker compose down -v --remove-orphans || true
    }
    trap cleanup EXIT
    loopback_compose docker-compose.yml
    # Fresh bootstrap must reach the final embedded migration version marker.
    python3 - <<'PY'
from pathlib import Path
p = Path('sql/imserver.sql')
s = p.read_text()
bad = '+CREATE TABLE IF NOT EXISTS `inboxmsgs_split`'
assert s.count(bad) == 1
p.write_text(s.replace(bad, bad[1:]))
PY
    git diff -- docker-compose.yml sql/imserver.sql > "$reports/fixture.patch"
    export BENCH_CLIENTS=20 BENCH_RATE=40 BENCH_GROUP_SENDERS=4
    export BENCH_WARMUP=5s BENCH_DURATION=20s BENCH_DELIVERY_GRACE=10s
    export BENCH_OUTPUT_DIR="$reports"
    docker compose config --format json | jq -e 'all(.services[].ports[]?; .host_ip == "127.0.0.1")'
    bash scripts/run-benchmark.sh all
    for scenario in private group; do
      jq -e --arg scenario "$scenario" '
        .workload.scenario == $scenario and
        all(.connections, .message_acknowledgements, .deliveries;
          .attempted > 0 and .succeeded == .attempted and .failed == 0)
      ' "$reports"/*-"$scenario".json > /dev/null
    done
    ;;
  openim)
    server_pid=
    cleanup() {
      if [[ -n "$server_pid" ]]; then kill "$server_pid" 2>/dev/null || true; fi
      docker compose ps -a > "$reports/container-status.txt" || true
      docker compose down -v --remove-orphans || true
    }
    trap cleanup EXIT
    loopback_compose docker-compose.yml
    sed -i 's/^queue: kafka$/queue: redis/' config/share.yml
    sed -i 's/^listenIP: 0\.0\.0\.0$/listenIP: 127.0.0.1/' config/openim-msggateway.yml
    sed -i 's/listenIP: 0\.0\.0\.0/listenIP: 127.0.0.1/' config/openim-api.yml
    git diff -- docker-compose.yml config > "$reports/fixture.patch"
    docker compose config --format json | jq -e 'all(.services.mongodb.ports[], .services.redis.ports[], .services.minio.ports[]; .host_ip == "127.0.0.1")'
    docker compose up -d mongodb redis minio
    go build -o "$RUNNER_TEMP/openim-source-server" ./cmd
    "$RUNNER_TEMP/openim-source-server" -c "$PWD/config" > "$reports/server.log" 2>&1 &
    server_pid=$!
    python3 "$script_dir/im-source-runtime.py"
    ;;
  *)
    printf 'Unknown product: %s\n' "$1" >&2
    exit 2
    ;;
esac
