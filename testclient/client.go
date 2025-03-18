package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// サーバーのアドレスとポートを指定
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("サーバーへの接続に失敗:", err)
	}
	defer conn.Close()

	// サーバーからのメッセージを受け取るゴルーチン
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Println("サーバーからの読み込みエラー:", err)
		}
	}()

	// 標準入力から読み取り、サーバーへ送信
	stdinScanner := bufio.NewScanner(os.Stdin)
	for stdinScanner.Scan() {
		text := stdinScanner.Text()
		if text == "/aaa" {
			// 例: 2バイトのビット列を定義して string に変換して送信する
			bitData := []byte{0x12, 0x23}
			_, err = fmt.Fprintln(conn, string(bitData))
		} else {
			_, err = fmt.Fprintln(conn, text)
		}
		if err != nil {
			// サーバーが接続を閉じた場合はエラーを無視してループを抜ける
			if opErr, ok := err.(*net.OpError); ok {
				if opErr.Err.Error() == "wsasend: An established connection was aborted by the software in your host machine" {
					break
				}
			}
			log.Println("送信エラー:", err)
			break
		}
	}
	if err := stdinScanner.Err(); err != nil {
		log.Println("標準入力の読み込みエラー:", err)
	}
}
