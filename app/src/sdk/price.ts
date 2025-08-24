import * as hl from "@nktkas/hyperliquid"
import Decimal from "decimal.js";

export interface ExchangeTokenMeta extends ExchangePrice {
    name: string;
    pair?: string;
    asset_index: number;
    leverage: number;
    volume: Decimal;
}

export interface ExchangePrice {
    price: string;
    mid_price: string;
    // 现货名称 HYPE/USDC或者@<coin_index>
    asset_name?: string;
    daily_change_percent: Decimal;
    daily_change: Decimal;
}

export class ExchangeTokens {
    private client: hl.InfoClient;

    // asset id 对应，用于交易
    coin_to_asset = new Map<string, number>();
    asset_to_coin = new Map<number, string>();

    // 属性 key为asset_index
    coin_to_perps_asset_ctx = new Map<number, hl.PerpsAssetCtx>();
    coin_to_spot_asset_ctx = new Map<number, hl.SpotAssetCtx>();
    coin_to_spot_token = new Map<number, hl.SpotUniverse>();
    coin_to_perps_token = new Map<number, hl.PerpsUniverse>();
    coin_to_spot_metadata = new Map<number, hl.SpotToken>();
    coin_to_perps_token_meta = new Map<number, hl.PerpsMeta>();

    // 倍数
    margintable = new Map<number, hl.MarginTable>()

    // 价格
    spot_coin_to_price = new Map<number, ExchangeTokenMeta>();
    perps_coin_to_price = new Map<number, ExchangeTokenMeta>();

    // 现货tokens 对应名称 key为token index / coin index
    coin_index_to_name = new Map<number, string>();
    coin_index_to_price = new Map<number, ExchangeTokenMeta>()

    constructor(network: "testnet" | "mainnet") {
        const transport = new hl.HttpTransport({
            isTestnet: network === "testnet",
        });
        this.client = new hl.InfoClient({
            transport,
        });

        //TODO: load storage


        // this.get_perps(); 
    }

    async add_pair_and_name_to_index_map(
    ) {
        const [meta, assetCtxs] = await this.client.spotMetaAndAssetCtxs()
        this.coin_index_to_name = new Map<number, string>();
        for (const token of meta.tokens) {
            this.coin_index_to_name.set(token.index, token.name);
            this.coin_to_spot_metadata.set(token.index, token);
        }

        for (const [index, asset] of meta.universe.entries()) {
            const spot_ind = 10000 + asset.index;

            // let name_to_ind = [asset.name, spot_ind]; 

            let token_1_name = this.coin_index_to_name.get(asset.tokens[0])
            let token_2_name = this.coin_index_to_name.get(asset.tokens[1])
            if (token_1_name == undefined || token_2_name == undefined) {
                continue;
            }

            const spot_name = `${token_1_name}/${token_2_name}`;
            this.coin_to_asset.set(spot_name, spot_ind);
            this.coin_to_asset.set(asset.name, spot_ind);

            this.asset_to_coin.set(spot_ind, spot_name);

            this.coin_to_spot_token.set(spot_ind, asset)

            const asset_ctx = assetCtxs.find((item) => item.coin == asset.name);
            this.coin_to_spot_asset_ctx.set(spot_ind, asset_ctx);
        }
        console.log(this);
        // this.client.extraAgents({ user: "0x5263ABaa3dd77dDD0870CdbCCE78e85d82Ea4c0c" }).then((value) =>
        //     console.log(value));
    }

