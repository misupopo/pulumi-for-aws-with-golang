# pulumi-for-aws-with-golang

## キーペアの作成
configの「keyPairName」は予めAWSに上げておき、
かつ作成のタイミングでキーペアがダウンロードできるので、.sshにダウンロードしておく。
そのあとで、pulumi upで構築。__

ec2インスタンスはそのキーペアに対して紐づく。

ec2インスタンスへのssh ログインはElastic IPアドレスに紐づくので、そのグローバルIPでログインができるようになる。


