package record

import(
  "os"
  "encoding/csv"
  "io"
  "strconv"
  "log"
  "os/exec"
  "fmt"
)

func Record(id string, times string, groupName string) string {

  msg := "" //返すメッセージ
  recordData := "" //書き込みデータ
  doUpdate := false //記録が更新したかどうか

  //data.csvに記録する
  //data.csvを読み込む
  csvFile, err := os.Open("../public/score.csv")
  if err != nil {
    log.Println(fmt.Sprintf("<Debug> ", err))
  }
  reader := csv.NewReader(csvFile)

  //groupNameの一致を探し，数値を比較する
  for {
    line, err := reader.Read()
    if err == io.EOF {
        break
    }
    //同時に書き込みデータを作成する
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
  file, err := os.Create("../public/score.csv")
  if err != nil{
    log.Println(fmt.Sprintf("<Debug> ", err))
  }
  defer file.Close()
  _, err = file.WriteString(recordData)
  if err != nil {
    log.Println(fmt.Sprintf("<Debug> ", err))
  }

  //csvファイルをgithubにpush
  if doUpdate {
    csvPush(id, groupName)
  }

  log.Println("<Info> id: " + id + ", record msg: " + msg)
  return msg

}

func csvPush(id string, groupName string){
  var err error
  //git add ../exp1_ranking/public/score.csv
  err = exec.Command("git", "add", "../public/score.csv").Run()
  if err != nil {
    log.Println(fmt.Sprintf("<Debug> ", err))
  }
  err = exec.Command("git", "commit", "-m", groupName + "の記録更新").Run()
  if err != nil {
    log.Println(fmt.Sprintf("<Debug> ", err))
  }
  err = exec.Command("git", "push").Run()
  if err != nil {
    log.Println(fmt.Sprintf("<Debug> ", err))
  }

  log.Println("<Info> id: " + id + ",git push new record")

}
