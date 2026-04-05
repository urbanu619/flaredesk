<template>
  <div class="table-box">
    <div class="card table-main">
      <!-- 表格头部 操作按钮 -->
      <div class="table-header">
        <div class="header-button-lf">
          <el-button type="success" @click="editUser( 0 ,{ })">新增</el-button>
        </div>
        <div class="header-button-ri">
          <el-button :icon="Refresh" circle @click="getList"/>
        </div>
      </div>

      <!-- 表格主体 -->
      <el-table :data="tableData" border size="small" row-key="id">
        <el-table-column prop="id" label="id" align="center"/>
        <el-table-column prop="parentId" label="父级id" width="80" align="center">
        </el-table-column>
        <el-table-column prop="name" label="名称" align="center" width="150">
        </el-table-column>
        <el-table-column prop="group" label="分组" align="center"/>
        <el-table-column prop="method" label="方法" align="center"/>
        <el-table-column prop="path" label="路径" align="center" show-overflow-tooltip  min-width="400"/>
        <el-table-column prop="address" label="操作" align="center" width="150" fixed="right">
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
        <el-form-item label="上级ID" > 
            <el-input v-model="drawerFrom.parentId" autocomplete="off" />
          </el-form-item>
          <el-form-item label="名称" prop="name">
            <el-input v-model="drawerFrom.name" autocomplete="off"/>
          </el-form-item>
          <el-form-item label="分组" prop="group">
            <el-input v-model="drawerFrom.group" autocomplete="off" :disabled="drawerType !== 0"/>
          </el-form-item>
          <el-form-item label="api路径" prop="path">
            <el-input v-model="drawerFrom.path" autocomplete="off" :disabled="drawerType !== 0"/>
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
import {onMounted, ref, reactive} from "vue";
import {Refresh} from "@element-plus/icons-vue";
import {apisCreate, apisSet, apisDelete, apisFind} from "@/api/modules/system";
import {ElMessageBox, ElNotification} from "element-plus";

const tableData = ref([]);

const searchForm = ref({

})

const pageable = reactive({
  pageNum: 1,
  pageSize: 1000,
  total: 0
})


onMounted(async () => {
  getList()
})

// 获取列表
const getList = async () => {
  tableData.value = []
  let params = {
    current: pageable.pageNum,
    pageSize: pageable.pageSize,
    order: "id desc",
    ...searchForm.value
  }

  const res = await apisFind(params)
  
  if (res.code === 200) {
    tableData.value = buildTree(res.data?.list || [])
    pageable.total = res.data.paging.total

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


// 删除
const del = async (row) => {
  ElMessageBox.confirm('确定删除当前api?', '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'danger',
  }).then(() => {
    apisDelete({id: row.id}).then(res => {
      if (res.code == 200) {
        ElNotification.success("操作成功")
        getList()
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
  parentId: "",//  integer  上级ID   必需
  name: "",//  string  必需 接口名称
  group: "",//  string  必需 分组
  path: "",//  string  必需 路径

})
const rules = ref({
  parentId: [{required: true, message: '不能为空', trigger: 'blur'}],
  name: [{required: true, message: '不能为空', trigger: 'blur'}],
  group: [{required: true, message: '不能为空', trigger: 'blur'}],
  path: [{required: true, message: '不能为空', trigger: 'blur'}],
})
const formRef = ref()
const editUser = (type, row) => {
  drawerType.value = type
  switch (type) {
    case 0:
      drawerFrom.value = {
        parentId: 0, // Interger 必需
        name: '', //  string  必需
        group: '', //  string  必需
        path: '', //  string 必需
      }
      drawerTitle.value = "新增apis分组"
      break;
    case 1:
      drawerFrom.value = { ...row }
      drawerFrom.value.parentId = row.id || 0
      drawerFrom.value.name = ""
      drawerFrom.value.group = row.group || ''
      drawerFrom.value.path = ""
      drawerTitle.value = "新增api"
      break;
    case 2:
      drawerFrom.value = { ...row }
      drawerFrom.value.id = row.id || 0
      drawerFrom.value.parentId = row.parentId || 0
      drawerFrom.value.name = row.name
      drawerFrom.value.group = row.group
      drawerFrom.value.path = row.path
      drawerTitle.value = "修改api"
      break;
  }
  drawerStatus.value = true
}

// 关闭弹框
const handleClose = () => {
  formRef.value.resetFields()
  drawerStatus.value = false
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
    let {id, parentId, name, group, path } = drawerFrom.value
    switch (drawerType.value) {
      case 0:
        params = {
          parentId: 0,
          name,
          group,
          path,
        }
        apisCreate(params).then(resp => {
          if (resp.code == 200) {
            formEl.resetFields();
            drawerStatus.value = false
            ElNotification.success("添加成功")
            getList()
          } else {
            ElNotification.error(res.msg)

          }
        })
        break;
      case 1:
        params = {
          parentId,
          name,
          group,
          path,
        }
        apisCreate(params).then(resp => {
          if (resp.code == 200) {
            formEl.resetFields();
            drawerStatus.value = false
            ElNotification.success("添加成功")
            getList()
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
          group,
          path,
        }
        apisSet(params).then(resp => {
          if (resp.code == 200) {
            formEl.resetFields();
            drawerStatus.value = false
            ElNotification.success("操作成功")
            getList()
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
