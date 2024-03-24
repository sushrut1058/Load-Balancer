# Custom Loadbalancer

## Overview

This Load Balancer project is designed to distribute incoming network traffic across multiple backend servers to enhance the reliability, scalability, and performance of web applications, microservices, or API endpoints. It supports both Layer 7 (Application Layer) load balancing with HTTP/HTTPS traffic and basic caching strategies, including an in-memory cache and Redis cache for optimized response times.

## Features

- **Layer 7 Load Balancing**: Routes HTTP/HTTPS traffic to backend servers based on the URL path and other HTTP headers.
- **Caching**: Supports basic in-memory caching and Redis caching to reduce backend load and improve response times.
- **Round-Robin Scheduling**: Distributes requests evenly across the pool of backend servers.
- **Health Checks**: WIP
- **Flexible Backend Management**: WIP

## Getting Started

### Prerequisites

- Go (version 1.15 or later recommended)
- Redis (optional, for Redis caching functionality)

### Usage Instructions

1. **Clone the repository**

   ```sh
   git clone https://github.com/yourusername/loadbalancer-project.git
   cd loadbalancer-project
   ```

2. **Edit the configuration file**
    Edit the `conf.json` according to your needs, here is a snapshot of the same 
   ```sh
   git clone https://github.com/yourusername/loadbalancer-project.git
   cd loadbalancer-project
   ```
   ![image](https://github.com/sushrut1058/loadbalancer/assets/62463384/a07d393a-38c8-4ce0-a38e-df89417e39da)
   - `port`: Port on which the loadbalancer service runs
   - `level`: L4/L7 (OSI model layer)
   - `strategy`: Routing strategy, currently only supports `round-robin`
   - `caching`: Caching mechanism, supports `none`, `baseline` for basic in-memory caching and `redis` for redis caching
   - `proto`: WIP
   - `serializer`: WIP
   - `servers`: Server URLs

3. **Dummy Client and Servers**

   - Use as many dummy servers as you need for testing. These dummy servers are just echo servers spitting out loremipsum, the command for starting dummy servers is:
   ```sh
   go run main.go 8xxx
   ```
   - Similarly, the dummy client can be started like this:
   ```sh
   go run main.go
   ```
Hit me up if you have any queries.
- Twitter: https://twitter.com/_localgh0st_
- Linkedin: https://www.linkedin.com/in/sushrut-hundikar-2371681b6/
