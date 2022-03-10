# 実験Aベンチマークサーバ
[![deploy](https://github.com/ohkilab/expA-benchmarkserver/actions/workflows/main.yml/badge.svg)](https://github.com/ohkilab/expA-benchmarkserver/actions/workflows/main.yml)

https://ohkilab.github.io/expA-benchmarkserver/index.html
## 導入方法
1. [公式サイト](https://go.dev/dl/)からgoをダウンロード
   > `$go version`でgoの存在確認
3. ターミナルを再起動
4. リポジトリからファイルをクローン<br>
   `$git clone git@github.com:ohkilab/expA-benchmarkserver.git`
## ベンチマークサーバ起動方法
1. `$go run main.go`でベンチマークサーバを起動
   > **注意** <br>
   main.goプログラムのディレクトリ(benchmarkserver/)をカレントディレクトリにする必要があります<br>
   `$cd benchmarkserver`でmaing.goをカレントディレクトリにする
2. `http://localhost:3000`または`http://<ipアドレス>:3000`でアクセス
