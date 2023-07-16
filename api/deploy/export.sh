# 引数の値を取得
arg=$1

ENV_FILE="./deploy/.env"

# use env
if [ "$arg" = "dev" ]; then
    ENV_FILE="./deploy/.env.development"
fi
if [[ -f $ENV_FILE ]]; then
  export $(cat $ENV_FILE | xargs)
fi

# nginx.test.conf ファイルを読み込む
config=`cat ../nginx.test.conf`
config=${config/server_name localhost;/server_name $HOST;}

config=${config/listen 443 ssl;/listen 80;}
config=${config/proxy_pass http:\/\/application:8080\/;/proxy_pass http:\/\/$PROXY_PASS\/;}
config=${config/ssl_certificate \/etc\/certs\/localhost.pem;/}
config=${config/ssl_certificate_key \/etc\/certs\/localhost-key.pem;/}

config=${config/error_log \/etc\/nginx\/conf.d\/log\/error.log;/error_log \/etc\/nginx\/conf.d\/error.log;}


echo "$config"
# 書き換えた設定をファイルに出力
echo "$config" > ./deploy/nginx.conf

echo "nginx.conf ファイルを作成しました。"

# nginx.conf binary export
scp -i deploy/dev/web.pem deploy/nginx.conf ubuntu@18.211.82.71:/home/ubuntu/


# go build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o deploy/main main.go

# go binary export
# scp -i deploy/dev/web.pem deploy/main ubuntu@18.211.82.71:/home/ubuntu/
