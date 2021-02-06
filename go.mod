module github.com/metaprov/modeld-go-sdk

go 1.13


replace (
	google.golang.org/grpc v1.32.0 => google.golang.org/grpc v1.29.1
	github.com/golang/protobuf v1.4.2 => github.com/golang/protobuf v1.3.4

)
require (
	github.com/golang/protobuf v1.4.2
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/pkg/errors v0.9.1
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20201009025420-dfb3f7c4e634 // indirect
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0

)
