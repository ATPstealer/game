<template>
  <MessageBlock
    v-if="messageData?.code"
    class="mb-4"
    v-bind="messageData"
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
import { usePostUserLogin, type UserPayload, type JsonResult } from '@/gen'

const emits = defineEmits<{
  (e: 'close'): void;
  (e: 'sign-up'): void;
}>()

const nickName = ref<string>('')
const pass = ref<string>('')
const messageData = ref<JsonResult>()

const { t } = useI18n()
const { setToken } = useUser()

const login = usePostUserLogin({
  mutation: {
    onSuccess: data => {
      messageData.value = data

      if (data && data?.code <= 0) {
        setToken(data as JsonResult & {data: {ttl: string; token: string}})
        setTimeout(() => {
          emits('close')
        }, 1000)
      }
    }
  }
})

const submit = () => {
  const userData = {
    nickName: nickName.value,
    password: sha256(pass.value).toString()
  } as UserPayload

  login.mutate({ data: { ...userData } })
}
</script>

<style scoped>

</style>
