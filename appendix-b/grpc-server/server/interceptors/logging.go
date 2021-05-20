package interceptors

import (
	"context"

	"github.com/practicalgo/code/appendix-b/grpc-server/server/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type wrappedServerStream struct {
	config *config.AppConfig
	grpc.ServerStream
}

func (s wrappedServerStream) SendMsg(m interface{}) error {
	s.config.Logger.Printf("Send msg called: %T", m)
	return s.ServerStream.SendMsg(m)
}

func (s wrappedServerStream) RecvMsg(m interface{}) error {
	s.config.Logger.Printf("Send msg called: %T", m)
	return s.ServerStream.RecvMsg(m)
}

func LoggingUnaryInterceptor(config *config.AppConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		l := config.Logger.Info().Str("method", info.FullMethod).Err(err)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			config.Logger.Debug().Msg("no metadata")
		} else {
			l = l.Strs("request-id", md.Get("Request-Id"))
		}
		l.Send()
		return resp, err
	}
}

func LoggingStreamInterceptor(config *config.AppConfig) grpc.StreamServerInterceptor {
	return func(srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		serverStream := wrappedServerStream{
			config:       config,
			ServerStream: stream,
		}
		err := handler(srv, serverStream)
		ctx := stream.Context()
		l := config.Logger.Info().Str("method", info.FullMethod).Err(err)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			config.Logger.Debug().Msg("no metadata")
		} else {
			l = l.Strs("request-id", md.Get("Request-Id"))
		}
		l.Send()
		return err
	}
}
