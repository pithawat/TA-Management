@echo off
set BASE_URL=http://host.docker.internal:8084
set TOKEN=

echo Running k6 load test with web dashboard via Docker...
echo Note: If your endpoint requires a token, set it in this script.

docker run --rm -i -p 5665:5665 -e BASE_URL=%BASE_URL% -e TOKEN=%TOKEN% grafana/k6 run --out web-dashboard - < script.js
pause
