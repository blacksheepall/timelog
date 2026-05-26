<template>
  <div class="space-y-6">
    <!-- 今日概览 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="bg-white shadow rounded-lg p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <ClockIcon class="h-8 w-8 text-blue-600" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Today's Logs</p>
            <p class="text-2xl font-semibold text-gray-900">{{ todayStats.count }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white shadow rounded-lg p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <PlayIcon class="h-8 w-8 text-green-600" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Active Sessions</p>
            <p class="text-2xl font-semibold text-gray-900">{{ todayStats.activeSessions }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white shadow rounded-lg p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <StopIcon class="h-8 w-8 text-red-600" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Total Time</p>
            <p class="text-2xl font-semibold text-gray-900">{{ todayStats.totalTime }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white shadow rounded-lg p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <TagIcon class="h-8 w-8 text-purple-600" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Categories Used</p>
            <p class="text-2xl font-semibold text-gray-900">{{ todayStats.categoriesUsed }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 分类时长统计 -->
    <div class="bg-white shadow rounded-lg">
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-medium text-gray-900">Category Duration Stats</h3>
      </div>
      <div v-if="loading" class="p-6 text-center">
        <div
          class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"
        ></div>
        <p class="mt-2 text-gray-600">Loading...</p>
      </div>
      <div v-else-if="categoryStats.length === 0" class="p-6 text-center text-gray-500">
        No category statistics available.
      </div>
      <div v-else class="divide-y divide-gray-200">
        <div v-for="stat in categoryStats" :key="stat.category.id" class="p-6 hover:bg-gray-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <span
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-white"
                :style="{ backgroundColor: stat.category.color }"
              >
                {{ stat.category.name }}
              </span>
              <div>
                <p class="text-sm font-medium text-gray-900">{{ stat.duration }}</p>
                <div class="mt-1 w-48 bg-gray-200 rounded-full h-2">
                  <div
                    class="bg-blue-600 h-2 rounded-full"
                    :style="{ width: stat.percentage + '%' }"
                  ></div>
                </div>
              </div>
            </div>
            <span class="text-sm font-semibold text-gray-700"> {{ stat.percentage }}% </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 最近的时间记录 -->
    <div class="bg-white shadow rounded-lg">
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-medium text-gray-900">Recent Time Logs</h3>
      </div>
      <div v-if="loading" class="p-6 text-center">
        <div
          class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"
        ></div>
        <p class="mt-2 text-gray-600">Loading...</p>
      </div>
      <div v-else-if="recentLogs.length === 0" class="p-6 text-center text-gray-500">
        No recent time logs found.
      </div>
      <div v-else class="divide-y divide-gray-200">
        <div v-for="log in recentLogs" :key="log.id" class="p-6 hover:bg-gray-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <span
                v-if="getCategory(log.category_id)"
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-white"
                :style="{ backgroundColor: getCategory(log.category_id)!.color }"
              >
                {{ getCategory(log.category_id)!.name }}
              </span>
              <span v-else class="text-gray-400">—</span>
              <div>
                <p class="text-sm font-medium text-gray-900">{{ log.remark || 'No remarks' }}</p>
                <p class="text-xs text-gray-500">
                  {{ formatDateTime(log.start_time) }} -
                  {{ log.end_time ? formatDateTime(log.end_time) : 'Ongoing' }}
                </p>
              </div>
            </div>
            <span class="text-sm text-gray-500">
              {{ calculateDuration(log.start_time, log.end_time) }}
            </span>
          </div>
        </div>
      </div>
      <div class="px-6 py-3 bg-gray-50 text-center">
        <router-link to="/timelogs" class="text-sm font-medium text-blue-600 hover:text-blue-500">
          View all time logs →
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue'
  import { ClockIcon, PlayIcon, StopIcon, TagIcon } from '@heroicons/vue/24/outline'
  import { timelogAPI, categoryAPI } from '@/api'
  import { formatDateTime, calculateDuration } from '@/utils/date'
  import { parseISO } from 'date-fns'
  import type { TimeLog, Category } from '@/types'

  const loading = ref(false)
  const recentLogs = ref<TimeLog[]>([])
  const todayLogs = ref<TimeLog[]>([])
  const categories = ref<Category[]>([])

  const getCategory = (categoryId: number): Category | undefined => {
    return categories.value.find(c => c.id === categoryId)
  }

  // 今日统计数据
  const todayStats = computed(() => {
    const activeSessions = todayLogs.value.filter(log => !log.end_time).length
    const categoriesUsed = new Set(todayLogs.value.map(log => log.category_id)).size

    // 计算总时间（包括ongoing记录）
    const totalMinutes = todayLogs.value.reduce((total, log) => {
      const start = parseISO(log.start_time)
      const end = log.end_time ? parseISO(log.end_time) : new Date()
      return total + (end.getTime() - start.getTime()) / (1000 * 60)
    }, 0)

    const totalTime =
      totalMinutes > 60
        ? `${Math.floor(totalMinutes / 60)}h ${Math.round(totalMinutes % 60)}m`
        : `${Math.round(totalMinutes)}m`

    return {
      count: todayLogs.value.length,
      activeSessions,
      totalTime,
      categoriesUsed,
    }
  })

  // 分类时长统计
  const categoryStats = computed(() => {
    const stats: Record<number, { category: Category; minutes: number }> = {}

    // 获取今天的开始和结束时间（UTC）
    const today = new Date()
    const todayLocalStr = today.toISOString().split('T')[0] // YYYY-MM-DD 本地日期
    const todayStart = parseISO(todayLocalStr + 'T00:00:00Z') // UTC开始时间
    const todayEnd = parseISO(todayLocalStr + 'T23:59:59.999Z') // UTC结束时间

    // 计算每个分类的总时长
    todayLogs.value.forEach(log => {
      const category = getCategory(log.category_id)
      if (!category) return

      const start = parseISO(log.start_time)
      const end = log.end_time ? parseISO(log.end_time) : todayEnd

      // 计算实际的开始时间（如果记录开始时间早于今天开始，使用今天开始）
      const actualStart = start < todayStart ? todayStart : start
      // 计算实际的结束时间（如果记录结束时间晚于今天结束，使用今天结束）
      const actualEnd = end > todayEnd ? todayEnd : end

      // 确保开始时间不晚于结束时间
      if (actualStart >= actualEnd) return

      // 计算时长（分钟）
      const minutes = (actualEnd.getTime() - actualStart.getTime()) / (1000 * 60)

      if (!stats[log.category_id]) {
        stats[log.category_id] = {
          category,
          minutes: 0,
        }
      }
      stats[log.category_id].minutes += minutes
    })

    // 计算总时长
    const totalMinutes = Object.values(stats).reduce((sum, stat) => sum + stat.minutes, 0)

    // 转换为显示格式并排序（时长长的在前）
    return Object.values(stats)
      .map(stat => {
        const percentage = totalMinutes > 0 ? Math.round((stat.minutes / totalMinutes) * 100) : 0
        const duration =
          stat.minutes > 60
            ? `${Math.floor(stat.minutes / 60)}h ${Math.round(stat.minutes % 60)}m`
            : `${Math.round(stat.minutes)}m`

        return {
          category: stat.category,
          minutes: stat.minutes,
          duration,
          percentage,
        }
      })
      .filter(stat => stat.minutes > 0) // 只显示时长大于0的分类
      .sort((a, b) => b.minutes - a.minutes) // 按时长降序排列
  })

  const loadRecentLogs = async () => {
    loading.value = true
    try {
      const response = await timelogAPI.getRecent(5)
      recentLogs.value = response.data || []

      // 加载今天的所有记录用于统计
      await loadTodayLogs()
    } catch (err) {
      console.error('Error loading recent logs:', err)
    } finally {
      loading.value = false
    }
  }

  const loadTodayLogs = async () => {
    try {
      // 获取今天的开始时间和结束时间（UTC）
      const today = new Date()
      const todayLocalStr = today.toISOString().split('T')[0] // YYYY-MM-DD 本地日期

      // 创建今天的开始时间（本地时间00:00:00）
      const todayStart = new Date(todayLocalStr + 'T00:00:00')
      const todayEnd = new Date(todayLocalStr + 'T23:59:59.999')

      // 转换为ISO字符串用于比较
      const todayStartIso = todayStart.toISOString()
      const todayEndIso = todayEnd.toISOString()

      // 获取所有时间记录，筛选在今天的记录（考虑跨天情况）
      const response = await timelogAPI.getAll()
      if (response.data) {
        todayLogs.value = response.data.filter(log => {
          // 只要记录的开始时间或结束时间在今天的范围内，就包含它
          const logStart = log.start_time
          const logEnd = log.end_time || new Date().toISOString()

          // 记录在今天开始，或者今天结束，或者在今天的范围内
          return (
            (logStart >= todayStartIso && logStart <= todayEndIso) || // 今天开始
            (logEnd >= todayStartIso && logEnd <= todayEndIso) || // 今天结束
            (logStart < todayStartIso && logEnd > todayEndIso) // 跨越今天
          )
        })
      }
    } catch (err) {
      console.error('Error loading today logs:', err)
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
    loadRecentLogs()
    loadCategories()
  })
</script>
