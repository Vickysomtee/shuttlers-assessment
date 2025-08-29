# Shuttlers Assesment
This document provides detailed documentation on the project.

## Deployments
The deploy.sh file handles the entire deployment process from building the image to running the application container.

### CI/CD integration
Integrating the deploy.sh script into a CI/CD pipeline will follow the following steps

- Ensure the pipeline points to the deployment branch
- ssh into the server
- git clone the repository into the server
- cd into project and run the deploy.sh script.

## Monitoring && Observality

### Logs

By default, go applications sends logs to stdout. Two tools will be used for aggregating and displaying the logs.

- Loki
  
    Loki is a log aggregation tool which can be used to aggregate logs directly from docker container. Loki has a driver that can be installed as docker plugin. This plugin collects the logs and pushes to a running Loki instance
    
    To install the loki docker plugin
    ```sh
    docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions
    ```

    To push the logs to a running Loki instance
    ```sh
    sudo nano /etc/docker/daemon.json
    ```

    Paste the below content in the file, and change localhost to the host address where the loki instance is running.
    ```json
    {
      "log-driver": "loki",
      "log-opts": {
        "loki-url": "http://localhost:3100/loki/api/v1/push",
        "loki-batch-size": "400"
      }
    }
    ```

- Grafana
  
    On the Grafana instance, add loki as a data source connection, aggrgate the container logs, and add to a Grafana Dashboard.

### Metrics
Cadvisor is a tool for exposing all container metric from the host service. Cadvisor integrates with prometheus which scrapes the metric from cadvisor which can then be visualized on Grafana.

cadvisor needs to run on the host machine where the application container is running. Once the cadvisor instance is running, it automatically collects and exposes the container metrics. Prometheus can then be configured to scrape metrics from cadvisor using the below prometheus config

```yaml
global:
  scrape_interval: 15s
  scrape_timeout: 10s
  evaluation_interval: 15s
alerting:
  alertmanagers:
    - static_configs:
      - targets: []
      scheme: http
      timeout: 10s
      api_version: v1
scrape_configs:
  - job_name: prometheus
    honor_timestamps: true
    scrape_interval: 15s
    scrape_timeout: 10s
    metrics_path: /metrics
    scheme: http
    static_configs:
      - targets:
        - localhost:9090
  - job_name: cadvisor
    scrape_interval: 5s
    static_configs:
      - tragets:
        - <localhost:8083> #IP:PORT of the host runningcadvisor
```


## Basic Auditing
This Audit was carried out using the [docker bench security tool](https://github.com/docker/docker-bench-security). Below is a script to carry out the test.

```
git clone https://github.com/docker/docker-bench-security.git
cd docker-bench-security
sudo ./docker-bench-security.sh
```
Below is a screenshot of the audit result

<img width="1347" height="836" alt="Screenshot From 2025-08-29 01-40-06" src="https://github.com/user-attachments/assets/18e42800-6778-4d4a-a40f-b1193491eeab" />


### Potential Vulnerability 1: Running the Docker daemon as root

**Issue:**
 - By default, the Docker daemon (dockerd) runs as the root user, and containers inherit root privileges inside their namespaces.

 - If a container escapes its isolation (via kernel exploit, misconfiguration, etc.), it could gain root access to the host system.

 - Mapping the host’s Docker socket (/var/run/docker.sock) into a container effectively gives that container root on the host.

**Solution:**

 - Use Rootless Docker (dockerd-rootless.sh) to run the Docker daemon as a non-root user.

 - Avoid mounting /var/run/docker.sock into containers unless absolutely necessary.

### Potential Vulnerability 1: Running the docker image as root

**Issue:**
 - When no user is creating at the process of creating the image, the container also inherits root priviledges.

**Solution:**
 - Always create a user for the container in the docker file.
   


## Service Level Objective (SLO) for "Hello, World!" Application

### SLO Definition

**99.9% availability over a rolling 30-day period**
This means the "Hello, World!" application should be accessible and responding successfully to requests 99.9% of the time in any given 30-day window, allowing for approximately 43 minutes of downtime per month.

### Service Level Indicators (SLIs)
**1. Request Success Rate**

**Definition:** Percentage of HTTP requests that return a successful response (2xx status codes) out of total requests

**Target:** ≥ 99.9% of requests should return HTTP 200 status
Measurement Window: Rolling 5-minute intervals, aggregated over 30 days

**2. Response Time (Latency)**

**Definition:** 95th percentile response time for successful requests

**Target:** ≥ 99.9% of successful requests should complete within 500ms at the 95th percentile
**Measurement Window:** Rolling 5-minute intervals, aggregated over 30 days

### Tracking Implementation
**Data Collection**
**Metrics to Capture:**

- HTTP response status codes (200, 4xx, 5xx)
- Request timestamps (start and end)
- Response times in milliseconds
- Request URLs and methods
- Error details for failed requests

**Monitoring Infrastructure**

- Application-level logging: Instrument the web server to log each request with timestamp, status code, and response time
- Request metrics: Collect health check results and request forwarding statistics
- Infrastructure monitoring: Track server CPU, memory, and network connectivity

### SLI Calculation Methods
**Success Rate SLI**

```
Success Rate = (Count of 2xx responses / Total request count) × 100
```

 - Measured every 5 minutes
 - Aggregated using a rolling 30-day window
 - Alert if success rate drops below 99.9% for more than 15 minutes

**Latency SLI:**

```
95th Percentile Latency = 95th percentile of response times for successful requests
```

 - Calculated every 5 minutes from response time histogram
 - Track trend over 30-day rolling window
 - Alert if p95 latency exceeds 500ms consistently

**Alerting and Reporting**

 - Real-time alerts: Trigger notifications when SLI thresholds are breached
 - Daily reports: Summary of SLO compliance and trending
 - Monthly reviews: Comprehensive analysis of SLO performance and potential improvements



  



