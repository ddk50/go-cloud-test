### 開発
```bash
$ docker-compose -f docker-compose.dev.yml up
```

#### キャッシュのクリアとか
```bash
$ docker-compose down --rmi all --volumes --remove-orphans
$ docker-compose -f docker-compose.dev.yml build --no-cache
$ docker-compose -f docker-compose.yml up
```

### Google Cloud Runへのデプロイ
こんな感じで

```bash
$ docker build -t gcr.io/pfj-test-434203/go-cloud-run .
$ docker build -t gcr.io/pfj-test-434203/go-cloud-run --no-cache .
$ docker push gcr.io/pfj-test-434203/go-cloud-run
```