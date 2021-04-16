bucketName = pulumi-state-object
currentTime = $(shell date +'%Y%m%d_%H%M%S')

define test_export_state
   make export-state
endef

# 例 20210416_110811
echo-time:
	echo $(shell date +'%Y%m%d_%H%M%S')

call-function:
	@$(call test_export_state)

# state fileをエクスポートする
export-state:
	pulumi stack export > ./state/stack_$(currentTime).json

# 例 import-state fileName=hogehoge.json
import-state:
	pulumi stack export > ./state/${fileName}.json

# s3のバケットリストを表示する
get-s3-buckets:
	aws s3 ls

# ターゲットのs3フォルダと同期を行う
sync-s3-buckets:
	aws s3 sync ./state s3://$(bucketName) --exclude ".gitkeep"

# pulumiのupを行った後にstateファイルをexportする
pulumi-up:
	pulumi up
	make export-state

# stack環境を選択する
pulumi-stacks:
	pulumi stack ls

