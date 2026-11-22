import { create } from 'zustand'

/** Режим оформления приложения. */
export type Theme = 'light' | 'dark'

/** Ключ localStorage, где хранится выбранная тема. */
export const THEME_STORAGE_KEY = 'pantheonv2.theme'

interface ThemeState {
  /** Текущая тема. */
  theme: Theme
  /** Установить конкретную тему и сохранить её. */
  setTheme: (theme: Theme) => void
  /** Переключить тему на противоположную и сохранить. */
  toggleTheme: () => void
}

/** Применяет тему на корне документа (переключает блок переменных в index.css). */
function applyTheme(theme: Theme) {
  document.documentElement.dataset.theme = theme
}

// Синхронно читаем сохранённую тему перед первым render, чтобы не было «мигания».
const initial: Theme =
  window.localStorage.getItem(THEME_STORAGE_KEY) === 'dark' ? 'dark' : 'light'
applyTheme(initial)

export const useThemeStore = create<ThemeState>((set, get) => ({
  theme: initial,
  setTheme: (theme) => {
    applyTheme(theme)
    window.localStorage.setItem(THEME_STORAGE_KEY, theme)
    set({ theme })
  },
  toggleTheme: () => {
    const next: Theme = get().theme === 'dark' ? 'light' : 'dark'
    applyTheme(next)
    window.localStorage.setItem(THEME_STORAGE_KEY, next)
    set({ theme: next })
  },
}))