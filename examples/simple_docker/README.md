# Simple Docker

Dockerfile backend with PostgreSQL and Redis.

## Benchmark

Create the Nozarashi resources once, then run three benchmark iterations:

```bash
sudo ../../bin/nozarashi create
./benchmark.sh
```

Pass an iteration count and output path when needed:

```bash
./benchmark.sh 5 ./result.csv
```

The benchmark preserves build caches and data volumes. It records build, healthy
startup, stop and total measurements for both tools in `result.csv`.
