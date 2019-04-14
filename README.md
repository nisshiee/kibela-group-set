# kibela-group-set

[Kibela](https://kibe.la/ja)のAPIを利用して、特定のフォルダ内の全記事を指定したグループに公開するという機能を提供します。

# 使い方

1. 変数の用意

    ```bash
    curl -sSL -o kibela-group-set.env https://raw.githubusercontent.com/nisshiee/kibela-group-set/master/.env.sample
    vi kibela-group-set.env
    ```

    ※ 各変数の詳細は後述
    
2. Docker Containerを実行

    ```bash
    docker run --env-file kibela-group-set.env nisshiee/kibela-group-set
    ```

## 変数について

- `KIBELA_TEAM`

    kibela-group-setを利用するKibelaのチームを指定します。
    
    `https://${KIBELA_TEAM}.kibe.la/` ←KibelaURLのこの部分に一致します。
    
- `KIBELA_TOKEN`

    Kibela Web APIの認証に利用するトークンです。以下のURLから取得してください。
    
    `https://${KIBELA_TEAM}.kibe.la/settings/access_tokens`
    
    `read`、`write`の両権限が必要です。

- `KIBELA_TARGET_FOLDER_ID`

    この変数で指定したフォルダ以下の全記事が対象になります。子フォルダ以下も再帰的に対象にします。
    
    このIDはAPIで使用する専用のIDを指定する必要があります。Web API Consoleを開き、以下のクエリを実行してIDを取得してください。
    
    Web API Console: `https://${KIBELA_TEAM}.kibe.la/api/console`
    
    ```
    query {
      folders(first: 100) {
        nodes {
          id
          fullName
        }
      }
    }
    ```

- `KIBELA_TARGET_GROUP_ID`

    この変数で指定したグループに向けて、対象の記事を公開します。
    
    このIDはAPIで使用する専用のIDを指定する必要があります。Web API Consoleを開き、以下のクエリを実行してIDを取得してください。
    
    Web API Console: `https://${KIBELA_TEAM}.kibe.la/api/console`

    ```
    query {
      groups(first: 100) {
        nodes {
          id
          name
        }
      }
    }
    ```
    
# 開発のしかた

- Go 1.12以降を想定しています
- Go Modulesを使用しています
    - `export GO111MODULE=on` が必要かもしれません
    - 詳細は[こちら](https://github.com/golang/go/wiki/Modules)
    
```bash
git clone https://github.com/nisshiee/kibela-group-set.git
cd kibela-group-set
go build *.go
```
