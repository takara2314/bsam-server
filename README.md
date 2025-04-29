<div align="center">
<a href="https://github.com/takara2314/bsam-server">
    <img src="./resources/logo.svg" width="128" height="128" alt="logo" />
</a>

# B-SAM Server - 視覚障がい者セーリング支援アプリ サーバー

![Language: Go](https://img.shields.io/badge/Language-Go-00add8?style=for-the-badge&logo=go)
![Package: Gorilla WebSocket](https://img.shields.io/badge/Package-Gorilla%20WebSocket-a1a1a1?style=for-the-badge)
![Framework: Gin](https://img.shields.io/badge/Framework-Gin-0090d1?style=for-the-badge)
![License: GPL-3.0](https://img.shields.io/badge/License-GPL%203.0-bd0000?style=for-the-badge)

</div>

視覚障がいのある方が、セーリング（ヨット競技）をより安全かつ楽しく行えるようにサポートするシステム「B-SAM（Blind Sailing Assist Mie）」のサーバーです。

コース上に設置されたブイに搭載されたスマートフォンから、そのブイの位置情報をリアルタイムで取得します。そして、この位置情報を競技者のスマートフォンに送り続け、常に最新のレース状況を把握できるようにします。これにより、視覚情報に頼ることなく、競技者はレースに集中することができます。

## 関連リポジトリ
[選手用アプリ（メイン）](https://github.com/takara2314/bsam)

[本部用アプリ](https://github.com/takara2314/bsam-admin)

[レースモニター（外部公開用）](https://github.com/takara2314/bsam-web)

## 前提
- Go 1.24.2
- Docker

## ライセンス
このプロジェクトは [GPL-3.0](./LICENSE) ライセンスの下で公開しています。

### 自由な利用と配布
ソフトウェアを自由に使用、修正、配布する権利が保証されています。
### ソースコードの公開
配布時にはソースコードを提供するか、入手方法を明示する必要があります。
### 派生作品の継承
派生作品も同じGPLv3ライセンスで公開しなければなりません（コピーレフト）。
### 特許権の取り扱い
ソフトウェアに含まれる特許の無償利用を認め、貢献者がユーザーに対して特許訴訟を起こすことを禁止しています。
### 商用利用
営利目的での使用や販売が可能ですが、ソースコードの公開や、派生物へのGPLv3適用などの条件を守る必要があります。

## 開発者
[濱口 宝 (Takara Hamaguchi)](https://github.com/takara2314)

<div align="center">
<small>
© 2022 NPO法人セイラビリティ三重
</small>
</div>
