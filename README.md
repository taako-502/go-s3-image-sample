# Go S3 Image Sampe

画面から S3 に画像をアップロードするサンプルコード。
http/template を利用する。

## バケットポリシー

以下のバケットポリシーを設定する。<br>

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowFullAccessToSampleDeveloper", // 任意の文字列
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::{code}:user/{user}"
            },
            "Action": "s3:*",
            "Resource": [
                "arn:aws:s3:::{bucket name}",
                "arn:aws:s3:::{bucket name}/*"
            ]
        }
    ]
}
```
