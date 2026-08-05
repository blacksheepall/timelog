<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-text-primary">学习情况</h1>
      <button
        @click="syncMaimemo"
        :disabled="syncing"
        class="inline-flex items-center px-4 py-2 text-sm font-medium text-bg-surface bg-brand border border-transparent rounded-md shadow-sm hover:bg-brand-hover focus:outline-none focus:ring-2 focus:ring-brand disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <ArrowPathIcon class="h-5 w-5 mr-2" :class="{ 'animate-spin': syncing }" />
        {{ syncing ? '同步中...' : '同步墨墨数据' }}
      </button>
    </div>

    <!-- 今日学习情况 -->
    <div class="bg-bg-surface shadow-md-md rounded-lg">
      <div class="px-6 py-4 border-b border-default">
        <h2 class="text-lg font-medium text-text-primary">今日学习情况</h2>
        <p class="text-sm text-text-muted mt-1">
          数据来自墨墨背单词，当日未在 App 内学习时无法获取。
        </p>
      </div>

      <div v-if="loading" class="p-8 text-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-brand mx-auto"></div>
        <p class="mt-2 text-text-secondary">加载中...</p>
      </div>

      <div v-else-if="studyMetrics.length === 0" class="p-8 text-center text-text-secondary">
        <BookOpenIcon class="h-12 w-12 mx-auto mb-4 text-text-tertiary" />
        <p>暂无学习数据。点击右上角「同步墨墨数据」获取。</p>
      </div>

      <div v-else class="p-6 space-y-6">
        <!-- 今日完成进度条 -->
        <div v-if="studyProgress">
          <div class="flex justify-between items-baseline text-sm mb-2">
            <span class="text-text-secondary">今日完成进度</span>
            <span class="text-text-primary font-medium">
              {{ studyProgress.finished }} / {{ studyProgress.total }}（{{ studyProgress.rate }}%）
            </span>
          </div>
          <div class="w-full bg-bg-elevated rounded-full h-2.5">
            <div
              class="bg-brand h-2.5 rounded-full transition-all"
              :style="{ width: studyProgress.percent + '%' }"
            ></div>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div
            v-for="metric in studyMetrics"
            :key="metric.id"
            class="border border-default rounded-lg p-4"
          >
            <p class="text-sm font-medium text-text-secondary">{{ metric.name }}</p>
            <p class="mt-2">
              <span class="text-3xl font-bold text-brand">
                {{ formatValue(metric.current_value) }}
              </span>
              <span class="text-text-muted ml-1">{{ metric.unit }}</span>
            </p>
            <p v-if="metric.last_recorded_at" class="text-xs text-text-muted mt-2">
              更新于 {{ formatDateTime(metric.last_recorded_at) }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- 指标完成情况 -->
    <div class="bg-bg-surface shadow-md-md rounded-lg">
      <div class="px-6 py-4 border-b border-default">
        <h2 class="text-lg font-medium text-text-primary">指标完成情况</h2>
        <p class="text-sm text-text-muted mt-1">启用中的约束关联的指标达成状态。</p>
      </div>

      <div v-if="loading" class="p-8 text-center text-text-secondary">加载中...</div>

      <div v-else-if="completions.length === 0" class="p-8 text-center text-text-secondary">
        <ClipboardDocumentCheckIcon class="h-12 w-12 mx-auto mb-4 text-text-tertiary" />
        <p>暂无关联指标的启用约束。</p>
      </div>

      <div v-else class="divide-y divide-default">
        <div
          v-for="item in completions"
          :key="item.constraint.id"
          class="px-6 py-4 flex items-center justify-between"
        >
          <div class="flex-1">
            <p class="text-sm font-medium text-text-primary">{{ item.constraint.description }}</p>
            <p class="text-sm text-text-secondary mt-1">
              <template v-if="item.evaluation">
                当前
                <span class="font-medium text-text-primary">{{ formatValue(item.evaluation.actual) }}</span>
                ，目标 {{ operatorSymbol(item.evaluation.operator) }}
                <span class="font-medium text-text-primary">{{ formatValue(item.evaluation.target) }}</span>
              </template>
              <template v-else>暂无指标数据</template>
            </p>
          </div>
          <span
            class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ml-4"
            :class="statusBadgeClass(item)"
          >
            {{ statusText(item) }}
          </span>
        </div>
      </div>
    </div>

    <!-- 最近学习记录 -->
    <div class="bg-bg-surface shadow-md-md rounded-lg">
      <div class="px-6 py-4 border-b border-default">
        <h2 class="text-lg font-medium text-text-primary">最近学习记录</h2>
      </div>

      <div v-if="loading" class="p-8 text-center text-text-secondary">加载中...</div>

      <div v-else-if="recentRecords.length === 0" class="p-8 text-center text-text-secondary">
        暂无记录
      </div>

      <table v-else class="min-w-full divide-y divide-default">
        <thead>
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase">指标</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase">数值</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase">来源</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase">时间</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-default">
          <tr v-for="(record, index) in recentRecords" :key="index">
            <td class="px-6 py-3 text-sm text-text-secondary">{{ record.metricName }}</td>
            <td class="px-6 py-3 text-sm text-text-primary font-medium">
              {{ formatValue(record.value) }} {{ record.unit }}
            </td>
            <td class="px-6 py-3 text-sm text-text-muted">{{ record.source || '-' }}</td>
            <td class="px-6 py-3 text-sm text-text-secondary">
              {{ formatDateTime(record.recorded_at || record.created_at) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import {
    ArrowPathIcon,
    BookOpenIcon,
    ClipboardDocumentCheckIcon,
  } from '@heroicons/vue/24/outline'
  import { ElMessage } from 'element-plus'
  import { constraintAPI, datasourceAPI, metricAPI } from '@/api'
  import type { Constraint, ConstraintEvaluation } from '@/gen/api/timelog/v1/constraint'
  import type { Metric, MetricRecord } from '@/gen/api/timelog/v1/metric'

  // 与后端 maimemo mapper 产出的指标名保持一致
  const STUDY_METRIC_NAMES = [
    '今日学习单词',
    '今日新学单词',
    '今日已完成单词',
    '今日学习时长',
    '今日应背单词',
    '今日完成率',
  ]

  interface CompletionItem {
    constraint: Constraint
    evaluation: ConstraintEvaluation | null
  }

  interface StudyRecord extends MetricRecord {
    metricName: string
    unit: string
  }

  const loading = ref(false)
  const syncing = ref(false)
  const metrics = ref<Metric[]>([])
  const completions = ref<CompletionItem[]>([])
  const recentRecords = ref<StudyRecord[]>([])

  const studyMetrics = computed(() =>
    metrics.value.filter(m => STUDY_METRIC_NAMES.includes(m.name))
  )

  // 今日完成进度：直接取「已完成 / 应背」两个指标计算，比读「今日完成率」更稳
  const studyProgress = computed(() => {
    const finished = metrics.value.find(m => m.name === '今日已完成单词')?.current_value
    const total = metrics.value.find(m => m.name === '今日应背单词')?.current_value
    if (finished == null || total == null || total <= 0) return null
    const rate = Math.round((finished / total) * 10000) / 100
    return { finished, total, rate, percent: Math.min(rate, 100) }
  })

  const loadData = async () => {
    loading.value = true
    try {
      const metricsResp = await metricAPI.getAll()
      metrics.value = metricsResp.data || []

      await Promise.all([loadCompletions(), loadRecentRecords()])
    } catch (err: any) {
      ElMessage.error(err.response?.data?.message || '加载失败')
    } finally {
      loading.value = false
    }
  }

  const loadCompletions = async () => {
    const resp = await constraintAPI.getAll(true)
    const constraints = (resp.data || []).filter(c => c.metric_id != null)

    const results = await Promise.allSettled(
      constraints.map(c => constraintAPI.evaluate(c.id))
    )
    completions.value = constraints.map((constraint, i) => {
      const result = results[i]
      return {
        constraint,
        evaluation:
          result.status === 'fulfilled' && result.value.data ? result.value.data : null,
      }
    })
  }

  const loadRecentRecords = async () => {
    const study = metrics.value.filter(m => STUDY_METRIC_NAMES.includes(m.name))
    const results = await Promise.allSettled(
      study.map(m => metricAPI.getRecords(m.id))
    )

    const merged: StudyRecord[] = []
    results.forEach((result, i) => {
      if (result.status !== 'fulfilled') return
      const metric = study[i]
      for (const record of result.value.data || []) {
        merged.push({ ...record, metricName: metric.name, unit: metric.unit })
      }
    })

    merged.sort((a, b) => recordTime(b) - recordTime(a))
    recentRecords.value = merged.slice(0, 10)
  }

  const recordTime = (record: MetricRecord) => {
    const raw = record.recorded_at || record.created_at
    return raw ? new Date(raw).getTime() : 0
  }

  const syncMaimemo = async () => {
    syncing.value = true
    try {
      const resp = await datasourceAPI.sync('maimemo')
      const result = resp.data
      if (result && result.failed > 0) {
        ElMessage.warning(`同步完成：成功 ${result.synced} 条，失败 ${result.failed} 条`)
      } else {
        ElMessage.success(`同步完成：成功 ${result?.synced ?? 0} 条`)
      }
      await loadData()
    } catch (err: any) {
      ElMessage.error(err.response?.data?.message || '同步失败')
    } finally {
      syncing.value = false
    }
  }

  const operatorSymbol = (op: string) => {
    const map: Record<string, string> = {
      gt: '>',
      gte: '≥',
      lt: '<',
      lte: '≤',
      eq: '=',
      ne: '≠',
    }
    return map[op] || op
  }

  const statusText = (item: CompletionItem) => {
    if (!item.evaluation) return '无数据'
    return item.evaluation.passed ? '通过' : '未通过'
  }

  const statusBadgeClass = (item: CompletionItem) => {
    if (!item.evaluation) return 'bg-bg-elevated text-text-muted'
    return item.evaluation.passed
      ? 'bg-success-bg text-success'
      : 'bg-danger-bg text-danger'
  }

  const formatValue = (value: number | undefined) => {
    if (value === undefined || value === null) return '-'
    return Number.isInteger(value) ? value.toString() : value.toFixed(2)
  }

  const formatDateTime = (dateString: string | undefined) => {
    if (!dateString) return ''
    return new Date(dateString).toLocaleString('zh-CN')
  }

  onMounted(() => {
    loadData()
  })
</script>
