package server

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func slogGinMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 次のハンドラー（実際のルーティング処理）を実行
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()

		// エラーがあれば取得
		errs := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// メソッドごとに独立したフィールドとして構造化ログを出力
		// HTTPメソッド、ステータス、レスポンスタイムなどの状態を外部からパッと見で把握できる
		logger.Info(
			"HTTP Request",
			slog.Int("status", status),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", c.ClientIP()),
			slog.Duration("latency", latency),
			slog.String("user-agent", c.Request.UserAgent()),
			slog.String("errors", errs),
		)
	}
}
