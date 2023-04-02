package discord

const (
	InformationCategory    = "録画-情報"
	LogChannel             = "動作ログ"
	UpdateChannel          = "更新情報"
	RecordingChannel       = "録画開始・完了"
	RecordingFailedChannel = "録画失敗"

	ProgramInformationCategory = "録画-番組情報"
	// サービス名でチャンネルが出来る

	NotifyAndScheduleCategory = "録画-通知・予約"
	AutoActionChannelName     = "自動検索"

	AutoActionChannelTopic          = `**このチャンネルにスレッドを投稿するとEPG更新時ルールに従い自動的に通知または録画されます。**`
	AutoActionChannelWelcomeMessage = `**このチャンネルにスレッドを投稿するとEPG更新時ルールに従い自動的に通知または録画されます。**

**投稿形式**
タイトル: 何でも構いません。
本文: 以下の形式で投稿します。指定の:は半角コロンを使ってください。全角半角、大文字小文字は同一視して検索されます。=の前後の空白は無視されます。
` + "```" + `
タイトル = ニュース
チャンネル= NHK
ジャンル=ニュース
` + "```" + `
（正確にはチャンネルというよりサービスの名前です）

**自動検索結果**
上記のルールにマッチする番組情報が追加されたとき、投稿にぶら下がる発言として、タイトル、番組情報発言へのリンクが投稿されます。
（この発言に:red_circle:でリアクションしても録画されません。リンク先の発言に:red_circle:を付けてください。）

**リアクション**
投稿の最初の発言にリアクションしておくと通知または自動録画が行われます。
:eyes: でリアクションしておくと自動検索でマッチした番組情報の投稿時にメンションが付きます。
:red_circle: でリアクションしておくと自動検索でマッチした番組情報の投稿時に録画が予約されます。
`
)
