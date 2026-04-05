<template>
  <div class="table-box">
    <div class="card table-search" v-show="isShowSearch">
      <el-form ref="elSearchFormRef" :inline="true" size="small" :model="searchForm" @keyup.enter="onSubmit">
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.userId" placeholder="用户ID" style="width:139px" />
        </el-form-item>
        <el-form-item label="UID">
          <el-input v-model="searchForm.uid" placeholder="UID" style="width:139px" />
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
        <el-table-column prop="userId" label="用户ID" align="center" :width="80" />
        <el-table-column prop="uid" label="交易所ID" align="center" :min-width="150" />
        <el-table-column prop="parentId" label="上级ID" align="center" :min-width="120" />
        <el-table-column prop="lockLevel" label="锁定团队等级" align="center" :min-width="120" />
        <el-table-column prop="acLevel" label="团队业绩等级" align="center" :min-width="120" />
        <el-table-column prop="level" label="团队等级" align="center" :min-width="120" />
        <el-table-column prop="inviteCount" label="直推人数" align="center" :min-width="120" />
        <el-table-column prop="teamCount" label="团队人数" align="center" :min-width="120" />
        <el-table-column prop="personAchievement" label="个人业绩USD" align="center" :min-width="150" />
        <el-table-column prop="teamTodayAchievement" label="当日团队业绩" align="center" :min-width="150" />
        <el-table-column prop="teamAchievement" label="团队总业绩" align="center" :min-width="150" />
        <el-table-column prop="largeRegionUserId" label="大区用户ID" align="center" :min-width="120" />
        <el-table-column prop="largeRegionAchievement" label="大区业绩" align="center" :min-width="150" />
        <el-table-column prop="fewTeamAchievement" label="小团队业绩" align="center" :min-width="150" />
        <el-table-column prop="giftPersonAchievement" label="赠送质押业绩" align="center" :min-width="150" />
        <el-table-column prop="giftLargeRegionAchievement" label="赠送大区业绩" align="center" :min-width="150" />
        <el-table-column prop="levelLargeRegionAchievement" label="等级赠送大区业绩" align="center" :min-width="150" />
        <el-table-column prop="updatedAt" label="更新时间" align="center" :min-width="160">
          <template #default="{ row }">
            <span v-if="row.updatedAt">{{ formatUnix(row.updatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" align="center" :width="100">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="openEdit(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>

      <Pagination :pageable="pageable" @handleCurrent="handleCurrent" />
    </div>
  </div>

  <ExportFieldFiltering ref="ExportFieldFilteringView" :derived_field="derived_field" @filesselectedfiles="filesselectedfiles" />
  
  <!-- 编辑对话框 -->
  <el-dialog v-model="dialogEdit" title="编辑用户配额" destroy-on-close center width="50%">
    <el-form ref="editFormRef" :model="editFormData" :rules="editRules" label-width="150px">
      <el-form-item label="用户ID" prop="userId">
        <el-input v-model="editFormData.userId" disabled />
      </el-form-item>
      <el-form-item label="锁定等级" prop="lockLevel">
        <el-input v-model="editFormData.lockLevel" type="number" placeholder="请输入锁定等级" clearable />
      </el-form-item>
      <el-form-item label="赠送质押业绩" prop="giftTeamAchievement">
        <el-input v-model="editFormData.giftTeamAchievement" type="number" placeholder="请输入赠送质押业绩" clearable />
      </el-form-item>
      <el-form-item label="赠送大区业绩" prop="giftLargeRegionAchievement">
        <el-input v-model="editFormData.giftLargeRegionAchievement" type="number" placeholder="请输入赠送大区业绩" clearable />
      </el-form-item>
      <el-form-item label="是否免除大区业绩" prop="isExemptionLargeAchievement">
        <el-radio-group v-model="editFormData.isExemptionLargeAchievement" size="small">
          <el-radio-button :value="true">是</el-radio-button>
          <el-radio-button :value="false">否</el-radio-button>
        </el-radio-group>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogEdit = false">取消</el-button>
        <el-button type="primary" @click="handleEditSubmit" :loading="editLoading">确定</el-button>
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
import { magicUserQuotaList, magicUserQuotaUpdate } from '@/api/modules/user.js'

const timeValue = ref([])
const isShowSearch = ref(true)
const searchForm = ref({ userId: null, uid: null })
const tableData = ref([])
const myTable = ref(null)
const pageable = reactive({ pageNum: 1, pageSize: 30, total: 0 })
const ExportFieldFilteringView = ref(null)
const derived_field = ref([])

// 编辑相关
const dialogEdit = ref(false)
const editFormRef = ref(null)
const editLoading = ref(false)
const originalData = ref(null) // 保存原始数据用于比较
const editFormData = ref({
  userId: '',
  lockLevel: null,
  giftTeamAchievement: null,
  giftLargeRegionAchievement: null,
  isExemptionLargeAchievement: false
})
const editRules = ref({
  userId: [{ required: true, message: '用户ID不能为空', trigger: 'blur' }],
  lockLevel: [{ required: true, message: '请输入锁定等级', trigger: 'blur' }],
  giftTeamAchievement: [{ required: true, message: '请输入赠送质押业绩', trigger: 'blur' }],
  giftLargeRegionAchievement: [{ required: true, message: '请输入赠送大区业绩', trigger: 'blur' }],
  isExemptionLargeAchievement: [{ required: true, message: '请选择是否免除大区业绩', trigger: 'change' }]
})

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
    // order: 'id desc',
    beginTime: timeValue.value?.[0] ? dayjs(timeValue.value[0]).startOf('day').unix() : null,
    endTime: timeValue.value?.[1] ? dayjs(timeValue.value[1]).endOf('day').unix() : null,
    ...searchForm.value
  }
  const res = await magicUserQuotaList(params)
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
  searchForm.value = { userId: null, uid: null }
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
  const res = await magicUserQuotaList(params)
  try {
    const downFile = {
      fileName: 'magic_user_quota',
      fileUrl: res.data.url
    }
    await downLoadFile(downFile)
    ElNotification.success('导出成功')
  } catch {
    ElNotification.error('导出失败')
  }
}

// 打开编辑对话框
const openEdit = (row) => {
  // 保存原始数据
  originalData.value = {
    userId: String(row.userId || ''),
    lockLevel: row.lockLevel || null,
    giftTeamAchievement: row.giftTeamAchievement || row.giftPersonAchievement || null,
    giftLargeRegionAchievement: row.giftLargeRegionAchievement || null,
    isExemptionLargeAchievement: row.isExemptionLargeAchievement === true || row.isExemptionLargeAchievement === 1
  }
  
  editFormData.value = {
    userId: String(row.userId || ''),
    lockLevel: row.lockLevel || null,
    giftTeamAchievement: row.giftTeamAchievement || row.giftPersonAchievement || null,
    giftLargeRegionAchievement: row.giftLargeRegionAchievement || null,
    isExemptionLargeAchievement: row.isExemptionLargeAchievement === true || row.isExemptionLargeAchievement === 1
  }
  dialogEdit.value = true
  // 清除表单验证
  if (editFormRef.value) {
    editFormRef.value.clearValidate()
  }
}

// 提交编辑
const handleEditSubmit = async () => {
  if (!editFormRef.value) return
  await editFormRef.value.validate(async (valid) => {
    if (!valid) return
    editLoading.value = true
    try {
      // 构建当前数据
      const currentData = {
        userId: editFormData.value.userId,
        lockLevel: editFormData.value.lockLevel ? Number(editFormData.value.lockLevel) : null,
        giftTeamAchievement: editFormData.value.giftTeamAchievement ? Number(editFormData.value.giftTeamAchievement) : null,
        giftLargeRegionAchievement: editFormData.value.giftLargeRegionAchievement ? Number(editFormData.value.giftLargeRegionAchievement) : null,
        isExemptionLargeAchievement: Boolean(editFormData.value.isExemptionLargeAchievement)
      }
      
      // 比较数据，只传递有变化的字段
      const params = {
        userId: currentData.userId // userId 必须传递
      }
      
      // 比较每个字段
      if (currentData.lockLevel !== originalData.value.lockLevel) {
        params.lockLevel = currentData.lockLevel
      }
      if (currentData.giftTeamAchievement !== originalData.value.giftTeamAchievement) {
        params.giftTeamAchievement = currentData.giftTeamAchievement
      }
      if (currentData.giftLargeRegionAchievement !== originalData.value.giftLargeRegionAchievement) {
        params.giftLargeRegionAchievement = currentData.giftLargeRegionAchievement
      }
      if (currentData.isExemptionLargeAchievement !== originalData.value.isExemptionLargeAchievement) {
        params.isExemptionLargeAchievement = currentData.isExemptionLargeAchievement
      }
      
      // 如果没有数据变化，提示用户
      if (Object.keys(params).length === 1) {
        ElNotification.warning('没有数据变化')
        editLoading.value = false
        return
      }
      
      const res = await magicUserQuotaUpdate(params)
      if (res.code === 200) {
        ElNotification.success('编辑成功')
        dialogEdit.value = false
        getList()
      } else {
        ElNotification.error(res.msg || '编辑失败')
      }
    } catch (error) {
      ElNotification.error('编辑失败')
    } finally {
      editLoading.value = false
    }
  })
}
</script>
