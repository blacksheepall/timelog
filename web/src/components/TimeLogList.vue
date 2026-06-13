<template>
  <div class="bg-bg-surface shadow-md-md rounded-lg">
    <div class="px-6 py-4 border-b border-default">
      <h2 class="text-lg font-medium text-text-primary">Time Logs</h2>
    </div>

    <div v-if="loading" class="p-6 text-center">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-brand"></div>
      <p class="mt-2 text-text-secondary">Loading...</p>
    </div>

    <div v-else-if="error" class="p-6 text-center text-danger">
      {{ error }}
    </div>

    <div v-else-if="timeLogs.length === 0" class="p-6 text-center text-text-secondary">
      No time logs found. Create your first one!
    </div>

    <div v-else class="overflow-x-auto">
      <table class="min-w-full divide-y divide-default">
        <thead class="bg-bg-elevated">
          <tr>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-text-secondary uppercase tracking-wider"
            >
              Start Time
            </th>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-text-secondary uppercase tracking-wider"
            >
              End Time
            </th>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-text-secondary uppercase tracking-wider"
            >
              Duration
            </th>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-text-secondary uppercase tracking-wider"
            >
              Category
            </th>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-text-secondary uppercase tracking-wider"
            >
              Remarks
            </th>
            <th
              class="px-6 py-3 text-right text-xs font-medium text-text-secondary uppercase tracking-wider"
            >
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-bg-surface divide-y divide-default">
          <tr v-for="log in timeLogs" :key="log.id" class="hover:bg-bg-elevated">
            <td class="px-6 py-4 whitespace-nowrap text-sm text-text-primary">
              {{ formatDateTime(log.start_time) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-text-primary">
              {{ log.end_time ? formatDateTime(log.end_time) : 'Ongoing' }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-text-primary">
              {{ calculateDuration(log.start_time, log.end_time) }}
            </td>
            <td class="px-6 py-4 text-sm text-text-primary">
              <span
                v-if="getCategory(log.category_id)"
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-bg-surface"
                :style="{ backgroundColor: getCategory(log.category_id)!.color }"
                :title="getCategory(log.category_id)!.description"
              >
                {{ getCategory(log.category_id)!.name }}
              </span>
              <span v-else class="text-text-tertiary">—</span>
            </td>
            <td class="px-6 py-4 text-sm text-text-primary max-w-xs truncate">
              {{ log.remark }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button @click="$emit('edit', log)" class="text-brand hover:text-blue-900 mr-4">
                Edit
              </button>
              <button @click="$emit('delete', log.id)" class="text-danger hover:text-red-900">
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
  import type { TimeLog, Category } from '@/types'
  import { formatDateTime, calculateDuration } from '@/utils/date'

  const props = defineProps<{
    timeLogs: TimeLog[]
    categories: Category[]
    loading: boolean
    error: string | null
  }>()

  const getCategory = (categoryId: number): Category | undefined => {
    return props.categories.find(c => c.id === categoryId)
  }

  defineEmits<{
    edit: [log: TimeLog]
    delete: [id: number]
  }>()
</script>
