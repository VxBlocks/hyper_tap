<script setup lang="ts">
import { showToast } from "vant";
import { ref } from "vue";
const showAddWalletModal = ref(false);
const selectWallet = ref("");
const alias = ref("");
const checkedFilters = ref([]);
const checkedInfo = ref([]);

const confirmAdd = () => {};
</script>
<template>
  <van-popup
    v-model:show="showAddWalletModal"
    position="bottom"
    round
    :style="{ height: '27%' }"
  >
    <div class="deposit-modal">
      <div class="modal-header">
        <h3>Add Wallet to Watchlist</h3>
      </div>
      <div class="modal-content">
        <div>
          <span>Select Wallet</span>
          <van-field v-model="selectWallet" size="normal" class="amount-field">
          </van-field>
        </div>
        <div>
          <span>Wallet Alias</span>
          <van-field v-model="alias" size="normal" class="amount-field">
          </van-field>
        </div>
        <div class="checkbox">
          <div class="label">Filters</div>
          <van-checkbox-group
            v-model="checkedInfo"
            shape="square"
            direction="horizontal"
            checked-color="#46F37D"
          >
            <van-checkbox name="a">现货</van-checkbox>
            <van-checkbox name="b">仓位</van-checkbox>
            <van-checkbox name="c">存取</van-checkbox>
            <van-checkbox name="d">挂单</van-checkbox>
            <van-checkbox name="e">TWAP</van-checkbox>
          </van-checkbox-group>
        </div>
        <div class="checkbox">
          <div class="label">Notification Method</div>
          <van-checkbox-group
            v-model="checkedInfo"
            shape="square"
            direction="horizontal"
            checked-color="#46F37D"
          >
            <van-checkbox name="a">APP</van-checkbox>
            <van-checkbox name="b">E-mail</van-checkbox>
          </van-checkbox-group>
        </div>
        <van-button
          size="small"
          color="#3a957f"
          round
          block
          @click="confirmAdd"
          class="add-btn"
        >
          添加
        </van-button>
      </div>
    </div>
  </van-popup>
</template>
<style scoped>
/* 存款弹窗 */
.deposit-modal {
  padding: 24px;
  :deep(.van-checkbox-group--horizontal) {
    gap: 10px;
  }
}

.modal-header {
  text-align: center;
  margin-bottom: 24px;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.modal-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.amount-field {
  border: 2px solid #d9d9d9;
  border-radius: 8px;
}

.checkbox {
  .label {
    margin-bottom: 5px;
  }
}

.add-btn {
  border: none;
  border-radius: 8px;
}
</style>
