# URL Shortener Project
This project is a lightweight URL shortener service, allowing users to convert long URLs into shorter, manageable links for easy sharing. Built with Go, it leverages the chi router for HTTP routing and the go-redis/redis/v8 package for storing URL mappings.

## Features
- Generate short URLs from long ones
- Redirect short URLs to their original long URLs

## Prerequisites
- Go (version 1.22)
- Docker and Docker Compose

## Usage
- To shorten a URL: http://localhost:8080/shorten?url=YOUR_LONG_URL
- To access a short URL: http://localhost:8080/r/YOUR_SHORT_URL
