<script setup lang="ts">
import { showToast } from "vant";
import { ref } from "vue";
const showSendModal = ref(false);
const depositAmount = ref("");
const maxAmount = ref("9,999.88"); 
const confirmDeposit = () => {
  if (!depositAmount.value) {
    showToast("请输入金额");
    return;
  }

  showToast(`存款 ${depositAmount.value} USDC 成功!`);
  showSendModal.value = false;
  depositAmount.value = "";
};
</script>
<template>
  <van-popup
    v-model:show="showSendModal"
    position="bottom"
    round
    :style="{ height: '27%' }"
  >
    <div class="deposit-modal">
      <div class="modal-header">
        <h3>发送代币</h3>
        <h5>将代币发送到Hyperliquid L1上的另一个账户</h5>
      </div>
      <div class="modal-content">
        <van-field
          v-model="depositAmount"
          size="normal"
          placeholder="目标地址"
          type="number"
          class="amount-field"
        >
        </van-field>
        <van-field
          v-model="depositAmount"
          size="normal"
          placeholder="金额(USDC)"
          type="number"
          class="amount-field"
        >
          <template #button> 最大值: {{ maxAmount }} </template>
        </van-field>

        <van-button
          size="small"
          color="#3a957f"
          round
          block
          @click="confirmDeposit"
          class="deposit-btn"
        >
          确认
        </van-button>
      </div>
    </div>
  </van-popup>
</template>
<style scoped>
/* 存款弹窗 */
.deposit-modal {
  padding: 24px;
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

.max-btn {
  border: none;
  color: #3a957f;
  font-size: 10px;
}

.amount-field {
  border: 2px solid #d9d9d9;
  border-radius: 8px;
}

.deposit-btn {
  border: none;
  border-radius: 8px;
}
</style>
