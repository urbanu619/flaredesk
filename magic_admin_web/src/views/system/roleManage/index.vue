<!--角色管理-->
<template>
  <div class="table-box">
    <div class="card table-main">
      <!-- 表格头部 操作按钮 -->
      <div class="table-header">
        <div class="header-button-lf">
          <el-button type="success" @click="edit( 'add' ,  {})">新增</el-button>
        </div>
        <div class="header-button-ri">
          <el-button :icon="Refresh" circle @click="refreshFn"/>
        </div>
      </div>
      <!-- 表格主体 -->
      <el-table :data="roleList" border size="small">
        <el-table-column prop="id" label="ID" align="center" width="80"/>
        <el-table-column prop="name" label="角色" align="center" />
        <el-table-column prop="menus" label="菜单权限id" align="center"/>
        <el-table-column prop="auths" label="api权限id" align="center"/>
        <el-table-column prop="address" label="操作" align="center">
          <template #default="scope">
            <el-button type="primary" link :disabled=" scope.row .id == 1 " @click="edit( 'edit', scope.row  )">编辑</el-button>
            <el-button type="primary" link :disabled=" scope.row .id == 1 " @click="del(scope.row)"> 删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>


  <!--新增 修改-->
  <el-drawer v-model="drawerStatus" :title="drawerTitle" :before-close="handleClose" :close-on-click-modal="false" size="50%">
    <div class="drawerView">
      <div class="drawerContent">
        <el-form ref="formRef" :model="drawerFrom" :rules="rules" label-width="auto">
          <el-form-item label="角色名称" prop="name">
            <el-input v-model="drawerFrom.name"/>
          </el-form-item>
          <el-form-item label="菜单权限">
            <div class="showMenuTree" @click="showMenuTree">{{ drawerFrom.menus }}</div>
          </el-form-item>
          <el-form-item label="api权限">
            <div class="showMenuTree" @click="showApisTree">{{ drawerFrom.auths }}</div>
          </el-form-item>
          <el-form-item label="备注" prop="">
            <el-input v-model="drawerFrom.desc"/>
          </el-form-item>
        </el-form>
      </div>
      <div class="drawerFooter">
        <el-button size="small" @click="cancalForm(formRef)">取消</el-button>
        <el-button size="small" @click="confrimForm(formRef)">保存</el-button>
      </div>
    </div>
  </el-drawer>

  <!--角色菜单权限树-->
  <el-drawer v-model="dialogFormVisible_menuAuth" title="菜单权限" size="50%" :close-on-click-modal="false">
    <div class="drawerView">
      <div class="drawerContent">
        <!--------------------->
        <div class="list" v-for="(item,index) in sysMenu" :key="index">
          <div class="titleView">
            <el-checkbox name="type" v-model="item.isSelect" @change="sysTitleIsSelectChange(item , index)">{{ item.name }}</el-checkbox>
          </div>
          <div class="children" style="display: flex">
            <div class="titleView" style="padding-left: 15px;" v-for="(a,b) in item.children" :key="a.id" @change="syschildrenIsSelectChange( item , a , index)">
              <div class="titleView">
                <el-checkbox name="type" v-model="a.isSelect">{{ a.name }}</el-checkbox>
              </div>
            </div>
          </div>
        </div>
        <!--------------------------->
      </div>
      <div class="drawerFooter">
        <el-button @click="popupMenuCancal">取消</el-button>
        <el-button type="primary" @click="popupMenuConfirm ">确定</el-button>
      </div>
    </div>
  </el-drawer>

  <!--角色api权限树-->
  <el-drawer v-model="dialogFormVisible_apiAuth" title="api权限" size="50%" :close-on-click-modal="false">
    <div class="drawerView">
      <div class="drawerContent">
        <!--------------------->
        <div class="list" v-for="(item,index) in apisData" :key="index">
          <div class="titleView">
            <el-checkbox name="type" v-model="item.isSelect" @change="sysTitleIsSelectChangeApis(item , index)">{{ item.name }}</el-checkbox>
          </div>
          <div class="children" style="display: flex">
            <div class="titleView" style="padding-left: 15px;" v-for="(a,b) in item.children" :key="a.id" @change="syschildrenIsSelectChangeApis( item , a , index)">
              <div class="titleView">
                <el-checkbox name="type" v-model="a.isSelect">{{ a.name }}</el-checkbox>
              </div>
            </div>
          </div>
        </div>
        <!--------------------------->
      </div>
      <div class="drawerFooter">
        <el-button @click="popupApiCancal">取消</el-button>
        <el-button type="primary" @click="popupApiConfirm ">确定</el-button>
      </div>
    </div>
  </el-drawer>

