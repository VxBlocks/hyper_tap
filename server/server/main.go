package main

import (
	"context"
	"hyperliquid-server/firebase"
	"hyperliquid-server/gen/login/v1/loginv1connect"
	"hyperliquid-server/gen/userwatch/v1/userwatchv1connect"
	"hyperliquid-server/handler"
	"hyperliquid-server/migrators"
	"hyperliquid-server/monitor"
	"logger"
	"net/http"
	"os"
	"time"

	"connectrpc.com/grpcreflect"

	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	_, err := firebase.InitFirebaseApp(context.Background())
	if err != nil {
		panic(err)
	}

	// panic(firebase.SendMessageTopic(context.Background(), "watch_address", firebase.GetFirebaseApp()))

	monitor.StartMonitor()

	err = migrators.Migrate()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	loggingInterceptor := connect.WithInterceptors(NewLoggingInterceptor())

	reflector := grpcreflect.NewStaticReflector(
		userwatchv1connect.UserWatchServiceName,
		loginv1connect.LoginServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	user_watch := handler.UserWatchHandler{}
	mux.Handle(userwatchv1connect.NewUserWatchServiceHandler(&user_watch, loggingInterceptor))

	login := handler.LoginHandler{}
	mux.Handle(loginv1connect.NewLoginServiceHandler(&login, loggingInterceptor))

	var addr = "0.0.0.0:" + os.Getenv("API_PORT")
	logger.Info("starting server", "listen", addr)

	err = http.ListenAndServe(
		addr,
		// Use h2c so we can serve HTTP/2 without TLS.
		withCORS(h2c.NewHandler(mux, &http2.Server{})),
	)
	if err != nil {
		panic(err)
	}
}

func NewLoggingInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			now := time.Now()
			res, err := next(ctx, req)
			if err != nil {
				logger.Error("rpc_error", "method", req.Spec().Procedure, "error", err)
			} else {
				logger.Info("rpc_call", "method", req.Spec().Procedure, "duration", time.Since(now))
			}
			return res, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func NewAuthInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Header().Get("Authorization") != "" {

			}
			now := time.Now()
			res, err := next(ctx, req)
			if err != nil {
				logger.Error("rpc_error", "service", req.Spec().Procedure, "error", err)
			} else {
				logger.Info("rpc_call", "service", req.Spec().Procedure, "duration", time.Since(now))
			}
			return res, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

// withCORS adds CORS support to a Connect HTTP handler.
func withCORS(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
		// AllowedOrigins: []string{"example.com"},
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: append(connectcors.AllowedHeaders(), "Authorization"),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return middleware.Handler(h)
}
