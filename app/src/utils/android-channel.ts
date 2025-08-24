
export async function initAndroidChannel() {

    // @ts-ignore
    window.handleNextMessage = () => {
        requestWebViewMsg()
    };

    setInterval(async () => {
        try {
            requestWebViewMsg()
        } catch (error) {
            // console.error("Error in notification loop:", error)
        }
    }, 500)
}

export async function requestWebViewMsg() {
    try {
        while (true) {
            let webViewMsg = getNextMsg()
            if (webViewMsg) {
                console.log("WebViewChannel", JSON.stringify(webViewMsg))
                // 处理消息
                handleWebViewMsg(webViewMsg)
            } else {
                break
            }
        }
    } catch (e) {
        console.log(e)
    }
}

interface WebViewMsg {
    id: number
    method: string
    params: string
}

function getNextMsg(): WebViewMsg | null {
    if (!(window as any).WebViewChannel) {
        return null
    }
    let msg = (window as any).WebViewChannel.nextMsg() as string | null
    if (msg) {
        let webViewMsg = JSON.parse(msg) as WebViewMsg
        return webViewMsg
    }
    return null
}

function addResult(id: number, result: any) {
    const json = JSON.stringify(result);
    (window as any).WebViewChannel.addMsg(id, json)
}

// Sleep function to pause execution for a given number of milliseconds
function sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function handleWebViewMsg(webViewMsg: WebViewMsg) {

    console.log("handleAnroidCall", `method = ${webViewMsg.method} params = ${JSON.stringify(webViewMsg.params)}`);

    const params = JSON.parse(webViewMsg.params);

    // 根据消息类型处理不同的逻辑
    switch (webViewMsg.method) {
        case "gotoNotification":
            gotoNotification();
            addResult(webViewMsg.id, true)
            break;
        case "saveFcmToken":
            console.log("保存FCM Token:", params.token);
            saveFcmToken(params.token);
            addResult(webViewMsg.id, true)
            break;
        // 可以添加更多case来处理其他消息类型
        default:
            console.warn("Unknown message method:", webViewMsg.method);
            addResult(webViewMsg.id, "false")
            return null
    }
}

function gotoNotification() {
    window.location.href = window.location.origin + "/#/account/index"
}

export function saveFcmToken(token: string) {
    localStorage.setItem("fcm_token", token)
}

export function getFcmToken() {
    return localStorage.getItem("fcm_token")
}
