import { createApp } from "vue";

import i18n from '@/utils/i18n/config.js';
import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query';
import { Locale } from 'vant'; 
import 'vant/lib/index.css';
import enUS from 'vant/es/locale/lang/en-US';
import App from "./App.vue";
import { setupPlugins } from './plugins';
import { setupRouter } from "./router";
import { setupStore } from './store';
import { initAndroidChannel } from "./utils/android-channel";
const queryClient = new QueryClient()
Locale.use('en-US', enUS);


initAndroidChannel()

const app = createApp(App);

setupPlugins(app);

// 安装pina store
setupStore(app);
// 安装router
setupRouter(app);
app.use(i18n)
app.use(VueQueryPlugin, { queryClient })
app.mount('#app');
