<template>
  <div class="h-screen flex flex-col">
    <div class="sticky-header">
      <BalanceCard />
      <div class="bg-white py-1">
        <van-tabs v-model:active="activeTab" class="w-full" swipeable>
          <van-tab v-for="tab in accountTabList" :name="tab.key" :title="tab.label" />
        </van-tabs>
      </div>
    </div>
    <!-- 优化后的动态组件渲染 -->

    <div class="flex-1 overflow-hidden bg-white px-4">
      <div class="h-full overflow-y-auto ">
        <component :is="tabComponentMap[activeTab]" :key="activeTab" />
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { onMounted, ref, watchEffect } from "vue";
import BalanceCard from "./components/balance-card.vue";
import SpotContent from "./components/spot/index.vue";
import PositionsContent1 from "./components/position/index.vue";
// import PendingOrdersContent from "./components/fopen-orders-content.vue";

import PendingOrdersContent from "./components/openorder/index.vue";
import ExecutedOrdersContent from "./components/fills-connect.vue";
import { useAccountTabs } from "@/enums/accountTabEnum";
import { current_network } from "@/constant";
import { useHomeStore } from "@/store/modules/home";
import { useAccountStore } from "@/store/modules/account";
import { useAppKitAccount } from "@reown/appkit/vue";
import { useI18n } from 'vue-i18n';
const { t } = useI18n();
const { accountTabList, AccountTabEnum } = useAccountTabs();
const activeTab = ref(AccountTabEnum.Positions);
const homeStore = useHomeStore()
const accountStore = useAccountStore()
const accountData = useAppKitAccount()

// 定义 tab 到组件的映射关系
const tabComponentMap = {
  [AccountTabEnum.Spot]: SpotContent,
  [AccountTabEnum.Positions]: PositionsContent1,
  [AccountTabEnum.PendingOrders]: PendingOrdersContent,
  [AccountTabEnum.ExecutedOrders]: ExecutedOrdersContent, 
};

watchEffect(() => {
  if (accountData.value && accountData.value.isConnected && accountData.value.address) {
    accountStore.initSubscribeClient(current_network, accountData.value.address)
  }
})

onMounted(() => {
  activeTab.value = AccountTabEnum.Positions;
  homeStore.initTokens(current_network)
});

</script>

<style scoped lang="less">
/* 固定的钱包卡片容器 */
.sticky-container {
  position: sticky;
  top: 0;
  z-index: 100;
  background: white;
}
</style>
