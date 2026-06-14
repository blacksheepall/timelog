<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-text-primary">约束管理</h1>
      <button
        @click="toggleForm"
        class="inline-flex items-center px-4 py-2 text-sm font-medium text-bg-surface bg-brand border border-transparent rounded-md shadow-sm hover:bg-brand-hover focus:outline-none focus:ring-2 focus:ring-brand"
      >
        <PlusIcon class="h-5 w-5 mr-2" />
        新建约束
      </button>
    </div>

    <!-- 约束创建/编辑表单 -->
    <div v-if="showForm" class="bg-bg-surface shadow-md-md rounded-lg p-6">
      <h2 class="text-lg font-medium text-text-primary mb-6">
        {{ isEditing ? '编辑约束' : '创建新约束' }}
      </h2>

      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div>
          <label for="description" class="block text-sm font-medium text-text-muted mb-2">
            约束描述 *
          </label>
          <textarea
            id="description"
            v-model="form.description"
            rows="3"
            required
            class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            placeholder="描述你的约束，比如：每天学习至少2小时..."
          ></textarea>
        </div>

        <div>
          <label for="punishment_quote" class="block text-sm font-medium text-text-muted mb-2">
            惩罚语录 *
          </label>
          <textarea
            id="punishment_quote"
            v-model="form.punishment_quote"
            rows="2"
            required
            class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            placeholder="如果没有遵守约束，对自己说的话..."
          ></textarea>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label for="start_date" class="block text-sm font-medium text-text-muted mb-2">
              开始日期 *
            </label>
            <input
              id="start_date"
              v-model="form.start_date"
              type="date"
              required
              class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            />
          </div>

          <div>
            <label for="end_date" class="block text-sm font-medium text-text-muted mb-2">
              结束日期
            </label>
            <input
              id="end_date"
              v-model="form.end_date"
              type="date"
              :min="form.start_date"
              class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            />
          </div>
        </div>

        <div class="border-t border-default pt-6">
          <h3 class="text-sm font-medium text-text-primary mb-4">关联指标（可选）</h3>
          <div
            v-if="metrics.length === 0"
            class="mb-4 text-sm text-text-muted bg-bg-elevated rounded-md p-3"
          >
            还没有指标，请先前往
            <router-link to="/metrics" class="text-brand hover:underline">指标页面</router-link>
            创建一个。
          </div>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div>
              <label for="metric_id" class="block text-sm font-medium text-text-muted mb-2">
                指标
              </label>
              <select
                id="metric_id"
                v-model="form.metric_id"
                class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand bg-bg-surface"
              >
                <option :value="null">不关联指标</option>
                <option v-for="m in metrics" :key="m.id" :value="m.id">
                  {{ m.name }} ({{ m.unit }})
                </option>
              </select>
            </div>

            <div>
              <label for="metric_operator" class="block text-sm font-medium text-text-muted mb-2">
                比较方式
              </label>
              <select
                id="metric_operator"
                v-model="form.metric_operator"
                class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand bg-bg-surface"
              >
                <option value="eq">等于</option>
                <option value="ne">不等于</option>
                <option value="gt">大于</option>
                <option value="gte">大于等于</option>
                <option value="lt">小于</option>
                <option value="lte">小于等于</option>
              </select>
            </div>

            <div>
              <label
                for="metric_target_value"
                class="block text-sm font-medium text-text-muted mb-2"
              >
                目标值
              </label>
              <input
                id="metric_target_value"
                v-model.number="form.metric_target_value"
                type="number"
                step="any"
                class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
                placeholder="例如：23.5"
              />
            </div>
          </div>
        </div>

        <div
          v-if="isEditing && !editingTask.is_active"
          class="bg-yellow-50 border border-yellow-200 rounded-md p-4"
        >
          <div class="flex">
            <div class="flex-shrink-0">
              <ExclamationTriangleIcon class="h-5 w-5 text-yellow-400" />
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-yellow-800">约束已完成</h3>
              <div class="mt-2 text-sm text-yellow-700">
                <p>此约束已标记为完成。您可以重新激活它或创建新的约束。</p>
              </div>
            </div>
          </div>
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
            {{ loading ? '保存中...' : isEditing ? '更新约束' : '创建约束' }}
          </button>
        </div>
      </form>
    </div>

    <!-- 约束列表 -->
    <div class="bg-bg-surface shadow-md-md rounded-lg">
      <div class="px-6 py-4 border-b border-default">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-medium text-text-primary">约束列表</h2>
          <div class="flex items-center space-x-4">
            <label class="flex items-center">
              <input
                v-model="showOnlyActive"
                type="checkbox"
                class="h-4 w-4 text-brand focus:ring-brand border-default rounded"
                @change="loadConstraints"
              />
              <span class="ml-2 text-sm text-text-muted">只显示活跃约束</span>
            </label>
          </div>
        </div>
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

      <div v-else-if="constraints.length === 0" class="p-8 text-center text-text-secondary">
        <DocumentTextIcon class="h-12 w-12 mx-auto mb-4 text-text-tertiary" />
        <p>暂无约束。创建你的第一个约束吧！</p>
      </div>

      <div v-else class="divide-y divide-default">
        <div
          v-for="constraint in constraints"
          :key="constraint.id"
          class="p-6 hover:bg-bg-elevated"
          :class="{ 'opacity-60': !constraint.is_active }"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <h3
                  class="text-lg font-medium"
                  :class="
                    constraint.is_active ? 'text-text-primary' : 'text-text-secondary line-through'
                  "
                >
                  {{ constraint.description }}
                </h3>
                <span
                  v-if="constraint.is_active"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"
                >
                  活跃
                </span>
                <span
                  v-else
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-bg-elevated text-text-primary"
                >
                  已完成
                </span>
              </div>

              <div class="bg-danger-bg border border-danger-border rounded-md p-3 mb-3">
                <div class="flex">
                  <div class="flex-shrink-0">
                    <ExclamationTriangleIcon class="h-5 w-5 text-red-400" />
                  </div>
                  <div class="ml-3">
                    <p class="text-sm text-red-700">{{ constraint.punishment_quote }}</p>
                  </div>
                </div>
              </div>

              <div class="flex items-center space-x-4 text-sm text-text-secondary">
                <span>开始日期: {{ formatDate(constraint.start_date) }}</span>
                <span v-if="constraint.end_date">
                  结束日期: {{ formatDate(constraint.end_date) }}
                </span>
                <span v-if="constraint.end_reason"> 结束理由: {{ constraint.end_reason }} </span>
              </div>

              <div
                v-if="constraint.metric_id"
                class="mt-3 inline-flex items-center px-3 py-1.5 rounded-md bg-brand-bg border border-brand-border text-sm"
              >
                <ChartBarIcon class="h-4 w-4 text-brand mr-2" />
                <span class="text-text-secondary">
                  指标：{{ metricMap[constraint.metric_id]?.name || '未知' }}
                  {{ operatorMap[constraint.metric_operator] || constraint.metric_operator }}
                  {{ constraint.metric_target_value }}
                  {{ metricMap[constraint.metric_id]?.unit || '' }}
                </span>
              </div>
            </div>

            <div class="flex items-center space-x-2 ml-4">
              <button
                v-if="constraint.metric_id"
                @click="evaluateConstraint(constraint)"
                class="inline-flex items-center px-3 py-1.5 text-sm font-medium text-brand bg-brand-bg border border-brand-border rounded-md hover:bg-brand-bg focus:outline-none focus:ring-2 focus:ring-brand"
              >
                <ChartBarIcon class="h-4 w-4 mr-1" />
                评估
              </button>
              <button
                v-else
                @click="editConstraint(constraint)"
                class="inline-flex items-center px-3 py-1.5 text-sm font-medium text-brand bg-brand-bg border border-brand-border rounded-md hover:bg-brand-bg focus:outline-none focus:ring-2 focus:ring-brand"
              >
                <ChartBarIcon class="h-4 w-4 mr-1" />
                关联指标
              </button>
              <button
                v-if="constraint.is_active"
                @click="completeConstraint(constraint)"
                class="inline-flex items-center px-3 py-1.5 text-sm font-medium text-green-700 bg-green-100 border border-green-300 rounded-md hover:bg-green-200 focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                <CheckCircleIcon class="h-4 w-4 mr-1" />
                完成
              </button>
              <button
                v-else
                @click="reactivateConstraint(constraint)"
                class="inline-flex items-center px-3 py-1.5 text-sm font-medium text-blue-700 bg-blue-100 border border-blue-300 rounded-md hover:bg-blue-200 focus:outline-none focus:ring-2 focus:ring-brand"
              >
                <ArrowPathIcon class="h-4 w-4 mr-1" />
                重新激活
              </button>
              <button
                @click="editConstraint(constraint)"
                class="inline-flex items-center px-3 py-1.5 text-sm font-medium text-text-muted bg-bg-surface border border-default rounded-md hover:bg-bg-elevated focus:outline-none focus:ring-2 focus:ring-brand"
              >
                <PencilIcon class="h-4 w-4" />
              </button>
              <button
                @click="deleteConstraint(constraint)"
                class="inline-flex items-center px-3 py-1.5 text-sm font-medium text-danger bg-bg-surface border border-danger-border rounded-md hover:bg-danger-bg focus:outline-none focus:ring-2 focus:ring-danger"
              >
                <TrashIcon class="h-4 w-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, reactive, computed, onMounted } from 'vue'
  import { useSettings } from '@/composables/useSettings'
  import {
    PlusIcon,
    PencilIcon,
    TrashIcon,
    CheckCircleIcon,
    ArrowPathIcon,
    ExclamationTriangleIcon,
    DocumentTextIcon,
    ChartBarIcon,
  } from '@heroicons/vue/24/outline'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { constraintAPI, metricAPI } from '@/api'

  const loading = ref(false)
  const error = ref(null)
  const showForm = ref(false)
  const editingTask = ref(null)
  const constraints = ref([])
  const metrics = ref([])

  // Use settings from composable
  const { timeLogShowOnlyActive: showOnlyActive } = useSettings()

  const isEditing = computed(() => !!editingTask.value)

  const metricMap = computed(() => {
    const map = {}
    for (const m of metrics.value) {
      map[m.id] = m
    }
    return map
  })

  const operatorMap = {
    eq: '等于',
    ne: '不等于',
    gt: '大于',
    gte: '大于等于',
    lt: '小于',
    lte: '小于等于',
  }

  const form = reactive({
    description: '',
    punishment_quote: '',
    start_date: '',
    end_date: '',
    metric_id: null,
    metric_operator: 'gte',
    metric_target_value: 0,
  })

  const resetForm = () => {
    form.description = ''
    form.punishment_quote = ''
    form.start_date = new Date().toISOString().split('T')[0] // Today's date
    form.end_date = ''
    form.metric_id = null
    form.metric_operator = 'gte'
    form.metric_target_value = 0
  }

  const loadEditingData = () => {
    if (editingTask.value) {
      form.description = editingTask.value.description
      form.punishment_quote = editingTask.value.punishment_quote
      form.start_date = editingTask.value.start_date.split('T')[0]
      form.end_date = editingTask.value.end_date ? editingTask.value.end_date.split('T')[0] : ''
      form.metric_id = editingTask.value.metric_id || null
      form.metric_operator = editingTask.value.metric_operator || 'gte'
      form.metric_target_value = editingTask.value.metric_target_value || 0
    } else {
      resetForm()
    }
  }

  const loadConstraints = async () => {
    loading.value = true
    error.value = null

    try {
      const response = await constraintAPI.getAll(showOnlyActive.value)
      constraints.value = response.data || []
    } catch (err) {
      error.value = err.response?.data?.message || '加载约束失败'
      ElMessage.error(error.value)
    } finally {
      loading.value = false
    }
  }

  const loadMetrics = async () => {
    try {
      const response = await metricAPI.getAll()
      metrics.value = response.data || []
    } catch (err) {
      console.error('加载指标失败', err)
    }
  }

  const handleSubmit = async () => {
    loading.value = true
    error.value = null

    try {
      const formData = {
        description: form.description,
        punishment_quote: form.punishment_quote,
        start_date: form.start_date,
        end_date: form.end_date || null,
        metric_id: form.metric_id || undefined,
        metric_operator: form.metric_id ? form.metric_operator : undefined,
        metric_target_value: form.metric_id ? form.metric_target_value : undefined,
      }

      if (isEditing.value) {
        await constraintAPI.update(editingTask.value.id, formData)
        ElMessage.success('约束更新成功')
      } else {
        await constraintAPI.create(formData)
        ElMessage.success('约束创建成功')
      }

      cancelEdit()
      loadConstraints()
    } catch (err) {
      error.value = err.response?.data?.message || '保存失败'
      ElMessage.error(error.value)
    } finally {
      loading.value = false
    }
  }

  const toggleForm = () => {
    showForm.value = !showForm.value
    if (showForm.value) {
      editingTask.value = null
      resetForm()
    }
  }

  const cancelEdit = () => {
    showForm.value = false
    editingTask.value = null
    resetForm()
  }

  const editConstraint = constraint => {
    editingTask.value = constraint
    loadEditingData()
    showForm.value = true
  }

  const completeConstraint = async constraint => {
    try {
      await ElMessageBox.prompt('请输入完成此约束的理由：', '完成约束', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputPattern: /\S+/,
        inputErrorMessage: '请输入完成理由',
      }).then(async ({ value }) => {
        await constraintAPI.complete(constraint.id, value)
        ElMessage.success('约束已完成')
        loadConstraints()
      })
    } catch (err) {
      if (err !== 'cancel') {
        ElMessage.error(err.response?.data?.message || '操作失败')
      }
    }
  }

  const reactivateConstraint = async constraint => {
    try {
      await ElMessageBox.confirm('确定要重新激活此约束吗？', '重新激活', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })

      await constraintAPI.reactivate(constraint.id)
      ElMessage.success('约束已重新激活')
      loadConstraints()
    } catch (err) {
      if (err !== 'cancel') {
        ElMessage.error(err.response?.data?.message || '操作失败')
      }
    }
  }

  const deleteConstraint = async constraint => {
    try {
      await ElMessageBox.confirm('确定要删除此约束吗？此操作不可恢复。', '删除约束', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })

      await constraintAPI.delete(constraint.id)
      ElMessage.success('约束已删除')
      loadConstraints()
    } catch (err) {
      if (err !== 'cancel') {
        ElMessage.error(err.response?.data?.message || '删除失败')
      }
    }
  }

  const evaluateConstraint = async constraint => {
    try {
      const response = await constraintAPI.evaluate(constraint.id)
      const result = response.data
      const metric = metricMap.value[constraint.metric_id]
      const status = result.passed ? '通过 ✅' : '未通过 ❌'
      ElMessageBox.alert(
        `<div class="space-y-2">
          <p><strong>状态：</strong>${status}</p>
          <p><strong>当前值：</strong>${result.actual} ${metric?.unit || ''}</p>
          <p><strong>目标值：</strong>${result.target} ${metric?.unit || ''}</p>
          <p><strong>规则：</strong>${operatorMap[result.operator] || result.operator}</p>
        </div>`,
        '约束评估结果',
        {
          confirmButtonText: '确定',
          dangerouslyUseHTMLString: true,
        }
      )
    } catch (err) {
      ElMessage.error(err.response?.data?.message || '评估失败')
    }
  }

  const formatDate = dateString => {
    if (!dateString) return ''
    return new Date(dateString).toLocaleDateString('zh-CN')
  }

  onMounted(() => {
    loadConstraints()
    loadMetrics()
  })
</script>
