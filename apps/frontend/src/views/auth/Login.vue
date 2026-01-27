<template>
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
        <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" class="login-form m-t-lg" @keyup.enter="handleLogin" label-position="top" v-loading="loading" element-loading-text="登录中...">
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
                    <div class="captcha-image" @click="refreshCaptcha" title="点击刷新验证码" role="button" tabindex="0">
                        <img :src="captchaImage" alt="验证码图片">
                    </div>
                </div>
            </el-form-item>
            <div class="form-footer flex-between m-b-lg">
                <el-checkbox v-model="loginForm.remember">记住我</el-checkbox>
                <el-link underline="never" type="primary" @click="showForgotDialog = true">
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
                <button type="button" class="social-btn wechat" @click="handleSocialLogin('wechat')"><el-icon>
                        <ChatDotRound />
                    </el-icon></button>
                <button type="button" class="social-btn qq" @click="handleSocialLogin('qq')"><el-icon>
                        <Opportunity />
                    </el-icon></button>
                <button type="button" class="social-btn dingtalk" @click="handleSocialLogin('dingtalk')"><el-icon>
                        <Share />
                    </el-icon></button>
            </div>
        </footer>
        <!-- 注册链接 -->
        <footer class="register-section m-t-xl">
            <span class="text-secondary">还没有账号？</span>
            <el-link underline="never" class="register-link" @click="showRegisterDialog = true">
                免费注册
            </el-link>
        </footer>
    </div>
    <!-- 弹窗：注册 -->
    <el-dialog v-model="showRegisterDialog" title="注册新账号" width="460px" :close-on-click-modal="false" destroy-on-close>
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
            <el-button @click="showRegisterDialog = false">取消</el-button>
            <el-button type="primary" @click="handleRegister" :loading="loading">立即注册</el-button>
        </template>
    </el-dialog>
    <!-- 弹窗：重置密码 -->
    <el-dialog v-model="showForgotDialog" title="重置密码" width="420px" :close-on-click-modal="false">
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
            <el-button @click="showForgotDialog = false">取消</el-button>
            <el-button type="primary">确认重置</el-button>
        </template>
    </el-dialog>
</template>
<script setup lang="ts">
import { getCaptchaApi, registerApi } from '@/api/auth'
import { useUserStore } from '@/stores/user'
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
const userStore = useUserStore()
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

        await userStore.login({
            account: loginForm.username,
            password: loginForm.password,
            captcha_id: captchaId.value,
            captcha_code: loginForm.captcha
        })

        loginForm.remember ? localStorage.setItem('username', loginForm.username) : localStorage.removeItem('username')
        ElMessage.success('登录成功，欢迎回来')

        const redirect = router.currentRoute.value.query.redirect as string || '/'
        router.push(redirect)
    } catch (err: any) {
        ElMessage.error(err.message || '登录失败')
        refreshCaptcha()
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
.login-card {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-radius: 16px;
    box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.07);
    padding: 40px;
    width: 100%;
    max-width: 440px;
    margin: 0 20px;
    animation: cardIn 0.5s ease-out;
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
    margin-bottom: 24px;
}

.accent-logo-box {
    width: 56px;
    height: 56px;
    background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
    border-radius: 12px;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.2);
}

.logo-icon {
    font-size: 32px;
}

.brand-title {
    font-size: 22px;
    font-weight: 700;
    color: #1e293b;
    margin-top: 12px;
}

.welcome-section {
    text-align: center;
    margin-bottom: 32px;
}

.welcome-title {
    font-size: 20px;
    font-weight: 600;
    color: #334155;
}

.welcome-subtitle {
    font-size: 14px;
    color: #64748b;
    margin-top: 8px;
}

.captcha-row {
    display: flex;
    gap: 12px;
    width: 100%;
}

.captcha-image {
    flex-shrink: 0;
    width: 120px;
    height: 40px;
    border-radius: 4px;
    border: 1px solid #e2e8f0;
    cursor: pointer;
    overflow: hidden;
}

.captcha-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.social-section {
    text-align: center;
}

.divider {
    position: relative;
    margin: 24px 0;
}

.divider::before,
.divider::after {
    content: "";
    position: absolute;
    top: 50%;
    width: 25%;
    height: 1px;
    background: #e2e8f0;
}

.divider::before {
    left: 0;
}

.divider::after {
    right: 0;
}

.social-icons {
    display: flex;
    justify-content: center;
    gap: 16px;
}

.social-btn {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    border: 1px solid #e2e8f0;
    background: white;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    color: #64748b;
    transition: all 0.3s;
}

.social-btn:hover {
    transform: translateY(-2px);
    border-color: #3b82f6;
    color: #3b82f6;
}

.register-section {
    text-align: center;
    margin-top: 24px;
    font-size: 14px;
}

.register-link {
    color: #3b82f6;
    font-weight: 600;
    margin-left: 4px;
}

.flex-between {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.w-full {
    width: 100%;
}

.m-t-lg {
    margin-top: 20px;
}

.m-t-xl {
    margin-top: 32px;
}

.m-t-md {
    margin-top: 12px;
}

.m-b-lg {
    margin-bottom: 20px;
}
</style>
