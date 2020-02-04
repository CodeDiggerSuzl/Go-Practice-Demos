# 命令行程序 Simple Media Player (SMP)

## 实现的功能

1. 音乐库功能, 查看, 添加, 删除里面的音乐曲目;
2. 播放音乐;
3. 支持 MP3 和 WAV, 也可以随时支持更多的音乐类型;
4. 退出程序.

## 接收的命令

运行后进入一个循环, 用于监听输入的状态, 接收如下命令:

1. 音乐库管理命令:`lib`, 包括 `list`,`add`,`remove`;
2. 播放管理:`play` 命令,`play` 后带歌曲名参数;
3. 退出程序:`q` 命令

## 模块

### 🎵 音乐库

音乐管理模块, 管理的对象为音乐. 包含如下信息:
- 唯一 id;
- 音乐名字;
- 艺术家名字;
- 音乐位置;
- 音乐类型 (MP3 或者 WAV)

定义音乐结构体
```go
type MusicEntry struct{
    Id string
    Name string
    Artist string
    Source string
    Type string
}
```
使用一个数组切片作为基础的储存结构, 其他操作都只是对这个数组切片的包装.

[manager.go](mlib/manager.go)

编写 `manager` 后要立马进行单元测试.

### ▶️ 音乐播放

音乐播放应该是一个很容易扩展的功能, 不应该在动代码的时候就大动代码.

设计一个简单的播放函数 `func Play(source, mtype string)`, 没直接传入 `MusicEntry` 是因为它包含了很多多余的信息. 本着最小的原则, 设计一个简单的接口.

```go
type Player interface{
    Play(source string)
}
```

### 主程序
和面向对象关系不大, 简写代码.

[main.go](./src/main.go)

### 遗留问题
1. 多任务的加入: 用户界面, 音乐播放,视频播放等功能. 可以使用`goroutine` 来进行完善.
2. 对播放的控制: 暂停,停止设置播放位置等, 使用`channel`和`goroutine`来实现.

## 小结
Go 看似简陋的面向对象的编程特性已经能够满足需求.