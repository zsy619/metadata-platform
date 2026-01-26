<template>
    <div class="login-container">
        <!-- 顶部装饰条 (特定于登录页) -->
        <div class="top-stripe" role="presentation"></div>
        <!-- 登录卡片 (复合 card-container 规范) -->
        <div class="login-card" role="main">
            <!-- 品牌头部 -->
            <header class="brand-header">
                <div class="accent-logo-box flex-center" aria-hidden="true">
                    <el-icon class="logo-icon">
                        <DataAnalysis />
                    </el-icon>
                </div>
                <h1 class="brand-title m-t-md">元数据平台</h1>
            </header>
            <!-- 欢迎文字 -->
            <section class="welcome-section m-t-lg">
                <h2 class="welcome-title text-primary">欢迎登录</h2>
                <p class="welcome-subtitle text-secondary">高效、安全的企业级数据管理系统</p>
            </section>
            <!-- 登录表单 -->
            <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" class="login-form m-t-lg" @keyup.enter="handleLogin" label-position="top">
                <el-form-item prop="username">
                    <el-input v-model="loginForm.username" placeholder="请输入用户名" size="large" aria-label="用户名">
                        <template #prefix>
                            <el-icon>
                                <User />
                            </el-icon>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item prop="password">
                    <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" size="large" show-password aria-label="密码">
                        <template #prefix>
                            <el-icon>
                                <Lock />
                            </el-icon>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item prop="captcha">
                    <div class="captcha-row">
                        <el-input v-model="loginForm.captcha" placeholder="验证码" size="large" class="captcha-input" aria-label="验证码">
                            <template #prefix>
                                <el-icon>
                                    <View />
                                </el-icon>
                            </template>
                        </el-input>
                        <div class="captcha-image" @click="refreshCaptcha" title="点击刷新验证码" role="button" tabindex="0" @keyup.enter="refreshCaptcha">
                            <img :src="captchaImage" alt="验证码图片">
                        </div>
                    </div>
                </el-form-item>
                <div class="form-footer flex-between m-b-lg">
                    <el-checkbox v-model="loginForm.remember">记住我</el-checkbox>
                    <el-link :underline="false" type="primary" @click="showForgotDialog = true">
                        忘记密码？
                    </el-link>
                </div>
                <el-button type="primary" size="large" class="login-button w-full" :loading="loading" @click="handleLogin">
                    {{ loading ? '登录中...' : '立即登录' }}
                </el-button>
            </el-form>
            <!-- 第三方登录 -->
            <footer class="social-section m-t-xl">
                <div class="divider">
                    <span class="text-secondary">或使用第三方账号登录</span>
                </div>
                <div class="social-icons m-t-lg">
                    <button type="button" class="social-btn wechat" @click="handleSocialLogin('wechat')" aria-label="微信登录">
                        <el-icon>
                            <ChatDotRound />
                        </el-icon>
                    </button>
                    <button type="button" class="social-btn qq" @click="handleSocialLogin('qq')" aria-label="QQ登录">
                        <el-icon>
                            <Opportunity />
                        </el-icon>
                    </button>
                    <button type="button" class="social-btn dingtalk" @click="handleSocialLogin('dingtalk')" aria-label="钉钉登录">
                        <el-icon>
                            <Share />
                        </el-icon>
                    </button>
                </div>
            </footer>
            <!-- 注册链接 -->
            <footer class="register-section m-t-xl" style="padding-top: 16px;">
                <span class="text-secondary">还没有账号？</span>
                <el-link :underline="false" style="color: var(--el-color-accent); font-weight: 700; font-size: 15px;" @click="showRegisterDialog = true">
                    免费注册
                </el-link>
            </footer>
        </div>
        <!-- 底部版权 -->
        <div class="login-footer m-t-lg">
            <p>© 2026 元数据管理平台 Copyright. All Rights Reserved.</p>
        </div>
        <!-- 弹窗：注册 -->
        <el-dialog v-model="showRegisterDialog" title="注册新账号" width="460px" :close-on-click-modal="false" destroy-on-close class="custom-dialog">
            <el-form ref="registerFormRef" :model="registerForm" :rules="registerRules" label-position="top">
                <el-form-item label="用户名" prop="username">
                    <el-input v-model="registerForm.username" placeholder="建议使用字母数字组合" />
                </el-form-item>
                <el-form-item label="邮箱" prop="email">
                    <el-input v-model="registerForm.email" placeholder="接收验证信息的邮箱" />
                </el-form-item>
                <el-form-item label="密码" prop="password">
                    <el-input v-model="registerForm.password" type="password" show-password placeholder="不少于6位" />
                </el-form-item>
                <el-form-item label="确认密码" prop="confirmPassword">
                    <el-input v-model="registerForm.confirmPassword" type="password" show-password placeholder="请再次输入密码" />
                </el-form-item>
            </el-form>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="showRegisterDialog = false">取消</el-button>
                    <el-button type="primary" @click="handleRegister" :loading="loading">立即注册</el-button>
                </div>
            </template>
        </el-dialog>
        <!-- 弹窗：重置密码 -->
        <el-dialog v-model="showForgotDialog" title="重置密码" width="420px" :close-on-click-modal="false" class="custom-dialog">
            <el-form label-position="top">
                <el-form-item label="注册邮箱">
                    <el-input placeholder="请输入您的注册邮箱" />
                </el-form-item>
                <el-form-item label="验证码">
                    <div class="captcha-row">
                        <el-input placeholder="6位验证码" />
                        <el-button type="primary">发送验证码</el-button>
                    </div>
                </el-form-item>
            </el-form>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="showForgotDialog = false">取消</el-button>
                    <el-button type="primary">确认重置</el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { getCaptchaApi, loginApi, registerApi } from '@/api/auth'
