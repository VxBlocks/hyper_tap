<!-- 个人中心当前委托 -->
<script lang="ts" setup>
import { computed, ref, watchEffect } from "vue";
import { useAccountStore } from "@/store/modules/account";
import { formatTimestamp } from "@/utils/time-utils";
import Empty from "@/components/empty.vue";


import * as hl from "@nktkas/hyperliquid";
import { formatNumber } from "@/utils/common";
import { useHomeStore } from "@/store/modules/home";
import { ExchangeHelper } from "@/utils/exchange-helper";
import { useAppKitAccount } from "@reown/appkit/vue";
import { privateKeyToAccount } from "viem/accounts";
import { FopenOrderState } from "@/store/types";
import { showConfirmDialog, showToast } from "vant";
const accountStore = useAccountStore();
const homeStore = useHomeStore();
const activeNames = ref([]);
const fopenOrders = computed(() => accountStore.getFopenOrderState)
const coinToAssetMap = computed(() => homeStore.getTokens.coin_to_asset)
const assetTocoinMap = computed(() => homeStore.getTokens.asset_to_coin)
const accountData = useAppKitAccount()
watchEffect(() => {
    console.log(fopenOrders.value);
    
})
function queryInfoClient() {
    try {
        let hasAddress = homeStore.currentAddressHasAgentddress(accountData.value.address)
        if (hasAddress) {
            let clientData = homeStore.getCurrentClientAddress(accountData.value.address)
            const agentAccount = privateKeyToAccount(clientData.agentPrivateKey as `0x${string}`);
            const transport = new hl.HttpTransport({ isTestnet: false })
            return new hl.ExchangeClient({
                wallet: agentAccount,
                transport: transport,
                isTestnet: false,
            });
        } else {
            throw new Error("Invalid private key");
        }
    } catch (error) {
        initAndQueryExchangeClient();
    }
}

async function initAndQueryExchangeClient() {
    try {
        let { agentAddress, agentPrivateKey } = await ExchangeHelper.initHyperliquid({});
        homeStore.addClientsAddress(accountData.value.address, agentAddress, agentPrivateKey)
    } catch (error) {
        console.error(error);
    }
}
const handleCancel = async (item: FopenOrderState) => {
    try {
        let client = queryInfoClient()
        let assetID = coinToAssetMap.value.get(item.coin)
        let result = await ExchangeHelper.cancelOrder(client, assetID, item.oid)
        console.log("sell:::", result);
        showToast('取消订单成功')
    } catch (error) {
        showToast(error.message || '取消订单失败')
    }

};

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
                            <span class="loss font-bold">{{ queryCoinName(item.coin) }}</span>
                        </div>
                        <div class="detail-col">
                            <p class="label">方向</p>
                            <!-- TODO 怎么计算 -->
                            <span class="value">Long</span>
                        </div>
                        <div class="detail-col">
                            <p class="label">{{ $t("labels.quantity") }}</p>
                            <span class="loss">
                                {{ formatNumber(item.origSz) }}
                            </span>
                        </div>
                    </div>
                </template>
                <div>
                    <div class="asset-title mb-4">
                        <div class="detail-col">
                            <p class="label">{{ $t("labels.time") }}</p>
                            <span class="value">{{ formatTimestamp(item.timestamp) }}</span>
                        </div>
                        <div class="detail-col">
                            <p class="label">类型</p>
                            <span class="value">{{ item.orderType }}</span>
                        </div>
                        <div class="detail-col">
                            <p class="label">触发条件</p>
                            <span class="value">{{ item.triggerCondition }}</span>
                        </div>
                    </div>
                    <div class="asset-title">
                        <div class="detail-col">
                            <p class="label">价格</p>
                            <span class="value">{{ formatNumber(item.limitPx) }}</span>
                        </div>
                        <div class="detail-col">
                            <p class="label">成交大小</p>
                            <!-- TODO 怎么获取 -->
                            <span class="value">--</span>
                        </div>
                        <div class="detail-col">
                            <p class="label">TP/SL</p>
                            <!-- TODO 怎么获取 -->
                            <span class="value">--</span>
                        </div>
                    </div>
                    <van-divider />
                    <div class="flex gap-10">
                        <div class="sent_button" @click="handleCancel(item)">取消</div>
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
