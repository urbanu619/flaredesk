<template>
  <div class="table-box">
    <!--    头部搜索-->
    <div class="card table-search" v-show="isShowSearch">
      <el-form ref="elSearchFormRef" :inline="true" size="small" :model="searchForm" class="demo-form-inline"
               @keyup.enter="onSubmit">
        <template v-if="showAllQuery">
          <!-- <el-form-item label="用户ID">
            <el-input v-model="searchForm.userId" placeholder="用户ID" clearable style="width: 139px;"/>
          </el-form-item>
          <el-form-item label="币种" prop="symbol">
            <el-select v-model="searchForm.symbol" placeholder="请选择" clearable style="width: 139px;">
              <el-option v-for="(item , index) in symbolList " :key="index" :label="item.label" :value="item.value"/>
            </el-select>
          </el-form-item>
          <el-form-item label="总资产" prop="totalAmount">
            <el-select v-model="searchForm.totalAmount" placeholder="请选择" clearable style="width: 139px;">
              <el-option label="总资产大于0" :value="0"/>
              <el-option label="全部" :value="null"/>
            </el-select>
          </el-form-item>
          <el-form-item label="可用资产" prop="balance">
            <el-select v-model="searchForm.balance" placeholder="请选择" clearable style="width: 139px;">
              <el-option label="可用资产大于0" :value="0"/>
              <el-option label="全部" :value="null"/>
            </el-select>
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
        <el-table-column prop="id" label="ID" align="center" width="80"/>
        <el-table-column prop="userId" label="用户ID" width="120" align="center">
          <template #default="scope">
            <span v-address-format="scope.row.userId"></span>
          </template>
        </el-table-column>
        <el-table-column prop="symbol" label="币种" width="80" align="center" />
        <el-table-column prop="totalAmount" label="总资产" align="center" min-width="160"/>
        <el-table-column prop="frozen" label="冻结资产" align="center" min-width="160"/>
        <el-table-column prop="balance" label="净资产" align="center" min-width="160"/>
        <el-table-column prop="version" label="事务版本" align="center" min-width="200"/>
        <el-table-column prop="createdAt" label="创建时间" align="center" min-width="120">
          <template #default="scope">
            {{dayjs(scope.row.createdAt * 1000).format("YYYY-MM-DD HH:mm:ss")}}
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" label="更新时间" align="center" min-width="120">
          <template #default="scope">
            {{dayjs(scope.row.updatedAt * 1000).format("YYYY-MM-DD HH:mm:ss")}}
          </template>
        </el-table-column>
        <!-- <el-table-column prop="" label="操作" min-width="100" align="center" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="openDialog(scope.row)">划扣</el-button>
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


    <!-- 划扣弹窗 -->
    <el-dialog v-model="dialogTransfer" title="划扣操作"  destroy-on-close center width="50%">
      <el-form ref="formRef" :model="transfer" :rules="rules" label-width="auto">
        <el-form-item label="地址ID" prop="userId">
        <el-input v-model="transfer.userId" placeholder="金额" disabled/>
      </el-form-item>
        <el-form-item label="操作币种" prop="symbol">
          <el-select v-model="transfer.symbol" placeholder="请选择币种" clearable class="form-item">
            <el-option v-for="(item , index) in symbolList " :key="item.value" :label="item.label" :value="item.value"/>
          </el-select>
        </el-form-item>
        <el-form-item label="划扣方式" prop="type">
          <el-select v-model="transfer.type" placeholder="请选择划扣方式" clearable class="form-item">
            <el-option v-for="(item , index) in transferMethodList " :key="item.value" :label="item.label" :value="item.value"/>
          </el-select>
        </el-form-item>
        <el-form-item label="金额" prop="amount">
          <el-input v-model="transfer.amount" placeholder="金额" clearable type="number" class="form-item"/>
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

<script setup>
import {onMounted, reactive, ref} from "vue";
import Pagination from "@/components/Pangination/Pagination.vue"

import {Refresh} from "@element-plus/icons-vue";
import { getUserAsset, getUserAssetExport, handlerAsset } from "@/api/modules/user.js";
import {useUserStore} from "@/stores/modules/user";
import {ElNotification} from "element-plus";
import { symbolList, transferMethodList } from '@/utils/dict'
import dayjs from "dayjs";
import { downLoadFile } from "@/api/downloadFile.js"

// ***************搜索框相关*****************************

// 表单规则校验
const rules = ref({
  userId: [{required: true, message: '不能为空', trigger: 'change',}],
  symbol: [{required: true, message: '不能为空', trigger: 'change',}],
  type: [{required: true, message: '不能为空', trigger: 'change',}],
  amount: [{required: true, message: '不能为空', trigger: 'change',}],
})

const showAllQuery = ref(true)
const isShowSearch = ref(true)

const userStore = useUserStore();

const searchForm = ref({
  userId: null, // 用户ID
  symbol: null, // 币种
  totalAmount: null, // 总资产大于
  balance: null, // 可用资产大于
})
const time = ref("")

const tableData = ref([])

const pageable = reactive({
  pageNum: 1,
  pageSize: 30,
  total: 0
}) 

// 划转对象
const transfer = reactive({
  userId: null, // 划转的userId
  symbol: null, // 币种
  type: null, // 类型 增加 减少 
  amount: null // 金额
})

const dialogTransfer = ref(false) // 划转弹窗开启/关闭

const formRef = ref(null) // 划转表单

const loading = ref(false) // loading

onMounted(() => {
  getList()
})


// 划扣弹窗打开
const openDialog = (row) => {
  transfer.symbol = null
  transfer.type = null
  transfer.amount = null
  transfer.userId = row.userId

  dialogTransfer.value = true
}

// 确定划转
const transferConfirm = () => {
  formRef.value?.validate(async (valid) => {
    if (valid) {
      const params = {
        userId: transfer.userId,
        symbol: transfer.symbol,
        amount: `${transfer.type}${transfer.amount}`,
      }
      let res = await handlerAsset(params)
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
// 分页回调
const handleCurrent = (data) => {
  pageable.pageNum = data.current;
  pageable.pageSize = data.pageSize;

  getList()
}
const onResetSearch = () => {
  searchForm.value.userId = null
  searchForm.value.symbol = null
  searchForm.value.totalAmount = null
  searchForm.value.balance = null

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
  let res = await getUserAsset(params)
  if (res.code === 200) {
    tableData.value = res.data.list
    pageable.total = res.data.paging.total

  } else {
    ElNotification.error(res.msg)
  }
};
 
</script>
<style scoped lang="scss">
.form-item{
  width: 100%;
  margin-bottom: 10px;
}

</style>
