#!/usr/bin/env bash
#
# Runs wrk at several concurrency levels and prints one comparable table.
#
# Latency and throughput need different concurrency: a saturated server reports
# queue depth, not service time. Read p50/p99 from the low rows and rps from the
# high ones.
#
#   bash devops/bench.sh                                  # default ladder on /houses
#   bash devops/bench.sh -s devops/bench.lua              # mixed traffic
#   bash devops/bench.sh -d 30s -c "10 25 50"             # custom duration/ladder
#   bash devops/bench.sh -u http://localhost:8080/api/v1/types

set -euo pipefail

URL="http://localhost:8080/api/v1/houses?page=1&page_size=20"
SCRIPT=""
DURATION="10s"
THREADS=4
LADDER="25 50 100 200"

usage() { sed -n '3,12p' "${BASH_SOURCE[0]}" | cut -c3-; exit 0; }

while getopts "u:s:d:t:c:h" opt; do
    case "$opt" in
        u) URL="$OPTARG" ;;
        s) SCRIPT="$OPTARG" ;;
        d) DURATION="$OPTARG" ;;
        t) THREADS="$OPTARG" ;;
        c) LADDER="$OPTARG" ;;
        h) usage ;;
        *) exit 2 ;;
    esac
done

command -v wrk >/dev/null || { echo "❌ wrk не установлен: brew install wrk" >&2; exit 1; }

BASE_URL="${URL%%/api/*}"
curl -fsS -o /dev/null "$BASE_URL/api/v1/healthcheck" \
    || { echo "❌ Приложение не отвечает на $BASE_URL — запустите make start" >&2; exit 1; }

wrk_args=(-t"$THREADS" -d"$DURATION" --latency -H "Accept-Encoding: gzip")
[ -n "$SCRIPT" ] && wrk_args+=(-s "$SCRIPT")

echo "target:   ${SCRIPT:-$URL}"
echo "duration: $DURATION per step, threads: $THREADS"
echo
printf "%-6s %-12s %-10s %-10s %-10s %s\n" "conn" "rps" "p50" "p90" "p99" "non-2xx"
printf -- "------------------------------------------------------------------\n"

for conn in $LADDER; do
    out="$(wrk "${wrk_args[@]}" -c"$conn" "$URL" 2>&1)"

    rps="$(awk '/Requests\/sec/{print $2}' <<<"$out")"
    p50="$(awk '$1=="50%"{print $2}' <<<"$out")"
    p90="$(awk '$1=="90%"{print $2}' <<<"$out")"
    p99="$(awk '$1=="99%"{print $2}' <<<"$out")"
    bad="$(awk '/Non-2xx or 3xx/{print $4}' <<<"$out")"

    printf "%-6s %-12s %-10s %-10s %-10s %s\n" \
        "$conn" "${rps:-—}" "${p50:-—}" "${p90:-—}" "${p99:-—}" "${bad:-0}"
done
