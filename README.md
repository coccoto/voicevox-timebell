# 🔔 VOICEVOX 時報アプリ

VOICEVOX を使用して指定した時刻に音声で時報をお知らせするモダンな Web アプリケーションです。

## 🤖 開発について

このプロジェクトは、AI Agent によって設計・開発されました。人間の介入を最小限に抑え、効率的かつ迅速に構築されています。

## ✨ 特徴

- 🕐 複数の時刻を設定可能な時報機能
- 🎤 VOICEVOX の多様な話者キャラクター対応
- 🌐 モダンで使いやすいWebインターフェース
- 🐳 Docker による簡単なセットアップ

## 🚀 クイックスタート

### 前提条件

- Docker
- Docker Compose

### 起動方法

1. **リポジトリをクローン**
   ```bash
   git clone <repository-url>
   cd voicevox-timebell
   ```

2. **Docker Compose でアプリケーションを起動**
   ```bash
   cd docker
   sudo docker compose up --build
   ```

3. **ブラウザでアクセス**
   ```
   http://localhost:8080
   ```

## 🔧 API エンドポイント

### 設定保存
- **エンドポイント**: `POST /api/save`

### 設定取得
- **エンドポイント**: `GET /api/config`

### 手動音声再生
- **エンドポイント**: `POST /api/play`
