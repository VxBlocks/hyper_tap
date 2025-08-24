<script setup lang="ts">
import { showToast } from 'vant';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAppKitAccount, useAppKitProvider, Provider } from '@reown/appkit/vue'
import { arbitrum } from "viem/chains";
import { createPublicClient, createWalletClient, custom, erc20Abi, formatUnits, http, parseUnits } from "viem";
import '/src/styles/dialog.css'

const { t } = useI18n()

const showDepositModal = defineModel<boolean>('show', { default: false })
const depositAmount = ref('')
const maxAmount = ref('0')
const handleLoading = ref(false)

const accountData = useAppKitAccount()
const provider = useAppKitProvider<Provider>("eip155");

// 查询用户USDC余额
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

// 获取用户最大可存款金额
const loadMaxAmount = async () => {
    if (accountData.value.address) {
        try {
            const { balance } = await queryMainnetUsdc(accountData.value.address)
            maxAmount.value = balance
        } catch (error) {
            console.error('Failed to load balance:', error)
            showToast(t('tip.load_balance_failed'))
        }
    }
}

const setMaxAmount = () => {
    depositAmount.value = maxAmount.value
}

// 确认存款
const confirmDeposit = async () => {
    handleLoading.value = true
    if (!depositAmount.value) {
        showToast(t('tip.enter_deposit_amount'))
        return
    } 
    try {
        const client = createWalletClient({
            chain: arbitrum,
            transport: custom(provider.walletProvider!),
        })

        const accounts = await client.getAddresses()
        const valueArg = parseUnits(depositAmount.value, 6) // USDC有6位小数

        const res = await client.writeContract({
            chain: null,
            account: accounts[0],
            address: "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
            abi: erc20Abi,
            functionName: "transfer",
            args: ["0x2Df1c51E09aECF9cacB7bc98cB1742757f163dF7", valueArg],
        })

        console.log('Deposit transaction:', res)
        showToast(t('tip.deposit_success'))
        showDepositModal.value = false
        depositAmount.value = ''
    } catch (error) {
        console.error('Deposit failed:', error)
        showToast(t('tip.deposit_failed'))
    } finally {
        handleLoading.value = false
    }
}

// 当模态框打开时加载最大金额
watch(showDepositModal, (newVal) => {
    if (newVal) {
        loadMaxAmount()
    }
})
</script>
<template>
    <van-popup v-model:show="showDepositModal" position="bottom" round :style="{ height: '160px' }">
        <div class="dialog">
            <div class="d_t">
                <h3>{{ t('labels.deposit') }}</h3>
            </div>
            <div class="d_msg">
                {{ t('labels.deposit_usdc_arbitrum') }}
            </div>
            <div class="d_c">
                <van-field v-model="depositAmount"  :placeholder="t('labels.amount')" type="number"
                    class="d_i_f">
                    <template #button>
                        <van-button size="mini" type="primary" plain @click="setMaxAmount" class="d_i_btn">
                            {{ t('labels.max_value') }}: {{ maxAmount }}
                        </van-button>
                    </template>
                </van-field>

                <van-button :loading="handleLoading" size="small" color="#3a957f" round block @click="confirmDeposit"
                    class="d_d_btn">
                    {{ t('btn.deposit') }}
                </van-button>
            </div>
        </div>
    </van-popup>
</template>
<style scoped>
/* 存款弹窗 */
.deposit-modal {
    padding: 24px;
}


.modal-header {
    text-align: center;
    margin-bottom: 24px;
}

.modal-header h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
}

.modal-content {
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.max-btn {
    border: none;
    color: #3a957f;
    font-size: 10px;
}

.amount-field {
    border: 2px solid #d9d9d9;
    border-radius: 8px;
}

.deposit-btn {
    border: none;
    border-radius: 8px;
}
</style>