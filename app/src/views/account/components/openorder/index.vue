<!--个人中心当前委托-->
<script setup lang="ts">
import { computed, ref } from "vue";
import { useAccountStore } from "@/store/modules/account";
import { useHomeStore } from "@/store/modules/home";
import { ExchangeHelper } from "@/utils/exchange-helper";
import { useAppKitAccount } from "@reown/appkit/vue";
import { showToast } from "vant";

const accountStore = useAccountStore();
const homeStore = useHomeStore();
const fopenOrders = computed(() => accountStore.getFopenOrderState)
const tokens = computed(() => homeStore.getTokens)
const accountData = useAppKitAccount()
const handleLoading = ref(false)

async function batchCancelOrders() {
    handleLoading.value = true
    try {
        let hasAddress = homeStore.currentAddressHasAgentddress(accountData.value.address)
        if (hasAddress) {
            let clientData = homeStore.getCurrentClientAddress(accountData.value.address)
            if (ExchangeHelper.queryHasAgentAddress(clientData.walletAddress as any, clientData.agentAddress as any)) {
                let result = await ExchangeHelper.batchRequestExchangeCancelMessage(
                    clientData.agentPrivateKey as `0x${string}`,
                    fopenOrders.value,
                    tokens.value,
                )
                console.log("sell:::", result);
                showToast("Success")
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

function openCancelOrder() {

}

</script>

<template>
    <div class="flex flex-col gap-y-1">
        <div class="flex justify-end mt-0 px-2" v-if="(fopenOrders && fopenOrders.length > 0)">
            <van-button :loading="handleLoading" color="#3A957F" @click="batchCancelOrders"
                class="min-w-[64px] rounded-button compact-button" size="mini" plain>{{ "Cancel All" }}</van-button>
        </div>

        <div class="position-list">
            <openorder-item v-for="(item, index) in fopenOrders" :item="item" :key="index" :name="index"
                @close-position="openCancelOrder" />
        </div>

    </div>
    <Empty v-if="!(fopenOrders && fopenOrders.length > 0)" />
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