</template>

<script setup name="roleManage">
import {onMounted, ref, reactive} from "vue";
import {Refresh} from "@element-plus/icons-vue";
import {getSysRoleList, getTreeList, sysRoleCreate, sysRoleSet, sysRoleDel, apisFind} from "../../../api/modules/system";
import {ElMessageBox, ElNotification} from "element-plus";
import {transformMenu} from "@/hooks/index.js"

import {useAuthStore} from "../../../stores/modules/auth.js"

const useAuth = useAuthStore()

const rules = ref({
  name: [{required: true, message: "请输入用户名", trigger: "blur"}],
})
const menuData = ref([])
const roleList = ref([])

const drawerStatus = ref(false)
const drawerTitle = ref("新增")
const formRef = ref()
const drawerFrom = ref({
  id: 0,
  name: "",// string  ID 编号  必需
  desc: "",//string  必需
  apis: [],// string  必需
  menus: [],// string

})

const dialogFormVisible_menuAuth = ref(false)
const dialogFormVisible_apiAuth = ref(false) // api权限弹窗
const userMenus = ref([])
const userApis = ref([]) // 用户拥有的apis权限
const apisData = ref([]) // 系统所有apis列表
const pageable = reactive({ // 查询api接口分页参数
  pageNum: 1,
  pageSize: 1000,
})

const sysMenu = ref([])
const personMenu = ref([])

onMounted(async () => {
  await useAuth.getAuthMenuList()
  get_SysRoleList()

})


const refreshFn = async () => {
  console.log("刷新")
  await useAuth.getAuthMenuList()
  get_SysRoleList()

}

// 获取apis列表
const getList = async () => {
  apisData.value = []
  let params = {
    current: pageable.pageNum,
    pageSize: pageable.pageSize,
    order: "id desc",
  }

  const res = await apisFind(params)
  
  if (res.code === 200) {
    return res
  } else {
    ElNotification.error(res.msg)
  }
};

// 创建多维数组
const buildTree = (data) => {
  const idMap = {};
  const tree = [];
  // 创建 id 到节点的映射，并初始化 children 数组
  data.forEach(item => {
    idMap[item.id] = { ...item, children: [] };
  });

  // 构建树结构
  data.forEach(item => {
    const node = idMap[item.id];
    if (item.parentId && idMap[item.parentId]) {
      idMap[item.parentId].children.push(node);
    } else if (item.parentId === 0 || item.parentId === null || item.parentId === undefined) {
      tree.push(node);
    }
  });

  return tree;
}


//获取角色列表
const get_SysRoleList = async () => {
  let res = await getSysRoleList({ order: "id desc" })
  if (res.code == 200) {
    roleList.value = res.data.list
  } else {
    ElNotification.error(res.msg)
  }
}

