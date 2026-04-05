<template>
  <div class="table-box">
    <div class="card table-search" v-show="isShowSearch">
      <el-form ref="elSearchFormRef" :inline="true" size="small" :model="searchForm" @keyup.enter="onSubmit">
        <el-form-item label="ID">
          <el-input v-model="searchForm.id" placeholder="ID" style="width:139px" />
        </el-form-item>
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.rewardUserId" placeholder="用户ID" style="width:139px" />
        </el-form-item>
        <el-form-item label="UID">
          <el-input v-model="searchForm.rewardUserUid" placeholder="UID" style="width:139px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.state" placeholder="状态" style="width:139px">
            <el-option label="待发放" value="waiting" />
            <el-option label="发放成功" value="success" />
          </el-select>
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
        <el-table-column prop="rawUserId" label="源用户ID" align="center" :min-width="120" />
        <el-table-column prop="rawUserUid" label="源用户UID" align="center" :min-width="150" />
        <el-table-column prop="rewardUserId" label="收益用户ID" align="center" :min-width="120" />
        <el-table-column prop="rewardUserUid" label="收益用户UID" align="center" :min-width="150" />
        <el-table-column prop="rewardSymbol" label="奖励币种" align="center" :min-width="120" />
        <el-table-column prop="symbolUsdPrice" label="币种单价" align="center" :min-width="120" />
        <el-table-column prop="capitalUsd" label="奖励基数USD" align="center" :min-width="150" />
        <el-table-column prop="rewardRatio" label="奖励比例" align="center" :min-width="120" />
        <el-table-column prop="rewardUsdValue" label="奖励USD价值" align="center" :min-width="150" />
        <el-table-column prop="rewardQuantity" label="奖励币种数量" align="center" :min-width="150" />
        <el-table-column prop="businessNumber" label="业务场景" align="center" :min-width="120" />
        <el-table-column prop="businessName" label="业务场景名" align="center" :min-width="150" />
        <el-table-column prop="contextName" label="上下文名" align="center" :min-width="150" />
        <el-table-column prop="contextValue" label="上下文值" align="center" :min-width="150" />
        <el-table-column prop="rewardPeriod" label="奖励期数" align="center" :min-width="150" />
        <el-table-column prop="rewardDate" label="奖励日期" align="center" :min-width="150" />
        <el-table-column prop="state" label="状态" align="center" :min-width="120">
          <template #default="{ row }">
            <span>{{ row.state === 'waiting' ? '待发放' : row.state === 'success' ? '发放成功' : row.state }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="describe" label="收益详情" align="center" :min-width="220" />
        <el-table-column prop="updatedAt" label="更新时间" align="center" :min-width="160">
          <template #default="{ row }">
            <span v-if="row.updatedAt">{{ formatUnix(row.updatedAt) }}</span>
          </template>
        </el-table-column>
      </el-table>

      <Pagination :pageable="pageable" @handleCurrent="handleCurrent" />
    </div>
  </div>

  <ExportFieldFiltering ref="ExportFieldFilteringView" :derived_field="derived_field" @filesselectedfiles="filesselectedfiles" />
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import dayjs from 'dayjs'
import { Refresh } from '@element-plus/icons-vue'
import { ElNotification } from 'element-plus'
import Pagination from '@/components/Pangination/Pagination.vue'
import ExportFieldFiltering from '@/components/ExportFieldFiltering/index.vue'
import { downLoadFile } from '@/api/downloadFile.js'
import { magicUserProfitRecordList } from '@/api/modules/user.js'

const timeValue = ref([])
const isShowSearch = ref(true)
const searchForm = ref({ id: null, rewardUserId: null, rewardUserUid: null, state: null })
const tableData = ref([])
const myTable = ref(null)
const pageable = reactive({ pageNum: 1, pageSize: 30, total: 0 })
const ExportFieldFilteringView = ref(null)
const derived_field = ref([])

const formatUnix = (val) => {
  if (!val) return ''
  return dayjs.unix(val).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  getList()
})

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
  const res = await magicUserProfitRecordList(params)
  if (res.code === 200) {
    tableData.value = res.data.list || []
    pageable.total = res.data.paging?.total || 0
    if (res.data.cols) {
      derived_field.value = res.data.cols
    }
  }
}

const onSubmit = () => {
  pageable.pageNum = 1
  getList()
}

const refresh = () => getList()

const onResetSearch = () => {
  searchForm.value = { id: null, rewardUserId: null, rewardUserUid: null, state: null }
  timeValue.value = []
  getList()
}

const handleCurrent = (data) => {
  pageable.pageNum = data.current
  pageable.pageSize = data.pageSize
  getList()
}

const exportTable = () => {
  ExportFieldFilteringView.value.show(derived_field.value)
}

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
  const res = await magicUserProfitRecordList(params)
  try {
    const downFile = {
      fileName: 'magic_user_profit_record',
      fileUrl: res.data.url
    }
    await downLoadFile(downFile)
    ElNotification.success('导出成功')
  } catch {
    ElNotification.error('导出失败')
  }
}
</script>
