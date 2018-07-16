# Linux网络编程
## 网络7层
- 应用层
- 表示层
- 会话层
- 传输层: TCP、UDP、NVP 等协议
- 网络层: IP、ICMP、ARP、反向ARP 等
- 数据链路层: 以太网、PDN等
- 物理层: 电话线、电缆等
## 进程控制
### 进程概念
- fork: 通过复制调用进程来创建新的进程
- exec: 通过一个新的程序覆盖原内存空间，实现进程的转变。各种exec系统调用之间的区别在于他们的参数不同
- wait:初级的进程同步措施，使一个进程等待知道另一个进程结束为止
- exit: 用来终止一个进程
### fork
系统调用fork是建立进程的基本操作。 新进程是调用fork进程的副本，新进程与原进程有一样的变量、运行一样的代码，并且获得同样的fd等。 
- 新进程是子进程，原进程是父进程，他们两个开始并发执行，都从fork之后的那句代码开始执行，因为同样的内存空间，PC(程序计数器)都指向同一行将要执行的代码
- fork返回值是-1(fork失败)， 0(表示当前是子进程)、pid(表示当前是父进程)，那返回值来判断是子进程还是父进程
``` c
#include <stdio.h> 
#include <unistd.h> 
main() { 
    pid_t pid;
    printf("Now only one process\n"); 
    printf("Calling fork...\n"); 
    pid=fork();
    if (!pid) {
        printf("I’m the child\n"); 
    }else if (pid>0) {
        printf("I’m the parent, child has pid %d\n",pid); 
    }else{ 
        printf("Fork fail!\n"); 
    }
}
```
 ![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/linux_socket_1.jpeg)
### exec
Exec不会产生一个新的进程，而是覆盖原进程，从新程序的入口开始执行，这样就形成了新的进程。 exec调用成功后不会有人和返回，后续代码也不会执行，后续代码执行表示exec执行失败。
- exec对应的系统调用有很多，他们大都是参数形式不同的区别
- execl系统调用第一个参数为可执行文件的路径，第二个为文件名称，后面可以像执行shell一样接任意多个参数，但最后需要一个NULL指针表示结束
``` c
// 把上面代码编译后命名为aa.out
#include <stdio.h> 
#include <unistd.h> 
main() { 
    pid_t pid;
    printf("Now only one process\n"); 
    printf("Calling fork...\n"); 
    pid=fork();
    if (!pid) {
        printf("I’m the child\n"); 
        execl("./aa.out", "aa", NULL);
        printf("exec fail\n"); 
    }else if (pid>0) {
        printf("I’m the parent, child has pid %d\n",pid); 
    }else{ 
        printf("Fork fail!\n"); 
    }
}
```
 ![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/linux_socket_2.jpeg)
### fork 与exec连用
就是上面那个例子，利用fork后返回的pid判断是父进程还是子进程，子进程的话调用exec执行另一个程序即可
下面两段代码, 父进程fork之后，子进程立即调用exec.out程序，然后父进程wait(NULL)等待子进程结束。
``` c
#include <stdio.h> 
#include <unistd.h> 
main() { 
    pid_t pid;
    printf("Now only one process\n"); 
    printf("Calling fork...\n"); 
    pid=fork();
    if (!pid) {
        printf("I’m the child\n"); 
        execl("./exec.out", "exec", NULL);
        printf("exec fail\n"); 
    }else if (pid>0) {
        wait(NULL);
        printf("I’m the parent not in exec, child has pid %d\n",pid); 
    }else{ 
        printf("Fork fail!\n"); 
    }
}
```
``` c
#include <stdio.h> 
#include <unistd.h> 
main() { 
    printf("I am exec program\n"); 
}
```
 ![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/linux_socket_3.jpeg)
