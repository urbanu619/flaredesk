export const staticRouter = [
  { path: "/", redirect: "/home/index" },
  { path: "/login", name: "login", component: () => import("@/views/login/index.vue"), meta: { title: "登录" } },
  {
    path: "/layout",
    name: "layout",
    component: () => import("@/layouts/index.vue"),
    redirect: "/home/index",
    children: [
      {
        path: "/home/index",
        name: "home",
        component: () => import("@/views/home/index.vue"),
        meta: {
          icon: "HomeFilled",
          isAffix: false,
          isFull: false,
          isHide: false,
          isKeepAlive: true,
          isLink: false,
          title: "首页"
        }
      },
      // ==================== 红包管理 ====================
      {
        path: "/redpacket",
        name: "redpacket",
        meta: {
          icon: "Wallet",
          isAffix: false,
          isFull: false,
          isHide: false,
          isKeepAlive: true,
          isLink: false,
          title: "红包管理"
        },
        children: [
          {
            path: "config",
            name: "redpacketConfig",
            component: () => import("@/views/redpacket/config/index.vue"),
            meta: { title: "红包配置", icon: "Menu", isKeepAlive: true }
          },
          {
            path: "send",
            name: "redpacketSend",
            component: () => import("@/views/redpacket/send/index.vue"),
            meta: { title: "手动发送", icon: "Menu", isKeepAlive: true }
          },
          {
            path: "record",
            name: "redpacketRecord",
            component: () => import("@/views/redpacket/record/index.vue"),
            meta: { title: "红包记录", icon: "Menu", isKeepAlive: true }
          },
          {
            path: "group",
            name: "telegramGroup",
            component: () => import("@/views/redpacket/group/index.vue"),
            meta: { title: "群组管理", icon: "Menu", isKeepAlive: true }
          },
          {
            path: "user-bind",
            name: "userBind",
            component: () => import("@/views/redpacket/user-bind/index.vue"),
            meta: { title: "用户绑定", icon: "Menu", isKeepAlive: true }
          },
          {
            path: "user-asset",
            name: "userAsset",
            component: () => import("@/views/redpacket/user-asset/index.vue"),
            meta: { title: "用户资产", icon: "Menu", isKeepAlive: true }
          }
        ]
      },
      // ==================== Cloudflare 管理 ====================
      {
        path: "/cloudflare",
        name: "cloudflare",
        meta: {
          icon: "Monitor",
          isAffix: false,
          isFull: false,
          isHide: false,
          isKeepAlive: true,
          isLink: false,
          title: "Cloudflare"
        },
        children: [
          {
            path: "account",
            name: "cfAccount",
            component: () => import("@/views/cloudflare/account/index.vue"),
            meta: { title: "账号管理", icon: "Key", isKeepAlive: true }
          },
          {
            path: "zones",
            name: "cfZones",
            component: () => import("@/views/cloudflare/zones/index.vue"),
            meta: { title: "Zone 列表", icon: "List", isKeepAlive: true }
          },
          {
            path: "dns",
            name: "cfDns",
            component: () => import("@/views/cloudflare/dns/index.vue"),
            meta: { title: "DNS 管理", icon: "Connection", isKeepAlive: true }
          },
          {
            path: "template",
            name: "cfDnsTemplate",
            component: () => import("@/views/cloudflare/template/index.vue"),
            meta: { title: "DNS 模板", icon: "Document", isKeepAlive: true }
          }
        ]
      },
      // ==================== 系统配置 ====================
      {
        path: "/system",
        name: "system",
        meta: {
          icon: "Setting",
          isAffix: false,
          isFull: false,
          isHide: false,
          isKeepAlive: true,
          isLink: false,
          title: "系统配置"
        },
        children: [
          {
            path: "accountManage",
            name: "accountManage",
            component: () => import("@/views/system/accountManage/index.vue"),
            meta: {
              title: "账号管理",
              icon: "Menu",
              isHide: false,
              isFull: false,
              isAffix: false,
              isKeepAlive: true
            }
          },
          {
            name: "角色管理",
            path: "roleManage",
            component: () => import("@/views/system/roleManage/index.vue"),
            meta: {
              title: "角色管理",
              icon: "Menu",
              isHide: false,
              isFull: false,
              isAffix: false,
              isKeepAlive: true
            },
            children: []
          },
          {
            name: "菜单管理",
            path: "menuMange",
            component: () => import("@/views/system/menuMange/index.vue"),
            meta: {
              title: "菜单管理",
              icon: "Menu",
              isHide: false,
              isFull: false,
              isAffix: false,
              isKeepAlive: true
            },
            children: []
          },
          {
            name: "apis管理",
            path: "apisManage",
            component: () => import("@/views/system/apisManage/index.vue"),
            meta: {
              title: "apis管理",
              icon: "Menu",
              isHide: false,
              isFull: false,
              isAffix: false,
              isKeepAlive: true
            },
            children: []
          }
        ]
      }
    ]
  },
  {
    path: "/404",
    name: "404",
    component: () => import("@/components/ErrorMessage/404.vue"),
    meta: { title: "404页面" }
  },
  {
    path: "/500",
    name: "500",
    component: () => import("@/components/ErrorMessage/500.vue"),
    meta: { title: "500页面" }
  },
  // Resolve refresh page, route warnings
  { path: "/:pathMatch(.*)*", component: () => import("@/components/ErrorMessage/404.vue") }
];
