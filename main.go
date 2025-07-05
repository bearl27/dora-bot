package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	// .env ファイルから環境変数をロード
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set.")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating Gemini client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	// 添削したい文章の入力
	fmt.Print("◯えもん: どうしたんだい？の▫太くん\n")
	scanner := NewScanner(os.Stdin)

	for {
		input, err := scanner.ReadLines()
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		textToEdit := strings.Join(input, "\n")

		if strings.ToLower(strings.TrimSpace(textToEdit)) == "exit" {
			fmt.Println("終了します。")
			break
		}

		if strings.TrimSpace(textToEdit) == "" {
			fmt.Println("◯えもん: なんだって？の▫太くん")
			continue
		}

		fmt.Println("\nテレテテッテレーン!")

		prompt := fmt.Sprintf(`文章の悩みを解決するようなドラえもんの秘密道具を考えてください。
		返信はまずに、以下のフォーマットでお願いします。
		「秘密道具名 !!!」

		ドラえもんっぽい口調で、秘密道具の機能の説明だけをしてください。

		文章:
		%s

		添削結果:
		`, textToEdit)

		// Gemini APIへのリクエスト
		resp, err := model.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			log.Printf("Error generating content: %v", err)
			fmt.Println("◯えもん: うーん、何か問題が発生したみたいだよ。もう一度試してみてね。")
			continue
		}

		// 添削結果の表示
		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			for _, part := range resp.Candidates[0].Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					fmt.Println(string(txt))
				}
			}
		} else {
			fmt.Println("◯えもん: うーん、何も返ってこなかったみたいだよ。もう一度試してみてね。")
		}
		fmt.Println("◯えもん: 何か他に悩みはないかい？の▫太くん")
		continue
	}
}

// Scanner struct to read multiple lines until an empty line
type Scanner struct {
	scanner *bufio.Scanner
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{scanner: bufio.NewScanner(r)}
}

func (s *Scanner) ReadLines() ([]string, error) {
	var lines []string
	for s.scanner.Scan() {
		line := s.scanner.Text()
		if line == "" { // 空行で入力終了
			break
		}
		lines = append(lines, line)
	}
	if err := s.scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
