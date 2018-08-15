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
信号量用来控制多个进程对共享资源使用的计数器，常被用于锁定保护机制。信号量对象实际是多个信号量的集合。内核中以数组形式实现。
##### 信号量数据结构
**sem**:
```c
struct sem{
    short sempid;   // 最近一次操作信号量的进程pid
    ushort semval;   // 信号量的计数值
    ushort  semncnt;  // 等待使用资源的进程数目
    ushort semzcnt; // 等待资源完全空闲的进程数目
}
```
**semid_qs**: semid_qs结构被用来存储每个信号量对象的有关信息。
```c
struct semid_ds{
    struct ipc_perm sem_perm; // 与ipc_perm类似
    __kernel_time_t sem_otime;  // 最后一次semop()操作的时间
    __kernel_time_t sem_ctime;   // 最后一次改动发生的时间
    struct sem *sem_base; // 信号量数组的起始地址
    struct sem_queue *sem_pending;  // 还没进行的操作
    struct sem_queue **sem_pending_lat; // 最后一个还没进行的操作
    struct sem_undo *undo; // undo请求的数目
    unsigned short sem_nsems; // 信号量数组的成员数目
}
```
**sembuf**: 定义信号量对象的基本操作
```c
struct sembuf{
    unsigned short sem_num;  // 接受操作的信号量在信号量数组中的序号
    short sem_op    // 定义了信号量可进行的操作(正、负、0)
    short sem_flg;  // 控制操作行为的标志
}
```
sem_op为负数，表示从指定信号量中减去相应的值，对应着获取信号量所监控的资源的操作。 没有指定sem_flg=IPC_NOWAIT标志，那么如果现有信号量数值小于sem_op绝对值(现有资源少于要获取的资源)，调用semop()函数的进程就会被阻塞直到信号量数值大于sem_op绝对值(有足够资源被释放)
sem_op正数，指定的信号量中加上相应的值，对应着释放信号量所监控的资源操作
sem_op为0，那么调用semop()函数的进程将会被阻塞直到对应信号量变为0。操作的实质就是等待信号量所监控的资源被全部使用，动态监控资源的使用并调整资源的分配，避免不必要的等待
##### 相关函数
**semget**: 建立新的信号量对象或者获取已有的对象的标识符
```
int semget (key_t key, int nsems, int semflg); 
// key和semflg与msgget中的参数一致
// nsems 指定了新生成的信号量对象中的信号量的数目，也就是数组长度
```
**semop**: 改变信号量对象中的各个信号量的状态
```c
* int semop ( int semid, struct sembuf *sops, unsigned nsops); 
// semid表示要操作的信号量对象
// sops 就是操作的内容
// nsops保存sops数组的长度，也即semop函数将要进行的操作个数
```
举个例子，如果有一个信号量监控一台最多处理10份作业的打印机，想打印机提交任务的时候。 需要先设置一个任务，在执行后再释放资源
```c
struct sembuf sem_get = { 0, -1, IPC_NOWAIT }; 
if((semop(sid, &sem_get, 1) == -1) {
    perror("semop"); 
}
// do.....
struct sembuf sem_release= { 0, 1, IPC_NOWAIT }; 
semop(sid,&sem_release,1); 
```
**semctl**:直接对信号量对象进行操作，但由于信号量对象是一个数组，需要制定操作的具体index，同时cmd比msgctl更多。 大同小异
```c
int semctl ( int semid, int semnum, int cmd, union semun arg ); 
```
#### 共享内存
被多个进程共享的内存。最快的同新方法，直接将信息映射到内存，省去其他IPC方法的中间步骤
##### 数据结构
**shmid_ds**: 
```c
struct shmid_ds { 
    struct ipc_perm  shm_perm;  // 权限信息
    int              shm_segsz;  // 共享内存的大小 字节为单位
    __kernel_time_t  shm_atime;  // 最近一次进程连接共享内存的时间
    __kernel_time_t  shm_dtime;  // 最近一次进程断开共享内存的时间
    __kernel_time_t  shm_ctime;  // 最近一次shmid_ds结构更改的时间
    __kernel_ipc_pid_t  shm_cpid; // 创建内存的pid
    __kernel_ipc_pid_t  shm_lpid;  // 最近一次连接内存的进程pid
    unsigned short      shm_nattch; // 与共享内存连接的进程数目
    unsigned short    shm_unused;
    void               *shm_unused2;
    void               *shm_unused3;
```
##### 函数操作
**sys_shmget**: 创建和获取已有的共享内存
```c
int shmget ( key_t key, int size, int shmflg); 
```
**shmat**: 当一个进程使用shmget函数获得共享内存标识符之后，可以使用shmat将共享内存映射到晋城自己的内存空间内。之后就像普通内存一样读写了
```c
int shmat ( int shmid, char *shmaddr, int shmflg); 
// shmaddr表示共享内存映射的地址，但要预先分配内存，不方便，因此一般都是置零，系统会自动分配一块未使用内存。
// 如果shmflg指定为SHM_RDONLY则是只读模式
```
**shmctl**:  直接对共享内存操作
**shmdt**: 当不需要共享内存时应该断开连接
```c
int shmdt ( char *shmaddr ); 
```
#### 共享内存与信号量结合使用
共享内存效率最高，但是要保证一致性原则，使用信号量来实现锁机制。 

