<template>
  <el-popover
    placement="bottom-start"
    :width="420"
    trigger="click"
    v-model:visible="popoverVisible"
  >
    <template #reference>
      <el-input
        v-model="displayValue"
        placeholder="请选择图标"
        readonly
        clearable
        @clear="handleClear"
      >
        <template #prefix>
          <font-awesome-icon :icon="modelValue" v-if="modelValue" />
        </template>
        <template #suffix>
          <font-awesome-icon icon="fa-solid fa-chevron-down" />
        </template>
      </el-input>
    </template>
    <div class="icon-picker">
      <div class="icon-search">
        <el-input
          v-model="searchText"
          :placeholder="`搜索图标 (共 ${totalCount} 个)`"
          clearable
          prefix-icon="Search"
        />
      </div>
      <div class="icon-tabs">
        <el-radio-group v-model="activeStyle" size="small">
          <el-radio-button value="solid">实心 ({{ solidCount }})</el-radio-button>
          <el-radio-button value="regular">常规 ({{ regularCount }})</el-radio-button>
          <el-radio-button value="brands">品牌 ({{ brandsCount }})</el-radio-button>
        </el-radio-group>
      </div>
      <div class="icon-category">
        <el-select v-model="activeCategory" size="small" placeholder="选择分类" clearable>
          <el-option v-for="cat in currentCategories" :key="cat.value" :label="cat.label" :value="cat.value" />
        </el-select>
      </div>
      <div class="icon-list">
        <div
          v-for="icon in filteredIcons"
          :key="icon"
          class="icon-item"
          :class="{ active: icon === modelValue }"
          @click="handleSelect(icon)"
          :title="icon"
        >
          <font-awesome-icon :icon="icon" />
        </div>
      </div>
      <div class="icon-footer">
        <span>显示 {{ filteredIcons.length }} 个图标</span>
      </div>
      <el-empty v-if="filteredIcons.length === 0" description="没有找到图标" :image-size="60" />
    </div>
  </el-popover>
</template>

<script setup lang="ts">
import { fab } from '@fortawesome/free-brands-svg-icons';
import { far } from '@fortawesome/free-regular-svg-icons';
import { fas } from '@fortawesome/free-solid-svg-icons';
import { computed, ref, watch } from 'vue';

const props = defineProps<{
  modelValue?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | undefined): void
}>()

const popoverVisible = ref(false)
const searchText = ref('')
const activeStyle = ref('solid')
const activeCategory = ref('')

const convertToFaFormat = (iconObj: any, prefix: string) => {
  if (iconObj && iconObj.iconName) {
    return `${prefix} fa-${iconObj.iconName}`
  }
  return null
}

const solidIcons: string[] = []
const regularIcons: string[] = []
const brandsIcons: string[] = []

Object.entries(fas).forEach(([key, icon]: [string, any]) => {
  if (key.startsWith('fa') && key.length > 2 && icon && icon.iconName) {
    solidIcons.push(convertToFaFormat(icon, 'fa-solid')!)
  }
})

Object.entries(far).forEach(([key, icon]: [string, any]) => {
  if (key.startsWith('fa') && key.length > 2 && icon && icon.iconName) {
    regularIcons.push(convertToFaFormat(icon, 'fa-regular')!)
  }
})

Object.entries(fab).forEach(([key, icon]: [string, any]) => {
  if (key.startsWith('fa') && key.length > 2 && icon && icon.iconName) {
    brandsIcons.push(convertToFaFormat(icon, 'fa-brands')!)
  }
})

const solidCount = solidIcons.length
const regularCount = regularIcons.length
const brandsCount = brandsIcons.length
const totalCount = solidCount + regularCount + brandsCount

const categories = {
  solid: [
    { value: 'arrows', label: '箭头方向' },
    { value: 'business', label: '商业办公' },
    { value: 'communication', label: '通讯交流' },
    { value: 'education', label: '教育学习' },
    { value: 'medical', label: '医疗健康' },
    { value: 'music', label: '音乐媒体' },
    { value: 'security', label: '安全锁' },
    { value: 'users', label: '用户人物' },
    { value: 'files', label: '文件文档' },
    { value: 'shopping', label: '购物支付' },
    { value: 'interface', label: '界面元素' },
    { value: 'coding', label: '代码开发' },
  ],
  regular: [
    { value: 'arrows', label: '箭头方向' },
    { value: 'business', label: '商业办公' },
    { value: 'communication', label: '通讯交流' },
    { value: 'files', label: '文件文档' },
    { value: 'interface', label: '界面元素' },
    { value: 'users', label: '用户人物' },
  ],
  brands: [
    { value: 'social', label: '社交媒体' },
    { value: 'tech', label: '科技平台' },
    { value: 'payment', label: '支付平台' },
    { value: 'browser', label: '浏览器' },
    { value: 'coding', label: '开发工具' },
    { value: 'design', label: '设计工具' },
  ]
}

const currentCategories = computed(() => categories[activeStyle.value as keyof typeof categories] || [])

