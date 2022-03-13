package ab

import (
  "io/ioutil"
  "math/rand"
  "strings"
  "time"
  "log"
  "os/exec"
  "regexp"
  "strconv"
)

//abコマンドで負荷をかけ，計測結果を返す
func Ab(url string) string {

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
