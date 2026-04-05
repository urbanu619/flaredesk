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
          <!-- <el-form-item label="币种" prop="symbol">
            <el-select v-model="searchForm.symbol" placeholder="请选择" clearable style="width: 139px;">
              <el-option v-for="(item , index) in symbolList " :key="index" :label="item.label" :value="item.value"/>
            </el-select>
          </el-form-item> -->
          <!-- <el-form-item label="业务场景">
            <el-select filterable v-model="searchForm.businessNumber" placeholder="请选择" clearable style="width: 139px;">
              <el-option v-for="(item , index) in businessNameList " :key="index" :label="item.label" :value="item.value"/>
            </el-select>
          </el-form-item>
          <el-form-item label="创建时间">
            <el-date-picker
              v-model="time"
              type="daterange"
              range-separator="-"
              start-placeholder="开始时间"
              end-placeholder="结束时间"/>
          </el-form-item> -->
        
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
        <!-- <div class="header-button-lf">
          <el-button type="success" @click="exportTable">导出</el-button>
          <el-button type="success" @click="exportTableDetail">导出详情</el-button>
        </div> -->
        <div class="header-button-ri">
          <el-button :icon="Refresh" circle @click="refresh"/>
          <el-button :icon="isShowSearch ? 'ArrowUpBold' :'ArrowDownBold'" circle @click="isShowSearch = !isShowSearch"/>
        </div>
      </div>

      <!-- 表格主体 -->
      <el-table :data="tableData" border size="small">
        <!-- <el-table-column prop="assetId" label="资产ID" align="center"/> -->
        <el-table-column prop="id" label="ID" width="80" align="center"/>
        <el-table-column prop="userId" label="用户ID" width="120" align="center">
          <template #default="scope">
            <span v-address-format="scope.row.userId"></span>
          </template>
        </el-table-column>
        <el-table-column prop="symbol" label="币种名" align="center"/>
        <el-table-column prop="balance" label="净资产" align="center" width="140"/>
        <el-table-column prop="frozen" label="冻结资产" align="center" width="140"/>
        <el-table-column prop="businessName" label="业务场景名称" align="center" width="140"/>
        <el-table-column prop="beforeAmount" label="交易前余额" align="center" width="100">
        </el-table-column>
        <el-table-column prop="amount" label="交易金额" align="center" width="120">
          <template #default="scope">
            <span :class="scope.row.amount > 0 ? 'green' : 'red'">
              {{ scope.row.amount }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="afterAmount" label="交易后余额" align="center" width="120"/>
        <!-- <el-table-column prop="contextName" label="上下文名" align="center" width="120"/> -->
        <!-- <el-table-column prop="contextValue" label="上下文值" align="center" width="120"/> -->
        <el-table-column prop="createdAt" label="创建时间" align="center" min-width="100">
          <template #default="scope">
            {{dayjs(scope.row.createdAt * 1000).format("YYYY-MM-DD HH:mm:ss")}}
          </template>
        </el-table-column>
        <!-- <el-table-column prop="updatedAt" label="更新时间" align="center" min-width="100">
          <template #default="scope">
            {{dayjs(scope.row.updatedAt * 1000).format("YYYY-MM-DD HH:mm:ss")}}
          </template>
        </el-table-column> -->
        <el-table-column prop="describe" label="备注" align="center" width="150" show-overflow-tooltip/>

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

<script setup name="question">
import {onMounted, ref, reactive} from "vue";
import dayjs from "dayjs";
import Pagination from "@/components/Pangination/Pagination.vue"
import { Refresh } from "@element-plus/icons-vue";

import { ElNotification} from "element-plus";
import { billRecord, billRecordExport, billRecordExportDetail } from "@/api/modules/finance.js";
import { symbolList, businessNameList } from '@/utils/dict'
import { downLoadFile } from "@/api/downloadFile.js"

// ***************搜索框相关*****************************

const showAllQuery = ref(true)
const isShowSearch = ref(true)
const time = ref([]) // 时间选择器


const searchForm = ref({
  address: null,
  userId: null,
  symbol: null,
  businessNumber: null,
})
const tableData = ref([])


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

// 重置
const onResetSearch = () => {
  searchForm.value.address = null
  searchForm.value.userId = null
  searchForm.value.symbol = null
  searchForm.value.businessNumber = null
  
  time.value = []

  getList()
}

// 分页回调
const handleCurrent = (data) => {
  pageable.pageNum = data.current;
  pageable.pageSize = data.pageSize;
  getList()
}

// 获取列表
const getList = async () => {
  tableData.value = []
  let params = {
    current: pageable.pageNum,
    pageSize: pageable.pageSize,
    order: "id desc",
    ...searchForm.value
  }
  if (time.value?.length) {  
    params.startDate = dayjs(time.value[0]).format("YYYY-MM-DD 00:00:00")
    params.endDate = dayjs(time.value[1]).format("YYYY-MM-DD 23:59:59")
  }
  let res = await billRecord(params)
  if (res.code === 200) {
    tableData.value = res.data.list
    pageable.total = res.data.paging.total

  } else {
    ElNotification.error(res.msg)
  }
};

</script>

<style scoped lang="scss">
.green{
  color: green;
}
.red{
  color: red;
}
</style>
