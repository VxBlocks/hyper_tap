<template>
  <div class="app-container">
    <!-- 核心修改：包裹 router-view 的容器 -->
    <routerView class="app-content">
      <template #default="{ Component, route }">
        <keep-alive v-if="keepAliveComponents" :include="keepAliveComponents">
          <component :is="Component" :key="route.fullPath" />
        </keep-alive>
        <component :is="Component" v-else :key="route.fullPath" />
      </template>
    </routerView>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRouteStore } from '@/store/modules/route';
import { useI18n } from 'vue-i18n';
const { t } = useI18n()

const active = ref(0);
const routeStore = useRouteStore();
const keepAliveComponents = computed(() => routeStore.keepAliveComponents);

</script>

<style scoped lang="less">
.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  padding-top: env(safe-area-inset-top);
  box-sizing: border-box;
}  
.app-content {
  flex: 1;
  overflow-x: hidden;
} 
</style>