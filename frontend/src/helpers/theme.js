import { ref } from 'vue'

const COOKIE_NAME = 'orbit_theme'
const COOKIE_MAX_AGE = 60 * 60 * 24 * 365
const DEFAULT_THEME = 'default'

export const THEMES = [
  {
    id: 'default',
    label: 'Default',
    description: 'Deep night with violet and mint signals',
    preview: { background: '#0b0d17', surface: '#1f2440', text: '#f4f5fa', primary: '#7c5cff', secondary: '#3ee6b0' },
  },
  {
    id: 'pink',
    label: 'Pink',
    description: 'Soft blush surfaces with hot pink accents',
    preview: { background: '#fff0f6', surface: '#ffffff', text: '#3b0f2a', primary: '#b5179e', secondary: '#d61a7a' },
  },
  {
    id: 'light',
    label: 'Light',
    description: 'Clean daylight with crisp indigo accents',
    preview: { background: '#f3f5fb', surface: '#ffffff', text: '#131726', primary: '#5b3fe0', secondary: '#0d9b74' },
  },
  {
    id: 'ocean',
    label: 'Ocean',
    description: 'Deep sea blues with glowing aqua highlights',
    preview: { background: '#041a24', surface: '#0d3a4f', text: '#eaf8fc', primary: '#2f7bff', secondary: '#19e0c8' },
  },
  {
    id: 'ember',
    label: 'Ember',
    description: 'Warm charcoal with burning orange and gold',
    preview: { background: '#150d0a', surface: '#33201a', text: '#fdf0e6', primary: '#f2600c', secondary: '#ffb020' },
  },
]

export const currentTheme = ref(DEFAULT_THEME)

function readCookie(name) {
  const entry = document.cookie.split('; ').find((row) => row.startsWith(`${name}=`))
  if (!entry) return null

  try {
    return decodeURIComponent(entry.slice(name.length + 1))
  } catch {
    return null
  }
}

function writeCookie(name, value) {
  const secure = window.location.protocol === 'https:' ? '; Secure' : ''
  document.cookie = `${name}=${encodeURIComponent(value)}; path=/; max-age=${COOKIE_MAX_AGE}; SameSite=Lax${secure}`
}

export function isValidTheme(id) {
  return THEMES.some((theme) => theme.id === id)
}

export function getSavedTheme() {
  const saved = readCookie(COOKIE_NAME)
  return isValidTheme(saved) ? saved : DEFAULT_THEME
}

export function applyTheme(id) {
  const theme = isValidTheme(id) ? id : DEFAULT_THEME
  document.documentElement.setAttribute('data-theme', theme)
  currentTheme.value = theme
  return theme
}

export function setTheme(id) {
  const theme = applyTheme(id)
  writeCookie(COOKIE_NAME, theme)
  return theme
}

export function initTheme() {
  return applyTheme(getSavedTheme())
}
