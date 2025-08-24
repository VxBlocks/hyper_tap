<!--个人中心仓位-->
<script setup lang="ts">
import { computed, ref } from "vue";
import { useAccountStore } from "@/store/modules/account";
import { AssetPosition } from "@nktkas/hyperliquid";
import { useHomeStore } from "@/store/modules/home";

import ClosePositionPopup from '@/components/modals/close-position.vue';
import { ExchangeHelper } from "@/utils/exchange-helper";
import { useAppKitAccount } from "@reown/appkit/vue";
import { showToast } from "vant";
import { useI18n } from "vue-i18n";
const { t } = useI18n();

const accountStore = useAccountStore();
const assetPositions = computed(() => accountStore.getClearinghouseInfo.assetPositions);
const homeStore = useHomeStore()
const spotAssetCtxs = computed(() => accountStore.getAssetCtxs)
const tokens = computed(() => homeStore.getTokens)
const accountData = useAppKitAccount()

// 弹窗控制
const showClosePositionPopup = ref(false);
const currentSelectedItem = ref<AssetPosition | null>(null);
const currentAmount = ref(0);
const handleLoading = ref(false);
// 打开弹窗
function openClosePosition(item: AssetPosition) {
    currentSelectedItem.value = item;
    currentAmount.value = Number(item.position.szi);
    showClosePositionPopup.value = true;
}

// 打开弹窗
function openTrade(item: AssetPosition) {
    showToast("In preparation...");
}


// 处理确认
function handleConfirm(data: any) {
    console.log('确认平仓:', data);
    // 关闭弹窗
    showClosePositionPopup.value = false;
}

// 处理取消
function handleCancel() {
    console.log('取消平仓');
    showClosePositionPopup.value = false;
}
async function batchClosePositions() {
    handleLoading.value = true
    try {
        let hasAddress = homeStore.currentAddressHasAgentddress(accountData.value.address)
        if (hasAddress) {
            let clientData = homeStore.getCurrentClientAddress(accountData.value.address)
            if (ExchangeHelper.queryHasAgentAddress(clientData.walletAddress as any, clientData.agentAddress as any)) {
                let result = await ExchangeHelper.batchRequestExchangeMessage(
                    clientData.agentPrivateKey as `0x${string}`,
                    assetPositions.value,
                    tokens.value,
                    spotAssetCtxs.value
                )
                console.log("sell:::", result);
                showToast("Orders submitted")
            } else {
                homeStore.clearClientsAddress()
                let { agentAddress, agentPrivateKey } = await ExchangeHelper.initHyperliquid({});
                homeStore.addClientsAddress(accountData.value.address, agentAddress, agentPrivateKey)
            } 
        } else {
            let { agentAddress, agentPrivateKey } = await ExchangeHelper.initHyperliquid({});
            homeStore.addClientsAddress(accountData.value.address, agentAddress, agentPrivateKey)
        }
    } catch (error) {
        console.log(error);
        showToast(error.message)
    } finally {
        handleLoading.value = false;
    }
} 
</script>

<template>
    <div class="flex flex-col gap-y-1">
        <div class="flex justify-end mt-0 px-2" v-if="(assetPositions && assetPositions.length > 0)">
            <van-button :loading="handleLoading" color="#3A957F" @click="batchClosePositions"
                class="min-w-[64px] rounded-button compact-button" size="mini" plain>{{ $t("btn.allclosepositions")
                }}</van-button>
        </div>

        <div class="position-list">
            <position-item v-for="(item, index) in assetPositions" :item="item" :key="index" :name="index"
                @close-position="openClosePosition" @trade="openTrade" />

            <!-- 单个弹窗组件，用于所有item -->
            <ClosePositionPopup v-model="showClosePositionPopup" :item="currentSelectedItem" :amount="currentAmount"
                @confirm="handleConfirm" @cancel="handleCancel" />
        </div>

    </div>
    <Empty v-if="!(assetPositions && assetPositions.length > 0)" />
</template>

<style lang="less" scoped>
/* 在 style 中添加 */
.rounded-button {
    border-radius: 6px;
}

:deep(.rounded-button) {
    border-radius: 6px;
}

/* 在你的 style 标签中添加 */
.compact-button {
    font-weight: bold;
    /* 文字加粗 */
    padding-top: 2px;
    /* 减小顶部距离 */
    padding-bottom: 2px;
    /* 减小底部距离 */
}

:deep(.compact-button) {
    font-weight: bold;
    padding-top: 2px;
    padding-bottom: 2px;
    height: auto;
    /* 允许按钮高度自适应内容 */
}
</style>