import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { getAuthToken, setAuthToken } from '@/utils/auth'
import './style.css'

async function bootstrap() {
  if (import.meta.env.DEV && !getAuthToken()) {
    try {
      const res = await fetch('/api/auth/dev-login', { method: 'POST' })
      const json = await res.json()
      const token = json?.data?.token
      if (token) {
        setAuthToken(token)
      }
    } catch (err) {
      console.warn('Dev auto-login failed:', err)
    }
  }

  const app = createApp(App)
  const pinia = createPinia()

  app.use(pinia)
  app.use(router)
  app.mount('#app')
}

bootstrap()
