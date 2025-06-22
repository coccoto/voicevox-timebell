# voicevox-timebell
自宅で使ってる時報アプリです。お気に入りの VOICEVOX のキャラクターが毎時時間をお知らせします。

## 動作環境
- Raspberry Pi 5
- Ubuntu 24.04

## 動作環境構築手順
1. プロジェクトをセットアップする。
    - git clone https://github.com/coccoto/voicevox-timebell.git
    - cd voicevox-timebell
    - sudo docker compose up --build
2. 設定画面 (http://localhost/) にアクセスして、以下の項目を設定する。
    - 通知する時間
    - キャラクター選択
    - 音声スタイル
3. 「テスト再生」ボタンを押下して、正常に音声が再生されることを確認する。

## 音量調整
- amixer sset 'Master' {volume}%

## API エンドポイント一覧

| メソッド | エンドポイント | 説明 |
|---------|---------------|------|
| GET | `/api/alert` | 時報を再生する。cron から定期的にリクエストされる。 |
| GET | `/api/config-read` | 設定を取得する。 |
| POST | `/api/config-register` | 設定を更新する。 |
| GET | `/api/speakers` | 利用可能な VOICEVOX のキャラクター一覧を取得する。 |
