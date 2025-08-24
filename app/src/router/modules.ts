import type { RouteRecordRaw } from 'vue-router'

const Layout = () => import('@/layout/index.vue')

const OtherLayout = () => import('@/layout/other-index.vue')

const routeModuleList: Array<RouteRecordRaw> = [
    {
        path: '/home',
        name: 'Home',
        redirect: '/home/index',
        component: Layout,
        meta: {
            title: '首页',
            hiddenHeader: true,
            icon: 'home-o',
        },
        children: [
            {
                path: 'index',
                name: 'HomePage',
                meta: {
                    keepAlive: true,
                },
                component: () => import('@/views/home/index.vue'),
            },
        ],
    },
    {
        path: '/trend',
        name: 'Trend',
        redirect: '/trend/index',
        component: Layout,
        meta: {
            title: '行情',
            hiddenHeader: true,
            icon: 'home-o',
        },
        children: [
            {
                path: 'index',
                name: 'TrendPage',
                meta: {
                    keepAlive: true,
                },
                component: () => import('@/views/trend/index.vue'),
            },
        ],
    },
    {
        path: '/sale',
        name: 'Sale',
        redirect: '/sale/index',
        component: Layout,
        meta: {
            title: '交易',
            hiddenHeader: true,
            icon: 'after-sale',
        },
        children: [
            {
                path: 'index',
                name: 'SalePage',
                meta: {
                    keepAlive: true,
                },
                component: () => import('@/views/sale/index.vue'),
            },
        ],
    },
    {
        path: '/account',
        name: 'Account',
        redirect: '/account/index',
        component: Layout,
        meta: {
            title: '资产',
            hiddenHeader: true,
            icon: 'home-user-circle-o',
        },
        children: [
            {
                path: 'index',
                name: 'AccountPage',
                meta: {
                    keepAlive: true,
                },
                component: () => import('@/views/account/index.vue'),
            },
        ],
    },
    {
        path: '/watchlist',
        name: 'Watch',
        component: () => import('@/views/watchlist/index.vue'),
        meta: {
            title: 'watchlist',
            hiddenHeader: true,
        },
    },
    {
        path: '/message_center',
        name: 'message_center',
        meta: {
            title: 'Message Center',
            hiddenHeader: true,
        },
        redirect: '/message_center/index',
        component: OtherLayout,
        children: [
            {
                path: 'index',
                name: 'MessagePage',
                meta: {
                    keepAlive: true,
                },
                component: () => import('@/views/message/index.vue'),
            },
        ],
    },
]

export default routeModuleList