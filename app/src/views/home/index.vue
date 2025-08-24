<template>
  <div class="h-screen flex flex-col overflow-hidden">
    <!-- 固定头部 -->
    <div class="sticky-header">
      <HomeHeaderCard />
      <div class="mailbox px-4">
        <van-badge :content="count" :show-zero="false">
          <img src="/mailbox.svg" alt="star" class="token_icon" @click="jumpToMessage" />
        </van-badge>
      </div>
    </div>

    <van-tabs v-model:active="activeTab" class="w-full text-xl" swipeable>
      <van-tab v-for="tab in tabs" :name="tab.key" :title="tab.label" :key="tab.key">
        <div class="scroll-container">
          <TokenTable :activeTab="activeTab" />
        </div>
      </van-tab>
    </van-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import HomeHeaderCard from "./components/home-header-card.vue";
import TokenTable from "@/components/tables/token-table.vue";
import { HomeTabEnum, useHomeTabs } from "@/enums/homeTabEnum";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";

const router = useRouter();
const { t } = useI18n();
const { homeTabList } = useHomeTabs();
const activeTab = ref(HomeTabEnum.Hot);
const tabs = homeTabList;
const count = ref(3);
function jumpToMessage() {
  router.push("/message_center");
}
</script>

<style scoped>
.sticky-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: white;
}

.mailbox {
  display: flex;
  justify-content: flex-end; 
  padding-top: 8px;
}

.scroll-container {
  height: calc(100vh - 244px);
  /* 根据实际头部高度调整 */
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}
</style>
