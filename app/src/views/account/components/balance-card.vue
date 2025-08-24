<script setup lang="ts">

import { useAppKitAccount, useAppKitProvider, Provider } from '@reown/appkit/vue'
import { showToast } from 'vant'
import { useI18n } from 'vue-i18n';
import { arbitrum } from "viem/chains";
import { ref, watchEffect, computed } from 'vue';
import DepositModal from '@/components/modals/deposit-modal.vue';
import { createPublicClient, createWalletClient, custom, erc20Abi, formatUnits, http, parseUnits } from "viem";
import { useAccountStore } from "@/store/modules/account"
import { AccountVault } from "@/utils/AccountVault"
import { BigNumberUtils } from "@/utils/big-number-utils"
import { useHomeStore } from '@/store/modules/home';

const { t } = useI18n();
const accountData = useAppKitAccount()

const showDepositModal = ref(false)

const handleDeposit = () => {
    showDepositModal.value = true
}

const handleExchange = () => {
    showToast(t('labels.btn_transfer_tip'))
}

const handleWithdraw = () => {
    showToast(t('labels.btn_withdraw_tip'))
}

function copyAddress() {
    try {
        navigator.clipboard.writeText(accountData.value.address)
        showToast("Copy Success!")
    } catch (error) {
    }
}
const accountStore = useAccountStore()
const homeStore = useHomeStore()
const tokens = computed(() => homeStore.getTokens)

const wsData = computed(() => accountStore.getWsData)

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

async function queryMainnetUsdc(address: string): Promise<{
    balance: string;
    decimals: number;
}> {

    const publicClient = createPublicClient({
        chain: arbitrum,
        transport: http("https://arb1.arbitrum.io/rpc"),
    });

    const balance = await publicClient.readContract({
        address: "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
        abi: erc20Abi,
        functionName: "balanceOf",
        args: [address as any],
    });

    const decimals = 6

    const formattedBalance = formatUnits(balance as bigint, decimals);

    return {
        balance: formattedBalance,
        decimals,
    };
}

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

        <div class="action-buttons">
            <div class="action-button" @click="handleDeposit">
                <van-icon name="gold-coin-o" size="24" />
                <span>{{ $t('btn.deposit') }}</span>
            </div>
            <div class="action-button" @click="handleExchange">
                <van-icon name="exchange" size="24" />
                <span>{{ $t('btn.transfer') }}</span>
            </div>
            <div class="action-button" @click="handleWithdraw">
                <van-icon name="cash-back-record" size="24" />
                <span>{{ $t('btn.withdraw') }}</span>
            </div>
        </div>
    </div>
    <DepositModal v-model:show="showDepositModal" />
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
    font-size: 18px;
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

.action-buttons {
    display: flex;
    gap: 16px;
    background-color: rgba(255, 255, 255, 0.1);
    margin-top: 14px;
    padding: 16px;
    border-radius: 12px;
}

.action-button {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    transition: opacity 0.2s;
}

.action-button:active {
    opacity: 0.7;
}

.action-button span {
    font-size: 12px;
    color: white;
}
</style>