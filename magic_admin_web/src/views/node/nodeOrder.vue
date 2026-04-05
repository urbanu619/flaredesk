<template>
  <div class="table-box">
    <div class="card table-search" v-show="isShowSearch">
      <el-form ref="elSearchFormRef" :inline="true" size="small" :model="searchForm" @keyup.enter="onSubmit">
        <el-form-item label="ID"><el-input v-model="searchForm.id" placeholder="ID" style="width:139px" /></el-form-item>
        <el-form-item label="时间"><el-date-picker v-model="timeValue" type="daterange" range-separator="-" start-placeholder="开始时间" end-placeholder="结束时间" /></el-form-item>
        <div class="operation"><el-button type="primary" icon="search" @click="onSubmit">查询</el-button><el-button icon="refresh" @click="onResetSearch">重置</el-button></div>
      </el-form>
    </div>
    <div class="card table-main">
      <div class="table-header">
        <div class="header-button-lf">
          <el-button type="primary" @click="openAddOrder">新增</el-button>
          <el-button type="success" @click="exportTable">导出</el-button>
        </div>
        <div class="header-button-ri"><el-button :icon="Refresh" circle @click="refresh" /><el-button :icon="isShowSearch ? 'ArrowUpBold' : 'ArrowDownBold'" circle @click="isShowSearch = !isShowSearch" /></div>
      </div>
      <el-table ref="myTable" :data="tableData" border size="small">
        <el-table-column prop="id" label="ID" align="center" :width="80" />
        <el-table-column prop="userId" label="用户ID" align="center" :min-width="100" />
        <el-table-column prop="uid" label="交易所ID" align="center" :min-width="140" />
        <el-table-column prop="productName" label="产品名称" align="center" :min-width="160" />
        <el-table-column prop="usdPrice" label="单价" align="center" :min-width="100" />
        <el-table-column prop="quantity" label="数量" align="center" :min-width="100" />
        <el-table-column prop="usdAmount" label="总价值" align="center" :min-width="120" />
        <el-table-column prop="payState" label="支付状态" align="center" :min-width="140">
          <template #default="{ row }">
            <el-tag :type="row.payState === 1 ? 'success' : row.payState === 0 ? 'warning' : 'danger'">
              {{ payStateText(row.payState) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="state" label="订单状态" align="center" :min-width="140">
          <template #default="{ row }">
            <el-tag :type="row.state === 1 ? 'success' : row.state === 0 ? 'info' : 'danger'">
              {{ orderStateText(row.state) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column v-for="(col, idx) in dynamicCols" :key="idx" :prop="col.json" :label="col.comment" align="center" :min-width="col.minWidth || 150">
          <template v-if="col.isTime" #default="{ row }"><span v-if="row[col.json]">{{ formatUnix(row[col.json]) }}</span></template>
        </el-table-column>
      </el-table>
      <Pagination :pageable="pageable" @handleCurrent="handleCurrent" />
    </div>
  </div>
  <ExportFieldFiltering ref="ExportFieldFilteringView" :derived_field="derived_field" @filesselectedfiles="filesselectedfiles" />

  <!-- 新增订单对话框 -->
  <el-dialog v-model="dialogAddOrderVisible" title="新增节点订单" destroy-on-close center width="40%">
    <el-form ref="addOrderFormRef" :model="addOrderForm" :rules="addOrderRules" label-width="120px">
      <el-form-item label="交易所UID" prop="uid">
        <el-input v-model="addOrderForm.uid" placeholder="请输入交易所UID" clearable />
      </el-form-item>
      <el-form-item label="产品ID" prop="productId">
        <el-input v-model="addOrderForm.productId" placeholder="请输入产品ID" clearable />
      </el-form-item>
      <el-form-item label="节点数量" prop="quantity">
        <el-input v-model="addOrderForm.quantity" type="number" placeholder="请输入节点数量" clearable />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogAddOrderVisible = false">取消</el-button>
        <el-button type="primary" :loading="addOrderLoading" @click="handleAddOrderSubmit">确定</el-button>
      </div>
    </template>
  </el-dialog>
</template>
<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import dayjs from 'dayjs'
import { Refresh } from '@element-plus/icons-vue'
import { ElNotification } from 'element-plus'
import Pagination from '@/components/Pangination/Pagination.vue'
import ExportFieldFiltering from '@/components/ExportFieldFiltering/index.vue'
import { downLoadFile } from '@/api/downloadFile.js'
import { nodeOrderList, nodeOrderCreate } from '@/api/modules/node.js'

const timeValue = ref([])
const isShowSearch = ref(true)
const searchForm = ref({ id: null })
const tableData = ref([])
const myTable = ref(null)
const pageable = reactive({ pageNum: 1, pageSize: 30, total: 0 })
const ExportFieldFilteringView = ref(null)
const derived_field = ref([])

// 新增订单相关
const dialogAddOrderVisible = ref(false)
const addOrderFormRef = ref(null)
const addOrderLoading = ref(false)
const addOrderForm = ref({
  productId: '',
  quantity: '',
  uid: ''
})
const addOrderRules = ref({
  uid: [{ required: true, message: '请输入交易所UID', trigger: 'blur' }],
  productId: [{ required: true, message: '请输入产品ID', trigger: 'blur' }],
  quantity: [{ required: true, message: '请输入节点数量', trigger: 'blur' }]
})

const formatUnix = (val) => { if (!val) return ''; return dayjs.unix(val).format('YYYY-MM-DD HH:mm:ss') }

// 支付状态文本
const payStateText = (val) => {
  if (val === 1) return '支付成功'
  if (val === 0) return '待支付'
  if (val === 2) return '支付失败'
  return '未知'
}

// 订单状态文本
const orderStateText = (val) => {
  if (val === 1) return '有效'
  if (val === 2) return '过期失效'
  if (val === 3) return '其他失效'
  if (val === 0) return '未知'
  return '未知'
}

const dynamicCols = computed(() =>
  derived_field.value
    .filter(c =>
      ![
        'id',
        'createdAt',
        'updatedAt',
        'recordId',
        'userId',
        'uid',
        'productId',
        'productName',
        'usdPrice',
        'quantity',
        'usdAmount',
        'payState',
        'state',
        'source'
      ].includes(c.json)
    )
    .map(c => ({
      ...c,
      isTime: /time|date|At/i.test(c.field),
      minWidth: /describe|remark|comment/i.test(c.json) ? 220 : (/id|uid|num/i.test(c.json) ? 120 : 150)
    }))
)

onMounted(() => { getList() })

const getList = async () => {
  tableData.value = []
  const params = {
    current: pageable.pageNum,
    pageSize: pageable.pageSize,
    order: 'id desc',
    beginTime: timeValue.value?.[0] ? dayjs(timeValue.value[0]).startOf('day').unix() : null,
    endTime: timeValue.value?.[1] ? dayjs(timeValue.value[1]).endOf('day').unix() : null,
    ...searchForm.value
  }
  const res = await nodeOrderList(params)
  if (res.code === 200) {
    tableData.value = res.data.list || []
    pageable.total = res.data.paging?.total || 0
    if (res.data.cols) derived_field.value = res.data.cols
  }
}

const onSubmit = () => { pageable.pageNum = 1; getList() }
const refresh = () => getList()
const onResetSearch = () => { searchForm.value = { id: null }; timeValue.value = []; getList() }
const handleCurrent = (data) => { pageable.pageNum = data.current; pageable.pageSize = data.pageSize; getList() }
const exportTable = () => { ExportFieldFilteringView.value.show(derived_field.value) }

const filesselectedfiles = async (str) => {
  const params = {
    current: pageable.pageNum,
    pageSize: pageable.pageSize,
    order: 'id desc',
    isExport: true,
    fields: str,
    beginTime: timeValue.value?.[0] ? dayjs(timeValue.value[0]).startOf('day').unix() : null,
    endTime: timeValue.value?.[1] ? dayjs(timeValue.value[1]).endOf('day').unix() : null,
    ...searchForm.value
  }
  const res = await nodeOrderList(params)
  try {
    await downLoadFile({ fileName: 'node_order', fileUrl: res.data.url })
    ElNotification.success('导出成功')
  } catch {
    ElNotification.error('导出失败')
  }
}

// 打开新增订单对话框
const openAddOrder = () => {
  addOrderForm.value = {
    productId: '',
    quantity: '',
    uid: ''
  }
  dialogAddOrderVisible.value = true
  if (addOrderFormRef.value) {
    addOrderFormRef.value.clearValidate()
  }
}

// 提交新增订单
const handleAddOrderSubmit = async () => {
  if (!addOrderFormRef.value) return
  await addOrderFormRef.value.validate(async (valid) => {
    if (!valid) return
    addOrderLoading.value = true
    try {
      const params = {
        productId: addOrderForm.value.productId,
        quantity: addOrderForm.value.quantity,
        uid: addOrderForm.value.uid
      }
      const res = await nodeOrderCreate(params)
      if (res.code === 200) {
        ElNotification.success('新增成功')
        dialogAddOrderVisible.value = false
        getList()
      } else {
        ElNotification.error(res.msg || '新增失败')
      }
    } catch (error) {
      ElNotification.error('新增失败')
    } finally {
      addOrderLoading.value = false
    }
  })
}
</script>
