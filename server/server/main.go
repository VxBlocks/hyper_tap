package main

import (
	"config"
	"context"
	"hyperliquid-server/firebase"
	"hyperliquid-server/gen/login/v1/loginv1connect"
	"hyperliquid-server/gen/news/v1/newsv1connect"
	"hyperliquid-server/gen/price_alert/v1/pricealertv1connect"
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

	// err = firebase.SendMessageTopic(context.Background(), "news", firebase.GetFirebaseApp())
	// if err != nil {
	// 	panic(err)
	// }
	// err = firebase.SendMessage(context.Background(), "cw7EFJ8MRCi46aY89SbbVu:APA91bEH9e9mZhMBNZ6a2l1mpFs7d0sV4o5Toe1aoP_PRVOqpOm3-_aRPna7c3f_2J8mOD0hbKDkW3mhe_fXp7c7gI23BGAfYzt0WayRLIdMy_FbKrHmSWo", firebase.GetFirebaseApp())
	// if err != nil {
	// 	panic(err)
	// }

	// return
	// monitor.StartMonitor()
	if config.C.Postgres.Url != "10.10.2.11" {
		monitor.StartNewsMonitor()
		monitor.StartPriceMonitor()
	}

	err = migrators.Migrate()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	loggingInterceptor := connect.WithInterceptors(NewLoggingInterceptor())

	reflector := grpcreflect.NewStaticReflector(
		userwatchv1connect.UserWatchServiceName,
		loginv1connect.LoginServiceName,
		newsv1connect.NewsServiceName,
		pricealertv1connect.PriceAlertServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	user_watch := handler.UserWatchHandler{}
	mux.Handle(userwatchv1connect.NewUserWatchServiceHandler(&user_watch, loggingInterceptor))

	login := handler.LoginHandler{}
	mux.Handle(loginv1connect.NewLoginServiceHandler(&login, loggingInterceptor))

	news := handler.NewsHandler{}
	mux.Handle(newsv1connect.NewNewsServiceHandler(&news, loggingInterceptor))

	priceAlert := handler.PriceAlertHandler{}
	mux.Handle(pricealertv1connect.NewPriceAlertServiceHandler(&priceAlert, loggingInterceptor))

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
		w.Header().Set("Content-Type", "application/text; charset=utf-8")
	})

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
