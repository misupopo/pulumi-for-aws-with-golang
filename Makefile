bucketName = pulumi-state-object
currentTime = $(shell date +'%Y%m%d_%H%M%S')

# 例 20210416_110811
echo-time:
	echo $(shell date +'%Y%m%d_%H%M%S')

# state fileをエクスポートする
export-state:
	pulumi stack export > ./state/stack_$(currentTime).json

import-state:
	pulumi stack export > ./state/stack_$(currentTime).json

pulumi-state-upload:


# s3のバケットリストを表示する
get-s3-buckets:
	aws s3 ls

# ターゲットのs3フォルダと同期を行う
sync-s3-buckets:
	aws s3 sync ./state s3://$(bucketName) --exclude ".gitkeep"
