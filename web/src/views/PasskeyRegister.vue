<template>
  <div class="min-h-[70vh] flex items-center justify-center px-4">
    <div class="w-full max-w-4xl grid lg:grid-cols-[1.1fr_0.9fr] gap-8">
      <section
        class="rounded-3xl border border-default bg-bg-surface/80 shadow-md-[0_20px_60px_rgba(15,23,42,0.08)] backdrop-blur"
      >
        <div class="p-8">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm uppercase tracking-[0.2em] text-text-secondary">Passkey Setup</p>
              <h2 class="mt-3 text-3xl font-semibold text-text-primary">绑定你的设备</h2>
              <p class="mt-2 text-text-secondary">输入临时密码，完成一次安全的 Passkey 绑定。</p>
            </div>
            <div
              class="hidden md:flex h-14 w-14 items-center justify-center rounded-2xl bg-text-primary text-white"
            >
              <KeyIcon class="h-7 w-7" />
            </div>
          </div>

          <div class="mt-8 grid gap-6">
            <div>
              <label class="block text-sm font-medium text-text-muted" for="temp-password">
                临时密码
              </label>
              <input
                id="temp-password"
                v-model="tempPassword"
                type="password"
                autocomplete="off"
                autocapitalize="off"
                spellcheck="false"
                class="mt-2 w-full rounded-xl border border-default bg-bg-surface px-4 py-3 text-text-primary shadow-sm focus:border-default focus:outline-none focus:ring-2 focus:ring-default"
                placeholder="从命令行获取的临时密码"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-text-muted" for="device-name">
                设备名称
              </label>
              <input
                id="device-name"
                v-model="deviceName"
                type="text"
                class="mt-2 w-full rounded-xl border border-default bg-bg-surface px-4 py-3 text-text-primary shadow-sm focus:border-default focus:outline-none focus:ring-2 focus:ring-default"
                placeholder="例如：Paul 的 MacBook"
              />
            </div>
            <button
              class="mt-2 inline-flex items-center justify-center rounded-xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white shadow-lg shadow-black/20 transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-60"
              :disabled="loading || !tempPassword"
              @click="handleRegister"
            >
              <span v-if="loading">正在绑定...</span>
              <span v-else>开始绑定</span>
            </button>
            <p
              v-if="error"
              class="rounded-xl border border-danger-border bg-danger-bg px-4 py-3 text-sm text-red-700"
            >
              {{ error }}
            </p>
          </div>
        </div>
      </section>

      <aside
        class="rounded-3xl bg-gradient-to-br from-slate-900 via-slate-800 to-slate-700 p-8 text-white"
      >
        <div class="flex items-center gap-3">
          <SparklesIcon class="h-5 w-5" />
          <p class="text-xs uppercase tracking-[0.2em] text-text-secondary">What happens</p>
        </div>
        <h3 class="mt-5 text-2xl font-semibold">一次绑定，长期免密</h3>
        <ul class="mt-6 space-y-4 text-sm text-text-secondary">
          <li>临时密码用于验证初始化权限，使用后立即失效。</li>
          <li>浏览器将提示你创建 Passkey，并保存到设备安全芯片。</li>
          <li>完成后即可使用 Passkey 登录，无需再输入密码。</li>
        </ul>
        <div class="mt-10 rounded-2xl border border-white/10 bg-bg-surface/5 p-4">
          <p class="text-sm text-text-secondary">需要在 HTTPS 或 localhost 环境中操作 WebAuthn。</p>
          <router-link
            to="/login"
            class="mt-3 inline-flex text-sm font-semibold text-white underline underline-offset-4"
          >
            已有 Passkey？去登录
          </router-link>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { KeyIcon, SparklesIcon } from '@heroicons/vue/24/outline'
  import { passkeyAPI } from '@/api'
  import { beginRegistration, isWebAuthnSupported } from '@/utils/webauthn'

  const router = useRouter()

  const tempPassword = ref('')
  const deviceName = ref('')
  const loading = ref(false)
  const error = ref('')

  const handleRegister = async () => {
    error.value = ''
    if (!isWebAuthnSupported()) {
      error.value = '当前浏览器不支持 Passkey/WebAuthn'
      return
    }

    loading.value = true
    try {
      console.log('Step 1: Calling registerBegin...')
      const beginResponse = await passkeyAPI.registerBegin(tempPassword.value, deviceName.value)
      console.log('Step 2: registerBegin response:', beginResponse)

      const { session_id, data } = beginResponse.data
      console.log('Step 3: Extracted session_id and data')

      console.log('Step 4: Calling beginRegistration with data:', data)
      const credential = await beginRegistration(data)
      console.log('Step 5: beginRegistration result:', credential)

      await passkeyAPI.registerFinish(session_id, credential, deviceName.value)
      console.log('Step 6: registerFinish success')

      router.push('/login')
    } catch (err: any) {
      console.error('Registration error:', err)
      error.value = err?.response?.data?.message || err?.message || '绑定失败'
    } finally {
      loading.value = false
    }
  }
</script>
