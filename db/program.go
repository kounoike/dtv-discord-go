package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	genreMap = []string{
		"ニュース・報道",     // 0x0
		"スポーツ",        // 0x1
		"情報・ワイドショー",   // 0x2
		"ドラマ",         // 0x3
		"音楽",          // 0x4
		"バラエティ",       // 0x5
		"映画",          // 0x6
		"アニメ・特撮",      // 0x7
		"ドキュメンタリー・教養", // 0x8
		"劇場・公演",       // 0x9
		"趣味・教育",       // 0xa
		"福祉",          // 0xb
		"予備",          // 0xc
		"予備",          // 0xd
		"拡張",          // 0xe
		"その他",         // 0xf
	}
	subGenreMap = [][]string{
		[]string{ // 0x0
			"定時・総合",     // 0x0
			"天気",        // 0x1
			"特集・ドキュメント", // 0x2
			"政治・国会",     // 0x3
			"経済・市況",     // 0x4
			"海外・国際",     // 0x5
			"解説",        // 0x6
			"討論・会談",     // 0x7
			"報道特番",      // 0x8
			"ローカル・地域",   // 0x9
			"交通",        // 0xa
			"",          // 0xb
			"",          // 0xc
			"",          // 0xd
			"",          // 0xe
			"その他",       // 0xf
		},
		[]string{ // 0x1
			"スポーツニュース",      // 0x0
			"野球",            // 0x1
			"サッカー",          // 0x2
			"ゴルフ",           // 0x3
			"その他の球技",        // 0x4
			"相撲・格闘技",        // 0x5
			"オリンピック・国際大会",   // 0x6
			"マラソン・陸上・水泳",    // 0x7
			"モータースポーツ",      // 0x8
			"マリン・ウィンタースポーツ", // 0x9
			"競馬・公営競技",       // 0xa
			"",              // 0xb
			"",              // 0xc
			"",              // 0xd
			"",              // 0xe
			"その他",           // 0xf
		},
		[]string{ // 0x2
			"芸能・ワイドショー", // 0x0
			"ファッション",    // 0x1
			"暮らし・住まい",   // 0x2
			"健康・医療",     // 0x3
			"ショッピング・通販", // 0x4
			"グルメ・料理",    // 0x5
			"イベント",      // 0x6
			"番組紹介・お知らせ", // 0x7
			"",          // 0x8
			"",          // 0x9
			"",          // 0xa
			"",          // 0xb
			"",          // 0xc
			"",          // 0xd
			"",          // 0xe
			"その他",       // 0xf
		},
		[]string{ // 0x3
			"国内ドラマ", // 0x0
			"海外ドラマ", // 0x1
			"時代劇",   // 0x2
			"",      // 0x3
			"",      // 0x4
			"",      // 0x5
			"",      // 0x6
			"",      // 0x7
			"",      // 0x8
			"",      // 0x9
			"",      // 0xa
			"",      // 0xb
			"",      // 0xc
			"",      // 0xd
			"",      // 0xe
			"その他",   // 0xf
		},
		[]string{ // 0x4
			"国内ロック・ポップス",      // 0x0
			"海外ロック・ポップス",      // 0x1
			"クラシック・オペラ",       // 0x2
			"ジャズ・フュージョン",      // 0x3
			"歌謡曲・演歌",          // 0x4
			"ライブ・コンサート",       // 0x5
			"ランキング・リクエスト",     // 0x6
			"カラオケ・のど自慢",       // 0x7
			"民謡・邦楽",           // 0x8
			"童謡・キッズ",          // 0x9
			"民族音楽・ワールドミュージック", // 0xa
			"",    // 0xb
			"",    // 0xc
			"",    // 0xd
			"",    // 0xe
			"その他", // 0xf
		},
		[]string{ // 0x5
			"クイズ",      // 0x0
			"ゲーム",      // 0x1
			"トークバラエティ", // 0x2
			"お笑い・コメディ", // 0x3
			"音楽バラエティ",  // 0x4
			"旅バラエティ",   // 0x5
			"料理バラエティ",  // 0x6
			"",         // 0x7
			"",         // 0x8
			"",         // 0x9
			"",         // 0xa
			"",         // 0xb
			"",         // 0xc
			"",         // 0xd
			"",         // 0xe
			"その他",      // 0xf
		},
		[]string{ // 0x6
			"洋画",  // 0x0
			"邦画",  // 0x1
			"アニメ", // 0x2
			"",    // 0x3
			"",    // 0x4
			"",    // 0x5
			"",    // 0x6
			"",    // 0x7
			"",    // 0x8
			"",    // 0x9
			"",    // 0xa
			"",    // 0xb
			"",    // 0xc
			"",    // 0xd
			"",    // 0xe
			"その他", // 0xf
		},
		[]string{ // 0x7
			"国内アニメ", // 0x0
			"海外アニメ", // 0x1
			"特撮",    // 0x2
			"",      // 0x3
			"",      // 0x4
			"",      // 0x5
			"",      // 0x6
			"",      // 0x7
			"",      // 0x8
			"",      // 0x9
			"",      // 0xa
			"",      // 0xb
			"",      // 0xc
			"",      // 0xd
			"",      // 0xe
			"その他",   // 0xf
		},
		[]string{ // 0x8
			"社会・時事",      // 0x0
			"歴史・紀行",      // 0x1
			"自然・動物・環境",   // 0x2
			"宇宙・科学・医学",   // 0x3
			"カルチャー・伝統文化", // 0x4
			"文学・文芸",      // 0x5
			"スポーツ",       // 0x6
			"ドキュメンタリー全般", // 0x7
			"インタビュー・討論",  // 0x8
			"",           // 0x9
			"",           // 0xa
			"",           // 0xb
			"",           // 0xc
			"",           // 0xd
			"",           // 0xe
			"その他",        // 0xf
		},
		[]string{ // 0x9
			"現代劇・新劇",  // 0x0
			"ミュージカル",  // 0x1
			"ダンス・バレエ", // 0x2
			"落語・演芸",   // 0x3
			"歌舞伎・古典",  // 0x4
			"",        // 0x5
			"",        // 0x6
			"",        // 0x7
			"",        // 0x8
			"",        // 0x9
			"",        // 0xa
			"",        // 0xb
			"",        // 0xc
			"",        // 0xd
			"",        // 0xe
			"その他",     // 0xf
		},
		[]string{ // 0xa
			"旅・釣り・アウトドア",   // 0x0
			"園芸・ペット・手芸",    // 0x1
			"音楽・美術・工芸",     // 0x2
			"囲碁・将棋",        // 0x3
			"麻雀・パチンコ",      // 0x4
			"車・オートバイ",      // 0x5
			"コンピュータ・ＴＶゲーム", // 0x6
			"会話・語学",        // 0x7
			"幼児・小学生",       // 0x8
			"中学生・高校生",      // 0x9
			"大学生・受験",       // 0xa
			"生涯教育・資格",      // 0xb
			"教育問題",         // 0xc
			"",             // 0xd
			"",             // 0xe
			"その他",          // 0xf
		},
		[]string{ // 0xb
			"高齢者",    // 0x0
			"障害者",    // 0x1
			"社会福祉",   // 0x2
			"ボランティア", // 0x3
			"手話",     // 0x4
			"文字(字幕)", // 0x5
			"音声解説",   // 0x6
			"",       // 0x7
			"",       // 0x8
			"",       // 0x9
			"",       // 0xa
			"",       // 0xb
			"",       // 0xc
			"",       // 0xd
			"",       // 0xe
			"その他",    // 0xf
		},
		[]string{ // 0xc
			"", // 0x0
			"", // 0x1
			"", // 0x2
			"", // 0x3
			"", // 0x4
			"", // 0x5
			"", // 0x6
			"", // 0x7
			"", // 0x8
			"", // 0x9
			"", // 0xa
			"", // 0xb
			"", // 0xc
			"", // 0xd
			"", // 0xe
			"", // 0xf
		},
		[]string{ // 0xd
			"", // 0x0
			"", // 0x1
			"", // 0x2
			"", // 0x3
			"", // 0x4
			"", // 0x5
			"", // 0x6
			"", // 0x7
			"", // 0x8
			"", // 0x9
			"", // 0xa
			"", // 0xb
			"", // 0xc
			"", // 0xd
			"", // 0xe
			"", // 0xf
		},
		[]string{ // 0xe
			"BS/地上デジタル放送用番組付属情報", // 0x0
			"広帯域 CS デジタル放送用拡張",   // 0x1
			"",             // 0x2
			"サーバー型番組付属情報",  // 0x3
			"IP 放送用番組付属情報", // 0x4
			"",             // 0x5
			"",             // 0x6
			"",             // 0x7
			"",             // 0x8
			"",             // 0x9
			"",             // 0xa
			"",             // 0xb
			"",             // 0xc
			"",             // 0xd
			"",             // 0xe
			"",             // 0xf
		},
		[]string{ // 0xf
			"",    // 0x0
			"",    // 0x1
			"",    // 0x2
			"",    // 0x3
			"",    // 0x4
			"",    // 0x5
			"",    // 0x6
			"",    // 0x7
			"",    // 0x8
			"",    // 0x9
			"",    // 0xa
			"",    // 0xb
			"",    // 0xc
			"",    // 0xd
			"",    // 0xe
			"その他", // 0xf
		},
	}
)

