<script setup lang="ts">
import { showToast } from "vant";
import { ref } from "vue";
import TokenInfo from "@/components/token-info.vue";
const searchList = ref([]);
const showSearchModal = ref(false);
const hotList = ref([
  {
    isCollected: false,
    tokenSymbol: "MMF",
    price: "1.221",
    amount: "11",
    interestRate: "11",
  },
  {
    isCollected: true,
    tokenSymbol: "BTC",
    price: "1.221",
    amount: "11",
    interestRate: "11",
  },
]);
function clickToShowToast() {
  showToast("搜索功能开发中...");
}
</script>
<template>
  <van-popup v-model:show="showSearchModal" position="bottom" round :style="{ height: '70%' }">
    <div class="search-title">Search</div>
    <van-search readonly shape="round" background="white" left-icon="" right-icon="search" @click="clickToShowToast" />
    <div v-if="searchList.length === 0" class="search-content">
      <span class="label">热门搜索</span>
      <div>
        <TokenInfo v-for="(item, index) in hotList" :key="index" :data="item" />
      </div>
    </div>
    <div v-if="searchList.length > 0" class="search-content">
      <span class="label">搜索结果</span>
      <div>
        <TokenInfo v-for="item in searchList" :key="item" :data="item" />
      </div>
    </div>
  </van-popup>
</template>
<style scoped>
/* 搜索弹窗 */
.search-title {
  padding: 24px 24px 5px 24px;
  font-weight: 600;
  font-size: 18px;
}

.search-content {
  padding: 0px 15px;

  .label {
    font-size: 16px;
    padding-bottom: 5px;
  }
}
</style>
