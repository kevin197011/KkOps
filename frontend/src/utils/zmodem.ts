// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

/**
 * Detect ZMODEM protocol sequence in terminal output
 */
export function detectZmodemSequence(data: string): {
  detected: boolean
  direction: 'upload' | 'download' | ''
  mode: 'zmodem' | 'trzsz' | ''
} {
  // Classic ZMODEM: **B (upload/rz) or **G (download/sz)
  if (data.includes('**B')) {
    return { detected: true, direction: 'upload', mode: 'zmodem' }
  }
  if (data.includes('**G')) {
    return { detected: true, direction: 'download', mode: 'zmodem' }
  }

  // trzsz protocol
  if (data.includes('::TRZSZ:')) {
    if (data.includes('TRANSFER:UPLOAD') || data.includes('trz')) {
      return { detected: true, direction: 'upload', mode: 'trzsz' }
    }
    if (data.includes('TRANSFER:DOWNLOAD') || data.includes('tsz')) {
      return { detected: true, direction: 'download', mode: 'trzsz' }
    }
  }
  return { detected: false, direction: '', mode: '' }
}

/**
 * Format bytes to human-readable string
 */
export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

/**
 * Calculate transfer progress percentage
 */
export function calculateProgress(transferred: number, total: number): number {
  if (total === 0) return 0
  return Math.min(100, Math.round((transferred / total) * 100))
}
