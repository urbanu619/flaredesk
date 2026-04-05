<template>
  <div class="table-box">
    <div class="card table-search" v-show="isShowSearch">
      <el-form ref="elSearchFormRef" :inline="true" size="small" :model="searchForm" @keyup.enter="onSubmit">
        <el-form-item label="ID"><el-input v-model="searchForm.id" placeholder="ID" style="width: 139px" /></el-form-item>
        <el-form-item label="产品名称">
          <el-input v-model="searchForm.name" placeholder="产品名称" style="width: 139px" />
        </el-form-item>
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
        <el-table-column prop="name" label="产品名称" align="center" :min-width="150" />
        <el-table-column prop="symbol" label="支付币种" align="center" :min-width="120" />
        <el-table-column prop="cycleDay" label="质押周期(天)" align="center" :min-width="130">
          <template #default="{ row }">
            <span>{{ row.cycleDay === 0 ? '活期' : row.cycleDay }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="pledgeMode" label="质押类型" align="center" :min-width="120">
          <template #default="{ row }">
            <el-tag :type="row.pledgeMode === 'current' ? 'success' : 'info'">
              {{ row.pledgeMode === 'current' ? '活期' : '定期' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="periodProfitRatio" label="期收益比例" align="center" :min-width="130">
          <template #default="{ row }">
            <span>{{ row.periodProfitRatio !== null && row.periodProfitRatio !== undefined ? row.periodProfitRatio : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="isQueue" label="是否开启排队" align="center" :min-width="140">
          <template #default="{ row }">
            <el-tag :type="row.isQueue === 1 || row.isQueue === true ? 'success' : 'info'">
              {{ row.isQueue === 1 || row.isQueue === true ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="quotaPerMinutes" label="每分钟释放额度" align="center" :min-width="150">
          <template #default="{ row }">
            <span>{{ row.quotaPerMinutes || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="currentAvailableQuota" label="当前可用额度" align="center" :min-width="150">
          <template #default="{ row }">
            <span>{{ row.currentAvailableQuota || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" align="center" :min-width="100" />
        <el-table-column prop="enable" label="是否有效" align="center" :min-width="120">
          <template #default="{ row }">
            <el-tag :type="row.enable === 1 || row.enable === true ? 'success' : 'danger'">
              {{ row.enable === 1 || row.enable === true ? '可质押' : '暂停质押' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" align="center" :min-width="160">
          <template #default="{ row }">
            <span v-if="row.createdAt">{{ formatUnix(row.createdAt) }}</span>
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
  <el-dialog v-model="dialogEditVisible" title="编辑产品配置" destroy-on-close center width="50%">
    <el-form ref="editFormRef" :model="editFormData" :rules="editFormRules" label-width="150px">
      <el-form-item label="产品名称">
        <el-input v-model="editFormData.name" disabled />
      </el-form-item>
      <el-form-item label="支付币种">
        <el-input v-model="editFormData.symbol" disabled />
      </el-form-item>
      <el-form-item label="质押周期(天)">
        <el-input
          v-model="editFormData.cycleDay"
          placeholder="请输入质押周期(0表示活期)"
        />
      </el-form-item>
      <el-form-item label="质押类型" prop="pledgeMode">
        <el-select v-model="editFormData.pledgeMode" placeholder="请选择质押类型" style="width: 100%">
          <el-option label="活期" value="current" />
          <el-option label="定期" value="regularly" />
        </el-select>
      </el-form-item>
      <el-form-item label="期收益比例" prop="periodProfitRatio">
        <el-input-number
          v-model="editFormData.periodProfitRatio"
          :precision="8"
          :step="0.00000001"
          placeholder="请输入期收益比例"
          controls-position="right"
          class="text-left-input"
          style="width: 100%"
        />
      </el-form-item>
      <el-form-item label="是否开启排队" prop="isQueue">
        <el-radio-group v-model="editFormData.isQueue" size="small">
          <el-radio-button :value="true">开启</el-radio-button>
          <el-radio-button :value="false">关闭</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="每分钟释放额度" prop="quotaPerMinutes">
        <el-input-number
          v-model="editFormData.quotaPerMinutes"
          :min="0"
          :precision="8"
          :step="0.00000001"
          placeholder="请输入每分钟释放生效额度"
          controls-position="right"
          class="text-left-input"
          style="width: 100%"
        />
      </el-form-item>
      <el-form-item label="排序" prop="sort">
        <el-input-number
          v-model="editFormData.sort"
          :min="0"
          placeholder="请输入排序"
          controls-position="right"
          class="text-left-input"
          style="width: 100%"
        />
      </el-form-item>
      <el-form-item label="是否有效" prop="enable">
        <el-radio-group v-model="editFormData.enable" size="small">
          <el-radio-button :value="true">可质押</el-radio-button>
          <el-radio-button :value="false">暂停质押</el-radio-button>
        </el-radio-group>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogEditVisible = false">取消</el-button>
        <el-button type="primary" :loading="editLoading" @click="handleEditSubmit">确定</el-button>
      </div>
    </template>
  </el-dialog>
</template>
<script setup>
import { ref, reactive, onMounted } from "vue";
import dayjs from "dayjs";
import { Refresh } from "@element-plus/icons-vue";
import { ElNotification } from "element-plus";
import Pagination from "@/components/Pangination/Pagination.vue";
import ExportFieldFiltering from "@/components/ExportFieldFiltering/index.vue";
import { downLoadFile } from "@/api/downloadFile.js";
import { magicStakeProductList, magicStakeProductUpdate } from "@/api/modules/service.js";

const timeValue = ref([]);
const isShowSearch = ref(true);
const searchForm = ref({ id: null, name: null });
const tableData = ref([]);
const myTable = ref(null);
const pageable = reactive({ pageNum: 1, pageSize: 30, total: 0 });
const ExportFieldFilteringView = ref(null);
const derived_field = ref([]);

const formatUnix = val => {
  if (!val) return "";
  return dayjs.unix(val).format("YYYY-MM-DD HH:mm:ss");
};

// 编辑相关
const dialogEditVisible = ref(false);
const editFormRef = ref(null);
const editLoading = ref(false);
const originalEditData = ref(null); // 保存原始数据用于比较
const editFormData = ref({
  // 产品ID
  id: null,
  // 展示用名称
  name: "",
  // 展示用支付币种
  symbol: "",
  /**
   * 质押周期天 0:活期
   */
  cycleDay: null,
  /**
   * 质押类型: current 活期 regularly 定期
   */
  pledgeMode: "",
  /**
   * 期收益比例
   */
  periodProfitRatio: null,
  /**
   * 产品是否开启排队
   */
  isQueue: null,
  /**
   * 每分钟释放生效额度
   */
  quotaPerMinutes: null,
  /**
   * 排序
   */
  sort: null,
  /**
   * 是否有效
   */
  enable: null
});
const editFormRules = ref({
  pledgeMode: [{ required: true, message: "请选择质押类型", trigger: "change" }],
  periodProfitRatio: [{ required: true, message: "请输入期收益比例", trigger: "blur" }],
  isQueue: [{ required: true, message: "请选择是否开启排队", trigger: "change" }],
  quotaPerMinutes: [{ required: true, message: "请输入每分钟释放生效额度", trigger: "blur" }],
  sort: [{ required: true, message: "请输入排序", trigger: "blur" }],
  enable: [{ required: true, message: "请选择是否有效", trigger: "change" }]
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
  const res = await magicStakeProductList(params);
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
  searchForm.value = { id: null, name: null };
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
  const res = await magicStakeProductList(params);
  try {
    await downLoadFile({ fileName: "magic_stake_product", fileUrl: res.data.url });
    ElNotification.success("导出成功");
  } catch {
    ElNotification.error("导出失败");
  }
};

// 打开编辑对话框
const openEdit = row => {
  originalEditData.value = {
    id: row.id,
    cycleDay: row.cycleDay,
    enable: row.enable === 1 || row.enable === true,
    isQueue: row.isQueue === 1 || row.isQueue === true,
    periodProfitRatio: row.periodProfitRatio,
    pledgeMode: row.pledgeMode || "",
    quotaPerMinutes: row.quotaPerMinutes ?? null,
    sort: row.sort
  };
  editFormData.value = {
    id: row.id,
    name: row.name || "",
    symbol: row.symbol || "",
    cycleDay: row.cycleDay ?? null,
    pledgeMode: row.pledgeMode || "",
    periodProfitRatio: row.periodProfitRatio ?? null,
    isQueue: row.isQueue === 1 || row.isQueue === true,
    quotaPerMinutes: row.quotaPerMinutes ?? null,
    sort: row.sort ?? null,
    enable: row.enable === 1 || row.enable === true
  };
  dialogEditVisible.value = true;
  if (editFormRef.value) {
    editFormRef.value.clearValidate();
  }
};

// 提交编辑
const handleEditSubmit = async () => {
  if (!editFormRef.value) return;
  await editFormRef.value.validate(async valid => {
    if (!valid) return;
    editLoading.value = true;
    try {
      // 只发送变化的数据
      const params = {
        id: editFormData.value.id
      };
      
      if (editFormData.value.pledgeMode !== originalEditData.value.pledgeMode) {
        params.pledgeMode = editFormData.value.pledgeMode;
      }
      if (editFormData.value.periodProfitRatio !== originalEditData.value.periodProfitRatio) {
        params.periodProfitRatio = Number(editFormData.value.periodProfitRatio);
      }
      if (editFormData.value.isQueue !== originalEditData.value.isQueue) {
        params.isQueue = editFormData.value.isQueue;
      }
      if (editFormData.value.quotaPerMinutes !== originalEditData.value.quotaPerMinutes) {
        params.quotaPerMinutes = Number(editFormData.value.quotaPerMinutes);
      }
      if (editFormData.value.sort !== originalEditData.value.sort) {
        params.sort = Number(editFormData.value.sort);
      }
      if (editFormData.value.enable !== originalEditData.value.enable) {
        params.enable = editFormData.value.enable;
      }

      const res = await magicStakeProductUpdate(params);
      if (res.code === 200) {
        ElNotification.success("编辑成功");
        dialogEditVisible.value = false;
        getList();
      } else {
        ElNotification.error(res.msg || "编辑失败");
      }
    } catch (error) {
      ElNotification.error("编辑失败");
    } finally {
      editLoading.value = false;
    }
  });
};
</script>

<style scoped>
:deep(.text-left-input .el-input__inner) {
  text-align: left;
}
</style>

