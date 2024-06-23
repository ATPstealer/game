<template>
  {{ props.building.equipmentEffect }}
</template>


<script setup lang="ts">
import type {Building} from "@/types/Buildings/index.interface";

interface Props {
  building: Building;
}

const props = defineProps<Props>()

</script>