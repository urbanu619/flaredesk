<template>
  <div class="table-box">
    <!--    头部搜索-->
    <div class="card table-search" v-show="isShowSearch">
      <el-form ref="elSearchFormRef" :inline="true" size="small" :model="searchForm" class="demo-form-inline"
               @keyup.enter="onSubmit">
        <template v-if="showAllQuery">
          <!-- <el-form-item label="订单ID">
            <el-input v-model="searchForm.orderId" placeholder="订单ID" clearable style="width: 139px;"/>
          </el-form-item>
          <el-form-item label="用户ID">
            <el-input v-model="searchForm.userId" placeholder="用户ID" clearable style="width: 139px;"/>
          </el-form-item>
          <el-form-item label="资金方向" prop="direction">
            <el-select v-model="searchForm.direction" placeholder="请选择" clearable style="width: 139px;">
              <el-option v-for="(item , index) in directionOption " :key="index" :label="item.label" :value="item.value"/>
            </el-select>
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-select v-model="searchForm.status" placeholder="请选择" clearable style="width: 139px;">
              <el-option v-for="(item , index) in ordertateOption?.slice(0,5) " :key="index" :label="item.label" :value="item.value"/>
            </el-select>
          </el-form-item>
          <el-form-item label="交易币种" prop="symbol">
            <el-select v-model="searchForm.symbol" placeholder="请选择" clearable style="width: 139px;">
              <el-option v-for="(item , index) in symbolList " :key="index" :label="item.label" :value="item.value"/>
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
        <div class="header-button-lf">
          <!-- <el-button type="success" @click="exportTable">导出</el-button> -->
        </div>
        <div class="header-button-ri">
          <el-button :icon="Refresh" circle @click="refresh"/>
          <el-button :icon="isShowSearch ? 'ArrowUpBold' :'ArrowDownBold'" circle @click="isShowSearch = !isShowSearch"/>
        </div>
      </div>

      <!-- 表格主体 -->
      <el-table :data="tableData" border size="small">
        <el-table-column prop="id" label="ID" width="80" align="center"/>
        <el-table-column prop="userId" label="用户ID" width="120" align="center">
         <template #default="scope">
            <span v-address-format="scope.row.userId"></span>
          </template>
        </el-table-column>
        <el-table-column prop="orderId" label="订单ID" width="120" align="center">
        <template #default="scope">
            <span v-address-format="scope.row.orderId"></span>
          </template>
        </el-table-column>
        <el-table-column prop="openId" label="UUID" width="120" align="center">
        <template #default="scope">
            <span v-address-format="scope.row.openId"></span>
          </template>
        </el-table-column>
        <!-- <el-table-column prop="coinId" label="币种ID" width="80" align="center"/> -->
        <el-table-column prop="symbol" label="币种名" width="80" align="center"/>
        <el-table-column prop="direction" label="资金方向" align="center">
          <template #default="scope">
            {{ valueToLabel(directionOption,scope.row.direction) }}
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="交易数量" align="center" width="120">
          <template #default="scope">
            <span :class="scope.row.direction == 1 || scope.row.direction == 3 || scope.row.direction == 5 ? 'green' : 'red'">
              <!-- {{ scope.row.direction == 1 ? '+' : '-' }} -->
              {{ scope.row.amount }}
            </span>
          </template>
        </el-table-column>
        <!-- <el-table-column prop="ratio" label="手续费比例" width="100" align="center"/>
        <el-table-column prop="fee" label="手续费数量" width="100" align="center"/>
        <el-table-column prop="fee" label="手续费数量" width="100" align="center"/>
        <el-table-column prop="realAmount" label="实际操作数量" width="120" align="center"/> -->

        <el-table-column prop="status" label="状态" align="center">
          <template #default="scope" >
            <span v-if="scope.row?.direction === 1">成功</span>
            <span v-else>{{ valueToLabel(ordertateOption,scope.row.status) }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="transferState" label="转出状态" align="center">
          <template #default="scope">
            {{ valueToLabel(transfertateOption,scope.row.transferState) }}
          </template>
        </el-table-column>

        <!-- <el-table-column prop="fromUserId" label="转出方ID" width="120" align="center"/>
        <el-table-column prop="fromUserCode" label="转出方邀请码" width="120" align="center"/>
        <el-table-column prop="toUserId" label="收款方ID" width="120" align="center"/>
        <el-table-column prop="toUserCode" label="收款方邀请码" width="120" align="center"/> -->
        
        <el-table-column prop="createdAt" label="创建时间" align="center" min-width="100">
          <template #default="scope">
            {{ scope.row.createdAt > 0 ? dayjs(scope.row.createdAt * 1000).format("YYYY-MM-DD HH:mm:ss") : ''}}
          </template>
        </el-table-column>
        <!-- <el-table-column prop="updatedAt" label="成交时间" align="center" min-width="100">
          <template #default="scope">
            {{ scope.row.updatedAt > 0 ? dayjs(scope.row.updatedAt * 1000).format("YYYY-MM-DD HH:mm:ss") : ''}}
          </template>
        </el-table-column> -->
        <el-table-column prop="describe" label="备注" align="center" width="180" show-overflow-tooltip/>
        <!-- <el-table-column prop="" label="操作" min-width="150" align="center" fixed="right">
          <template #default="scope">
            <el-button v-if="scope.row?.status === 0" type="primary" link @click="openDialog(scope.row)">提现审核</el-button>
          </template>
        </el-table-column> -->
      </el-table>

      <!-- 分页组件 -->
      <slot name="pagination">
        <Pagination
          :pageable="pageable"
          @handleCurrent="handleCurrent"
        />
      </slot>
    </div>


    <!-- 提现审核弹窗 -->
    <el-dialog v-model="dialogTransfer" title="提现审核"  destroy-on-close center width="50%">
      <el-form ref="formRef" :model="transfer" :rules="rules" label-width="auto">
        <el-form-item label="orderId" prop="txId">
          <el-input v-model="transfer.txId" placeholder="orderId" disabled/>
        </el-form-item>
        <el-form-item label="审核" prop="step">
          <el-select v-model="transfer.step" placeholder="请选择审核结果" clearable class="form-item">
            <!-- 1 通过 2 驳回 -->
            <el-option label="通过" :value="1"/>
            <el-option label="驳回" :value="2"/>
          </el-select>
        </el-form-item>
        <el-form-item label="备注" prop="desc">
          <el-input v-model="transfer.desc" placeholder="请输入备注" clearable class="form-item"/>
        </el-form-item>
      </el-form>
    
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogTransfer = false">取消</el-button>
          <el-button type="primary" @click="transferConfirm" :loading="loading">
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>

  </div>
  


</template>

<script setup name="question">
import {onMounted, ref, reactive} from "vue";
import dayjs from "dayjs";
import Pagination from "@/components/Pangination/Pagination.vue"
import { Refresh } from "@element-plus/icons-vue";

import { ElNotification } from "element-plus";
import {appTransferRecord, appTransferRecordExport, handlerTransfer} from "@/api/modules/finance.js";
import { ordertateOption, transfertateOption, directionOption, symbolList, } from '@/utils/dict'
import { valueToLabel } from '@/utils/index'
import { downLoadFile } from "@/api/downloadFile.js"
 

// ***************搜索框相关*****************************

// 表单规则校验
const rules = ref({
  txId: [{required: true, message: '不能为空', trigger: 'change',}],
  step: [{required: true, message: '不能为空', trigger: 'change',}],
  desc: [{required: true, message: '不能为空', trigger: 'change',}],
})

const showAllQuery = ref(true)
const isShowSearch = ref(true)
const time = ref() // 时间选择器

const searchForm = ref({
  address: null,
  userId: null,
  orderId: null,
  direction: null,
  status: null,
  symbol: null,
})
const tableData = ref([])


// 划转对象
const transfer = reactive({
  txId: null, // 订单id
  step: null, // 1 通过 2 驳回
  desc: null // 备注
})

const dialogTransfer = ref(false) // 划转弹窗开启/关闭

const formRef = ref(null) // 划转表单

const loading = ref(false) // loading

const pageable = reactive({
  pageNum: 1,
  pageSize: 30,
  total: 0
})


onMounted(() => { 
  getList()
})


// 审核弹窗打开
const openDialog = (row) => {
  transfer.txId = row.orderId

  dialogTransfer.value = true
}

// 审核提交
const transferConfirm = () => {
  formRef.value?.validate(async (valid) => {
    if (valid) {
      const params = {
       ...transfer
      }
      let res = await handlerTransfer(params)
      if (res.code === 200) {
        ElNotification.success('操作成功')
        dialogTransfer.value = false
        getList()
      } else {
        // ElNotification.error(res.msg)
      }
      console.log(params, 'submit!')
    } else {
      console.log('error submit!', fields)
    }
  })
}


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
  searchForm.value.orderId = null
  searchForm.value.direction = null
  searchForm.value.status = null
  searchForm.value.symbol = null

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
  let res = await appTransferRecord(params)
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
