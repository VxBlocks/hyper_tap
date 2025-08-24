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

    <!-- 标签栏 -->
    <van-tabbar v-model="active" route class="tabbar" active-color="#46CCB9">
      <van-tabbar-item replace to="/home" icon="home-o">{{ t("menu.home") }}</van-tabbar-item>
      <van-tabbar-item replace to="/trend" icon="chart-trending-o">{{ t("menu.markets") }}</van-tabbar-item>
      <van-tabbar-item replace to="/sale" icon="after-sale">{{ t("menu.trade") }}</van-tabbar-item>
      <van-tabbar-item replace to="/account" icon="user-circle-o">{{ t("menu.my") }}</van-tabbar-item>
    </van-tabbar>
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
  // padding-top: env(safe-area-inset-top);
  box-sizing: border-box;
}

.app-nav-bar {

  // padding-top: env(safe-area-inset-top);
  // :deep(.van-nav-bar) {
  //   padding-top: env(safe-area-inset-top);
  // }
}

.app-content {
  flex: 1;
  overflow-x: hidden;
} 
</style>