const categoryKeywords: Record<string, string[]> = {
  arrows: ['arrow', 'angle', 'chevron', 'caret', 'turn', 'direction', 'up', 'down', 'left', 'right', 'rotate', 'sync', 'refresh', 'exchange', 'expand', 'compress', 'arrows'],
  business: ['briefcase', 'building', 'business', 'chart', 'coin', 'money', 'dollar', 'credit', 'bank', 'wallet', 'receipt', 'invoice', 'handshake', 'suitcase', 'store', 'shop', 'industry', 'factory'],
  communication: ['comment', 'message', 'envelope', 'mail', 'phone', 'chat', 'bell', 'notification', 'inbox', 'paper-plane', 'at', 'address', 'contact', 'speech'],
  education: ['book', 'graduation', 'school', 'university', 'lesson', 'pencil', 'pen', 'note', 'journal', 'diploma', 'certificate', 'learn', 'teach', 'student'],
  medical: ['hospital', 'medical', 'health', 'heart', 'pulse', 'prescription', 'pills', 'capsule', 'stethoscope', 'user-doctor', 'nurse', 'ambulance', 'first-aid', 'thermometer'],
  music: ['music', 'play', 'pause', 'stop', 'forward', 'backward', 'volume', 'speaker', 'headphones', 'microphone', 'radio', 'podcast', 'record', 'guitar', 'drum'],
  security: ['lock', 'unlock', 'key', 'shield', 'protect', 'secure', 'password', 'fingerprint', 'eye', 'visibility', 'camera', 'surveillance', 'guard', 'safety'],
  users: ['user', 'person', 'people', 'team', 'group', 'family', 'child', 'man', 'woman', 'avatar', 'profile', 'account', 'id', 'identity'],
  files: ['file', 'folder', 'document', 'archive', 'zip', 'pdf', 'excel', 'word', 'powerpoint', 'image', 'photo', 'video', 'audio', 'download', 'upload', 'cloud', 'disk', 'drive', 'save'],
  shopping: ['cart', 'shop', 'store', 'bag', 'basket', 'tag', 'price', 'sale', 'discount', 'gift', 'coupon', 'receipt', 'credit-card', 'payment', 'checkout', 'product'],
  interface: ['menu', 'bars', 'list', 'grid', 'table', 'dashboard', 'panel', 'sidebar', 'navbar', 'button', 'toggle', 'switch', 'slider', 'progress', 'spinner', 'loading', 'plus', 'minus', 'times', 'check', 'search', 'filter', 'sort', 'settings', 'cog', 'gear', 'tools', 'wrench'],
  coding: ['code', 'terminal', 'console', 'command', 'git', 'github', 'branch', 'merge', 'pull', 'push', 'commit', 'database', 'server', 'cloud', 'api', 'bug', 'debug', 'laptop', 'computer', 'desktop', 'monitor'],
  social: ['facebook', 'twitter', 'instagram', 'linkedin', 'youtube', 'tiktok', 'snapchat', 'pinterest', 'reddit', 'tumblr', 'whatsapp', 'telegram', 'discord', 'wechat', 'weibo', 'qq'],
  tech: ['apple', 'google', 'microsoft', 'amazon', 'android', 'windows', 'linux', 'ubuntu', 'fedora', 'centos', 'docker', 'kubernetes', 'aws', 'azure', 'firebase'],
  payment: ['paypal', 'stripe', 'visa', 'mastercard', 'alipay', 'wechat-pay', 'apple-pay', 'google-pay', 'bitcoin', 'ethereum'],
  browser: ['chrome', 'firefox', 'safari', 'edge', 'opera', 'ie', 'browser', 'internet-explorer'],
  design: ['figma', 'sketch', 'adobe', 'photoshop', 'illustrator', 'xd', 'canva', 'invision', 'behance', 'dribbble'],
}

const currentIcons = computed(() => {
  switch (activeStyle.value) {
    case 'solid':
      return solidIcons
    case 'regular':
      return regularIcons
    case 'brands':
      return brandsIcons
    default:
      return solidIcons
  }
})

const filteredIcons = computed(() => {
  let icons = currentIcons.value
  
  if (activeCategory.value) {
    const keywords = categoryKeywords[activeCategory.value] || []
    icons = icons.filter(icon => {
      const iconName = icon.toLowerCase()
      return keywords.some(keyword => iconName.includes(keyword))
    })
  }
  
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    icons = icons.filter(icon => icon.toLowerCase().includes(search))
  }
  
  return icons.slice(0, 500)
})

const displayValue = computed(() => props.modelValue || '')

const handleSelect = (icon: string) => {
  emit('update:modelValue', icon)
  popoverVisible.value = false
}

const handleClear = () => {
  emit('update:modelValue', '')
}

watch(popoverVisible, (val) => {
  if (!val) {
    searchText.value = ''
    activeCategory.value = ''
  }
})
</script>

<style scoped>
.icon-picker {
  max-height: 400px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.icon-search {
  margin-bottom: 10px;
}
.icon-tabs {
  margin-bottom: 10px;
}
.icon-tabs :deep(.el-radio-group) {
  flex-wrap: wrap;
}
.icon-category {
  margin-bottom: 10px;
}
.icon-list {
  display: grid;
  grid-template-columns: repeat(10, 1fr);
  gap: 6px;
  max-height: 280px;
  overflow-y: auto;
  padding: 4px;
}
.icon-item {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 16px;
}
.icon-item:hover {
  background-color: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
}
.icon-item.active {
  background-color: var(--el-color-primary);
  color: white;
}
.icon-footer {
  text-align: center;
  padding: 8px 0;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
