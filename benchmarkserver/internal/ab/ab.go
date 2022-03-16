package ab

import (
  "io/ioutil"
  "math/rand"
  "time"
  "strings"
  "log"
  "os/exec"
  "regexp"
  "strconv"
  "fmt"
  "os"
)

//abコマンドで負荷をかけ，計測結果を返す
func Ab(logfile *os.File, id string, url string) (string, string) {

  var measureTimes float64 //計測時間の合計
  measureTimes = 0

  //ランダムタグで検索
  //(改ざんチェック)
  //タグファイル（randomtag.txt）からランダムにタグを抽出し，そのタグでabコマンドを実行する
  file, err := ioutil.ReadFile("./data/randomtag.txt")

  if err != nil {
    log.Println(fmt.Sprintf("<Debug> ", err))
  }

  //randomtag.txtを改行で配列に分割し，分割した配列の中からランダムでひとつを選択する
  randomTags := strings.Split(string(file), "\n")
  rand.Seed(time.Now().UnixNano())
  randomTag := randomTags[rand.Intn(len(randomTags))]
  log.Println("<Info> id: " + id + ", selected tag: " + randomTag)
  fmt.Fprintln(logfile, time.Now().Format("2006/01/02 15:04:05") + "<Info> id: " + id + ", selected tag: " + randomTag)

  //選択されたタグを使用してabコマンドを実行
  //http://192.168.1.101/~username/directory/progC.php?tag=fiat
  //-c -nを変更する
  out, err := exec.Command("ab", "-c", "1", "-n", "1", url + "?tag=" + randomTag).Output()

  if err != nil {
    //urlが不明
    log.Println(fmt.Sprintf("<Error> id: " + id + " execCmd(ab -c 1 -n 1 " + url + "?tag=" + randomTag + ")" , err))
    fmt.Fprintln(logfile, time.Now().Format("2006/01/02 15:04:05") + fmt.Sprintf("<Error> id: " + id + " execCmd(ab -c 1 -n 1 " + url + "?tag=" + randomTag + ")" , err))

    return "URLが不明です", "0.00"
  }

  execRes := string(out)
  //abコマンドの結果を:と改行で分割する
  reg := "[:\n]"
  splitExecRes := regexp.MustCompile(reg).Split(execRes, -1)
  //分割したものからRequests per secondを探す
  //次にあるのが計測値なので，i+1して指定，空白で分割し，数値のみ取り出す
  //例：Requests per second:    720.46 [#/sec] (mean)
  for i, s := range splitExecRes {
    if s == "Requests per second" {
      ss := strings.Split(splitExecRes[i + 1], " ")

      log.Println("<Info> id: " + id + ", Requests per second: " + ss[len(ss) - 3])
      fmt.Fprintln(logfile, time.Now().Format("2006/01/02 15:04:05") + "<Info> id: " + id + ", Requests per second: " + ss[len(ss) - 3])

      //float64に変換して加算
      measureTime, _ := strconv.ParseFloat(ss[len(ss) - 3], 64)
      measureTimes += measureTime
      break
    }
  }

  //curlでhtmlを取得し，imgタグ内の.staticflickr.comの数が100個あるか数える
  //htmlが正常か簡易的にチェック
  if !Checkhtml(logfile, id, url, randomTag) {
    return "HTMLファイルが改ざんされている可能性があります", "0.00"
  }

  /*
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

  */

  //文字列にして返す
  return "", strconv.FormatFloat(measureTimes, 'f', 2, 64)
}

//htmlファイルが簡易的に正常かどうか確認する
func Checkhtml(logfile *os.File, id string, url string, tag string) bool {
  //.staticflickr.comという文字列が何個あるか確認する
  //.staticflickr.comは，Flickrサーバ上の画像URL	http://farm5.staticflickr.com/40～略～m.jpgで使われている

  count := 0

  //curlでhtmlを取得する
  out, err := exec.Command("curl", url + "?tag=" + tag).Output()

  if err != nil {
    log.Println(fmt.Sprintf("<Error> id: " + id + " execCmd(curl " + url + "?tag=" + tag + ")" , err))
    fmt.Fprintln(logfile, time.Now().Format("2006/01/02 15:04:05") + fmt.Sprintf("<Error> id: " + id + " execCmd(curl " + url + "?tag=" + tag + ")" , err))
    return false
  }

  html := string(out)

  //"<"でファイルを分割する
  reg := "[<]"
  splitHtml := regexp.MustCompile(reg).Split(html, -1)
  //分割したものから .staticflickr.comが含まれているか確認する
  for _, s := range splitHtml {
    if strings.Contains(s, ".staticflickr.com") {
    //if strings.Contains(s, "html") {
      count++
    }
  }

  //.staticflickr.comが一定個以上あった場合，正常そう
  if(count > 5){
    log.Println(fmt.Sprintf("<Info> id: " + id + ", htmlchek Success: .staticflickr.com num: ", count))
    fmt.Fprintln(logfile, time.Now().Format("2006/01/02 15:04:05") + fmt.Sprintf("<Info> id: " + id + ", htmlchek Success: .staticflickr.com num: ", count))
    return true
  }else{
    log.Println(fmt.Sprintf("<Info> id: " + id + ", htmlchek Failure: .staticflickr.com num: ", count))
    fmt.Fprintln(logfile, time.Now().Format("2006/01/02 15:04:05") + fmt.Sprintf("<Info> id: " + id + ", htmlchek Failure: .staticflickr.com num: ", count))
    return false
  }
}
