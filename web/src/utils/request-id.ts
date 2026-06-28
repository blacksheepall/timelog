const REQUEST_ID_HEADER = 'X-Request-ID'

let lastRequestID = ''

export function generateRequestID(): string {
  if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
    try {
      return crypto.randomUUID()
    } catch {
      // fall through to manual generation
    }
  }

  // Minimal v4-style UUID fallback.
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
    const r = (Math.random() * 16) | 0
    const v = c === 'x' ? r : (r & 0x3) | 0x8
    return v.toString(16)
  })
}

export function getLastRequestID(): string {
  return lastRequestID
}

export function setLastRequestID(id: string): void {
  lastRequestID = id
}

export function getRequestIDHeader(): string {
  return REQUEST_ID_HEADER
}
