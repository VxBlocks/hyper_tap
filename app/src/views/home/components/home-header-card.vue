<script setup lang="ts">
import { onMounted, watchEffect } from "vue";
import BalanceCard from "./balance-card.vue";
import ConnectButton from './connect-button.vue';
import { useAppKitAccount, useDisconnect } from '@reown/appkit/vue';
import { useHomeStore } from "@/store/modules/home";
import { ExchangeHelper } from "@/utils/exchange-helper";
import { ApiService } from "@/utils/api-service";
import { current_network } from "@/constant";
import { useAccountStore } from "@/store/modules/account";
const accountData = useAppKitAccount();
const { disconnect } = useDisconnect();
const homeStore = useHomeStore()
const accountStore = useAccountStore()
const TAG = "[Views]::HeaderCard::"
const queryExchClient = async () => {
    let hasAddress = homeStore.currentAddressHasAgentddress(accountData.value.address)
    if (!hasAddress) {
        try {
            let { agentAddress, agentPrivateKey } = await ExchangeHelper.initHyperliquid({});
            homeStore.addClientsAddress(accountData.value.address, agentAddress, agentPrivateKey)

        } catch (error) {
        }
    }
};

async function handleError(error: any) {
    if (error && error.reason && error.reason == "rejected") {
        await disconnect();
        homeStore.clearClientsAddress()
    }
}
watchEffect(async () => {
    if (accountData.value && accountData.value.isConnected && accountData.value.address) {
        homeStore.initTokens(current_network)
        accountStore.initSubscribeClient(current_network, accountData.value.address)
        await queryExchClient()
        // console.log("loginAndRegisterFcm", accountData.value.address)
        // try {
        //     await ApiService.loginAndRegisterFcm({ address: accountData.value.address })
        // } catch (error) {
        //     console.log("loginAndRegisterFcm error", JSON.stringify(error))
        // } finally {
        //     console.log("loginAndRegisterFcm ok", accountData.value.address)
        // }
    }
})
onMounted(async () => {
    await queryExchClient()
})
</script>
<template>
    <div v-if="accountData && accountData.isConnected && accountData.address">
        <BalanceCard />
    </div>
    <div v-else>
        <ConnectButton />
    </div>
</template>
<style scoped lang="scss"></style>