//HTMLが読み込まれたとき
$(document).ready(function(){
  //計測開始ボタンクリックアクション
	$('#startMeasureBtn').on('click', function(){

    //inputタグの入力
    console.log($('input[name="url"]').val())
    //selectタグの入力
    console.log($('[name="groupName"] option:selected').val())

    //入力フォームを非表示にし，計測中を表示
    $('#topPage').toggle();
    $('#startedMeasure').toggle();

    //ajax urlとgroupNameを/measureに送る
		$.ajax({
			type: "POST",
      //送信先URL
			url: "measure",
			data: {
        //送信データ
				"url": $('input[name="url"]').val(),
				"groupName": $('[name="groupName"] option:selected').val(),
			},
      //サーバから受け取るデータ(data)の形式
			dataType: "json",
      //受け取り成功時
			success: function(data){
				//計測中を非表示にし，計測結果を表示する
				console.log(data.Time)
				console.log(data.Msg)
				setTimeout(function(){
					$('#startedMeasure').toggle();
					$('#MeasureTime').text('Requests per second:' + data.Time)
					$('#Msg').text(data.Msg)
					$('#measureResult').toggle();
				}, 3000);

			}
		});

	});

	//結果画面にあるトップへボタンを押したとき
	//非表示・表示を切り替える
	$('#restartBtn').on('click', function(){
		$('#measureResult').toggle();
		$('#topPage').toggle();
	});

});
