import { App } from 'vue';
import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import { createRouterGuards } from './routerGuards';
import { ErrorPageRoute, RootRoute } from './base';
import { useRouteStoreWithOut } from '@/store/modules/route';
import routeModuleList from './modules';

export const constantRouter: RouteRecordRaw[] = [RootRoute, ErrorPageRoute]

const routeStore = useRouteStoreWithOut()

routeStore.setMenus(routeModuleList)
routeStore.setRouters(constantRouter.concat(routeModuleList))

const router = createRouter({
  history: createWebHashHistory(''),
  routes: constantRouter.concat(...routeModuleList),
  strict: true,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})
export function setupRouter(app: App) {
  app.use(router);
  // 创建路由守卫
  createRouterGuards(router);
}
export default router;
