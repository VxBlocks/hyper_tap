import { ExchangeTokens } from "@/sdk/price"; 
import { BigNumberUtils } from "@/utils/big-number-utils";
export async function get_best_rise_tokens(tokens: ExchangeTokens) {
    return tokens
        .tokens_order_by_change_percent()
        .reverse()
        .slice(0, 10)
        .map((item) => {
            return {
                symbol: item.name,
                pair: item.pair ? item.pair : "USD",
                leverage: item.leverage,
                price: item.price,
                volume:item.volume.toFixed(2),
                change: item.daily_change_percent.toFixed(2),
            };
        });
} 
export async function get_worst_rise_tokens(tokens: ExchangeTokens) {
    return tokens
        .tokens_order_by_change_percent()
        .slice(0, 10)
        .map((item) => {
            return {
                symbol: item.name,
                pair: item.pair ? item.pair : "USD",
                leverage: item.leverage,
                price: item.price,
                volume:item.volume.toFixed(2),
                change: item.daily_change_percent.toFixed(2),
            };
        });
}

export async function get_hot_tokens(tokens: ExchangeTokens) {
    return tokens
        .token_order_open_interest()
        .reverse()
        .slice(0, 10)
        .map((item) => {
            return {
                symbol: item.name,
                pair: item.pair ? item.pair : "USD",
                leverage: item.leverage,
                price: item.price,
                volume:item.volume.toFixed(2),
                change: item.daily_change_percent.toFixed(2),
            };
        });
}

export async function get_chosed_tokens(tokens: ExchangeTokens) {
    return tokens
        .get_token_price([5, 1, 159, 5, 200, 25, 165, 122, 123, 124, 123, 124, 123, 124])
        .map((item) => {
            return {
                symbol: item.name,
                pair: item.pair ? item.pair : "USD",
                leverage: item.leverage,
                price: item.price,
                volume:item.volume.toFixed(2),
                change: item.daily_change_percent.toFixed(2),
            };
        });
}

export async function get_day_vlm(tokens: ExchangeTokens) {
    return tokens
        .token_order_day_ntl_vlm()
        .reverse()
        .slice(0, 10)
        .map((item) => {
            return {
                symbol: item.name,
                pair: item.pair ? item.pair : "USD",
                leverage: item.leverage,
                price: item.price, 
                volume:item.volume.toFixed(2),
                change: item.daily_change_percent.toFixed(2),
            };
        });
}

export async function get_new_tokens(tokens: ExchangeTokens) {
    return tokens
        .token_order_index()
        .reverse()
        .slice(0, 20)
        .map((item) => {
            return {
                symbol: item.name,
                pair: item.pair ? item.pair : "USDC",
                leverage: item.leverage,
                price: item.price,
                volume:item.volume.toFixed(2),
                change: item.daily_change_percent.toFixed(2),
            };
        });
}
