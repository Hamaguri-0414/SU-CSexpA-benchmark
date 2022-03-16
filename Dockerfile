#goのベース環境をもってくる
FROM golang:1.17-alpine as builder

#アップデートとgitのインストール
RUN apk update && apk add git alpine-sdk

#abコマンドのインストール
RUN apk --no-cache add apache2-utils

#sshコマンドのインストール
RUN apk add openssh

#タイムゾーンの設定
RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

#sshキーの追加
ADD .ssh /root/.ssh
RUN chmod 600 /root/.ssh*

#ワーキングディレクトリの設定
WORKDIR /go/src
#githubからclone
RUN git clone git@github.com:ohkilab/SU-CSexpA-benchmark.git
RUN git config --global user.email "aiba@sec.inf.shizuoka.ac.jp"
RUN git config --global user.name "expA-benchmark-container"

#ベンチマークサーバを起動
#CMD ["go","run","main.go"]