<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-text-primary">指标管理</h1>
      <button
        @click="toggleForm"
        class="inline-flex items-center px-4 py-2 text-sm font-medium text-bg-surface bg-brand border border-transparent rounded-md shadow-sm hover:bg-brand-hover focus:outline-none focus:ring-2 focus:ring-brand"
      >
        <PlusIcon class="h-5 w-5 mr-2" />
        新建指标
      </button>
    </div>

    <!-- 指标创建/编辑表单 -->
    <div v-if="showForm" class="bg-bg-surface shadow-md-md rounded-lg p-6">
      <h2 class="text-lg font-medium text-text-primary mb-6">
        {{ isEditing ? '编辑指标' : '创建新指标' }}
      </h2>

      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div>
          <label for="name" class="block text-sm font-medium text-text-muted mb-2">
            指标名称 *
          </label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            required
            class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            placeholder="例如：入睡时间"
          />
        </div>

        <div>
          <label for="description" class="block text-sm font-medium text-text-muted mb-2">
            描述
          </label>
          <textarea
            id="description"
            v-model="form.description"
            rows="2"
            class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            placeholder="说明该指标的含义和记录方式"
          ></textarea>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label for="metric_type" class="block text-sm font-medium text-text-muted mb-2">
              指标类型 *
            </label>
            <select
              id="metric_type"
              v-model="form.metric_type"
              required
              class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand bg-bg-surface"
            >
              <option value="numeric">数值</option>
              <option value="count">计数</option>
              <option value="time">时间</option>
              <option value="boolean">是否</option>
            </select>
          </div>

          <div>
            <label for="unit" class="block text-sm font-medium text-text-muted mb-2">
              单位 *
            </label>
            <input
              id="unit"
              v-model="form.unit"
              type="text"
              required
              class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
              placeholder="例如：分钟、小时、次"
            />
          </div>
        </div>

        <div>
          <label for="extra" class="block text-sm font-medium text-text-muted mb-2">
            扩展信息
          </label>
          <input
            id="extra"
            v-model="form.extra"
            type="text"
            class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            placeholder="可选的 JSON 或备注"
          />
        </div>

        <div class="flex justify-end space-x-3">
          <button
            type="button"
            @click="cancelEdit"
            class="px-4 py-2 text-sm font-medium text-text-muted bg-bg-surface border border-default rounded-md shadow-sm hover:bg-bg-elevated focus:outline-none focus:ring-2 focus:ring-brand"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="px-4 py-2 text-sm font-medium text-bg-surface bg-brand border border-transparent rounded-md shadow-sm hover:bg-brand-hover focus:outline-none focus:ring-2 focus:ring-brand disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ loading ? '保存中...' : isEditing ? '更新指标' : '创建指标' }}
          </button>
        </div>
      </form>
    </div>

    <!-- 指标列表 -->
    <div class="bg-bg-surface shadow-md-md rounded-lg">
      <div class="px-6 py-4 border-b border-default">
        <h2 class="text-lg font-medium text-text-primary">指标列表</h2>
        <p class="text-sm text-text-muted mt-1">
          指标数值由 MCP / 自动化脚本写入，这里仅管理元数据并查看当前值。
        </p>
      </div>

      <div v-if="loading" class="p-8 text-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-brand mx-auto"></div>
        <p class="mt-2 text-text-secondary">加载中...</p>
      </div>

      <div v-else-if="error" class="p-8 text-center">
        <div class="text-danger">
          <ExclamationTriangleIcon class="h-8 w-8 mx-auto mb-2" />
          <p>{{ error }}</p>
        </div>
      </div>

      <div v-else-if="metrics.length === 0" class="p-8 text-center text-text-secondary">
        <DocumentTextIcon class="h-12 w-12 mx-auto mb-4 text-text-tertiary" />
        <p>暂无指标。创建你的第一个指标吧！</p>
      </div>

      <div v-else class="divide-y divide-default">
        <div v-for="metric in metrics" :key="metric.id" class="p-6 hover:bg-bg-elevated">
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <h3 class="text-lg font-medium text-text-primary">
                  {{ metric.name }}
                </h3>
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-bg-elevated text-text-primary"
                >
                  {{ formatMetricType(metric.metric_type) }}
                </span>
              </div>

              <p v-if="metric.description" class="text-sm text-text-secondary mb-3">
                {{ metric.description }}
              </p>

              <div class="flex flex-wrap items-center gap-4 text-sm text-text-secondary">
                <span class="inline-flex items-center">
                  <span class="text-2xl font-bold text-brand mr-1">
                    {{ formatValue(metric.current_value) }}
                  </span>
                  <span class="text-text-muted">{{ metric.unit }}</span>
                </span>
                <span v-if="metric.last_recorded_at">
                  最近更新: {{ formatDateTime(metric.last_recorded_at) }}
                </span>
              </div>
            </div>

            <div class="flex items-center space-x-2 ml-4">
              <button
                @click="viewRecords(metric)"
                class="inline-flex items-center px-3 py-1.5 text-sm font-medium text-text-secondary bg-bg-surface border border-default rounded-md hover:bg-bg-elevated focus:outline-none focus:ring-2 focus:ring-brand"
              >
                <ChartBarIcon class="h-4 w-4 mr-1" />
                记录
              </button>
              <button
                @click="editMetric(metric)"
                class="inline-flex items-center px-3 py-1.5 text-sm font-medium text-text-muted bg-bg-surface border border-default rounded-md hover:bg-bg-elevated focus:outline-none focus:ring-2 focus:ring-brand"
              >
                <PencilIcon class="h-4 w-4" />
              </button>
              <button
                @click="deleteMetric(metric)"
                class="inline-flex items-center px-3 py-1.5 text-sm font-medium text-danger bg-bg-surface border border-danger-border rounded-md hover:bg-danger-bg focus:outline-none focus:ring-2 focus:ring-danger"
              >
                <TrashIcon class="h-4 w-4" />
              </button>
            </div>
          </div>

          <!-- 指标记录列表 -->
          <div v-if="selectedMetricId === metric.id" class="mt-4 border-t border-default pt-4">
            <div v-if="recordsLoading" class="p-4 text-center text-text-secondary">
              加载记录中...
            </div>
            <div v-else-if="recordsError" class="p-4 text-center text-danger">
              {{ recordsError }}
            </div>
            <div v-else-if="records.length === 0" class="p-4 text-center text-text-secondary">
              暂无记录
            </div>
            <table v-else class="min-w-full divide-y divide-default">
              <thead>
                <tr>
                  <th class="px-3 py-2 text-left text-xs font-medium text-text-muted uppercase">
                    时间
                  </th>
                  <th class="px-3 py-2 text-left text-xs font-medium text-text-muted uppercase">
                    数值
                  </th>
                  <th class="px-3 py-2 text-left text-xs font-medium text-text-muted uppercase">
                    来源
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-default">
                <tr v-for="record in records" :key="record.id">
                  <td class="px-3 py-2 text-sm text-text-secondary">
                    {{ formatDateTime(record.recorded_at || record.created_at) }}
                  </td>
                  <td class="px-3 py-2 text-sm text-text-primary font-medium">
                    {{ formatValue(record.value) }} {{ metric.unit }}
                  </td>
                  <td class="px-3 py-2 text-sm text-text-muted">
                    {{ record.source || '-' }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, computed, onMounted } from 'vue'
  import {
    PlusIcon,
    PencilIcon,
    TrashIcon,
    ExclamationTriangleIcon,
    DocumentTextIcon,
    ChartBarIcon,
  } from '@heroicons/vue/24/outline'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { metricAPI } from '@/api'
  import type {
    Metric,
    MetricRecord,
    CreateMetricRequest,
    UpdateMetricRequest,
  } from '@/gen/api/timelog/v1/metric'

  const loading = ref(false)
  const error = ref('')
  const showForm = ref(false)
  const editingMetric = ref<Metric | null>(null)
  const metrics = ref<Metric[]>([])

  const selectedMetricId = ref<number | null>(null)
  const records = ref<MetricRecord[]>([])
  const recordsLoading = ref(false)
  const recordsError = ref('')

  const isEditing = computed(() => !!editingMetric.value)

  const form = reactive({
    name: '',
    description: '',
    metric_type: 'numeric',
    unit: '',
    extra: '',
  })

  const resetForm = () => {
    form.name = ''
    form.description = ''
    form.metric_type = 'numeric'
    form.unit = ''
    form.extra = ''
  }

  const loadEditingData = () => {
    if (editingMetric.value) {
      form.name = editingMetric.value.name
      form.description = editingMetric.value.description || ''
      form.metric_type = editingMetric.value.metric_type
      form.unit = editingMetric.value.unit
      form.extra = editingMetric.value.extra || ''
    } else {
      resetForm()
    }
  }

  const loadMetrics = async () => {
    loading.value = true
    error.value = ''

    try {
      const response = await metricAPI.getAll()
      metrics.value = response.data || []
    } catch (err: any) {
      error.value = err.response?.data?.message || '加载指标失败'
      ElMessage.error(error.value)
    } finally {
      loading.value = false
    }
  }

  const handleSubmit = async () => {
    loading.value = true
    error.value = ''

    try {
      if (isEditing.value && editingMetric.value) {
        const payload: UpdateMetricRequest = {
          id: editingMetric.value.id,
          name: form.name,
          description: form.description,
          metric_type: form.metric_type,
          unit: form.unit,
          extra: form.extra || undefined,
        }
        await metricAPI.update(editingMetric.value.id, payload)
        ElMessage.success('指标更新成功')
      } else {
        const payload: CreateMetricRequest = {
          name: form.name,
          description: form.description,
          metric_type: form.metric_type,
          unit: form.unit,
          extra: form.extra || undefined,
        }
        await metricAPI.create(payload)
        ElMessage.success('指标创建成功')
      }

      cancelEdit()
      loadMetrics()
    } catch (err: any) {
      error.value = err.response?.data?.message || '保存失败'
      ElMessage.error(error.value)
    } finally {
      loading.value = false
    }
  }

  const toggleForm = () => {
    showForm.value = !showForm.value
    if (showForm.value) {
      editingMetric.value = null
      resetForm()
    }
  }

  const cancelEdit = () => {
    showForm.value = false
    editingMetric.value = null
    resetForm()
  }

  const editMetric = (metric: Metric) => {
    editingMetric.value = metric
    loadEditingData()
    showForm.value = true
  }

  const viewRecords = async (metric: Metric) => {
    if (selectedMetricId.value === metric.id) {
      selectedMetricId.value = null
      return
    }

    selectedMetricId.value = metric.id
    recordsLoading.value = true
    recordsError.value = ''
    records.value = []

    try {
      const response = await metricAPI.getRecords(metric.id)
      records.value = response.data || []
    } catch (err: any) {
      recordsError.value = err.response?.data?.message || '加载记录失败'
      ElMessage.error(recordsError.value)
    } finally {
      recordsLoading.value = false
    }
  }

  const deleteMetric = async (metric: Metric) => {
    try {
      await ElMessageBox.confirm(
        '确定要删除此指标吗？关联的约束将失效。此操作不可恢复。',
        '删除指标',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }
      )

      await metricAPI.delete(metric.id)
      ElMessage.success('指标已删除')
      if (selectedMetricId.value === metric.id) {
        selectedMetricId.value = null
      }
      loadMetrics()
    } catch (err: any) {
      if (err !== 'cancel') {
        ElMessage.error(err.response?.data?.message || '删除失败')
      }
    }
  }

  const formatMetricType = (type: string) => {
    const map: Record<string, string> = {
      numeric: '数值',
      count: '计数',
      time: '时间',
      boolean: '是否',
    }
    return map[type] || type
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
    loadMetrics()
  })
</script>
