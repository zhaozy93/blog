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
execlp("sort", "sort", NULL); 
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
mknod("/tmp/sampleFIFO",s_IFIFO|0666,0) 
```
[Linux进程间通信——使用命名管道 - CSDN博客](https://blog.csdn.net/ljianhui/article/details/10202699)
#### 有名管道注意事项
- 有名管道必须同时有至少一个读、一个写进程
- 一个进程试图向一个没有读进程的管道写内容则报SIGPIPE错误
- 有名管道分阻塞和非阻塞两种
- 每次发送内容有大小限制、512字节

### 文件和记录锁定
**文件锁定**: 锁定整个文件
**记录锁定**:锁定文件的某一特定部分，从文件的某一相对位置开始的一段连续的字节流，不同于其他以强制性记录结构阻止文件的操作系统，因此记录锁也成为范围锁，对文件的某个范围的锁定。
**咨询式锁定**: 当某一进程对某一文件进行了咨询式锁定后，其他访问该文件的进程被操作系统告知共享文件已经上了锁，但不阻止他们继续对锁定文件的操作。 也就是说，只要有文件的读写权，咨询锁便可以忽略。
**强制锁定**: 某一共享文件被强制锁定后，操作系统之后对每个读写文件的请求检查，只要有锁就拒绝。
#### System V咨询锁
```c
#include <unistd.h> 
int lockf(int fd, int function, long size);
```
- size指明了从文件当前位置开始的一段连续区域长度，size为0时，锁定记录将由当前位置一直到文件尾部。
- function 可以是
    - F_ULOCK为一个锁定区域解锁
    - F_LOCK锁定一个区域
    - F_TLOCK测试并锁定一个区域
    - F_TEST 测试一个区域是否上锁
如果使用F_LOCK，且对应文件已经上锁，那么将会阻塞至文件解锁为止。采用F_TLOCK如果文件已经上锁那么将直接返回不阻塞。
#### BSD咨询锁
```c
#include <sys/file.h> 
int flock(int fd, int operation);
```
operation可以设定为以下几个值
- LOCK_SH 共享锁
- LOCK_EX 互斥锁
- LOCK_UN 解锁
- LOCK_NB 当文件已被锁定时不阻塞
- LOCK_SH  | LOCK_NB 非阻塞性共享锁
- LOCK_EX  | LOCK_NB 非阻塞性互斥锁
BSD有互斥锁和共享锁两种，多个进程可以用于属于同一文件的共享锁，但是某个文件不能存在多个互斥锁或者互斥锁与共享锁并存的状况。
对于已锁定文件，如果采用非阻塞锁定方式，则立即返回调用失败
#### linux其他锁技术
创建辅助文件以表示进程对共享文件的锁操作是Linux其他上锁技术的基本点。如果辅助文件存在，资源便被其他进程锁定了，否则，进程就可以创建辅助文件以对资源上锁。
基于文件一般就要两步，先判断锁是否存在(文件)，不存在再加锁(创建文件)
```c
if((fd=open(file,0))<0){
    fd=creat(file,0644); 
}
```
但是所有的两步操作都应该注意 如果中间进程被CPU打断了如何处理的问题，会不会出错
##### link方法
```c
#include <unistd.h> 
int link(char* existingpath, char* newpath); 
unlink(tempfile) 
```
采用link方式 ，如果文件的新链接名字已经存在，系统调用link便会失败，解锁时unlike这个临时文件即可
##### create方法
如果文件已经存在且不允许写，那么对该文件调用create便会错误返回。因此创建一个禁止写操作的临时文件。
该方法有个缺陷就是root权限对所有都有读写权限
##### open方法
如果open打开文件时增加选项O_CREATE，又设置O_EXCL，如果文件存在，那么调用就失败了，表示加锁了，否则就创建文件返回成功。
##### 总结
- 花费时间更长，因为要多个系统调用
- 需要多用刀一个辅助文件
- 系统崩溃后辅助文件依旧存在
- 多进程不知道如何通知了，不能轮询尝试加锁，应该释放锁时自动通知正在等待的进程
- 不释放锁便结束进程
#### System V IPC
在UNIX system V中引入了几种新的进程通讯方式，即消息队列(Message Queues), 信号量(semaphores) 和共享内存(shared memory)，统称为 System V IPC。
- IPC具体实例在内核中以对象的形式出现，称为IPC对象
- 每个IPC对象在内核中有唯一的标识符，在每一类中标志符唯一，不同类别可以出现相同标识符
- 标识符内核使用，程序中通过关键词key访问IPC对象，关键词也必须唯一
- 关键字的生成很重要，可以采用
```c
// 保证pathname 和 proj名称相同，两个进程就会获得同样的key
key_t ftok(char *pathname, char proj)
```
##### ipcs命令
ipcs 命令显示系统内核的IPC对象状况
- -q 只显示消息队列
- -m 只显示共享内存
- -s 只显示信号量
##### IPC对象结构
```c
struct ipc_perm{
    key_t key;  // 关键字
    ushort uid;  // uid
    ushort gid;   // gid
    ushort cuid;   // 创建者uid
    ushort cgid;   // 创建者gid
    ushort mode;   // 权限
    ushort seq;    // 频率信息忽略
}
```
##### 消息队列
###### 消息的结构msgbuf
```c
struct msgbuf{
    long mtype;   // message类型，区分消息类型
    char mtext[1];  // 消息数据的内容
}
```
mtext虽然在定义中只是一个字符数组，事实上，它的对应部分可以是任意的数据类型，甚至是多个数据类型的集合。 但是加mtype所占的4个字节，每个msgbuf最多占用4056个字节。 例如
```c
struct msgbuf{
    long mtype;   // message类型，区分消息类型
    long request_id;  // 消息数据的内容 部分
    struct client info; // 消息数据的内容 部分
}
```
###### 消息链表节点msg
消息队列在内核中以消息链表的形式出现，链表中每个节点结构就是msg
```c
struct msg{
    struct msg *msg_next; //下一个节点的指针
    long msg_type;  // 和msgbuf中mtype 一样
    char *msg_spot; // 消息内容即msgbuf中mtext在内存中位置
    time_t msg_stime;
    short msg_ts; // 消息长度
}
```
###### 消息链表
内核中的每个消息队列对象都对应一个msgqid_ds结构的数据来保存。
```c
struct msgqid_ds{
    struct ipc_perm msg_perm;
    struct msg *msg_firts;
    struct msg *msg_last;
    __kernel_time_t msg_stime; // 最近一次接受消息的时间
    __kernel_time_t msg_rtime; // 最近一次取出消息时间
    __kernel_time_t msg_ctime; // 最近一次队列发生改动的时间
    struct wait_queue *wwait; // 等待队列的指针
    struct wait_queue *rwait; // 等待队列的指针
    unsigned short msg_cbytes; // number of bytes 
    unsigned short msg_qnum; // number of msg
    unsigned short msg_qybytes; // 占用内存的最大字节数
    __kernel_ipc_pid_t msg_lspid; // 最后一次发送消息的进程pid
    __kernel_ipc_pid_t msg_lrpid; // 最后一次取消息的进程pid
}
```
###### msgget
`int msgget(key_t key, int msgflag)`
- msgget用来创建新的消息队列或获取已有的消息队列。
- msgflag控制操作
    - IPC_CREAT: 消息队列不存在，则创建，否则进行打开操作
    - IPC_EXCL: 消息队列不存在则创建，存在则产生错误。必须和IPC_CREAT一起使用
    - 0600: 除了前面的符号，还可以最后加权限以 | 分隔
###### msgsnd
`int msgsnd(int msqid, struct msgbuf *msgp, int msgsz, int msgflag)`
- msqid:消息队列id
- msgbuf: 发送的内容体
- msgsz: 消息长度， sizeof(struct buf) - sizeof(long)
- msgflag: 控制位 
    - 0: 忽略控制位
    - IPC_NOWAIT: 队列满则不写入，不阻塞 立即返回
###### msgrcv 
`int msgrcv(int msqid, struct msgbuf *msgp, int msgsz, long mtype, int msgflag)`
- mtype: 表示消息类型，如果值为0则返回最旧的消息，否则寻找队列中与之匹配的消息并返回。
- msgflag与msgsnd的msgflag一致
###### msgctl 直接控制消息队列的行为
`int msgctl(int msgqid, int cmd, struct msqid_ds *buf)`
- cmd是要对消息队列的操作
    - IPC_STAT: 去处保存的消息队列数据，存入第三个参数
    - IPC_SET:设定消息队列的msg_perm成员，值由buf给出
    - IPC_EMID:内核中删除消息队列

