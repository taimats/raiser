package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	fp, ok := os.LookupEnv("FRONT_BASE_URL")
	if !ok {
		log.Fatal("パスのURLを設定ください")
	}
	bp, ok := os.LookupEnv("BACK_BASE_URL")
	if !ok {
		log.Fatal("パスのURLを設定ください")
	}
	ticker := time.NewTicker(5 * time.Minute)
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	fmt.Printf("\x1b[32m%s\x1b[0m\n", "ヘルスチェックを開始します...")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
L:
	for {
		select {
		case <-ticker.C:
			if err := healthCheck(fp); err != nil {
				l.Error("REQUEST_ERROR", slog.String("front_error", err.Error()))
			}
			if err := healthCheck(bp); err != nil {
				l.Error("REQUEST_ERROR", slog.String("back_error", err.Error()))
			}
			l.Info("REQUEST_OK")
		case <-ctx.Done():
			break L
		}
	}
	fmt.Printf("\x1b[32m%s\x1b[0m\n", "ヘルスチェックを終了しました!!")
}

func healthCheck(path string) error {
	res, err := http.Get(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	io.Copy(io.Discard, res.Body)
	return nil
}
