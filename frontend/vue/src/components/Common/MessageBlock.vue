<template>
  <Message
    :pt="{
      wrapper: {
        class: '!p-2'
      }
    }"
    :severity="code <= 0 ? 'success' : 'error'"
  >
    {{ message }}
  </Message>
</template>

<script setup lang="ts">
import isString from 'lodash/isString'
import Message from 'primevue/message'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

interface Props {
  code: number;
  values?: Record<string, string | number>;
  text?: string;
  status?: string;
}

const props = defineProps<Props>()
const { t, locale } = useI18n()

const codeText = ref<string>(t(`codes.${props.code.toString()}`))

watch(locale, () => {
  codeText.value = t(`codes.${props.code.toString()}`)
})

watch(() => props.code, () => {
  codeText.value = t(`codes.${props.code.toString()}`)
})

const message = computed(() => {
  if (props.values) {
    // Перестраховка, если приходит string
    let values = {} as Record<string, string | number>
    if (isString(props.values)) {
      values = JSON.parse(props.values as string)
    } else {
      values = props.values
    }

    return codeText.value.replace(/%\w+/g, match => {
      const key = match.slice(1)
      const value = values?.[key]

      return value.toString() || match
    })
  }

  return codeText.value
})
</script>

<style scoped>

</style>