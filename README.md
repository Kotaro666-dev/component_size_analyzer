# 概要

プロジェクト内のコンポーネントのサイズを判断する上で使用する目的で作られています。

> コンポーネントのサイズを決定するのに便利なメトリクスの1つに、コンポーネント内の総ステートメント数（名前空間またはディレクトリ内に含まれるソースファイル内のステートメントの合計）がある。

[「ソフトウェアアーキテクチャ・ハードパーツ」, P.85 より引用](https://www.oreilly.co.jp/books/9784814400065/)

# 使い方

## 1. サンプルの設定ファイルを埋めてください

path: ./config.sample.json

```json
{
  "rootDir": "",
  "statement_character": ""
}
```

- rootDir: 分析したいプロジェクトのルートディレクトリ
- statement_character: 各言語のステートメント終端を指し示す特殊文字
  - 例: C や Dart などではセミコロン、Python や Ruby などでは改行文字

## 2. 設定ファイルを作成してください

```bash
$ make init
cp config.sample.json config.json
```

## 3. プログラムを実行する

```bash
$ make run
```

# 出力結果

path: result.csv

```csv
コンポーネント名, 名前空間, パーセント, 総ステートメント数, ファイル数
componentA, namespace.componentA, 0.11, 9, 1
componentB, namespace.componentB, 0.46, 38, 4
componentC, namespace.componentC, 2.10, 321, 10
componentD, namespace.componentD, 0.75, 62, 1
...
```
