<template>
  <div class="table-box">
    <div class="card table-search" v-show="isShowSearch">
      <el-form ref="elSearchFormRef" :inline="true" size="small" :model="searchForm" @keyup.enter="onSubmit">
        <el-form-item label="ID">
          <el-input v-model="searchForm.id" placeholder="ID" style="width:139px" />
        </el-form-item>
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.userId" placeholder="用户ID" style="width:139px" />
        </el-form-item>
        <el-form-item label="时间">
          <el-date-picker v-model="timeValue" type="daterange" range-separator="-" start-placeholder="开始时间" end-placeholder="结束时间" />
        </el-form-item>
        <div class="operation">
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onResetSearch">重置</el-button>
        </div>
      </el-form>
    </div>

    <div class="card table-main">
      <div class="table-header">
        <div class="header-button-lf">
          <el-button type="success" @click="exportTable">导出</el-button>
        </div>
        <div class="header-button-ri">
          <el-button :icon="Refresh" circle @click="refresh" />
          <el-button :icon="isShowSearch ? 'ArrowUpBold' : 'ArrowDownBold'" circle @click="isShowSearch = !isShowSearch" />
        </div>
      </div>

      <el-table ref="myTable" :data="tableData" border size="small">
        <el-table-column prop="id" label="ID" align="center" :width="80" />
        <el-table-column prop="userId" label="用户ID" align="center" :min-width="120" />
        <el-table-column prop="uid" label="UID" align="center" :min-width="150" />
        <el-table-column prop="symbol" label="币种" align="center" :min-width="120" />
        <el-table-column prop="amount" label="金额" align="center" :min-width="150" />
        <el-table-column prop="type" label="类型" align="center" :min-width="120" />
        <el-table-column prop="state" label="状态" align="center" :min-width="120" />
        <el-table-column prop="remark" label="备注" align="center" :min-width="220" />
        <el-table-column prop="updatedAt" label="更新时间" align="center" :min-width="160">
          <template #default="{ row }">
            <span v-if="row.updatedAt">{{ formatUnix(row.updatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="审核状态" align="center" :min-width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusInfo(row.status).type">
              {{ getStatusInfo(row.status).label }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" align="center" :width="120">
          <template #default="{ row }">
            <el-button
              v-if="row.direction === 2 && row.status === 0"
              type="primary"
              size="small"
              @click="openWithdrawAudit(row)"
            >
              提现审核
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <Pagination :pageable="pageable" @handleCurrent="handleCurrent" />
    </div>
  </div>

  <ExportFieldFiltering ref="ExportFieldFilteringView" :derived_field="derived_field" @filesselectedfiles="filesselectedfiles" />

  <!-- 提现审核弹窗 -->
  <el-dialog
    v-model="dialogWithdrawAuditVisible"
    title="提现审核"
    width="480px"
    center
    destroy-on-close
  >
    <el-form
      ref="withdrawAuditFormRef"
      :model="withdrawAuditForm"
      :rules="withdrawAuditRules"
      label-width="120px"
    >
      <el-form-item label="提现订单ID" prop="orderId">
        <el-input v-model="withdrawAuditForm.orderId" disabled />
      </el-form-item>
      <el-form-item label="审核结果" prop="step">
        <el-radio-group v-model="withdrawAuditForm.step" size="small">
          <el-radio-button label="1">通过</el-radio-button>
          <el-radio-button label="2">驳回</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="驳回原因" prop="desc">
        <el-input
          v-model="withdrawAuditForm.desc"
          type="textarea"
          :rows="3"
          placeholder="请输入驳回原因（通过可不填）"
          clearable
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogWithdrawAuditVisible = false">取消</el-button>
        <el-button type="primary" :loading="withdrawAuditLoading" @click="handleWithdrawAuditSubmit">
          确定
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import dayjs from 'dayjs'
import { Refresh } from '@element-plus/icons-vue'
import { ElNotification } from 'element-plus'
import Pagination from '@/components/Pangination/Pagination.vue'
import ExportFieldFiltering from '@/components/ExportFieldFiltering/index.vue'
import { downLoadFile } from '@/api/downloadFile.js'
import { magicAssetRwRecordList, handlerTransfer } from '@/api/modules/finance.js'

const timeValue = ref([])
const isShowSearch = ref(true)
const searchForm = ref({ id: null, userId: null })
const tableData = ref([])
const myTable = ref(null)
const pageable = reactive({ pageNum: 1, pageSize: 30, total: 0 })
const ExportFieldFilteringView = ref(null)
const derived_field = ref([])

// 提现审核弹窗相关
const dialogWithdrawAuditVisible = ref(false)
const withdrawAuditFormRef = ref(null)
const withdrawAuditLoading = ref(false)
const withdrawAuditForm = ref({
  orderId: '',
  step: '1', // 1 通过，非1 驳回，这里用 '2' 表示驳回
  desc: ''
})
const withdrawAuditRules = ref({
  orderId: [{ required: true, message: '提现订单ID不能为空', trigger: 'blur' }],
  step: [{ required: true, message: '请选择审核结果', trigger: 'change' }]
  // desc 在驳回时单独校验
})

const formatUnix = (val) => {
  if (!val) return ''
  return dayjs.unix(val).format('YYYY-MM-DD HH:mm:ss')
}

// 审核状态映射
const statusMap = {
  0: { label: '待审核', type: 'warning' },
  1: { label: '审核通过', type: 'success' },
  2: { label: '完成转出/转入', type: 'success' }
}

// 获取审核状态信息
const getStatusInfo = (status) => {
  if (statusMap[status]) {
    return statusMap[status]
  }
  if (status >= 3) {
    return { label: '审核不通过/失败', type: 'danger' }
  }
  return { label: '未知', type: 'info' }
}

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
  const res = await magicAssetRwRecordList(params)
  if (res.code === 200) {
    tableData.value = res.data.list || []
    pageable.total = res.data.paging?.total || 0
    if (res.data.cols) derived_field.value = res.data.cols
  }
}

const onSubmit = () => { pageable.pageNum = 1; getList() }
const refresh = () => getList()
const onResetSearch = () => { searchForm.value = { id: null, userId: null }; timeValue.value = []; getList() }
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
  const res = await magicAssetRwRecordList(params)
  try {
    await downLoadFile({ fileName: 'magic_asset_rw_record', fileUrl: res.data.url })
    ElNotification.success('导出成功')
  } catch {
    ElNotification.error('导出失败')
  }
}

// 打开提现审核弹窗（仅 withdrawState == 0 的记录可点击）
const openWithdrawAudit = (row) => {
  withdrawAuditForm.value = {
    // 提现订单ID：优先使用 orderId，其次使用当前记录 id
    orderId: String(row.orderId || row.id || ''),
    step: '1',
    desc: ''
  }
  dialogWithdrawAuditVisible.value = true
  if (withdrawAuditFormRef.value) {
    withdrawAuditFormRef.value.clearValidate()
  }
}

// 提交提现审核
const handleWithdrawAuditSubmit = async () => {
  if (!withdrawAuditFormRef.value) return
  await withdrawAuditFormRef.value.validate(async (valid) => {
    if (!valid) return
    // 驳回时 desc 必填
    if (withdrawAuditForm.value.step !== '1' && !withdrawAuditForm.value.desc.trim()) {
      ElNotification.warning('请输入驳回原因')
      return
    }
    withdrawAuditLoading.value = true
    try {
      const params = {
        orderId: withdrawAuditForm.value.orderId,
        step: withdrawAuditForm.value.step,
        desc: withdrawAuditForm.value.desc || ''
      }
      const res = await handlerTransfer(params)
      if (res.code === 200) {
        ElNotification.success('审核成功')
        dialogWithdrawAuditVisible.value = false
        getList()
      } else {
        ElNotification.error(res.msg || '审核失败')
      }
    } catch (error) {
      ElNotification.error('审核失败')
    } finally {
      withdrawAuditLoading.value = false
    }
  })
}
</script>
