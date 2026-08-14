天気の子
====

## Overview

個人遊び用のLinebot. 

![Demo](./docs/sendWeatherDemo.gif)

## Description
- 毎朝6時に天気情報配信
- 天気を聞くと設定された地域の天気予報を教えてくれる
- 返信イベントは[ojichat](https://github.com/greymd/ojichat)搭載

## Requirement
- Go 1.25 or later
- MongoDB

以前は `github.com/yuki9431/logger` / `mongohelper` / `weather` を外部リポジトリとして
参照していたが、`internal/` 配下に取り込んだため個別の取得は不要になった。

## Contribution
1. Fork ([https://github.com/yuki9431/myLinebot](https://github.com/yuki9431/myLinebot))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Create new Pull Request


## Author
[Dillen H. Tomida](https://github.com/yuki9431)
