# ikascrew テスト計画(2026-07-10)

対象: core / plugin / volumes / powermate / ikasbox / stream / client / server の8リポジトリ。
目的: GUI・実デバイス・実ネットワークに依存しないユニットテストを整備し、`go test ./...` が各リポジトリで安定して通る状態にする。

## 共通方針

- 標準 `testing` パッケージのテーブル駆動テスト。新規依存ライブラリは追加しない。
- gocv はこの環境でビルド・実行可能(確認済み)。`gocv.NewMat` 等で合成フレームを作ってテストしてよい。ただし **ウィンドウ表示(highgui)は使わない**。
- 除外対象: OpenCV ウィンドウ描画(core/window, server の描画)、Ebitengine のゲームループ(client/window の描画)、PowerMate 実デバイス入出力、外部サイトへのスクレイピング(stream の agouti 系)。
- ファイル I/O が必要なテストは `t.TempDir()` を使う。SQLite は一時ディレクトリの DB ファイルで実体テスト。
- ネットワークは loopback 限定(httptest / bufconn / 127.0.0.1 UDP)。マルチキャストのテストは環境依存で不安定なら `t.Skip` 可。
- `.claude/worktrees/` 以下は過去の作業コピーなので無視する。
- コミットはしない(実装のみ。レビュー後に人間がコミット)。

## リポジトリ別計画

### core
- `util/`: ファイル列挙・パス処理・動画ユーティリティの純粋関数をテーブル駆動でテスト。
- `effect.go` / `transition.go` / `image.go` / `video.go`: 合成 Mat を使った変換・状態遷移のテスト。
- `sync/`: 既存テストの拡充(境界値)。
- `multicast/`: メッセージの組み立て/解釈部分の単体テスト。実 UDP は 127.0.0.1 で送受信できる範囲のみ、不安定なら skip。
- `window/`, `cmd/`: 対象外。

### plugin
- `video/video.go`: `Normalize`(旧語彙 image/countdown の吸収)とレジストリの登録・解決。
- `video/param/`: JSON params のパース(正常系・非 JSON の旧形式・不正 JSON)。
- 各プラグインの `New(param)`: countdown / terminal / telop の既存テスト拡充、image / file は一時ファイル(gocv imwrite で生成)を使った検証、output。
- `effect/light`, `effect/speed`, `transition/switch`: 合成 Mat に適用して出力の性質(サイズ・チャンネル・値域)を検証。
- `cap.go` / `version.go`: 可能な範囲で。

### volumes
- `volumes.go`: 値の管理ロジック(範囲、増減、正規化)を既存 `volumes_test.go` を拡充してカバー。termbox の UI 部分は対象外。

### powermate
- `action_test.go` の拡充: HID イベント → アクション変換のテーブル駆動テスト(回転方向・押下・境界値)。
- `listen_windows.go` / `listen_linux.go`: 実デバイスが必要なので対象外。

### ikasbox
- `db/`: 一時 SQLite で content / group / project / project_group の CRUD・paging をテスト。既存の `content_thumbnail_test.go` を参考に。
- `handler/api/`: httptest + 一時 DB でエンドポイントのテスト(既存 `api_test.go`, `content_test.go` を拡充)。特に content register 時の params 検証パス。
- `config/`: functional options のテスト。
- `contentimport/`: 一時ディレクトリでのインポートロジック。
- React SPA / cmd は対象外。

### stream
- `m3u.go`: 既存テスト拡充(パースの異常系)。
- `utils/`: search / utils の純粋関数。
- `plugins/local/`: 一時ディレクトリでのローカルプラグイン動作。
- `handler/`: httptest で可能な範囲。websocket と `_cmd/`・streamsb(スクレイピング)は対象外。

### client
- `tool/`: 既存 `contents_test.go` の拡充、project / tool のロジック(ikasbox への HTTP は httptest でモック)。
- `config/`: functional options。
- `request.go`: gRPC リクエスト組み立てのロジック(bufconn またはフェイクで可能な範囲)。
- `discover.go`: マルチキャスト受信のフォールバック(localhost:55555)ロジックを切り出せる範囲で。
- `window/`(Ebitengine 描画)、`powermate.go`、`input.go` の実デバイス部分は対象外。

### server
- `config/`: functional options。
- `video_gen.go` / `handler.go`: work file(config.json)の ID→Type/Params 解決ロジックを一時ファイルでテスト。gRPC ハンドラは bufconn で描画に依存しない範囲(PutVolume 等)のみ。
- `stream.go` / `opening.go`: 純粋ロジック部分。
- OpenCV ウィンドウ描画(`window.go`)は対象外。

## 実施方法

各リポジトリに Sonnet サブエージェントを割り当てて並列実装。各エージェントは
1. リポジトリの CLAUDE.md(あれば)と対象コードを読む
2. 上記計画に沿ってテストを実装
3. `go test ./...` をリポジトリルートで実行し全パスを確認
4. 追加テストの一覧とカバーできなかった箇所を報告
