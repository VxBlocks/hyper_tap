import { useI18n } from "vue-i18n";

export enum HomeTabEnum {
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
    { key: HomeTabEnum.Watchlist, label: t("tabs.favorites") },
    { key: HomeTabEnum.Hot, label: t("tabs.hot") },
    { key: HomeTabEnum.New, label: t("tabs.new") },
    { key: HomeTabEnum.Gainers, label: t("tabs.gainers") },
    { key: HomeTabEnum.Losers, label: t("tabs.losers") },
    { key: HomeTabEnum.Volume, label: t("tabs.volume") }
  ];
  
  return { homeTabList, HomeTabEnum };
}