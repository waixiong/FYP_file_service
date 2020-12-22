package middleware

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"google.golang.org/grpc"
)

// https://awesomeopensource.com/project/grpc-ecosystem/go-grpc-middleware

// AddInterceptor returns grpc.Server config option.
func AddInterceptor(unaries []grpc.UnaryServerInterceptor, streams []grpc.StreamServerInterceptor, opts []grpc.ServerOption) []grpc.ServerOption {
	// Add unary interceptor
	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		unaries...,
	))

	// Add stream interceptor (added as an example here)
	opts = append(opts, grpc_middleware.WithStreamServerChain(
		streams...,
	))

	return opts
}
