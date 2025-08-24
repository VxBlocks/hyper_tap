import { ExchangeTokens } from "@/sdk/price";
import { BigNumberUtils } from "./big-number-utils";
import * as hl from "@nktkas/hyperliquid";
import { generatePrivateKey, privateKeyToAccount } from "viem/accounts";
import { Provider, useAppKitProvider } from "@reown/appkit/vue";
import { BrowserProvider } from "ethers";
import { AgentData, FopenOrderState } from "@/store/types";
import { AssetPosition } from "@nktkas/hyperliquid";
import { OrderParams } from "@nktkas/hyperliquid/script/src/types/mod";
import { queryObjects } from "v8";

export const ExchangeHelper = {
    /**
     * 
     * @param meta  
     * @param asset 
     * @param isBuy
     * @param price
     * @param slippage
     */
    get_exchange_token_item(tokens: ExchangeTokens, asset: string, price: string, slippage: number, isBuy: boolean) {
        try {
            let coinToAsset = tokens.coin_to_asset;
            let assetID = coinToAsset.get(asset);
            let perpsToken = tokens.coin_to_perps_token.get(assetID);
            let sz_decimals = perpsToken.szDecimals;
            let max_decimals = assetID < 10000 ? 6 : 8;
            let price_decimals = Math.max(max_decimals - sz_decimals)
            let slippage_factor = isBuy ? 1 + slippage : 1 - slippage;
            let mpx = BigNumberUtils.multiply(price, slippage_factor)
            // Round to the correct number of decimal places and significant figures
            let px = BigNumberUtils.roundToSignificantAndDecimal(mpx, 5, price_decimals);
            return { px, sz_decimals }
        } catch (error) {
            throw error;
        }
    },
    /**
     * 
     * @param client 
     * @param coinSymbol 币名
     * @param szi  数量
     * @param tokens 
     * @param price 价格
     * @param slippage 滑点
     * @param isBuy  买/卖
     * @returns 
     */
    async requestExchange(client: hl.ExchangeClient, coinSymbol: string, szi: number, tokens: ExchangeTokens, price: string, slippage: number, isBuy: boolean) {
        try {
            let coinToAsset = tokens.coin_to_asset;
            let assetID = coinToAsset.get(coinSymbol);
            let { px, sz_decimals } = ExchangeHelper.get_exchange_token_item(tokens, coinSymbol, price, slippage, isBuy)
            let size = BigNumberUtils.round(Math.abs(szi), sz_decimals)
            let result = await client.order({
                orders: [{
                    a: assetID, //BTC
                    b: false, // buy
                    p: px,
                    s: size,
                    r: false,
                    t: { limit: { tif: "Gtc" } }
                }],
                grouping: "na",
            });
            return result;
        } catch (error) {
            throw error;
        }
    },
    async cancelOrder(client: hl.ExchangeClient, asset: number, oid: number) {
        try {
            const result = await client.cancel({
                cancels: [
                    { a: asset, o: oid },
                ],
            });
            return result;
        } catch (error) {
            throw error;
        }
    },
    async batchRequestExchangeCancelMessage(
        agentPrivateKey: `0x${string}`,
        openOrders: FopenOrderState[], tokens: ExchangeTokens) {
        try {
            const agentAccount = privateKeyToAccount(agentPrivateKey);
            const transport = new hl.HttpTransport({ isTestnet: false })
            let client = new hl.ExchangeClient({
                wallet: agentAccount,
                transport: transport,
                isTestnet: false,
            });
            let cancelOrders = this.buildCancelOrders(openOrders, tokens);
            const result = await client.cancel({
                cancels: cancelOrders,
            });
            return result;
        } catch (error) {
            throw error;
        }
    },

    buildCancelOrders(openOrders: FopenOrderState[], tokens: ExchangeTokens) {
        let cancelOrders = []
        let coinToAsset = tokens.coin_to_asset;
        openOrders && openOrders.forEach((item, index) => {
            let asset = coinToAsset.get(item.coin);
            let oid = item.oid;
            cancelOrders.push({
                a: asset,
                o: oid,
            })
        })
        return cancelOrders;
    },
    async getExchangeClient({ agentPrivateKey, isTestnet = false }: { agentPrivateKey: `0x${string}`, isTestnet?: boolean }): Promise<hl.ExchangeClient> {
        try {
            const agentAccount = privateKeyToAccount(agentPrivateKey);
            const transport = new hl.HttpTransport({ isTestnet })
            return new hl.ExchangeClient({
                wallet: agentAccount,
                transport: transport,
                isTestnet,
            });
        } catch (error) {
            throw error;
        }
    },
    async initHyperliquid({ isTestnet = false }: { isTestnet?: boolean }): Promise<AgentData> {
        const provider = useAppKitProvider<Provider>("eip155");
        const browserProvider = new BrowserProvider(provider.walletProvider!);
        try {
            const transport = new hl.HttpTransport({ isTestnet })
            const tempExchClient = new hl.ExchangeClient({
                wallet: await browserProvider.getSigner(),
                transport: transport,
                isTestnet,
            });
            await this.registerTechworld(tempExchClient)
            const privateKey = generatePrivateKey();
            const account = privateKeyToAccount(privateKey);
            const agentAddr = account.address;
            const agentName = "MyAgent";
            await tempExchClient.approveAgent({
                agentAddress: agentAddr,
                agentName,
            })
            // let agents = await infoClient.extraAgents({
            //     user: agentAddr,
            // }); 
            return {
                validUntil: 0,
                agentAddress: agentAddr,
                agentPrivateKey: privateKey,
            };
        } catch (error) {
            throw (error)
        }
    },
    async registerTechworld(tempExchClient: any) {
        try {
            await tempExchClient.registerReferrer({
                code: "TECHWORLD"
            })
        } catch (error) {
        }
    },

    async requestExchangeMessage(agentPrivateKey: `0x${string}`, coinSymbol: string, szi: number, tokens: ExchangeTokens, price: string, slippage: number, isBuy: boolean, isLimit: boolean) {
        try {
            const agentAccount = privateKeyToAccount(agentPrivateKey);
            const transport = new hl.HttpTransport({ isTestnet: false })
            let exchangeClient = new hl.ExchangeClient({
                wallet: agentAccount,
                transport: transport,
                isTestnet: false,
            });
            let coinToAsset = tokens.coin_to_asset;
            let assetID = coinToAsset.get(coinSymbol);
            let { px, sz_decimals } = ExchangeHelper.get_exchange_token_item(tokens, coinSymbol, price, slippage, isBuy)
            let size = BigNumberUtils.round(Math.abs(szi), sz_decimals)
            let result = await exchangeClient.order({
                orders: [{
                    a: assetID, //BTC
                    b: isBuy, // buy
                    p: px,
                    s: size,
                    r: true, // 只做空
                    t: { limit: { tif: isLimit ? "Gtc" : "FrontendMarket" } }
                }],
                grouping: "na",
            });
            return result;
        } catch (error) {
            throw error;
        }
    },

    async queryHasAgentAddress(userAddress: `0x${string}`, agentAddress: `0x${string}`) {
        let hasAgentAddress = false
        try {
            const transport = new hl.HttpTransport({ isTestnet: false })
            let client = new hl.InfoClient({
                transport,
            });
            let req = await client.extraAgents({ user: userAddress })
            let findAgent = req.find((item) => item.address === agentAddress)
            if (findAgent) {
                hasAgentAddress = true
            }
        } catch (error) {
            throw error;
        }
        return hasAgentAddress
    },

    async batchRequestExchangeMessage(
        agentPrivateKey: `0x${string}`,
        positions: AssetPosition[],
        tokens: ExchangeTokens,
        perpsAssetCtx: hl.PerpsAssetCtx[]) {
        try {
            const agentAccount = privateKeyToAccount(agentPrivateKey);
            const transport = new hl.HttpTransport({ isTestnet: false })
            let exchangeClient = new hl.ExchangeClient({
                wallet: agentAccount,
                transport: transport,
                isTestnet: false,
            });
            let orders = await this.handleOrders(positions, tokens, perpsAssetCtx)
            let result = await exchangeClient.order({
                orders: orders,
                grouping: "na",
            });
            return result;
        } catch (error) {
            throw error;
        }
    },
    async handleOrders(positions: AssetPosition[], tokens: ExchangeTokens, perpsAssetCtx: hl.PerpsAssetCtx[]) {
        let orders = [] as OrderParams[]
        try {
            let coinToAsset = tokens.coin_to_asset;
            positions.forEach((item) => {
                let assetID = coinToAsset.get(item.position.coin);
                let currentAssetCtxs = perpsAssetCtx[assetID]
                let price = currentAssetCtxs.markPx
                let szi = Number(item.position.szi)
                let isBuy = Number(item.position.szi) < 0 ? true : false;
                let { px, sz_decimals } = ExchangeHelper.get_exchange_token_item(tokens, item.position.coin, price, 0.08, isBuy)
                let size = BigNumberUtils.round(Math.abs(szi), sz_decimals)
                orders.push({
                    a: assetID, //BTC
                    b: isBuy, // buy
                    p: px,
                    s: size,
                    r: true, // 只做空
                    t: { limit: { tif: "FrontendMarket" } }
                })
            })
        } catch (error) {

        }

        return orders
    },


}