## 通信协议
通讯协议用于协调不同网络设备之间的信息交互，建立设备之间互相识别的好、有价值的信息机制，如XNS、TCP/IP等
跳过一堆过时的专有网络通信协议
### TCP/IP
TCP/IP实际上是一个一起工作的通信家族，为忘记数据通信提供通路。分为三部分
- Internet协议(IP)
- 传输控制协议TCP和用户数据报文协议UDP
- 处于TCP和UDP之上的一组协议专门开发的应用程序
**网络层协议**:  包括Internet协议IP、网际控制报文协议ICMP和地址识别协议ARP
    - Internet协议(IP): 被设计成互联分组交换通信网，形成一个国际通信环境。负责在源主机和目的地主机之间传输来自较高层软件的成为数据报文的数据块，提供非连接型传递服务。
    - 网际控制报文协议(ICMP): 不是IP层部分，但直接和IP层一起工作，报告网络上的某些出错情况
    - 地址识别协议(ARP): 不是网络层部分，处于IP和数据链路层之间，在32位IP地址和48位局域网地址之间执行翻译的协议。
**传输层协议**: 
    - 传输控制协议TCP: 基于IP提供非连接型传递服务，TCP为应用程序提供可靠地面向连接的传输层服务。 为用户进程之间的对话负责，确保两个以上进程之间的可靠通信。
        1. 监听输入对话建立请求
        2. 请求另一网络站点对话
        3. 可靠的发送和接收数据
        4. 适度的关闭对话
    - 用户数据报文协议UDP: 不可靠的非连接性传输层协议，允许两者之间传送数据而不必建立对话，不适用端对端差错校验。
**应用程序**: 包含Telnet、FTP、TFTP、简单的文件传送协议SMTP和域名服务DNS
#### Internet协议IP
提供子网的互联，形成较大的网际，使不同的子网之间能传输数据。 网际由许多自治系统组成，每个系统是一个中央管理的网络或一系列子网。
IP规定包如何从一个子网路由选择到另一个子网。自治系统中每个节点具有唯一的IP地址。IP使用本身的帧头和检查来确保数据报文的正确传送。路由选择表列出了子网上各种不同节点之间的通路和通路开销，如果个别节点之间有较多的通路，则可选择最经济的一条。如果数据包较大，目的地不能接受，则将他分成较小的段。
##### 数据传送
主要目的是为数据输入/输出网络提供基本算法，为高层协议提供无连接的传送服务。意味着IP将数据传递给接收站点以前不在传输站点和接收站点之间建立对话。只是封装和传递数据，不想发送者或接受者报告包的状态，不处理故障。由高层TCP负责执行消除差错。
IP协议将正确格式化的数据包传递给目的地站点，不期待状态回应。由于是无连接的协议，因此可能接受、发从错误序列的数据。
IP接受并格式化数据，以便传输到数据链路层。此外IP还检索来自数据链路的数据，并将它发送给高层。
IP协议不关心包内的数据类型，只知道把IP帧头的控制协议加到高层协议(TCPUDP)所接受的数据上，并试图把他传递给网络上某个节点
##### 分段包
包太大会造成不能一次传输，因此需要切分包再组装
 - 分段控制、16位: 识别、标志和分段偏移
 - 生存时间TTL、8位:防止出现错误的包再网际路由器之间不断循环
 - 协议字段: 指定高级协议
 - 校验和: 16位循环冗余检验，确保帧头的完整性
 - IP选项字段
 - IP源和目的地址字段
