<template>
  <div class="group-container">
    <ProTable
      ref="proTableRef"
      :columns="columns"
      :request-api="getGroupList"
      :search-columns="searchColumns"
      :toolbar-buttons="toolbarButtons"
      :show-pagination="true"
    >
      <!-- 状态列 -->
      <template #status="{ row }">
        <el-tag :type="getStatusType(row.status)">
          {{ getStatusText(row.status) }}
        </el-tag>
      </template>

      <!-- 操作列 -->
      <template #actions="{ row }">
        <el-button link type="primary" @click="handleView(row)">查看</el-button>
        <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
        <el-button
          link
          type="danger"
          @click="handleDelete(row)"
          :disabled="row.status === 3"
        >
          删除
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
        label-width="120px"
      >
        <el-form-item label="群组ID" prop="chatId">
          <el-input
            v-model.number="formData.chatId"
            placeholder="请输入Telegram群组ID"
            :disabled="isEdit"
          />
        </el-form-item>
        <el-form-item label="群组标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入群组标题" />
        </el-form-item>
        <el-form-item label="群组用户名">
          <el-input v-model="formData.username" placeholder="请输入群组用户名（@开头）" />
        </el-form-item>
        <el-form-item label="群组类型">
          <el-select v-model="formData.chatType" placeholder="请选择群组类型">
            <el-option label="普通群组" value="group" />
            <el-option label="超级群组" value="supergroup" />
            <el-option label="频道" value="channel" />
          </el-select>
        </el-form-item>
        <el-form-item label="成员数量">
          <el-input-number
            v-model="formData.memberCount"
            :min="0"
            placeholder="成员数量"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="formData.status" placeholder="请选择状态">
            <el-option label="正常" :value="1" />
            <el-option label="已禁用" :value="2" />
            <el-option label="已退出" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="群组描述">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入群组描述"
          />
        </el-form-item>
        <el-form-item label="备注">
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
    <el-dialog v-model="viewDialogVisible" title="群组详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="群组ID">
          {{ viewData.chatId }}
        </el-descriptions-item>
        <el-descriptions-item label="群组标题">
          {{ viewData.title }}
        </el-descriptions-item>
        <el-descriptions-item label="群组用户名">
          {{ viewData.username || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="群组类型">
          {{ getChatTypeText(viewData.chatType) }}
        </el-descriptions-item>
        <el-descriptions-item label="成员数量">
          {{ viewData.memberCount || 0 }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(viewData.status)">
            {{ getStatusText(viewData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="机器人加入时间" :span="2">
          {{ formatTime(viewData.botJoinedAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">
          {{ formatTime(viewData.createdAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="更新时间" :span="2">
          {{ formatTime(viewData.updatedAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="群组描述" :span="2">
          {{ viewData.description || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">
          {{ viewData.remark || '-' }}
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
  getGroupList,
  createGroup,
  updateGroup,
  deleteGroup,
  syncGroupsFromBot
} from '@/api/modules/redpacket';

const proTableRef = ref();
const formRef = ref();
const dialogVisible = ref(false);
const viewDialogVisible = ref(false);
const dialogTitle = ref('');
const isEdit = ref(false);

const formData = reactive({
  id: null,
  chatId: null,
  title: '',
  username: '',
  chatType: 'supergroup',
  memberCount: 0,
  status: 1,
  description: '',
  remark: ''
});

const viewData = ref({});

// 表格列配置
const columns = [
  { prop: 'id', label: 'ID', width: 80 },
  { prop: 'chatId', label: '群组ID', width: 150 },
  { prop: 'title', label: '群组标题', minWidth: 200 },
  { prop: 'username', label: '用户名', width: 150 },
  { prop: 'chatType', label: '类型', width: 100 },
  { prop: 'memberCount', label: '成员数', width: 100 },
  { prop: 'status', label: '状态', width: 100, slot: 'status' }
];

// 搜索列配置
const searchColumns = [
  { prop: 'title', label: '群组标题', type: 'input' },
  { prop: 'username', label: '用户名', type: 'input' },
  {
    prop: 'status',
    label: '状态',
    type: 'select',
    options: [
      { label: '正常', value: 1 },
      { label: '已禁用', value: 2 },
      { label: '已退出', value: 3 }
    ]
  }
];

// 工具栏按钮
const toolbarButtons = [
  {
    label: '新增群组',
    type: 'primary',
    icon: 'Plus',
    onClick: handleAdd
  },
  {
    label: '同步群组',
    type: 'success',
    icon: 'Refresh',
    onClick: handleSync
  }
];

// 表单验证规则
const formRules = {
  chatId: [{ required: true, message: '请输入群组ID', trigger: 'blur' }],
  title: [{ required: true, message: '请输入群组标题', trigger: 'blur' }]
};

// 获取状态文本
const getStatusText = status => {
  const statusMap = {
    0: '已失效',
    1: '正常',
    2: '已禁用',
    3: '已退出'
  };
  return statusMap[status] || '未知';
};

// 获取状态类型
const getStatusType = status => {
  const typeMap = {
    0: 'danger',
    1: 'success',
    2: 'warning',
    3: 'info'
  };
  return typeMap[status] || '';
};

// 获取群组类型文本
const getChatTypeText = type => {
  const typeMap = {
    group: '普通群组',
    supergroup: '超级群组',
    channel: '频道'
  };
  return typeMap[type] || type || '-';
};

// 格式化时间
const formatTime = timestamp => {
  if (!timestamp) return '-';
  const date = new Date(timestamp * 1000);
  return date.toLocaleString('zh-CN');
};

// 新增
function handleAdd() {
  dialogTitle.value = '新增群组';
  isEdit.value = false;
  resetForm();
  dialogVisible.value = true;
}

// 编辑
function handleEdit(row) {
  dialogTitle.value = '编辑群组';
  isEdit.value = true;
  Object.assign(formData, row);
  dialogVisible.value = true;
}

// 查看
function handleView(row) {
  viewData.value = { ...row };
  viewDialogVisible.value = true;
}

// 删除
async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(
      `确定要删除群组"${row.title}"吗？此操作会将群组状态标记为已退出。`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    );

    await deleteGroup(row.id);
    ElMessage.success('删除成功');
    proTableRef.value.refresh();
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败');
    }
  }
}

// 提交表单
async function handleSubmit() {
  try {
    await formRef.value.validate();

    const api = isEdit.value ? updateGroup : createGroup;
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

// 同步群组
async function handleSync() {
  try {
    await ElMessageBox.confirm(
      '确定要从Telegram Bot同步群组信息吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }
    );

    await syncGroupsFromBot();
    ElMessage.success('同步成功');
    proTableRef.value.refresh();
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '同步失败');
    }
  }
}

// 重置表单
function resetForm() {
  Object.assign(formData, {
    id: null,
    chatId: null,
    title: '',
    username: '',
    chatType: 'supergroup',
    memberCount: 0,
    status: 1,
    description: '',
    remark: ''
  });
  formRef.value?.clearValidate();
}

// 关闭对话框
function handleDialogClose() {
  resetForm();
}
</script>

<style scoped lang="scss">
.group-container {
  padding: 20px;
}
</style>
