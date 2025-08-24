import { ExchangeTokens } from "@/sdk/price"
import { BigNumberUtils } from "./big-number-utils"
import * as hl from "@nktkas/hyperliquid";

export class AccountVault {
    static getAccountVaultBalance(data: hl.WsWebData2, tokens: ExchangeTokens) {
        const perps_balance = data.clearinghouseState?.marginSummary.accountValue || 0

        const spot_balances = data && data.spotState && data.spotState.balances ? data.spotState.balances : []

        const spotAssetCtxs = data.spotAssetCtxs

        let spot_balance = "0"
        for (const balance of spot_balances) {
            const price = this.queryLatestPrice(balance.token, tokens, spotAssetCtxs)
            const usdcbalance = BigNumberUtils.multiply(balance.total, price)
            spot_balance = BigNumberUtils.add(spot_balance, usdcbalance)
        }

        return BigNumberUtils.add(perps_balance, spot_balance)
    }

    static queryLatestPrice(token: number, tokens: ExchangeTokens, spotAssetCtxs: hl.SpotAssetCtx[]) {
        let latestPrice = "0"
        if (!tokens.coin_index_to_price) {
            return latestPrice
        }
        let priceData = tokens.coin_index_to_price.get(token)
        if (!(priceData && priceData.asset_name)) {
            return latestPrice
        }
        let coin = priceData.asset_name
        let currentAssetCtxs = spotAssetCtxs.find((item) => item.coin === coin)
        if (currentAssetCtxs) {
            latestPrice = currentAssetCtxs.markPx
        }
        return latestPrice
    }
}