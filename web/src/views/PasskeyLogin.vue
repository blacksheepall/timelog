<template>
  <div class="min-h-[70vh] flex items-center justify-center px-4">
    <div class="w-full max-w-5xl grid lg:grid-cols-[0.9fr_1.1fr] gap-8">
      <aside
        class="rounded-3xl border border-default bg-bg-surface p-8 shadow-md-[0_18px_45px_rgba(15,23,42,0.08)]"
      >
        <p class="text-xs uppercase tracking-[0.2em] text-text-secondary">Welcome back</p>
        <h2 class="mt-3 text-3xl font-semibold text-text-primary">用 Passkey 登录</h2>
        <p class="mt-2 text-text-secondary">无密码、更快速的身份验证。</p>
        <div class="mt-8 space-y-4 text-sm text-text-secondary">
          <div class="flex items-start gap-3">
            <ShieldCheckIcon class="mt-0.5 h-5 w-5 text-emerald-500" />
            <p>系统会提示你使用设备上的生物识别或安全密钥。</p>
          </div>
          <div class="flex items-start gap-3">
            <ShieldCheckIcon class="mt-0.5 h-5 w-5 text-emerald-500" />
            <p>验证成功后将自动进入仪表盘。</p>
          </div>
        </div>
        <div
          class="mt-8 rounded-2xl border border-default bg-bg-elevated p-4 text-xs text-text-secondary"
        >
          WebAuthn 需要 HTTPS 或 localhost 环境。
        </div>
      </aside>

      <section
        class="rounded-3xl bg-gradient-to-br from-slate-900 via-slate-800 to-slate-700 p-10 text-white"
      >
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs uppercase tracking-[0.2em] text-text-secondary">Passkey Login</p>
            <h3 class="mt-4 text-2xl font-semibold">开始验证</h3>
          </div>
          <div class="hidden md:flex h-14 w-14 items-center justify-center rounded-2xl bg-bg-surface/10">
            <LockClosedIcon class="h-7 w-7" />
          </div>
        </div>

        <button
          class="mt-10 inline-flex w-full items-center justify-center rounded-xl bg-bg-surface px-5 py-3 text-sm font-semibold text-text-primary shadow-lg shadow-black/20 transition hover:bg-bg-elevated disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="loading"
          @click="handleLogin"
        >
          <span v-if="loading">验证中...</span>
          <span v-else>使用 Passkey 登录</span>
        </button>

        <p
          v-if="error"
          class="mt-6 rounded-xl border border-white/20 bg-bg-surface/10 px-4 py-3 text-sm"
        >
          {{ error }}
        </p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { LockClosedIcon, ShieldCheckIcon } from '@heroicons/vue/24/outline'
  import { passkeyAPI } from '@/api'
  import { beginLogin, isWebAuthnSupported } from '@/utils/webauthn'
  import { setAuthToken } from '@/utils/auth'

  const router = useRouter()
  const loading = ref(false)
  const error = ref('')

  const handleLogin = async () => {
    error.value = ''
    if (!isWebAuthnSupported()) {
      error.value = '当前浏览器不支持 Passkey/WebAuthn'
      return
    }

    loading.value = true
    try {
      const beginResponse = await passkeyAPI.loginBegin()
      const { session_id, data } = beginResponse.data
      const assertion = await beginLogin(data)
      const finishResponse = await passkeyAPI.loginFinish(session_id, assertion)
      setAuthToken(finishResponse.data.token)
      router.push('/')
    } catch (err: any) {
      error.value = err?.response?.data?.message || err?.message || '登录失败'
    } finally {
      loading.value = false
    }
  }
</script>
