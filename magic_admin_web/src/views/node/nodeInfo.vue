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
        <el-table-column prop="nodeName" label="节点名称" align="center" :min-width="180" />
        <el-table-column prop="supportPaymentSymbols" label="支持的支付币种" align="center" :min-width="160" />
        <el-table-column prop="usdPrice" label="单价" align="center" :min-width="120" />
        <el-table-column prop="soldQuantity" label="已售出数量" align="center" :min-width="120" />
        <el-table-column prop="upperQuantityLimit" label="可售出上限" align="center" :min-width="120" />
        <el-table-column prop="limitTime" label="结束时间" align="center" :min-width="160">
          <template #default="{ row }"
            ><span v-if="row.limitTime">{{ formatUnix(row.limitTime) }}</span></template
          >
        </el-table-column>
        <el-table-column prop="sort" label="排序" align="center" :min-width="100" />
        <el-table-column prop="enable" label="是否有效" align="center" :min-width="120">
          <template #default="{ row }">
            <span>{{ row.enable === true || row.enable === 1 ? '有效' : '无效' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="isDisplay" label="是否展示" align="center" :min-width="120">
          <template #default="{ row }">
            <span>{{ row.isDisplay === true || row.isDisplay === 1 ? '展示' : '不展示' }}</span>
          </template>
        </el-table-column>
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
  
  <!-- 编辑对话框 -->
  <el-dialog v-model="dialogEdit" title="编辑节点信息" destroy-on-close center width="60%">
    <el-form ref="editFormRef" :model="editFormData" :rules="editRules" label-width="165px">
      <el-form-item label="节点ID" prop="id">
        <el-input v-model="editFormData.id" disabled />
      </el-form-item>
      <el-form-item label="单价" prop="usdPrice">
        <el-input v-model="editFormData.usdPrice" placeholder="请输入单价" clearable />
      </el-form-item>
      <el-form-item label="支持的支付币种" prop="supportPaymentSymbols">
        <el-input v-model="editFormData.supportPaymentSymbols" placeholder="请输入支持的支付币种，用逗号分割" clearable />
      </el-form-item>
      <el-form-item label="平台最大限制" prop="upperQuantityLimit">
        <el-input v-model="editFormData.upperQuantityLimit" placeholder="请输入平台最大限制" clearable />
      </el-form-item>
      <el-form-item label="产品用户购买有效上限" prop="personBuyLimit">
        <el-input v-model="editFormData.personBuyLimit" type="number" placeholder="请输入产品用户购买有效上限" clearable />
      </el-form-item>
      <el-form-item label="分红比例" prop="dividendRatio">
        <el-input v-model="editFormData.dividendRatio" placeholder="请输入分红比例" clearable />
      </el-form-item>
      <el-form-item label="分红交易所UID" prop="dividendUserUid">
        <el-input v-model="editFormData.dividendUserUid" placeholder="请输入分红交易所UID" clearable />
      </el-form-item>
      <el-form-item label="最近一层返佣比例" prop="rebateRatio">
        <el-input v-model="editFormData.rebateRatio" placeholder="请输入最近一层返佣比例" clearable />
      </el-form-item>
      <el-form-item label="排序" prop="sort">
        <el-input v-model="editFormData.sort" type="number" placeholder="请输入排序" clearable />
      </el-form-item>
      <el-form-item label="售卖结束时间" prop="limitTime">
        <el-date-picker
          v-model="editLimitTime"
          type="datetime"
          placeholder="请选择售卖结束时间"
          style="width: 100%"
        />
      </el-form-item>
      <el-form-item label="是否有效" prop="enable">
        <el-radio-group v-model="editFormData.enable" size="small" clearable>
          <el-radio-button :value="true">开放购买</el-radio-button>
          <el-radio-button :value="false">未开放购买</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="是否展示" prop="isDisplay">
        <el-radio-group v-model="editFormData.isDisplay" size="small" clearable>
          <el-radio-button :value="1">开放展示</el-radio-button>
          <el-radio-button :value="0">不开放展示</el-radio-button>
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
import { ref, reactive, onMounted, computed } from "vue";
import dayjs from "dayjs";
import { Refresh } from "@element-plus/icons-vue";
import { ElNotification } from "element-plus";
import Pagination from "@/components/Pangination/Pagination.vue";
import ExportFieldFiltering from "@/components/ExportFieldFiltering/index.vue";
import { downLoadFile } from "@/api/downloadFile.js";
import { nodeInfoList, nodeInfoUpdate } from "@/api/modules/node.js";
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

// 编辑相关
const dialogEdit = ref(false);
const editFormRef = ref(null);
const editLoading = ref(false);
const editLimitTime = ref(null);
const editFormData = ref({
  id: null,
  dividendRatio: "",
  dividendUserUid: "",
  enable: null,
  isDisplay: null,
  limitTime: null,
  personBuyLimit: null,
  rebateRatio: "",
  sort: null,
  supportPaymentSymbols: "",
  upperQuantityLimit: "",
  usdPrice: ""
});
const editRules = ref({
  usdPrice: [{ required: true, message: "请输入单价", trigger: "blur" }],
  supportPaymentSymbols: [{ required: true, message: "请输入支持的支付币种", trigger: "blur" }],
  upperQuantityLimit: [{ required: true, message: "请输入平台最大限制", trigger: "blur" }],
  personBuyLimit: [{ required: true, message: "请输入产品用户购买有效上限", trigger: "blur" }],
  dividendRatio: [{ required: true, message: "请输入分红比例", trigger: "blur" }],
  dividendUserUid: [{ required: true, message: "请输入分红交易所UID", trigger: "blur" }],
  rebateRatio: [{ required: true, message: "请输入最近一层返佣比例", trigger: "blur" }],
  sort: [{ required: true, message: "请输入排序", trigger: "blur" }],
  enable: [{ required: true, message: "请选择是否有效", trigger: "change" }],
  isDisplay: [{ required: true, message: "请选择是否展示", trigger: "change" }]
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
  const res = await nodeInfoList(params);
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
  const res = await nodeInfoList(params);
  try {
    await downLoadFile({ fileName: "node_info", fileUrl: res.data.url });
    ElNotification.success("导出成功");
  } catch {
    ElNotification.error("导出失败");
  }
};

// 打开编辑对话框
const openEdit = row => {
  editFormData.value = {
    id: row.id,
    dividendRatio: row.dividendRatio || "",
    dividendUserUid: row.dividendUserUid || "",
    enable: row.enable !== undefined ? (row.enable === 1 || row.enable === true ? true : false) : null,
    isDisplay: row.isDisplay !== undefined ? (row.isDisplay === 1 || row.isDisplay === true ? 1 : 0) : null,
    limitTime: row.limitTime || null,
    personBuyLimit: row.personBuyLimit || null,
    rebateRatio: row.rebateRatio || "",
    sort: row.sort || null,
    supportPaymentSymbols: row.supportPaymentSymbols || "",
    upperQuantityLimit: row.upperQuantityLimit || "",
    usdPrice: row.usdPrice || ""
  };
  // 设置日期选择器的值（Unix时间戳转换为日期对象）
  editLimitTime.value = row.limitTime ? dayjs.unix(row.limitTime).toDate() : null;
  // 同步到表单数据
  editFormData.value.limitTime = row.limitTime || null;
  dialogEdit.value = true;
};

// 提交编辑
const handleEditSubmit = async () => {
  if (!editFormRef.value) return;
  // 验证日期选择器
  if (!editLimitTime.value) {
    ElNotification.warning("请选择售卖结束时间");
    return;
  }
  await editFormRef.value.validate(async valid => {
    if (!valid) return;
    editLoading.value = true;
    try {
      // 将日期选择器的值转换为Unix时间戳
      const params = {
        ...editFormData.value,
        limitTime: editLimitTime.value ? dayjs(editLimitTime.value).unix() : null,
        personBuyLimit: editFormData.value.personBuyLimit ? Number(editFormData.value.personBuyLimit) : null,
        sort: editFormData.value.sort ? Number(editFormData.value.sort) : null
      };
      const res = await nodeInfoUpdate(params);
      if (res.code === 200) {
        ElNotification.success("编辑成功");
        dialogEdit.value = false;
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
