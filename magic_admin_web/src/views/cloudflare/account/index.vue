<template>
  <div class="page-container">
    <ProTable
      ref="proTableRef"
      :columns="columns"
      :request-api="getCfAccountList"
      :search-columns="searchColumns"
      :toolbar-buttons="toolbarButtons"
    >
      <template #status="{ row }">
        <el-tag :type="row.status === 1 ? 'success' : 'danger'">
          {{ row.status === 1 ? '正常' : '禁用' }}
        </el-tag>
      </template>
      <template #apiToken="{ row }">
        <span>{{ row.apiToken ? row.apiToken.slice(0, 8) + '••••••••' : '-' }}</span>
      </template>
      <template #actions="{ row }">
        <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
        <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
      </template>
    </ProTable>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑账号' : '新增账号'" width="520px">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="110px">
        <el-form-item label="账号名称" prop="name">
          <el-input v-model="formData.name" placeholder="便于识别的名称" />
        </el-form-item>
        <el-form-item label="CF 邮箱">
          <el-input v-model="formData.email" placeholder="Cloudflare 账号邮箱（选填）" />
        </el-form-item>
        <el-form-item label="API Token" prop="apiToken">
          <el-input v-model="formData.apiToken" type="password" show-password placeholder="Cloudflare API Token" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formData.remark" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item v-if="isEdit" label="状态">
          <el-select v-model="formData.status">
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="2" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import ProTable from '@/components/ProTable/index.vue';
import { getCfAccountList, createCfAccount, updateCfAccount, deleteCfAccount } from '@/api/modules/cloudflare';

const proTableRef = ref();
const formRef = ref();
const dialogVisible = ref(false);
const isEdit = ref(false);

const formData = reactive({ id: null, name: '', email: '', apiToken: '', remark: '', status: 1 });

const rules = {
  name: [{ required: true, message: '请输入账号名称', trigger: 'blur' }],
  apiToken: [{ required: true, message: '请输入 API Token', trigger: 'blur' }]
};

const columns = [
  { prop: 'id', label: 'ID', width: 70 },
  { prop: 'name', label: '账号名称', minWidth: 150 },
  { prop: 'email', label: '邮箱', minWidth: 180 },
  { prop: 'apiToken', label: 'API Token', minWidth: 160, slot: 'apiToken' },
  { prop: 'status', label: '状态', width: 90, slot: 'status' },
  { prop: 'remark', label: '备注', minWidth: 150 },
  { prop: 'actions', label: '操作', width: 130, fixed: 'right', slot: 'actions' }
];

const searchColumns = [
  { prop: 'name', label: '账号名称', type: 'input' },
  { prop: 'status', label: '状态', type: 'select', options: [{ label: '正常', value: 1 }, { label: '禁用', value: 2 }] }
];

const toolbarButtons = [{ label: '新增账号', type: 'primary', icon: 'Plus', onClick: handleAdd }];

function handleAdd() {
  isEdit.value = false;
  Object.assign(formData, { id: null, name: '', email: '', apiToken: '', remark: '', status: 1 });
  dialogVisible.value = true;
}

function handleEdit(row) {
  isEdit.value = true;
  Object.assign(formData, { ...row, apiToken: '' });
  dialogVisible.value = true;
}

async function handleDelete(row) {
  await ElMessageBox.confirm(`确定删除账号「${row.name}」？`, '提示', { type: 'warning' });
  await deleteCfAccount(row.id);
  ElMessage.success('删除成功');
  proTableRef.value.refresh();
}

async function handleSubmit() {
  await formRef.value.validate();
  const api = isEdit.value ? updateCfAccount : createCfAccount;
  await api(formData);
  ElMessage.success(isEdit.value ? '更新成功' : '创建成功');
  dialogVisible.value = false;
  proTableRef.value.refresh();
}
</script>

<style scoped lang="scss">
.page-container { padding: 20px; }
</style>
