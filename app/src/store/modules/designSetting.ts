import { defineStore } from 'pinia';
import { store } from '@/store';
import designSetting, { DesignSettingState } from '@/settings/designSetting';
const { darkMode, appTheme, appThemeList, isPageAnimate, pageAnimateType } = designSetting


export const useDesignSettingStore = defineStore("app-design-setting", {
    state: (): DesignSettingState => ({
        darkMode,
        appTheme,
        appThemeList,
        isPageAnimate,
        pageAnimateType,
    }),
    getters: {
        getDarkMode(): 'light' | 'dark' {
            return this.darkMode
        },
        getAppTheme(): string {
            return this.appTheme
        },
        getAppThemeList(): string[] {
            return this.appThemeList
        },
        getIsPageAnimate(): boolean {
            return this.isPageAnimate
        },
        getPageAnimateType(): string {
            return this.pageAnimateType
        },
    },
    actions: {
        setDarkMode(mode: 'light' | 'dark'): void {
            this.darkMode = mode
        },
        setPageAnimateType(type: string): void {
            this.pageAnimateType = type
        },
    },
});

// Need to be used outside the setup
export function useDesignSettingWithOut() {
  return useDesignSettingStore(store)
}
