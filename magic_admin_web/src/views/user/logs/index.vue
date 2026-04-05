<template>
  <div class="table-box">
    <!--    头部搜索-->
    <div class="card table-search" v-show="isShowSearch">
      <el-form ref="elSearchFormRef" :inline="true" size="small" :model="searchForm" class="demo-form-inline"
               @keyup.enter="onSubmit">
        <template v-if="showAllQuery">
          <!-- <el-form-item label="用户ID">
            <el-input v-model="searchForm.userId" placeholder="用户ID" clearable style="width: 139px;"/>
          </el-form-item> -->
          <!-- <el-form-item label="登陆IP">
            <el-input v-model="searchForm.ip" placeholder="登陆IP" clearable style="width: 139px;"/>
          </el-form-item> -->

          <!-- <el-form-item label="地址">
            <el-input v-model="searchForm.address" placeholder="用户地址" clearable style="width: 139px;"/>
          </el-form-item>

          <el-form-item label="地址ID">
            <el-input v-model="searchForm.userId" placeholder="地址ID" clearable style="width: 139px;"/>
          </el-form-item>-->

          <!-- <el-form-item label="时间">
            <el-date-picker
              v-model="time"
              type="daterange"
              range-separator="-"
              start-placeholder="开始时间"
              end-placeholder="结束时间"/>
          </el-form-item>  -->
        </template>
        <div class="operation">
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onResetSearch">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery = true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery = false" v-else>收起</el-button>
        </div>
      </el-form>
    </div>

    <div class="card table-main">
      <!-- 表格头部 操作按钮 -->
      <div class="table-header">
        <div class="header-button-lf">
          <!--          <el-button type="success">新增</el-button>-->
        </div>
        <div class="header-button-ri">
          <el-button :icon="Refresh" circle @click="refresh"/>
          <el-button :icon="isShowSearch ? 'ArrowUpBold' :'ArrowDownBold'" circle @click="isShowSearch = !isShowSearch"/>
        </div>
      </div>


      <!-- 表格主体 -->
      <el-table :data="tableData" border size="small" align="center">
        <el-table-column prop="id" label="ID" width="80" align="center"/>
        <el-table-column prop="uid" label="交易所ID" width="80" align="center"/>
        <el-table-column prop="userId" label="用户ID" width="80" align="center"/>
        <el-table-column prop="fund_amount" label="联创基金" width="120" align="center"/>
        <el-table-column prop="isLegacy" label="是否老用户" align="center" min-width="140">
          <template #default="scope">
            <el-tag
              :type="scope.row.isLegacy ? 'primary' : 'danger'"
              disable-transitions
              >
              {{ scope.row.isLegacy ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="isSponsor" label="是否保荐商" align="center" min-width="120">
          <template #default="scope">
            <el-tag
              :type="scope.row.isSponsor ? 'primary' : 'danger'"
              disable-transitions
              >
              {{ scope.row.isSponsor ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="state" label="导入状态" min-width="100"/>
        <el-table-column prop="createdAt" label="创建时间" align="center" min-width="100">
          <template #default="scope">
            {{dayjs(scope.row.createdAt * 1000).format("YYYY-MM-DD HH:mm:ss")}}
          </template>
        </el-table-column>
      </el-table>
      <!-- 分页组件 -->
      <slot name="pagination">
        <Pagination
          :pageable="pageable"
          @handleCurrent="handleCurrent"
        />
      </slot>
    </div>

  </div>


</template>

<script setup name="menuMange">
import {onMounted, reactive, ref} from "vue";
import Pagination from "@/components/Pangination/Pagination.vue"
import dayjs from "dayjs";

import {Refresh, Operation, Search} from "@element-plus/icons-vue";
import { getUserLogs} from "@/api/modules/user.js";
import {ElMessageBox, ElNotification} from "element-plus";

// ***************搜索框相关*****************************
const showAllQuery = ref(true)
const isShowSearch = ref(true)
const searchForm = ref({
  // address: null, // 地址
  userId: null, 
  ip: null, // 登录IP
})

const time = ref([]) // 时间选择器

const tableData = ref([])

const proTable = ref();
const menuData = ref([]);


const userRow = ref({})
const pageable = reactive({
  pageNum: 1,
  pageSize: 30,
  total: 0
})


onMounted(() => {
  getList()
})
// 头部查询按钮
const onSubmit = () => {
  console.log('submit!')

  getList()
}
// 刷新
const refresh = () => {
  getList()
}
// 分页回调
const handleCurrent = (data) => {
  pageable.pageNum = data.current;
  pageable.pageSize = data.pageSize;

  getList()
}
const onResetSearch = () => {
  // searchForm.value.address = null
  searchForm.value.userId = null
  searchForm.value.ip = null

  time.value = []
  getList()
}


// 获取用户列表
const getList = async () => {
  tableData.value = []
  let params = {
    current: pageable.pageNum,
    pageSize: pageable.pageSize,
    order: "id desc",
    ...searchForm.value
  }

  if (time.value?.length) { 
    params.startDate = dayjs(time.value[0]).unix()
    params.endDate = dayjs(time.value[1]).unix()
  }

  let res = await getUserLogs(params)
  if (res.code === 200) {
    tableData.value = res.data.list
    pageable.total = res.data.paging.total

  } else {
    ElNotification.error(res.msg)
  }
};
// 复制
const chsoeCopy = (row) => {                            //这里的scope就是每一行的row
  console.log(row);
  if (!row) {
    ElNotification.error("复制失败")
    return
  }

  const save = function (e) {
    e.clipboardData.setData("text/plain", row);
    e.preventDefault();                     // 阻止默认行为
  };
  document.addEventListener("copy", save);  // 添加一个copy事件
  document.execCommand("copy");             // 执行copy方法
  // 复制成功提示
  ElNotification.success("复制成功")
}


// 启用禁用
const starUsewr = (row, type) => {
  console.log("启用禁用")
}
// 删除用户信息
const deleteAccount = (row) => {
  console.log("删除")
  // await useHandleData(deleteUser, { id: [params.id] }, `删除【${params.username}】用户`);
  // proTable.value.getTableList();


};


// 导出 start
const dialogRef = ref()
const downloadFile = () => {
  console.log("导出")
}


const handleSizeChange = () => {

}

const handleCurrentChange = () => {

}

// 导出 end


</script>
<style scoped lang="scss">


</style>
