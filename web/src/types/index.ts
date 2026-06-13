export type { ApiResponse } from './api'
import type {
  Category,
  CategoryTreeNode as GeneratedCategoryTreeNode,
} from '@/gen/api/timelog/v1/category'

export type {
  Category,
  CreateCategoryRequest,
  UpdateCategoryRequest,
} from '@/gen/api/timelog/v1/category'
export type {
  Constraint,
  CreateConstraintRequest,
  UpdateConstraintRequest,
} from '@/gen/api/timelog/v1/constraint'
export type {
  Task,
  CreateTaskRequest,
  UpdateTaskRequest,
  TaskStats,
} from '@/gen/api/timelog/v1/task'
export type {
  PasskeyBeginPayload,
  PasskeyCredential,
  PasskeyLoginResponse,
} from '@/gen/api/timelog/v1/auth'
export type {
  Timelog,
  CreateTimelogRequest,
  UpdateTimelogRequest,
} from '@/gen/api/timelog/v1/timelog'
export type { Timelog as TimeLog } from '@/gen/api/timelog/v1/timelog'
export type { CreateTimelogRequest as CreateTimeLogRequest } from '@/gen/api/timelog/v1/timelog'
export type { UpdateTimelogRequest as UpdateTimeLogRequest } from '@/gen/api/timelog/v1/timelog'
export type { PasskeyBeginPayload as PasskeyBeginResponse } from '@/gen/api/timelog/v1/auth'
export type { Category as Tag } from '@/gen/api/timelog/v1/category'

export interface CategoryNode extends Omit<GeneratedCategoryTreeNode, 'category' | 'children'> {
  category: Category
  children: CategoryNode[]
}
