<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Tab Navigation -->
    <div class="bg-white px-4 py-3">
      <div class="flex justify-between">
        <div class="flex space-x-6">
          <button
            v-for="tab in tabs"
            :key="tab"
            @click="activeTab = tab"
            :class="[
              'text-sm font-sm',
              activeTab === tab && 'text-[#000] !font-bold',
            ]"
          >
            {{ tab }}
          </button>
        </div>
        <button class="add_wallet" @click="handleAdd">+ Add Wallet</button>
      </div>
    </div>
    <!-- Pull to Refresh -->
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <!-- Crypto List -->
      <div>
        <!-- Tab内容切换动画 -->
        <transition name="fade" mode="out-in">
          <div :key="activeTab" class="transition-container">
            <!-- 使用 van-cell-group 和 van-cell 实现表格效果 -->
            <div v-if="activeTab === 'Messages'">
              <WalletMessageContent />
            </div>
            <div v-else-if="activeTab === 'List'">
              <WalletListContent />
            </div>
          </div>
        </transition>
      </div>
    </van-pull-refresh>
    <AddWalletModal v-model:show="showAddWalletModal" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { onMounted } from "vue";
import NavBar from "@/components/navbar.vue";

import { useAccountStore } from "@/store/modules/account";
import { useAppKitAccount } from "@reown/appkit/vue";
import { current_network } from "@/constant";
import { useRouter } from "vue-router";
import AddWalletModal from "@/components/modals/add-wallet-modal.vue";
import WalletListContent from "./walletlist.vue";
import WalletMessageContent from "./wallet-message.vue";

const showAddWalletModal = ref(false);
const router = useRouter();

const activeTab = ref("Messages");
const tabs = ["Messages", "List"];

const refreshing = ref(false);

const activeNames = ref([]);
const handleAdd = () => {
  showAddWalletModal.value = true;
};
async function refresh_data() {}
// 下拉刷新处理函数
const onRefresh = () => {
  refreshing.value = true;
  refresh_data().finally(() => {
    refreshing.value = false;
  });
};

onMounted(async () => {
  activeTab.value = "Messages";
  // 初始加载数据
  refreshing.value = true;
  await refresh_data();
  refreshing.value = false;
});
</script>

<style scoped>
/* 确保列宽一致 */
.flex-1 {
  flex: 1;
}

.collapse_list {
  background-color: #f8f9fa;
  margin-bottom: 3px;
}

:deep(.van-collapse-item__title--expanded:after) {
  display: none;
}

/* Tab切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.transition-container {
  transition: opacity 0.3s ease;
}
.add_wallet {
  background-color: #46f37d;
  border-radius: 4px;
  font-size: 14px;
  padding: 3px 10px;
  align: "center";
  cursor: "pointer";
}
</style>
