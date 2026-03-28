# k6 Load Test

This directory contains a k6 load test script (`script.js`) specifically tailored for testing the `GetAllCourse` endpoint which utilizes Redis caching.

## Prerequisites
- [Docker](https://docs.docker.com/get-docker/) must be installed on your system.
**(Alternatively, if you have [k6](https://k6.io/docs/get-started/installation/) installed locally, you can run the manual command below).**

## How to Run

We use k6's built-in web dashboard (available in k6 v0.39.0+) to generate reports.

### Using Batch Script via Docker (Windows)
```cmd
.\run.bat
```
*(This will pull the `grafana/k6` image and run the script inside a Docker container. It uses `host.docker.internal` to connect to your local backend.)*

### Manual Command (Local k6)
```bash
# If you don't use the bat file and have k6 installed
k6 run --out web-dashboard script.js
```

### Authentication
If the `GetAllCourse` endpoint requires authorization, you can set `TOKEN` in `run.bat` or pass it as an environment variable: **run on CMD terminal
```bash
docker run --rm -i -p 5665:5665 -v %cd%:/home/k6 -e BASE_URL=http://host.docker.internal:8084 -e TOKEN=YOUR_JWT_TOKEN grafana/k6 run --out web-dashboard - < script.js
```

## Viewing the Dashboard
Once the test starts, the terminal will display a URL (usually `http://localhost:5665`). Open that URL in your browser to view the real-time dashboard report.
