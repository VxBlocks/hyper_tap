<script setup lang="ts">
import { showToast } from 'vant';
import { ref } from 'vue';
const showExchandeModal = ref(false)
const exchangeAmount = ref('')
// const exchangeType = ref('Spot->Contract')
const maxAmount = ref('9,999.88')

const setMaxAmount = () => {
    exchangeAmount.value = maxAmount.value
}

const confirmExchange = () => {
    if (!exchangeAmount.value) {
        showToast('请输入兑换金额')
        return
    }

    showToast(`兑换 ${exchangeAmount.value} 成功!`)
    showExchandeModal.value = false
    exchangeAmount.value = ''
}
</script>
<template>
    <van-popup v-model:show="showExchandeModal" position="bottom" round :style="{ height: '27%' }">
        <div class="exchange-modal">
            <div class="modal-header">
                <h3>从转账USDC</h3>
            </div>
            <div class="modal-description">
                在您的永久合约账户和现货账户之间转移USDC。
            </div>
            <div class="transfer-direction">
                <span>现货</span>
                <van-icon name="exchange" class="swap-icon" />
                <span>永久合约</span>
            </div>
            <div class="modal-content">
                <van-field v-model="exchangeAmount" size="normal" placeholder="金额(USDC)" type="number"
                    class="amount-field">
                    <template #button>
                        <van-button size="mini" type="primary" plain @click="setMaxAmount" class="max-btn">
                            最大值: {{ maxAmount }}
                        </van-button>
                    </template>
                </van-field>

                <van-button size="small" color="#3a957f" round block @click="confirmExchange" class="exchange-btn">
                    存款
                </van-button>
            </div>
        </div>
    </van-popup>
</template>
<style scoped>
/* 存款弹窗 */
.exchange-modal {
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

.exchange-btn {
    border: none;
    border-radius: 8px;
}
</style>