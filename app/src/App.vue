<template>
  <router-view #="{ Component }">
    <component :is="Component" />
  </router-view>
</template>

<script setup lang="ts">

import {
  createAppKit,
} from '@reown/appkit/vue'
import { wagmiAdapter, networks, projectId } from '@/utils/config'
import { arbitrumSepolia, arbitrum } from '@reown/appkit/networks'
import { onMounted } from 'vue'
import { openUrl } from '@tauri-apps/plugin-opener'
import { current_network } from './constant'
import { useHomeStore } from './store/modules/home'

const homeStore = useHomeStore()
function init() {
  createAppKit({
    adapters: [wagmiAdapter],
    networks,
    projectId,
    defaultNetwork: current_network == "mainnet" ? arbitrum : arbitrumSepolia,
    themeMode: 'light',
    metadata: {
      name: 'AppKit Vue Example',
      description: 'AppKit Vue Example',
      url: 'https://reown.com/appkit',
      icons: ['https://avatars.githubusercontent.com/u/179229932?s=200&v=4']
    },
    themeVariables: {
      '--w3m-accent': '#000000',
    }
  })
}
function initListenter() {
  const temp = window.open;
  window.open = (url: string, target, features) => {
    if (url) {
      openUrl(url);
      return null
    } else {
      return temp(url, target, features)
    }
  };
}
onMounted(async () => {
  initListenter()
  init()
  homeStore.initTokens(current_network)
})

</script>

<style lang="less" scoped> 
  @import './styles/index.less';
</style>