### shell执行原理
- shell接到前台命令时如ls就调用fork、exec、同时shell作为父进程就wait等待exec调用的ls命令结束
- shell接到后台命令时，则fork、exec之后并行，shell作为父进程不进行wait等待
### wait 进程的同步
wait接受一个指针参数，如果传入NULL则表示忽略，否则指针被赋值为第一个结束的子进程的结束状态，即exit退出时的状态值。
### 特殊情况
- 子进程结束、父进程并不wait: 子进程结束，进入过度状态，不占用内核资源，但是在进程表中有一项记录，当父进程之后进入wait时便可取到这条记录，然后系统会删除这条记录
- 子进程结束前，父进程已经不存在: 子进程交给系统的初始化进程所属(更改父进程)
### 进程属性
- PID:  进程ID不在叙述，getpid()。    0 是进程调度的进程。 1是init进程
- PGRD: 进程组ID。 最初从fork、exec那里继承进程组ID，但也可以调用setpgrp()来更换新组，新的组ID就是当前进程的PID，同时自己变成组首。更换组id有好处是将进程的生命周期拉长，就像nohup启动程序那样
- 环境变量，环境变量是以NULL结尾的字符串的集合。 
``` c
#include <stdio.h> 
#include <unistd.h> 

extern char** environ;
main() { 
    char** env=environ;
    while (*env) { 
        printf("%s\n",*env++);
    }
}
```
- 进程的实际用户和实际用户组: 永远是启动该进程的用户和用户组，不可更改
- 进程的有效用户和有效用户组: 被用来确认某个进程是否有访问某些文件的权限。 root进程可以随意更改，但普通进程只能更改为当前用户和用户组。
- 进程的资源: 可以获取和设定进程对资源的使用状况和使用限制。 getrlimit 和setrlimit。 只有root权限可以提供进程资源的限制，普通权限只能降低
- 进程优先级: nice(-1)。 Root权限可以提高进程优先级，传一个负数， 普通进程只能传入一个正数，降低进程的优先级
### 守护进程
独立于任何终端之外的进程，在后台默默执行，不会被终端的中断所打断。
守护进程的创建方式
    - 通过/etc/rc.d创建，很多系统服务是通过这种方式，会获得超级用户权限
    - 网络服务通过inetd方式创建，监听各种网络请求
    - cron启动的定时任务也是一种守护进程
    - at启动的处理程序
    - nohup启动的进程
