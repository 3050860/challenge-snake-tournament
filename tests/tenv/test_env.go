package tenv

import "github.com/testcontainers/testcontainers-go"

type TestEnvironment struct {
	Net    *testcontainers.DockerNetwork
	Config Config
}

type testEnvironmentRequest struct {
	db bool
}

type TestEnvironmentOption func(req *testEnvironmentRequest)

func (opt TestEnvironmentOption) Customize(req *testEnvironmentRequest) {
	opt(req)
}

func WithDB() TestEnvironmentOption {
	return func(req *testEnvironmentRequest) {
		req.db = true
	}
}
