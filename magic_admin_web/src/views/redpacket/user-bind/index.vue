<template>
  <div class="user-bind-container">
    <ProTable
      ref="proTableRef"
      :columns="columns"
      :request-api="getUserBindList"
      :search-columns="searchColumns"
      :toolbar-buttons="toolbarButtons"
      :show-pagination="true"
    >
      <!-- 绑定状态列 -->
      <template #bindStatus="{ row }">
        <el-tag :type="getStatusType(row.bindStatus)">
          {{ getStatusText(row.bindStatus) }}
        </el-tag>
      </template>

      <!-- 绑定时间列 -->
      <template #bindTime="{ row }">
        {{ formatTime(row.bindTime) }}
      </template>

      <!-- 操作列 -->
      <template #actions="{ row }">
        <el-button link type="primary" @click="handleView(row)">查看</el-button>
        <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
        <el-button
          link
          type="danger"
          @click="handleDelete(row)"
          :disabled="row.bindStatus === 0"
        >
          解绑
        </el-button>
      </template>
    </ProTable>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="140px"
      >
        <el-form-item label="平台用户ID" prop="userId">
          <el-input
            v-model.number="formData.userId"
            placeholder="请输入平台用户ID"
          />
        </el-form-item>
        <el-form-item label="Telegram ID" prop="telegramId">
          <el-input
            v-model.number="formData.telegramId"
            placeholder="请输入Telegram用户ID（可选）"
          />
        </el-form-item>
        <el-form-item label="Telegram用户名" prop="telegramUsername">
          <el-input
            v-model="formData.telegramUsername"
            placeholder="请输入Telegram用户名（不含@）"
          />
        </el-form-item>
        <el-form-item label="Telegram名字">
          <el-input
            v-model="formData.telegramFirstName"
            placeholder="请输入Telegram名字（可选）"
          />
        </el-form-item>
        <el-form-item label="绑定状态" v-if="isEdit">
          <el-select v-model="formData.bindStatus" placeholder="请选择状态">
            <el-option label="已绑定" :value="1" />
            <el-option label="已解绑" :value="0" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 查看对话框 -->
    <el-dialog v-model="viewDialogVisible" title="绑定详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">
          {{ viewData.id }}
        </el-descriptions-item>
        <el-descriptions-item label="平台用户ID">
          {{ viewData.userId }}
        </el-descriptions-item>
        <el-descriptions-item label="Telegram ID">
          {{ viewData.telegramId || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="Telegram用户名">
          {{ viewData.telegramUsername || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="Telegram名字">
          {{ viewData.telegramFirstName || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="绑定状态">
          <el-tag :type="getStatusType(viewData.bindStatus)">
            {{ getStatusText(viewData.bindStatus) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="绑定时间" :span="2">
          {{ formatTime(viewData.bindTime) }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">
          {{ formatTimestamp(viewData.createdAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="更新时间" :span="2">
          {{ formatTimestamp(viewData.updatedAt) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import ProTable from '@/components/ProTable/index.vue';
import {
  getUserBindList,
  createUserBind,
  updateUserBind,
  deleteUserBind
} from '@/api/modules/redpacket';

const proTableRef = ref();
const formRef = ref();
const dialogVisible = ref(false);
const viewDialogVisible = ref(false);
const dialogTitle = ref('');
const isEdit = ref(false);

const formData = reactive({
  id: null,
  userId: null,
  telegramId: null,
  telegramUsername: '',
  telegramFirstName: '',
  bindStatus: 1
});

const viewData = ref({});

// 表格列配置
const columns = [
  { prop: 'id', label: 'ID', width: 80 },
  { prop: 'userId', label: '平台用户ID', width: 120 },
  { prop: 'telegramId', label: 'Telegram ID', width: 150 },
  { prop: 'telegramUsername', label: 'Telegram用户名', minWidth: 150 },
  { prop: 'telegramFirstName', label: 'Telegram名字', width: 150 },
  { prop: 'bindStatus', label: '绑定状态', width: 100, slot: 'bindStatus' },
  { prop: 'bindTime', label: '绑定时间', width: 180, slot: 'bindTime' }
];

// 搜索列配置
const searchColumns = [
  { prop: 'user_id', label: '平台用户ID', type: 'input' },
  { prop: 'telegram_username', label: 'Telegram用户名', type: 'input' },
  {
    prop: 'bind_status',
    label: '绑定状态',
    type: 'select',
    options: [
      { label: '已绑定', value: 1 },
      { label: '已解绑', value: 0 }
    ]
  }
];

// 工具栏按钮
const toolbarButtons = [
  {
    label: '新增绑定',
    type: 'primary',
    icon: 'Plus',
    onClick: handleAdd
  }
];

// 表单验证规则
const formRules = {
  userId: [{ required: true, message: '请输入平台用户ID', trigger: 'blur' }],
  telegramUsername: [{ required: true, message: '请输入Telegram用户名', trigger: 'blur' }]
};

// 获取状态文本
const getStatusText = status => {
  const statusMap = {
    1: '已绑定',
    0: '已解绑'
  };
  return statusMap[status] !== undefined ? statusMap[status] : '未知';
};

// 获取状态类型
const getStatusType = status => {
  const typeMap = {
    1: 'success',
    0: 'info'
  };
  return typeMap[status] || '';
};

// 格式化时间戳（秒）
const formatTimestamp = timestamp => {
  if (!timestamp) return '-';
  const date = new Date(timestamp * 1000);
  return date.toLocaleString('zh-CN');
};

// 格式化时间（ISO 字符串或时间戳）
const formatTime = time => {
  if (!time) return '-';
  const date = new Date(time);
  if (isNaN(date.getTime())) return '-';
  return date.toLocaleString('zh-CN');
};

// 新增
function handleAdd() {
  dialogTitle.value = '新增绑定';
  isEdit.value = false;
  resetForm();
  dialogVisible.value = true;
}

// 编辑
function handleEdit(row) {
  dialogTitle.value = '编辑绑定';
  isEdit.value = true;
  Object.assign(formData, {
    id: row.id,
    userId: row.userId,
    telegramId: row.telegramId,
    telegramUsername: row.telegramUsername,
    telegramFirstName: row.telegramFirstName,
    bindStatus: row.bindStatus
  });
  dialogVisible.value = true;
}

// 查看
function handleView(row) {
  viewData.value = { ...row };
  viewDialogVisible.value = true;
}

// 解绑
async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(
      `确定要解绑用户 "${row.telegramUsername}" 吗？此操作会将绑定状态改为已解绑。`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    );

    await deleteUserBind(row.id);
    ElMessage.success('解绑成功');
    proTableRef.value.refresh();
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '解绑失败');
    }
  }
}

// 提交表单
async function handleSubmit() {
  try {
    await formRef.value.validate();

    const api = isEdit.value ? updateUserBind : createUserBind;
    await api(formData);

    ElMessage.success(isEdit.value ? '更新成功' : '创建成功');
    dialogVisible.value = false;
    proTableRef.value.refresh();
  } catch (error) {
    if (error !== false) {
      ElMessage.error(error.message || '操作失败');
    }
  }
}

// 重置表单
function resetForm() {
  Object.assign(formData, {
    id: null,
    userId: null,
    telegramId: null,
    telegramUsername: '',
    telegramFirstName: '',
    bindStatus: 1
  });
  formRef.value?.clearValidate();
}

// 关闭对话框
function handleDialogClose() {
  resetForm();
}
</script>

<style scoped lang="scss">
.user-bind-container {
  padding: 20px;
}
</style>
