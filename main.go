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
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("ポート番号を設定ください")
	}
	srv := startServer(port)
	fmt.Printf("\x1b[32m%s\x1b[0m\n", "サーバーが起動しました!!")

	path, ok := os.LookupEnv("FRONT_BASE_URL")
	if !ok {
		log.Fatal("パスのURLを設定ください")
	}
	ticker := time.NewTicker(10 * time.Second)
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	fmt.Printf("\x1b[32m%s\x1b[0m\n", "ヘルスチェックを開始します...")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
L:
	for {
		select {
		case <-ticker.C:
			if err := healthCheck(path); err != nil {
				l.Error("REQUEST_ERROR", slog.String("err", err.Error()))
			}
			l.Info("REQUEST_OK")
		case <-ctx.Done():
			break L
		}
	}
	<-ctx.Done()
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("サーバーのシャットダウンに失敗: (error: %s)", err.Error())
	}
	fmt.Printf("\x1b[32m%s\x1b[0m\n", "サーバーがシャットダウンしました!!")
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

func startServer(port string) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%s", port)}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("サーバーの起動に失敗:%v", err)
		}
	}()
	return srv
}
