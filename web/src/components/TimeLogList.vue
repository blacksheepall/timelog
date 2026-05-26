<template>
  <div class="bg-white shadow rounded-lg">
    <div class="px-6 py-4 border-b border-gray-200">
      <h2 class="text-lg font-medium text-gray-900">Time Logs</h2>
    </div>

    <div v-if="loading" class="p-6 text-center">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <p class="mt-2 text-gray-600">Loading...</p>
    </div>

    <div v-else-if="error" class="p-6 text-center text-red-600">
      {{ error }}
    </div>

    <div v-else-if="timeLogs.length === 0" class="p-6 text-center text-gray-500">
      No time logs found. Create your first one!
    </div>

    <div v-else class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              Start Time
            </th>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              End Time
            </th>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              Duration
            </th>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              Category
            </th>
            <th
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              Remarks
            </th>
            <th
              class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="log in timeLogs" :key="log.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ formatDateTime(log.start_time) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ log.end_time ? formatDateTime(log.end_time) : 'Ongoing' }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ calculateDuration(log.start_time, log.end_time) }}
            </td>
            <td class="px-6 py-4 text-sm text-gray-900">
              <span
                v-if="getCategory(log.category_id)"
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-white"
                :style="{ backgroundColor: getCategory(log.category_id)!.color }"
                :title="getCategory(log.category_id)!.description"
              >
                {{ getCategory(log.category_id)!.name }}
              </span>
              <span v-else class="text-gray-400">—</span>
            </td>
            <td class="px-6 py-4 text-sm text-gray-900 max-w-xs truncate">
              {{ log.remark }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button @click="$emit('edit', log)" class="text-blue-600 hover:text-blue-900 mr-4">
                Edit
              </button>
              <button @click="$emit('delete', log.id)" class="text-red-600 hover:text-red-900">
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
