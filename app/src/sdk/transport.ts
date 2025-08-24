import { DescService } from "@bufbuild/protobuf";
import { Client, createClient as connectCreateClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

const transport = createConnectTransport({
    baseUrl: "https://hl.coolrc.me",
});

export default function createClient<T extends DescService>(service: T): Client<T> {
    return connectCreateClient(service, transport);
}
