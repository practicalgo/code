// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RepoClient is the client API for Repo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RepoClient interface {
	GetRepos(ctx context.Context, in *RepoGetRequest, opts ...grpc.CallOption) (Repo_GetReposClient, error)
}

type repoClient struct {
	cc grpc.ClientConnInterface
}

func NewRepoClient(cc grpc.ClientConnInterface) RepoClient {
	return &repoClient{cc}
}

func (c *repoClient) GetRepos(ctx context.Context, in *RepoGetRequest, opts ...grpc.CallOption) (Repo_GetReposClient, error) {
	stream, err := c.cc.NewStream(ctx, &Repo_ServiceDesc.Streams[0], "/Repo/GetRepos", opts...)
	if err != nil {
		return nil, err
	}
	x := &repoGetReposClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Repo_GetReposClient interface {
	Recv() (*RepoGetReply, error)
	grpc.ClientStream
}

type repoGetReposClient struct {
	grpc.ClientStream
}

func (x *repoGetReposClient) Recv() (*RepoGetReply, error) {
	m := new(RepoGetReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RepoServer is the server API for Repo service.
// All implementations must embed UnimplementedRepoServer
// for forward compatibility
type RepoServer interface {
	GetRepos(*RepoGetRequest, Repo_GetReposServer) error
	mustEmbedUnimplementedRepoServer()
}

// UnimplementedRepoServer must be embedded to have forward compatible implementations.
type UnimplementedRepoServer struct {
}

func (UnimplementedRepoServer) GetRepos(*RepoGetRequest, Repo_GetReposServer) error {
	return status.Errorf(codes.Unimplemented, "method GetRepos not implemented")
}
func (UnimplementedRepoServer) mustEmbedUnimplementedRepoServer() {}

// UnsafeRepoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RepoServer will
// result in compilation errors.
type UnsafeRepoServer interface {
	mustEmbedUnimplementedRepoServer()
}

func RegisterRepoServer(s grpc.ServiceRegistrar, srv RepoServer) {
	s.RegisterService(&Repo_ServiceDesc, srv)
}

func _Repo_GetRepos_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RepoGetRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RepoServer).GetRepos(m, &repoGetReposServer{stream})
}

type Repo_GetReposServer interface {
	Send(*RepoGetReply) error
	grpc.ServerStream
}

type repoGetReposServer struct {
	grpc.ServerStream
}

func (x *repoGetReposServer) Send(m *RepoGetReply) error {
	return x.ServerStream.SendMsg(m)
}

// Repo_ServiceDesc is the grpc.ServiceDesc for Repo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Repo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Repo",
	HandlerType: (*RepoServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetRepos",
			Handler:       _Repo_GetRepos_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "repositories.proto",
}
