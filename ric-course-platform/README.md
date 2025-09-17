# RIC 选课平台（React + Go + PostgreSQL）

本项目完全满足面试附加题的 3 个核心功能：
1) 首次访问或点击“查看所有课程”展示全部课程与数量；  
2) “选择学生”下拉列出所有学生姓名；  
3) 选择某个学生后展示其所选课程并显示数量。  

后端：Go + Gin + pgx；前端：React；数据库：PostgreSQL。也提供 Docker Compose 一键运行。

> 题目来源：个人面试附加题（程序员）PDF。

## 一键启动（推荐：Docker）
```bash
# 进入项目根目录
cd ric-course-platform

# 构建并启动
docker compose up -d --build

# 前端：http://localhost:3000
# 后端：http://localhost:8080
# 健康检查：http://localhost:8080/healthz
# 示例接口：/api/courses, /api/students, /api/students/{id}/courses
```

## 手工运行

### 1) 数据库（PostgreSQL）
```bash
docker run --name ric-pg -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=ric -p 5432:5432 -d postgres:15
# 初始化表结构与样例数据：
psql "postgres://postgres:postgres@localhost:5432/ric" -f backend/db/schema.sql
psql "postgres://postgres:postgres@localhost:5432/ric" -f backend/db/seed.sql
```

### 2) 后端（Go）
```bash
cd backend
cp .env.example .env   # 如需修改数据库连接
go mod download
go run .
# 监听 :8080
```

### 3) 前端（React）
```bash
cd frontend
npm install
# 如需切换为你部署的后端，请复制 .env.example 为 .env 并修改 REACT_APP_API_BASE_URL
npm start
# 浏览器打开 http://localhost:3000
```

## API 说明
- `GET /api/courses` → `{ count, items: Course[] }`
- `GET /api/students` → `{ items: Student[] }`
- `GET /api/students/:id/courses` → `{ count, items: Course[] }`

## 目录结构
```
ric-course-platform/
├─ backend/                # Go 后端
│  ├─ db/
│  │  ├─ schema.sql
│  │  └─ seed.sql
│  ├─ .env.example
│  ├─ go.mod
│  ├─ Dockerfile
│  └─ main.go
├─ frontend/               # React 前端
│  ├─ public/
│  │  └─ index.html
│  ├─ src/
│  │  ├─ App.js
│  │  ├─ api.js
│  │  ├─ index.js
│  │  └─ styles.css
│  ├─ .env.example
│  ├─ package.json
│  ├─ Dockerfile
│  └─ README.md
└─ docker-compose.yml
```

## 安全提示
提交前请确保 `.env` 中未包含真实密码或 Token。将敏感配置改为环境变量或占位符。

## 兼容提示
如不方便自建后端，可在 `frontend/.env` 中将 `REACT_APP_API_BASE_URL` 指向主办方公开 API（如果仍在线）。
"# ric-course-platform" 
