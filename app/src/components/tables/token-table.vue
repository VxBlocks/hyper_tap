<script setup lang="ts">
import { HomeTabEnum } from "@/enums/homeTabEnum";
import { computed, onMounted, ref, watch } from "vue";
import { useHomeStore } from "@/store/modules/home";
import { current_network } from "@/constant";
import TokenInfo from "@/components/token-info.vue";

const homeStore = useHomeStore();
const tokenList = computed(() => homeStore.getTokenList);

const props = defineProps<{
  activeTab: HomeTabEnum;
}>();
async function onRefreshTokenList() {
  await homeStore.queryTokenListByType(props.activeTab, current_network);
}

async function initTokensData() {
  await homeStore.queryTokenListByType(props.activeTab, current_network);
}
onMounted(async () => {
  await initTokensData();
});

watch(
  () => props.activeTab,
  async () => {
    await onRefreshTokenList();
  }
);
</script>

<template>
  <div>
    <!-- 表头 -->
    <div class="crypto-table-header sticky-header">
      <div class="describe px-4">
        <div class="gap-[8px] flex">
          <div>{{ $t("token-tabs.pair") }} / {{ $t("token-tabs.day_volume") }}</div>
        </div>
        <div>
          <span>{{ $t("token-tabs.lprice") }} / {{ $t("token-tabs.change") }}</span>
        </div>
      </div>
    </div>
    <div class="px-4">
      <TokenInfo v-for="(item, index) in tokenList" :key="index" :data="item" />
    </div>
  </div>
</template>

<style scoped>
.sticky-header {
  position: sticky;
  top: 0;
  z-index: 10;
  background-color: white;
  flex-shrink: 0;
}

.describe {
  color: #78817e;
  font-size: 12px;
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
}
</style>
