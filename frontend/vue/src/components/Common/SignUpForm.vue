<template>
  <MessageBlock
    v-if="messageData?.code"
    v-bind="messageData"
    class="mb-4"
  />
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1">
      <label
        class="block text-gray-700 text-sm font-medium"
        for="nickname"
      >
        {{ t('account.nickname') }}
      </label>
      <InputText
        v-model="nickName"
        aria-describedby="nickname-help"
        class="w-full"
        :class="{'p-invalid': !nickNameValid}"
        input-id="nickname"
      />
      <small
        v-if="!nickNameValid"
        id="nickname-help"
        class="text-red-500"
      >
        {{ t('account.errors.nickname') }}
      </small>
    </div>
    <div class="flex flex-col gap-1">
      <label
        class="block text-gray-700 text-sm font-medium"
        for="email"
      >
        {{ t('account.email') }}
      </label>
      <InputText
        v-model="email"
        aria-describedby="email-help"
        class="w-full"
        :class="{'p-invalid': !emailValid}"
        input-id="email"
      />
      <small
        v-if="!emailValid"
        id="email-help"
        class="text-red-500"
      >
        {{ t('account.errors.email') }}
      </small>
    </div>
    <div class="flex flex-col gap-1">
      <label
        class="block text-gray-700 text-sm font-medium"
        for="password"
      >
        {{ t('account.password') }}
      </label>
      <Password
        v-model="pass1"
        aria-describedby="pass-help"
        class="w-full"
        :class="{'p-invalid': !strongPassword}"
        :feedback="false"
        input-class="w-full"
        input-id="password"
      />
      <small
        v-if="!strongPassword"
        id="pass-help"
        class="text-red-500"
      >
        {{ t('account.errors.password') }}
      </small>
    </div>
    <div class="flex flex-col gap-1">
      <label
        class="block text-gray-700 text-sm font-medium"
        for="confirm"
      >
        {{ t('account.confirmPass') }}
      </label>
      <Password
        v-model="pass2"
        aria-describedby="confirm-help"
        class="w-full"
        :class="{'p-invalid': !matchPasswords}"
        :feedback="false"
        input-class="w-full"
        input-id="confirm"
      />
      <small
        v-if="!matchPasswords"
        id="confirm-help"
        class="text-red-500"
      >
        {{ t('account.errors.passMatch') }}
      </small>
    </div>
    <div class="flex gap-8 items-center justify-between">
      <Button
        :disabled="!validForm"
        :label="t('account.signup')"
        severity="info"
        type="submit"
        @click="submit"
      />
      <Button
        class="text-xs md:text-sm max-w-[50%]"
        :label="t('account.haveAccount')"
        severity="info"
        text
        @click="emits('log-in')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useMutation } from '@tanstack/vue-query'
import sha256 from 'crypto-js/sha256'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { JsonResult } from '@/api'
import { postUserCreateMutation } from '@/api/@tanstack/vue-query.gen'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import type { BackData } from '@/types'

const emits = defineEmits<{
  (e: 'close'): void;
  (e: 'log-in'): void;
}>()

const nickName = ref<string>('nomel')
const email = ref<string>('nomelnomel@gmail.com')
const pass1 = ref<string>('qweASD123')
const pass2 = ref<string>('qweASD123')
const messageData = ref<JsonResult>()

const { t } = useI18n()

const nickNameValid = computed(() => !nickName.value || nickName.value.length > 3)

const emailValid = computed(() =>
  !email.value ||
    email.value.toLowerCase().match(
      /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
    )
)

const strongPassword = computed(() => {
  if (!pass1.value) {
    return true
  }

  if (pass1.value.length < 8) {
    return false
  }

  const uppercaseRegex = /[A-Z]/
  const lowercaseRegex = /[a-z]/
  if (!uppercaseRegex.test(pass1.value) || !lowercaseRegex.test(pass1.value)) {
    return false
  }

  const digitRegex = /\d/

  return digitRegex.test(pass1.value)
})

const matchPasswords = computed(() => {
  if (!pass1.value || !pass2.value) {
    return true
  }

  return pass1.value === pass2.value
})

const validForm = computed(() => {
  return nickName.value && nickNameValid.value &&
      email.value && emailValid.value &&
      pass1.value && strongPassword.value &&
      pass2.value && matchPasswords.value
})

const signUp = useMutation({
  ...postUserCreateMutation(),
  onSuccess: (data: JsonResult) => {
    messageData.value = data

    if (data && data.code === -1) {
      setTimeout(() => {
        emits('log-in')
      }, 1000)
    }
  }
})
const submit = () => {
  messageData.value = {} as BackData
  const userData = {
    nickName: nickName.value,
    email: email.value,
    password: sha256(pass1.value).toString()
  }

  signUp.mutate({ body: { ...userData } })
}
</script>

<style scoped>

</style>
