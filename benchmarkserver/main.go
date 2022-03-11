package main

import (
  "log"
  "net/http"
  "os/exec"
  "text/template"
  "encoding/json"
	"fmt"
  "regexp"
  "strings"
  "io"
  "os"
  "encoding/csv"
  "strconv"
  //"reflect" reflect.TypeOf(t)
  "io/ioutil"
  "math/rand"
  "time"
)

//main画面
func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./html/index.html"))
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
	ret.Time = ab(url)

  //計算結果を記録する
  ret.Msg = record(ret.Time, groupName)

	// 構造体をJSON文字列化する
	jsonBytes, _ := json.Marshal(ret)
  // 返す
  fmt.Fprintf(w, string(jsonBytes))

}

//abコマンドで負荷をかけ，計測結果を返す
func ab(url string) string {

  var file []uint8
  var execRes string
  var out []byte
  var reg string
  var splitExecRes []string
  var measureTimes float64
  measureTimes = 0

  //ランダムタグで検索
  //改ざんチェック
  file, _ = ioutil.ReadFile("./data/randomtag.txt")
  randomTags := strings.Split(string(file), "\n")
  rand.Seed(time.Now().UnixNano())
  randomTag := randomTags[rand.Intn(len(randomTags))]
  log.Println(randomTag)
  //http://192.168.1.101/~username/directory/progC.php?tag=fiat
  out, _ = exec.Command("ab", "-c", "1", "-n", "1", url + "?tag=" + randomTag).Output()
  execRes = string(out)
  //abコマンドの結果を:と改行で分割する
  reg = "[:\n]"
  splitExecRes = regexp.MustCompile(reg).Split(execRes, -1)
  //分割したものからRComplete requestsを探す
  //次にあるのが計測値なので，i+1して指定，空白で分割し，数値のみ取り出す
  //例：Complete requests:      1
  for i, s := range splitExecRes {
    if s == "Complete requests" {
      ss := strings.Split(splitExecRes[i + 1], " ")
      //正常にレスポンスが帰ってこなかった
      //検索タグ以外のタグで検索されているかも
      if ss[len(ss) - 1] == "0" {
        //改ざんされているよ通知etc
        log.Println("不正")
      }else{
        log.Println("正常")
      }
      break
    }
  }

  //複数タグで検索し，計測
  file, _ = ioutil.ReadFile("./data/searchtag.txt")
  tags := strings.Split(string(file), "\n")
  for i, s := range tags {
    //検索タグ数（本番は30とかにしたい）
    if i == 2 || i > len(tags){
      break
    }
    log.Println(s)
    //-c -nを変更する
    out, _ = exec.Command("ab", "-c", "1", "-n", "1", url + "?tag=" + s).Output()
    execRes = string(out)
    //abコマンドの結果を:と改行で分割する
    reg = "[:\n]"
    splitExecRes = regexp.MustCompile(reg).Split(execRes, -1)
    //分割したものからRequests per secondを探す
    //次にあるのが計測値なので，j+1して指定，空白で分割し，数値のみ取り出す
    //例：Requests per second:    720.46 [#/sec] (mean)
    for j, s := range splitExecRes {
      if s == "Requests per second" {
        ss := strings.Split(splitExecRes[j + 1], " ")
        log.Println(ss[len(ss) - 3])
        //float64に変換して加算
        measureTime, _ := strconv.ParseFloat(ss[len(ss) - 3], 64)
        measureTimes += measureTime
        break
      }
    }

  }

  //文字列にして返す
  return strconv.FormatFloat(measureTimes, 'f', 2, 64)
}

func record(times string, groupName string) string {

  msg := ""
  recordData := ""
  doUpdate := false

  //data.csvに記録する
  //data.csvを読み込む
  csvFile, _ := os.Open("../public/score.csv")
  reader := csv.NewReader(csvFile)

  //groupNameの一致を探し，数値を比較する
  for {
    line, err := reader.Read()
    if err == io.EOF {
        break
    }
    //書き込みデータを作成する
    recordData += line[0] + ","
    //グループ名を探し，計測時間を比較
    if line[0] == groupName {
      nowData, _ := strconv.ParseFloat(times, 64)
      highData, _ := strconv.ParseFloat(line[1], 64)
      if nowData > highData {
        recordData += times + "\n"
        msg = "記録更新！！！"
        doUpdate = true
      }else{
        recordData += line[1] + "\n"
        msg = "記録更新ならず，現在の最高値：" + line[1]
      }
    }else{
      recordData += line[1] + "\n"
    }
  }
  csvFile.Close()
  //ファイル書き込み
  file, _ := os.Create("../public/score.csv")
  defer file.Close()
  _, err := file.WriteString(recordData)
  if err != nil {
    log.Println(err)
  }

  //csvファイルをgithubにpush
  if doUpdate {
    csvPush(groupName)
  }

  return msg

}



func csvPush(groupName string){
  var err error
  //git add ../exp1_ranking/public/score.csv
  err = exec.Command("git", "add", "../public/score.csv").Run()
  if err != nil {
    log.Println(err)
  }
  err = exec.Command("git", "commit", "-m", groupName + "の記録更新").Run()
  if err != nil {
    log.Println(err)
  }
  err = exec.Command("git", "push").Run()
  if err != nil {
    log.Println(err)
  }

}



func main() {
  // css、scriptフォルダにアクセスできるようにする
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
  http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("script/"))))
  http.Handle("/gif/", http.StripPrefix("/gif", http.FileServer(http.Dir("gif/"))))

  //ルーティング設定。"/"というアクセスがきたらstaticディレクトリのコンテンツを表示させる
  http.HandleFunc("/", rootHandler)
  http.HandleFunc("/measure", measureHandler)

  log.Println("Listening...")
  // 3000ポートでサーバーを立ち上げる
  http.ListenAndServe(":3000", nil)
}
