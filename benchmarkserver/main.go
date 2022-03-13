package main

import (
  "log"
  "net/http"
  "text/template"
  "encoding/json"
	"fmt"
  //"reflect" reflect.TypeOf(t)
  "benchmarkserver/pkg/ab"
  "benchmarkserver/pkg/record"
)

//main画面
func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./web/html/index.html"))
	tmpl.Execute(w, nil)
}


// ajax戻り値のJSON用構造体
type Param struct {
	Time string
  Msg string
}

//フォームからの入力を処理 index.jsから受け取る
func measureHandler(w http.ResponseWriter, r *http.Request) {

  //返すJSONデータ変数
  var ret Param

  //フォームを解析
  r.ParseForm()

  url := r.Form["url"][0]
  groupName := r.Form["groupName"][0]

  log.Println(url)
  log.Println(groupName)

  //エラー処理　再入力させる
  if(len(url) == 0){
    w.Write([]byte("URLを入力してください"))
    return
  }

  //計算結果を戻り値に入れる
	ret.Time = ab.Ab(url)

  //計算結果を記録する
  ret.Msg = record.Record(ret.Time, groupName)

	// 構造体をJSON文字列化する
	jsonBytes, _ := json.Marshal(ret)
  // 返す
  fmt.Fprintf(w, string(jsonBytes))


}


func main() {
  // webフォルダにアクセスできるようにする
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css/"))))
  http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("./web/script/"))))
  http.Handle("/gif/", http.StripPrefix("/gif/", http.FileServer(http.Dir("./web/gif/"))))

  //ルーティング設定。"/"というアクセスがきたらstaticディレクトリのコンテンツを表示させる
  http.HandleFunc("/", rootHandler)
  http.HandleFunc("/measure", measureHandler)

  log.Println("Listening...")
  // 3000ポートでサーバーを立ち上げる
  http.ListenAndServe(":3000", nil)
}
