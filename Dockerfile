ARG USER_NAME
ARG USER_ID
ARG USER_GID
ARG PASSWORD

FROM node:18-alpine as app

WORKDIR /workspace/app
COPY ./app ./

RUN npm i -g npm && \
    npm i && \
    npm run export

FROM golang:1.19-alpine as api

WORKDIR /workspace/api/

COPY ./api ./
COPY --from=app /workspace/app/out ./out

RUN apk update && apk add build-base && \
  go mod tidy && \
  go build -o main main.go

FROM alpine:3.18

ARG USER_NAME
ARG USER_ID
ARG USER_GID
ARG PASSWORD

RUN addgroup -g ${USER_GID} -S ${USER_NAME} && \
    adduser  -u ${USER_ID} -G ${USER_NAME}  -s /bin/sh -D ${USER_NAME} &&  \
    echo ${USER_NAME}:${PASSWORD} | chpasswd;

WORKDIR /application

COPY --from=api /workspace/api/main ./

RUN chown -R ${USER_ID}:${USER_GID} ./

USER ${USER_ID}:${USER_GID}

CMD [ "./main" ]