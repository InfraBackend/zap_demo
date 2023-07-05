package __

import (
  "context"
  "zap_demo/simple_zap"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, in *ChatMessage) (*ChatMessage, error) {
  simple_zap.WithCtx(ctx).Sugar().Info("Reach SayHello service!")
  return &ChatMessage{Body: "Hello From SayHello service!"}, nil
}
