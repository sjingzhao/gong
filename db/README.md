# DB Package

数据库访问层，使用 [sqlc](https://github.com/sqlc-dev/sqlc) 生成类型安全的 Go 代码。

## 目录结构

```
db/
├── schema.sql      # 数据库表结构定义
├── query.sql       # SQL 查询语句
├── sqlc.yaml       # sqlc 配置文件
└── README.md       # 本文件
```

## 快速开始

### 1. 安装 sqlc

```bash
# macOS
brew install sqlc

# 或使用 Go 安装
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### 2. 生成 Go 代码

```bash
cd db
sqlc generate
```

### 3. 使用示例

```go
package main

import (
    "context"
    "database/sql"
    "your-module/db"
)

func main() {
    conn, err := sql.Open("postgres", "your-dsn")
    if err != nil {
        panic(err)
    }
    
    queries := db.New(conn)
    ctx := context.Background()
    
    // 获取用户
    user, err := queries.GetUser(ctx, 1)
    if err != nil {
        panic(err)
    }
    
    // 创建用户
    newUser, err := queries.CreateUser(ctx, db.CreateUserParams{
        ID:    2,
        Name:  "John Doe",
        Email: sql.NullString{String: "john@example.com", Valid: true},
    })
    
    // 列出用户
    users, err := queries.ListUsers(ctx, db.ListUsersParams{
        Limit:  10,
        Offset: 0,
    })
}
```

## 可用的查询

### Users 表

| 方法名 | 类型 | 描述 |
|--------|------|------|
| `GetUser` | `one` | 根据 ID 获取单个用户 |
| `CreateUser` | `one` | 创建新用户 |
| `ListUsers` | `many` | 分页列出所有用户 |
| `UpdateUserEmail` | `exec` | 更新用户邮箱 |

### Authors 表

| 方法名 | 类型 | 描述 |
|--------|------|------|
| `DeleteAuthor` | `exec` | 根据 ID 删除作者 |

## 数据库 Schema

### authors 表

```sql
CREATE TABLE authors (
    id   INT PRIMARY KEY NOT NULL,
    name CHAR(50)        NOT NULL,
    bio  TEXT
);
```

### users 表

```sql
CREATE TABLE users (
    id         INT PRIMARY KEY NOT NULL,
    name       CHAR(50)        NOT NULL,
    email      TEXT,
    created_at INT
);
```

## 配置说明

[sqlc.yaml](sqlc.yaml) 配置：

- **engine**: PostgreSQL
- **queries**: query.sql
- **schema**: schema.sql
- **gen/go/package**: db
- **gen/go/out**: db (生成的代码输出到当前目录)

## 注意事项

1. **参数占位符**: 使用 PostgreSQL 风格的 `$1`, `$2`, `$3`...
2. **返回类型**: 
   - `:one` - 返回单条记录
   - `:many` - 返回多条记录
   - `:exec` - 执行操作，无返回值
3. **生成的文件**: 运行 `sqlc generate` 后会在当前目录生成 `db.go`, `models.go`, `query.sql.go` 等文件

## 开发流程

1. 修改 `schema.sql` 添加/修改表结构
2. 修改 `query.sql` 添加/修改 SQL 查询
3. 运行 `sqlc generate` 生成新的 Go 代码
4. 在业务代码中使用生成的类型安全方法
