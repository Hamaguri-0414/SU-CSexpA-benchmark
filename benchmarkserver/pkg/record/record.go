package record

import(
  "os"
  "encoding/csv"
  "io"
  "strconv"
  "log"
  "os/exec"

)

func Record(times string, groupName string) string {

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
