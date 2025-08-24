
<script setup lang="ts">
import { useRouteStore } from "@/store/modules/route";
import { computed, ComputedRef } from "vue";
import { RouteRecordRaw, useRoute } from "vue-router";

const routeStore = useRouteStore()
// 需要缓存的路由组件
const keepAliveComponents = computed(() => routeStore.keepAliveComponents)
const currentRoute = useRoute()

const getTitle = computed(() => currentRoute.meta.title as string)

 

const getShowHeader = computed(() => !currentRoute.meta.hiddenHeader)

</script>

<template>
  <!-- <van-tabbar route class="tabbar">
      <van-tabbar-item
        v-for="menu in getMenus"
        :key="menu.name"
        replace
        :to="menu.path"
      >
        <template #icon>
          <i :class="menu.meta?.icon" />
        </template>
        {{ menu.meta?.title }}
      </van-tabbar-item>
    </van-tabbar> -->
</template>


<style scoped lang="less">
.my-navbar {
  background-color: #503cfe;
  :deep(.van-nav-bar__title) {
    color: #fff !important;
  }
}
</style>
