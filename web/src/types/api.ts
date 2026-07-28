/**
 * Frozen API response envelope — mirrors router/apiresponse.go on the backend.
 *
 * Contract:
 * - Success: HTTP 200, body { data: T, message: string, status: 200 }
 * - Error: non-2xx, body { data: null, message: string, status: <same code> }
 * - The frontend reads HTTP status (axios) and `message` on errors only;
 *   envelope `status` is convention, never read.
 * - Backend enforces this shape via router.TestEnvelopeContract.
 *
 * Any change must stay in sync with router/apiresponse.go.
 */
export interface ApiResponse<T> {
  data: T
  message?: string
  status: number
}
