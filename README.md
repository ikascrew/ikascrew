I am ikascrew.
I am a program born to transform "VJ System".

# Overview

"ikascrew" is a software designed to make VJ easy with PC.
VJ can be done just by preparing the personal computer and MP4.

Currently it only supports Linux, but I hope to support Windows and Mac in the future as well.

You can also simplify the operation by preparing "Joy Stick" and "Powermate".

# System Requirements

CPU:Intel core i7
RAM:4GB

# Install

install opencv
use gocv.io/x/gocv


## windows

https://github.com/ikascrew/ikascrew/releases

## Create Project

ikascrew init [project]

ex:
   ikascrew init sample

## Use

ikascrew play [project]

    -> ikascrew server [project]
    -> ikascrew client


# Develop

go get github.com/ikascrew/ikascrew

# 現在の状況

2019年12月31日に大きな変更を開始しています。
現在、joystickとpowermateでコントロールしていますが、
それはクライアントの状況だけにして、クライアントから来るgRPCの通信のみで実行できる形にします。

gRPCを中心にした構成にし、サーバ機で動画を溜め込む仕組みに変更していきたいと思っています。
また、VideoととともにEffectを追加し、動画の切り替えエンジンも見直しを図っていきます。

opencvへのアクセスをgocvに変更したことで変更しやすい部分も出てくると思いますので
頑張っていきたいとは思っています。


プロジェクトはデータベースに設定しておき、サーバは常にそこを見て、
クライアントはそこから情報を取り出すようにしていきます。

実行時にオリジナルのVideo、Effectを追加できるように、一度コンパイルして実行する形式を取ろうと考えています。

   videos {
       "file" : "github.com/ikascrew/video/file",
       "countdown" : "github.com/ikascrew/plugins/countdown"
   } ,
   effects {
       "switch : "github.com/ikascrew/effect/switch"
   } 

って感じで設定？

またそれに伴い、オープニングのムービーはソフト独自の動画を作成する予定です。
