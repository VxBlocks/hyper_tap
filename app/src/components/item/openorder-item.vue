<script setup lang="ts">
import { useAccountStore } from '@/store/modules/account';
import { useHomeStore } from '@/store/modules/home';
import { FopenOrderState } from '@/store/types';
import { formatNumber, } from "@/utils/common";
import { formatTimestamp } from "@/utils/time-utils";
import { useAppKitAccount } from '@reown/appkit/vue';
import { privateKeyToAccount } from 'viem/accounts';
import { computed } from 'vue';
import { useI18n } from "vue-i18n";

import * as hl from "@nktkas/hyperliquid";
import { ExchangeHelper } from '@/utils/exchange-helper';
import { showToast } from 'vant';
import { BigNumberUtils } from '@/utils/big-number-utils';
const homeStore = useHomeStore();
const coinToAssetMap = computed(() => homeStore.getTokens.coin_to_asset)
const assetTocoinMap = computed(() => homeStore.getTokens.asset_to_coin)
const { t } = useI18n();
const accountData = useAppKitAccount()

// 定义 props
const props = defineProps<{
    item: FopenOrderState;
}>()
// 定义 emits
const emit = defineEmits<{
    (e: 'cancel-order', item: FopenOrderState): void;
}>()
function cancelOrder(item: FopenOrderState) {
    emit('cancel-order', item);
}

function queryCoinName(asset: string) {
    if (coinToAssetMap.value && coinToAssetMap.value.get(asset)) {
        let coin = coinToAssetMap.value.get(asset);
        if (assetTocoinMap.value && assetTocoinMap.value.get(coin)) {
            return assetTocoinMap.value.get(coin);
        }
    }
    return asset;
}
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
    }
}

const handleCancel = async (item: FopenOrderState) => {
    try {
        let hasAddress = homeStore.currentAddressHasAgentddress(accountData.value.address)
        if (hasAddress) {
            let clientData = homeStore.getCurrentClientAddress(accountData.value.address)
            if (ExchangeHelper.queryHasAgentAddress(clientData.walletAddress as any, clientData.agentAddress as any)) {
                const agentAccount = privateKeyToAccount(clientData.agentPrivateKey as `0x${string}`);
                const transport = new hl.HttpTransport({ isTestnet: false })
                let client = new hl.ExchangeClient({
                    wallet: agentAccount,
                    transport: transport,
                    isTestnet: false,
                });
                let assetID = coinToAssetMap.value.get(item.coin)
                let result = await ExchangeHelper.cancelOrder(client, assetID, item.oid)
                showToast('Success')
            } else {
                homeStore.clearClientsAddress()
                let { agentAddress, agentPrivateKey } = await ExchangeHelper.initHyperliquid({});
                homeStore.addClientsAddress(accountData.value.address, agentAddress, agentPrivateKey)
            }
        } else {
            homeStore.clearClientsAddress()
            let { agentAddress, agentPrivateKey } = await ExchangeHelper.initHyperliquid({});
            homeStore.addClientsAddress(accountData.value.address, agentAddress, agentPrivateKey)
        }

    } catch (error) {
        showToast(error.message || 'Failed')
    }

};

// 预先导入所有图标
function queryIconUrl(item: FopenOrderState) {
    let coinName = queryCoinName(item.coin);
    return `/crypto/svg/color/${coinName}.svg`
}

function queryOrderValue() {
    let quantity = props.item.origSz;
    let price = props.item.limitPx;
    return BigNumberUtils.multiply(quantity, price)
}

</script>

