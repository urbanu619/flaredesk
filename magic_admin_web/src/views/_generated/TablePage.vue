<template>
  <div class="table-box">
    <div class="card table-search" v-show="isShowSearch">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        size="small"
        :model="searchForm"
        class="demo-form-inline"
        @keyup.enter="onSubmit"
      >
        <!-- 其他筛选项 -->
        <template v-if="showAllQuery">
          <template v-for="(f, idx) in searchFields" :key="idx">
            <el-form-item :label="f.label">
              <component
                :is="f.type === 'select' ? 'el-select' : 'el-input'"
                v-model="searchForm[f.key]"
                :placeholder="f.label"
                style="width:139px"
              >
                <el-option v-if="f.options" v-for="opt in f.options" :key="opt.value" :label="opt.label" :value="opt.value" />
              </component>
            </el-form-item>
          </template>

          <!-- 时间筛选 -->
          <el-form-item label="时间">
            <el-date-picker v-model="timeValue" type="daterange" range-separator="-" start-placeholder="开始时间" end-placeholder="结束时间" />
          </el-form-item>
        </template>

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
        <el-table-column v-for="col in visibleCols" :key="col.json" :prop="col.json" :label="col.comment || col.field" align="center" :min-width="col.minWidth" :width="col.width">
          <template #default="{ row }">
            <span v-if="col.isTime">{{ formatUnix(row[col.json]) }}</span>
            <span v-else-if="col.render === 'valueToLabel'">{{ valueToLabel(col.dict || [], row[col.json]) }}</span>
            <span v-else>{{ row[col.json] }}</span>
          </template>
        </el-table-column>

        <el-table-column fixed="right" label="操作" align="center" :min-width="80">
          <template #default="scope">
            <el-button type="primary" link @click="$emit('edit', scope.row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>

      <slot name="pagination">
        <Pagination :pageable="pageable" @handleCurrent="handleCurrent" />
      </slot>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import dayjs from 'dayjs'
import Pagination from '@/components/Pangination/Pagination.vue'
import { Refresh } from '@element-plus/icons-vue'
import { valueToLabel as v2l } from '@/utils/index.ts'

const props = defineProps({
  listApi: { type: Function, required: true },
  title: { type: String, default: '' },
  searchFields: { type: Array, default: () => [] }
})

const timeValue = ref([])
const isShowSearch = ref(true)
const showAllQuery = ref(true)
const searchForm = ref({})
const tableData = ref([])
const myTable = ref(null)
const pageable = reactive({ pageNum: 1, pageSize: 30, total: 0 })
const derived_field = ref([])

const formatUnix = val => {
  if (!val) return ''
  return dayjs.unix(val).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(()=>{ getList() })

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
  const res = await props.listApi(params)
  if (res.code === 200) {
    tableData.value = res.data.list || []
    pageable.total = res.data.paging?.total || 0
    derived_field.value = res.data.cols || []
  }
}

const onSubmit = ()=>{ pageable.pageNum = 1; getList() }
const refresh = ()=> getList()
const onResetSearch = ()=>{ searchForm.value = {}; timeValue.value = []; getList() }
const handleCurrent = (data)=>{ pageable.pageNum = data.current; pageable.pageSize = data.pageSize; getList() }

const exportTable = ()=>{
  const fields = derived_field.value.map(f=>f.json).join(',')
  const params = { current: pageable.pageNum, pageSize: pageable.pageSize, order: 'id desc', isExport: true, fields, ...searchForm.value }
  props.listApi(params).then(res=>{
    if (res.code === 200) {
      window.open(res.data.url || '')
    }
  })
}

const visibleCols = computed(()=>{
  const cols = derived_field.value || []
  const timeCols = []
  const normalCols = []
  cols.forEach(c=>{
    const json = c.json
    const comment = c.comment || c.field
    const col = { json, comment }
    const isTime = /time|date|At|Date|createdAt|updatedAt/i.test(c.field)
    col.isTime = isTime
    if (json === 'id' || json === 'Id') { col.width = 80 }
    if (!col.width) {
      if (/id|uid|num|count|number|price|amount|quantity|ratio|usd/i.test(json)) col.minWidth = 120
      else if (/describe|remark|comment|content|detail|context/i.test(json)) col.minWidth = 220
      else col.minWidth = 150
    }
    if (isTime) timeCols.push(col)
    else normalCols.push(col)
  })
  const filtered = normalCols.filter(c=>!(c.json==='createdAt' && normalCols.some(x=>x.json==='updatedAt')))
  return [...filtered, ...timeCols]
})

function valueToLabel(dict, val){ return v2l(dict, val) }
</script>

<style scoped>
.table-box{ }
</style>
