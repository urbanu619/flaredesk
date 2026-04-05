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
          <el-button type="primary" @click="openAdd">新增</el-button>
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
        <el-table-column prop="name" label="banner名称" align="center" :min-width="180" />
        <el-table-column prop="url" label="bannerURL" align="center" :min-width="220" />
        <el-table-column prop="sort" label="排序" align="center" :min-width="100" />
        <el-table-column prop="remark" label="备注信息" align="center" :min-width="220" />
        <el-table-column prop="updatedAt" label="更新时间" align="center" :min-width="160">
          <template #default="{ row }"
            ><span v-if="row.updatedAt">{{ formatUnix(row.updatedAt) }}</span></template
          >
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
  
  <!-- 新增/编辑对话框 -->
  <el-dialog v-model="dialogVisible" :title="dialogTitle" destroy-on-close center width="50%">
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
      <el-form-item label="ID" prop="id" v-if="isEdit">
        <el-input v-model="formData.id" disabled />
      </el-form-item>
      <el-form-item label="图片名称" prop="name">
        <el-input v-model="formData.name" placeholder="请输入图片名称" clearable />
      </el-form-item>
      <el-form-item label="URL" prop="url">
        <el-input v-model="formData.url" placeholder="请输入URL" clearable />
      </el-form-item>
      <el-form-item label="排序" prop="sort">
        <el-input v-model="formData.sort" type="number" placeholder="请输入排序" clearable />
      </el-form-item>
      <el-form-item label="备注" prop="remark">
        <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注" clearable />
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
import { nodeBannerList, nodeBannerCreate, nodeBannerUpdate } from "@/api/modules/node.js";
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

// 新增/编辑相关
const dialogVisible = ref(false);
const dialogTitle = ref("新增");
const isEdit = ref(false);
const formRef = ref(null);
const submitLoading = ref(false);
const formData = ref({
  id: null,
  name: "",
  url: "",
  sort: null,
  remark: ""
});
const formRules = ref({
  name: [{ required: true, message: "请输入图片名称", trigger: "blur" }],
  url: [{ required: true, message: "请输入URL", trigger: "blur" }],
  sort: [{ required: true, message: "请输入排序", trigger: "blur" }],
  remark: [{ required: false, message: "请输入备注", trigger: "blur" }]
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
  const res = await nodeBannerList(params);
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
  const res = await nodeBannerList(params);
  try {
    await downLoadFile({ fileName: "node_banner", fileUrl: res.data.url });
    ElNotification.success("导出成功");
  } catch {
    ElNotification.error("导出失败");
  }
};

// 打开新增对话框
const openAdd = () => {
  isEdit.value = false;
  dialogTitle.value = "新增";
  formData.value = {
    id: null,
    name: "",
    url: "",
    sort: null,
    remark: ""
  };
  dialogVisible.value = true;
  // 清除表单验证
  if (formRef.value) {
    formRef.value.clearValidate();
  }
};

// 打开编辑对话框
const openEdit = row => {
  isEdit.value = true;
  dialogTitle.value = "编辑";
  formData.value = {
    id: row.id,
    name: row.name || "",
    url: row.url || "",
    sort: row.sort || null,
    remark: row.remark || ""
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
        ...formData.value,
        sort: formData.value.sort ? Number(formData.value.sort) : null
      };
      const res = isEdit.value ? await nodeBannerUpdate(params) : await nodeBannerCreate(params);
      if (res.code === 200) {
        ElNotification.success(isEdit.value ? "编辑成功" : "新增成功");
        dialogVisible.value = false;
        getList();
      } else {
        ElNotification.error(res.msg || (isEdit.value ? "编辑失败" : "新增失败"));
      }
    } catch (error) {
      ElNotification.error(isEdit.value ? "编辑失败" : "新增失败");
    } finally {
      submitLoading.value = false;
    }
  });
};
</script>
