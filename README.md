# ドラ ◯ もんと会話してみた（Go・Gemini API・ハンズオン）

# はじめに

もし、あの青くて丸い友達、ドラえもんがのび太くんの悩みを解決してくれるように、あなたの悩みを解決してくれるとしたら？そんな夢のような体験を、Google の最新 AI モデル Gemini 1.5 Flash と Go 言語を使って実現してみました！

# 事前準備

今回は以下のものを使いますが，説明は行いません．事前に準備しておいてください．

- Go のダウンロード
- エディタ
- API キーの説明

# 1. GeminiAPI の取得

[Gemini API - Google AI for Developers](https://ai.google.dev/gemini-api/docs?hl=ja) を開き，以下の表の動作順に進み API キーを取得します．

|     |                                                         web の画像                                                         |                      動作                       |
| :-: | :------------------------------------------------------------------------------------------------------------------------: | :---------------------------------------------: |
|  1  | ![image.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3748983/23fa5f59-629b-41c3-9ded-394e19397750.png) |       「GeminiAPI キーを取得する」を押す        |
|  2  | ![image.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3748983/e6c36f54-5401-4c85-97f4-61456bf8a962.png) |            「API キーを作成」を押す             |
|  3  | ![image.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3748983/27305215-d531-4e67-bc1f-da4dfa45b912.png) | 　「新しいプロジェクトで API キーを作成」を押す |
|  4  | ![image.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3748983/558ed908-dc8e-4dd0-9033-3b590288c8b1.png) |       作成された API キーをコピーしておく       |

# 2. 開発

最終的には以下の tree のようにします．

```
dora-bot
├── .env
├── go.mod
├── go.sum
└── main.go
```

## 2-1. フォルダ作成

まず，フォルダを作成します．

```bash
mkdir dora-bot  ## フォルダ作成
cd dora-bot     ## フォルダ配下に入る
```

## 2-2. 環境変数設定

先ほど取得した API キーを環境変数として置いておきます

```bash
touch .env
```

この作成した.env 内に以下を記述して API キーを保存します

```
GEMINI_API_KEY=YOUR_API_KEY
```

> YOUR_API_KEY には自分の API キーをペーストしてください

## 2-3.go の準備

```
go mod init github.com/bearl27/dora-bot
go get github.com/google/generative-ai-go/genai
go get github.com/joho/godotenv
```

# 3. 開発

## 3-1. main.go の作成

以下の内容で`main.go`ファイルを作成します。

```go
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

	// ドラえもんの挨拶
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
```

## 3-2. プログラムの解説

### 主要な機能

1. **環境変数の読み込み**: `.env`ファイルから API キーを安全に読み込みます
2. **Gemini API クライアントの初期化**: Google Gemini 1.5 Flash モデルを使用します
3. **複数行入力の対応**: 空行が入力されるまで文章を受け付けます
4. **ドラえもん風レスポンス**: AI が秘密道具を使って悩みを解決してくれます

### 技術的なポイント

- **genai パッケージ**: Google の Generative AI Go SDK を使用
- **godotenv**: 環境変数を安全に管理
- **カスタム Scanner**: 複数行入力に対応した入力処理

## 3-3. 実行方法

プロジェクトのルートディレクトリで以下のコマンドを実行します。

```bash
go run main.go
```

# 4. 使用結果

実際にドラえもんボットを使用してみましょう！

## 4-1. 実行例

```
$ go run main.go
◯えもん: どうしたんだい？の▫太くん

> プレゼン資料の文章が固くて、聞き手の心に響かない
>

テレテテッテレーン!

ハートフル文章変換マシーン !!!

のび太くん、それは「ハートフル文章変換マシーン」だよ！この道具はね、固くて堅苦しい文章を、聞き手の心にスッと入り込むような温かい表現に変えてくれるんだ。専門用語は身近な例えに変えて、数字だけの説明には人の気持ちを込めて、一方的な説明から聞き手との対話みたいな文章にしてくれるよ。これを使えば、君のプレゼン資料もきっと聞いてる人の心に響く、親しみやすいものになるはずだよ！

◯えもん: 何か他に悩みはないかい？の▫太くん

> exit

終了します。
```

## 4-2. 特徴的な機能

### 🤖 ドラえもんらしい口調

- 「のび太くん」「〜だよ」「〜なんだ」などの親しみやすい表現
- 「テレテテッテレーン!」という効果音で秘密道具の登場を演出

### 🔧 創造的な秘密道具

- ユーザーの悩みに応じたオリジナルの秘密道具を提案
- 単なる添削ではなく、クリエイティブな解決方法を提供

### 💬 対話型インターフェース

- 複数行の入力に対応
- 連続した相談が可能
- `exit`で終了

## 4-3. 他の活用例

このドラえもんボットは文章の悩み以外にも様々な用途で活用できます：

- **メール文章の改善**: ビジネスメールを親しみやすい表現に
- **SNS投稿の最適化**: より魅力的な投稿文の提案
- **レポート作成支援**: 読みやすい文章構造のアドバイス
- **創作活動のサポート**: 小説やブログの文章改善

# 5. 技術的な学習ポイント

このプロジェクトを通して学べる技術要素：

## 5-1. Go言語の特徴
- **シンプルな構文**: C言語の影響を受けた読みやすいコード
- **並行処理**: Goルーチンによる効率的な処理（今回は未使用）
- **標準ライブラリ**: 豊富な標準パッケージ

## 5-2. 外部API連携
- **認証**: APIキーによる認証方式
- **HTTPクライアント**: RESTful APIとの通信
- **エラーハンドリング**: 適切な例外処理

## 5-3. 環境管理
- **dotenv**: 環境変数の安全な管理
- **設定の分離**: 開発・本番環境の設定分離

# 6. 今後の拡張案

## 6-1. 機能拡張
- **音声入力対応**: 音声認識APIとの連携
- **GUI化**: デスクトップアプリケーション化
- **Web化**: Webアプリケーションとして公開
- **履歴機能**: 過去の相談内容の保存

## 6-2. AI機能強化
- **感情分析**: 文章の感情を分析して適切な道具を選択
- **学習機能**: ユーザーの好みを学習してパーソナライズ
- **マルチモーダル**: 画像や音声にも対応

## 6-3. 他キャラクター対応
- **アンパンマン**: 正義感あふれるアドバイス
- **ピカチュウ**: ポケモン風の可愛い返答
- **カスタムキャラクター**: ユーザー独自のAIキャラクター作成

# まとめ

今回の体験を通して、AI が単なる情報提供だけでなく、クリエイティブな発想や擬人化された対話を通じて、私たちの日常の悩みを楽しく解決してくれる可能性を強く感じました。

## 📚 学んだこと

- **Go言語の実践的な活用**: 外部APIとの連携や環境変数管理
- **Gemini APIの可能性**: 高度な自然言語処理の活用
- **AIペルソナの創造**: キャラクター性を持たせたAIの実装
- **ユーザー体験の設計**: 対話型インターフェースの重要性

## 🚀 プロジェクトの意義

このプロジェクトは技術学習だけでなく、AIとの新しい関わり方を提案しています：

1. **親しみやすいAI**: 堅苦しい技術的な応答ではなく、愛らしいキャラクターとしてのAI
2. **創造的な問題解決**: 既存の解決法にとらわれない、新しいアプローチの提案
3. **エンターテイメント性**: 学習や問題解決を楽しい体験として提供

## 🎯 次のステップ

もし文章に悩んだら、あなたも自分だけの「ドラえもん」を作って、秘密道具の力を借りてみてはいかがでしょうか？ きっと、想像以上の面白いアドバイスがもらえるはずです！

**Happy Coding! 🎉**

---

*このプロジェクトがあなたのAI開発の冒険の始まりになることを願っています。「ドラえもん」のように、無限の可能性を秘めたAIの世界を一緒に探検しましょう！*
