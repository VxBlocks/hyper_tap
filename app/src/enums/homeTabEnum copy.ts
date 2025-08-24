import { useI18n } from "vue-i18n";

// src/enums/HomeTabKey.ts
export enum HomeTabKey {
  Watchlist = "Watchlist",
  Hot = "Hot",
  New = "New",
  Gainers = "Gainers",
  Losers = "Losers",
  Volume = "Volume"
}

export function useHomeTabs() {
  const { t } = useI18n();
  
  const homeTabList = [
    { key: HomeTabKey.Watchlist, label: t("tabs.favorites") },
    { key: HomeTabKey.Hot, label: t("tabs.hot") },
    { key: HomeTabKey.New, label: t("tabs.new") },
    { key: HomeTabKey.Gainers, label: t("tabs.gainers") },
    { key: HomeTabKey.Losers, label: t("tabs.losers") },
    { key: HomeTabKey.Volume, label: t("tabs.volume") }
  ];
  
  return { homeTabList, HomeTabKey };
}