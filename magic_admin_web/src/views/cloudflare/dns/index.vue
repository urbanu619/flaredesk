<template>
  <div class="page-container">
    <!-- 顶部选择器 -->
    <el-card class="selector-card">
      <el-row :gutter="12" align="middle">
        <el-col :span="4">
          <el-select v-model="selectedAccountId" placeholder="选择 CF 账号" @change="onAccountChange" style="width:100%">
            <el-option v-for="a in accounts" :key="a.id" :label="a.name" :value="a.id" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select v-model="selectedZoneId" placeholder="选择域名 (Zone)" :disabled="!selectedAccountId" @change="loadRecords" filterable style="width:100%">
            <el-option v-for="z in zones" :key="z.id" :label="z.name" :value="z.id" />
          </el-select>
        </el-col>
        <el-col :span="14" style="display:flex;flex-wrap:wrap;gap:6px">
          <el-button type="primary" :disabled="!selectedZoneId" @click="showAddDialog">新增记录</el-button>
          <el-button :disabled="!selectedZoneId" @click="showSingleBatchDialog">批量新增</el-button>
          <el-button type="success" :disabled="zones.length === 0" @click="showCrossBatchDialog">跨域名批量新增</el-button>
          <el-button type="danger" :disabled="zones.length === 0" @click="showCrossDeleteDialog">跨域名批量删除</el-button>
          <el-button type="warning" :disabled="zones.length === 0" @click="showCrossProxyDialog">跨域名橙云切换</el-button>
          <el-button :disabled="!selectedZoneId" @click="loadRecords">刷新</el-button>
          <el-button type="warning" plain :disabled="selectedIds.length === 0" @click="showBatchProxy(true)">批量开橙云</el-button>
          <el-button type="info" plain :disabled="selectedIds.length === 0" @click="showBatchProxy(false)">批量关橙云</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- DNS 记录表格 -->
    <el-card class="table-card">
      <el-table v-loading="loading" :data="records" @selection-change="onSelectionChange" border stripe>
        <el-table-column type="selection" width="50" />
        <el-table-column prop="type" label="类型" width="80" />
        <el-table-column prop="name" label="名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip />
        <el-table-column prop="ttl" label="TTL" width="80">
          <template #default="{ row }">{{ row.ttl === 1 ? 'Auto' : row.ttl }}</template>
        </el-table-column>
        <el-table-column label="橙云" width="80">
          <template #default="{ row }">
            <el-tag :type="row.proxied ? 'warning' : 'info'" size="small">{{ row.proxied ? '开启' : '关闭' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="batch-tip" v-if="selectedIds.length > 0">已选 {{ selectedIds.length }} 条记录</div>
    </el-card>

    <!-- 新增/编辑记录 -->
    <el-dialog v-model="recordDialog" :title="editRecord ? '编辑 DNS 记录' : '新增 DNS 记录'" width="520px">
      <el-form ref="recordFormRef" :model="recordForm" :rules="recordRules" label-width="90px">
        <el-form-item label="类型" prop="type">
          <el-select v-model="recordForm.type" style="width:100%">
            <el-option v-for="t in dnsTypes" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称" prop="name"><el-input v-model="recordForm.name" placeholder="如 @ 或 sub" /></el-form-item>
        <el-form-item label="内容" prop="content"><el-input v-model="recordForm.content" placeholder="如 1.2.3.4" /></el-form-item>
        <el-form-item label="TTL">
          <el-select v-model="recordForm.ttl" style="width:100%">
            <el-option label="Auto" :value="1" /><el-option label="60s" :value="60" />
            <el-option label="5min" :value="300" /><el-option label="1h" :value="3600" />
          </el-select>
        </el-form-item>
        <el-form-item label="橙云代理"><el-switch v-model="recordForm.proxied" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="recordDialog = false">取消</el-button>
        <el-button type="primary" @click="submitRecord">确定</el-button>
      </template>
    </el-dialog>

    <!-- 当前域名批量新增 -->
    <el-dialog v-model="batchDialog" title="批量新增 DNS 记录（当前域名）" width="660px">
      <p class="tip-text">每行一条，格式：<code>类型 名称 内容 [proxied]</code></p>
      <pre class="tip-example">A  @  1.2.3.4  true
CNAME  www  example.com</pre>
      <el-input v-model="batchText" type="textarea" :rows="10" placeholder="A  @  1.2.3.4  true" />
      <template #footer>
        <el-button @click="batchDialog = false">取消</el-button>
        <el-button type="primary" @click="submitBatch">批量提交</el-button>
      </template>
    </el-dialog>

    <!-- 跨域名批量新增 -->
    <el-dialog v-model="crossBatchDialog" title="跨域名批量新增 DNS 记录" width="760px">
      <el-form label-width="90px">
        <el-form-item label="选择域名">
          <div style="width:100%">
            <div class="zone-actions">
              <el-button size="small" @click="crossBatchZoneIds = zones.map(z=>z.id)">全选</el-button>
              <el-button size="small" @click="crossBatchZoneIds = []">清空</el-button>
              <span class="zone-count">已选 {{ crossBatchZoneIds.length }} / {{ zones.length }}</span>
            </div>
            <el-select v-model="crossBatchZoneIds" multiple filterable collapse-tags collapse-tags-tooltip placeholder="选择域名" style="width:100%">
              <el-option v-for="z in zones" :key="z.id" :label="z.name" :value="z.id" />
            </el-select>
          </div>
        </el-form-item>
        <el-form-item label="DNS 记录">
          <p class="tip-text" style="margin:0 0 4px">每行一条，格式：<code>类型 名称 内容 [proxied]</code></p>
          <el-input v-model="crossBatchText" type="textarea" :rows="7" placeholder="A  @  1.2.3.4  true" style="width:100%" />
        </el-form-item>
      </el-form>
      <ProgressBlock v-if="progress.show" :progress="progress" />
      <template #footer>
        <el-button @click="crossBatchDialog = false" :disabled="progress.running">取消</el-button>
        <el-button type="primary" @click="submitCrossBatch" :loading="progress.running">开始添加</el-button>
      </template>
    </el-dialog>

    <!-- 跨域名批量删除 -->
    <el-dialog v-model="crossDeleteDialog" title="跨域名批量删除 DNS 记录" width="660px">
      <el-alert type="warning" :closable="false" style="margin-bottom:16px">
        <template #title>此操作将删除所选域名中匹配条件的所有 DNS 记录，不可恢复！</template>
      </el-alert>
      <el-form label-width="90px">
        <el-form-item label="选择域名">
          <div style="width:100%">
            <div class="zone-actions">
              <el-button size="small" @click="crossDeleteZoneIds = zones.map(z=>z.id)">全选</el-button>
              <el-button size="small" @click="crossDeleteZoneIds = []">清空</el-button>
              <span class="zone-count">已选 {{ crossDeleteZoneIds.length }} / {{ zones.length }}</span>
            </div>
            <el-select v-model="crossDeleteZoneIds" multiple filterable collapse-tags collapse-tags-tooltip placeholder="选择域名" style="width:100%">
              <el-option v-for="z in zones" :key="z.id" :label="z.name" :value="z.id" />
            </el-select>
          </div>
        </el-form-item>
        <el-form-item label="记录类型"><el-select v-model="crossDeleteType" placeholder="不限" clearable style="width:100%"><el-option v-for="t in dnsTypes" :key="t" :label="t" :value="t" /></el-select></el-form-item>
        <el-form-item label="记录名称"><el-input v-model="crossDeleteName" placeholder="如 @ 或 www，不填则匹配所有" /></el-form-item>
      </el-form>
      <ProgressBlock v-if="delProgress.show" :progress="delProgress" />
      <template #footer>
        <el-button @click="crossDeleteDialog = false" :disabled="delProgress.running">取消</el-button>
        <el-button type="danger" @click="submitCrossDelete" :loading="delProgress.running">确认删除</el-button>
      </template>
    </el-dialog>

    <!-- 跨域名橙云切换 -->
    <el-dialog v-model="crossProxyDialog" title="跨域名批量切换橙云" width="660px">
      <el-form label-width="90px">
        <el-form-item label="选择域名">
          <div style="width:100%">
            <div class="zone-actions">
              <el-button size="small" @click="crossProxyZoneIds = zones.map(z=>z.id)">全选</el-button>
              <el-button size="small" @click="crossProxyZoneIds = []">清空</el-button>
              <span class="zone-count">已选 {{ crossProxyZoneIds.length }} / {{ zones.length }}</span>
            </div>
            <el-select v-model="crossProxyZoneIds" multiple filterable collapse-tags collapse-tags-tooltip placeholder="选择域名" style="width:100%">
              <el-option v-for="z in zones" :key="z.id" :label="z.name" :value="z.id" />
            </el-select>
          </div>
        </el-form-item>
        <el-form-item label="记录类型"><el-select v-model="crossProxyType" placeholder="不限（只处理A/AAAA/CNAME）" clearable style="width:100%"><el-option v-for="t in ['A','AAAA','CNAME']" :key="t" :label="t" :value="t" /></el-select></el-form-item>
        <el-form-item label="记录名称"><el-input v-model="crossProxyName" placeholder="不填则匹配所有" /></el-form-item>
        <el-form-item label="橙云状态">
          <el-radio-group v-model="crossProxyEnabled">
            <el-radio :value="true">开启橙云</el-radio>
            <el-radio :value="false">关闭橙云</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <ProgressBlock v-if="proxyProgress.show" :progress="proxyProgress" />
      <template #footer>
        <el-button @click="crossProxyDialog = false" :disabled="proxyProgress.running">取消</el-button>
        <el-button type="primary" @click="submitCrossProxy" :loading="proxyProgress.running">确认执行</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, defineComponent, h } from 'vue';
import { useRoute } from 'vue-router';
import { ElMessage, ElMessageBox, ElProgress } from 'element-plus';
import {
  getCfAccountList, getCfRecords,
  getLocalZoneAll,
  createCfRecord, batchCreateCfRecord, updateCfRecord,
  deleteCfRecord, toggleCfProxy,
  crossZoneDeleteRecords, crossZoneToggleProxy
} from '@/api/modules/cloudflare';

// 进度条组件
const ProgressBlock = defineComponent({
  props: { progress: Object },
  setup(props) {
    return () => h('div', { class: 'progress-block' }, [
      h(ElProgress, { percentage: props.progress.pct, status: props.progress.pct === 100 ? 'success' : '' }),
      h('p', { class: 'progress-msg' }, props.progress.msg)
    ]);
  }
});

const accounts = ref([]);
const zones = ref([]);
const records = ref([]);
const selectedAccountId = ref(null);
const selectedZoneId = ref(null);
const selectedIds = ref([]);
const loading = ref(false);

const recordDialog = ref(false);
const batchDialog = ref(false);
const crossBatchDialog = ref(false);
const crossDeleteDialog = ref(false);
const crossProxyDialog = ref(false);

const editRecord = ref(null);
const recordFormRef = ref();
const batchText = ref('');

// 跨域名批量新增
const crossBatchZoneIds = ref([]);
const crossBatchText = ref('');
const progress = ref({ show: false, running: false, pct: 0, msg: '' });

// 跨域名批量删除
const crossDeleteZoneIds = ref([]);
const crossDeleteType = ref('');
const crossDeleteName = ref('');
const delProgress = ref({ show: false, running: false, pct: 0, msg: '' });

// 跨域名橙云切换
const crossProxyZoneIds = ref([]);
const crossProxyType = ref('');
const crossProxyName = ref('');
const crossProxyEnabled = ref(true);
const proxyProgress = ref({ show: false, running: false, pct: 0, msg: '' });

const dnsTypes = ['A', 'AAAA', 'CNAME', 'TXT', 'MX', 'NS', 'SRV', 'CAA'];
const recordForm = ref({ type: 'A', name: '', content: '', ttl: 1, proxied: false });
const recordRules = {
  type: [{ required: true, message: '请选择类型' }],
  name: [{ required: true, message: '请输入名称' }],
  content: [{ required: true, message: '请输入内容' }]
};

const route = useRoute();

onMounted(async () => {
  const res = await getCfAccountList({ pageSize: 100 });
  accounts.value = res?.data?.list || [];

  // 从 Zone 列表跳转过来时，自动选中账号和域名
  const { accountId, zoneId } = route.query;
  if (accountId && zoneId) {
    selectedAccountId.value = Number(accountId);
    await onAccountChange();
    selectedZoneId.value = zoneId;
    await loadRecords();
  }
});

async function onAccountChange() {
  selectedZoneId.value = null;
  zones.value = [];
  records.value = [];
  crossBatchZoneIds.value = [];
  crossDeleteZoneIds.value = [];
  crossProxyZoneIds.value = [];
  if (!selectedAccountId.value) return;
  try {
    const res = await getLocalZoneAll(selectedAccountId.value);
    const result = res?.data;
    if (!Array.isArray(result)) {
      ElMessage.error('获取域名失败，请先在「Zone 列表」页同步域名');
      return;
    }
    // 本地 zone 用 zoneId 作为选择值
    zones.value = result.map(z => ({ ...z, id: z.zoneId }));
    if (zones.value.length === 0) {
      ElMessage.warning('本地没有域名数据，请先在「Zone 列表」页点击「从 CF 同步」');
    } else {
      ElMessage.success(`已加载 ${zones.value.length} 个域名`);
    }
  } catch (e) {
    ElMessage.error('获取域名异常：' + e.message);
  }
}

async function loadRecords() {
  if (!selectedZoneId.value) return;
  loading.value = true;
  try {
    const res = await getCfRecords({ account_id: selectedAccountId.value, zone_id: selectedZoneId.value });
    const result = res?.data?.result;
    records.value = Array.isArray(result) ? result : [];
    if (!Array.isArray(result)) ElMessage.error('获取记录失败：' + JSON.stringify(res?.data));
  } catch (e) {
    ElMessage.error('获取记录异常：' + e.message);
    records.value = [];
  } finally {
    loading.value = false;
  }
}

function onSelectionChange(rows) { selectedIds.value = rows.map(r => r.id); }

function showAddDialog() {
  editRecord.value = null;
  recordForm.value = { type: 'A', name: '', content: '', ttl: 1, proxied: false };
  recordDialog.value = true;
}
function showSingleBatchDialog() { batchText.value = ''; batchDialog.value = true; }
function showCrossBatchDialog() {
  crossBatchText.value = '';
  crossBatchZoneIds.value = [];
  progress.value = { show: false, running: false, pct: 0, msg: '' };
  crossBatchDialog.value = true;
}
function showCrossDeleteDialog() {
  crossDeleteZoneIds.value = [];
  crossDeleteType.value = '';
  crossDeleteName.value = '';
  delProgress.value = { show: false, running: false, pct: 0, msg: '' };
  crossDeleteDialog.value = true;
}
function showCrossProxyDialog() {
  crossProxyZoneIds.value = [];
  crossProxyType.value = '';
  crossProxyName.value = '';
  crossProxyEnabled.value = true;
  proxyProgress.value = { show: false, running: false, pct: 0, msg: '' };
  crossProxyDialog.value = true;
}

function handleEdit(row) {
  editRecord.value = row;
  recordForm.value = { type: row.type, name: row.name, content: row.content, ttl: row.ttl, proxied: row.proxied };
  recordDialog.value = true;
}

async function handleDelete(row) {
  await ElMessageBox.confirm(`确定删除记录 ${row.name}？`, '提示', { type: 'warning' });
  await deleteCfRecord({ account_id: selectedAccountId.value, zone_id: selectedZoneId.value, record_id: row.id });
  ElMessage.success('删除成功');
  loadRecords();
}

async function submitRecord() {
  await recordFormRef.value.validate();
  if (editRecord.value) {
    await updateCfRecord({ accountId: selectedAccountId.value, zoneId: selectedZoneId.value, recordId: editRecord.value.id, ...recordForm.value });
    ElMessage.success('更新成功');
  } else {
    await createCfRecord({ accountId: selectedAccountId.value, zoneId: selectedZoneId.value, ...recordForm.value });
    ElMessage.success('创建成功');
  }
  recordDialog.value = false;
  loadRecords();
}

function parseRecordsText(text) {
  return text.trim().split('\n').filter(l => l.trim()).map(line => {
    const parts = line.trim().split(/\s+/);
    return { type: parts[0] || 'A', name: parts[1] || '', content: parts[2] || '', ttl: 1, proxied: parts[3] === 'true' };
  }).filter(r => r.name && r.content);
}

async function submitBatch() {
  const list = parseRecordsText(batchText.value);
  if (!list.length) { ElMessage.warning('没有有效记录'); return; }
  const res = await batchCreateCfRecord({ accountId: selectedAccountId.value, zoneId: selectedZoneId.value, records: list });
  ElMessage.success(`完成：成功 ${res?.data?.success}，失败 ${res?.data?.fail}`);
  batchDialog.value = false;
  batchText.value = '';
  loadRecords();
}

async function submitCrossBatch() {
  if (!crossBatchZoneIds.value.length) { ElMessage.warning('请选择至少一个域名'); return; }
  const list = parseRecordsText(crossBatchText.value);
  if (!list.length) { ElMessage.warning('没有有效记录'); return; }
  const total = crossBatchZoneIds.value.length;
  progress.value = { show: true, running: true, pct: 0, msg: '开始处理...' };
  let ok = 0, fail = 0;
  for (let i = 0; i < crossBatchZoneIds.value.length; i++) {
    const zoneId = crossBatchZoneIds.value[i];
    const zoneName = zones.value.find(z => z.id === zoneId)?.name || zoneId;
    progress.value.msg = `处理 ${zoneName} (${i + 1}/${total})`;
    try {
      const res = await batchCreateCfRecord({ accountId: selectedAccountId.value, zoneId, records: list });
      ok += res?.data?.success || 0;
      fail += res?.data?.fail || 0;
    } catch { fail += list.length; }
    progress.value.pct = Math.round(((i + 1) / total) * 100);
  }
  progress.value = { show: true, running: false, pct: 100, msg: `完成：${total} 个域名，成功 ${ok} 条，失败 ${fail} 条` };
  ElMessage.success(progress.value.msg);
  if (selectedZoneId.value && crossBatchZoneIds.value.includes(selectedZoneId.value)) loadRecords();
}

async function submitCrossDelete() {
  if (!crossDeleteZoneIds.value.length) { ElMessage.warning('请选择至少一个域名'); return; }
  if (!crossDeleteType.value && !crossDeleteName.value) { ElMessage.warning('类型和名称至少填一个'); return; }
  await ElMessageBox.confirm(`确定删除 ${crossDeleteZoneIds.value.length} 个域名中匹配的 DNS 记录？此操作不可恢复！`, '危险操作', { type: 'error' });
  delProgress.value = { show: true, running: true, pct: 0, msg: '正在删除...' };
  try {
    const res = await crossZoneDeleteRecords({
      accountId: selectedAccountId.value,
      zoneIds: crossDeleteZoneIds.value,
      type: crossDeleteType.value,
      name: crossDeleteName.value
    });
    delProgress.value = { show: true, running: false, pct: 100, msg: `完成：删除 ${res?.data?.deleted} 条，失败 ${res?.data?.fail} 条` };
    ElMessage.success(delProgress.value.msg);
  } catch (e) {
    delProgress.value.running = false;
    ElMessage.error('操作失败：' + e.message);
  }
  if (selectedZoneId.value && crossDeleteZoneIds.value.includes(selectedZoneId.value)) loadRecords();
}

async function submitCrossProxy() {
  if (!crossProxyZoneIds.value.length) { ElMessage.warning('请选择至少一个域名'); return; }
  proxyProgress.value = { show: true, running: true, pct: 0, msg: '正在处理...' };
  try {
    const res = await crossZoneToggleProxy({
      accountId: selectedAccountId.value,
      zoneIds: crossProxyZoneIds.value,
      type: crossProxyType.value,
      name: crossProxyName.value,
      proxied: crossProxyEnabled.value
    });
    proxyProgress.value = { show: true, running: false, pct: 100, msg: `完成：更新 ${res?.data?.updated} 条，失败 ${res?.data?.fail} 条` };
    ElMessage.success(proxyProgress.value.msg);
  } catch (e) {
    proxyProgress.value.running = false;
    ElMessage.error('操作失败：' + e.message);
  }
  if (selectedZoneId.value && crossProxyZoneIds.value.includes(selectedZoneId.value)) loadRecords();
}

async function showBatchProxy(proxied) {
  await ElMessageBox.confirm(`确定${proxied ? '开启' : '关闭'}选中 ${selectedIds.value.length} 条记录的橙云代理？`, '提示', { type: 'warning' });
  const res = await toggleCfProxy({ accountId: selectedAccountId.value, zoneId: selectedZoneId.value, recordIds: selectedIds.value, proxied });
  ElMessage.success(`完成：成功 ${res?.data?.success}，失败 ${res?.data?.fail}`);
  loadRecords();
}
</script>

<style scoped lang="scss">
.page-container { padding: 20px; display: flex; flex-direction: column; gap: 16px; }
.batch-tip { margin-top: 12px; color: #666; font-size: 13px; }
.tip-text { color: #666; margin-bottom: 6px; font-size: 13px; }
.tip-example { background: #f5f5f5; padding: 8px 12px; border-radius: 4px; font-size: 12px; margin-bottom: 10px; }
.zone-actions { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.zone-count { color: #666; font-size: 13px; }
.progress-block { margin-top: 16px; }
.progress-msg { margin: 8px 0 0; color: #666; font-size: 13px; }
</style>
