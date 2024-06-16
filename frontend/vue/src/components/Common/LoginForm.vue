<template>
  <MessageBlock
    v-if="messageData?.code"
    v-bind="messageData"
    class="mb-4"
  />
  <div class="flex flex-col gap-8">
    <div class="p-float-label mt-4">
      <InputText
        v-model="nickName"
        class="w-full"
        input-id="nickname"
      />
      <label
        class="block text-gray-700 text-sm font-medium"
        for="nickname"
      >
        {{ t('account.nickname') }}
      </label>
    </div>
    <div class="p-float-label">
      <Password
        v-model="pass"
        class="w-full"
        :feedback="false"
        input-class="w-full"
        input-id="password"
        @keydown.enter="submit"
      />
      <label
        class="block text-gray-700 text-sm font-medium"
        for="password"
      >
        {{ t('account.password') }}
      </label>
    </div>
    <div class="flex items-center justify-between">
      <Button
        :label="t('account.login')"
        severity="info"
        type="submit"
        @click="submit"
      />
      <Button
        :label="t('account.signup')"
        severity="info"
        text
        @click="emits('sign-up')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import sha256 from 'crypto-js/sha256'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useUser } from '@/composables/useUser'
import type { BackData } from '@/types'

const emits = defineEmits<{
  (e: 'close'): void;
  (e: 'sign-up'): void;
}>()

const nickName = ref<string>('')
const pass = ref<string>('')
const messageData = ref<BackData>()

const { t } = useI18n()
const { logIn } = useUser()
const submit = () => {
  const userData = {
    nickName: nickName.value,
    password: sha256(pass.value).toString()
  }
  const { data, onFetchFinally } = logIn(userData)

  onFetchFinally(() => {
    messageData.value = data.value

    if (data.value?.code <= 0) {
      setTimeout(() => {
        emits('close')
      }, 1000)
    }
  })
}
</script>

<style scoped>

</style>
