package main

import (
   "context"
   "fmt"
   "github.com/gorilla/mux"
   "go.uber.org/zap"
   "google.golang.org/grpc"
   "google.golang.org/grpc/metadata"
   "net"
   "net/http"
   "zap_demo/http_server"
   proto "zap_demo/proto"
   "zap_demo/simple_zap"

   "github.com/google/uuid"
)

// http 拦截器, 给 context 注入一个带 traceId的 拦截器
func traceIdInterceptor(next http.Handler) http.Handler {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      trace_id := uuid.New().String()
      ctx := simple_zap.NewCtx(r.Context(), zap.String("traceId", trace_id)) // 给 ctx 注入一个 Field(内含日志打印的 k-v对, k 为 "traceId", v 就是 traceId)
      ctx = context.WithValue(ctx, "traceId", trace_id)                      // 给 ctx 注入 withValue k-v 对, 用户 grpc 调用时取出 traceId, 通过
      r2 := r.WithContext(ctx)
      next.ServeHTTP(w, r2)
   })
}

// grpc 拦截器, 尝试获取 traceId
func unaryTraceIdIterceptor() grpc.UnaryServerInterceptor {
   return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
      md, ok := metadata.FromIncomingContext(ctx)
      if ok {
         trace_id := md.Get("traceId")[0]
         ctx = simple_zap.NewCtx(ctx, zap.String("traceId", trace_id))
      } else {
         trace_id := uuid.New().String()
         ctx = simple_zap.NewCtx(ctx, zap.String("traceId", trace_id))
      }
      resp, err = handler(ctx, req)
      return resp, err
   }
}

func main() {

   //http server
   r := mux.NewRouter()
   r.Use(traceIdInterceptor)
   r.HandleFunc("/hello", http_server.HelloHandler)
   r.HandleFunc("/hello_to_grpc", http_server.HelloToGrpcHandler)

   http.Handle("/", r)

   http_err_chan := make(chan error)
   go func() {
      http_err_chan <- http.ListenAndServe(":2000", nil)
      fmt.Println("err http")
   }()

   // grpc server
   lis, err := net.Listen("tcp", ":8000")
   if err != nil {
      fmt.Printf("Fail to listen: %v", err)
   }
   fmt.Println("start listen---")
   s := proto.Server{}

   grpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryTraceIdIterceptor()))

   proto.RegisterChatServiceServer(grpcServer, &s)

   grpc_err_chan := make(chan error)
   go func() {
      grpc_err_chan <- grpcServer.Serve(lis)

   }()

   err = <-grpc_err_chan
   if err != nil {
      fmt.Println(err)
   }

}
