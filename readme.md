# sjq

JSONとそれ以外を分離するGo製のCLIツールです。

## 使い方

```sh
cat gochamaze.log | sjq
# gochamze.logのJSONとしてパースできる部分のみをsrtdoutに出力する
{ "timestamp": "2023-04-01T12:00:00Z", "message": "Hello, world!" }
{ "timestamp": "2023-04-01T12:00:01Z", "message": "Goodbye, world!" }

cat gochamaze.log | sjq -n
# gochamze.logのJSONとしてはパースできない部分のみをstdoutに出力する
application started

cat gochamaze.log | sjq -j json.log -n non-json.log
# json.logにJSONとしてパースできる部分のみをstdoutに出力する
# non-json.logにJSONとしてはパースできない部分のみをstdoutに出力する
```

gochamaze.logは以下のようなテキストファイルです。

```
{ "timestamp": "2023-04-01T12:00:00Z", "message": "Hello, world!" }
application started
{ "timestamp": "2023-04-01T12:00:01Z", "message": "Goodbye, world!" }
```
