<template>
  <div class="table-box">
    <div class="card table-search" v-show="isShowSearch">
      <el-form ref="elSearchFormRef" :inline="true" size="small" :model="searchForm" @keyup.enter="onSubmit">
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.userId" placeholder="用户ID" style="width: 139px" />
        </el-form-item>
        <el-form-item label="UID">
          <el-input v-model="searchForm.uid" placeholder="UID" style="width: 139px" />
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
        <div class="header-button-lf">
          <el-button type="success" @click="exportTable">导出</el-button>
        </div>
        <div class="header-button-ri">
          <el-button :icon="Refresh" circle @click="refresh" />
          <el-button :icon="isShowSearch ? 'ArrowUpBold' : 'ArrowDownBold'" circle @click="isShowSearch = !isShowSearch" />
        </div>
      </div>

      <el-table ref="myTable" :data="tableData" border size="small">
        <el-table-column prop="userId" label="用户ID" align="center" :min-width="120" />
        <el-table-column prop="uid" label="交易所ID" align="center" :min-width="150" />
        <!-- 节点收益 -->
        <el-table-column prop="todayNodeProfitQuantity" label="当日节点奖励" align="center" :min-width="150">
          <template #default="{ row }">
            <span>{{ row.todayNodeProfitQuantity !== null && row.todayNodeProfitQuantity !== undefined ? row.todayNodeProfitQuantity : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="nodeProfitQuantity" label="累计节点产出" align="center" :min-width="150">
          <template #default="{ row }">
            <span>{{ row.nodeProfitQuantity !== null && row.nodeProfitQuantity !== undefined ? row.nodeProfitQuantity : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="nodeProfitUsdAmount" label="累计节点产出USD" align="center" :min-width="170">
          <template #default="{ row }">
            <span>{{ row.nodeProfitUsdAmount !== null && row.nodeProfitUsdAmount !== undefined ? row.nodeProfitUsdAmount : '-' }}</span>
          </template>
        </el-table-column>
        <!-- 动态收益 -->
        <el-table-column prop="todayDynamicProfitQuantity" label="当日动态奖励" align="center" :min-width="150">
          <template #default="{ row }">
            <span>{{ row.todayDynamicProfitQuantity !== null && row.todayDynamicProfitQuantity !== undefined ? row.todayDynamicProfitQuantity : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="dynamicProfitQuantity" label="累计动态奖励" align="center" :min-width="150">
          <template #default="{ row }">
            <span>{{ row.dynamicProfitQuantity !== null && row.dynamicProfitQuantity !== undefined ? row.dynamicProfitQuantity : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="dynamicProfitUsdAmount" label="累计动态奖励USD" align="center" :min-width="170">
          <template #default="{ row }">
            <span>{{ row.dynamicProfitUsdAmount !== null && row.dynamicProfitUsdAmount !== undefined ? row.dynamicProfitUsdAmount : '-' }}</span>
          </template>
        </el-table-column>
        <!-- 静态收益 -->
        <el-table-column prop="staticProfitQuantity" label="累计静态产出" align="center" :min-width="150">
          <template #default="{ row }">
            <span>{{ row.staticProfitQuantity !== null && row.staticProfitQuantity !== undefined ? row.staticProfitQuantity : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="staticProfitUsdAmount" label="累计静态产出USD" align="center" :min-width="170">
          <template #default="{ row }">
            <span>{{ row.staticProfitUsdAmount !== null && row.staticProfitUsdAmount !== undefined ? row.staticProfitUsdAmount : '-' }}</span>
          </template>
        </el-table-column>
        <!-- 总收益 -->
        <el-table-column prop="totalProfitQuantity" label="总收益币" align="center" :min-width="130">
          <template #default="{ row }">
            <span>{{ row.totalProfitQuantity !== null && row.totalProfitQuantity !== undefined ? row.totalProfitQuantity : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="totalProfitUsd" label="总收益USD" align="center" :min-width="130">
          <template #default="{ row }">
            <span>{{ row.totalProfitUsd !== null && row.totalProfitUsd !== undefined ? row.totalProfitUsd : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" label="更新时间" align="center" :min-width="160">
          <template #default="{ row }">
            <span v-if="row.updatedAt">{{ formatUnix(row.updatedAt) }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
      </el-table>

      <Pagination :pageable="pageable" @handleCurrent="handleCurrent" />
    </div>
  </div>

  <ExportFieldFiltering ref="ExportFieldFilteringView" :derived_field="derived_field" @filesselectedfiles="filesselectedfiles" />
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import dayjs from "dayjs";
import { Refresh } from "@element-plus/icons-vue";
import { ElNotification } from "element-plus";
import Pagination from "@/components/Pangination/Pagination.vue";
import ExportFieldFiltering from "@/components/ExportFieldFiltering/index.vue";
import { downLoadFile } from "@/api/downloadFile.js";
import { magicUserProfitList } from "@/api/modules/user.js";

const timeValue = ref([]);
const isShowSearch = ref(true);
const searchForm = ref({ userId: null, uid: null });
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
    // order: "id desc",
    beginTime: timeValue.value?.[0] ? dayjs(timeValue.value[0]).startOf("day").unix() : null,
    endTime: timeValue.value?.[1] ? dayjs(timeValue.value[1]).endOf("day").unix() : null,
    ...searchForm.value
  };
  const res = await magicUserProfitList(params);
  if (res.code === 200) {
    tableData.value = res.data.list || [];
    pageable.total = res.data.paging?.total || 0;
    if (res.data.cols) {
      derived_field.value = res.data.cols;
    }
  }
};

const onSubmit = () => {
  pageable.pageNum = 1;
  getList();
};

const refresh = () => getList();

const onResetSearch = () => {
  searchForm.value = { userId: null, uid: null };
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
  const res = await magicUserProfitList(params);
  try {
    const downFile = {
      fileName: "magic_user_profit",
      fileUrl: res.data.url
    };
    await downLoadFile(downFile);
    ElNotification.success("导出成功");
  } catch {
    ElNotification.error("导出失败");
  }
};
</script>
