import { useI18n } from "vue-i18n";

export enum AccountTabEnum {
  Positions = "Positions",
  TWAP = "TWAP",
  Spot = "Spot",
  PendingOrders = "PendingOrders",
  ExecutedOrders = "ExecutedOrders",
  OrderHistory = "OrderHistory",
}

export function useAccountTabs() {
  const { t } = useI18n();

  const accountTabList = [
    { key: AccountTabEnum.Positions, label: t("my-tabs.positions") },
    { key: AccountTabEnum.Spot, label: t("my-tabs.spot") },
    { key: AccountTabEnum.PendingOrders, label: t("my-tabs.open-orders") },
    { key: AccountTabEnum.ExecutedOrders, label: t("my-tabs.trade-history") }, 
  ];

  return { accountTabList, AccountTabEnum };
}