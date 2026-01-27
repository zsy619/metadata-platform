<template>
    <div class="auth-wrapper">
        <!-- 背景装饰层 -->
        <div class="auth-bg">
            <div class="bg-shape shape-1"></div>
            <div class="bg-shape shape-2"></div>
            <div class="bg-shape shape-3"></div>
        </div>
        <div class="auth-container">
            <div class="auth-content">
                <router-view />
            </div>
            <footer class="auth-footer">
                <p>© 2026 元数据管理平台 · 企业级数据资产治理中心</p>
                <div class="footer-links">
                    <span>隐私政策</span>
                    <span class="divider">|</span>
                    <span>服务条款</span>
                    <span class="divider">|</span>
                    <span>帮助文档</span>
                </div>
            </footer>
        </div>
    </div>
</template>
<script setup lang="ts">
</script>
<style scoped>
.auth-wrapper {
    position: relative;
    min-height: 100vh;
    width: 100%;
    background: #0f172a;
    overflow: hidden;
    display: flex;
    justify-content: center;
}

/* 动态背景 */
.auth-bg {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 1;
    overflow: hidden;
}

.bg-shape {
    position: absolute;
    filter: blur(80px);
    opacity: 0.4;
    border-radius: 50%;
    animation: move 20s infinite alternate;
}

.shape-1 {
    width: 500px;
    height: 500px;
    background: #3b82f6;
    top: -100px;
    left: -100px;
}

.shape-2 {
    width: 400px;
    height: 400px;
    background: #8b5cf6;
    bottom: -50px;
    right: -50px;
    animation-duration: 25s;
}

.shape-3 {
    width: 300px;
    height: 300px;
    background: #0ea5e9;
    top: 40%;
    right: 15%;
    animation-duration: 30s;
}

@keyframes move {
    from {
        transform: translate(0, 0);
    }

    to {
        transform: translate(100px, 50px) scale(1.1);
    }
}

.auth-container {
    position: relative;
    z-index: 10;
    width: 100%;
    max-width: 500px;
    display: flex;
    flex-direction: column;
    padding: 40px 20px;
}

.auth-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
}

.auth-footer {
    text-align: center;
    color: rgba(255, 255, 255, 0.45);
    font-size: 13px;
    padding-top: 40px;
}

.footer-links {
    margin-top: 12px;
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 12px;
}

.footer-links span {
    cursor: pointer;
    transition: color 0.3s;
}

.footer-links span:hover:not(.divider) {
    color: #3b82f6;
}

.divider {
    color: rgba(255, 255, 255, 0.2);
    cursor: default !important;
}

@media (max-width: 640px) {
    .auth-container {
        padding: 20px 15px;
    }
}
</style>
