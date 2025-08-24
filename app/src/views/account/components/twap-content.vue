<!-- 个人中心TWAP -->
<script lang="ts" setup>
import { computed, ref } from "vue";
import { useAccountStore } from "@/store/modules/account";

const accountStore = useAccountStore();
const activeNames = ref([]);
const balances = computed(() => accountStore.getSpotInfo.balances);
// todo
const datas = false;

const handleTerminate = () => {
};

function queryCoinName() { }
</script>
<template>
  <div v-if="datas">
    <van-collapse v-model="activeNames" accordion>
      <van-collapse-item v-for="(item, index) in [1, 2]" :key="index" :name="index" class="collapse_list">
        <template #title>
          <div class="asset-title">
            <div class="detail-col">
              <p class="label">{{ $t("labels.coin") }}</p>
              <span class="value">HYPE</span>
            </div>
            <div class="detail-col">
              <p class="label">{{ $t("labels.quantity") }}</p>
              <span class="loss">45.28</span>
            </div>
            <div class="detail-col">
              <p class="label">{{ $t("labels.executed_size") }}</p>
              <span class="loss">
                0 ETH
              </span>
            </div>
          </div>
        </template>
        <div>
          <div class="asset-title">
            <div class="detail-col" :style="{ height: '50%' }">
              <p class="label">{{ $t("labels.average_price") }}</p>
              <span class="value">HYPE</span>
            </div>
            <div class="detail-col" :style="{ height: '50%' }">
              <p class="label">{{ $t("labels.running_time") }}/{{ $t("labels.total") }}</p>
              <span class="value">45.28</span>
            </div>
          </div>
          <div class="asset-title mt-4">
            <div class="detail-col" :style="{ height: '50%' }">
              <p class="label">{{ $t("labels.create_time") }}</p>
              <span class="value">HYPE</span>
            </div>
            <div class="detail-col" :style="{ height: '50%' }">
              <p class="label">{{ $t("labels.only_reduce") }}</p>
              <span class="value">45.28</span>
            </div>
          </div>
          <van-divider />
          <div class="flex gap-10">
            <div class="sent_button" @click="handleTerminate">Terminate</div>
          </div>
        </div>
      </van-collapse-item>
    </van-collapse>
  </div>
  <div v-else>
    <coming-soon />
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
      font-size: 10px;
      color: #D9D9D9;
      font-weight: bold;
    }

    .value {
      font-size: 10px;
      color: #333;
      font-weight: bold;
    }

    .loss {
      font-size: 10px;
      color: #3a957f;
      font-weight: bold;
    }

  }
}

.sent_button {
  color: #46ccb9;
  cursor: pointer;
}
</style>
