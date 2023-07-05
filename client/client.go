package main

import (
   "context"
   "net/http"
   "zap_demo/simple_zap"
)

func main() {
   _, err := http.Get("http://localhost:2000/hello")
   if err != nil {
      simple_zap.WithCtx(context.Background()).Sugar().Warn(err)
      return
   }

   _, err = http.Get("http://localhost:2000/hello_to_grpc")
   if err != nil {
      simple_zap.WithCtx(context.Background()).Sugar().Warn(err)
      return
   }

}
 
