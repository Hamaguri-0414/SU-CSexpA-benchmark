package main

import (
  "log"
  "net/http"
  "text/template"
  "encoding/json"
	"fmt"
  "time"
  "os"
  //"reflect"
  //reflect.TypeOf(t)
  "benchmarkserver/internal/ab"
  "benchmarkserver/internal/record"
  "github.com/rs/xid"
)

// ajax戻り値のJSON用構造体
type Param struct {
	Time string
  Msg string
}

func main() {
  // webフォルダにアクセスできるようにする
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css/"))))
  http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("./web/script/"))))
  http.Handle("/gif/", http.StripPrefix("/gif/", http.FileServer(http.Dir("./web/gif/"))))

  //ルーティング設定 "/"というアクセスがきたら rootHandlerを呼び出す
  http.HandleFunc("/", rootHandler)
  http.HandleFunc("/measure", measureHandler)

  log.Println("Listening...")
  // 3000ポートでサーバーを立ち上げる
  http.ListenAndServe(":3000", nil)
}

//main画面
func rootHandler(w http.ResponseWriter, r *http.Request) {
  //index.htmlを表示させる
	tmpl := template.Must(template.ParseFiles("./web/html/index.html"))
	tmpl.Execute(w, nil)
}

//フォームからの入力を処理 index.jsから受け取る
func measureHandler(w http.ResponseWriter, r *http.Request) {

  //ログファイルを開く
  logfile := logfileOpenPush()
  defer logfile.Close()

  //index.jsに返すJSONデータ変数
  var ret Param
  //POSTデータのフォームを解析
  r.ParseForm()

  url := r.Form["url"][0]
  groupName := r.Form["groupName"][0]

  //idを設定(logを対応づけるため)
  guid := xid.New()
  log.Println("<Info> request URL: " + url + ", GroupName: " + groupName + ", id: " + guid.String())
  fmt.Fprintln(logfile, time.Now().Format("2006/01/02 15:04:05") + "<Info> request URL: " + url + ", GroupName: " + groupName + ", id: " + guid.String())

  //abコマンドで負荷をかける．計測時間を返す
	ret.Msg, ret.Time = ab.Ab(logfile, guid.String(), url)

  //計算結果を記録する
  if ret.Msg == "" {
    ret.Msg = record.Record(logfile, guid.String(), ret.Time, groupName)
  }

	// 構造体をJSON文字列化する
	jsonBytes, _ := json.Marshal(ret)
  // index.jsに返す
  fmt.Fprintf(w, string(jsonBytes))
}

//ログファイルを開く，ログファイルをgithubにpushする
func logfileOpenPush() *os.File {

  //ログファイルを開く(logを記録するファイル)
  logfile, err := os.OpenFile("data/log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
  if err != nil {
    log.Println(fmt.Sprintf("<Debug> ", err))
  }
  return logfile
}
