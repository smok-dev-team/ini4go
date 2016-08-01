## Config
ini 配置文件读取、写入组件。

* Section - 支持分组;
* 多文件 - 可一次读取多个文件;
* 变量 - 支持变量替换;
* List - 支持读取重复的 key, 其值为一个 list;
* 注释 - 读取、写入注释;
* 默认值 - 读取值的时候, 可以设定默认值。

##### 读取文件

```
var r = New(false)
r.LoadFiles("./test.conf")
fmt.Println(r.GetValue("s1", "sk1"))
```

##### 写入文件

```
var r = New(false)
r.SetValue("s1", "p1", "v1")
r.MustSection("s1").MustOption("p2").SetValue("v2")
r.MustSection("s2").MustOption("p2").SetValue("v2")
fmt.Println(r.WriteToFile("./output.conf"))
```