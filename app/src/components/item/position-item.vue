<script setup lang="ts">
import { useAccountStore } from '@/store/modules/account';
import { useHomeStore } from '@/store/modules/home';
import { BigNumberUtils } from '@/utils/big-number-utils';
import { AssetPosition } from '@nktkas/hyperliquid';
import { formatNumber } from "@/utils/common";
import { computed } from 'vue';
import { useI18n } from "vue-i18n";
const accountStore = useAccountStore();
const homeStore = useHomeStore()
const spotAssetCtxs = computed(() => accountStore.getAssetCtxs)
const coinToAsset = computed(() => homeStore.getTokens.coin_to_asset)
const { t } = useI18n();

// 定义 props
const props = defineProps<{     // 控制显示隐藏
    item: AssetPosition;          // 仓位数据 
}>()
// 定义 emits
const emit = defineEmits<{
    (e: 'close-position', item: AssetPosition): void;
    (e: 'trade', item: AssetPosition): void;
}>()
// 处理关闭仓位按钮点击事件
function handleClosePosition(item: AssetPosition) {
    emit('close-position', item);
}
// 处理交易按钮点击事件
function handleTrade(item: AssetPosition) {
    emit('trade', item);
}


function queryMarkPrice(coin: string) {
    let index = coinToAsset.value.get(coin)
    let currentAssetCtxs = spotAssetCtxs.value[index]
    return currentAssetCtxs.markPx
}

// 盈亏
function queryPnlAmount(data: AssetPosition) {
    let markPx = queryMarkPrice(data.position.coin) // 最新标记价格
    let entryPx = data.position.entryPx // 入手平均价格
    let szi = data.position.szi // 持仓数量
    return BigNumberUtils.round(BigNumberUtils.multiply(BigNumberUtils.subtract(markPx, entryPx), szi))
}
// 盈亏比例
function queryPnlScale(data: AssetPosition) {
    let pnlAmount = queryPnlAmount(data) // 盈亏金额
    let markPx = queryMarkPrice(data.position.coin) // 最新标记价格
    let entryValue = BigNumberUtils.multiply(markPx, Math.abs(Number(data.position.szi))) // 仓位价值
    let leverageValue = data.position.leverage.value // 倍数
    let scale = BigNumberUtils.multiply(BigNumberUtils.divide(pnlAmount, entryValue), leverageValue)
    return BigNumberUtils.multiply(BigNumberUtils.round(scale, 4), 100)
}

// 预先导入所有图标
function queryIconUrl(item: AssetPosition) {
    return `/crypto/svg/color/${item.position.coin}.svg`
}
</script>

<template>
    <div class="position-card">
        <div class="position-header">
            <div class="pair-info">
                <img :src="queryIconUrl(item)" alt="star" onerror="this.src='/crypto/svg/color/BTC.svg'"
                    class="token_icon" />
                <span class="pair-name">{{ item.position.coin }}</span>
                <div :class="Number(item.position.szi) < 0 ? 'leverage-tag2' : 'leverage-tag1'">
                    {{ item.position.leverage.value + "x" }}
                </div>
                <div :class="Number(item.position.szi) < 0 ? 'leverage-tag2' : 'leverage-tag1'">
                    {{ Number(item.position.szi) > 0 ? t('badge.long') : t('badge.short') }}
                </div>
            </div>
            <div :class="(Number(queryPnlAmount(item)) > 0 ? 'pnl positive' : 'pnl negative')">
                {{ (Number(queryPnlAmount(item)) > 0 ? "+$" : "-$") + Math.abs(Number(queryPnlAmount(item))) }} ({{
                    queryPnlScale(item) +
                    "%" }})
            </div>
        </div>
        <div class="mb-2">
            <van-row class="w-full mb-1">
                <van-col span="8">
                    <div class="detail-item">
                        <div class="label">{{ $t("labels.quantity") }}</div>
                        <div :class="Number(item.position.szi) > 0 ? 'value highlight' : 'value red'">
                            {{ formatNumber(Math.abs(Number(item.position.szi))) }}</div>
                    </div>
                </van-col>
                <van-col span="10">
                    <div class="detail-item">
                        <div class="label">{{ $t("labels.position_value") }}</div>
                        <div class="value">${{
                            formatNumber(BigNumberUtils.floor(item.position.positionValue, 4)) }}
                        </div>
                    </div>
                </van-col>
                <van-col span="6">
                    <div class="detail-item">
                        <div class="label">{{ $t("labels.margin") }}</div>
                        <div class="value">${{ formatNumber(BigNumberUtils.floor(item.position.marginUsed, 4))
                        }}
                        </div>
                    </div>
                </van-col>
            </van-row>
            <van-row class="w-full">
                <van-col span="8">
                    <div class="detail-item">
                        <div class="label">{{ $t("labels.entry") }}/{{ $t("labels.mark") }}</div>
                        <div class="value">
                            <span style="color: #979797;">{{ formatNumber(item.position.entryPx) }}/</span>
                            <span :class="Number(item.position.szi) > 0 ? 'value red' : 'value highlight'">{{
                                formatNumber(queryMarkPrice(item.position.coin)) }}</span>
                        </div>
                    </div>
                </van-col>
                <van-col span="10">
                    <div class="detail-item">
                        <div class="label">{{ $t("labels.liq_price") }}</div>
                        <div class="value">--</div>
                    </div>
                </van-col>
                <van-col span="6">
                    <div class="detail-item">
                        <div class="label">{{ $t("labels.tp_sl") }}</div>
                        <div class="value">--</div>
                    </div>
                </van-col>
            </van-row>
        </div>
        <div class="flex justify-between">
            <div class="flex gap-2">
                <van-button plain color="#3A957F" size="mini" class="min-w-[59px] rounded-button compact-button"
                    @click="handleClosePosition(item)">
                    {{ $t('btn.close_position') }}
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
.position-card {
    background-color: #FBFBFB;
    border-radius: 12px;
    padding: 10px 16px;
    margin-bottom: 16px;
}

.position-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
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
    color: #333;
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