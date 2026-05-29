# Test Container for Snake Tournament Game Service

This directory contains the test container setup for the snake tournament game service.

## Overview

This test container provides an isolated environment for running tests against the game service. It includes:
- A Dockerfile for building the test environment
- Docker Compose configuration for test dependencies
- Test execution scripts

## Build Instructions

To build the test container:

1. Make sure you have Docker installed
2. Run: `docker build -t snake-tournament-test .`

To run tests in the container:

1. `docker run snake-tournament-test`

## Test Structure

The test container includes:
- All necessary dependencies for testing the game service
- A configured environment for running unit and integration tests
- Test helpers and utilities for consistent test execution

## Running Tests

Tests should be run via the container to ensure consistency with the development environment.
