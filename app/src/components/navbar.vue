<template>
  <van-nav-bar class="my-navbar" left-arrow @click-left="router.back">
    <template #title>
      {{ getTitle }}
    </template>
    <template #left>
      <van-icon name="arrow-left" size="18px" color="#000000" />
    </template>
    <template #right>
      <slot name="right" />
    </template>
  </van-nav-bar>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";

const router = useRouter();
const currentRoute = useRoute();

const getTitle = computed(() => currentRoute.meta.title as string);
</script>

<style scoped lang="less">
.my-navbar {
  :deep(.van-nav-bar__title) {
    color: #000 !important;
  }
}
</style>
