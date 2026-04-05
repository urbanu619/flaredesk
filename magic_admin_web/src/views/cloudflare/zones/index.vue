<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-select v-model="selectedAccountId" placeholder="选择 CF 账号" @change="loadZones" style="width:200px">
            <el-option v-for="a in accounts" :key="a.id" :label="a.name" :value="a.id" />
          </el-select>
          <el-input v-model="searchName" placeholder="搜索域名" clearable style="width:200px" @input="loadZones" />
          <el-button type="primary" :loading="syncing" :disabled="!selectedAccountId" @click="handleSync">
            从 CF 同步
          </el-button>
          <el-button :disabled="!selectedAccountId" @click="loadZones">刷新</el-button>
          <span v-if="total" class="zone-total">共 {{ total }} 个域名</span>
        </div>
      </template>

      <el-table v-loading="loading" :data="zones" border stripe>
        <el-table-column prop="name" label="域名" min-width="160" sortable />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? '已激活' : row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="暂停" width="80">
          <template #default="{ row }">
            <el-tag :type="row.paused ? 'warning' : 'info'" size="small">{{ row.paused ? '暂停' : '正常' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="planName" label="套餐" width="100" />
        <el-table-column label="NS 服务器" min-width="260">
          <template #default="{ row }">
            <div v-for="ns in parseNs(row.nameServers)" :key="ns" class="ns-item">{{ ns }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="activatedOn" label="激活时间" width="175">
          <template #default="{ row }">{{ formatDate(row.activatedOn) }}</template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="140" show-overflow-tooltip />
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="goToDns(row)">DNS记录</el-button>
            <el-button link type="info" @click="handleRemark(row)">备注</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @change="loadZones"
        />
      </div>
    </el-card>

    <el-dialog v-model="remarkDialog" title="编辑备注" width="400px">
      <el-input v-model="remarkText" type="textarea" :rows="3" placeholder="输入备注" />
      <template #footer>
        <el-button @click="remarkDialog = false">取消</el-button>
        <el-button type="primary" @click="submitRemark">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { getCfAccountList, syncCfZones, getLocalZoneList, updateZoneRemark } from '@/api/modules/cloudflare';

const router = useRouter();
const accounts = ref([]);
const zones = ref([]);
const selectedAccountId = ref(null);
const loading = ref(false);
const syncing = ref(false);
const searchName = ref('');
const page = ref(1);
const pageSize = ref(50);
const total = ref(0);

const remarkDialog = ref(false);
const remarkText = ref('');
const remarkRow = ref(null);

onMounted(async () => {
  const res = await getCfAccountList({ pageSize: 100 });
  accounts.value = res?.data?.list || [];
  if (accounts.value.length === 1) {
    selectedAccountId.value = accounts.value[0].id;
    loadZones();
  }
});

async function loadZones() {
  if (!selectedAccountId.value) return;
  loading.value = true;
  try {
    const res = await getLocalZoneList({
      accountId: selectedAccountId.value,
      name: searchName.value,
      pageNum: page.value,
      pageSize: pageSize.value
    });
    zones.value = res?.data?.list || [];
    total.value = Number(res?.data?.total) || 0;
  } finally {
    loading.value = false;
  }
}

async function handleSync() {
  syncing.value = true;
  try {
    const res = await syncCfZones(selectedAccountId.value);
    ElMessage.success(`同步完成：新增 ${res?.data?.created}，更新 ${res?.data?.updated}，共 ${res?.data?.total} 个域名`);
    page.value = 1;
    loadZones();
  } catch (e) {
    ElMessage.error('同步失败：' + e.message);
  } finally {
    syncing.value = false;
  }
}

function parseNs(nsJson) {
  try { return JSON.parse(nsJson) || []; } catch { return []; }
}

function formatDate(str) {
  if (!str) return '-';
  return new Date(str).toLocaleString('zh-CN', { hour12: false }).replace(/\//g, '-');
}

function goToDns(row) {
  router.push({ name: 'cfDns', query: { accountId: selectedAccountId.value, zoneId: row.zoneId, zoneName: row.name } });
}

function handleRemark(row) {
  remarkRow.value = row;
  remarkText.value = row.remark || '';
  remarkDialog.value = true;
}

async function submitRemark() {
  await updateZoneRemark({ id: remarkRow.value.id, remark: remarkText.value });
  remarkRow.value.remark = remarkText.value;
  remarkDialog.value = false;
  ElMessage.success('备注已保存');
}
</script>

<style scoped lang="scss">
.page-container { padding: 20px; }
.card-header { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.zone-total { color: #666; font-size: 13px; }
.ns-item { font-size: 12px; color: #555; line-height: 1.7; }
.pagination-wrap { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
