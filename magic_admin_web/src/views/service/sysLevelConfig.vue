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
        <el-table-column prop="level" label="用户等级" align="center" :min-width="120" />
        <el-table-column prop="icon" label="等级图标" align="center" :min-width="120">
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
        <el-table-column prop="levelName" label="等级名称" align="center" :min-width="150" />
        <el-table-column prop="personAchievement" label="个人业绩要求" align="center" :min-width="150" />
        <el-table-column prop="teamAchievement" label="团队业绩要求" align="center" :min-width="150" />
        <el-table-column prop="fewTeamAchievement" label="小团队业绩要求" align="center" :min-width="150" />
        <el-table-column prop="staticRatio" label="总静态收益占比" align="center" :min-width="120" />
        <el-table-column prop="avgRatio" label="均分占比" align="center" :min-width="120" />
        <el-table-column prop="weightedRatio" label="加权占比" align="center" :min-width="120" />
        <el-table-column prop="giftLargeRegionAchievement" label="赠送大区业绩" align="center" :min-width="150" />
        <el-table-column fixed="right" label="操作" align="center" :width="120">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="openEdit(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
      <Pagination :pageable="pageable" @handleCurrent="handleCurrent" />
    </div>
  </div>
  <ExportFieldFiltering ref="ExportFieldFilteringView" :derived_field="derived_field" @filesselectedfiles="filesselectedfiles" />

  <!-- 编辑等级配置对话框 -->
  <el-dialog v-model="dialogEditVisible" title="编辑等级配置" destroy-on-close center width="50%">
    <el-form ref="editFormRef" :model="editFormData" :rules="editFormRules" label-width="140px">
      <el-form-item label="等级" prop="level">
        <el-input v-model.number="editFormData.level" type="number" placeholder="请输入等级" disabled />
      </el-form-item>
      <el-form-item label="图标" prop="icon">
        <el-input v-model="editFormData.icon" placeholder="请输入图标URL" clearable />
      </el-form-item>
      <el-form-item label="等级名称" prop="levelName">
        <el-input v-model="editFormData.levelName" placeholder="请输入等级名称" clearable />
      </el-form-item>
      <el-form-item label="升级奖励U" prop="levelUpgradeUsdProfit">
        <el-input v-model.number="editFormData.levelUpgradeUsdProfit" type="number" placeholder="请输入升级奖励U" clearable />
      </el-form-item>
      <el-form-item label="个人业绩要求" prop="personAchievement">
        <el-input v-model.number="editFormData.personAchievement" type="number" placeholder="请输入个人业绩要求" clearable />
      </el-form-item>
      <el-form-item label="团队总业绩要求" prop="teamAchievement">
        <el-input v-model.number="editFormData.teamAchievement" type="number" placeholder="请输入团队总业绩要求" clearable />
      </el-form-item>
      <el-form-item label="小区业绩要求" prop="fewTeamAchievement">
        <el-input v-model.number="editFormData.fewTeamAchievement" type="number" placeholder="请输入小区业绩要求" clearable />
      </el-form-item>
      <el-form-item label="总静态占比" prop="staticRatio">
        <el-input v-model.number="editFormData.staticRatio" type="number" placeholder="请输入总静态占比" clearable />
      </el-form-item>
      <el-form-item label="均分占比" prop="avgRatio">
        <el-input v-model.number="editFormData.avgRatio" type="number" step="0.01" min="0" max="1" placeholder="请输入0-1之间的数字" clearable />
      </el-form-item>
      <el-form-item label="加权分配占比" prop="weightedRatio">
        <el-input v-model.number="editFormData.weightedRatio" type="number" placeholder="请输入加权分配占比" clearable />
      </el-form-item>
      <el-form-item label="升级赠送大区业绩" prop="giftLargeRegionAchievement">
        <el-input v-model.number="editFormData.giftLargeRegionAchievement" type="number" placeholder="请输入升级赠送大区业绩" clearable />
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
import { ref, reactive, onMounted, computed } from "vue";
import dayjs from "dayjs";
import { Refresh } from "@element-plus/icons-vue";
import { ElNotification } from "element-plus";
import Pagination from "@/components/Pangination/Pagination.vue";
import ExportFieldFiltering from "@/components/ExportFieldFiltering/index.vue";
import { downLoadFile } from "@/api/downloadFile.js";
import { sysLevelConfigList, sysLevelConfigCreate } from "@/api/modules/service.js";
const timeValue = ref([]);
const isShowSearch = ref(true);
const searchForm = ref({ id: null });
const tableData = ref([]);
const myTable = ref(null);
const pageable = reactive({ pageNum: 1, pageSize: 30, total: 0 });
const ExportFieldFilteringView = ref(null);
const derived_field = ref([]);

