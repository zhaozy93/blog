### 管道
管道就是将一个程序的输出和另一个程序的输入。 当进程创建一个管道的时候，系统内核同时为该进程设立一对文件句柄(流)，一个用来从该管道读取数据(read)，一个用来从管道输入数据(write)。
在Linux系统内核中，每一个管道都是一个inode节点来表示，但是看不到而已。
#### pipe()系统调用
int pipe(int fd[2]) 调用pipe会返回一个参数， 0表示管道创建成功了。
 需要传入两个长度的int作为文件句柄，0表示读、1表示写，这样便可以进行进程间通信了
```c
// fork后子进程会集成父进程的文件描述符，因此可以直接调用
#include <sys/wait.h>
#include <stdio.h>
#include <unistd.h> 
#include <sys/types.h>


int main() {
    int fd[2], nbytes; 
    pid_t childpid;
    pipe(fd);
    char string[] = "Hello, world!\n"; 
    char readbuffer[80];

    if((childpid = fork()) == -1) {
        printf("%s\n", "fork error");
        exit(1);
    }

    if(childpid == 0){
        printf("%s\n", "i am child");
        close(fd[0]);
        write(fd[1], string, strlen(string));
    }else{
        printf("%s\n", "i am father");
        close(fd[1]);
        nbytes = read(fd[0], readbuffer, sizeof(readbuffer)); 
        printf("Received string: %s", readbuffer);
    }

    exit(0);
}
```
#### dup()系统调用
刚刚的pipe只是创建了管道，如果子进程执行了exec，那么创建的文件描述符肯定不会被继承，也就无法进行痛心了，但是标准输入和标准输入被继承了。 也就是说可以让刚刚创建的管道进行重定向到标准IO即可实现与exec创建的新进程进行通信。
int dup(fd[0]) 
创建的子进程执行dup(fd[0])则可以将未使用的输入句柄重定向标准输入，这样父进程只要把通信内容写入标准输出即可完成通信
```c
if(childpid == 0){
        printf("%s\n", "i am child");
        close(fd[0]);
        dup(fd[0]);
        execlp("sort", "sort", NULL);
        // write(fd[1], string, strlen(string));
    }
```
#### dup2()
原来的dup在调用前需要先close掉对应的句柄，在调用dup，如果中间进了一个中断，那就报错了，进程没有stdin。 dup2则将步骤归一了。 
 int dup2( int oldfd, int newfd ); 
```c
// 将fd[0]重定向到0(stdin)
dup2(0, fd[0]); 
execlp(“sort”, “sort”, NULL); 
```
#### 注意事项
- pipe操作在fork之前
- 关闭不需要的句柄
- 管道只能实现父子之间进程的通信
#### popen()/pclose()
就是封装了的fork、exec、dup节省代码，一步完成不被中断。
FILE *popen ( char *command, char *type); 
int pclose( FILE *stream )
优势在于省代码，command可以是任意符合shell命令的指令。
### 有名管道
管道由于只能在父子进程之间通信，虽然是inode节点形式存在，但对外不可见，也就无法完成普遍意义的进程间通信。有名管道以一种特殊的设备文件形式存在于文件系统中，这样不仅具备了管道的通信功能，又有了普通文件的有点(被多个进程共享，同时访问、长期存在)
其实就是创建一个文本，一个在读、一个在写，这样就实现了通信
#### 有名管道创建方式
mknod sampleFIFO p 
mkfifo – m 0666 sampleFIFO 
如果使用mknod方式，之后还要进行chmod权限的设置
例如：
```c
int mknod( char *pathname, mode_t mode, dev_t dev); 
mknod(“/tmp/sampleFIFO”,s_IFIFO|0666,0) 
```
[Linux进程间通信——使用命名管道 - CSDN博客](https://blog.csdn.net/ljianhui/article/details/10202699)
