## デプロイ

```bash
$ gcloud functions deploy Messages --runtime go113 --trigger-http --allow-unauthenticated
```
### API のテスト

```bash
$ curl -X GET https://us-central1-glossom-bulletin-board-sample.cloudfunctions.net/Messages
[{"name":"ヤンマガ読者","text":"漫画は面白いです。","timestamp":"2021-03-24 21:00:00"},{"name":"Glossom社員","text":"そうですね。","timestamp":"2021-03-24 21:00:01"}]
$ curl -X POST \
 	-d "name=ヤンマガチーム" -d "text=これからも楽しみにしていてくださいね。" \
 	-w '%{http_code}' \
 	https://us-central1-glossom-bulletin-board-sample.cloudfunctions.net/Messages
201
$ curl -X PUT -w '%{http_code}' https://us-central1-glossom-bulletin-board-sample.cloudfunctions.net/Messages
405
```

## 情報の確認

```bash
$ gcloud functions describe Message
```