import {defineStore} from "pinia";
import {ref} from "vue";
import {getAuthButtonListApi, getAuthMenuListApi} from "@/api/modules/login";
import {getFlatMenuList, getShowMenuList, getAllBreadcrumbList} from "@/utils";
import {getTreeList} from "../../api/modules/system";
import {ElMessage} from "element-plus";
import {transformMenu} from "@/hooks/index.js";

export const useAuthStore = defineStore("halfsky-auth", () => {
    // 所有菜单列表
    const authMenuAllList = ref([])
    // 菜单权限列表
    const authMenuList = ref([])
    // 当前页面的 router name，用来做按钮权限筛选
    const routeName = ref("")

    //设置全部菜单
    function setAuthMenuAllList(data) {
      this.authMenuAllList = data
    }
  
    // 隐藏isHide为true的菜单项 
    function removeHiddenItems(data) {
      return data
        .filter(item => !item.meta?.isHide)
        .map(item => {
          const newItem = { ...item };
          if (newItem.children && Array.isArray(newItem.children)) {
            newItem.children = removeHiddenItems(newItem.children);
          }
          return newItem;
        });
    }

    //设置用户菜单
    function setAuthMenuList(data) {
      this.authMenuList = data
    }

    // 所有菜单权限列表
    function getAuthMenuListGet() {
      return authMenuAllList.value;
    }

    // 菜单权限列表 ==> 这里的菜单没有经过任何处理
    function authMenuListGet() {
      return authMenuList.value;
    }

    // 菜单权限列表 ==> 左侧菜单栏渲染，需要剔除 isHide == true
    function showMenuListGet() {
      // console.log("左侧菜单", authMenuList.value)
      const newMenu = removeHiddenItems(authMenuList.value)
      console.log("左侧菜单1", newMenu)
      return newMenu
    }


    // 菜单权限列表 ==> 扁平化之后的一维数组菜单，主要用来添加动态路由
    function flatMenuListGet() {
      return getFlatMenuList(authMenuList.value)
    }


    // 递归处理后的所有面包屑导航列表
    function breadcrumbListGet() {
      return getAllBreadcrumbList(authMenuList.value)
    }

    async function getAuthMenuList() {

      let {data} = await getTreeList()
      let {sysMenu, personMenu} = data

      let sysMenuList = transformMenu(sysMenu, item => ({
        id: item.id,
        children: item.children,
        path: item.router,
        name: item.name,
        component: item.component,
        redirect: "",
        meta: {
          icon: item.icon, // "Menu",
          title: item.name, //"账号管理"
          isLink: item.isLink || false,
          isHide: item.isHide || false,
          isFull: item.isFull || false,
          isAffix: item.isAffix || false,
          isKeepAlive: true
        }
      }));
      let personMenuList = transformMenu(personMenu, item => ({
        id: item.id,
        children: item.children,
        path: item.router,
        name: item.name,
        component: item.component,
        redirect: "",
        meta: {
          icon: item.icon, // "Menu",
          title: item.name, //"账号管理"
          isLink: item.isLink || false,
          isHide: item.isHide || false,
          isFull: item.isFull || false,
          isAffix: item.isAffix || false,
          isKeepAlive: true
        }
      }));
      console.log("全部菜单结构", sysMenuList);
      console.log("管理员菜单结构", personMenuList);
      authMenuAllList.value = sysMenuList;  //全部
      authMenuList.value = personMenuList;   //个人
    }

// Set RouteName
    async function setRouteName(name) {
      this.routeName = name;
    }

    function clearMenusList() {
      authMenuAllList.value = []
      authMenuList.value = []
    }


    return {
      authMenuAllList,
      authMenuList,
      routeName,
      getAuthMenuListGet,
      authMenuListGet,
      showMenuListGet,
      flatMenuListGet,
      breadcrumbListGet,
      getAuthMenuList,
      setRouteName,
      setAuthMenuAllList,
      setAuthMenuList,
      clearMenusList
    }

  }, {
    persist: true,
  }
)