import {
    ChatDotRound,
    DataAnalysis,
    Lock,
    Opportunity, Share,
    User,
    View
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loginFormRef = ref()
const registerFormRef = ref()
const loading = ref(false)
const showRegisterDialog = ref(false)
const showForgotDialog = ref(false)
const captchaImage = ref('')
const captchaId = ref('')

const loginForm = reactive({ username: '', password: '', captcha: '', remember: false })
const loginRules = reactive({
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
    captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
})

const registerForm = reactive({ username: '', email: '', password: '', confirmPassword: '' })
const registerRules = reactive({
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    email: [{ required: true, type: 'email', message: '请输入正确的邮箱', trigger: 'blur' }],
    password: [{ required: true, min: 6, message: '密码不能少于6位', trigger: 'blur' }],
    confirmPassword: [
        { required: true, message: '请再次输入密码', trigger: 'blur' },
        { validator: (_rule: any, value: string, cb: any) => value !== registerForm.password ? cb(new Error('两次输入密码不一致')) : cb(), trigger: 'blur' }
    ]
})

const refreshCaptcha = async () => {
    try {
        const res = await getCaptchaApi()
        captchaImage.value = res.data.pic_path
        captchaId.value = res.data.captcha_id
    } catch (err) {
        ElMessage.error('加载验证码失败')
    }
}

const handleLogin = async () => {
    if (!loginFormRef.value || loading.value) return
    try {
        await loginFormRef.value.validate()
        loading.value = true
        const res = await loginApi(loginForm.username, loginForm.password, captchaId.value, loginForm.captcha)
        localStorage.setItem('token', res.data?.access_token || '')
        loginForm.remember ? localStorage.setItem('username', loginForm.username) : localStorage.removeItem('username')
        ElMessage.success('登录成功，欢迎回来')
        router.push('/')
    } catch (err: any) {
        ElMessage.error(err.message || '登录失败')
    } finally {
        loading.value = false
    }
}

const handleRegister = async () => {
    if (!registerFormRef.value) return
    try {
        await registerFormRef.value.validate()
        loading.value = true
        await registerApi(registerForm)
        ElMessage.success('注册成功，请登录')
        showRegisterDialog.value = false
    } catch (err: any) {
        ElMessage.error('注册失败')
    } finally {
        loading.value = false
    }
}

const handleSocialLogin = (type: string) => ElMessage.info(`正在连接${type}授权服务...`)

onMounted(() => {
    const user = localStorage.getItem('username')
    if (user) { loginForm.username = user; loginForm.remember = true }
    refreshCaptcha()
})
</script>
<style scoped>
.login-container {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    position: relative;
    overflow: hidden;
}

.top-stripe {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 360px;
    background-color: var(--el-color-primary);
    background-image: linear-gradient(135deg, var(--el-color-primary) 0%, var(--el-color-primary-dark-2) 100%);
    z-index: 1;
}

.login-card {
    width: 100%;
    max-width: 440px;
    background: var(--el-bg-color);
    border-radius: var(--el-border-radius-medium);
    box-shadow: var(--el-box-shadow-dark);
    padding: 48px;
    position: relative;
    z-index: 10;
    animation: cardIn 0.6s var(--el-transition-function-ease-in-out-bezier);
}

@keyframes cardIn {
    from {
        opacity: 0;
        transform: translateY(20px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.brand-header {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.accent-logo-box {
    width: 64px;
    height: 64px;
    transition: transform 0.5s ease;
}

.accent-logo-box:hover {
    transform: rotate(10deg) scale(1.1);
}

.logo-icon {
    font-size: 36px;
}

.welcome-section {
    text-align: center;
}

.captcha-row {
    display: flex;
    gap: 12px;
    align-items: center;
}

.captcha-image {
    width: 120px;
    height: 40px;
    border-radius: var(--el-border-radius-base);
    overflow: hidden;
    cursor: pointer;
    border: 1px solid var(--el-border-color);
    flex-shrink: 0;
}

.captcha-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

/* 后端 social-section 和 divider 样式由 global components.css 驱动 */
.social-section {
    text-align: center;
}

.social-icons {
    display: flex;
    justify-content: center;
    gap: 20px;
}

.social-btn {
    width: 44px;
    height: 44px;
    border: 1px solid var(--el-border-color);
    border-radius: 50%;
    background: var(--el-bg-color);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all var(--el-transition-duration);
    color: var(--el-text-color-regular);
    font-size: 22px;
}

.social-btn:hover {
    color: #FFFFFF;
    border-color: transparent;
    transform: translateY(-2px) scale(1.1);
}

.social-btn.wechat:hover {
    background-color: #07C160;
    box-shadow: 0 2px 8px rgba(7, 193, 96, 0.3);
}

.social-btn.qq:hover {
    background-color: #409EFF;
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.social-btn.dingtalk:hover {
    background-color: #0077FA;
    box-shadow: 0 2px 8px rgba(0, 119, 250, 0.3);
}

.register-section {
    text-align: center;
    font-size: 14px;
}

.login-footer {
    text-align: center;
    color: rgba(255, 255, 255, 0.7);
    font-size: 12px;
}

.flex-between {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.w-full {
    width: 100%;
}

@media (max-width: 768px) {
    .top-stripe {
        height: 180px;
    }

    .login-card {
        padding: 32px 24px;
        margin: 0 16px;
        box-shadow: var(--el-box-shadow);
    }

    .login-footer {
        color: var(--el-text-color-secondary);
    }
}
</style>