<template>
    <div class="openorder-card">
        <div class="openorder-header">
            <div class="flex w-full justify-between">
                <div class="pair-info">
                    <img :src="queryIconUrl(item)" alt="star" onerror="this.src='/crypto/svg/color/BTC.svg'"
                        class="token_icon" />
                    <span class="pair-name">{{ queryCoinName(item.coin) }}</span>
                </div>
                <van-button plain color="#3A957F" size="mini" class="min-w-[59px] rounded-button compact-button"
                    @click="handleCancel(item)">
                    {{ "Cancel" }}
                </van-button>
            </div>
        </div>
        <van-row class="w-full mb-1">
            <van-col span="8">
                <div class="flex gap-2">
                    <div :class="item.side === 'A' ? 'leverage-tag2' : 'leverage-tag1'">
                        {{ item.side === 'A' ? t('badge.short') : t('badge.long') }}
                    </div>
                    <div class="timevalue">
                        {{ formatTimestamp(item.timestamp) }}
                    </div>
                </div>

            </van-col>
        </van-row>
        <div class="mb-2">
            <van-row class="w-full mb-1">
                <van-col span="8">
                    <div class="detail-item">
                        <div class="label">{{ "Quantity" }}</div>
                        <div class="value">
                            {{ formatNumber(item.origSz) }}</div>
                    </div>
                </van-col>
                <van-col span="10">
                    <div class="detail-item">
                        <div class="label">{{ "Transaction quantity" }}</div>
                        <div class="value">--</div>
                    </div>
                </van-col>
                <van-col span="6">
                    <div class="detail-item">
                        <div class="label">{{ "Price" }}</div>
                        <div class="value">${{ formatNumber(item.limitPx) }}
                        </div>
                    </div>
                </van-col>
            </van-row>
            <van-row class="w-full">
                <van-col span="8">
                    <div class="detail-item">
                        <div class="label">{{ "Reduce Only" }}</div>
                        <div class="value">
                            <span style="color: #979797;">{{ item.reduceOnly ? "Yes" : "No" }}</span>
                        </div>
                    </div>
                </van-col>
                <van-col span="10">
                    <div class="detail-item">
                        <div class="label">{{ "TP/SL" }}</div>
                        <div class="value">--</div>
                    </div>
                </van-col>
                <van-col span="6">
                    <div class="detail-item">
                        <div class="label">{{ "Order Value" }}</div>
                        <div class="value">{{ queryOrderValue() }}</div>
                    </div>
                </van-col>
            </van-row>
        </div>
    </div>
</template>

<style lang="css" scoped>
.openorder-card {
    background-color: #FBFBFB;
    border-radius: 12px;
    padding: 10px 16px;
    margin-bottom: 16px;
}

.openorder-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 4px;
}

.pair-info {
    display: flex;
    align-items: center;
    gap: 4px;
}

.pair-name {
    font-size: 10px;
    font-weight: 600;
    color: #333;
}

.timevalue {
    color: #D9D9D9;
    font-size: 10px;
}

.leverage-tag1 {
    padding: 2px 8px;
    border-radius: 2px;
    background-color: #c5ede7;
    font-size: 8px;
    color: #46CCB9;
    font-weight: bold;
}

.leverage-tag2 {
    padding: 2px 8px;
    border-radius: 2px;
    background-color: #fac4cc;
    font-size: 8px;
    color: #F7435D;
    font-weight: bold;
}

.direction-tag {
    margin-left: 4px;
}

.pnl {
    font-size: 10px;
    font-weight: 600;
}

.positive {
    color: #3A957F;
}

.negative {
    color: #F7435D;
}

.position-details {
    margin-bottom: 20px;
}

.detail-item {
    flex: 1;
    text-align: left;
}


.label {
    font-size: 10px;
    color: #D9D9D9;
    font-weight: bold;
}

.value {
    font-size: 10px;
}

.highlight {
    color: #46CCB9;
}

.valuered {
    color: #F75D43;
}

.red {
    color: #F75D43;
}

.gray {
    color: #999;
}

.action-buttons {
    display: flex;
    gap: 8px;
    align-items: center;
}

.token_icon {
    width: 12px;
    height: 12px;
}

.trade-btn {
    flex: 1;
    background-color: #10b981;
}

:deep(.van-button--small) {
    height: 32px;
    font-size: 8px;
}

:deep(.van-tag--mini) {
    padding: 2px 6px;
    font-size: 8px;
}
</style>