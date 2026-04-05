<template>
  <div class="table-box">
    <div class="card table-main">
      <!-- 表格头部 操作按钮 -->
      <div class="table-header">
        <div class="header-button-lf">
          <el-button type="success" @click="editUser( 0 ,{ })">新增</el-button>
        </div>
        <div class="header-button-ri">
          <el-button :icon="Refresh" circle @click="getTableList"/>
        </div>
      </div>



      <!-- 表格主体 -->
      <el-table :data="menuData" border size="small" row-key="name">
        <el-table-column prop="mate.title" label="菜单名称" width="150">
          <template #default="scope"> {{ scope.row.meta.title }}</template>
        </el-table-column>
        <el-table-column prop="" label="图标" width="60" align="center">
          <template #default="scope">
            <el-icon :size="18">
              <component :is="scope.row.meta.icon"></component>
            </el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="mate.icon" label="图标名称" width="100" align="center">
          <template #default="scope"> {{ scope.row.meta.icon }}</template>
        </el-table-column>
        <el-table-column prop="path" label="菜单路径" align="center"/>
        <el-table-column prop="component" label="页面路径" align="center"/>
        <el-table-column prop="address" label="操作" align="center" width="150">
          <template #default="scope">
            <el-button type="primary" link @click="editUser( 1 ,scope.row)"> 新增</el-button>
            <el-button type="primary" link @click="editUser( 2 ,scope.row)"> 编辑</el-button>
            <el-button type="primary" link @click="del(scope.row)"> 删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>


  <!--编辑-->
  <el-drawer
    v-model="drawerStatus"
    :title="drawerTitle"
    :before-close="handleClose"
    size="50%"
    :close-on-click-modal="false"

  >
    <div class="drawerView">
      <div class="drawerContent">
        <el-form ref="formRef"
                 :model="drawerFrom"
                 :rules="rules"
                 label-width="auto"
        >
          <el-form-item label="上级菜单" v-if="drawerType == 1   "> {{ drawerFrom.parentName }}</el-form-item>

          <el-form-item label="图标" prop="icon">
            <SelectIcon v-model:icon-value="drawerFrom.icon"/>
          </el-form-item>
          <el-form-item label="菜单名称" prop="name">
            <el-input v-model="drawerFrom.name" autocomplete="off"/>
          </el-form-item>
          <el-form-item label="路由地址" prop="router">
            <el-input v-model="drawerFrom.router" autocomplete="off"/>
          </el-form-item>
          <el-form-item label="页面路径" prop="component">
            <el-input v-model="drawerFrom.component" autocomplete="off"/>
          </el-form-item>
          <el-form-item label="排序" prop="sort">
            <el-input v-model="drawerFrom.sort" autocomplete="off"/>
          </el-form-item>
          <el-form-item label="隐藏菜单" prop="isHide">
            <el-radio-group v-model="drawerFrom.isHide" size="small">
              <el-radio-button :value="true" label="是" />
              <el-radio-button :value="false" label="否" />
            </el-radio-group>
          </el-form-item>

        </el-form>
      </div>
      <div class="drawerFooter">
        <el-button size="small" @click="cancalForm(formRef)">取消</el-button>
        <el-button size="small" @click="confrimForm(formRef)">保存</el-button>
      </div>


    </div>


  </el-drawer>


</template>

<script setup name="menuMange">
import {onMounted, ref, watch} from "vue";
import authMenuList from "@/assets/json/authMenuList.json";
import {Refresh} from "@element-plus/icons-vue";
import {sysMenuCreate, sysMenuDel, sysMenuSet} from "../../../api/modules/system";
import {ElMessageBox, ElNotification} from "element-plus";
import SelectIcon from "@/components/SelectIcon/index.vue";
import {transformMenu} from "@/hooks/index.js"

import {useAuthStore} from "../../../stores/modules/auth.js"

const useAuth = useAuthStore()
const proTable = ref();
const menuData = ref([]);


const userRow = ref({})


onMounted(async () => {
  await useAuth.getAuthMenuList()
  menuData.value = useAuth.getAuthMenuListGet()
  console.log(menuData.value)
})
const getTableList = async () => {
  await useAuth.getAuthMenuList()

}

// ProTable 实例
watch(() => useAuth.getAuthMenuListGet(), (res) => {
  //系统菜单发生改变
  menuData.value = res
})

// 启用禁用
const starUsewr = (row, type) => {
  console.log("启用禁用")
}
// 删除
const del = async (row) => {
  console.log("删除")

  ElMessageBox.confirm('确定删除当前菜单?', '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'danger',
  }).then(() => {
    sysMenuDel({id: row.id}).then(res => {
      if (res.code == 200) {

        ElNotification.success("操作成功")
        useAuth.getAuthMenuList()
      } else {
        ElNotification.error(res.msg)
      }
    })
  }).catch(() => {
    ElNotification.success("取消删除")
  })


};


