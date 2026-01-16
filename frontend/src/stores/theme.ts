// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { create } from 'zustand'

type ThemeMode = 'light' | 'dark'

interface ThemeState {
  mode: ThemeMode
  setMode: (mode: ThemeMode) => void
  toggleMode: () => void
}

// Initialize from localStorage
const storedMode = localStorage.getItem('theme-mode') as ThemeMode | null
const initialMode = storedMode === 'light' || storedMode === 'dark' ? storedMode : 'light'

export const useThemeStore = create<ThemeState>((set) => ({
  mode: initialMode,
  setMode: (mode) => {
    localStorage.setItem('theme-mode', mode)
    set({ mode })
  },
  toggleMode: () => {
    set((state) => {
      const newMode = state.mode === 'light' ? 'dark' : 'light'
      localStorage.setItem('theme-mode', newMode)
      return { mode: newMode }
    })
  },
}))