// 新增 修改
const edit = (type, row) => {
  console.log("编辑")
  switch (type) {
    case "add":
      drawerFrom.value.name = ''
      drawerFrom.value.desc = ''
      drawerFrom.value.auths = []
      drawerFrom.value.menus = []
      drawerFrom.value.namne = ''
      userMenus.value = []
      userApis.value = []
      drawerTitle.value = "新增角色"
      break;
    case "edit":
      drawerTitle.value = "修改角色"
      drawerFrom.value = { ...row }
      drawerFrom.value.id = row.id
      drawerFrom.value.auths = row.auths
      userMenus.value = row.menus
      userApis.value = row.auths

      break;
  }
  drawerStatus.value = true
  console.log("sysMenu：", sysMenu.value)
  console.log("personMenu：", personMenu.value)
  console.log("拥有权限：", userMenus.value)
  console.log("拥有API权限：", userApis.value)


}
// 删除
const del = (row) => {
  console.log("删除")
  if (row.id == 1) {
    ElNotification.error("禁止操作")
    return;
  }
  ElMessageBox.confirm('删除当前角色?', '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'danger',
  }).then(() => {
    sysRoleDel({id: row.id}).then(res => {
      if (res.code == 200) {
        ElNotification.success("操作成功")
        refreshFn()
      } else {
        ElNotification.error(res.msg)
      }
    })
  }).catch(() => {
    ElNotification.success("取消")
  })


};
const handleClose = () => {
  formRef.value.resetFields()
  drawerStatus.value = false
}


const cancalForm = (formEl) => {
  console.log("取消")
  if (!formEl) return
  formEl.resetFields()
  drawerStatus.value = false
}
const confrimForm = (formEl) => {
  console.log("保存")
  if (!formEl) return
  formEl.validate((valid) => {
    if (valid) {
      let params = {}
      let {id, name, desc, apis, menus, auths} = drawerFrom.value
      switch (drawerTitle.value) {
        case "新增角色":
          params = {
            name,
            desc,
            apis: ["*"],// string  必需
            menus
          }
          sysRoleCreate(params).then(res => {
            if (res.code == 200) {
              ElNotification.success("操作成功")
              refreshFn()
            } else {
              ElNotification.error(res.msg)
            }
          })
          break;
        case "修改角色":
          params = {
            id,
            name,
            desc,
            apis: auths,// string  必需
            menus
          }
          sysRoleSet(params).then(res => {
            if (res.code == 200) {
              ElNotification.success("操作成功")
              refreshFn()
            } else {
              ElNotification.error(res.msg)
            }
          })
          break;
      }
      drawerStatus.value = false
    } else {
      console.log('error submit!')
    }
  })
}
// apis权限选择
const showApisTree = async () => {
  
  const res = await getList()
  apisData.value = buildTree(res.data?.list || [])
  console.log("apisMenu：", apisData.value)
  addIsSelectProperty(apisData.value);   //给所有节点设置isSelect  ：false
  console.log("apisMenu：", apisData.value)

  console.log("api-Ids：", getAllIds(apisData.value))

  if (userApis.value[0] == '*') {
    setIsSelectProperty(apisData.value, getAllIds(apisData.value))
  } else {
    setIsSelectProperty(apisData.value, userApis.value)
  }
  console.log("apisMenu：", apisData.value)
  dialogFormVisible_apiAuth.value = true
}


// 菜单权限选择

const menusList = ref([])   //菜单权限整合
const showMenuTree = () => {
  sysMenu.value = useAuth.getAuthMenuListGet()
  personMenu.value = useAuth.authMenuListGet()
  console.log("sysMenu：", sysMenu.value)
  addIsSelectProperty(sysMenu.value);   //给所有节点设置isSelect  ：false
  console.log("sysMenu：", sysMenu.value)

  console.log("Ids：", getAllIds(sysMenu.value))

  if (userMenus.value[0] == '*') {
    setIsSelectProperty(sysMenu.value, getAllIds(sysMenu.value))
  } else {
    setIsSelectProperty(sysMenu.value, userMenus.value)
  }
  console.log("sysMenu：", sysMenu.value)
  dialogFormVisible_menuAuth.value = true
}


// 菜单权限选择 取消
const popupMenuCancal = (formEl) => {
  // formEl.resetFields();
  drawerStatus.value = false
  dialogFormVisible_menuAuth.value = false
}
// 菜单权限选择 确定保存
const popupMenuConfirm = () => {

  drawerFrom.value.menus = getAllIsSelectIds(sysMenu.value)
  console.log("Ids：", drawerFrom.value.menus)
  dialogFormVisible_menuAuth.value = false
}

