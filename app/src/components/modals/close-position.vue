<!-- 平仓弹出框组件 -->
<script setup lang="ts">

import { computed, onMounted, ref, watch, watchEffect } from 'vue';
import { AssetPosition } from '@nktkas/hyperliquid';
import { useAccountStore } from '@/store/modules/account';
import { useHomeStore } from '@/store/modules/home';
import { BigNumberUtils } from '@/utils/big-number-utils';
import { useAppKitAccount } from '@reown/appkit/vue';
import { ExchangeHelper } from '@/utils/exchange-helper';
import { showToast } from 'vant';
const homeStore = useHomeStore()
const accountStore = useAccountStore()
const spotAssetCtxs = computed(() => accountStore.getAssetCtxs)
const coinToAsset = computed(() => homeStore.getTokens.coin_to_asset)
const tokens = computed(() => homeStore.getTokens)

const accountData = useAppKitAccount()

// 定义 props
const props = defineProps<{
    modelValue: boolean;           // 控制显示隐藏
    item?: AssetPosition;          // 仓位数据
    amount?: number;               // 数量
}>()
const price = ref("")
const quantity = ref("0")
const selectedType = ref('left');
const progress = ref(0);
const handleLoading = ref(false);
onMounted(() => {
    price.value = "";
    selectedType.value = 'left';
    progress.value = 0;
})
watchEffect(() => {
    if (props.item) {
        let szi = props.item.position.szi;
        price.value = "";
        quantity.value = Math.abs(Number(szi)).toString();
        progress.value = 0;
    }
})
watch(selectedType, () => {
    if (props.item) {
        let szi = props.item.position.szi;
        price.value = "";
        if (selectedType.value === 'left') {
            price.value = queryMarkPrice(props.item.position.coin);
        }
        quantity.value = Math.abs(Number(szi)).toString();
        progress.value = 0;
    }
})
function queryMarkPrice(coin: string) {
    let index = coinToAsset.value.get(coin)
    let currentAssetCtxs = spotAssetCtxs.value[index]
    return currentAssetCtxs.markPx
}
function handleMind() {
    let market = queryMarkPrice(props.item.position.coin)
    price.value = market
}

async function handleConfirm(isLimit: boolean) {
    console.log("handleConfirm", isLimit);
    handleLoading.value = true;
    try {
        let exchangePrice = price.value
        if (!isLimit) {
            exchangePrice = queryMarkPrice(props.item.position.coin)
        }
        let hasAddress = homeStore.currentAddressHasAgentddress(accountData.value.address)
        if (hasAddress) {
            let clientData = homeStore.getCurrentClientAddress(accountData.value.address)
            if (ExchangeHelper.queryHasAgentAddress(clientData.walletAddress as any, clientData.agentAddress as any)) {
                let isBuy = Number(props.item.position.szi) < 0 ? true : false;
                let result = await ExchangeHelper.requestExchangeMessage(
                    clientData.agentPrivateKey as `0x${string}`,
                    props.item.position.coin,
                    Number(quantity.value),
                    tokens.value,
                    exchangePrice,
                    isLimit ? 0 : 0.08,
                    isBuy,
                    isLimit)
                console.log("sell:::", result);
                showBottom.value = false;
            } else {
                homeStore.clearClientsAddress()
                let { agentAddress, agentPrivateKey } = await ExchangeHelper.initHyperliquid({});
                homeStore.addClientsAddress(accountData.value.address, agentAddress, agentPrivateKey)
            }

        }


    } catch (error) {
        console.log(error);
        showToast(error.message)
    } finally {
        handleLoading.value = false;
    }
}


// 定义 emits
const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void;
    (e: 'confirm', data: any): void;
    (e: 'cancel'): void;
}>()

// 计算属性，用于 v-model 双向绑定
const showBottom = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
})

function handlePriceInput(value: string) {
    progress.value = 0
    price.value = value
}

function handleChange(changeType: string) {
    selectedType.value = changeType;
}
watch(progress, () => {
    if (progress.value != 0) {
        let scale = BigNumberUtils.divide(progress.value, 100);
        let szi = Math.abs(Number(props.item.position.szi));
        let newSzi = BigNumberUtils.multiply(szi, scale);
        quantity.value = newSzi;
    } else {
        quantity.value = Math.abs(Number(props.item.position.szi)).toString();
    }
})

// 拖动开始
const onDragStart = () => {
    console.log('开始拖动');
};

// 拖动结束
const onDragEnd = () => {
    console.log('结束拖动');
};

// 进度变化
const onProgressChange = (value: number) => {
    console.log('进度变化:', value);
};

</script>

