## デプロイ

```bash
$ gcloud functions deploy Messages --runtime go113 --trigger-http --allow-unauthenticated
```

## ローカルでのテスト

```bash
$ curl -X GET https://(project name).cloudfunctions.net/Messages
Hello, World!
$ curl -X PUT https://(project name).cloudfunctions.net/Messages
403 - Forbidden
$ curl -X POST https://(project name).cloudfunctions.net/Messages
405 - Method Not Allowed
```

## 情報の確認

```bash
$ gcloud functions describe Message
```