    async get_perps() {
        const [meta, assetCtxs] = await this.client.metaAndAssetCtxs()
        let coin_to_asset = new Map<string, number>();
        let asset_to_coin = new Map<number, string>();

        let coin_to_perps_asset_ctx = new Map<number, hl.PerpsAssetCtx>();

        let coin_to_perps_token = new Map<number, hl.PerpsUniverse>();

        for (const [id, table] of meta.marginTables) {
            this.margintable.set(id, table)
        }

        for (const [index, coin] of meta.universe.entries()) {
            coin_to_asset.set(coin.name, index)
            asset_to_coin.set(index, coin.name)
            coin_to_perps_asset_ctx.set(index, assetCtxs[index])
            coin_to_perps_token.set(index, coin)
        }

        this.coin_to_spot_asset_ctx = new Map<number, hl.SpotAssetCtx>();
        this.asset_to_coin = asset_to_coin;
        this.coin_to_asset = coin_to_asset;


        await this.add_pair_and_name_to_index_map();

        this.coin_to_perps_token = coin_to_perps_token;
        this.coin_to_perps_asset_ctx = coin_to_perps_asset_ctx;

        this.spot_coin_to_price = new Map<number, ExchangeTokenMeta>();

        for (const [asset_ind, asset_ctx] of this.coin_to_spot_asset_ctx.entries()) {
            const price = this.price_from_shared_asset(asset_ctx);
            const coin = this.coin_to_spot_token.get(asset_ind)!;
            const coin1 = coin.tokens[0];
            const coin2 = coin.tokens[1];

            const coin1_name = this.coin_index_to_name.get(coin1)!;
            const coin2_name = this.coin_index_to_name.get(coin2)!;

            let supply1 = new Decimal(asset_ctx.dayNtlVlm) 
            const price_and_meta = {
                ...price,
                name: coin1_name,
                pair: coin2_name,
                asset_index: asset_ind,
                leverage: 1,
                volume: supply1
            }
            this.spot_coin_to_price.set(asset_ind, price_and_meta);
            if (coin2_name == "USDC") {
                this.coin_index_to_price.set(coin1, price_and_meta)
            }
        }

        this.perps_coin_to_price = new Map();
        for (const [asset_ind, asset_ctx] of coin_to_perps_asset_ctx.entries()) {
            const price = this.price_from_shared_asset(asset_ctx);
            let supply1 = new Decimal(asset_ctx.dayNtlVlm) 
            const price_and_meta = {
                ...price,
                name: this.asset_to_coin.get(asset_ind)!,
                asset_index: asset_ind,
                leverage: this.coin_to_perps_token.get(asset_ind)!.maxLeverage,
                volume: supply1
            }
            this.perps_coin_to_price.set(asset_ind, price_and_meta);
        }
    }
    price_from_shared_asset(asset_ctx: hl.SharedAssetCtx): ExchangePrice {
        const price = new Decimal(asset_ctx.markPx);
        const prevDayPx = new Decimal(asset_ctx.prevDayPx);
        const change = price.minus(prevDayPx);
        let changePct = change.div(prevDayPx).mul(100);
        if (prevDayPx.isZero()) {
            changePct = new Decimal(100);
        }

        return {
            price: asset_ctx.markPx,
            mid_price: asset_ctx.midPx,
            // @ts-ignore
            asset_name: asset_ctx.coin,
            daily_change_percent: changePct,
            daily_change: change,
        };
    };

    get_token_price(tokens: number[]): ExchangeTokenMeta[] {
        let res: ExchangeTokenMeta[] = [];

        for (const index of tokens) {
            let price = this.spot_coin_to_price.get(index);
            if (price) {
                res.push(price);
            }
            price = this.perps_coin_to_price.get(index);
            if (price) {
                res.push(price);
            }
        }

        return res;
    }

    tokens_order_by_change_percent(): ExchangeTokenMeta[] {

        let prices = Array.from(this.perps_coin_to_price.values());

        let sorted = prices.sort((a, b) => a.daily_change_percent.comparedTo(b.daily_change_percent));

        return sorted
    }

    // 订单量
    token_order_open_interest(): ExchangeTokenMeta[] {
        let metadatas = Array.from(this.coin_to_perps_asset_ctx.entries());

        let sorted = metadatas.sort((a, b) => {
            let supply1 = new Decimal(a[1].openInterest)
            let price1 = new Decimal(a[1].markPx)
            let supply2 = new Decimal(b[1].openInterest)
            let price2 = new Decimal(b[1].markPx)

            let usdc1 = supply1.mul(price1)
            let usdc2 = supply2.mul(price2)
            return usdc1.comparedTo(usdc2)
        });

        let prices = sorted.map((a) => {
            return this.perps_coin_to_price.get(a[0])!;
        })

        return prices
    }

    //当日交易量
    token_order_day_ntl_vlm(): ExchangeTokenMeta[] {
        let metadatas = Array.from(this.coin_to_perps_asset_ctx.entries());

        let sorted = metadatas.sort((a, b) => {
            let supply1 = new Decimal(a[1].dayNtlVlm) 
            let supply2 = new Decimal(b[1].dayNtlVlm)  
            return supply1.comparedTo(supply2)
        }); 
        let prices = sorted.map((a) => {
            return this.perps_coin_to_price.get(a[0])!;
        })

        return prices
    }

    token_order_index(): ExchangeTokenMeta[] {
        let prices = Array.from(this.perps_coin_to_price.entries());

        let sorted = prices.sort((a, b) => Decimal(a[0]).comparedTo(Decimal(b[0])));
        return sorted.map(x => x[1])
    }
}

