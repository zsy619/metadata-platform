<template>
    <div class="navbar">
        <div class="hamburger-container" @click="toggleSideBar">
            <el-icon :class="{ 'is-active': sidebar.opened }" class="hamburger">
                <Fold v-if="sidebar.opened" />
                <Expand v-else />
            </el-icon>
        </div>
        <breadcrumb class="breadcrumb-container" />
        <div class="right-menu">
            <div class="right-menu-item">
                <el-dropdown trigger="click" @command="handleLangCommand">
                    <span class="el-dropdown-link">
                        <el-icon class="right-icon">
                            <i class="fas fa-language"></i> <!-- Or just Text -->
                            <span style="font-size: 14px; margin-left: 4px">{{ currentLang === 'zh-CN' ? '中文' : 'En' }}</span>
                        </el-icon>
                    </span>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item command="zh-CN" :disabled="currentLang === 'zh-CN'">中文</el-dropdown-item>
                            <el-dropdown-item command="en-US" :disabled="currentLang === 'en-US'">English</el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </div>
            <div class="right-menu-item">
                <el-tooltip :content="$t('header.fullScreen')" effect="dark" placement="bottom">
                    <el-icon class="right-icon" @click="toggleFullScreen">
                        <FullScreen />
                    </el-icon>
                </el-tooltip>
            </div>
            <div class="right-menu-item">
                <el-tooltip :content="$t('theme.switch')" effect="dark" placement="bottom">
                    <el-icon class="right-icon" @click="toggleTheme">
                        <Moon v-if="isDark" />
                        <Sunny v-else />
                    </el-icon>
                </el-tooltip>
            </div>
            <el-dropdown class="avatar-container" trigger="click" @command="handleCommand">
                <div class="avatar-wrapper">
                    <img :src="avatar + '?imageView2/1/w/80/h/80'" class="user-avatar">
                    <span class="user-name">{{ userName }}</span>
                    <el-icon class="el-icon--right">
                        <CaretBottom />
                    </el-icon>
                </div>
                <template #dropdown>
                    <el-dropdown-menu class="user-dropdown">
                        <router-link to="/">
                            <el-dropdown-item>首页</el-dropdown-item>
                        </router-link>
                        <router-link to="/user/profile">
                            <el-dropdown-item>{{ $t('header.profile') }}</el-dropdown-item>
                        </router-link>
                        <el-dropdown-item divided command="logout">
                            <span style="display:block;">{{ $t('header.logout') }}</span>
                        </el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </div>
    </div>
</template>
<script setup lang="ts">
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { CaretBottom, Expand, Fold, FullScreen, Moon, Sunny } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import Breadcrumb from './Breadcrumb.vue'

const appStore = useAppStore()
const userStore = useUserStore()
const router = useRouter()

const sidebar = computed(() => appStore.sidebar)
const avatar = computed(() => userStore.avatar)
const userName = computed(() => userStore.userName)
const isDark = computed(() => appStore.theme === 'dark')
const currentLang = computed(() => appStore.language)

const toggleSideBar = () => {
    appStore.toggleSidebar()
}

const toggleFullScreen = () => {
    if (!document.fullscreenElement) {
        document.documentElement.requestFullscreen()
    } else {
        if (document.exitFullscreen) {
            document.exitFullscreen()
        }
    }
}

const toggleTheme = () => {
    const newTheme = isDark.value ? 'light' : 'dark'
    appStore.setTheme(newTheme)
}

const handleLangCommand = (lang: string) => {
    appStore.setLanguage(lang)
    ElMessage.success(lang === 'zh-CN' ? '已切换为中文' : 'Switched to English')
}

const handleCommand = (command: string) => {
    if (command === 'logout') {
        logout()
    }
}

const logout = async () => {
    ElMessageBox.confirm('确定注销并退出系统吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(async () => {
        await userStore.logout()
        router.push(`/login?redirect=${router.currentRoute.value.fullPath}`)
        ElMessage.success('退出成功')
    }).catch(() => { })
}
</script>
<style scoped>
.navbar {
    height: 50px;
    overflow: hidden;
    position: relative;
    background: #fff;
    box-shadow: 0 1px 4px rgba(0, 21, 41, .08);
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background .3s;
    -webkit-tap-highlight-color: transparent;
    display: flex;
    align-items: center;
    padding: 0 15px;
}

.hamburger-container:hover {
    background: rgba(0, 0, 0, .025)
}

.hamburger {
    font-size: 20px;
}

.breadcrumb-container {
    float: left;
    flex: 1;
}

.right-menu {
    float: right;
    height: 100%;
    line-height: 50px;
    display: flex;
    align-items: center;
}

.right-menu:focus {
    outline: none;
}

.avatar-container {
    margin-right: 30px;
}

.avatar-wrapper {
    margin-top: 5px;
    position: relative;
    cursor: pointer;
    display: flex;
    align-items: center;
}

.user-avatar {
    cursor: pointer;
    width: 40px;
    height: 40px;
    border-radius: 10px;
    margin-right: 10px;
}

.user-name {
    font-size: 14px;
    color: #333;
    margin-right: 5px;
}

.right-menu-item {
    display: inline-block;
    padding: 0 8px;
    height: 100%;
    font-size: 18px;
    color: #5a5e66;
    vertical-align: text-bottom;
    cursor: pointer;
    transition: background .3s;
}

.right-menu-item:hover {
    background: rgba(0, 0, 0, .025)
}

.right-icon {
    font-size: 20px;
    vertical-align: middle;
}
</style>
