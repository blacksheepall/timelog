import axios from 'axios'
import { clearAuthToken, getAuthToken } from '@/utils/auth'
import { generateRequestID, getLastRequestID, setLastRequestID } from '@/utils/request-id'
import type { ApiResponse } from '@/types/api'
import type { CategoryNode } from '@/types'
import type {
  Timelog,
  CreateTimelogRequest,
  UpdateTimelogRequest,
} from '@/gen/api/timelog/v1/timelog'
import type {
  Category,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  MoveCategoryRequest,
} from '@/gen/api/timelog/v1/category'
import type {
  Task,
  CreateTaskRequest,
  UpdateTaskRequest,
  TaskStats,
} from '@/gen/api/timelog/v1/task'
import type {
  Constraint,
  CreateConstraintRequest,
  UpdateConstraintRequest,
} from '@/gen/api/timelog/v1/constraint'
import type {
  Metric,
  MetricRecord,
  CreateMetricRequest,
  UpdateMetricRequest,
} from '@/gen/api/timelog/v1/metric'
import type {
  PasskeyBeginPayload,
  PasskeyCredential,
  PasskeyLoginResponse,
} from '@/gen/api/timelog/v1/auth'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(config => {
  const token = getAuthToken()
  if (token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${token}`
  }

  const reqID = generateRequestID()
  setLastRequestID(reqID)
  config.headers = config.headers || {}
  config.headers['X-Request-ID'] = reqID

  return config
})

// Store notification handler to be set by App.vue
let notificationHandler: ((type: 'success' | 'error', message: string) => void) | null = null

export const setNotificationHandler = (
  handler: (type: 'success' | 'error', message: string) => void
) => {
  notificationHandler = handler
}

api.interceptors.response.use(
  response => {
    const reqID = response.headers?.['x-request-id']
    if (reqID) {
      setLastRequestID(reqID as string)
    }
    return response
  },
  error => {
    const reqID =
      error.response?.headers?.['x-request-id'] ||
      error.config?.headers?.['X-Request-ID'] ||
      getLastRequestID()
    console.error(`API Error (request_id: ${reqID}):`, error)

    // Check if the error is a timeout or cancellation error
    // These can have various code/message combinations depending on the scenario
    const isTimeout =
      error.code === 'ECONNABORTED' || // Request abort (timeout)
      error.message?.toLowerCase().includes('timeout') // Case-insensitive message check

    const isCancellation = axios.isCancel(error) // Explicit cancellation via axios cancel token/AbortController

    if (isTimeout && !isCancellation && notificationHandler) {
      // Show browser notification for timeout (but not for explicit cancellations)
      notificationHandler('error', 'Request timeout - The server took too long to respond')
    }

    // Handle 401 Unauthorized globally: clear token and redirect to login
    if (error.response?.status === 401) {
      clearAuthToken()
      if (typeof window !== 'undefined' && window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }

    return Promise.reject(error)
  }
)

export const timelogAPI = {
  getAll: (): Promise<ApiResponse<Timelog[]>> => api.get('/timelogs').then(res => res.data),

  getRecent: (limit: number = 5): Promise<ApiResponse<Timelog[]>> =>
    api
      .get(`/timelogs?limit=${limit}&order=${encodeURIComponent('created_at DESC')}`)
      .then(res => res.data),

  getById: (id: number): Promise<ApiResponse<Timelog>> =>
    api.get(`/timelogs/${id}`).then(res => res.data),

  create: (data: CreateTimelogRequest): Promise<ApiResponse<Timelog>> =>
    api.post('/timelogs', data).then(res => res.data),

  update: (id: number, data: UpdateTimelogRequest): Promise<ApiResponse<Timelog>> =>
    api.put(`/timelogs/${id}`, data).then(res => res.data),

  delete: (id: number): Promise<ApiResponse<null>> =>
    api.delete(`/timelogs/${id}`).then(res => res.data),
}

export const categoryAPI = {
  getAll: (): Promise<ApiResponse<Category[]>> => api.get('/categories').then(res => res.data),

  getTree: (): Promise<ApiResponse<CategoryNode[]>> =>
    api.get('/categories/tree').then(res => res.data),

  getByLevel: (level: number): Promise<ApiResponse<Category[]>> =>
    api.get(`/categories?level=${level}`).then(res => res.data),

  getByParentId: (parentId: number): Promise<ApiResponse<Category[]>> =>
    api.get(`/categories?parent_id=${parentId}`).then(res => res.data),

  getById: (id: number): Promise<ApiResponse<Category>> =>
    api.get(`/categories/${id}`).then(res => res.data),

  create: (data: CreateCategoryRequest): Promise<ApiResponse<Category>> =>
    api.post('/categories', data).then(res => res.data),

  update: (id: number, data: UpdateCategoryRequest): Promise<ApiResponse<Category>> =>
    api.put(`/categories/${id}`, data).then(res => res.data),

  delete: (id: number): Promise<ApiResponse<null>> =>
    api.delete(`/categories/${id}`).then(res => res.data),

  move: (id: number, parentId: number | null): Promise<ApiResponse<null>> => {
    const body: MoveCategoryRequest = parentId === null ? {} : { parent_id: parentId }
    return api.post(`/categories/${id}/move`, body).then(res => res.data)
  },
}

// Backward compatibility - tagAPI is deprecated, use categoryAPI instead
/** @deprecated Use categoryAPI instead */
export const tagAPI = categoryAPI

export const taskAPI = {
  getAll: (
    date?: string,
    includeSuspended?: boolean,
    includeCompleted?: boolean
  ): Promise<ApiResponse<Task[]>> => {
    const params = new URLSearchParams()
    if (date) params.append('date', date)
    if (includeSuspended) params.append('include_suspended', 'true')
    if (includeCompleted) params.append('include_completed', 'true')
    const url = `/tasks${params.toString() ? '?' + params.toString() : ''}`
    return api.get(url).then(res => res.data)
  },

  getById: (id: number): Promise<ApiResponse<Task>> =>
    api.get(`/tasks/${id}`).then(res => res.data),

  create: (data: CreateTaskRequest): Promise<ApiResponse<Task>> =>
    api.post('/tasks', data).then(res => res.data),

  update: (id: number, data: UpdateTaskRequest): Promise<ApiResponse<Task>> =>
    api.put(`/tasks/${id}`, data).then(res => res.data),

  delete: (id: number): Promise<ApiResponse<null>> =>
    api.delete(`/tasks/${id}`).then(res => res.data),

  complete: (id: number): Promise<ApiResponse<null>> =>
    api.post(`/tasks/${id}/complete`).then(res => res.data),

  incomplete: (id: number): Promise<ApiResponse<null>> =>
    api.post(`/tasks/${id}/incomplete`).then(res => res.data),

  suspend: (id: number): Promise<ApiResponse<null>> =>
    api.post(`/tasks/${id}/suspend`).then(res => res.data),

  unsuspend: (id: number): Promise<ApiResponse<null>> =>
    api.post(`/tasks/${id}/unsuspend`).then(res => res.data),

  getStats: (date: string): Promise<ApiResponse<TaskStats>> =>
    api.get(`/tasks/stats/${date}`).then(res => res.data),
}

export const constraintAPI = {
  getAll: (active?: boolean): Promise<ApiResponse<Constraint[]>> => {
    const url = active !== undefined ? `/constraints?active=${active}` : '/constraints'
    return api.get(url).then(res => res.data)
  },

  getById: (id: number): Promise<ApiResponse<Constraint>> =>
    api.get(`/constraints/${id}`).then(res => res.data),

  create: (data: CreateConstraintRequest): Promise<ApiResponse<Constraint>> =>
    api.post('/constraints', data).then(res => res.data),

  update: (id: number, data: UpdateConstraintRequest): Promise<ApiResponse<Constraint>> =>
    api.put(`/constraints/${id}`, data).then(res => res.data),

  delete: (id: number): Promise<ApiResponse<null>> =>
    api.delete(`/constraints/${id}`).then(res => res.data),

  complete: (id: number, endReason: string): Promise<ApiResponse<null>> =>
    api.post(`/constraints/${id}/complete`, { end_reason: endReason }).then(res => res.data),

  reactivate: (id: number): Promise<ApiResponse<null>> =>
    api.post(`/constraints/${id}/reactivate`).then(res => res.data),

  evaluate: (id: number): Promise<ApiResponse<any>> =>
    api.get(`/constraints/${id}/evaluation`).then(res => res.data),
}

export const metricAPI = {
  getAll: (): Promise<ApiResponse<Metric[]>> => api.get('/metrics').then(res => res.data),

  getById: (id: number): Promise<ApiResponse<Metric>> =>
    api.get(`/metrics/${id}`).then(res => res.data),

  getRecords: (id: number): Promise<ApiResponse<MetricRecord[]>> =>
    api.get(`/metrics/${id}/records`).then(res => res.data),

  create: (data: CreateMetricRequest): Promise<ApiResponse<Metric>> =>
    api.post('/metrics', data).then(res => res.data),

  update: (id: number, data: UpdateMetricRequest): Promise<ApiResponse<Metric>> =>
    api.put(`/metrics/${id}`, data).then(res => res.data),

  delete: (id: number): Promise<ApiResponse<null>> =>
    api.delete(`/metrics/${id}`).then(res => res.data),
}

export const configAPI = {
  getTimezone: (): Promise<ApiResponse<{ timezone: string }>> =>
    api.get('/timezone').then(res => res.data),
}

export const passkeyAPI = {
  registerBegin: (
    tempPassword: string,
    deviceName?: string
  ): Promise<ApiResponse<PasskeyBeginPayload>> =>
    api
      .post('/passkey/register/begin', {
        temp_password: tempPassword,
        device_name: deviceName,
      })
      .then(res => res.data),

  registerFinish: (
    sessionId: string,
    response: any,
    deviceName?: string
  ): Promise<ApiResponse<PasskeyCredential>> =>
    api
      .post('/passkey/register/finish', {
        session_id: sessionId,
        response,
        device_name: deviceName,
      })
      .then(res => res.data),

  loginBegin: (): Promise<ApiResponse<PasskeyBeginPayload>> =>
    api.post('/passkey/login/begin').then(res => res.data),

  loginFinish: (sessionId: string, response: any): Promise<ApiResponse<PasskeyLoginResponse>> =>
    api
      .post('/passkey/login/finish', {
        session_id: sessionId,
        response,
      })
      .then(res => res.data),

  listCredentials: (): Promise<ApiResponse<PasskeyCredential[]>> =>
    api.get('/passkey/credentials').then(res => res.data),

  deleteCredential: (id: number): Promise<ApiResponse<null>> =>
    api.delete(`/passkey/credentials/${id}`).then(res => res.data),
}

export default api
