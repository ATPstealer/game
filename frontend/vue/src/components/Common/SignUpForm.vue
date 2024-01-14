<template>
  <MessageBlock
    v-if="message"
    :message="message"
    class="mb-4"
  />
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1">
      <label
        for="nickname"
        class="block text-gray-700 text-sm font-medium"
      >
        nickName
      </label>
      <InputText
        v-model="nickName"
        class="w-full"
        input-id="nickname"
        aria-describedby="nickname-help"
        :class="{'p-invalid': !nickNameValid}"
      />
      <small
        v-if="!nickNameValid"
        id="nickname-help"
        class="text-red-500"
      >
        Fill in the field nickname
      </small>
    </div>
    <div class="flex flex-col gap-1">
      <label
        for="email"
        class="block text-gray-700 text-sm font-medium"
      >
        Email
      </label>
      <InputText
        v-model="email"
        class="w-full"
        input-id="email"
        :class="{'p-invalid': !emailValid}"
        aria-describedby="email-help"
      />
      <small
        v-if="!emailValid"
        id="email-help"
        class="text-red-500"
      >
        Enter correct email address
      </small>
    </div>
    <div class="flex flex-col gap-1">
      <label
        for="password"
        class="block text-gray-700 text-sm font-medium"
      >
        Password
      </label>
      <Password
        v-model="pass1"
        input-class="w-full"
        class="w-full"
        :feedback="false"
        input-id="password"
        aria-describedby="pass-help"
        :class="{'p-invalid': !strongPassword}"
      />
      <small
        v-if="!strongPassword"
        id="pass-help"
        class="text-red-500"
      >
        Password must include lowercase and uppercase later, digits and be more than 8 symbols
      </small>
    </div>
    <div class="flex flex-col gap-1">
      <label
        for="confirm"
        class="block text-gray-700 text-sm font-medium"
      >
        Confirm Password
      </label>
      <Password
        v-model="pass2"
        input-class="w-full"
        class="w-full"
        :feedback="false"
        input-id="confirm"
        aria-describedby="confirm-help"
        :class="{'p-invalid': !matchPasswords}"
      />
      <small
        v-if="!matchPasswords"
        id="confirm-help"
        class="text-red-500"
      >
        Passwords doesn't match
      </small>
    </div>
    <div class="flex gap-8 items-center justify-between">
      <Button
        type="submit"
        severity="info"
        label="SignUp"
        @click="submit"
        :disabled="!validForm"
      />
      <Button
        severity="info"
        text
        label="Already have an account? Log in"
        @click="emits('log-in')"
        class="text-xs md:text-sm"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import sha256 from 'crypto-js/sha256'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import { computed, ref } from 'vue'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useUser } from '@/composables/useUser'
import type { DataMessage } from '@/types'

const emits = defineEmits<{
  (e: 'close'): void;
  (e: 'log-in'): void;
}>()

const nickName = ref<string>('')
const email = ref<string>('')
const pass1 = ref<string>('')
const pass2 = ref<string>('')
const message = ref<DataMessage | null>(null)

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

const { signUp } = useUser()
const submit = () => {
  message.value = null
  const userData = {
    nickName: nickName.value,
    email: email.value,
    password: sha256(pass1.value).toString()
  }
  const { dataMessage, onFetchFinally } = signUp(userData)

  onFetchFinally(() => {
    message.value = dataMessage.value

    if (message.value && message.value?.status === 'success') {
      setTimeout(() => {
        emits('log-in')
      }, 1000)
    }
  })
}
</script>

<style scoped>

</style>
