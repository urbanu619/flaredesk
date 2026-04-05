# 推广文章草稿

## 标题（二选一）
- 《我用 Go + Vue 做了个 Cloudflare DNS 批量管理面板，CF 控制台做不到的事它都能做》
- 《管理几十个域名太痛了，我开源了一个 Cloudflare DNS 批量操作工具》

---

## 正文

最近在管理大量域名，每次要批量改 DNS 记录都要在 Cloudflare 控制台一个一个点，或者自己写脚本调 API。

忍无可忍，花了几天写了个面板：**Flaredesk** — 专门解决 CF 控制台做不到的批量操作。

GitHub：https://github.com/urbanu619/flaredesk

### 解决了什么问题

**CF 控制台的痛点：**
- 多个账号要来回切换
- 批量添加 DNS 记录只能一条一条
- 想把 50 个域名的 A 记录同时改掉？不可能
- 橙云代理开关要一个个点

**Flaredesk 能做：**
- 多账号统一管理，一个面板全搞定
- 跨域名批量新增：同样的记录，一次性写入 N 个域名
- 跨域名批量删除：按类型/名称，批量清理所有域名的记录
- 一键橙云切换：多个域名同时开启/关闭 CF 代理
- DNS 模板：把常用的记录组合保存下来，新域名直接套用

### 技术栈

- 后端：Go + Gin + GORM
- 前端：Vue3 + Element Plus
- 数据库：MySQL + Redis

### 安全性

所有数据在你自己的服务器上，API Token 不经过任何第三方。开源代码，自己审查。

### MCP Server（给 AI 用户）

同时发布了 MCP Server，可以在 Claude / Cursor 里直接用自然语言管理 DNS：

```bash
# 在 Claude Desktop 配置文件中添加
{
  "mcpServers": {
    "flaredesk": {
      "command": "npx",
      "args": ["flaredesk-mcp"],
      "env": { "CF_API_TOKEN": "你的Token" }
    }
  }
}
```

然后直接跟 Claude 说："帮我把所有域名的 www A 记录指向 1.2.3.4"，它会调用工具批量完成。

npm：https://www.npmjs.com/package/flaredesk-mcp

### 部署

自部署，5 分钟跑起来。需要 Go 1.21+、MySQL、Redis。

详细步骤见 README：https://github.com/urbanu619/flaredesk

---

欢迎 Star，有问题开 Issue。做站群、搞出海、管大量域名的同学应该用得上。

---

## 发布渠道

### V2EX（推荐节点：程序员 / 分享创造）
直接贴正文，标题简短有力，结尾不要太推销。

### 掘金
加上代码截图效果更好，可以补充技术实现细节。

### Twitter/X
缩短版：
> 开源了一个 Cloudflare DNS 批量管理面板 Flaredesk
>
> 解决 CF 控制台做不到的：跨域名批量操作、橙云批量切换、DNS 模板
>
> 还有 MCP Server，在 Claude 里直接用自然语言管 DNS
>
> GitHub: https://github.com/urbanu619/flaredesk
> npm: https://www.npmjs.com/package/flaredesk-mcp

### 目标行业圈（微信群/Telegram）
针对博彩/游戏运维群，改成：
> 管几百个域名的工具，Cloudflare DNS 批量操作面板，开源自部署，数据在自己服务器上。
> 跨域名批量改记录、橙云批量切换、DNS 模板一键套用。
> [GitHub 链接]