##### IP寻址
#### 网际控制报文协议ICMP
IP需要它帮助传输差错和控制报文。
- 一种报文是回应请求，用于测试目的地是否可达，此外回应请求报文还跟踪响应时间，以便确定路线的平均延时。
- 源断开报文，拥塞控制方法，IP网关接受了较多的包，超出控制能力，则靠ICMP进行摆脱
- 路由选择变化请求，例如重定向数据请求
- 用于计时打印请求和确认。用于估算网络上的平均往复时延。
#### 传输控制协议TCP
TCP具有严格的内装差错检验算法来确保数据的完整性。TCP是传输层协议，目的是允许数据同网络上的另外站点进行可靠的交换。提供端口编号的译码，易识别主机的应用程序，从而完成可靠地数据传输。
TCP使用顺序编号和确认信息同网络上另外的站点交谈。接收站点可能接受两个同样的包，使用带有确认信息的顺序编号，这种处理称全双工。连接的每一端都必须考虑到另一端的需要而维持自己的顺序编号。
TCP是面向自己的顺序协议，包内的每一个字节被分配一个顺序编号，并分配给每包一个顺序编号。
#### 用户数据报文协议UDP
无连接的、不可靠的传输服务。接收数据时不想发送方提供确认信息，不提供输入包的顺序，出现丢包或重复包也不会像发送方发出差错报文，和IP协议很像。 UDP主要是分配和管理端口编号，以正确无误的识别运行在网络站点上的个别应用程序。


## berkeley套接字
### Socket
每一个Socket都用一个半相关描述
{ 协议， 本地地址， 本地端口 }
一个完整的Socket则用一个相关描述
{ 协议， 本地地址， 本地端口， 远程地址， 远程端口}
每一个Socket有一个本地的唯一Socket号，由操作系统分配。
Socket是面向客户-服务器模型设计的，针对客户和服务器程序提供不同的Socket系统调用。 可以随机申请一个Socket号，服务器拥有全局公认的Socket，任何客户都可以向他发出连接请求和信息请求。

### Socket三种类型
#### 流式套接字 SOCK_STREAM
提供可靠地、面向连接的通讯流，如果通过流式套接字发送了顺序数据：1、2,那么到达远端的数据顺序也是1、2。
流式套接字使用了TCP协议，TCP确保数据传输是正确的，并且是顺序的。
Telnet协议和浏览器都是基于流式连接。
 
图4

#### 数据报套接字 SOCK_DGRAM
定义了一种无连接的服务，数据通过相互独立的报文进行传输，无序的，不保证可靠、无差错。 原始套接字允许对底层协议如IP、ICMP直接访问。
    - 发送一个数据包，他可能不到达
    - 可能不同顺序到达
    - 到达了可能内容存在错误
数据报套接字使用IP和UDP。 UDP不像流式套接字那样维护一个打开的连接，只是把数据打包，贴上IP，然后发出去，整个过程不建立连接。
tftp、bootp基于UDP
事实上每一种UDP程序都有自己的数据确认协议，例如TFTP程序会对每一个请求返回一个ACK包，如果客户端5s拿不到ACK返回则会重发一次。

图5