如果是shell前台启动守护进程那就需要子进程去更改组id，否则shell进程死掉时，系统会请求同组id的所有进程。
## 进程间通信
### 信号
最古老的方式，系统通通知一个或多个进程异步事件的发生。信号不止能从内核发往进程，也可以进程发往进程。 但信号不能携带信息，只能依靠固有的N中信号来表达各自的含义。
进程在接到某些信号时，就会停止工作，相当于临时加入一个exit命令。发生核心转贮，即把内存映像写入当前目录的core文件中去，包含执行最后状态的全部变量、寄存器和内核中的控制信息等，linux的调试程序gdb知道如何解读这个文件，这样便可以方便调试。 
系统调用abort()就是自己给自己发送一个非正常退出的命令，然后触发核心转贮，成为一种程序调试的辅助手段。
#### 信号的处理
大部分的信号都是终止进程的作用，但有时候我们在处理重要数据的时候可以暂时接管除了SIGKILL 和SIGSTOP以外的信号的处理动作
int signal (int sig, __sighandler_t handler);
第一个参数是信号的类型，第二个是处理函数 SIG_IGN 表示忽略该信号、SIG_DFL 表示恢复该信号的默认处理动作
``` c
// 执行代码发现ctrl+c失效了，因为程序忽略了SIGINT信号
#include <stdlib.h>
#include <stdio.h> 
#include <signal.h>
int main(void) {
    signal(SIGINT,SIG_IGN); 
    printf("xixi\n");
    sleep(10); 
}
```
#### 信号的继承
- 子进程在fork后会继承父进程的信号处理方式，但子进程调用exec后则会恢复信号的默认处理动作
#### 信号系统调用
正常情况下，系统调用时不能被打断的，除了几个例外，如wait、pause、以及慢速设备的读取
#### 连续处理信号
- 连续收到多个相同信号则一个一个串行处理，不会提前中止正在执行的那个信号,并且最多积累一个等待执行，不会积累多个
- 执行信号A时收到信号B则中止信号A的处理，转向处理信号B，处理完信号B再回头处理信号A。
- 同种信号不能积累！！！不能积累。 也就是说最多有两个，一个在执行、一个在等待，不可能出现两个在等待
#### 信号缺点
- 不能保证顺序
- 不能保证不丢失(不积累)
#### 进程间发送信号
通过系统调用kill(int pid, int sig)来完成
- pid为0则发给当前所有同pid组的进程
- pid为-1则由进程号从高到低所有进程(受权限控制)
- pid为具体的单个数值
- 普通权限的进程只能想同一个用户的进程发信号，也就是意味着不能向其他用户的进程发信号，root可以
#### alarm(int time)系统调用
在指定的时间秒后发送一个SIGALARM信号给自己，同样不能积累，第二次调用会覆盖第一次调用，但如果第二次是alarm(0)是取消，但如果是具体数字则会在减去第一次调用的时间消耗。
#### pause()系统调用
让进程等待，直到接收到一个信号为止
#### setjmp() 与longjmp()
int setjmp(jmp_buf env)来保存当前程序的执行位置，通过保存堆栈环境实现
void longjmp(jmp_buf env, int val)   跳转至setjmp位置时，val就是执行setjmp的返回值。虽然没啥用
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
### System V IPC
在UNIX system V中引入了几种新的进程通讯方式，即消息队列(Message Queues), 信号量(semaphores) 和共享内存(shared memory)，统称为 System V IPC。
- IPC具体实例在内核中以对象的形式出现，称为IPC对象
- 每个IPC对象在内核中有唯一的标识符，在每一类中标志符唯一，不同类别可以出现相同标识符
- 标识符内核使用，程序中通过关键词key访问IPC对象，关键词也必须唯一
- 关键字的生成很重要，可以采用
```c
// 保证pathname 和 proj名称相同，两个进程就会获得同样的key
key_t ftok(char *pathname, char proj)
```
#### ipcs命令
ipcs 命令显示系统内核的IPC对象状况
- -q 只显示消息队列
- -m 只显示共享内存
- -s 只显示信号量
#### IPC对象结构
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
#### 消息队列
##### 消息的结构msgbuf
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
##### 消息链表节点msg
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
##### 消息链表
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
##### msgget
`int msgget(key_t key, int msgflag)`
- msgget用来创建新的消息队列或获取已有的消息队列。
- msgflag控制操作
    - IPC_CREAT: 消息队列不存在，则创建，否则进行打开操作
    - IPC_EXCL: 消息队列不存在则创建，存在则产生错误。必须和IPC_CREAT一起使用
    - 0600: 除了前面的符号，还可以最后加权限以 | 分隔
##### msgsnd
`int msgsnd(int msqid, struct msgbuf *msgp, int msgsz, int msgflag)`
- msqid:消息队列id
- msgbuf: 发送的内容体
- msgsz: 消息长度， sizeof(struct buf) - sizeof(long)
- msgflag: 控制位 
    - 0: 忽略控制位
    - IPC_NOWAIT: 队列满则不写入，不阻塞 立即返回
##### msgrcv 
`int msgrcv(int msqid, struct msgbuf *msgp, int msgsz, long mtype, int msgflag)`
- mtype: 表示消息类型，如果值为0则返回最旧的消息，否则寻找队列中与之匹配的消息并返回。
- msgflag与msgsnd的msgflag一致
##### msgctl 直接控制消息队列的行为
`int msgctl(int msgqid, int cmd, struct msqid_ds *buf)`
- cmd是要对消息队列的操作
    - IPC_STAT: 去处保存的消息队列数据，存入第三个参数
    - IPC_SET:设定消息队列的msg_perm成员，值由buf给出
    - IPC_EMID:内核中删除消息队列

#### 信号量
信号量用来控制多个进程对共享资源使用的计数器，常被用于锁定保护机制。

