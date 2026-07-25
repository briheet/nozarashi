#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
REPOSITORY_ROOT="$(CDPATH= cd -- "$SCRIPT_DIR/../.." && pwd)"

NOZARASHI_BIN="${NOZARASHI_BIN:-$REPOSITORY_ROOT/bin/nozarashi}"
COMPOSE_FILE="${COMPOSE_FILE:-$SCRIPT_DIR/compose.yaml}"
HEALTH_URL="${HEALTH_URL:-http://127.0.0.1:18080/api/v1/health}"
HEALTH_TIMEOUT="${HEALTH_TIMEOUT:-45}"

cd "$SCRIPT_DIR"

waitForHealth() {
	local deadline=$((SECONDS + HEALTH_TIMEOUT))

	while ((SECONDS < deadline)); do
		if curl --fail --silent --output /dev/null "$HEALTH_URL"; then
			return 0
		fi
		sleep 0.2
	done

	echo "health check timed out after ${HEALTH_TIMEOUT}s: $HEALTH_URL" >&2
	return 1
}

runPhase() {
	local tool="$1"
	local phase="$2"

	case "$tool:$phase" in
	nozarashi:build)
		"$NOZARASHI_BIN" build
		;;
	nozarashi:start)
		"$NOZARASHI_BIN" up
		waitForHealth
		;;
	nozarashi:stop)
		"$NOZARASHI_BIN" down
		;;
	docker:build)
		docker compose --file "$COMPOSE_FILE" down --remove-orphans
		docker compose --file "$COMPOSE_FILE" pull postgres redis
		docker compose --file "$COMPOSE_FILE" build backend
		;;
	docker:start)
		docker compose --file "$COMPOSE_FILE" up --detach --no-build
		waitForHealth
		;;
	docker:stop)
		docker compose --file "$COMPOSE_FILE" stop
		;;
	*)
		echo "unknown benchmark phase: $tool $phase" >&2
		return 1
		;;
	esac
}

# Internal phase execution is measured by the parent process.
if [[ "${1:-}" == "__phase" ]]; then
	runPhase "$2" "$3"
	exit
fi

ITERATIONS="${1:-3}"
RESULT_FILE="${2:-$SCRIPT_DIR/result.csv}"

if [[ ! "$ITERATIONS" =~ ^[1-9][0-9]*$ ]]; then
	echo "iterations must be a positive integer" >&2
	exit 1
fi

for commandName in container curl docker; do
	if ! command -v "$commandName" >/dev/null 2>&1; then
		echo "required command is not available: $commandName" >&2
		exit 1
	fi
done

if [[ ! -x "$NOZARASHI_BIN" ]]; then
	echo "nozarashi binary is not executable: $NOZARASHI_BIN" >&2
	exit 1
fi

docker compose --file "$COMPOSE_FILE" config --quiet

if ! container system dns list | grep -Fxq "simple-docker"; then
	echo "project DNS is missing; run sudo ../../bin/nozarashi create first" >&2
	exit 1
fi

# Prepare Nozarashi resources outside the measured phases.
"$NOZARASHI_BIN" create >/dev/null

# Stop benchmark-owned services before checking whether another stack owns the port.
"$NOZARASHI_BIN" down >/dev/null 2>&1 || true
docker compose --file "$COMPOSE_FILE" stop >/dev/null 2>&1 || true

if curl --silent --output /dev/null --max-time 1 "$HEALTH_URL"; then
	echo "port 18080 is already serving another stack; stop it before benchmarking" >&2
	exit 1
fi

mkdir -p "$(dirname -- "$RESULT_FILE")"

writeCSVRow() {
	local value
	local separator=""

	for value in "$@"; do
		value="${value//\"/\"\"}"
		printf '%s"%s"' "$separator" "$value" >>"$RESULT_FILE"
		separator=","
	done
	printf '\n' >>"$RESULT_FILE"
}

: >"$RESULT_FILE"
writeCSVRow \
	"timestamp_utc" \
	"host_os" \
	"host_arch" \
	"tool" \
	"tool_version" \
	"iteration" \
	"phase" \
	"command" \
	"exit_code" \
	"real_seconds" \
	"user_seconds" \
	"system_seconds" \
	"max_rss_bytes" \
	"health_url" \
	"health_status" \
	"cache_policy" \
	"volume_policy"

HOST_OS="$(uname -s)"
HOST_ARCH="$(uname -m)"
DOCKER_VERSION="$(docker compose version --short)"
NOZARASHI_VERSION="$(git -C "$REPOSITORY_ROOT" rev-parse --short HEAD 2>/dev/null || printf 'development')"
if [[ -n "$(git -C "$REPOSITORY_ROOT" status --porcelain 2>/dev/null)" ]]; then
	NOZARASHI_VERSION="${NOZARASHI_VERSION}-dirty"
