<template>
  <div class="red-packet-config">
    <ProTable
      ref="proTable"
      :columns="columns"
      :request-api="requestApi"
      :search-columns="searchColumns"
      :toolbar-buttons="toolbarButtons"
      :show-pagination="true"
      :actions-width="300"
    >
      <!-- 配置类型 -->
      <template #configType="{ row }">
        <el-tag :type="row.configType === 1 ? 'success' : 'info'">
          {{ row.configType === 1 ? '定时红包' : '手动触发' }}
        </el-tag>
      </template>

      <!-- 红包类型 -->
      <template #packetType="{ row }">
        <el-tag :type="row.packetType === 1 ? '' : 'warning'">
          {{ row.packetType === 1 ? '普通红包' : '手气红包' }}
        </el-tag>
      </template>

      <!-- 状态 -->
      <template #status="{ row }">
        <el-switch
          v-model="row.status"
          :active-value="1"
          :inactive-value="2"
          @change="handleStatusChange(row)"
        />
      </template>

      <!-- 操作 -->
      <template #actions="{ row }">
        <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
        <el-button link type="primary" @click="handleView(row)">查看</el-button>
        <el-button
          v-if="row.configType === 2"
          link
          type="success"
          @click="handleSend(row)"
        >
          立即发送
        </el-button>
        <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
      </template>
    </ProTable>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px">
        <el-form-item label="配置名称" prop="configName">
          <el-input v-model="formData.configName" placeholder="请输入配置名称" />
        </el-form-item>

        <el-form-item label="配置类型" prop="configType">
          <el-radio-group v-model="formData.configType">
            <el-radio :label="1">定时红包</el-radio>
            <el-radio :label="2">手动触发</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="选择群组" prop="groupId">
          <el-select
            v-model="formData.groupId"
            placeholder="请选择Telegram群组"
            filterable
            style="width: 100%"
            @change="handleGroupChange"
          >
            <el-option
              v-for="group in groupList"
              :key="group.chatId"
              :label="group.title"
              :value="String(group.chatId)"
            >
              <span>{{ group.title }}</span>
              <span style="float: right; color: #999; font-size: 12px">{{ group.chatId }}</span>
            </el-option>
          </el-select>
          <div v-if="formData.groupId" class="form-tip">群组ID: {{ formData.groupId }}</div>
        </el-form-item>

        <el-form-item label="红包类型" prop="packetType">
          <el-radio-group v-model="formData.packetType">
            <el-radio :label="1">普通红包</el-radio>
            <el-radio :label="2">手气红包</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="红包金额" prop="totalAmount">
          <el-input-number
            v-model="formData.totalAmount"
            :min="0"
            :precision="2"
            placeholder="请输入红包总金额"
          />
          <span class="ml-2">{{ formData.symbol }}</span>
        </el-form-item>

        <el-form-item label="红包个数" prop="totalCount">
          <el-input-number
            v-model="formData.totalCount"
            :min="1"
            :max="100"
            placeholder="请输入红包个数"
          />
        </el-form-item>

        <el-form-item label="消息语言" prop="lang">
          <el-select v-model="formData.lang" placeholder="请选择消息语言" style="width: 100%">
            <el-option label="🇻🇳 Tiếng Việt (越南语)" value="vi" />
            <el-option label="🇮🇩 Bahasa Indonesia (印尼语)" value="id" />
            <el-option label="🇺🇸 English (英语)" value="en" />
            <el-option label="🇨🇳 中文" value="zh" />
          </el-select>
        </el-form-item>

        <el-form-item label="过期时间" prop="expireMinutes">
          <el-input-number
            v-model="formData.expireMinutes"
            :min="1"
            :max="1440"
            placeholder="红包过期时间"
          />
          <span class="ml-2">分钟</span>
        </el-form-item>

        <el-form-item v-if="formData.packetType === 2" label="单人最大金额" prop="maxGrabAmount">
          <el-input-number
            v-model="formData.maxGrabAmount"
            :min="0"
            :precision="2"
            placeholder="0=不限制"
          />
          <span class="ml-2">{{ formData.symbol }}（0=不限制）</span>
        </el-form-item>

        <el-form-item label="祝福语" prop="blessingWords">
          <el-input
            v-model="formData.blessingWords"
            type="textarea"
            :rows="3"
            placeholder="请输入祝福语"
          />
        </el-form-item>

        <!-- 定时红包特有字段 -->
        <template v-if="formData.configType === 1">
          <el-form-item label="Cron表达式" prop="cronExpr">
            <el-input
              v-model="formData.cronExpr"
              placeholder="例如: 0 0 12 * * * (每天12点)"
            />
            <div class="form-tip">
              格式: 秒 分 时 日 月 周<br>
              示例: 0 0 12 * * * 表示每天中午12点
            </div>
          </el-form-item>

          <el-form-item label="开始日期" prop="startDate">
            <el-date-picker
              v-model="formData.startDate"
              type="datetime"
              placeholder="选择开始日期"
            />
          </el-form-item>

          <el-form-item label="结束日期" prop="endDate">
            <el-date-picker
              v-model="formData.endDate"
              type="datetime"
              placeholder="选择结束日期"
            />
          </el-form-item>
        </template>

        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="formData.remark"
            type="textarea"
            :rows="2"
            placeholder="请输入备注"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 查看对话框 -->
    <el-dialog v-model="viewDialogVisible" title="配置详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ viewData.id }}</el-descriptions-item>
        <el-descriptions-item label="配置名称">{{ viewData.configName }}</el-descriptions-item>
        <el-descriptions-item label="配置类型">
          <el-tag :type="viewData.configType === 1 ? 'success' : 'info'">
            {{ viewData.configType === 1 ? '定时红包' : '手动触发' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="红包类型">
          <el-tag :type="viewData.packetType === 1 ? '' : 'warning'">
            {{ viewData.packetType === 1 ? '普通红包' : '手气红包' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="群组ID">{{ viewData.groupId }}</el-descriptions-item>
        <el-descriptions-item label="群组名称">{{ viewData.groupName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="红包金额">{{ viewData.totalAmount }} {{ viewData.symbol || 'VND' }}</el-descriptions-item>
        <el-descriptions-item label="红包个数">{{ viewData.totalCount }}</el-descriptions-item>
        <el-descriptions-item label="消息语言">{{ {vi:'🇻🇳 越南语',id:'🇮🇩 印尼语',en:'🇺🇸 英语',zh:'🇨🇳 中文'}[viewData.lang] || '🇻🇳 越南语' }}</el-descriptions-item>
        <el-descriptions-item label="过期时间">{{ viewData.expireMinutes || 10 }} 分钟</el-descriptions-item>
        <el-descriptions-item label="单人上限">{{ viewData.maxGrabAmount > 0 ? viewData.maxGrabAmount + ' ' + (viewData.symbol || 'VND') : '不限制' }}</el-descriptions-item>
        <el-descriptions-item label="祝福语" :span="2">{{ viewData.blessingWords || '-' }}</el-descriptions-item>
        <el-descriptions-item v-if="viewData.configType === 1" label="Cron表达式">{{ viewData.cronExpr || '-' }}</el-descriptions-item>
        <el-descriptions-item v-if="viewData.configType === 1" label="开始日期">{{ viewData.startDate ? new Date(viewData.startDate).toLocaleString() : '-' }}</el-descriptions-item>
        <el-descriptions-item v-if="viewData.configType === 1" label="结束日期">{{ viewData.endDate ? new Date(viewData.endDate).toLocaleString() : '-' }}</el-descriptions-item>
        <el-descriptions-item label="执行次数">{{ viewData.execCount || 0 }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ viewData.remark || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import ProTable from '@/components/ProTable/index.vue';
import {
  getRedPacketConfigList,
  createRedPacketConfig,
  updateRedPacketConfig,
  deleteRedPacketConfig,
  toggleRedPacketConfigStatus,
  sendRedPacketManual,
  getGroupList
} from '@/api/modules/redpacket';

// 表格实例
const proTable = ref();

// 群组列表
const groupList = ref([]);

// 加载群组列表
async function loadGroupList() {
  try {
    const res = await getGroupList({ current: 1, pageSize: 200, status: 1 });
    if (res.data && res.data.list) {
      groupList.value = res.data.list;
    } else if (res.data && Array.isArray(res.data)) {
      groupList.value = res.data;
    }
  } catch (e) {
    console.error('加载群组列表失败', e);
  }
}

// 语言 → 默认币种映射
const langSymbolMap = { vi: 'VND', id: 'IDR', en: 'USD', zh: 'CNY' };

// 语言切换时自动联动币种
watch(() => formData.lang, (lang) => {
  formData.symbol = langSymbolMap[lang] || 'VND';
});

onMounted(() => {
  loadGroupList();
});

// 群组选择变化 — 自动填充 groupName
function handleGroupChange(chatId) {
  const group = groupList.value.find(g => String(g.chatId) === String(chatId));
  formData.groupName = group ? group.title : '';
}

// 后端返回格式适配
const requestApi = async params => {
  const res = await getRedPacketConfigList(params);
  if (res.data) {
    return { data: { list: res.data.list, paging: res.data.paging } };
  }
  return res;
};

// 表格列配置
const columns = [
  { prop: 'id', label: 'ID', width: 80 },
  { prop: 'configName', label: '配置名称', minWidth: 120 },
  { prop: 'configType', label: '配置类型', width: 100, slot: 'configType' },
  { prop: 'groupName', label: '群组名称', width: 140 },
  { prop: 'packetType', label: '红包类型', width: 100, slot: 'packetType' },
  { prop: 'totalAmount', label: '红包金额', width: 100 },
  { prop: 'totalCount', label: '红包个数', width: 80 },
  { prop: 'expireMinutes', label: '过期(分钟)', width: 100 },
  { prop: 'maxGrabAmount', label: '单人上限', width: 100 },
  { prop: 'lang', label: '语言', width: 80 },
  { prop: 'cronExpr', label: 'Cron表达式', width: 140 },
  { prop: 'execCount', label: '执行次数', width: 80 },
  { prop: 'status', label: '状态', width: 80, slot: 'status' }
];

// 搜索列配置
const searchColumns = [
  { prop: 'configName', label: '配置名称', type: 'input' },
  {
    prop: 'configType',
    label: '配置类型',
    type: 'select',
    options: [
      { label: '定时红包', value: 1 },
      { label: '手动触发', value: 2 }
    ]
  },
  {
    prop: 'packetType',
    label: '红包类型',
    type: 'select',
    options: [
      { label: '普通红包', value: 1 },
      { label: '手气红包', value: 2 }
    ]
  }
];

// 工具栏按钮
const toolbarButtons = [
  {
    label: '新增配置',
    type: 'primary',
    icon: 'Plus',
    onClick: handleAdd
  }
];

// 对话框相关
const dialogVisible = ref(false);
const dialogTitle = ref('');
const viewDialogVisible = ref(false);
const viewData = ref({});
const formRef = ref();
const formData = reactive({
  id: null,
  configName: '',
  configType: 1,
  groupId: '',
  groupName: '',
  packetType: 1,
  totalAmount: 0,
  totalCount: 1,
  symbol: 'VND',
  lang: 'vi',
  expireMinutes: 10,
  maxGrabAmount: 0,
  blessingWords: '',
  cronExpr: '',
  startDate: null,
  endDate: null,
  remark: ''
});

// 表单验证规则
const rules = {
  configName: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
  configType: [{ required: true, message: '请选择配置类型', trigger: 'change' }],
  groupId: [{ required: true, message: '请选择群组', trigger: 'change' }],
  packetType: [{ required: true, message: '请选择红包类型', trigger: 'change' }],
  totalAmount: [{ required: true, message: '请输入红包金额', trigger: 'blur' }],
  totalCount: [{ required: true, message: '请输入红包个数', trigger: 'blur' }],
  cronExpr: [
    {
      trigger: 'blur',
      validator: (rule, value, callback) => {
        if (formData.configType === 1 && !value) {
          callback(new Error('定时红包必须配置Cron表达式'));
        } else {
          callback();
        }
      }
    }
  ]
};

// 新增
function handleAdd() {
  dialogTitle.value = '新增红包配置';
  Object.assign(formData, {
    id: null,
    configName: '',
    configType: 1,
    groupId: '',
    groupName: '',
    packetType: 1,
    totalAmount: 0,
    totalCount: 1,
    symbol: 'VND',
    lang: 'vi',
    expireMinutes: 10,
    maxGrabAmount: 0,
    blessingWords: '',
    cronExpr: '',
    startDate: null,
    endDate: null,
    remark: ''
  });
  dialogVisible.value = true;
}

// 辅助函数：检测是否为有效日期（排除零值时间戳）
function isValidDate(val) {
  if (!val) return false;
  const d = new Date(val);
  return !isNaN(d.getTime()) && d.getFullYear() >= 2000;
}

// 编辑 — 只复制表单相关字段
function handleEdit(row) {
  dialogTitle.value = '编辑红包配置';
  Object.assign(formData, {
    id: row.id,
    configName: row.configName,
    configType: row.configType,
    groupId: String(row.groupId),
    groupName: row.groupName,
    packetType: row.packetType,
    totalAmount: row.totalAmount,
    totalCount: row.totalCount,
    symbol: row.symbol || 'VND',
    lang: row.lang || 'vi',
    expireMinutes: row.expireMinutes || 10,
    maxGrabAmount: row.maxGrabAmount || 0,
    blessingWords: row.blessingWords,
    cronExpr: row.cronExpr,
    startDate: isValidDate(row.startDate) ? row.startDate : null,
    endDate: isValidDate(row.endDate) ? row.endDate : null,
    remark: row.remark
  });
  dialogVisible.value = true;
}

// 查看
function handleView(row) {
  const data = { ...row };
  data.startDate = isValidDate(data.startDate) ? data.startDate : null;
  data.endDate = isValidDate(data.endDate) ? data.endDate : null;
  viewData.value = data;
  viewDialogVisible.value = true;
}

// 删除
function handleDelete(row) {
  ElMessageBox.confirm('确定要删除这条配置吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    await deleteRedPacketConfig(row.id);
    ElMessage.success('删除成功');
    proTable.value.refresh();
  });
}

// 状态切换
async function handleStatusChange(row) {
  if (!row.id) return;
  try {
    await toggleRedPacketConfigStatus(row.id, row.status);
    ElMessage.success(row.status === 1 ? '启用成功' : '禁用成功');
    proTable.value.refresh();
  } catch (error) {
    row.status = row.status === 1 ? 2 : 1;
    ElMessage.error('操作失败');
  }
}

// 立即发送（仅手动触发类型）
function handleSend(row) {
  ElMessageBox.confirm('确定要立即发送这个红包吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    await sendRedPacketManual(row.id);
    ElMessage.success('红包发送成功');
    proTable.value.refresh();
  });
}

// 提交表单
function handleSubmit() {
  formRef.value.validate(async valid => {
    if (!valid) return;

    try {
      // 只提交表单相关字段
      const data = {
        id: formData.id,
        configName: formData.configName,
        configType: formData.configType,
        groupId: formData.groupId,
        groupName: formData.groupName,
        packetType: formData.packetType,
        totalAmount: formData.totalAmount,
        totalCount: formData.totalCount,
        symbol: formData.symbol || 'VND',
        lang: formData.lang || 'vi',
        expireMinutes: formData.expireMinutes || 10,
        maxGrabAmount: formData.maxGrabAmount || 0,
        blessingWords: formData.blessingWords,
        cronExpr: formData.cronExpr,
        startDate: formData.startDate,
        endDate: formData.endDate,
        remark: formData.remark
      };
      // 后端 time.Time 类型需要 ISO 8601 格式字符串
      if (data.startDate) {
        data.startDate = new Date(data.startDate).toISOString();
      } else {
        data.startDate = null;
      }
      if (data.endDate) {
        data.endDate = new Date(data.endDate).toISOString();
      } else {
        data.endDate = null;
      }

      if (data.id) {
        await updateRedPacketConfig(data);
        ElMessage.success('更新成功');
      } else {
        await createRedPacketConfig(data);
        ElMessage.success('创建成功');
      }

      dialogVisible.value = false;
      proTable.value.refresh();
    } catch (error) {
      ElMessage.error(error.message || '操作失败');
    }
  });
}

// 对话框关闭
function handleDialogClose() {
  formRef.value?.resetFields();
}
</script>

<style scoped lang="scss">
.red-packet-config {
  padding: 20px;

  .form-tip {
    font-size: 12px;
    color: #999;
    margin-top: 5px;
    line-height: 1.5;
  }

  .ml-2 {
    margin-left: 8px;
  }
}
</style>
