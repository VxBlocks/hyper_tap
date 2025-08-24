import { defineStore } from 'pinia';
import { Storage } from '@/utils/Storage';
import { BalanceViewData, ClearinghouseInfo, FillsState, FopenOrderState, SpotInfo } from '../types';
import { ACCOUNT_CLEARING_HOUSE_INFO, ACCOUNT_FILLS, ACCOUNT_FOPEN_ORDERS, ACCOUNT_SPOT_INFO } from '../mutation-types';
import * as hl from "@nktkas/hyperliquid";
import { store } from '@/store';
import { SpotBalance } from '@nktkas/hyperliquid';
import { BigNumberUtils } from '@/utils/big-number-utils';
import { ExchangeTokens } from '@/sdk/price';
import { EXCHANGE_TOKENS } from '../mutation-types';

interface AccountState {
    spotInfo: SpotInfo,
    clearinghouseInfo: ClearinghouseInfo,
    fillsState: FillsState,
    fopenOrders: FopenOrderState[]
    loading: boolean
    subscribe_client: hl.SubscriptionClient | null
    spotAssetCtxs: hl.SpotAssetCtx[] | null // 现货
    assetCtxs: hl.PerpsAssetCtx[] | null // 永续  
    balanceViewDatas: BalanceViewData[]
    wsData: hl.WsWebData2 | null
}

export const useAccountStore = defineStore("account", {
    state: (): AccountState => ({
        spotInfo: null,
        clearinghouseInfo: null,
        fillsState: null,
        fopenOrders: [],
        loading: false,
        subscribe_client: null,
        spotAssetCtxs: [],
        assetCtxs: [],
        balanceViewDatas: [],
        wsData: null
    }),
    getters: {
        getWsData(): hl.WsWebData2 | null {
            return this.wsData
        },
        getSpotInfo(): SpotInfo {
            return this.spotInfo || Storage.get(ACCOUNT_SPOT_INFO, '') || {}
        },
        getClearinghouseInfo(): ClearinghouseInfo {
            return this.clearinghouseInfo || Storage.get(ACCOUNT_CLEARING_HOUSE_INFO, '') || {}
        },
        getFillsState(): FillsState {
            return this.fillsState || Storage.get(ACCOUNT_FILLS, '') || {}
        },
        getFopenOrderState(): FopenOrderState[] {
            return this.fopenOrders || Storage.get(ACCOUNT_FOPEN_ORDERS, '') || []
        },
        getSpotAssetCtxs(): hl.SpotAssetCtx[] {
            return this.spotAssetCtxs || []
        },
        getAssetCtxs(): hl.PerpsAssetCtx[] {
            return this.assetCtxs || []
        },
        getBalanceViewDatas(): BalanceViewData[] {
            return this.balanceViewDatas || []
        }
    },
    actions: {
        async initSubscribeClient(network: "testnet" | "mainnet", address: string) {
            console.log("initSubscribeClient address ::: ", address);
            try {
                let isTestnet = (network === 'testnet');
                const transport = new hl.HttpTransport({
                    isTestnet,
                });
                const subTransport = new hl.WebSocketTransport({
                    url: isTestnet ? "wss://api-ui.hyperliquid-testnet.xyz/ws" : "wss://api-ui.hyperliquid.xyz/ws",
                });
                this.client = new hl.InfoClient({
                    transport,
                });
                this.subscribe_client = new hl.SubscriptionClient({
                    transport: subTransport,
                });
                await this.subscribe_client.webData2({ user: address }, (data: hl.WsWebData2) => {
                    this.spotInfo = data?.spotState
                    this.clearinghouseInfo = data?.clearinghouseState
                    this.fopenOrders = data?.openOrders;
                    this.assetCtxs = data?.assetCtxs;
                    this.spotAssetCtxs = data?.spotAssetCtxs;
                    this.balanceViewDatas = this.handleBalanceViewDatas(data)
                    this.wsData = data
                });
                await this.subscribe_client.userFills({ user: address }, (data: hl.WsWebData2) => {
                    this.fillsState = data;
                });
            } catch (error) {
                console.error("initSubscribeClient error", error);
            }
        },
        handleBalanceViewDatas(data: hl.WsWebData2) {
            console.log("handleBalanceViewDatas::::", data);
            let balanceViewDatas: BalanceViewData[] = []
            if (!data || !data.spotState) {
                return balanceViewDatas
            }
            let balances = data.spotState.balances ?? []
            balances.length > 0 && balances.map((item: SpotBalance) => {
                balanceViewDatas.push({
                    currency: item.coin,
                    token: item.token,
                    type: "SPOT",
                    usdcValue: "0",
                    pnl: "",
                    contract: "--",
                    entryNtl: item.entryNtl,
                    availableBalance: BigNumberUtils.subtract(item.total, item.hold),
                    totalBalance: item.total
                })
            })
            return balanceViewDatas
        },
    },
});



// Need to be used outside the setup
export function useAccountStoreWithOut() {
    return useAccountStore(store)
} 