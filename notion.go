package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CreatePageRequest struct {
	Parent struct {
		PageID string `json:"page_id"`
	} `json:"parent"`
	Markdown string `json:"markdown"`
}

func main() {
	// 1. 環境変数を読む
	token := os.Getenv("NOTION_TOKEN")
	parentPageID := os.Getenv("NOTION_PARENT_PAGE_ID")

	if token == "" || parentPageID == "" {
		fmt.Println("NOTION_TOKEN または NOTION_PARENT_PAGE_ID が未設定です")
		return
	}

	// 2. Notionに送りたい内容を作る
	var reqBody CreatePageRequest
	reqBody.Parent.PageID = parentPageID
	reqBody.Markdown = `
	# Goから作成したページ
	
	<table>
	<tr><td>名前</td><td>年齢</td></tr>
	<tr><td>田中</td><td>20</td></tr>
	<tr><td>佐藤</td><td>21</td></tr>
	</table>
	`
	

	// 3. JSONに変換する
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("JSON変換エラー:", err)
		return
	}

	// 4. HTTPリクエストを作る
	req, err := http.NewRequest("POST", "https://api.notion.com/v1/pages", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("リクエスト作成エラー:", err)
		return
	}

	// 5. ヘッダーを付ける
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Notion-Version", "2026-03-11")
	req.Header.Set("Content-Type", "application/json")

	// 6. リクエストを送る
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("送信エラー:", err)
		return
	}
	defer resp.Body.Close()

	// 7. 結果を読む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("レスポンス読取エラー:", err)
		return
	}

	fmt.Println("status:", resp.Status)
	fmt.Println("response:")
	fmt.Println(string(body))
}