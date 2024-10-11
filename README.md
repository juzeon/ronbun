# ronbun

计算机会议索引&论文助手，会议分类、论文爬虫、本地论文搜索、论文全文翻译工具箱。

## 功能

- 按[ccf-deadlines](https://github.com/ccfddl/ccf-deadlines)论文分类：会议分类、CCF Rank
- dblp 论文爬虫、论文摘要爬虫
- 论文关键词搜索、基于向量的语义搜索
- 论文全文翻译，双语对照，输出阅读器视图，支持 LaTeX 公式渲染

## 编译

需要 Go 1.23+

```bash
git clone https://github.com/juzeon/ronbun.git
cd ronbun/
go build
```

## 使用

### 依赖

本项目依赖[ccf-deadlines](https://github.com/ccfddl/ccf-deadlines)，需要先 clone：

```bash
git clone https://github.com/ccfddl/ccf-deadlines.git
```

记录下路径，假设这里是`D:\dev\ccf-deadlines`。

在摘要 Text Embedding 和向量搜索时需要用到 jina 的 Embedding API，可以在这里申请：<https://jina.ai/embeddings/>

![](https://public.ptree.top/ShareX/2024/10/10/1728558627/fstutMhIND.png)

清除浏览数据和 Cookie 刷新后这个 API Key 会自动生成新的，可以无限获取，每个都有免费 1M Token 的额度。建议先申请 10 个。

### 本体

第一次使用需要运行一次初始化：

```bash
./ronbun
# 然后打开User Home目录下的.ronbun/config.yml文件，填入配置文件
```

`config.yml`：

```yaml
ccf_path: "D:\\dev\\ccf-deadlines" # 改成 ccf-deadlines 的路径
jina_tokens: # 填入你的 jina token
  - jina_b480ecf26bb11111111111111111RqVNBhptko1BbEUbuJImbl
  - jina_726ceca975d322222222222222ZwvNZYc5cKL9luIIpHKPRgr
  - jina_1e1ff02f590f333333333333333387pl9aD0aN5cOeqyx
concurrency: 20 # 多线程爬虫的并发度（线程数），比如 20 个线程
search_limit: 20 # 向量搜索时显示的论文数量限制，比如显示前 20 个最相关的
grobid_urls: # Grobid 服务的 URL，用于全文翻译时识别 PDF；以下网址来自 https://github.com/binary-husky/gpt_academic，也可以自己搭建
  - https://qingxu98-grobid.hf.space
  - https://qingxu98-grobid2.hf.space
  - https://qingxu98-grobid3.hf.space
  - https://qingxu98-grobid4.hf.space
  - https://qingxu98-grobid5.hf.space
  - https://qingxu98-grobid6.hf.space
  - https://qingxu98-grobid7.hf.space
  - https://qingxu98-grobid8.hf.space
openai: # OpenAI API 的信息，用于全文翻译
  endpoint: https://api.openai.com
  model: gpt-4o-mini
  key: sk-splvM85y111111111111111111aB1DdAf
```

然后再运行工具查看帮助。

数据库文件为`.ronbun/ronbun.db`，可以使用 Navicat 打开。

## 子命令

### update-list - 从 dblp 爬论文标题和链接

运行之后会提示选择会议分类、CCF 等级，等等，然后可以开始多线程爬取。数据存到数据库中。

dblp 只有会议、年份、论文标题、链接等等信息，所以等下我们还要根据 doi 链接跳转过去爬取摘要（Abstract）。

### update-abstract - 从 doi 源站爬论文摘要

目标论文是数据库里所有没有摘要的论文。

目前实现的 doi 源站：

- `dl.acm.org`：有请求频率限制
- `ieeexplore.ieee.org`：基本没有请求频率限制
- `link.springer.com`：无请求频率限制
- `www.usenix.org`：无请求频率限制

有请求频率限制的可以自行配置代理 IP 池。推荐项目：

<https://github.com/honmashironeko/ProxyCat>

<https://github.com/jhao104/proxy_pool>

公共代理可能不太稳定，要寻求稳定的话，建议自己使用 Tor 和 V2Ray/Xray 搭建代理 IP 池，教程请看：<https://blog.skyju.cc/post/v2ray-tor-proxy-ip-pool/>

dl.acm.org 频率限制比较严格，用 Tor 代理 IP 池可以秒杀。但 ieeexplore.ieee.org 屏蔽了 Tor 的 IP，好在没什么频率限制，直连就行。

路由规则的详细配置请参考 Xray 的文档：<https://xtls.github.io/>

### update-embedding - 生成论文摘要 Text Embedding

这一步会调用 jina 的 API 生成论文摘要的 Text Embedding 向量，存到数据库里。目标论文是数据库里所有有摘要但没 Embedding 向量的论文。

### search - 标题关键词搜索

根据标题搜索论文，可以选择过滤会议分类、CCF 等级等。

[结果示例](https://public.ptree.top/ShareX/manual/Search%20for%20serverless.html)

### search-vec - 语义搜索

文档搜文档，基于数据库里的论文摘要。将输入文档的向量和数据库里的向量匹配，取前 20 个最近似的论文。

[结果示例](https://public.ptree.top/ShareX/manual/Search%20by%20document%202024-10-10%2019%2129%2158.html)

### translate - 全文翻译

先调用 Grobid 从 PDF 中提取文本，再调用 OpenAI API 大模型全文翻译。输出双语对照的结果，特别按阅读器视图进行优化，并支持 LaTeX 公式渲染

[结果示例](https://public.ptree.top/ShareX/manual/Translation%20for%20fast24-li.pdf.html)