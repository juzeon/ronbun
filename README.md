# ronbun

计算机会议索引&论文助手，会议分类、论文爬虫、本地论文搜索工具箱。

## 功能

- 按[ccf-deadlines](https://github.com/ccfddl/ccf-deadlines)论文分类：会议分类、CCF Rank
- dblp论文爬虫、论文摘要爬虫
- 论文关键词搜索、基于向量的语义搜索

## 编译

需要Go 1.23+

```bash
git clone https://github.com/juzeon/ronbun.git
cd ronbun/
go build
```

## 使用

### 依赖

本项目依赖[ccf-deadlines](https://github.com/ccfddl/ccf-deadlines)，需要先clone：

```bash
git clone https://github.com/ccfddl/ccf-deadlines.git
```

记录下路径，假设这里是`D:\dev\ccf-deadlines`。

在摘要Text Embedding和向量搜索时需要用到jina的Embedding API，可以在这里申请：<https://jina.ai/embeddings/>

![](https://public.ptree.top/ShareX/2024/10/10/1728558627/fstutMhIND.png)

清除浏览数据和Cookie刷新后这个API Key会自动生成新的，可以无限获取，每个都有免费1M Token的额度。建议先申请10个。

### 本体

第一次使用需要运行一次初始化：

```bash
./ronbun
# 然后打开User Home目录下的.ronbun/config.yml文件，填入配置文件
```

`config.yml`：

```yaml
ccf_path: "D:\\dev\\ccf-deadlines" # 改成ccf-deadlines的路径
jina_tokens: # 填入你的jina token
  - jina_b480ecf26bb11111111111111111RqVNBhptko1BbEUbuJImbl
  - jina_726ceca975d322222222222222ZwvNZYc5cKL9luIIpHKPRgr
  - jina_1e1ff02f590f333333333333333387pl9aD0aN5cOeqyx
```

然后再运行工具查看帮助。

数据库文件为`.ronbun/ronbun.db`，可以使用Navicat打开。

## 子命令

### update-list - 从dblp爬论文标题和链接

运行之后会提示选择会议分类、CCF等级，等等，然后可以开始多线程爬取。数据存到数据库中。

dblp只有会议、年份、论文标题、链接等等信息，所以等下我们还要根据doi链接跳转过去爬取摘要（Abstract）。

### update-abstract - 从doi源站爬论文摘要

目标论文是数据库里所有没有摘要的论文。

目前实现的doi源站：

- dl.acm.org：有请求频率限制
- ieeexplore.ieee.org：基本没有请求频率限制
- link.springer.com：无请求频率限制
- www.usenix.org：无请求频率限制

有请求频率限制的可以自行配置代理IP池。推荐项目：

<https://github.com/honmashironeko/ProxyCat>

<https://github.com/jhao104/proxy_pool>

公共代理可能不太稳定，要寻求稳定的话，建议自己使用Tor和V2Ray/Xray搭建代理IP池，教程请看：<https://blog.skyju.cc/post/v2ray-tor-proxy-ip-pool/>

dl.acm.org频率限制比较严格，用Tor代理IP池可以秒杀。但ieeexplore.ieee.org屏蔽了Tor的IP，好在没什么频率限制，直连就行。

路由规则的详细配置请参考Xray的文档：<https://xtls.github.io/>

### update-embedding - 生成论文摘要Text Embedding

这一步会调用jina的API生成论文摘要的Text Embedding向量，存到数据库里。目标论文是数据库里所有有摘要但没Embedding向量的论文。

### search - 标题关键词搜索

根据标题搜索论文，可以选择过滤会议分类、CCF等级等。

### search-vec - 语义搜索

文档搜文档，基于数据库里的论文摘要。将输入文档的向量和数据库里的向量匹配，取前20个最近似的论文。