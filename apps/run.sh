#!/bin/bash

# 前后端统一运行脚本
# 支持：
#   ./run.sh start         - 启动前后端服务
#   ./run.sh start backend - 只启动后端服务
#   ./run.sh start frontend - 只启动前端服务
#   ./run.sh stop          - 停止前后端服务
#   ./run.sh status        - 查看服务状态
#   ./run.sh logs          - 查看服务日志

# 配置项
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/backend"
FRONTEND_DIR="$SCRIPT_DIR/frontend"
BACKEND_PORT=8080
FRONTEND_PORT=5173
BACKEND_LOG="$SCRIPT_DIR/logs/backend.log"
FRONTEND_LOG="$SCRIPT_DIR/logs/frontend.log"
PID_FILE="$SCRIPT_DIR/run.pid"

# 创建日志目录
mkdir -p "$SCRIPT_DIR/logs"

# 检查环境函数
check_env() {
    echo "=== 环境检查 ==="

    # 检查 Node.js 版本
    if command -v node &> /dev/null; then
        NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
        if [ "$NODE_VERSION" -lt 18 ]; then
            echo "❌ Node.js 版本过低 (当前: $(node -v), 要求: >= 18)"
            echo "   请升级 Node.js: nvm install 18 && nvm use 18"
            return 1
        else
            echo "✅ Node.js: $(node -v)"
        fi
    else
        echo "❌ 未安装 Node.js"
        return 1
    fi

    # 检查 Go
    if command -v go &> /dev/null; then
        echo "✅ Go: $(go version)"
    else
        echo "❌ 未安装 Go"
        return 1
    fi

    return 0
}

# 等待数据库就绪
wait_for_db() {
    echo "等待数据库连接..."
    local max_attempts=30
    local attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if nc -z localhost 3306 2>/dev/null; then
            echo "✅ 数据库已就绪"
            return 0
        fi
        attempt=$((attempt + 1))
        echo "  尝试 $attempt/$max_attempts..."
        sleep 1
    done

    echo "❌ 数据库连接超时"
    echo "   请启动 MySQL 或检查配置: $BACKEND_DIR/.env"
    return 1
}

# 函数：启动后端服务
start_backend() {
    echo "正在启动后端服务..."
    cd "$BACKEND_DIR" || exit 1

    # 确保 PID 文件存在
    touch "$PID_FILE"

    # 启动后端并保存 PID
    nohup go run cmd/main.go > "$BACKEND_LOG" 2>&1 &
    BACKEND_PID=$!

    # 验证进程是否启动成功
    sleep 1
    if ! ps -p "$BACKEND_PID" > /dev/null 2>&1; then
        echo "❌ 后端服务启动失败，请查看日志: $BACKEND_LOG"
        echo ""
        echo "最近 10 行日志:"
        tail -n 10 "$BACKEND_LOG"
        return 1
    fi

    echo "backend:$BACKEND_PID" >> "$PID_FILE"
    echo "✅ 后端服务启动成功，PID: $BACKEND_PID"
    echo "   日志: $BACKEND_LOG"
}

# 函数：启动前端服务
start_frontend() {
    echo "正在启动前端服务..."
    cd "$FRONTEND_DIR" || exit 1

    # 确保 PID 文件存在
    touch "$PID_FILE"

    # 检查 node_modules 是否存在
    if [ ! -d "node_modules" ]; then
        echo "首次启动，正在安装依赖..."
        if ! npm install; then
            echo "❌ 依赖安装失败"
            return 1
        fi
    fi

    # 启动前端
    nohup npm run dev > "$FRONTEND_LOG" 2>&1 &
    FRONTEND_PID=$!

    # 验证进程是否启动成功
    sleep 2
    if ! ps -p "$FRONTEND_PID" > /dev/null 2>&1; then
        echo "❌ 前端服务启动失败，请查看日志: $FRONTEND_LOG"
        echo ""
        echo "最近 10 行日志:"
        tail -n 10 "$FRONTEND_LOG"
        return 1
    fi

    echo "frontend:$FRONTEND_PID" >> "$PID_FILE"
    echo "✅ 前端服务启动成功，PID: $FRONTEND_PID"
    echo "   日志: $FRONTEND_LOG"
}

# 函数：停止后端服务
stop_backend() {
    if [ -f "$PID_FILE" ] && grep -q "backend:" "$PID_FILE"; then
        BACKEND_PID=$(grep "backend:" "$PID_FILE" | cut -d":" -f2)
        if ps -p "$BACKEND_PID" > /dev/null 2>&1; then
            echo "正在停止后端服务(PID: $BACKEND_PID)..."
            kill "$BACKEND_PID"
            sleep 2
            if ps -p "$BACKEND_PID" > /dev/null 2>&1; then
                echo "强制停止后端服务..."
                kill -9 "$BACKEND_PID"
            fi
        fi
        # 从PID文件中删除后端PID
        sed -i.bak "/backend:/d" "$PID_FILE" && rm -f "$PID_FILE.bak"
        echo "后端服务已停止"
    fi
}

