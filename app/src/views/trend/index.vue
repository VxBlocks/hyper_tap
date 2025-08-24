<template>
  <div class="h-screen flex flex-col overflow-hidden">
    <!-- 固定头部 -->
    <div class="sticky-header">
      <!-- Tab Navigation -->
      <div class="relative bg-white px-4 pt-3 pb-1 border-b border-gray-200">
        <div class="flex space-x-6 justify-center">
          <button v-for="tab in tabs" :key="tab.key" @click="activeTab = tab.key" :class="[
            'text-sm font-medium',
            activeTab === tab.key && 'text-[#000] !font-bold',
          ]">
            {{ $t(tab.translationKey) }}
          </button>
        </div>
        <div class="absolute right-[15px] top-[14px]">
          <img src="/search.png" alt="search" @click="clickToSearch" />
        </div>
      </div>
    </div>
    <van-tabs v-model:active="activeTab1" class="w-full text-xl" swipeable>
      <van-tab v-for="tab in homeTabList" :name="tab.key" :title="tab.label" :key="tab.key">
        <div class="scroll-container">
          <TokenTable :activeTab="activeTab1" />
        </div>
      </van-tab>
    </van-tabs>
    <SearchModal v-model:show="showSearchModal" />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onMounted } from "vue";
import SearchModal from "@/components/modals/search-mark-modal.vue";
import { useI18n } from 'vue-i18n';
import { useHomeTabs } from "@/enums/homeTabEnum";

const { HomeTabEnum, homeTabList } = useHomeTabs();

const { t } = useI18n();
const showSearchModal = ref(false);
const refreshing = ref(false);

const activeTab = ref("contract");
const tabs = [
  { key: "contract", translationKey: "tabs.contract" },
  { key: "spot", translationKey: "tabs.spot" }
];


const activeTab1 = ref(HomeTabEnum.Gainers);
const clickToSearch = () => {
  showSearchModal.value = true;
};
async function refresh_data() { }
const onRefresh = () => {
  refreshing.value = true;
  refresh_data().finally(() => {
    refreshing.value = false;
  });
};

onMounted(async () => {
  activeTab.value = "contract";
  // 初始加载数据
  refreshing.value = true;
  await refresh_data();
  refreshing.value = false;
});
</script>

<style scoped>
.sticky-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: white;
}

.scroll-container {
  height: calc(100vh - 188px);
  /* 根据实际头部高度调整 */
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}
</style>