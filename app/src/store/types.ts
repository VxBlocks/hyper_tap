import * as hl from "@nktkas/hyperliquid"

/**
 * 现货账户信息
 */
export interface SpotInfo extends hl.SpotClearinghouseState {

}

/**
 * 永续账户信息
 */
export interface ClearinghouseInfo extends hl.PerpsClearinghouseState {

}

export interface AccountBalanceData {
    spotInfo: SpotInfo,
    clearinghouseInfo: ClearinghouseInfo,
}

/**
 * 已完成订单
 */
export interface FillsState extends hl.WsUserFills {

}

/**
 * 未完成订单
 */
export interface FopenOrderState extends hl.FrontendOrder {

}


export interface ExchangeTokenItem {
    symbol: string;
    pair: string;
    leverage: number;
    price: string;
    change: string;
    index: number;
}

export interface ClientAddress {
    walletAddress: string;
    agentAddress: string;
    agentPrivateKey: string;
}

export interface AgentData { 
    agentAddress: string;
    agentPrivateKey: string;
    validUntil:number
}

export interface ShowModalState {
    index: number
    showModal: boolean
    amount: string // 数量
    price?: string
}

export interface ShowSaleModal {
    index: number
    showModal: boolean 
}

export interface BalanceViewData {
    currency: string, // 币种 
    token:number,
    type: string // 类型
    usdcValue:string // USDC价值
    pnl:string // 盈亏
    contract:string // 合约
    availableBalance:string // 可用余额
    totalBalance:string // 总余额
    entryNtl:string
    
}