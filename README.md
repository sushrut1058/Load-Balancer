# Reverse Proxy 

## Overview

This is a high-performance reverse proxy implemented in Go, designed to efficiently handle HTTP requests and forward them to backend servers. The project supports concurrent processing with worker pools and incorporates a custom queue to manage incoming requests.
It supports both Layer 7 (Application Layer) and Layer 4 (Transport Layer) load balancing with HTTP/HTTPS traffic and caching mechanisms, including in-memory caching and also Redis caching for reliability.

## Features
- **Concurrent Request Handling**: Uses a worker pool to process multiple requests simultaneously.
- **Request Queueing**: Employs Go's in- built thread-safe channels as queue to manage and distribute requests to available workers.
- **Worker Management**: Dynamic management of worker availability, with a mechanism to handle busy states.
- **Layer 7 Load Balancing**: Routes HTTP/HTTPS traffic to backend servers based on the URL path and other HTTP headers.
- **Layer 4 Load Balancing**: Routes traffic using only transport layer information making it a lower latency load balancing option  
- **Caching**: Supports basic in-memory caching and Redis caching to reduce backend load and improve response times (Note: for simplicity, caching is enabled only for GET requests due to the state changing nature of POST requests)
- **Round-Robin & Weighted Round-Robin Scheduling**: Distributes requests according to the backend servers' pre-defined capacity to handle load.
- **Health Checks**: WIP

## Getting Started

### Prerequisites

- Go (version 1.15 or later recommended)
- Redis (optional, for Redis caching functionality)

### Usage Instructions

1. **Clone the repository**

   ```sh
   git clone https://github.com/sushrut1058/Reverse-Proxy.git
   cd Reverse-Proxy/reverse-proxy
   ```

2. **Edit the configuration file**
    Edit the `config.json` according to your needs.
   ```sh
   {
    "port":"8080",
    "level":"L7",
    "strategy":"round-robin",
    "caching":"redis",
    "proto":"tcp",
    "serializer":"none",
    "servers":{
        "http://localhost:8081":3,
        "http://localhost:8082":2
    },
    "maxWorkers":5,
    "cache-ignore":[]
   }
   ```
   
   - `port`: Port on which the loadbalancer service runs
   - `level`: L4/L7 (OSI model layer)
   - `strategy`: Routing strategy, currently supports `round-robin` and `weighted-round-robin`
   - `caching`: Caching mechanism, supports `none`, `baseline` for basic in-memory caching and `redis` for redis caching
   - `servers`: Server URLs and their load handling capacities (number of requests they can handle at once safely)
   - `maxWorkers`: Number of worker threads that will pick requests from channel, it is recommended that the number worker threads be greater than or equal to the total load handling capacity of all backend servers
   - `cache-ignore`: The http methods for which requests will not be cached
   - `proto`: WIP
   - `serializer`: WIP

3. **Dummy Client and Servers**

   - Use as many servers as you need for testing. These dummy servers are just echo servers spitting out random text, the command for starting dummy servers is given below.
   ```sh
   go run main.go 8xxx
   ```
   - Similarly, the test client can be started in a similar manner. There is are two test clients, `concurrent_client` and `client` but it is recommended to use `concurrent_client` instead, due to its ability to emulate multiple clients which is ideal for testing a load-balancer
   ```sh
   go run main.go
   ```
    
Hit me up if you have any queries.
- Twitter: https://twitter.com/_localgh0st_
- Linkedin: https://www.linkedin.com/in/sushrut-hundikar-2371681b6/
