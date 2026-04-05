<template>
  <div class="table-box">
    <div class="card table-search" v-show="isShowSearch">
      <el-form ref="elSearchFormRef" :inline="true" size="small" :model="searchForm" @keyup.enter="onSubmit">
        <el-form-item label="ID"><el-input v-model="searchForm.id" placeholder="ID" style="width: 139px" /></el-form-item>
        <el-form-item label="时间"
          ><el-date-picker
            v-model="timeValue"
            type="daterange"
            range-separator="-"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
        /></el-form-item>
        <div class="operation">
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button
          ><el-button icon="refresh" @click="onResetSearch">重置</el-button>
        </div>
      </el-form>
    </div>
    <div class="card table-main">
      <div class="table-header">
        <div class="header-button-lf"><el-button type="success" @click="exportTable">导出</el-button></div>
        <div class="header-button-ri">
          <el-button :icon="Refresh" circle @click="refresh" /><el-button
            :icon="isShowSearch ? 'ArrowUpBold' : 'ArrowDownBold'"
            circle
            @click="isShowSearch = !isShowSearch"
          />
        </div>
      </div>
      <el-table ref="myTable" :data="tableData" border size="small">
        <el-table-column prop="id" label="ID" align="center" :width="80" />
        <el-table-column
          v-for="(col, idx) in dynamicCols"
          :key="idx"
          :prop="col.json"
          :label="col.comment"
          align="center"
          :min-width="col.minWidth || 150"
        >
          <template v-if="col.isTime" #default="{ row }"
            ><span v-if="row[col.json]">{{ formatUnix(row[col.json]) }}</span></template
          >
        </el-table-column>
      </el-table>
      <Pagination :pageable="pageable" @handleCurrent="handleCurrent" />
    </div>
  </div>
  <ExportFieldFiltering ref="ExportFieldFilteringView" :derived_field="derived_field" @filesselectedfiles="filesselectedfiles" />
</template>
<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import dayjs from "dayjs";
import { Refresh } from "@element-plus/icons-vue";
import { ElNotification } from "element-plus";
import Pagination from "@/components/Pangination/Pagination.vue";
import ExportFieldFiltering from "@/components/ExportFieldFiltering/index.vue";
import { downLoadFile } from "@/api/downloadFile.js";
import { nodeOrderPaymentsList } from "@/api/modules/node.js";
const timeValue = ref([]);
const isShowSearch = ref(true);
const searchForm = ref({ id: null });
const tableData = ref([]);
const myTable = ref(null);
const pageable = reactive({ pageNum: 1, pageSize: 30, total: 0 });
const ExportFieldFilteringView = ref(null);
const derived_field = ref([]);
const formatUnix = val => {
  if (!val) return "";
  return dayjs.unix(val).format("YYYY-MM-DD HH:mm:ss");
};
const dynamicCols = computed(() =>
  derived_field.value
    .filter(c => c.json !== "id")
    .map(c => ({
      ...c,
      isTime: /time|date|At/i.test(c.field),
      minWidth: /describe|remark|comment/i.test(c.json) ? 220 : /id|uid|num/i.test(c.json) ? 120 : 150
    }))
);
onMounted(() => {
  getList();
});
const getList = async () => {
  tableData.value = [];
  const params = {
    current: pageable.pageNum,
    pageSize: pageable.pageSize,
    order: "id desc",
    beginTime: timeValue.value?.[0] ? dayjs(timeValue.value[0]).startOf("day").unix() : null,
    endTime: timeValue.value?.[1] ? dayjs(timeValue.value[1]).endOf("day").unix() : null,
    ...searchForm.value
  };
  const res = await nodeOrderPaymentsList(params);
  if (res.code === 200) {
    tableData.value = res.data.list || [];
    pageable.total = res.data.paging?.total || 0;
    if (res.data.cols) derived_field.value = res.data.cols;
  }
};
const onSubmit = () => {
  pageable.pageNum = 1;
  getList();
};
const refresh = () => getList();
const onResetSearch = () => {
  searchForm.value = { id: null };
  timeValue.value = [];
  getList();
};
const handleCurrent = data => {
  pageable.pageNum = data.current;
  pageable.pageSize = data.pageSize;
  getList();
};
const exportTable = () => {
  ExportFieldFilteringView.value.show(derived_field.value);
};
const filesselectedfiles = async str => {
  const params = {
    current: pageable.pageNum,
    pageSize: pageable.pageSize,
    order: "id desc",
    isExport: true,
    fields: str,
    beginTime: timeValue.value?.[0] ? dayjs(timeValue.value[0]).startOf("day").unix() : null,
    endTime: timeValue.value?.[1] ? dayjs(timeValue.value[1]).endOf("day").unix() : null,
    ...searchForm.value
  };
  const res = await nodeOrderPaymentsList(params);
  try {
    await downLoadFile({ fileName: "node_order_payments", fileUrl: res.data.url });
    ElNotification.success("导出成功");
  } catch {
    ElNotification.error("导出失败");
  }
};
</script>