<template>
    <van-popup v-model:show="showBottom" round position="bottom" :style="{ height: '340px', }">
        <div class="popup-content">
            <div class="flex justify-center mt-1">
                <span class="title">{{ $t("labels.close_position") }}</span>
            </div>
            <div class="flex justify-center mt-1">
                <TrapezoidButtons :left-text='$t("labels.market_price")' :right-text='$t("labels.limit_price")'
                    :initial-selected="selectedType" @change="handleChange" />
            </div>
            <div class="flex justify-center mt-2 smallSize">
                {{ selectedType === 'left' ? $t("labels.close_position_msg") : $t("labels.limit_close_position_msg") }}
            </div>
            <div class="flex flex-col gap-3 mt-3">
                <van-field v-if='selectedType != "left"' v-model="price" class="custom-field" input-align="right" center
                    :label="$t('labels.close_position_price')" label-width="100px"
                    style="border: 1px solid #F2F3F6; border-radius: 8px;">
                    <template #button>
                        <span style="color: #3A957F;cursor: pointer;font-size: 10px;" @click="handleMind"> Mind </span>
                    </template>
                    <template #error-message></template>
                </van-field>
                <van-field v-else class="custom-field" input-align="right" center :label="$t('labels.price')"
                    label-width="100px" style="border: 1px solid #FFFFFF; border-radius: 8px;">
                    <template #button>
                        <span style="color: #A0A0A0;cursor: pointer;font-size: 10px;"> {{ $t("labels.market_price") }}
                        </span>
                    </template>
                    <template #error-message></template>
                </van-field>
                <van-field :label="$t('labels.quantity')" @input="handlePriceInput" v-model="quantity"
                    class="custom-field" type="number" input-align="right" size="normal"
                    style="border: 1px solid #F2F3F6; font-size: 10px; border-radius: 8px;">
                    <!-- <template #extra>{{ "ETH" }}</template> -->
                </van-field>

                <div class="flex mx-1">
                    <van-slider v-model="progress" :min="0" :max="100" bar-height="4px" active-color="#3A957F"
                        @change="onProgressChange" @drag-start="onDragStart" @drag-end="onDragEnd">
                        <template #button>
                            <div class="custom-button">{{ progress + "%" }}</div>
                        </template></van-slider>
                </div>
                <div class="flex justify-center message">
                    {{ $t("labels.close_buttom_msg") }} : {{ item.position.unrealizedPnl }}
                </div>
                <div v-if="selectedType === 'left'" class="flex justify-center message">
                    <van-button :loading="handleLoading" size="mini"
                        :disabled="!(Number(quantity) && Number(quantity) <= Math.abs(Number(props.item.position.szi)))"
                        color="#3A957F" class="rounded-button " block @click="handleConfirm(false)">{{
                            $t("labels.market_close_position")
                        }}</van-button>
                </div>
                <div v-else class="flex justify-center message">
                    <van-button :loading="handleLoading" size="mini"
                        :disabled="!(Number(price) && Number(quantity) && Number(quantity) <= Math.abs(Number(props.item.position.szi)))"
                        color="#3A957F" class="rounded-button " block @click="handleConfirm(true)">{{
                            $t("labels.limit_close_position")
                        }}</van-button>
                </div>
            </div>
        </div>
    </van-popup>
</template>

<style scoped lang="less">
.popup-content {
    height: 100%;
    padding: 16px;
    display: flex;
    flex-direction: column;
    padding-bottom: env(safe-area-inset-bottom);
}

.button-container {
    display: flex;
    background-color: #f5f5f5;
    background-image:
        linear-gradient(#e0e0e0 1px, transparent 1px),
        linear-gradient(90deg, #e0e0e0 1px, transparent 1px);
    background-size: 10px 10px;
    padding: 8px;
    border-radius: 4px;
    width: fit-content;
}

.button {
    padding: 12px 24px;
    font-weight: bold;
    position: relative;
    cursor: pointer;
    text-align: center;
    min-width: 80px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.left-button {
    background-color: #e0e0e0;
    color: #333;
    /* 直角梯形 - 右上角切掉 */
    clip-path: polygon(0 0, calc(100% - 20px) 0, 100% 100%, 0 100%);
    margin-right: -10px;
    /* 使两个按钮重叠 */
    z-index: 1;
}

.right-button {
    background-color: #3ac7a8;
    color: white;
    /* 直角梯形 - 左下角切掉 */
    clip-path: polygon(0 0, 100% 0, 100% 100%, 20px 100%);
    margin-left: -10px;
    /* 使两个按钮重叠 */
}

.left-button.active {
    background-color: #d0d0d0;
}

.right-button.active {
    background-color: #2db598;
}


.title {
    font-size: 16px;
    font-weight: bold;
    margin-bottom: 20px;
}

.smallSize {
    font-size: 10px;
}

.message {
    color: #3A957F;
    font-size: 10px;
}

.labelTitle {
    font-size: 10px;
    color: #78817E;
}

.custom-field {
    :deep(.van-field__label) {
        font-size: 10px !important;
        color: #A0A0A0;
    }

    --van-field-input-text-color: #3A957F;
    --van-field-text-color: #3A957F;
}


/* 在 style 中添加 */
.rounded-button {
    border-radius: 8px;
    height: 30px;
    font-weight: 600;
    font-size: 10px;
}

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

.custom-button {
    min-width: 26px;
    padding: 0 4px;
    color: #fff;
    font-size: 10px;
    line-height: 18px;
    text-align: center;
    background-color: #3A957F;
    border-radius: 100px;
}
</style>