from faster_whisper import WhisperModel
import sys

initial_prompt = """そうだ。今日はピクニックしない？天気もいいし、絶好のピクニック日和だと思う。いいですね。
では、準備をはじめましょうか。そうしよう！どこに行く？そうですね。三ツ池公園なんか良いんじゃないかな。
今の時期なら桜が綺麗だしね。じゃあそれで決まり！わかりました。電車だと550円掛かるみたいです。
少し時間が掛かりますが、歩いた方が健康的かもしれません。
"""

model_size = sys.argv[1]

# Run on CPU with INT8
model = WhisperModel(model_size, device="cpu", compute_type="int8")

segments, info = model.transcribe(sys.argv[2], beam_size=5, language="ja", initial_prompt=initial_prompt, vad_filter=True)

with open(sys.argv[3], "w", encoding="utf-8") as f:
    for segment in segments:
        # print("[%.2fs -> %.2fs] %s" % (segment.start, segment.end, segment.text))
        print(segment.text, file=f)
