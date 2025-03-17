import json
import os
import winsound
# 環境変数 WORKSPACE_FOLDER の値を取得
workspace_folder = os.environ.get('WORKSPACE_FOLDER')

# 値を表示
print(f"workspace_folder: {workspace_folder}")

# Read app_config.json file
with open(f'{workspace_folder}/cmd/app/app_config.json', 'r', encoding='utf-8') as file:
    config_data = json.load(file)

# Convert JSON data to dictionary
config_dict = dict(config_data)

app_name = config_dict.get('Name')
app_version = config_dict.get('Version')

print(f"app_name: {app_name}")
print(f"app_version: {app_version}")

all_rebuild = "-a" if os.environ.get('ENV') == 'prod' else ""
env_name = "dev" if os.environ.get('ENV') == 'dev' else "prod" if os.environ.get('ENV') == 'prod' else "stg"

# Build command
# -o 出力フォルダ
# -trimpath ビルドパスを削除
# -v ビルドログを出力
# -a 全ての依存関係を再ビルド
# -buildmode=exe 実行可能ファイルを生成
# -ldflags "-s -w" バイナリサイズを小さくする
# -H=windowsgui コンソールを表示しない
# -linkmode external -extldflags '-static -Wl,cmd/app/app.res' リソースを埋め込む
build_command = f"go build -o {workspace_folder}/build/{app_name}_{app_version}.exe -trimpath " \
                f"-v {all_rebuild} -buildmode=exe -ldflags \"-s -w -H=windowsgui -X main.env={env_name} " \
                f"-linkmode external -extldflags '-static -Wl,{workspace_folder}/cmd/app/app.res'\" "

print(f"build_command: {build_command}")

os.system(build_command)

# Play beep sound
winsound.PlaySound("SystemAsterisk", winsound.SND_ALIAS)
