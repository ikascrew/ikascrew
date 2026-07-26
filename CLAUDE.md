# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working in this directory.

> ## ⚠️ 2026-07-26 のデータ消失について
>
> `D:\Go\Projects\ikascrew` を誤って `rm -rf` したため、**各リポジトリの作業ツリーは GitHub の master から復元した**。復元時点の各リポジトリ最終コミットは **2026-07-05〜07-08**。
> **7/8 以降の未 push 作業は失われている。** このファイルの記述は消失前(7/26 時点)のものを保持しているため、**下記の項目はドキュメントにはあるがコードには存在しない**:
>
> - `core.Frame` 型による gocv 依存の分離 — 復元コードでは `ikasbox` / `server` がまだ `gocv.io/x/gocv` を直接 import している(3箇所)
> - `plugin/video/param` パッケージ(JSON params のパース層)
> - server の同居モード `-ikasbox` / `-db` フラグ
>
> 一方、`video.Normalize` による語彙正規化、`ikasbox content register` サブコマンド、`contents.params` カラム、work file 基準の `Effect` 解決、`core/multicast` は復元コードにも**存在する**。
>
> これらの再実装が必要。以後は**こまめに push すること**。

## このディレクトリについて

**ikascrew** — VJ(ビデオジョッキー)用の映像再生システムを構成する複数の Go リポジトリを並べた作業ディレクトリ。中心となるのは次の3パッケージで、それぞれに詳細な `CLAUDE.md` がある — **各パッケージ内で作業するときは必ずそのパッケージの CLAUDE.md を読むこと。**

| パッケージ | 役割 | ポート |
|---|---|---|
| `ikasbox/` | 動画・画像コンテンツの管理(Go サーバー + React SPA、SQLite、GoCV サムネイル生成) | HTTP :5555 |
| `server/` | 動画の表示。OpenCV ウィンドウに映像を描画し、gRPC で遠隔制御される再生エンジン | gRPC :55555 |
| `client/` | 表示のコントロール。Ebitengine 製の操作 UI(サムネイル一覧・Next キュー・ボリューム)から server に gRPC で指示を送る | — |

残りの `core`, `plugin`, `pb`, `powermate`, `volumes` は上記3つが依存する共有ライブラリ:

- `pb/` — gRPC プロトコル定義。server/client 間のプロトコル変更はここを触る必要がある。
- `core/` — Window・multicast などの共通基盤。
- `plugin/` — server が使う `core.Video` 実装("file", "img", "cd", "terminal")。
- `powermate/`, `volumes/` — client の `go.mod` が `replace` でローカル参照している(変更は両リポジトリでコミット/プッシュが必要)。

**このディレクトリ配下はこの8リポジトリのみ**。VJ システムと無関係だった旧・実験リポジトリ(`brender`, `generator`, `pose`, `ikamera`, `xbox`, `go-mp4`, `gocv`, および `ikascrew-stream`, `ikascrew-mp4`, `ikascrew-viewer`, `ikascrew-util`)は 2026-07-26 に `D:\Go\Projects\` 直下へ分離済み。ここに新しいディレクトリを増やす場合は VJ システムの構成要素かを確認すること。

このディレクトリ自体も `github.com/ikascrew/ikascrew` として git 管理されている(`main` ブランチ)。ただし追跡するのは `README.md` / `CLAUDE.md` のみで、8つのサブリポジトリは `.gitignore` で除外している。

## 全体のデータフロー

```
[準備フェーズ]  ikasbox (HTTP :5555)
    ├─ ika-server create <project-id> → server/.server/config.json (ID→パスのマップ)
    └─ ika-client create <project-id> → client/.client/contents.json + サムネイル

[実行フェーズ]  ikasbox 不要
    client (Ebitengine UI) ──gRPC :55555 (Effect/Switch/PutVolume/Sync)──▶ server (OpenCV 描画)
