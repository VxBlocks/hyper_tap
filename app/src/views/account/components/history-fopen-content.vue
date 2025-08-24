<!-- 个人中心历史委托 -->
<script lang="ts" setup>
import { computed, ref } from "vue";
import { useAccountStore } from "@/store/modules/account";
import Empty from "@/components/empty.vue";
import { useHomeStore } from "@/store/modules/home";

const accountStore = useAccountStore();
const homeStore = useHomeStore();
const activeNames = ref([]);
const fopenOrders = computed(() => accountStore.getFopenOrderState);
const coinToAssetMap = computed(() => homeStore.getTokens.coin_to_asset);
const assetTocoinMap = computed(() => homeStore.getTokens.asset_to_coin); 
const handleSend = () => { };

function queryCoinName(asset: string) {
  if (coinToAssetMap.value && coinToAssetMap.value.get(asset)) {
    let coin = coinToAssetMap.value.get(asset);
    if (assetTocoinMap.value && assetTocoinMap.value.get(coin)) {
      return assetTocoinMap.value.get(coin);
    }
  }
  return asset;
}
</script>
<template>
  <div v-if="fopenOrders && fopenOrders.length > 0">
    <van-collapse v-model="activeNames" accordion>
      <van-collapse-item v-for="(item, index) in fopenOrders" :key="index" :name="index" class="collapse_list">
        <template #title>
          <div class="asset-title">
            <div class="detail-col">
              <p class="label">{{ $t("labels.coin") }}</p>
              <span class="value">HYPE</span>
            </div>
            <div class="detail-col">
              <p class="label">方向</p>
              <span class="value">45.28</span>
            </div>
            <div class="detail-col">
              <p class="label">{{ $t("labels.quantity") }}</p>
              <span class="value">45.28</span>
            </div>
          </div>
        </template>
        <div>
          <div class="asset-title">
            <div class="detail-col">
              <p class="label">{{ $t("labels.time") }}</p>
              <span class="value">HYPE</span>
            </div>
            <div class="detail-col">
              <p class="label">类型</p>
              <span class="value">45.28</span>
            </div>
            <div class="detail-col">
              <p class="label">触发条件</p>
              <span class="value">45.28</span>
            </div>
          </div>
          <div class="asset-title mt-4">
            <div class="detail-col">
              <p class="label">价格</p>
              <span class="value">HYPE</span>
            </div>
            <div class="detail-col">
              <p class="label">成交大小</p>
              <span class="value">HYPE</span>
            </div>
          </div>
          <div class="asset-title mt-4">
            <div class="detail-col">
              <p class="label">状态</p>
              <span class="value">HYPE</span>
            </div>
          </div>
        </div>
      </van-collapse-item>
    </van-collapse>
  </div>
  <div v-else>
    <Empty />
  </div>
</template>

<style lang="scss">
.collapse_list {
  background-color: #f8f9fa;
  margin-bottom: 3px;
}

.asset-title {
  display: flex;
  width: 100%;

  .detail-col {
    width: 30%;

    .label {
      color: #d9d9d9;
    }

    .value {
      color: #000000;
    }

    .loss {
      color: #3a957f;
    }

    .profit {
      color: #f7435d;
    }
  }
}

.sent_button {
  color: #46ccb9;
  cursor: pointer;
}
</style>
