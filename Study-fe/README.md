# Study Questions Manager - 前端

基于 React + Vite 构建的学习题目管理系统前端。

## 功能特性

- 📝 题目管理：增删改查
- 🔍 搜索功能：支持题目、答案、分类搜索
- 📂 分类筛选：按分类查看题目
- 📱 响应式设计：支持移动端
- 🎨 现代化UI：渐变色彩、卡片设计
- ⚡ 实时预览：表单输入实时预览

## 技术栈

- React 18
- React Router v6
- Axios
- Vite
- CSS3 (原生样式)

## 快速开始

### 安装依赖
```bash
npm install
```

### 启动开发服务器
```bash
npm run dev
```

项目将在 http://localhost:3000 启动

### 构建生产版本
```bash
npm run build
```

## 项目结构

```
src/
├── pages/           # 页面组件
│   ├── QuestionList.jsx    # 题目列表页
│   ├── QuestionForm.jsx    # 添加题目页
│   └── QuestionEdit.jsx    # 编辑题目页
├── api.js          # API 接口封装
├── App.jsx         # 主应用组件
├── main.jsx        # 入口文件
└── index.css       # 全局样式
```

## API 接口

前端通过 axios 与后端 Go 服务交互：

- `GET /question` - 获取所有题目
- `GET /question/:category` - 按分类获取题目
- `POST /question` - 创建题目
- `PUT /question/:id` - 更新题目
- `DELETE /question/:id` - 删除题目

## 页面路由

- `/` - 题目列表页
- `/create` - 添加题目页
- `/edit/:id` - 编辑题目页
- `/category/:category` - 分类题目页

## 注意事项

1. 确保后端服务运行在 `http://localhost:8080`
2. 前端开发服务器默认运行在 `http://localhost:3000`
3. 已配置代理，前端请求会自动转发到后端
4. 支持热重载，修改代码后自动刷新页面

## 开发说明

- 使用 ES6+ 语法
- 采用函数式组件 + Hooks
- 响应式设计，支持移动端
- 统一的错误处理和用户反馈
- 代码结构清晰，易于维护 