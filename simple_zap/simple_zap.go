package simple_zap

import (
   "context"
   "go.uber.org/zap"
   "go.uber.org/zap/zapcore"
   "gopkg.in/natefinch/lumberjack.v2"
)

const loggerCtxKey = "logger_ctx_key"

var Logger *zap.Logger

func init() {
   hook := lumberjack.Logger{
      Filename:   "/zap_demo/log/simple_zap.log",
      MaxSize:    1,    // 日志文件最大体积(MB)
      MaxAge:     1,    //最多维持几天的日志
      MaxBackups: 1,    //最多保持多少个日志文件
      LocalTime:  true, // 使用本地时间, 而不是UTC时间
      Compress:   true, // 压缩日志
   }
   encoder_config := zapcore.EncoderConfig{
      MessageKey:     "msg",                          // 日志内容 的 key
      LevelKey:       "level",                        // 日志 level 的 key
      TimeKey:        "time",                         // 日志时间 的 key
      CallerKey:      "lineNum",                      // 日志产生的文件及其行数 的 key
      FunctionKey:    "func",                         // 日志产生的函数 的 key
      LineEnding:     zapcore.DefaultLineEnding,      // 回车符换行
      EncodeLevel:    zapcore.LowercaseLevelEncoder,  // level小写: info,debug,warn等 而不是 Info, Debug,Warn等
      EncodeTime:     zapcore.ISO8601TimeEncoder,     // 时间格式: "2006-01-02T15:04:05.000Z0700"
      EncodeDuration: zapcore.SecondsDurationEncoder, // 时间戳用float64型,更加准确, 另一种是NanosDurationEncoder int64
      EncodeCaller:   zapcore.ShortCallerEncoder,     // 产生日志文件的路径格式: 包名/文件名:行号

   }
   caller := zap.AddCaller() //日志打印输出 文件名, 行号, 函数名

   development := zap.Development()  // 可输出 dpanic, panic 级别的日志
 

   field := zap.Fields() // 负责给日志生成一个个 k-v 对

   var syncers []zapcore.WriteSyncer // io writer
   syncers = append(syncers, zapcore.AddSync(&hook))

   atomic_level := zap.NewAtomicLevel()  // 设置日志 level
   atomic_level.SetLevel(zap.DebugLevel) // 打印 debug, info, warn,eror, depanic,panic,fetal 全部级别日志

   core := zapcore.NewCore(zapcore.NewJSONEncoder(encoder_config), zapcore.NewMultiWriteSyncer(syncers...), atomic_level)

   Logger = zap.New(core, caller, development, field)

}

 
// 给 ctx 注入一个 looger, logger 中包含Field(内含日志打印的 k-v对)
func NewCtx(ctx context.Context, fields ...zapcore.Field) context.Context {
   return context.WithValue(ctx, loggerCtxKey, WithCtx(ctx).With(fields...))
}

//  尝试从 context 中获取带有 traceId Field的 logge
func WithCtx(ctx context.Context) *zap.Logger {
   if ctx == nil {
      return Logger
   }
   ctx_logger, ok := ctx.Value(loggerCtxKey).(*zap.Logger)
   if ok {
      return ctx_logger
   }
   return Logger
}
