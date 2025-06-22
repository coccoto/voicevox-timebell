#!/bin/bash

# cron をバックグラウンドで起動する
cron

# Go アプリをフォアグラウンドで起動する
exec /app/dist/app
