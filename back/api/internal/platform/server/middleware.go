package server

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"svem-chat-api/internal/platform/config"

	"github.com/gin-gonic/gin"
)

func Start() {
	logFile, err := os.OpenFile("hemochart.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		panic("ログファイルのオープンに失敗しました: " + err.Error())
	}

	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	logger := slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	router := gin.New()
	router.Use(slogGinMiddleware(logger))
	router.Use(gin.Recovery())

	conf := config.GetConfig()
	srv := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: router,
	}

	routing(router)

	go func() {
		slog.Info("サーバーを起動します", slog.String("port", conf.Port))
		// ErrServerClosed は「正常なシャットダウン」なのでエラーとして扱わない
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("サーバーの異常終了", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	// SIGINT (Ctrl+C) と SIGTERM (DockerやKubernetesからの終了命令) を監視
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("終了シグナルを受信しました。サーバーを安全に停止します...")

	// この間に、現在処理中のリクエストを最後まで完了させる
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("サーバーの強制終了", slog.String("error", err.Error()))
	}

	slog.Info("サーバーが正常に停止しました")
}
