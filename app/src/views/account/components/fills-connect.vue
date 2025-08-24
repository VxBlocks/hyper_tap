<!-- 个人中心历史成交 -->
<script lang="ts" setup>
import { computed, ref } from "vue";
import { useAccountStore } from "@/store/modules/account";
import { formatTimestamp } from "@/utils/time-utils";
import { formatNumber } from "@/utils/common";
import { BigNumberUtils } from "@/utils/big-number-utils";
import { useHomeStore } from "@/store/modules/home";
const accountStore = useAccountStore(); 
const homeStore = useHomeStore();
const activeNames = ref([]);
const fillOrders = computed(() => accountStore.getFillsState.fills);

const coinToAssetMap = computed(() => homeStore.getTokens.coin_to_asset);
const assetTocoinMap = computed(() => homeStore.getTokens.asset_to_coin);
const handleSend = () => { };

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
    <div v-if="fillOrders && fillOrders.length > 0">
        <van-collapse v-model="activeNames" accordion>
            <van-collapse-item v-for="(item, index) in (fillOrders.sort((a, b) => a.time - b.time))" :key="index"
                :name="index">
                <template #title>
                    <div class="asset-title">
                        <div class="detail-col">
                            <p class="label">{{ $t("labels.coin") }}</p>
                            <span class="loss font-bold">{{ queryCoinName(item.coin) }}</span>
                        </div>
                        <div class="detail-col">
                            <p class="label">{{ $t("labels.time") }}</p> 
                            <span class="value">{{ formatTimestamp(item.time) }}</span>
                        </div>
                        <div class="detail-col">
                            <p class="label">{{ $t("labels.quantity") }}</p>
                            <span class="loss">
                                {{ formatNumber(item.sz) }}
                            </span>
                        </div>
                    </div>
                </template>
                <div>
                    <div class="asset-title mb-4">
                        <div class="detail-col">
                            <p class="label">{{ $t("labels.direction") }}</p>
                            <span class="loss">{{ item.dir }}</span>
                        </div>
                        <div class="detail-col">
                            <p class="label">{{ $t("labels.price") }}</p>
                            <span class="value">{{ formatNumber(item.px) }}</span>
                        </div>
                        <div class="detail-col">
                            <p class="label">{{ $t("labels.trade_value") }}</p>
                            <span class="value">{{
                                formatNumber(BigNumberUtils.multiply(item.px, item.sz))
                                }}</span>
                        </div>
                    </div>
                    <div class="asset-title">
                        <div class="detail-col">
                            <p class="label">{{ $t("labels.closed_pnl") }}</p> 
                            <span class="value">{{ item.closedPnl }}</span>
                        </div>
                        <div class="detail-col">
                            <p class="label">{{ $t("labels.fee") }}</p>
                            <span class="value">{{ item.fee }}</span>
                        </div>
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
</style>