// 编辑相关
const dialogEditVisible = ref(false);
const editFormRef = ref(null);
const editLoading = ref(false);
const editFormData = ref({
  id: null,
  // 等级
  level: null,
  // 图标
  icon: "",
  // 等级名称
  levelName: "",
  // 升级奖励U
  levelUpgradeUsdProfit: null,
  // 个人业绩要求
  personAchievement: null,
  // 团队总业绩要求
  teamAchievement: null,
  // 小区业绩要求
  fewTeamAchievement: null,
  // 总静态占比
  staticRatio: null,
  // 均分占比
  avgRatio: null,
  // 加权分配占比
  weightedRatio: null,
  // 升级赠送大区业绩
  giftLargeRegionAchievement: null
});
const editFormRules = ref({
  level: [{ required: true, message: "请输入等级", trigger: "blur" }],
  icon: [{ required: true, message: "请输入图标", trigger: "blur" }],
  levelName: [{ required: true, message: "请输入等级名称", trigger: "blur" }],
  levelUpgradeUsdProfit: [{ required: true, message: "请输入升级奖励U", trigger: "blur" }],
  personAchievement: [{ required: true, message: "请输入个人业绩要求", trigger: "blur" }],
  teamAchievement: [{ required: true, message: "请输入团队总业绩要求", trigger: "blur" }],
  fewTeamAchievement: [{ required: true, message: "请输入小区业绩要求", trigger: "blur" }],
  staticRatio: [{ required: true, message: "请输入总静态占比", trigger: "blur" }],
  avgRatio: [
    { required: true, message: "请输入均分占比", trigger: "blur" },
    { type: "number", min: 0, max: 1, message: "均分占比必须在0-1之间", trigger: "blur" }
  ],
  weightedRatio: [{ required: true, message: "请输入加权分配占比", trigger: "blur" }],
  giftLargeRegionAchievement: [{ required: true, message: "请输入升级赠送大区业绩", trigger: "blur" }]
});

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
  const res = await sysLevelConfigList(params);
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
  const res = await sysLevelConfigList(params);
  try {
    await downLoadFile({ fileName: "sys_level_config", fileUrl: res.data.url });
    ElNotification.success("导出成功");
  } catch {
    ElNotification.error("导出失败");
  }
};

// 打开编辑对话框
const openEdit = (row) => {
  editFormData.value = {
    id: row.id || null,
    level: row.level ?? null,
    icon: row.icon || "",
    levelName: row.levelName || "",
    levelUpgradeUsdProfit: row.levelUpgradeUsdProfit ?? null,
    personAchievement: row.personAchievement ?? null,
    teamAchievement: row.teamAchievement ?? null,
    fewTeamAchievement: row.fewTeamAchievement ?? null,
    staticRatio: row.staticRatio ?? null,
    avgRatio: row.avgRatio != null ? Number(row.avgRatio) : null,
    weightedRatio: row.weightedRatio ?? null,
    giftLargeRegionAchievement: row.giftLargeRegionAchievement ?? null
  };
  dialogEditVisible.value = true;
  if (editFormRef.value) {
    editFormRef.value.clearValidate();
  }
};

// 提交编辑
const handleEditSubmit = async () => {
  if (!editFormRef.value) return;
  await editFormRef.value.validate(async (valid) => {
    if (!valid) return;
    editLoading.value = true;
    try {
      const params = {
        id: editFormData.value.id,
        // 等级
        level: editFormData.value.level,
        // 图标
        icon: editFormData.value.icon,
        // 等级名称
        levelName: editFormData.value.levelName,
        // 升级奖励U
        levelUpgradeUsdProfit: editFormData.value.levelUpgradeUsdProfit,
        // 个人业绩要求
        personAchievement: editFormData.value.personAchievement,
        // 团队总业绩要求
        teamAchievement: editFormData.value.teamAchievement,
        // 小区业绩要求
        fewTeamAchievement: editFormData.value.fewTeamAchievement,
        // 总静态占比
        staticRatio: editFormData.value.staticRatio,
        // 均分占比
        avgRatio: editFormData.value.avgRatio,
        // 加权分配占比
        weightedRatio: editFormData.value.weightedRatio,
        // 升级赠送大区业绩
        giftLargeRegionAchievement: editFormData.value.giftLargeRegionAchievement
      };
      const res = await sysLevelConfigCreate(params);
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