// api权限选择 取消
const popupApiCancal = (formEl) => {
  drawerStatus.value = false
  dialogFormVisible_apiAuth.value = false
}
// api权限选择 确定保存
const popupApiConfirm = () => {

  drawerFrom.value.auths = getAllIsSelectIds(apisData.value)
  console.log("Ids：", drawerFrom.value.auths)
  dialogFormVisible_apiAuth.value = false
}

const addIsSelectProperty = (arr) => {
  if (!Array.isArray(arr)) return; // 确保输入是数组
  for (const item of arr) {
    if (typeof item === 'object' && item !== null) {
      item.isSelect = false; // 添加 isSelect 属性
      // 递归处理子数组（通常通过 'children' 属性）
      if (Array.isArray(item.children)) {
        addIsSelectProperty(item.children);
      }
    }
  }
}

const sysTitleIsSelectChange = (item, index) => {
  sysMenu.value.forEach((m, n) => {
    if (m.id == item.id) {
      if (m.children.length != 0) {
        m.children.forEach((a, b) => {
          a.isSelect = m.isSelect
        })
      }
    }
  })
}
// 一级菜单处理
const syschildrenIsSelectChange = (item, a, index) => {
  let c_l_t = 0
  item.children.forEach((v) => {
    if (v.isSelect) {
      c_l_t++;
    }
  })
  if (c_l_t != 0) {
    sysMenu.value.forEach((m, n) => {
      if (m.id == item.id) {
        m.isSelect = true
      }
    })
  } else {
    sysMenu.value.forEach((m, n) => {
      if (m.id == item.id) {
        m.isSelect = false
      }
    })
  }
}
const sysTitleIsSelectChangeApis = (item, index) => {
  apisData.value.forEach((m, n) => {
    if (m.id == item.id) {
      if (m.children.length != 0) {
        m.children.forEach((a, b) => {
          a.isSelect = m.isSelect
        })
      }
    }
  })
}
// 一级菜单处理
const syschildrenIsSelectChangeApis = (item, a, index) => {
  let c_l_t = 0
  item.children.forEach((v) => {
    if (v.isSelect) {
      c_l_t++;
    }
  })
  if (c_l_t != 0) {
    apisData.value.forEach((m, n) => {
      if (m.id == item.id) {
        m.isSelect = true
      }
    })
  } else {
    apisData.value.forEach((m, n) => {
      if (m.id == item.id) {
        m.isSelect = false
      }
    })
  }
}

// 获取所有ID（递归函数）
const getAllIds = (data) => {
  let ids = [];

  function traverse(items) {
    items.forEach(item => {
      ids.push(item.id);
      if (item.children && item.children.length > 0) {
        traverse(item.children);
      }
    });
  }

  traverse(data);
  return ids;
}


// 匹配锁有菜单id相同 isSelect赋值true
const setIsSelectProperty = (arr, arr1) => {
  if (!Array.isArray(arr)) return; // 确保输入是数组
  for (const item of arr) {
    arr1.forEach(id => {
      if (item.id == id) {
        item.isSelect = true;
      }
    })
    // 递归处理子数组（通常通过 'children' 属性）
    if (Array.isArray(item.children)) {
      setIsSelectProperty(item.children, arr1);
    }
  }
}

// 获取所有isSelect == true 的 ID（递归函数）
const getAllIsSelectIds = (data) => {
  let ids = [];

  function traverse(items) {
    items.forEach(item => {
      if (item.isSelect) {
        ids.push(item.id);
      }
      if (item.children && item.children.length > 0) {
        traverse(item.children);
      }
    });
  }

  traverse(data);
  return ids;
}


</script>
<style scoped lang="scss">
.drawerView {
  display: flex;
  flex-direction: column;
  height: 100%;

  .drawerContent {
    flex: 1;
    overflow: auto;

    .showMenuTree {
      cursor: pointer;
    }

    .list {
      margin-bottom: 5px;
      border-bottom: 1px dashed rgba(0, 0, 0, 0.5);

      .children {
        margin-bottom: 10px;

        .item {
          .children {
            .item {
            }
          }
        }
      }
    }


  }

  .drawerFooter {
    height: 40px;
    flex-shrink: 0;
    display: flex;
    justify-content: right;

  }
}

</style>
