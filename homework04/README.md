# 博客实验室

在 `homework04` 目录执行 `go run .`，或将 IDE 的工作目录设为该目录后运行。默认页面地址为 <http://localhost:8888/index>，实际端口取自 `config.yaml`。新增页面后需要停止旧进程并重新启动。

页面包含注册、登录、文章管理、评论操作、接口清单和最近一次请求响应。HTML 通过 `go:embed` 嵌入程序，无外部前端依赖；修改 HTML 后也需要重新编译运行。

## 接口清单

| 方法 | 路径 | 请求 / 功能 |
| --- | --- | --- |
| GET | `/index` | 操作页面，无需登录 |
| POST | `/api/v1/users/register` | `username`、`email`、`password` |
| POST | `/api/v1/users/login` | `username`、`password`，返回 `data.authorization` |
| GET | `/api/v1/posts` | 当前登录用户的文章列表，空结果为 `[]` |
| POST | `/api/v1/posts` | `title`、`content`，作者来自登录身份 |
| GET | `/api/v1/posts/:id` | 自己的文章详情 |
| PUT | `/api/v1/posts/:id` | `title`、`content`，完整替换，以路径 ID 为准 |
| DELETE | `/api/v1/posts/:id` | 删除自己的文章 |
| POST | `/api/v1/comments` | `postId`（正整数）、`content`，返回新评论 ID |
| GET | `/api/v1/comments?postId=1` | 指定文章的评论列表，需登录，空结果为 `[]` |
| DELETE | `/api/v1/comments/:id` | 删除自己的评论 |

文章和评论接口均要求 `Authorization: <JWT>`，不要添加 `Bearer` 前缀。页面登录后自动添加此请求头；Token 只保存在页面内存中，刷新页面需要重新登录。“清除本页登录状态”是本地操作，不会撤销已经签发的 JWT。

“我的文章”会直接显示每篇文章对应的评论，文章列表及详情的 `data` 中包含 `comments` 数组。独立评论列表使用文章 ID 查询，显示评论人、内容和毫秒时间；只为自己的评论显示删除按钮。发表评论、删除评论后自动刷新两个列表。登录用户可以查询已有文章的评论，删除仍限评论本人。列表返回独立响应对象，不暴露关联用户的密码等字段。删除前页面会请求确认。

响应使用 `code`、`msg`、`data` / `error`。文章时间格式为 `2006-01-02 15:04:05.000`（服务器本地时间）。常见状态为 400 无效 ID、401 未登录或凭据错误、403 非本人资源、404 资源不存在、409 重复账户、422 参数校验失败。

## 验证

从仓库根目录运行 `go test ./homework04/...`。集成测试使用临时 SQLite 数据库，覆盖全部业务接口及身份冒用、跨用户操作的拒绝，不改动本地 `data/blogs.db`。
