<template>
  <div class="space-y-6">
    <h1 class="text-2xl font-bold text-text-primary">Statistics</h1>

    <!-- 时间过滤器 -->
    <div class="bg-bg-surface shadow-md-md rounded-lg p-6">
      <h2 class="text-lg font-medium text-text-primary mb-4">Time Range</h2>

      <!-- 快捷选项 -->
      <div class="flex flex-wrap gap-2 mb-4">
        <button
          @click="setToday"
          class="px-3 py-1 text-sm font-medium text-blue-700 bg-blue-100 rounded-md hover:bg-blue-200 focus:outline-none focus:ring-2 focus:ring-brand"
        >
          Today
        </button>
        <button
          @click="setThisWeek"
          class="px-3 py-1 text-sm font-medium text-blue-700 bg-blue-100 rounded-md hover:bg-blue-200 focus:outline-none focus:ring-2 focus:ring-brand"
        >
          This Week
        </button>
        <button
          @click="setLast30Days"
          class="px-3 py-1 text-sm font-medium text-blue-700 bg-blue-100 rounded-md hover:bg-blue-200 focus:outline-none focus:ring-2 focus:ring-brand"
        >
          Last 30 Days
        </button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label for="start_date" class="block text-sm font-medium text-text-muted mb-2">
            Start Date
          </label>
          <input
            id="start_date"
            v-model="dateRange.start"
            type="date"
            class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
          />
        </div>
        <div>
          <label for="end_date" class="block text-sm font-medium text-text-muted mb-2">
            End Date
          </label>
          <input
            id="end_date"
            v-model="dateRange.end"
            type="date"
            class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
          />
        </div>
        <div class="flex items-end">
          <button
            @click="applyFilter"
            class="w-full px-4 py-2 text-sm font-medium text-bg-surface bg-brand border border-transparent rounded-md shadow-sm hover:bg-brand-hover focus:outline-none focus:ring-2 focus:ring-brand"
          >
            Apply Filter
          </button>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="bg-bg-surface shadow-md-md rounded-lg p-6 text-center">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-brand"></div>
      <p class="mt-2 text-text-secondary">Loading statistics...</p>
    </div>

    <!-- 统计概览 -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="bg-bg-surface shadow-md-md rounded-lg p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <ClockIcon class="h-8 w-8 text-brand" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-text-secondary">Total Logs</p>
            <p class="text-2xl font-semibold text-text-primary">{{ statistics.totalLogs }}</p>
          </div>
        </div>
      </div>

      <div class="bg-bg-surface shadow-md-md rounded-lg p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <ClockIcon class="h-8 w-8 text-success" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-text-secondary">Total Time</p>
            <p class="text-2xl font-semibold text-text-primary">{{ statistics.totalTime }}</p>
          </div>
        </div>
      </div>

      <div class="bg-bg-surface shadow-md-md rounded-lg p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <ClockIcon class="h-8 w-8 text-yellow-600" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-text-secondary">Average Session</p>
            <p class="text-2xl font-semibold text-text-primary">{{ statistics.averageSession }}</p>
          </div>
        </div>
      </div>

      <div class="bg-bg-surface shadow-md-md rounded-lg p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <PlayIcon class="h-8 w-8 text-danger" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-text-secondary">Active Sessions</p>
            <p class="text-2xl font-semibold text-text-primary">{{ statistics.activeSessions }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 按标签统计 -->
    <div v-if="!loading" class="bg-bg-surface shadow-md-md rounded-lg">
      <div class="px-6 py-4 border-b border-default">
        <h2 class="text-lg font-medium text-text-primary">Time by Tags</h2>
      </div>
      <div v-if="tagStats.length === 0" class="p-6 text-center text-text-secondary">
        No data available for the selected time range.
      </div>
      <div v-else class="divide-y divide-default">
        <div v-for="stat in tagStats" :key="stat.category.id" class="p-6">
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center space-x-3">
              <span
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-bg-surface"
                :style="{ backgroundColor: stat.category.color }"
              >
                {{ stat.category.name }}
              </span>
              <span class="text-sm text-text-secondary">{{ stat.count }} sessions</span>
            </div>
            <span class="text-sm font-medium text-text-primary">{{ stat.totalTime }}</span>
          </div>
          <div class="w-full bg-bg-elevated rounded-full h-2">
            <div
              class="h-2 rounded-full"
              :style="{
                backgroundColor: stat.category.color,
                width: `${stat.percentage}%`,
              }"
            ></div>
          </div>
          <div class="mt-1 text-xs text-text-secondary">
            {{ stat.percentage.toFixed(1) }}% of total time
          </div>
        </div>
      </div>
    </div>

    <!-- 每日时间分布 -->
    <div v-if="!loading" class="bg-bg-surface shadow-md-md rounded-lg">
      <div class="px-6 py-4 border-b border-default">
        <h2 class="text-lg font-medium text-text-primary">Daily Time Distribution</h2>
      </div>
      <div v-if="dailyStats.length === 0" class="p-6 text-center text-text-secondary">
        No data available for the selected time range.
      </div>
      <div v-else class="p-6">
        <div class="space-y-4">
          <div v-for="day in dailyStats" :key="day.date" class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <span class="text-sm font-medium text-text-primary w-24">{{ day.date }}</span>
              <div class="w-48 bg-bg-elevated rounded-full h-2">
                <div
                  class="h-2 bg-brand rounded-full"
                  :style="{ width: `${(day.minutes / maxDailyMinutes) * 100}%` }"
                ></div>
              </div>
            </div>
            <span class="text-sm text-text-secondary">{{ day.timeFormatted }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, inject } from 'vue'
  import { ClockIcon, PlayIcon } from '@heroicons/vue/24/outline'
  import { timelogAPI, categoryAPI } from '@/api'
  import { formatDate } from '@/utils/date'
  import type { TimeLog, Category } from '@/types'

  // 注入全局通知功能
  const showNotification = inject('showNotification') as (
    type: 'success' | 'error',
    message: string
  ) => void

  const loading = ref(false)
  const timeLogs = ref<TimeLog[]>([])
  const categories = ref<Category[]>([])

  const getCategory = (categoryId: number): Category | undefined => {
    return categories.value.find(c => c.id === categoryId)
  }

  // 日期范围
  const dateRange = ref({
    start: formatDate(new Date(Date.now() - 30 * 24 * 60 * 60 * 1000)), // 30天前
    end: formatDate(new Date()),
  })

  // 过滤后的时间记录
  const filteredLogs = computed(() => {
    if (!dateRange.value.start || !dateRange.value.end) return timeLogs.value

    const start = new Date(dateRange.value.start)
    const end = new Date(dateRange.value.end)
    end.setHours(23, 59, 59, 999) // 包含结束日期的整天

    return timeLogs.value.filter(log => {
      const logDate = new Date(log.start_time)
      return logDate >= start && logDate <= end
    })
  })

  // 统计概览
  const statistics = computed(() => {
    const logs = filteredLogs.value
    const completedLogs = logs.filter(log => log.end_time)

    // 计算总时间（分钟）
    const totalMinutes = completedLogs.reduce((total, log) => {
      const start = new Date(log.start_time)
      const end = new Date(log.end_time!)
      return total + (end.getTime() - start.getTime()) / (1000 * 60)
    }, 0)

    // 格式化总时间
    const totalTime =
      totalMinutes > 60
        ? `${Math.floor(totalMinutes / 60)}h ${Math.round(totalMinutes % 60)}m`
        : `${Math.round(totalMinutes)}m`

    // 平均会话时间
    const averageMinutes = completedLogs.length > 0 ? totalMinutes / completedLogs.length : 0
    const averageSession =
      averageMinutes > 60
        ? `${Math.floor(averageMinutes / 60)}h ${Math.round(averageMinutes % 60)}m`
        : `${Math.round(averageMinutes)}m`

    // 活跃会话数
    const activeSessions = logs.filter(log => !log.end_time).length

    return {
      totalLogs: logs.length,
      totalTime,
      averageSession,
      activeSessions,
    }
  })

  // 按标签统计
  const tagStats = computed(() => {
    const logs = filteredLogs.value.filter(log => log.end_time)
    const tagMap = new Map<number, { category: any; minutes: number; count: number }>()

    let totalMinutes = 0

    logs.forEach(log => {
      const category = getCategory(log.category_id)
      if (!category) return

      const start = new Date(log.start_time)
      const end = new Date(log.end_time!)
      const minutes = (end.getTime() - start.getTime()) / (1000 * 60)

      totalMinutes += minutes

      if (tagMap.has(log.category_id)) {
        const existing = tagMap.get(log.category_id)!
        existing.minutes += minutes
        existing.count += 1
      } else {
        tagMap.set(log.category_id, {
          category,
          minutes,
          count: 1,
        })
      }
    })

    return Array.from(tagMap.values())
      .map(stat => ({
        ...stat,
        totalTime:
          stat.minutes > 60
            ? `${Math.floor(stat.minutes / 60)}h ${Math.round(stat.minutes % 60)}m`
            : `${Math.round(stat.minutes)}m`,
        percentage: totalMinutes > 0 ? (stat.minutes / totalMinutes) * 100 : 0,
      }))
      .sort((a, b) => b.minutes - a.minutes)
  })

  // 每日统计
  const dailyStats = computed(() => {
    const logs = filteredLogs.value.filter(log => log.end_time)
    const dailyMap = new Map<string, number>()

    logs.forEach(log => {
      const date = new Date(log.start_time).toISOString().split('T')[0]
      const start = new Date(log.start_time)
      const end = new Date(log.end_time!)
      const minutes = (end.getTime() - start.getTime()) / (1000 * 60)

      if (dailyMap.has(date)) {
        dailyMap.set(date, dailyMap.get(date)! + minutes)
      } else {
        dailyMap.set(date, minutes)
      }
    })

    return Array.from(dailyMap.entries())
      .map(([date, minutes]) => ({
        date: new Date(date).toLocaleDateString(),
        minutes,
        timeFormatted:
          minutes > 60
            ? `${Math.floor(minutes / 60)}h ${Math.round(minutes % 60)}m`
            : `${Math.round(minutes)}m`,
      }))
      .sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime())
  })

  // 每日最大分钟数（用于进度条）
  const maxDailyMinutes = computed(() => {
    return Math.max(...dailyStats.value.map(day => day.minutes), 1)
  })

  const loadTimeLogs = async () => {
    loading.value = true

    try {
      const response = await timelogAPI.getAll()
      timeLogs.value = response.data || []
    } catch (err) {
      console.error('Error loading time logs:', err)
      showNotification('error', 'Failed to load time logs')
    } finally {
      loading.value = false
    }
  }

  const applyFilter = () => {
    // 过滤逻辑已经在computed中处理，这里只需要触发重新计算
    console.log('Filter applied:', dateRange.value)
  }

  // 设置今天的日期范围
  const setToday = () => {
    const today = new Date()
    dateRange.value = {
      start: formatDate(today),
      end: formatDate(today),
    }
  }

  // 设置本周的日期范围（周一到周日）
  const setThisWeek = () => {
    const today = new Date()
    const dayOfWeek = today.getDay() // 0 = Sunday, 1 = Monday, ..., 6 = Saturday

    // 计算本周一的日期
    const monday = new Date(today)
    monday.setDate(today.getDate() - (dayOfWeek === 0 ? 6 : dayOfWeek - 1))
    monday.setHours(0, 0, 0, 0)

    // 计算本周日的日期
    const sunday = new Date(monday)
    sunday.setDate(monday.getDate() + 6)
    sunday.setHours(23, 59, 59, 999)

    dateRange.value = {
      start: formatDate(monday),
      end: formatDate(sunday),
    }
  }

  // 设置最近30天的日期范围
  const setLast30Days = () => {
    const today = new Date()
    const thirtyDaysAgo = new Date(today)
    thirtyDaysAgo.setDate(today.getDate() - 30)

    dateRange.value = {
      start: formatDate(thirtyDaysAgo),
      end: formatDate(today),
    }
  }

  const loadCategories = async () => {
    try {
      const response = await categoryAPI.getAll()
      categories.value = response.data || []
    } catch (err) {
      console.error('Error loading categories:', err)
    }
  }

  onMounted(() => {
    loadTimeLogs()
    loadCategories()
  })
</script>