# 函数：停止前端服务
stop_frontend() {
    if [ -f "$PID_FILE" ] && grep -q "frontend:" "$PID_FILE"; then
        FRONTEND_PID=$(grep "frontend:" "$PID_FILE" | cut -d":" -f2)
        if ps -p "$FRONTEND_PID" > /dev/null 2>&1; then
            echo "正在停止前端服务(PID: $FRONTEND_PID)..."
            kill "$FRONTEND_PID"
            sleep 2
            if ps -p "$FRONTEND_PID" > /dev/null 2>&1; then
                echo "强制停止前端服务..."
                kill -9 "$FRONTEND_PID"
            fi
        fi
        # 从PID文件中删除前端PID
        sed -i.bak "/frontend:/d" "$PID_FILE" && rm -f "$PID_FILE.bak"
        echo "前端服务已停止"
    fi
}

# 函数：查看服务状态
status() {
    echo "=== 服务状态 ==="

    # 检查后端服务
    if [ -f "$PID_FILE" ] && grep -q "backend:" "$PID_FILE"; then
        BACKEND_PID=$(grep "backend:" "$PID_FILE" | cut -d":" -f2)
        if ps -p "$BACKEND_PID" > /dev/null 2>&1; then
            echo "后端服务: 运行中 (PID: $BACKEND_PID)"
            echo "  访问地址: http://localhost:$BACKEND_PORT"
        else
            echo "后端服务: 已停止 (PID: $BACKEND_PID 已失效)"
            sed -i.bak "/backend:/d" "$PID_FILE" && rm -f "$PID_FILE.bak"
        fi
    else
        echo "后端服务: 未运行"
    fi

    # 检查前端服务
    if [ -f "$PID_FILE" ] && grep -q "frontend:" "$PID_FILE"; then
        FRONTEND_PID=$(grep "frontend:" "$PID_FILE" | cut -d":" -f2)
        if ps -p "$FRONTEND_PID" > /dev/null 2>&1; then
            echo "前端服务: 运行中 (PID: $FRONTEND_PID)"
            echo "  访问地址: http://localhost:$FRONTEND_PORT"
        else
            echo "前端服务: 已停止 (PID: $FRONTEND_PID 已失效)"
            sed -i.bak "/frontend:/d" "$PID_FILE" && rm -f "$PID_FILE.bak"
        fi
    else
        echo "前端服务: 未运行"
    fi
}

# 函数：查看日志
logs() {
    if [ "$1" = "backend" ]; then
        echo "=== 后端服务日志 ==="
        tail -f "$BACKEND_LOG"
    elif [ "$1" = "frontend" ]; then
        echo "=== 前端服务日志 ==="
        tail -f "$FRONTEND_LOG"
    else
        echo "=== 服务日志 ==="
        echo "查看后端日志: $0 logs backend"
        echo "查看前端日志: $0 logs frontend"
    fi
}

# 主逻辑
case "$1" in
    start)
        # 环境检查
        if ! check_env; then
            echo ""
            echo "环境检查失败，请先解决上述问题"
            exit 1
        fi
        echo ""

        case "$2" in
            backend)
                # 等待数据库就绪
                if ! wait_for_db; then
                    exit 1
                fi
                echo ""
                stop_backend
                if ! start_backend; then
                    exit 1
                fi
                sleep 1
                status
                ;;
            frontend)
                stop_frontend
                if ! start_frontend; then
                    exit 1
                fi
                sleep 1
                status
                ;;
            *)
                # 等待数据库就绪
                if ! wait_for_db; then
                    exit 1
                fi
                echo ""
                # 停止现有服务
                stop_backend
                stop_frontend
                # 启动所有服务
                if ! start_backend; then
                    exit 1
                fi
                # 等待后端启动完成
                sleep 3
                if ! start_frontend; then
                    exit 1
                fi
                sleep 1
                # 显示状态
                echo ""
                status
                ;;
        esac
        ;;
    stop)
        case "$2" in
            backend)
                stop_backend
                ;;
            frontend)
                stop_frontend
                ;;
            *)
                stop_backend
                stop_frontend
                echo "所有服务已停止"
                ;;
        esac
        ;;
    status)
        status
        ;;
    logs)
        logs "$2"
        ;;
    *)
        echo "用法: $0 {start|stop|status|logs} [backend|frontend]"
        echo "  start    - 启动服务"
        echo "  stop     - 停止服务"
        echo "  status   - 查看服务状态"
        echo "  logs     - 查看服务日志"
        exit 1
        ;;
esac

exit 0
