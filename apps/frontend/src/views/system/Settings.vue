<template>
  <div class="system-settings">
    <div class="page-header">
      <h1>系统设置</h1>
    </div>

    <el-card class="settings-card">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本设置" name="basic">
          <el-form :model="basicSettings" label-width="120px">
            <el-form-item label="系统名称">
              <el-input v-model="basicSettings.systemName" placeholder="请输入系统名称" />
            </el-form-item>
            <el-form-item label="系统描述">
              <el-input
                v-model="basicSettings.systemDescription"
                type="textarea"
                :rows="3"
                placeholder="请输入系统描述"
              />
            </el-form-item>
            <el-form-item label="Logo URL">
              <el-input v-model="basicSettings.logoUrl" placeholder="请输入Logo图片URL" />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="安全设置" name="security">
          <el-form :model="securitySettings" label-width="120px">
            <el-form-item label="会话超时">
              <el-input-number v-model="securitySettings.sessionTimeout" :min="5" :max="1440" />
              <span style="margin-left: 10px">分钟</span>
            </el-form-item>
            <el-form-item label="密码强度">
              <el-select v-model="securitySettings.passwordStrength">
                <el-option label="低" value="low" />
                <el-option label="中" value="medium" />
                <el-option label="高" value="high" />
              </el-select>
            </el-form-item>
            <el-form-item label="两步验证">
              <el-switch v-model="securitySettings.twoFactorAuth" />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="通知设置" name="notification">
          <el-form :model="notificationSettings" label-width="120px">
            <el-form-item label="邮件通知">
              <el-switch v-model="notificationSettings.emailEnabled" />
            </el-form-item>
            <el-form-item label="SMTP 服务器">
              <el-input v-model="notificationSettings.smtpHost" placeholder="smtp.example.com" />
            </el-form-item>
            <el-form-item label="SMTP 端口">
              <el-input-number v-model="notificationSettings.smtpPort" :min="1" :max="65535" />
            </el-form-item>
            <el-form-item label="发送邮箱">
              <el-input v-model="notificationSettings.smtpFrom" placeholder="noreply@example.com" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <div class="form-actions">
        <el-button @click="handleReset">重置</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">保存设置</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

defineOptions({
  name: 'SystemSettings'
})

const activeTab = ref('basic')
const saving = ref(false)

const basicSettings = ref({
  systemName: '元数据管理平台',
  systemDescription: '',
  logoUrl: ''
})

const securitySettings = ref({
  sessionTimeout: 30,
  passwordStrength: 'medium',
  twoFactorAuth: false
})

const notificationSettings = ref({
  emailEnabled: false,
  smtpHost: '',
  smtpPort: 587,
  smtpFrom: ''
})

const handleSave = async () => {
  saving.value = true
  try {
    // TODO: 调用 API 保存设置
    await new Promise(resolve => setTimeout(resolve, 500))
    ElMessage.success('设置保存成功')
  } catch (error) {
    ElMessage.error('设置保存失败')
  } finally {
    saving.value = false
  }
}

const handleReset = () => {
  basicSettings.value = {
    systemName: '元数据管理平台',
    systemDescription: '',
    logoUrl: ''
  }
  securitySettings.value = {
    sessionTimeout: 30,
    passwordStrength: 'medium',
    twoFactorAuth: false
  }
  notificationSettings.value = {
    emailEnabled: false,
    smtpHost: '',
    smtpPort: 587,
    smtpFrom: ''
  }
  ElMessage.info('已重置为默认设置')
}
</script>

<style scoped>
.system-settings {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h1 {
  margin: 0;
  font-size: 24px;
  color: #303133;
}

.settings-card {
  margin-bottom: 20px;
}

.form-actions {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
  text-align: right;
}

.form-actions .el-button {
  margin-left: 10px;
}
</style>
