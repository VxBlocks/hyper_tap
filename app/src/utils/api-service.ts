import { LoginService } from "@/sdk/backend-sdk/login/v1/login_pb";
import createClient from "@/sdk/transport";
import { Provider, useAppKitProvider } from "@reown/appkit/vue";
import { BrowserProvider } from "ethers";
import { getFcmToken } from "./android-channel";

import { Storage } from '@/utils/Storage';

export class ApiService {

    static async loginAndRegisterFcm({ address }: { address: string }) {
        const loginClient = createClient(LoginService)

        if (await this.shouldlogin({ address })) {
            const provider = useAppKitProvider<Provider>("eip155");
            const browserProvider = new BrowserProvider(provider.walletProvider!);
            const msg = "Hello Hyperliquid"
            const signature = await (await browserProvider.getSigner()).signMessage(msg)
            const loginResp = await loginClient.login({
                address,
                msg,
                signature
            })
            const session = loginResp.session;
            Storage.set('session', session)
        }
        if (getFcmToken()) {
            await loginClient.registerFcm({
                fcmToken: getFcmToken()
            }, {
                headers: { 'Authorization': Storage.get('session', '') }
            })
        }
    }

    static async shouldlogin({ address }: { address: string }): Promise<boolean> {
        try {
            const session = Storage.get('session', '')
            if (!session) {
                return true
            }
            const loginClient = createClient(LoginService)
            const loginResp = await loginClient.loginIsValid({
                session,
                address,
            })
            return !loginResp.valid
        } catch (error) {
            return true
        }
    }
}