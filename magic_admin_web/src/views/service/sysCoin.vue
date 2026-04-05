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
        <div class="header-button-lf">
          <el-button type="success" @click="exportTable">导出</el-button>
        </div>
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
        <el-table-column prop="symbol" label="币种名" align="center" :min-width="120" />
        <el-table-column prop="icon" label="图标" align="center" :min-width="120">
          <template #default="{ row }">
            <el-image
              v-if="row.icon"
              :src="row.icon"
              style="width: 40px; height: 40px"
              fit="contain"
              :preview-src-list="[row.icon]"
              :initial-index="0"
              preview-teleported
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="minWithdrawAmount" label="提现最小限额" align="center" :min-width="130">
          <template #default="{ row }">
            <span>{{ row.minWithdrawAmount !== null && row.minWithdrawAmount !== undefined ? row.minWithdrawAmount : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="maxWithdrawAmount" label="提现最大限额" align="center" :min-width="130">
          <template #default="{ row }">
            <span>{{ row.maxWithdrawAmount === -1 ? '无限制' : (row.maxWithdrawAmount !== null && row.maxWithdrawAmount !== undefined ? row.maxWithdrawAmount : '-') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="withdrawRatio" label="手续费比例" align="center" :min-width="120">
          <template #default="{ row }">
            <span>{{ row.withdrawRatio !== null && row.withdrawRatio !== undefined ? row.withdrawRatio : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="withdrawBaseAmount" label="单次手续费" align="center" :min-width="120">
          <template #default="{ row }">
            <span>{{ row.withdrawBaseAmount !== null && row.withdrawBaseAmount !== undefined ? row.withdrawBaseAmount : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="withdrawNoauditLimit" label="免审额度" align="center" :min-width="120">
          <template #default="{ row }">
            <span>{{ row.withdrawNoauditLimit === 0 ? '需审核' : (row.withdrawNoauditLimit !== null && row.withdrawNoauditLimit !== undefined ? row.withdrawNoauditLimit : '-') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="withdrawStatus" label="是否允许提现" align="center" :min-width="140">
          <template #default="{ row }">
            <el-tag :type="row.withdrawStatus === 1 ? 'success' : 'danger'">
              {{ row.withdrawStatus === 1 ? "允许" : "禁止" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enable" label="是否有效" align="center" :min-width="120">
          <template #default="{ row }">
            <el-tag :type="row.enable ? 'success' : 'info'">
              {{ valueToLabel(booleanOption, row.enable) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" align="center" :min-width="160">
          <template #default="{ row }">
            <span v-if="row.createdAt">{{ formatUnix(row.createdAt) }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" label="更新时间" align="center" :min-width="160">
          <template #default="{ row }">
            <span v-if="row.updatedAt">{{ formatUnix(row.updatedAt) }}</span>
            <span v-else>-</span>
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
  <el-dialog v-model="dialogVisible" title="编辑币种" destroy-on-close center width="60%">
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="150px">
      <el-form-item label="代币名" prop="symbol">
        <el-input v-model="formData.symbol" placeholder="请输入代币名" clearable />
      </el-form-item>
      <el-form-item label="图标" prop="icon">
        <el-input v-model="formData.icon" placeholder="请输入图标URL" clearable />
      </el-form-item>
      <el-form-item label="单笔最小" prop="minWithdrawAmount">
        <el-input v-model="formData.minWithdrawAmount" type="number" placeholder="请输入单笔最小金额" clearable />
      </el-form-item>
      <el-form-item label="单笔最大" prop="maxWithdrawAmount">
        <el-input v-model="formData.maxWithdrawAmount" type="number" placeholder="请输入单笔最大金额" clearable />
      </el-form-item>
      <el-form-item label="单笔基础手续费" prop="withdrawBaseAmount">
        <el-input v-model="formData.withdrawBaseAmount" type="number" placeholder="请输入单笔基础手续费" clearable />
      </el-form-item>
      <el-form-item label="手续费比例" prop="withdrawRatio">
        <el-input v-model="formData.withdrawRatio" type="number" placeholder="请输入手续费比例" clearable />
      </el-form-item>
      <el-form-item label="单笔免审数量" prop="withdrawNoauditLimit">
        <el-input v-model="formData.withdrawNoauditLimit" type="number" placeholder="请输入单笔免审数量" clearable />
      </el-form-item>
      <el-form-item label="状态" prop="withdrawStatus">
        <el-radio-group v-model="formData.withdrawStatus" size="small">
          <el-radio-button :value="true">允许提现</el-radio-button>
          <el-radio-button :value="false">禁止提现</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="是否启用" prop="enable">
        <el-radio-group v-model="formData.enable" size="small">
          <el-radio-button :value="true">启用</el-radio-button>
          <el-radio-button :value="false">禁用</el-radio-button>
        </el-radio-group>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </div>
    </template>
  </el-dialog>
</template>
<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import dayjs from "dayjs";
import { Refresh } from "@element-plus/icons-vue";
import { ElNotification } from "element-plus";
import Pagination from "@/components/Pangination/Pagination.vue";
import ExportFieldFiltering from "@/components/ExportFieldFiltering/index.vue";
import { downLoadFile } from "@/api/downloadFile.js";
import { sysCoinList, sysCoinUpdate } from "@/api/modules/service.js";
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

// 新增相关
const dialogVisible = ref(false);
const formRef = ref(null);
const submitLoading = ref(false);
const formData = ref({
  id: null,
  symbol: "",
  icon: "",
  minWithdrawAmount: null,
  maxWithdrawAmount: null,
  withdrawBaseAmount: null,
  withdrawRatio: null,
  withdrawNoauditLimit: null,
  withdrawStatus: true,
  enable: true
});
const formRules = ref({
  symbol: [{ required: true, message: "请输入代币名", trigger: "blur" }],
  icon: [{ required: true, message: "请输入图标URL", trigger: "blur" }],
  minWithdrawAmount: [{ required: true, message: "请输入单笔最小金额", trigger: "blur" }],
  maxWithdrawAmount: [{ required: true, message: "请输入单笔最大金额", trigger: "blur" }],
  withdrawBaseAmount: [{ required: true, message: "请输入单笔基础手续费", trigger: "blur" }],
  withdrawRatio: [{ required: true, message: "请输入手续费比例", trigger: "blur" }],
  withdrawNoauditLimit: [{ required: true, message: "请输入单笔免审数量", trigger: "blur" }],
  withdrawStatus: [{ required: true, message: "请选择状态", trigger: "change" }],
  enable: [{ required: true, message: "请选择是否启用", trigger: "change" }]
});
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
  const res = await sysCoinList(params);
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
  const res = await sysCoinList(params);
  try {
    await downLoadFile({ fileName: "sys_coin", fileUrl: res.data.url });
    ElNotification.success("导出成功");
  } catch {
    ElNotification.error("导出失败");
  }
};

// 打开编辑对话框
const openEdit = row => {
  formData.value = {
    id: row.id,
    symbol: row.symbol || "",
    icon: row.icon || "",
    minWithdrawAmount: row.minWithdrawAmount ?? row.min_withdraw_amount ?? null,
    maxWithdrawAmount: row.maxWithdrawAmount ?? row.max_withdraw_amount ?? null,
    withdrawBaseAmount: row.withdrawBaseAmount ?? row.withdraw_base_amount ?? null,
    withdrawRatio: row.withdrawRatio ?? row.withdraw_ratio ?? null,
    withdrawNoauditLimit: row.withdrawNoauditLimit ?? row.withdraw_noaudit_limit ?? null,
    withdrawStatus: row.withdrawStatus === 1 || row.withdrawStatus === true,
    enable: row.enable === 1 || row.enable === true
  };
  dialogVisible.value = true;
  // 清除表单验证
  if (formRef.value) {
    formRef.value.clearValidate();
  }
};

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return;
  await formRef.value.validate(async valid => {
    if (!valid) return;
    submitLoading.value = true;
    try {
      const params = {
        id: formData.value.id,
        symbol: formData.value.symbol,
        icon: formData.value.icon,
        minWithdrawAmount: formData.value.minWithdrawAmount ? Number(formData.value.minWithdrawAmount) : null,
        maxWithdrawAmount: formData.value.maxWithdrawAmount ? Number(formData.value.maxWithdrawAmount) : null,
        withdrawBaseAmount: formData.value.withdrawBaseAmount ? Number(formData.value.withdrawBaseAmount) : null,
        withdrawRatio: formData.value.withdrawRatio ? Number(formData.value.withdrawRatio) : null,
        withdrawNoauditLimit: formData.value.withdrawNoauditLimit ? Number(formData.value.withdrawNoauditLimit) : null,
        withdrawStatus: Boolean(formData.value.withdrawStatus),
        enable: Boolean(formData.value.enable)
      };
      const res = await sysCoinUpdate(params);
      if (res.code === 200) {
        ElNotification.success("新增成功");
        dialogVisible.value = false;
        getList();
      } else {
        ElNotification.error(res.msg || "编辑失败");
      }
    } catch (error) {
      ElNotification.error("编辑失败");
    } finally {
      submitLoading.value = false;
    }
  });
};
</script>
