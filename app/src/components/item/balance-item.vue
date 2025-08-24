<script setup lang="ts">
import { useAccountStore } from '@/store/modules/account';
import { useHomeStore } from '@/store/modules/home';
import { BalanceViewData } from '@/store/types';
import { BigNumberUtils } from '@/utils/big-number-utils';
import { computed } from 'vue';
import { useI18n } from "vue-i18n";
import { formatNumber, formatNumberWithCommas } from "@/utils/common";
const accountStore = useAccountStore();
const homeStore = useHomeStore()
const coinIndexToPrice = computed(() => homeStore.getTokens.coin_index_to_price)
const spotAssetCtxs = computed(() => accountStore.getSpotAssetCtxs)
const { t } = useI18n();

// 定义 props
const props = defineProps<{     // 控制显示隐藏
    item: BalanceViewData;          // 现货数据
}>()
// 定义 emits
const emit = defineEmits<{
    (e: 'sell', item: BalanceViewData): void;
    (e: 'trade', item: BalanceViewData): void;
}>()
function handleSell(item: BalanceViewData) {
    emit('sell', item);
}
// 处理交易按钮点击事件
function handleTrade(item: BalanceViewData) {
    emit('trade', item);
}


function queryUSDCValue(item: BalanceViewData) {
    let currentToken = item.token
    if (currentToken == 0) {
        return item.totalBalance
    }
    return queryTotalBalance(item)
}
function queryLatestPrice(token: number) {
    let latestPrice = "0"
    if (!coinIndexToPrice.value) {
        return latestPrice
    }
    let priceData = coinIndexToPrice.value.get(token)
    if (!(priceData && priceData.asset_name)) {
        return latestPrice
    }
    let coin = priceData.asset_name
    let currentAssetCtxs = spotAssetCtxs.value.find((item) => item.coin === coin)
    if (currentAssetCtxs) {
        latestPrice = currentAssetCtxs.markPx
    }
    return latestPrice
}

function queryTotalBalance(item: BalanceViewData) {
    let latestPrice = queryLatestPrice(item.token)
    return BigNumberUtils.multiply(item.totalBalance, latestPrice)
}

// 盈亏
function queryLRoe(item: BalanceViewData) {
    if (item.token === 0) {
        return "0"
    }
    let usdcValue = queryUSDCValue(item)
    let entryNtlValue = item.entryNtl
    let a = BigNumberUtils.subtract(usdcValue, entryNtlValue)
    return BigNumberUtils.floor(a, 2)

}
// 比例
function queryPnl(item: BalanceViewData) {
    let rol = queryLRoe(item)
    let entryNtlValue = item.entryNtl
    console.log("queryPnl::::", rol, entryNtlValue);
    let a = BigNumberUtils.divide(Math.abs(Number(rol)), entryNtlValue)
    let b = BigNumberUtils.multiply(a, 100)
    return BigNumberUtils.floor(b, 2)
}


// 预先导入所有图标
function queryIconUrl(item: BalanceViewData) {
    return `/crypto/svg/color/${item.currency}.svg`
}

</script>

<template>
    <div class="balance-card">
        <div class="balance-header">
            <div class="pair-info">
                <img :src="queryIconUrl(item)" alt="star" onerror="this.src='/crypto/svg/color/BTC.svg'"
                    class="token_icon" />
                <span class="pair-name">{{ item.currency }}</span>
                <div class="leverage-tag1">
                    {{ item.type }}
                </div>
            </div>
            <div v-if="item.token != 0" :class="(Number(queryLRoe(item)) > 0 ? 'pnl positive' : 'pnl negative')">
                {{ queryLRoe(item) + "$" }} ({{ (queryPnl(item)) + "%" }})
            </div>
        </div>
        <div class="mb-2">
            <van-row class="w-full mb-2">
                <van-col span="8">
                    <div class="detail-item">
                        <div class="label">{{ $t("labels.total_balance") }}</div>
                        <div class="value">
                            {{ formatNumber(item.totalBalance) }}</div>
                    </div>
                </van-col>
                <van-col span="10">
                    <div class="detail-item">
                        <div class="label">{{ $t("labels.available_balance") }}</div>
                        <div class="value">{{ formatNumberWithCommas(item.availableBalance) }}
                        </div>
                    </div>
                </van-col>
                <van-col span="6">
                    <div class="detail-item">
                        <div class="label">{{ $t("labels.usd_value") }}</div>
                        <div class="value">${{ formatNumber(BigNumberUtils.floor(queryUSDCValue(item), 4)) }}
                        </div>
                    </div>
                </van-col>
            </van-row>
        </div>
        <div class="flex justify-between">
            <div class="flex gap-2">
                <van-button plain color="#3A957F" size="mini" class="min-w-[59px] rounded-button compact-button"
                    @click="handleSell(item)">
                    {{ $t("btn.sell") }}
                </van-button>
            </div>
            <van-button color="#3A957F" size="mini" @click="handleTrade"
                class="min-w-[59px] rounded-button compact-button">
                {{ $t('btn.trade') }}
            </van-button>
        </div>
    </div>
</template>

<style lang="css" scoped>
.balance-card {
    background-color: #FBFBFB;
    border-radius: 12px;
    padding: 10px 16px;
    margin-bottom: 16px;
}

.balance-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
}

.pair-info {
    display: flex;
    align-items: center;
    gap: 8px;
} 
 
.pair-name {
    font-size: 10px;
    font-weight: 600;
    color: #333;
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
    font-weight: bold;
    color: #D9D9D9;
    margin-bottom: 2px;
}

.value {
    font-size: 10px;
    color: #000000;
    font-weight: bold;
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

.token_icon {
    width: 12px;
    height: 12px;
}

.gray {
    color: #999;
}

.action-buttons {
    display: flex;
    gap: 8px;
    align-items: center;
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