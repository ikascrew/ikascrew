# ikascrew

**A VJ (video jockey) system written in Go — manage your clips, drive them from a controller, project them full-screen.**

**ikascrew は Go で書かれた VJ(ビデオジョッキー)システムです。** 映像素材を管理し、コントローラーから操作して、フルスクリーンに投影します。

This repository is the **umbrella** for the system: it holds the overall documentation and points at the component repositories. It contains no Go code of its own.

このリポジトリはシステム全体の**傘**にあたり、全体ドキュメントと各コンポーネントへの案内だけを持ちます。Go のコードは含みません。

---

## Components / 構成

The system is split into three applications and five shared libraries, each in its own repository under [`github.com/ikascrew`](https://github.com/ikascrew).

システムは3つのアプリケーションと5つの共有ライブラリに分かれ、それぞれが独立したリポジトリです。

### Applications / アプリケーション

| Repository | Role | Port |
|---|---|---|
| [`ikasbox`](https://github.com/ikascrew/ikasbox) | **Manage.** Content library for videos and images — Go server + React SPA, SQLite metadata, thumbnail generation.<br>**管理。** 動画・画像コンテンツのライブラリ。 | HTTP `:5555` |
| [`server`](https://github.com/ikascrew/server) | **Show.** Playback engine that draws to an OpenCV window and is driven remotely over gRPC.<br>**表示。** OpenCV ウィンドウに描画する再生エンジン。 | gRPC `:55555` |
| [`client`](https://github.com/ikascrew/client) | **Control.** Ebitengine operator UI — thumbnail grid, Next queue, volume — sending commands to `server`.<br>**操作。** Ebitengine 製のコントロール UI。 | — |

### Shared libraries / 共有ライブラリ

| Repository | Role |
|---|---|
| [`core`](https://github.com/ikascrew/core) | Window handling, frame types, UDP multicast discovery / Window・フレーム型・マルチキャスト探索 |
| [`plugin`](https://github.com/ikascrew/plugin) | Video source implementations: `file`, `img`, `cd`, `terminal` / 映像ソース実装 |
| [`pb`](https://github.com/ikascrew/pb) | gRPC protocol definitions shared by `server` and `client` / gRPC プロトコル定義 |
| [`powermate`](https://github.com/ikascrew/powermate) | Griffin PowerMate dial input / PowerMate ダイヤル入力 |
| [`volumes`](https://github.com/ikascrew/volumes) | Value/volume management used by the operator UI / 値の管理ロジック |

---

## How it works / 全体のデータフロー

```
[Preparation / 準備フェーズ]      ikasbox (HTTP :5555)
    ├─ ika-server create <project-id>  →  server/.server/config.json
    └─ ika-client create <project-id>  →  client/.client/contents.json + thumbnails

[Performance / 実行フェーズ]      ikasbox not required / ikasbox 不要
    client (Ebitengine UI) ──gRPC :55555──▶ server (OpenCV rendering)
                              Effect / Switch / PutVolume / Sync
```

1. Import your material into **ikasbox** as a group, then turn that group into a project. Metadata lives in SQLite.
   **ikasbox** に素材をグループとして import し、そのグループをプロジェクト化します。
2. **server** and **client** each run their own `create` subcommand against ikasbox to pull the project down into a local work file. Both must be created from the **same project** so that content ID → path mappings agree.
   両者は同じプロジェクトから `create` する必要があります（コンテンツ ID とパスの対応を一致させるため）。
3. At showtime (`start`) ikasbox is no longer needed. `client` drives `server` over gRPC:
   本番時は ikasbox 不要で、client が gRPC で server を制御します。
   - `Effect` — load and push a clip by content ID / コンテンツ ID で映像をロード・プッシュ
   - `Switch` — crossfade to next/prev / next・prev でクロスフェード切替
   - `PutVolume` — send Volume (crossfade) / Light / Wait values / 3種の値を送る
   - `Sync` — toggle fullscreen ⇄ windowed / フルスクリーン切替
4. `server` and `ikasbox` announce themselves over UDP multicast (`224.0.0.224:15496`); `client` listens at startup and configures its gRPC target automatically, falling back to `localhost:55555`.
   server と ikasbox はマルチキャストで自身を告知し、client は起動時にこれを受信して接続先を自動設定します。

---

## Requirements / 動作要件

- **Go 1.26+** (the newest component modules declare `go 1.26.1`)
- **OpenCV 4.x** — required by [gocv](https://gocv.io/) to build `core`, `plugin`, `ikasbox` and `server`.
  `pb`, `powermate` and `volumes` build without it.
  `core` / `plugin` / `ikasbox` / `server` のビルドに OpenCV が必要です。
- A GPU/display capable of full-screen video playback / フルスクリーン再生できる表示環境

---

## Getting started / 起動手順

The component repositories expect to sit side by side in one directory, because several of them reference each other with local `replace` directives.

各リポジトリは互いに `replace` でローカル参照するため、**同一ディレクトリに横並びで**配置してください。

```sh
mkdir ikascrew && cd ikascrew
for r in ikasbox server client core plugin pb powermate volumes; do
  git clone https://github.com/ikascrew/$r.git
done
```

```powershell
# PowerShell
mkdir ikascrew; cd ikascrew
'ikasbox','server','client','core','plugin','pb','powermate','volumes' |
  ForEach-Object { git clone "https://github.com/ikascrew/$_.git" }
```

Then:

```sh
# 1. Content management / コンテンツ管理 (first run: init → group import → create project)
cd ikasbox/cmd && go run main.go start          # :5555

# 2. Preparation / 準備 (once, while ikasbox is running)
cd server && go run ./cmd/ika-server create <project-id>
cd client && go run ./cmd/ika-client create <project-id>

# 3. Performance / 本番 (ikasbox not needed)
cd server && go run ./cmd/ika-server start      # gRPC :55555 + OpenCV window
cd client && go run ./cmd/ika-client start      # operator UI (-pm enables the PowerMate dial)
```

---

## Repository layout / このリポジトリについて

This repository tracks only the shared documentation:

```
README.md      ← you are here / このファイル
CLAUDE.md      ← guidance for Claude Code across the whole system
docs/          ← cross-repository documents / リポジトリ横断のドキュメント
```

The eight component directories are ignored via `.gitignore` — clone them yourself as shown above. Each carries its own `README.md`, and the three applications also carry a detailed `CLAUDE.md`.

8つのコンポーネントディレクトリは `.gitignore` で除外しています。各リポジトリは個別に clone してください。

---

## History / 沿革

ikascrew started in 2017 as a single monolithic repository (still visible in this repository's `master` branch) and was later split into the component repositories described above.

ikascrew は 2017 年に単一のモノリシックなリポジトリとして始まり(このリポジトリの `master` ブランチに当時のコードが残っています)、その後現在の複数リポジトリ構成に分割されました。
