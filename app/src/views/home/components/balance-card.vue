<script setup lang="ts">
import CopyIcon from "@/components/CopyIcon.vue"
import { useAccountStore } from "@/store/modules/account"
import { useHomeStore } from "@/store/modules/home"
import { AccountVault } from "@/utils/AccountVault"
import { BigNumberUtils } from "@/utils/big-number-utils"
import { useAppKitAccount } from "@reown/appkit/vue"
import { showToast } from 'vant'
import { computed, ref, watchEffect } from "vue"
const accountData = useAppKitAccount()

const tokens = computed(() => homeStore.getTokens)

const homeStore = useHomeStore()
const accountStore = useAccountStore()

const wsData = computed(() => accountStore.getWsData)
function copyAddress() {
    try {
        navigator.clipboard.writeText(accountData.value.address)
        showToast("Copy Success!")
    } catch (error) {
    }
}
const currentBalance = ref("")
watchEffect(() => {
    if (wsData.value && tokens.value) {
        let balance = AccountVault.getAccountVaultBalance(wsData.value, tokens.value)
        if (balance) {
            currentBalance.value = BigNumberUtils.floor(BigNumberUtils.round(balance), 2)
        } else {
            currentBalance.value = "--"
        }

    }
}) 
</script>
<template>
    <div class="wallet-header">
        <div class="flex justify-between">
            <div class="wallet-info">
                <div class="wallet-avatar">
                    <span>W</span>
                </div>
                <div class="flex flex-col gap-1">
                    <span class="wallet-name">Wallet001</span>
                    <div class="flex items-center gap-1">
                        <span class="wallet-address">{{ accountData && accountData.address
                            ? accountData.address.substring(0, 10) + "..." +
                            accountData.address.substring(accountData.address.length - 10, accountData.address.length) :
                            ""
                            }}</span>
                        <div @click="copyAddress">
                            <CopyIcon />
                        </div>
                    </div>
                </div>
            </div>
            <div class="balance-section">
                <div class="balance-label">{{ $t('labels.estimated_assets') }}</div>
                <div class="balance-amount">{{ currentBalance }}</div>
            </div>
        </div>
    </div>
</template>
<style scoped lang="css">
.wallet-header {
    margin: 0px 16px;
    background-color: #2c2c2c;
    color: white;
    padding: 16px;
    border-radius: 10px;
    overflow: hidden;
}

.wallet-info {
    display: flex;
    align-items: center;
    gap: 8px;
}

.wallet-avatar {
    width: 40px;
    height: 40px;
    background-color: #4caf50;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    font-weight: bold;
    color: white;
}

.wallet-name {
    font-size: 14px;
    font-weight: 500;
}

.wallet-address {
    font-size: 8px;
}

.balance-section {
    text-align: right;
}

.balance-label {
    font-size: 8px;
    color: white;
}

.balance-amount {
    font-size: 20px;
    font-weight: bold;
    color: white;
}
</style>