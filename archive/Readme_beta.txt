-------------------------------------------------
-------------------------------------------------

　　MuModelViewer

　　ベータ版用追加Readme

　　　　　　　　　　　　　　　　　miu200521358

-------------------------------------------------
-------------------------------------------------

----------------------------------------------------------------
■　配布元
----------------------------------------------------------------

　・Discord「miuの実験室」
　　　https://discord.gg/MW2Bn47aCN
　・Discord「MMDの集い　Discord支部」
　　　https://discord.gg/wBcwFreHJ8

　※基本的にベータ版は上記二箇所でのみ配布しておりますので、上記以外で見かけたらお手数ですがご連絡下さい。


----------------------------------------------------------------
■　β版をご利用いただくにあたってのお願い
----------------------------------------------------------------

・この度はβ版テスターにご参加いただき、ありがとうございます
　β版をご利用いただくのにあたって、下記点をお願いいたします。

・不具合報告、改善要望、大歓迎です。
　要望についてはお応えできるかは分かりませんが…

・不具合報告の場合、下記をご報告ください
　・ツール名
　・ベータ版の番号
　・エラー発生時にダイアログが表示されている場合、そこのエラーメッセージ
　　・ダイアログが表示された時点でクリップボードにエラーメッセージがコピーされています
　・一般配布されているモデルの場合、お迎えできるURL

　・miuの実験室 - 00_質問・相談・報告
　　https://discord.gg/MW2Bn47aCN
　　・カテゴリに参加する必要がありますので、「カテゴリ参加申請」から申請を出してください


・ベータ版の扱いは、リリース版と同様にお願いします。
　自作発言とツールの再配布だけNG。

・ベータ版を使ってみて良かったら、ぜひ公開して、宣伝してください！
　励みになります！公開先はどこでもOKです。
　その際に、Twitterアカウント（@miu200521358）を添えていただけたら、喜んで拝見に伺います


----------------------------------------------------------------
■　履歴
----------------------------------------------------------------

MuModelViewer_1.1.0_beta_04 (2025/03/25)
- 基幹ライブラリ修正
  - VSync無効設定 (デバイス側の制御は非対応)
  - カメラ同期とオーバーレイを排他に
  - 初期FPS制限を 30fps に変更
  - 回転付与の計算を拡張版Slerpを使うよう修正
  - Buffer 系の削除処理のエラーハンドリング追加
  - 物理デバッグ時のメモリリーク修正
  - モデルとモーションのパスを変更したタイミングで、モデルの再読み込みを行うよう修正
  - 文字列のコピペなどでもモデルの読み込みを行うよう修正

MuModelViewer_1.1.0_beta_03 (2025/03/24)
- 基幹ライブラリ修正
    - ボーン表示がバグっていたのを修正
    - 物理リセット時のデフォーム計算がミスってたので修正
    - 再生ボタン押下時のボタン文言を切り替えるよう修正
    - カメラ同期機能復活
    - オーバーレイ表示機能復活
        - カメラを自動的にそれっぽくメインウィンドウに合わせるようにしてみた

MuModelViewer_1.1.0_beta_02 (2025/03/18)
- 基幹ライブラリ修正
    - 頂点が割り当てられていない材質がある場合に、エラーが発生する問題を修正
    - テクスチャ読み込み失敗時に、エラーログの出力を追加
    - 日本語以外の言語設定時に、不足文言があるとクラッシュする問題を修正
- 材質リストを材質テーブルビューに変更
    - テクスチャなどが読み込めないものであるかどうかが分かるように表示情報追加

MuModelViewer_1.1.0_beta_01 (2025/03/17)
- 基幹ライブラリ作り直し
    - 速度・安定度UP
    - コントローラーウィンドウとビューワーウィンドウの移動や最前面、最小化などを連動させるよう機能追加
    - 移動付与ボーンの移動がうまくいってなかったのを修正
- 画像データの読み込み拡張
    - 読み込みに失敗した場合、他の拡張子でもデコードを試す
- スフィアファイルの対応拡張子に sph を追加
- 圧縮バイナリ形式のXファイルの読み込みに対応
- 材質の表示ON OFF機能追加

MuModelViewer_1.0.0_beta_02 (2024/12/15)
- アイコン調整
- Xファイル読み取り追加

MuModelViewer_1.0.0_beta_01 (2024/12/13)
- ベータ版公開

