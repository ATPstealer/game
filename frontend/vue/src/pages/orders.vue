<template>
  <Layout>
    <div v-if="!isFetching">
      <h2 class="font-bold text-2xl">
        Orders
      </h2>
      <div class="grid grid-cols-5 gap-2 font-bold">
        <h3 class="col-span-1">
          Cell
        </h3>
        <h3 class="col-span-1">
          Resource
        </h3>
        <h3 class="col-span-1">
          Amount
        </h3>
        <h3 class="col-span-1">
          Price
        </h3>
        <h3 class="col-span-1">
          Type (Close)
        </h3>
      </div>
      <div
        v-for="order in orders"
        :key="order.id"
        class="grid grid-cols-5 "
      >
        <p class="col-span-1">
          {{ order.x }}:{{ order.y }}
        </p>
        <p class="col-span-1">
          {{ order.resourceName }}
        </p>
        <p class="col-span-1">
          {{ order.amount }}
        </p>
        <p class="col-span-1">
          {{ order.priceForUnit }}$
        </p>
        <Button
          text
          :label="order.sell ? 'sell' : 'buy'"
          @click="closeOrder(order.id)"
          class="text-red-600 font-bold col-span-1"
        />
      </div>
    </div>
    <Loading v-else />
  </Layout>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import Layout from '@/components/Common/Layout.vue'
import Loading from '@/components/Common/Loading.vue'
import { useGetData } from '@/composables/useGetData'
import { useOrders } from '@/composables/useOrders'
import type { Order } from '@/types'

const { data: orders, isFetching } = useGetData<Order[]>('/market/order/my')
const { closeOrder } = useOrders()
</script>

<style scoped>

</style>
