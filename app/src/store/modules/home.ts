import { defineStore } from 'pinia';
import { store } from '@/store';
import { ExchangeTokens } from '@/sdk/price';
import { Storage } from '@/utils/Storage';
import { CLIENT_ADDRESS, EXCHANGE_TOKEN_LIST, EXCHANGE_TOKENS } from '../mutation-types';
import { ClientAddress, ExchangeTokenItem } from '../types';
import { HomeTabEnum } from '@/enums/homeTabEnum';
import { get_best_rise_tokens, get_chosed_tokens, get_day_vlm, get_hot_tokens, get_new_tokens, get_worst_rise_tokens } from '@/sdk/utils';

interface HomeState {
    loading: boolean;
    tokens: ExchangeTokens
    loadingToken: Promise<void> | null,
    tokenList: ExchangeTokenItem[]
    clients: ClientAddress[]
}

export const useHomeStore = defineStore("home", {
    state: (): HomeState => ({
        loading: true,
        tokens: undefined,
        loadingToken: null,
        tokenList: [],
        clients: []
    }),
    getters: {
        getTokens(): ExchangeTokens {
            return this.tokens || Storage.get(EXCHANGE_TOKENS, '') || {}
        },
        getTokenList(): ExchangeTokenItem[] {
            return this.tokenList || Storage.get(EXCHANGE_TOKEN_LIST, '') || {}
        },
        getClientsAddress(): ClientAddress[] {
            return this.walletAddress || Storage.get(CLIENT_ADDRESS, '') || {}
        },
        hasAgentddress(): boolean {
            let clients = Storage.get(CLIENT_ADDRESS)
            if (clients && clients.length > 0) {
                return true
            }
            return false
        },
        getLoading(): boolean {
            return this.loading
        },
    },
    actions: {
        addClientsAddress(address: string, agentAddress: string, agentPrivateKey: string) {
            let hasAddress = this.currentAddressHasAgentddress(address)
            if (!hasAddress) {
                let newClientsAddress = { walletAddress: address, agentAddress, agentPrivateKey }
                let newClients = [newClientsAddress, ...this.clients]
                Storage.set(CLIENT_ADDRESS, newClients)
            }
        },
        clearClientsAddress() {
            return Storage.remove(CLIENT_ADDRESS)
        },

        addSession(session: string) {
            Storage.set('session', session)
        },

        getSessionHeader() {
            return {
                'Authorization': Storage.get('session', '')
            }
        },

        getSession(): string {
            return Storage.get('session', '')
        },

        currentAddressHasAgentddress(address: string): boolean {
            let clients = Storage.get(CLIENT_ADDRESS)
            if (clients && clients.length > 0) {
                let findClient = clients.find(item => item.walletAddress === address)
                if (findClient && findClient.agentAddress && findClient.agentPrivateKey) {
                    return true
                }
            }
            return false
        },
        getCurrentClientAddress(address: string): ClientAddress {
            let clients = Storage.get(CLIENT_ADDRESS)
            if (clients && clients.length > 0) {
                let findClient = clients.find(item => item.walletAddress === address)
                if (findClient && findClient.agentAddress && findClient.agentPrivateKey) {
                    return findClient
                }
            }
            return null
        },
        setLoading(loading: boolean) {
            this.loading = loading
        },
        async initTokens(network: "testnet" | "mainnet") {
            if (!this.tokens) {
                const tokens = new ExchangeTokens(network);
                this.loadingToken = tokens.get_perps()
                this.tokens = tokens
            }
            if (this.loadingToken != null) {
                await this.loadingToken
                this.loadingToken = null
            }
        },
        async queryTokenListByType(tab: HomeTabEnum, network: "testnet" | "mainnet") {
            this.loading = true
            try {
                if (!this.tokens) {
                    const tokens = new ExchangeTokens(network);
                    this.loadingToken = tokens.get_perps()
                    this.tokens = tokens
                }
                if (this.loadingToken != null) {
                    await this.loadingToken
                    this.loadingToken = null
                }
                switch (tab) {
                    case HomeTabEnum.Watchlist:
                        this.tokenList = await get_chosed_tokens(this.tokens)
                        break;
                    case HomeTabEnum.Hot:
                        this.tokenList = await get_hot_tokens(this.tokens)
                        break;
                    case HomeTabEnum.New:
                        this.tokenList = await get_new_tokens(this.tokens)
                        break;
                    case HomeTabEnum.Gainers:
                        this.tokenList = await get_best_rise_tokens(this.tokens)
                        break;
                    case HomeTabEnum.Losers:
                        this.tokenList = await get_worst_rise_tokens(this.tokens)
                        break;
                    case HomeTabEnum.Volume:
                        this.tokenList = await get_day_vlm(this.tokens)
                        break;
                    default:
                        this.tokenList = await get_chosed_tokens(this.tokens)
                }
            } catch (error) {
                this.tokenList = []
                console.log("get_tokens error", error);
            } finally {
                this.loading = false
            }

        },
    },
});
export function useHomeStoreWithOut() {
    return useHomeStore(store)
}