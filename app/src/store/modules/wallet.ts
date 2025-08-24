import { defineStore } from 'pinia';
import { Storage } from '@/utils/Storage';
import { ACCOUNT_SPOT_INFO, WALLET_CURRENT_ADDRESS } from '../mutation-types';

interface WalletState {
  currentAddress: string;
}

export const useWalletStore = defineStore("wallet", {
  state: (): WalletState => ({
    currentAddress: '',
  }),
  getters: {
    getSpotInfo(): string {
      return this.currentAddress || Storage.get(WALLET_CURRENT_ADDRESS, '') || {}
    },
  },
  actions: {
  },
});
