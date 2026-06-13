<template>
  <div class="bg-bg-surface shadow-md-md rounded-lg p-6">
    <h2 class="text-lg font-medium text-text-primary mb-6">
      {{ isEditing ? 'Edit Time Log' : 'Create New Time Log' }}
    </h2>

    <form @submit.prevent="handleSubmit" class="space-y-6">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label for="start_time" class="block text-sm font-medium text-text-muted mb-2">
            Start Time *
          </label>
          <div class="flex space-x-2">
            <input
              id="start_time"
              ref="startTimeInput"
              v-model="form.start_time"
              type="datetime-local"
              required
              class="flex-1 px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            />
            <button
              type="button"
              @click="setCurrentTime('start')"
              class="px-3 py-2 text-xs font-medium text-brand bg-brand-bg border border-blue-200 rounded-md hover:bg-blue-100 focus:outline-none focus:ring-2 focus:ring-brand whitespace-nowrap transition-colors"
              title="Set to current time"
            >
              <div class="flex items-center space-x-1">
                <ClockIcon class="h-4 w-4" />
                <span class="hidden sm:inline">Now</span>
              </div>
            </button>
          </div>
        </div>

        <div>
          <label for="end_time" class="block text-sm font-medium text-text-muted mb-2">
            End Time
          </label>
          <div class="flex space-x-2">
            <input
              id="end_time"
              v-model="form.end_time"
              type="datetime-local"
              class="flex-1 px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            />
            <button
              type="button"
              @click="setCurrentTime('end')"
              class="px-3 py-2 text-xs font-medium text-brand bg-brand-bg border border-blue-200 rounded-md hover:bg-blue-100 focus:outline-none focus:ring-2 focus:ring-brand whitespace-nowrap transition-colors"
              title="Set to current time"
            >
              <div class="flex items-center space-x-1">
                <ClockIcon class="h-4 w-4" />
                <span class="hidden sm:inline">Now</span>
              </div>
            </button>
            <button
              v-if="form.end_time"
              type="button"
              @click="clearEndTime"
              class="px-3 py-2 text-xs font-medium text-text-secondary bg-bg-elevated border border-default rounded-md hover:bg-bg-base focus:outline-none focus:ring-2 focus:ring-brand transition-colors"
              title="Clear end time"
            >
              <XMarkIcon class="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>

      <div>
        <label for="category_id" class="block text-sm font-medium text-text-muted mb-2">
          Category *
        </label>
        <select
          id="category_id"
          v-model="form.category_id"
          required
          class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
        >
          <option value="" disabled>Select a category</option>
          <option
            v-for="category in availableCategories"
            :key="category.id"
            :value="category.id"
            :style="{ color: category.color }"
          >
            {{ category.path === '/' ? '' : category.path.replace(/\//g, ' / ') + ' / '
            }}{{ category.name }}
          </option>
        </select>
      </div>

      <div>
        <label for="task_id" class="block text-sm font-medium text-text-muted mb-2">
          Associated Task (Optional)
        </label>
        <select
          id="task_id"
          v-model="form.task_id"
          class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
        >
          <option value="" disabled>Select a task (optional)</option>
          <option value="">No task</option>
          <option v-for="task in availableTasks" :key="task.id" :value="task.id">
            {{ task.title }} ({{ task.category?.name }})
          </option>
        </select>
      </div>

      <div>
        <label for="remark" class="block text-sm font-medium text-text-muted mb-2"> Remark </label>
        <textarea
          id="remark"
          v-model="form.remark"
          rows="3"
          placeholder="Add any notes or description..."
          class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
        ></textarea>
      </div>

      <div class="flex justify-end space-x-4">
        <button
          v-if="isEditing"
          type="button"
          @click="handleCancel"
          class="px-4 py-2 text-sm font-medium text-text-muted bg-bg-surface border border-default rounded-md shadow-sm hover:bg-bg-elevated focus:outline-none focus:ring-2 focus:ring-brand"
        >
          Cancel
        </button>
        <button
          type="submit"
          :disabled="submitting"
          class="px-4 py-2 text-sm font-medium text-bg-surface bg-brand border border-transparent rounded-md shadow-sm hover:bg-brand-hover focus:outline-none focus:ring-2 focus:ring-brand disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ submitting ? 'Saving...' : isEditing ? 'Update' : 'Create' }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
  import { reactive, watch, computed, ref, nextTick } from 'vue'
  import { ClockIcon, XMarkIcon } from '@heroicons/vue/24/outline'
  import type { TimeLog, Category, Task, CreateTimeLogRequest, UpdateTimeLogRequest } from '@/types'
  import { formatDateTimeLocal, formatUTCToLocal, formatLocalToUTC } from '@/utils/date'

  const props = defineProps<{
    editingLog?: TimeLog
    submitting: boolean
    availableCategories: Category[]
    availableTasks: Task[]
    lastEndTime?: string | null
  }>()

  const emit = defineEmits<{
    create: [data: CreateTimeLogRequest]
    update: [data: UpdateTimeLogRequest]
    cancel: []
  }>()

  const isEditing = computed(() => !!props.editingLog)

  // Template ref for auto-focus
  const startTimeInput = ref<HTMLInputElement>()

  const form = reactive<{
    start_time: string
    end_time: string
    category_id: number | ''
    task_id: number | ''
    remark: string
  }>({
    start_time: '',
    end_time: '',
    category_id: '',
    task_id: '',
    remark: '',
  })

  const resetForm = () => {
    // 如果有上一个 timelog 的结束时间，使用它作为默认的开始时间
    if (props.lastEndTime) {
      form.start_time = formatUTCToLocal(props.lastEndTime)
    } else {
      form.start_time = formatDateTimeLocal()
    }
    form.end_time = ''
    form.category_id = ''
    form.task_id = ''
    form.remark = ''
  }

  const setCurrentTime = (field: 'start' | 'end') => {
    const currentTime = formatDateTimeLocal()
    if (field === 'start') {
      form.start_time = currentTime
    } else {
      form.end_time = currentTime
    }
  }

  const clearEndTime = () => {
    form.end_time = ''
  }

  const loadEditingData = async () => {
    if (props.editingLog) {
      form.start_time = formatUTCToLocal(props.editingLog.start_time)
      form.end_time = props.editingLog.end_time ? formatUTCToLocal(props.editingLog.end_time) : ''
      form.category_id = props.editingLog.category_id
      form.task_id = props.editingLog.task_id || ''
      form.remark = props.editingLog.remark || ''
    } else {
      resetForm()
    }

    // Auto-focus the start time input
    await nextTick()
    startTimeInput.value?.focus()
  }

  const handleSubmit = () => {
    if (form.category_id === '') {
      return // Category is required
    }

    const payload = {
      start_time: formatLocalToUTC(form.start_time),
      end_time: form.end_time ? formatLocalToUTC(form.end_time) : undefined,
      category_id: Number(form.category_id),
      task_id: form.task_id ? Number(form.task_id) : undefined,
      remark: form.remark,
    }

    if (isEditing.value) {
      emit('update', payload)
    } else {
      emit('create', payload)
    }
  }

  const handleCancel = () => {
    emit('cancel')
    resetForm()
  }

  watch(() => props.editingLog, loadEditingData, { immediate: true })
</script>
