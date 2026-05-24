export interface Category {
  id: number
  name: string
  color: string
  description: string
  parent_id?: number | null
  level: number
  sort_order: number
  path: string
  created_at: string
  updated_at: string
}

export interface CategoryNode {
  category: Category
  children?: CategoryNode[]
}

export interface TimeLog {
  id: number
  start_time: string
  end_time?: string | null
  category_id: number
  category: Category
  task_id?: number | null
  task?: Task | null
  remark?: string | null
  created_at: string
  updated_at: string
  deleted_at?: string | null
}

export interface CreateTimeLogRequest {
  start_time: string
  end_time?: string
  category_id: number
  task_id?: number | null
  remark: string
}

export interface UpdateTimeLogRequest {
  start_time?: string
  end_time?: string
  category_id?: number
  task_id?: number | null
  remark?: string
}

export interface Task {
  id: number
  title: string
  description: string
  category_id: number
  category?: Category
  due_date: string
  estimated_minutes: number
  is_completed: boolean
  completed_at?: string | null
  is_suspended: boolean
  created_at: string
  updated_at: string
}

export interface CreateTaskRequest {
  title: string
  description?: string
  category_id: number
  due_date: string
  estimated_minutes: number
}

export interface UpdateTaskRequest {
  title?: string
  description?: string
  category_id?: number
  due_date?: string
  estimated_minutes?: number
  is_completed?: boolean
}

export interface TaskStats {
  total_tasks: number
  completed_tasks: number
  completion_rate: number
}

export interface ApiResponse<T> {
  data: T
  message?: string
  status: number
}

export interface PasskeyBeginResponse<T> {
  session_id: string
  data: T
}

export interface PasskeyCredential {
  id: number
  device_name?: string
  created_at: string
}

export interface PasskeyLoginResponse {
  token: string
  token_type: string
  expires_in: number
}

export interface Constraint {
  id: number
  description: string
  end_reason?: string
  punishment_quote: string
  start_date: string
  end_date?: string | null
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateConstraintRequest {
  description: string
  punishment_quote: string
  start_date: string
  end_date?: string
}

export interface UpdateConstraintRequest {
  description?: string
  punishment_quote?: string
  start_date?: string
  end_date?: string
}

// Backward compatibility - Tag is deprecated, use Category instead
/** @deprecated Use Category instead */
export type Tag = Category