#### 原始套接字
主要应用于一些协议的开发，对底层的操作，功能强大，但基本用不到

### TCP/IP
传输控制协议/网络协议，作为软件的网络组成部件而设计。
#### 控制数据的协议
TCP以连接为基础，两台计算机必须先建立一个连接，才能真正传输数据。
UDP无连接服务，数据直接发送而不必建立网络连接。和TCP比占用带宽少。
#### 数据路由协议
IP 处理实际上传输数据
ICMP处理IP状态信息
RIP决定信息传输的最佳路由路线协议中的一个
OSPF决定路由的协议
#### 用户服务
FTP 从一台计算机上传输一个或多个文件到另一台计算机
TELNET 允许一个远程登陆到一台机器
#### 其他协议
RPC 远程的应用程序通过简单的有效地手段练习本地的应用程序
SMTP 专门为电子邮件在多机器中传输的协议，平时发邮件的SMTP服务器提供的必然服务
SNMP 超级用户准备的服务，超级用户可以通过它来进行简单的网络管理

### 套接字地址
它是通过标准的UNIX文件描述符和其他的程序通讯的一个方法
#### 数据如何传输
数据切分成包，包的数据头被第一层(HTTP协议)协议加上第一层协议数据，然后整个包再被下层协议包装一次(UDP)，然后再被下层协议包装(IP协议)，最后被底层硬件层包装一层信息(Ethernet信息头)

### 系统调用
int socket (int domain, int type, int protocol)  创建socket描述符
close(sockfd)
int shutdown(int sockfd, int how)
#### 有连接的
```
int bind(int sockfd,  struct sockaddr *my_addr, int addrlen);  将一个socket描述符和系统一个端口绑定
int listen(int sockfd, int backlog)  创建等待连接请求， backlog为等待最大队列长度
int accept(int sockfd, void *addr, int *addrlen) 接收连接请求，之后会产生一个新的fd，原fd不变，新的fd可以进行send()和recv()操作
int send(int sockfd, const void *msg, int len, int flags) 与远程连接建立的socket 发送信息到远程机器
int recv(int sockfd, void *buf, int len, unsigned int flags) 读取远程连接发来的数据，指定缓冲区长度
int connect (int sockfd, struct sockaddr *serv_addr, int addrlen) 连接到远程服务器
intgetpeername(int sockfd, struct sockaddr *addr, int *addrlen)获取远程连接的是谁
```
#### 无连接的
`int sendto(int sockfd, const void *msg, int len, unsigned int flags, const struct sockaddr *to, int tolen) `与send类似，但是无提前建立好的连接符，但是知道对方的ip和端口
int  recvfrom() 与send类似
### IO模式
对于网络请求，数据从网络层拷贝到内核缓存区，再从内核缓冲区拷贝到程序的数据区

图11

#### 阻塞IO

图6


#### 非阻塞IO
在请求IO操作时，立即返回一个错误，然后程序之后主动循环去请求IO是否完成(polling操作)，但这极大浪费CPU

图7

#### IO多路复用
调用`select`或者`poll`函数，在调用他们的时候阻塞，在得到某个可以操作的文件描述符之后再调用`recvfrom`，优势在于可以一次监听多个文件描述符，只要其中有一个是可操作状态之后，`select`就会返回

图8

#### 信号驱动IO
在调用IO操作时，立即返回。 等待IO操作准备就绪后，发送一个信号(SIGIO)给程序，程序在收到信号后再进行IO处理操作。

图9

#### 异步IO
程序告诉内核要进行什么IO操作，然后立即返回，内核在完成全部操作后通知程序，完成了IO操作
与信号驱动区别
    - 信号驱动在文件描述符可以被操作时通知程序
    - 异步IO在内核完成所有IO操作后通知程序

#### fcntl()
将一个套接字设置为非阻塞模式，之后需要不断轮序，如果无可读取数据则会立刻返回-1
```
sockfd = socket(AF_INET, SOCK_STREAM, 0);
fcntl(sockfd, F_SETFL, O_NONBLOCK)
```