fi

LAST_REAL="0"
LAST_USER="0"
LAST_SYSTEM="0"
LAST_RSS="0"

runTimedPhase() {
	local tool="$1"
	local toolVersion="$2"
	local iteration="$3"
	local phase="$4"
	local commandLabel="$5"
	local metricsFile
	local exitCode
	local realSeconds
	local userSeconds
	local systemSeconds
	local maxRSS
	local healthStatus=""

	metricsFile="$(mktemp -t nozarashi-benchmark.XXXXXX)"

	set +e
	/usr/bin/time -lp \
		sh -c 'exec "$@" 2>&3' sh \
		"$SCRIPT_DIR/benchmark.sh" "__phase" "$tool" "$phase" \
		3>&2 2>"$metricsFile"
	exitCode=$?
	set -e

	realSeconds="$(awk '$1 == "real" { print $2; exit }' "$metricsFile")"
	userSeconds="$(awk '$1 == "user" { print $2; exit }' "$metricsFile")"
	systemSeconds="$(awk '$1 == "sys" { print $2; exit }' "$metricsFile")"
	maxRSS="$(awk '/maximum resident set size/ { print $1; exit }' "$metricsFile")"
	rm -f "$metricsFile"

	realSeconds="${realSeconds:-0}"
	userSeconds="${userSeconds:-0}"
	systemSeconds="${systemSeconds:-0}"
	maxRSS="${maxRSS:-0}"

	if [[ "$phase" == "start" ]]; then
		if ((exitCode == 0)); then
			healthStatus="healthy"
		else
			healthStatus="failed"
		fi
	fi

	writeCSVRow \
		"$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
		"$HOST_OS" \
		"$HOST_ARCH" \
		"$tool" \
		"$toolVersion" \
		"$iteration" \
		"$phase" \
		"$commandLabel" \
		"$exitCode" \
		"$realSeconds" \
		"$userSeconds" \
		"$systemSeconds" \
		"$maxRSS" \
		"$HEALTH_URL" \
		"$healthStatus" \
		"preserved" \
		"preserved"

	LAST_REAL="$realSeconds"
	LAST_USER="$userSeconds"
	LAST_SYSTEM="$systemSeconds"
	LAST_RSS="$maxRSS"

	return "$exitCode"
}

benchmarkTool() {
	local tool="$1"
	local toolVersion="$2"
	local iteration="$3"
	local totalReal="0"
	local totalUser="0"
	local totalSystem="0"
	local maximumRSS="0"

	local phases=(
		"build"
		"start"
		"stop"
	)

	local labels=(
		"$tool build"
		"$tool start and wait for health"
		"$tool stop"
	)

	for phaseIndex in "${!phases[@]}"; do
		local phase="${phases[$phaseIndex]}"

		if ! runTimedPhase \
			"$tool" \
			"$toolVersion" \
			"$iteration" \
			"$phase" \
			"${labels[$phaseIndex]}"; then
			echo "$tool $phase failed during iteration $iteration" >&2
			return 1
		fi

		totalReal="$(awk -v total="$totalReal" -v value="$LAST_REAL" 'BEGIN { printf "%.2f", total + value }')"
		totalUser="$(awk -v total="$totalUser" -v value="$LAST_USER" 'BEGIN { printf "%.2f", total + value }')"
		totalSystem="$(awk -v total="$totalSystem" -v value="$LAST_SYSTEM" 'BEGIN { printf "%.2f", total + value }')"
		maximumRSS="$(awk -v maximum="$maximumRSS" -v value="$LAST_RSS" 'BEGIN { print (value > maximum ? value : maximum) }')"
	done

	writeCSVRow \
		"$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
		"$HOST_OS" \
		"$HOST_ARCH" \
		"$tool" \
		"$toolVersion" \
		"$iteration" \
		"total" \
		"$tool build + start + stop" \
		"0" \
		"$totalReal" \
		"$totalUser" \
		"$totalSystem" \
		"$maximumRSS" \
		"$HEALTH_URL" \
		"healthy" \
		"preserved" \
		"preserved"
}

for ((iteration = 1; iteration <= ITERATIONS; iteration++)); do
	echo "benchmark iteration $iteration/$ITERATIONS"

	# Alternate tool order to reduce cache and thermal ordering bias.
	if ((iteration % 2 == 1)); then
		benchmarkTool "nozarashi" "$NOZARASHI_VERSION" "$iteration"
		benchmarkTool "docker" "$DOCKER_VERSION" "$iteration"
	else
		benchmarkTool "docker" "$DOCKER_VERSION" "$iteration"
		benchmarkTool "nozarashi" "$NOZARASHI_VERSION" "$iteration"
	fi
done

echo "benchmark results written to $RESULT_FILE"
