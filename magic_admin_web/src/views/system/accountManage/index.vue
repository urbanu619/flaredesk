<!--管理员列表-->
<template>
  <div class="table-box">
    <div class="card table-main">





      <!-- 表格头部 操作按钮 -->
      <div class="table-header">
        <div class="header-button-lf">
          <el-button type="success" @click="editUser(0,{})">新增</el-button>
        </div>
        <div class="header-button-ri">
          <el-button :icon="Refresh" circle @click="refreshFn"/>
        </div>
      </div>




      <!-- 表格主体 -->
      <el-table :data="tableData" border size="small">
        <el-table-column prop="id" label="ID" align="center" width="80"/>
        <el-table-column prop="username" label="账号" align="center"/>
        <el-table-column prop="nickname" label="昵称" align="center"/>
        <el-table-column prop="nickname" label="头像" align="center">
          <template #default="scope">
            <div style="display: flex;justify-content: center;align-items: center;">
              <el-avatar :size="30" :src="scope.row.avatar"/>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="role_name" label="角色" align="center"></el-table-column>
        <el-table-column prop="" label="操作" width="250" align="center">
          <template #default="scope">
            <el-button type="primary" link :disabled="scope.row.id == 1 " @click="editUser(1,scope.row)"> 编辑</el-button>
            <el-button type="danger" link v-if="!scope.row.lock "
                       :disabled="scope.row.id == 1 "
                       @click="starUsewr(scope.row , 0 )"> 禁用
            </el-button>
            <el-button type="primary" link v-if="scope.row.lock"
                       :disabled="scope.row.id == 1 "
                       @click="starUsewr(scope.row , 1)"> 启用
            </el-button>
            <el-button type="danger" link :disabled="scope.row.id == 1 " @click="resetPwd(scope.row )"> 重置密码</el-button>
            <el-button type="danger" link :disabled="scope.row.id == 1 " @click="deleteAccount(scope.row)"> 删除</el-button>
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

          <el-form-item label="账号" prop="username">
            <el-input v-model="drawerFrom.username"/>
          </el-form-item>
          <el-form-item label="昵称" prop="nickname">
            <el-input v-model="drawerFrom.nickname"/>
          </el-form-item>
          <el-form-item label="头像" prop="avatar">
            <upImage @upImgSuccess="upImgSuccess" :avatar="drawerFrom.avatar"></upImage>
          </el-form-item>
          <el-form-item label="角色" prop="roleId">
            <el-select v-model="drawerFrom.roleId" placeholder="请选择">
              <el-option v-for="(item , index) in roleList " :key="index" :label="item.name" :value="item.id"/>
            </el-select>

          </el-form-item>
          <el-form-item label="密码" :prop="drawerTitle=='新增管理员' ?'password':'' ">
            <el-input type="password" v-model="drawerFrom.password"/>
          </el-form-item>
          <el-form-item label="重置密码">
            <el-switch v-model="drawerFrom.resetGoogleKey"/>
          </el-form-item>
          <el-form-item label="锁定">
            <el-switch v-model="drawerFrom.lock"/>
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
import {onMounted, reactive, ref} from "vue";
import {Refresh} from "@element-plus/icons-vue";
import {getSysRoleList, getUserList, sysCreateUser, sysDelUser, sysEditUser} from "../../../api/modules/system";
import {ElNotification, ElMessageBox} from "element-plus";
import upImage from "../../../components/Upload/Img.vue"


const onSubmit = () => {
  console.log('submit!')
}

const roleList = ref([])
const tableData = ref([])

const proTable = ref();
const userRow = ref({})

onMounted(() => {
  get_SysRoleList()
})
const refreshFn = () => {
  console.log("刷新")
  get_SysRoleList()
}
//获取角色
const get_SysRoleList = async () => {
  let res = await getSysRoleList({ order: "id desc" })
  if (res.code == 200) {
    roleList.value = res.data.list
    console.log(roleList.value)
    getTableList()
  } else {
    ElNotification.error(res.msg)
  }


}
//获取管理员
const getTableList = async () => {
  let res = await getUserList({order: "id desc"})
  if (res.code == 200) {
    res.data.list.forEach((v) => {
      roleList.value.forEach(a => {
        if (v.role_id == a.id) {
          v.role_name = a.name
        }
      })

    })
    tableData.value = res.data.list

  } else {
    ElNotification.error(res.msg)
  }


};

// 启用禁用
const starUsewr = (row, type) => {
  console.log("启用禁用", row.id)
  if (row.id == 1) {
    ElNotification.error("禁止操作")
    return
  }


}
// 删除用户信息
const deleteAccount = (row) => {
  console.log("删除")
  if (row.id == 1) {
    ElNotification.error("禁止操作")
    return
  }
  ElMessageBox.confirm(`删除账号${row.username}吗?`, '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'danger',
  }).then(() => {
    let params = {
      id: row.id
    }
    sysDelUser(params).then(res => {
      if (resp.code == 200) {
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
//重置密码
const resetPwd = () => {
  ElMessageBox.confirm('重置管理员密码?', '提示', {
      confirmButtonText: '重置',
      cancelButtonText: '取消',
      type: 'danger',
    }
  )
    .then(() => {

    })
    .catch(() => {

    })
}


//新增 编辑 弹框  start
const drawerStatus = ref(false)
const drawerTitle = ref("")

const drawerFrom = ref({
  userId: "",
  avatar: "",
  username: "",
  nickname: "",
  password: "",
  resetGoogleKey: false,
  lock: false,
  roleId: ""
})

const rules = ref({
  username: [{required: true, message: "请输入昵称", trigger: "blur"}],
  nickname: [{required: true, message: "请输入昵称", trigger: "blur"}],
  password: [{required: true, message: "请输入密码", trigger: "blur"}],
  roleId: [{required: true, message: "请选择", trigger: "blur"}]
})
const formRef = ref()

const editUser = (type, row) => {
  switch (type) {
    case 0:
      drawerFrom.value = {
        userId: "",
        avatar: "",
        username: "",
        nickname: "",
        password: "",
        roleId: "",
        resetGoogleKey: false,
        lock: false
      }
      drawerTitle.value = "新增管理员"
      break;
    case 1:
      drawerFrom.value = { ...row }
      drawerFrom.value.roleId = row.role_id
      drawerTitle.value = "修改管理员"
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

    switch (drawerTitle.value) {
      case "新增管理员":
        let {avatar, username, nickname, password, roleId} = drawerFrom.value
        params = {
          avatar, username, nickname, password, roleId
        }
        sysCreateUser(params).then(resp => {
          if (resp.code == 200) {
            formEl.resetFields();
            drawerStatus.value = false
            ElNotification.success("添加成功")
            refreshFn()
          } else {
            ElNotification.error(res.msg)

          }
        })


        break;
      case "修改管理员":

        params = {
          userId: drawerFrom.value.userId,
          avatar: drawerFrom.value.avatar,
          nickname: drawerFrom.value.nickname,
          password: drawerFrom.value.password,
          resetGoogleKey: drawerFrom.value.resetGoogleKey,
          lock: drawerFrom.value.lock,
        }
        sysEditUser(params).then(resp => {
          if (resp.code == 200) {
            formEl.resetFields();
            drawerStatus.value = false
            ElNotification.success("操作成功")
            refreshFn()
          } else {
            ElNotification.error(res.msg)
          }

        })
        break;
    }
  })
}


//新增 编辑 弹框  end


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

