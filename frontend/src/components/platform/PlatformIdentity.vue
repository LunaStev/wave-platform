<script setup lang="ts">
defineProps<{ service: string }>()
</script>

<template>
  <RouterLink class="platform-identity" to="/" aria-label="Wave Platform">
    <img src="/img/wave.svg" alt="Wave" />
    <span>Wave</span>
    <small v-if="service">{{ service }}</small>
  </RouterLink>
</template>
