<template>
  <div>
    <div
      class="p-4 hover:bg-bg-elevated flex items-center justify-between"
      :style="{ paddingLeft: `${level * 2 + 1}rem` }"
    >
      <div class="flex items-center space-x-3">
        <!-- 展开/折叠图标 -->
        <button
          v-if="node.children && node.children.length > 0"
          @click="toggleExpand"
          class="text-text-tertiary hover:text-text-secondary"
        >
          <ChevronDownIcon v-if="expanded" class="h-4 w-4" />
          <ChevronRightIcon v-else class="h-4 w-4" />
        </button>
        <span v-else class="w-4"></span>

        <span
          class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium text-bg-surface"
          :style="{ backgroundColor: node.category.color }"
        >
          {{ node.category.name }}
        </span>

        <div class="text-sm">
          <p class="text-text-primary" v-if="node.category.description">
            {{ node.category.description }}
          </p>
          <p class="text-xs text-text-secondary">
            Level {{ node.category.level }}
            <span v-if="node.children && node.children.length > 0"
              >· {{ node.children.length }} sub-categories</span
            >
          </p>
        </div>
      </div>

      <div class="flex items-center space-x-2">
        <button
          v-if="node.category.level < 3"
          @click="$emit('move', node.category)"
          class="text-text-secondary hover:text-text-primary text-sm font-medium"
          title="Move category"
        >
          Move
        </button>
        <button
          @click="$emit('edit', node.category)"
          class="text-brand hover:text-blue-900 text-sm font-medium"
        >
          Edit
        </button>
        <button
          @click="$emit('delete', node.category.id)"
          class="text-danger hover:text-red-900 text-sm font-medium"
        >
          Delete
        </button>
      </div>
    </div>

    <!-- 子分类 -->
    <div v-if="expanded && node.children && node.children.length > 0">
      <CategoryTreeNode
        v-for="child in node.children"
        :key="child.category.id"
        :node="child"
        :level="level + 1"
        @edit="$emit('edit', $event)"
        @delete="$emit('delete', $event)"
        @move="$emit('move', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { ChevronDownIcon, ChevronRightIcon } from '@heroicons/vue/24/outline'
  import type { CategoryNode, Category } from '@/types'

  interface Props {
    node: CategoryNode
    level: number
  }

  defineProps<Props>()

  defineEmits<{
    edit: [category: Category]
    delete: [id: number]
    move: [category: Category]
  }>()

  const expanded = ref(true)

  const toggleExpand = () => {
    expanded.value = !expanded.value
  }
</script>
