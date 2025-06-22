# voicevox-timebell
自宅で使ってる時報アプリです。お気に入りの VOICEVOX のキャラクターが時間をお知らせします。

## 動作環境
- Raspberry Pi 5
- Ubuntu 24.04

## 動作手順
- git clone https://github.com/coccoto/voicevox-timebell.git
- cd voicevox-timebell
- sudo docker compose up --build

## 音量調整
- amixer sset 'Master' {volume}%