//编辑

const drawerStatus = ref(false)
const drawerTitle = ref("")
const drawerType = ref(0)
const drawerFrom = ref({
  id: "",
  parentId: "",//  integer  上级ID  必需
  parentName: "",//  integer  上级name  必需
  name: "",//  string  必需
  icon: "",//  string  必需
  router: "",//  string  必需
  sort: "",//  integer 排序
  isHide: false,//  boolean 是否隐藏菜单
  component: "", //页面路径

})
const rules = ref({
  parentId: [{required: true, message: '不能为空', trigger: 'blur'}],
  name: [{required: true, message: '不能为空', trigger: 'blur'}],
  icon: [{required: true, message: '不能为空', trigger: 'blur'}],
  router: [{required: true, message: '不能为空', trigger: 'blur'}],
  component: [{ required: true, message: '不能为空', trigger: 'blur' }],
  sort: [{required: true, message: '不能为空', trigger: 'blur'}],
})
const formRef = ref()
const editUser = (type, row) => {
  console.log("row", row)
  drawerType.value = type
  switch (type) {
    case 0:
      drawerFrom.value = {
        parentId: 0,
        parentName: "一级菜单",
        name: "",//  string  必需
        icon: "",//  string  必需
        router: "",//  string  必需
        sort: 0,//  integer
        isHide: false,
        component: "",
        isKeepAlive: true
      }
      drawerTitle.value = "新增菜单"
      break;
    case 1:
      drawerFrom.value = { ...row }
      drawerFrom.value.icon = ""
      drawerFrom.value.parentId = row.id || 0
      drawerFrom.value.parentName = row.meta.title
      drawerFrom.value.name = ""
      drawerFrom.value.component = ""
      drawerFrom.value.isKeepAlive = true
      drawerFrom.value.sort = row.sort || 0
      drawerFrom.value.isHide = row.meta?.isHide || false
      drawerTitle.value = "新增菜单"
      break;
    case 2:
      drawerFrom.value = { ...row }
      drawerFrom.value.parentId = row.parent_id || 0
      drawerFrom.value.id = row.id
      drawerFrom.value.icon = row.meta.icon
      drawerFrom.value.parentName = row.meta.title
      drawerFrom.value.router = row.path
      drawerFrom.value.component = row.component
      drawerFrom.value.isKeepAlive = true
      drawerFrom.value.sort = row.sort || 0
      drawerFrom.value.isHide = row.meta?.isHide || false
      drawerTitle.value = "修改菜单"
      break;
  }
  drawerStatus.value = true
}


// 关闭弹框
const handleClose = () => {
  formRef.value.resetFields()
  drawerStatus.value = false
}
// 图片上传成功回调
const upImgSuccess = (data) => {
  drawerFrom.value.avatar = data
}
// 关闭弹框
const cancalForm = (formEl) => {
  formEl.resetFields()
  drawerStatus.value = false
}
// 弹框保存
const confrimForm = (formEl) => {
  if (!formEl) return
  formEl.validate((valid) => {
    if (!valid) return
    let params = {}
    let {parentId, id, parentName, name, icon, router, component, sort, isHide, isKeepAlive} = drawerFrom.value
    switch (drawerType.value) {
      case 0:
        params = {
          parentId: 0,
          name,
          icon,
          router,
          component,
          sort: +sort,
          isHide,
          isKeepAlive
        }
        sysMenuCreate(params).then(resp => {
          if (resp.code == 200) {
            formEl.resetFields();
            drawerStatus.value = false
            ElNotification.success("添加成功")
            useAuth.getAuthMenuList()
          } else {
            ElNotification.error(res.msg)

          }
        })
        break;
      case 1:
        params = {
          parentId,
          name,
          icon,
          router,
          component,
          sort: +sort,
          isHide,
          isKeepAlive
        }
        sysMenuCreate(params).then(resp => {
          if (resp.code == 200) {
            formEl.resetFields();
            drawerStatus.value = false
            ElNotification.success("修改成功")
            useAuth.getAuthMenuList()
          } else {
            ElNotification.error(res.msg)

          }
        })
        break;
      case 2:

        params = {
          id,
          parentId,
          name,
          icon,
          router,
          component,
          sort: +sort,
          isHide: isHide,
          isKeepAlive
        }
        sysMenuSet(params).then(resp => {
          if (resp.code == 200) {
            formEl.resetFields();
            drawerStatus.value = false
            ElNotification.success("操作成功")
            useAuth.getAuthMenuList()
          } else {
            ElNotification.error(res.msg)
          }

        })
        break;
    }
  })
}


</script>
<style scoped lang="scss">
.table-box {
  position: relative;

  .searchView {
    margin-bottom: 10px;
  }
}

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
      border-bottom: 1px dashed;

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
