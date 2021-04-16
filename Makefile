
# state fileをエクスポートする
export-state:
	pulumi stack export > ./state/stack.json

import-state:
	pulumi stack export > ./state/stack.json


