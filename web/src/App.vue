<template>
  <div id="app" class="min-h-screen bg-bg-base">
    <!-- 顶部导航栏 -->
    <header class="bg-bg-surface shadow-md-md">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-6">
          <div class="flex items-center">
            <h1 class="text-3xl font-bold text-text-primary mr-8">TimeLog</h1>
            <!-- 导航菜单 -->
            <nav class="hidden md:flex space-x-8">
              <router-link
                to="/"
                class="text-text-secondary hover:text-text-primary px-3 py-2 text-sm font-medium transition-colors"
                :class="{
                  'text-brand font-semibold': $route.name === 'Home',
                }"
              >
                Dashboard
              </router-link>
              <router-link
                to="/timelogs"
                class="text-text-secondary hover:text-text-primary px-3 py-2 text-sm font-medium transition-colors"
                :class="{
                  'text-brand font-semibold': $route.name === 'TimeLog',
                }"
              >
                Time Logs
              </router-link>
              <router-link
                to="/tasks"
                class="text-text-secondary hover:text-text-primary px-3 py-2 text-sm font-medium transition-colors"
                :class="{
                  'text-brand font-semibold': $route.name === 'Tasks',
                }"
              >
                Tasks
              </router-link>
              <router-link
                to="/categories"
                class="text-text-secondary hover:text-text-primary px-3 py-2 text-sm font-medium transition-colors"
                :class="{
                  'text-brand font-semibold': $route.name === 'Categories',
                }"
              >
                Categories
              </router-link>
              <router-link
                to="/statistics"
                class="text-text-secondary hover:text-text-primary px-3 py-2 text-sm font-medium transition-colors"
                :class="{
                  'text-brand font-semibold': $route.name === 'Statistics',
                }"
              >
                Statistics
              </router-link>
              <router-link
                to="/constraints"
                class="text-text-secondary hover:text-text-primary px-3 py-2 text-sm font-medium transition-colors"
                :class="{
                  'text-brand font-semibold': $route.name === 'Constraints',
                }"
              >
                约束
              </router-link>
              <router-link
                to="/metrics"
                class="text-text-secondary hover:text-text-primary px-3 py-2 text-sm font-medium transition-colors"
                :class="{
                  'text-brand font-semibold': $route.name === 'Metrics',
                }"
              >
                指标
              </router-link>
            </nav>
          </div>

          <div class="hidden md:flex items-center gap-4">
            <button
              @click="themeCycle"
              class="inline-flex items-center gap-2 rounded-full border border-default px-3 py-2 text-sm font-medium text-text-secondary transition hover:border-default hover:text-text-primary"
              :title="themeLabel"
            >
              <component :is="themeIcon" class="h-4 w-4" />
              <span class="text-xs">{{ themeLabel }}</span>
            </button>
            <router-link
              to="/passkey/register"
              class="inline-flex items-center gap-2 rounded-full border border-default px-4 py-2 text-sm font-medium text-text-secondary transition hover:border-default hover:text-text-primary"
            >
              绑定设备
            </router-link>
            <button
              class="inline-flex items-center gap-2 rounded-full bg-text-primary px-4 py-2 text-sm font-semibold text-bg-surface transition hover:bg-text-muted"
              @click="handleLogout"
            >
              退出
            </button>
          </div>

          <!-- 移动端菜单按钮 -->
          <button
            @click="mobileMenuOpen = !mobileMenuOpen"
            class="md:hidden inline-flex items-center justify-center p-2 rounded-md text-text-tertiary hover:text-text-secondary hover:bg-bg-elevated"
          >
            <Bars3Icon v-if="!mobileMenuOpen" class="h-6 w-6" />
            <XMarkIcon v-else class="h-6 w-6" />
          </button>
        </div>

        <!-- 移动端导航菜单 -->
        <div v-if="mobileMenuOpen" class="md:hidden border-t border-default py-4">
          <nav class="space-y-1">
            <button
              class="block w-full px-3 py-2 text-left text-base font-medium text-text-secondary hover:text-text-primary hover:bg-bg-elevated transition-colors"
              @click="handleThemeToggle"
            >
              <span class="flex items-center gap-2">
                <component :is="themeIcon" class="h-5 w-5" />
                主题: {{ themeLabel }}
              </span>
            </button>
            <router-link
              to="/passkey/register"
              class="block px-3 py-2 text-base font-medium text-text-secondary hover:text-text-primary hover:bg-bg-elevated transition-colors"
              @click="mobileMenuOpen = false"
            >
              绑定设备
            </router-link>
            <router-link
              to="/"
              class="block px-3 py-2 text-base font-medium text-text-secondary hover:text-text-primary hover:bg-bg-elevated transition-colors"
              :class="{
                'text-brand bg-brand-bg': $route.name === 'Home',
              }"
              @click="mobileMenuOpen = false"
            >
              Dashboard
            </router-link>
            <router-link
              to="/timelogs"
              class="block px-3 py-2 text-base font-medium text-text-secondary hover:text-text-primary hover:bg-bg-elevated transition-colors"
              :class="{
                'text-brand bg-brand-bg': $route.name === 'TimeLog',
              }"
              @click="mobileMenuOpen = false"
            >
              Time Logs
            </router-link>
            <router-link
              to="/tasks"
              class="block px-3 py-2 text-base font-medium text-text-secondary hover:text-text-primary hover:bg-bg-elevated transition-colors"
              :class="{
                'text-brand bg-brand-bg': $route.name === 'Tasks',
              }"
              @click="mobileMenuOpen = false"
            >
              Tasks
            </router-link>
            <router-link
              to="/categories"
              class="block px-3 py-2 text-base font-medium text-text-secondary hover:text-text-primary hover:bg-bg-elevated transition-colors"
              :class="{
                'text-brand bg-brand-bg': $route.name === 'Categories',
              }"
              @click="mobileMenuOpen = false"
            >
              Categories
            </router-link>
            <router-link
              to="/statistics"
              class="block px-3 py-2 text-base font-medium text-text-secondary hover:text-text-primary hover:bg-bg-elevated transition-colors"
              :class="{
                'text-brand bg-brand-bg': $route.name === 'Statistics',
              }"
              @click="mobileMenuOpen = false"
            >
              Statistics
            </router-link>
            <router-link
              to="/constraints"
              class="block px-3 py-2 text-base font-medium text-text-secondary hover:text-text-primary hover:bg-bg-elevated transition-colors"
              :class="{
                'text-brand bg-brand-bg': $route.name === 'Constraints',
              }"
              @click="mobileMenuOpen = false"
            >
              约束
            </router-link>
            <router-link
              to="/metrics"
              class="block px-3 py-2 text-base font-medium text-text-secondary hover:text-text-primary hover:bg-bg-elevated transition-colors"
              :class="{
                'text-brand bg-brand-bg': $route.name === 'Metrics',
              }"
              @click="mobileMenuOpen = false"
            >
              指标
            </router-link>
            <button
              class="block w-full px-3 py-2 text-left text-base font-medium text-text-secondary hover:text-text-primary hover:bg-bg-elevated transition-colors"
              @click="handleLogout"
            >
              退出
            </button>
          </nav>
        </div>
      </div>
    </header>

    <!-- 主要内容区域 -->
    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <router-view />
    </main>

    <!-- 全局通知组件 -->
    <div
      v-if="notification.show"
      class="fixed bottom-4 right-4 bg-bg-surface border border-default rounded-lg shadow-md-md p-4 max-w-sm z-50"
      :class="{
        'border-success-border bg-success-bg': notification.type === 'success',
        'border-danger-border bg-danger-bg': notification.type === 'error',
      }"
    >
      <div class="flex items-center">
        <CheckCircleIcon v-if="notification.type === 'success'" class="h-5 w-5 text-success mr-2" />
        <XCircleIcon v-if="notification.type === 'error'" class="h-5 w-5 text-danger mr-2" />
        <p
          class="text-sm font-medium"
          :class="{
            'text-green-800': notification.type === 'success',
            'text-red-800': notification.type === 'error',
          }"
        >
          {{ notification.message }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, provide, onMounted, onUnmounted, computed } from 'vue'
  import {
    CheckCircleIcon,
    XCircleIcon,
    Bars3Icon,
    XMarkIcon,
    MoonIcon,
    SunIcon,
    ComputerDesktopIcon,
  } from '@heroicons/vue/24/outline'
  import { useSettings } from '@/composables/useSettings'
  import { setNotificationHandler } from '@/api'
  import { clearAuthToken } from '@/utils/auth'
  import { useRouter } from 'vue-router'

  // 移动端菜单状态
  const mobileMenuOpen = ref(false)

  // 全局通知系统
  const notification = reactive({
    show: false,
    type: 'success' as 'success' | 'error',
    message: '',
  })

  const showNotification = (type: 'success' | 'error', message: string) => {
    notification.type = type
    notification.message = message
    notification.show = true

    setTimeout(() => {
      notification.show = false
    }, 3000)
  }

  // 通过provide向子组件提供全局通知功能
  provide('showNotification', showNotification)

  const REMINDER_INTERVAL_MS = 25 * 60 * 1000
  let reminderTimer: number | undefined

  const sendReminderNotification = () => {
    try {
      new Notification('该记录TimeLog了', {
        body: '又过去了25分钟，别忘了记录你的时间日志。',
        tag: 'timelog-reminder',
      })
    } catch (error) {
      console.error('Failed to send reminder notification', error)
    }
  }

  const startReminderTimer = () => {
    if (reminderTimer) {
      window.clearInterval(reminderTimer)
    }
    reminderTimer = window.setInterval(sendReminderNotification, REMINDER_INTERVAL_MS)
  }

  const initSystemNotifications = () => {
    if (typeof window === 'undefined' || !('Notification' in window)) {
      showNotification('error', '当前浏览器不支持系统通知提醒')
      return
    }

    if (Notification.permission === 'granted') {
      startReminderTimer()
      return
    }

    if (Notification.permission === 'denied') {
      showNotification('error', '请在浏览器设置中允许通知以启用提醒')
      return
    }

    Notification.requestPermission().then(permission => {
      if (permission === 'granted') {
        startReminderTimer()
        showNotification('success', '已启用每25分钟一次的记录提醒')
      } else {
        showNotification('error', '未授予通知权限，无法启用提醒')
      }
    })
  }

  const router = useRouter()

  const handleLogout = () => {
    clearAuthToken()
    router.push('/login')
  }

  const handleThemeToggle = () => {
    themeCycle()
    mobileMenuOpen.value = false
  }

  // Theme toggle
  const { theme, updateSetting } = useSettings()

  const themeCycle = () => {
    const cycle = ['light', 'dark', 'auto'] as const
    const current = theme.value
    const next = cycle[(cycle.indexOf(current) + 1) % cycle.length]
    updateSetting('theme', next)

    const html = document.documentElement
    const systemDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    const isDark = next === 'dark' || (next === 'auto' && systemDark)
    if (isDark) {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }
  }

  const themeIcon = computed(() => {
    if (theme.value === 'dark') return MoonIcon
    if (theme.value === 'light') return SunIcon
    return ComputerDesktopIcon
  })

  const themeLabel = computed(() => {
    if (theme.value === 'dark') return '暗色'
    if (theme.value === 'light') return '亮色'
    return '自动'
  })

  // Initialize settings on app mount
  onMounted(() => {
    const { loadSettings, theme } = useSettings()
    loadSettings()

    // Theme switching logic
    const applyTheme = (t: string) => {
      const html = document.documentElement
      const systemDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      const isDark = t === 'dark' || (t === 'auto' && systemDark)
      if (isDark) {
        html.classList.add('dark')
      } else {
        html.classList.remove('dark')
      }
    }

    applyTheme(theme.value)

    // Watch for system theme changes when auto
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
      if (theme.value === 'auto') {
        applyTheme('auto')
      }
    })

    initSystemNotifications()

    // Register notification handler for API timeout errors
    setNotificationHandler(showNotification)
  })

  onUnmounted(() => {
    if (reminderTimer) {
      window.clearInterval(reminderTimer)
    }
  })
</script>