```

1. **ikasbox** にコンテンツをグループとして import し、グループをプロジェクト化して管理する。メタデータは SQLite(`ikasbox.db`)。
2. **server** と **client** は各自の `create` サブコマンドで ikasbox からプロジェクト情報を取得し、ローカルのワークファイルに落とす。**コンテンツ ID→パスのマッピングは server/client 間で一致している必要がある**(同じプロジェクトで create すること)。
3. 本番(`start`)時は ikasbox への依存なし。client が gRPC(`github.com/ikascrew/pb`)で server を制御する:
   - `Effect` — コンテンツ ID を指定して動画をロード・プッシュ(型と JSON params は server が work file から解決)
   - `Switch` — next/prev でクロスフェード切替
   - `PutVolume` — Volume(クロスフェード値)/ Light / Wait の3種の値を送る
   - `Sync` — フルスクリーン⇔ウィンドウ切替
4. server と ikasbox は `ikascrew/core/multicast` の UDP マルチキャスト(224.0.0.224:15496)で自身を告知し、client は起動時にこれを受信して gRPC 接続先を自動設定する(見つからない場合は `localhost:55555` にフォールバック)。マルチキャストは全インターフェースへの送信/join で実装されている(Windows の複数 NIC 対策)。ただし `create` 時の ikasbox 接続先(client/server とも)は `localhost:5555` 固定のまま。

## コンテンツの型と JSON params(一貫パイプライン)

コンテンツの型は `plugin/video` レジストリの正語彙 **`file` / `img` / `cd` / `terminal`** に統一されている(`video.Normalize` が旧語彙 "image"/"countdown" を吸収)。プラグイン固有のパラメータは **JSON 文字列**で、次の規約に従う:

- **運ぶ層は素通し**: ikasbox の DB(`contents.params` TEXT)→ `/project/content/list` API → server の work file → プラグイン、と不透明な文字列のまま流れる。中間層で params をパースしない。
- **解釈は各プラグインの `New(param string)` だけ**が行う(非 JSON は旧来の生文字列=パスやテキストとして解釈)。
- **検証は登録時に前倒し**: `ikasbox content register <group-id> <name> <type> [params-json]` が実際にプラグインを生成して params を検証し、サムネイルもプラグインの描画フレームから作る(実体ファイル不要の「生成型コンテンツ」)。
- server の `Effect` はコンテンツ ID から work file の Type/Params で解決するため、**pb(gRPC プロトコル)は変更不要**。client が送る Type は旧 work file の救済用でしかない。

例: `content register 1 "NewYear" cd "{\"target\":\"2027-01-01T00:00:00+09:00\",\"text\":\"HNY\"}"`

型ごとの params スキーマと新プラグインの追加手順は **`plugin/CLAUDE.md`** に明文化されている(語彙・param 形式の本家は plugin リポジトリ)。

## 共通の前提・慣習

- **OpenCV(gocv)** が ikasbox と server のビルドに必須。
- 設定は各パッケージとも functional options によるグローバルシングルトン(`config.Set(opts...)` / `config.Get()`)パターン。
- エラーは `golang.org/x/xerrors` でラップ。コメント・コミットメッセージは日本語が多く、コミット件名は `fix:` / `feat:` / `perf:` / `chore:` プレフィックス。
- テストはほぼ存在しない(ikasbox の `db/` に1ファイルのみ)。
- 各リポジトリは GitHub の同一 org(`github.com/ikascrew/*`)配下で、pseudo-version 依存で相互参照している。

## 起動手順(最短)

```
# 1. コンテンツ管理(初回は init → group import → project 作成)
cd ikasbox/cmd && go run main.go start          # :5555

# 2. 準備(ikasbox 起動中に一度だけ)
cd server && go run ./cmd/ika-server create <project-id>
cd client && go run ./cmd/ika-client create <project-id>

# 3. 本番(ikasbox 不要)
cd server && go run ./cmd/ika-server start      # gRPC :55555 + OpenCV ウィンドウ
cd client && go run ./cmd/ika-client start      # 操作 UI(-pm で PowerMate ダイヤル有効)
```
