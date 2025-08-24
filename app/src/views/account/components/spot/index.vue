<!--个人中心余额-->
<script setup lang="ts">
import { useAccountStore } from '@/store/modules/account';
import { useHomeStore } from '@/store/modules/home';
import { computed, ref } from 'vue';
import { BalanceViewData } from '@/store/types';
import { showToast } from 'vant';
const accountStore = useAccountStore();
const homeStore = useHomeStore();
const balances = computed(() => accountStore.getBalanceViewDatas)
const handleLoading = ref(false);

// 打开弹窗
function openSell(item: BalanceViewData) {
    showToast("In preparation...");
}

// 打开弹窗
function openTrade(item: BalanceViewData) {
    showToast("In preparation...");
}

function batchSell() {
    showToast("In preparation...");
} 
</script>
<template>
    <div class="flex flex-col gap-y-1">
        <div class="flex justify-end mt-0 px-2" v-if="(balances && balances.length > 0)">
            <van-button :loading="handleLoading" color="#3A957F" @click="batchSell"
                class="min-w-[64px] rounded-button compact-button" size="mini" plain>{{ "Sell All" }}</van-button>
        </div>
        <div class="balance-list">
            <balance-item v-for="(item, index) in balances" :item="item" :key="index" :name="index" @sell="openSell"
                @trade="openTrade" />
        </div>

    </div>
</template>

<style lang="css" scoped>
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