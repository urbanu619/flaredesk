<template>
  <div class="table-box">
    <div class="card table-search" v-show="isShowSearch">
      <el-form ref="elSearchFormRef" :inline="true" size="small" :model="searchForm" @keyup.enter="onSubmit">
        <el-form-item label="ID">
          <el-input v-model="searchForm.id" placeholder="ID" style="width: 139px" />
        </el-form-item>
        <el-form-item label="时间">
          <el-date-picker
            v-model="timeValue"
            type="daterange"
            range-separator="-"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
        </el-form-item>
        <div class="operation">
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onResetSearch">重置</el-button>
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
        <el-table-column prop="periodNo" label="期编号" align="center" :min-width="150" />
        <el-table-column prop="jobState" label="状态" align="center" :min-width="120">
          <template #default="{ row }">
            <el-tag type="info" v-if="row.jobState">{{ valueToLabel(stakePeriodJobStateOption, row.jobState) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="price" label="当期单价" align="center" :min-width="120" />
        <el-table-column prop="staticSymbolProfit" label="期静态币总收益" align="center" :min-width="120" />
        <el-table-column prop="staticUsdProfit" label="期USD总收益" align="center" :min-width="120" />
        <el-table-column prop="levelProfit" label="动态等级奖励" align="center" :min-width="120" />
        <el-table-column prop="levelProfitIsCal" label="等级奖励是否计算完成" align="center" :min-width="120">
          <template #default="{ row }">{{ valueToLabel(booleanOption, row.levelProfitIsCal) }}</template>
        </el-table-column>
        <el-table-column prop="desc" label="错误信息" align="center" :min-width="220" />
        <el-table-column prop="updatedAt" label="更新时间" align="center" :min-width="160">
          <template #default="{ row }"
            ><span v-if="row.updatedAt">{{ formatUnix(row.updatedAt) }}</span></template
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
import { sysStakePeriodJobList } from "@/api/modules/task.js";
import { stakePeriodJobStateOption } from "@/utils/dict";
import { valueToLabel } from "@/utils";
const booleanOption = [
  { label: "否", value: false },
  { label: "是", value: true }
];

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
  const res = await sysStakePeriodJobList(params);
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
  const res = await sysStakePeriodJobList(params);
  try {
    await downLoadFile({ fileName: "sys_stake_period_job", fileUrl: res.data.url });
    ElNotification.success("导出成功");
  } catch {
    ElNotification.error("导出失败");
  }
};
</script>
