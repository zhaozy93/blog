# Go 包管理GOPATH、Vendor
今天入职了新公司，看了golang一个项目的README.md，里面介绍了vendor的使用规范以及govendor的使用。 之前知道goroot、gopath，对vendor了解，但没有仔细研究，临时恶补一下。

## cannot find package
src/myapp/app.go:4:8: cannot find package "deep" in any of:
	/usr/local/go/src/deep (from $GOROOT)
	/Users/zyzhao/Desktop/pro/src/deep (from $GOPATH)
这是一种很常见的报错，表明在编译时缺少包依赖，因此大家对GOROOT和GOPATH比较熟悉了。
## GOROOT
最简单，其实就是golang的安装路径，如果去/usr/local/go 里面会发现我们的标准包都在里面。
## GOPATH
既然GOROOT是golang的安装路径，相当于官方仓库，那对于每一个项目应该有一个工作仓库，那GOPATH其实就相当于我们每个项目的工作仓库。 因此我们经常看到在make时重写 
``` shell
export GOPATH=$GOPATH:`pwd`
```
GOPATH和bin路径有些类似，也是使用:来连接多个文件夹。
## Vendor
Vendor是今天的主角，可以把vendor看作是单个项目的专有仓库，毕竟GOPATH还是可以贡献的，所有build项目都会以GOPATH来寻找依赖包，具有共享性。 但vendor是与项目文件平级的，只能被当前和下级项目使用，所以专有性更强。 关于vendor包的专有性(为什么只有本项目和下级项目可以使用)后面会介绍道。

## 依赖寻找优先级
既然有三个依赖包的存放位置，那么肯定要有一个生效优先级顺序的问题。
Vendor > GOROOT > GOPATH

### 优先级顺序证明
建一个errors的文件夹，包含文件errors.go
``` go
package errors

import "fmt"

const ERRORS = 1 

func init(){
  fmt.Println("error deep")
}
```
在main函数中增加` _ = errors.ERRORS`
#### Step 1 Vendor > GOROOT
当把包放在vendor中，编译无误。 证明vendor中的errors包对标准包即GOROOT中的errors产生了遮蔽效果，即优先级 Vendor > GOROOT
#### Step 2 GOROOT > GOPATH
当把包放在GOPATH中，编译报错 `undefined: errors.ERRORS``。 证明GOPATH中的errors包对标准包即GOROOT中的errors没有产生了遮蔽效果，因为标准包中没有ERRORS这个输出变量。即优先级  GOROOT > GOPATH

## Vendor内部优先级顺序
先来看一组目录
```
----mypro
   ---- main.go
   ---- vendor
      ----- deep1
         ------deep1.go
         ------vendor
            ------deepdeep1
               ------deepdeep1.go
      ----- deep2
         -------deep2.go
         --------vendor
             -------deepdeep2
                 ---------deepdeep2.go
```

梳理一下 各个项目的pwd
```
mypro: ./mypro
deep1: ./mypro/vendor/deep1
deepdeep1: ./mypro/vendor/deep1/vendor/deepdeep1
deep2: ./mypro/vendor/deep2
deepdeep2: ./mypro/vendor/deep2/vendor/deepdeep2
```
deepdeep1项目中的依赖寻找路径是这样的
> ./mypro/vendor/deep1/vendor/deepdeep1/vendor   
> ./mypro/vendor/deep1/vendor  
> ./mypro/vendor/  
> GOROOT  
> GOPATH  
通过这个vendor的寻找优先级可以发现寻找vendor包只会向上级寻找，直到顶层，再寻找GOROOT和GOPTATH。
因此对于同一依赖项目，但是依赖不同版本，可以使用vendor来解决，但是不要这样做，当依赖特定版本的依赖时，一定是项目本身有问题。

