# Go包管理之Modules
在1.5引入vendor机制后，go在1.11正式引入完善的包管理机制modules，并希望在1.13中作为稳定特征引入。 当前在1.11和1.12中作为可选机制，默认未开启。

## Go项目文件工程路径
目前我们一个普通的go项目要遵循下面的文件目录结构，包含一个bin文件夹用来存放生成的可执行文件，以及存放源码的src文件夹。在src文件夹中有我们的main文件，以及各种依赖到的依赖包。 
```
-- mypro
  -- bin
  -- src
     -- mypro
       -- main.go
     -- vendor 
       -- deeps....
     -- github.com
       -- deeps....
```
如果我们的项目只是一个单文件项目， 像这样简单的项目我们也必须建立一个标准文件结构目录，因为引入了非官方包geohash_go(假设当前非GOPATH路径)。
```golang
package main

import (
	"fmt"
	"github.com/zhaozy93/geohash_go"
)

func main() {
	hashKey, err := geohash_go.EnGeoHash(float64(39.928167), float64(116.389550), 10)
	lat, lng, err := geohash_go.DeGeoHash(hashKey)
	neigh, err := geohash_go.GetNeighbour(hashKey)
	fmt.Println(hashKey, lat, lng, neigh, err)
}
```
## 尝试Modules
Modules作为更标准的包管理已经解决了上述问题，但我们仍然推荐大家使用标准的文件目录结构管理项目。 我们来进行一次最简单的尝试
### 建立 go.mod
在任意目录下我们使用 go mod init mypro来创建一个项目，项目名称为mypro。
会发现目录下多了一个go.mod文件，文件内容为
```
module mypro

go 1.12
```
### 编写main文件
与上述main文件一致即可
### go build
使用go build来编译文件，注意我们此时并不在GOPATH目录下。
执行go build会得到`go: finding github.com/zhaozy93/geohash_go latest`
会发现go.mod文件多了`require github.com/zhaozy93/geohash_go v0.0.0-20190723023836-e24355b62ab7`，同时多了go.sum文件。
同时发现在GOROOT目录下多了pkg/mod/github.com/zhaozy93/geohash_go@v0.0.0-20190723023836-e24355b62ab7， 说明我们的modules包管理将真正的包依赖下载到了goroot目录下，并没有放在项目文件目录下。

## 理解Modules
刚才的练手算是对modules有了一个最简单的认识，在项目文件中引入import依赖，在执行go run、go build命令时会主动将依赖包下载到goroot的pkg目录下，同时会带有版本号，刚刚的`@v0.0.0-20190723023836-e24355b62ab7`就是版本号，因此可以不同项目依赖不同版本。
#### 指定依赖版本号
但有时我们也希望指定对应的版本号，可以直接修改`go.mod`文件`require foo v1.2.3`或者在文件路径下主动执行`go get foo@v1.2.3`命令即可。
#### 添加本地依赖包
有时候我们可能不止需要外部包，一个项目依赖多个本地包共同构成，因此引入本地依赖包也是一种需求。
之前需要将依赖包放在GOPATH或者GOROOT中，因此我们需要对不同的项目执行频繁的修改GOPATH。
假设文件为main.go的项目需要依赖同级目录的deep1依赖
```
-- mypro
  -- go.mod
  -- main.go
  -- deep1 
    -- deep1.go
```
那么只需要修改go.mod,使用replace相当于为deep1依赖指定路径而不是使用go get从互联网获取。
```
require deep1 v1.0.0

replace deep1 => ./deep1
```

## 更多Q&A
欢迎参考 [Modules · golang/go Wiki · GitHub](https://github.com/golang/go/wiki/Modules#how-do-i-use-vendoring-with-modules-is-vendoring-going-away)，例如正在考虑vendor和module如何兼容使用的你,虽然答案是不可以。