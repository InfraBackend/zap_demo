package http_server

import (
  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"
  "net/http"
  proto "zap_demo/proto"
  "zap_demo/simple_zap"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
  simple_zap.WithCtx(r.Context()).Sugar().Info("Reach HelloHandler.")
}

func HelloToGrpcHandler(w http.ResponseWriter, r *http.Request) {
  simple_zap.WithCtx(r.Context()).Sugar().Info("Reach HelloToGrpcHandler.")

  var conn *grpc.ClientConn
  conn, err := grpc.Dial(":8000", grpc.WithInsecure())
  if err != nil {
     simple_zap.WithCtx(r.Context()).Sugar().Error(err)
  }

  defer conn.Close()

  // 获得一个 ChatService 的 client
  c := proto.NewChatServiceClient(conn)

  // grpc 调用远程的 SayHello
  trace_id := r.Context().Value("traceId").(string)
  grpc_ctx := metadata.NewOutgoingContext(r.Context(), metadata.Pairs("traceId", trace_id))
  response, err := c.SayHello(grpc_ctx, &proto.ChatMessage{Body: "Hello From Client!"})
  if err != nil {
     simple_zap.WithCtx(r.Context()).Sugar().Warn(err)
  }

  simple_zap.WithCtx(r.Context()).Sugar().Debugf("Response from server: %s", response.Body)

}
