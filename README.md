# Go チュートリアルのサンプルアプリケーション

これは、次の教材で作られたサンプルアプリケーションです。
[_Go チュートリアル_](https://gotutorial.jp/)
（第7版）

## ライセンス

[Go チュートリアル](https://gotutorial.jp/)内にある
ソースコードはMITライセンスとBeerwareライセンスのもとで公開されています。
詳細は [LICENSE.md](LICENSE.md) をご覧ください。

## 使い方

このアプリケーションを動かす場合は、まずはリポジトリをクローンしてください。

次に、依存パッケージをダウンロードします。

```
$ go mod tidy
```

`templ` CLIを`go.mod`と同じバージョンでインストールします。

```
$ go install github.com/a-h/templ/cmd/templ@v0.3.833
```

templのコード生成を実行します。

```
$ templ generate
```

サーバーを起動します。

```
$ go run cmd/server/main.go
```

最後に、テストを実行してうまく動いているかどうか確認してください。

```
$ go test ./...
```

詳しくは、[_Go チュートリアル_](https://gotutorial.jp/)
を参考にしてください。