func (p *Program) UnmarshalJSON(b []byte) error {
	type program Program
	var pp program
	err := json.Unmarshal(b, &pp)
	if err != nil {
		return err
	}

	*p = (Program)(pp)
	p.Json = b
	genres := jsoniter.Get(b, "genres")
	if genres.Size() > 0 {
		genreAny := genres.Get(0)
		lv1 := genreAny.Get("lv1").ToInt()
		lv2 := genreAny.Get("lv2").ToInt()
		p.Genre = fmt.Sprintf("%s-%s", genreMap[lv1], subGenreMap[lv1][lv2])
	}

	return nil
}

func (q *Queries) InsertProgram(ctx context.Context, p Program) error {
	args := createProgramParams{
		ID:          p.ID,
		Json:        p.Json,
		EventID:     p.EventID,
		ServiceID:   p.ServiceID,
		NetworkID:   p.NetworkID,
		StartAt:     p.StartAt,
		Duration:    p.Duration,
		IsFree:      p.IsFree,
		Name:        p.Name,
		Description: p.Description,
		Genre:       p.Genre,
	}
	return q.createProgram(ctx, args)
}

func (q *Queries) UpdateProgram(ctx context.Context, p Program) error {
	args := updateProgramParams{
		ID:          p.ID,
		Json:        p.Json,
		EventID:     p.EventID,
		ServiceID:   p.ServiceID,
		NetworkID:   p.NetworkID,
		StartAt:     p.StartAt,
		Duration:    p.Duration,
		IsFree:      p.IsFree,
		Name:        p.Name,
		Description: p.Description,
	}
	return q.updateProgram(ctx, args)
}

func (p *Program) StartTime() time.Time {
	return time.Unix(p.StartAt/1000, (p.StartAt%1000)*1000)
}
