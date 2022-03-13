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
)

//abコマンドで負荷をかけ，計測結果を返す
func Ab(id string, url string) (string, string) {

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

  //http://192.168.1.101/~username/directory/progC.php?tag=fiat
  //-c -nを変更する
  out, err := exec.Command("ab", "-c", "1", "-n", "1", url + "?tag=" + randomTag).Output()

  if err != nil {
    //urlが不明
    log.Println(fmt.Sprintf("<Error> execCmd(ab -c 1 -n 1 " + url + "?tag=" + randomTag + ")" , err))
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

      //float64に変換して加算
      measureTime, _ := strconv.ParseFloat(ss[len(ss) - 3], 64)
      measureTimes += measureTime
      break
    }
  }

  //curlでhtmlを取得し，<img>の数が100個あるか数える

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
