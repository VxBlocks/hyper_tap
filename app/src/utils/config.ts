
import { current_network } from '@/constant'
import { WagmiAdapter } from '@reown/appkit-adapter-wagmi'
import { arbitrumSepolia, arbitrum, type AppKitNetwork } from '@reown/appkit/networks'

export const projectId = "3b44d83935daeb453a39630ed3e52ebd" // this is a public projectId only to use on localhost
if (!projectId) {
    throw new Error('VITE_PROJECT_ID is not set')
}

export const networks: [AppKitNetwork, ...AppKitNetwork[]] = [current_network == "mainnet" ? arbitrum : arbitrumSepolia]

export const wagmiAdapter = new WagmiAdapter({
    networks,
    projectId
})
export const config = wagmiAdapter.wagmiConfig