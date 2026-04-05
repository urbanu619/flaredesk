<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="title">DNS 记录模板</span>
          <el-button type="primary" @click="handleAdd">新增模板</el-button>
        </div>
      </template>

      <el-table v-loading="loading" :data="templates" border stripe>
        <el-table-column prop="name" label="模板名称" min-width="160" />
        <el-table-column label="记录数" width="90">
          <template #default="{ row }">{{ parseRecords(row.records).length }} 条</template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="200" show-overflow-tooltip />
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">{{ formatTs(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="success" @click="handleApply(row)">应用</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑模板 -->
    <el-dialog v-model="editDialog" :title="editingId ? '编辑模板' : '新增模板'" width="680px">
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="90px">
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="formData.name" placeholder="如：基础A记录模板" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formData.remark" placeholder="可选" />
        </el-form-item>
        <el-form-item label="DNS 记录" prop="recordsText">
          <p class="tip-text">每行一条，格式：<code>类型 名称 内容 [proxied]</code></p>
          <pre class="tip-example">A  @  1.2.3.4  true
A  www  1.2.3.4  true
CNAME  mail  mail.example.com</pre>
          <el-input v-model="formData.recordsText" type="textarea" :rows="10" placeholder="A  @  1.2.3.4  true" style="width:100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialog = false">取消</el-button>
        <el-button type="primary" @click="submitForm">保存</el-button>
      </template>
    </el-dialog>

    <!-- 应用模板到域名 -->
    <el-dialog v-model="applyDialog" :title="`应用模板：${applyTemplate?.name}`" width="700px">
      <el-form label-width="90px">
        <el-form-item label="选择账号">
          <el-select v-model="applyAccountId" placeholder="选择 CF 账号" @change="onApplyAccountChange" style="width:100%">
            <el-option v-for="a in accounts" :key="a.id" :label="a.name" :value="a.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="选择域名" v-if="applyZones.length">
          <div style="width:100%">
            <div class="zone-actions">
              <el-button size="small" @click="applyZoneIds = applyZones.map(z=>z.id)">全选</el-button>
              <el-button size="small" @click="applyZoneIds = []">清空</el-button>
              <span class="zone-count">已选 {{ applyZoneIds.length }} / {{ applyZones.length }}</span>
            </div>
            <el-select v-model="applyZoneIds" multiple filterable collapse-tags collapse-tags-tooltip placeholder="选择目标域名" style="width:100%">
              <el-option v-for="z in applyZones" :key="z.id" :label="z.name" :value="z.id" />
            </el-select>
          </div>
        </el-form-item>
        <el-form-item label="记录预览" v-if="applyTemplate">
          <div class="records-preview">
            <div v-for="(r, i) in parseRecords(applyTemplate.records)" :key="i" class="record-row">
              <el-tag size="small" type="info">{{ r.type }}</el-tag>
              <span class="record-name">{{ r.name }}</span>
              <span class="record-content">{{ r.content }}</span>
              <el-tag size="small" :type="r.proxied ? 'warning' : 'info'">{{ r.proxied ? '橙云' : '灰云' }}</el-tag>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <ProgressBlock v-if="applyProgress.show" :progress="applyProgress" />
      <template #footer>
        <el-button @click="applyDialog = false" :disabled="applyProgress.running">取消</el-button>
        <el-button type="primary" @click="submitApply" :loading="applyProgress.running">开始应用</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, defineComponent, h } from 'vue';
import { ElMessage, ElMessageBox, ElProgress } from 'element-plus';
import {
  getDnsTemplateList, createDnsTemplate, updateDnsTemplate, deleteDnsTemplate,
  getCfAccountList, getLocalZoneAll, batchCreateCfRecord
} from '@/api/modules/cloudflare';

const ProgressBlock = defineComponent({
  props: { progress: Object },
  setup(props) {
    return () => h('div', { class: 'progress-block' }, [
      h(ElProgress, { percentage: props.progress.pct, status: props.progress.pct === 100 ? 'success' : '' }),
      h('p', { class: 'progress-msg' }, props.progress.msg)
    ]);
  }
});

const loading = ref(false);
const templates = ref([]);
const editDialog = ref(false);
const editingId = ref(null);
const formRef = ref();
const formData = ref({ name: '', remark: '', recordsText: '' });
const formRules = {
  name: [{ required: true, message: '请输入模板名称' }],
  recordsText: [{ required: true, message: '请输入 DNS 记录' }]
};

// 应用模板
const applyDialog = ref(false);
const applyTemplate = ref(null);
const accounts = ref([]);
const applyAccountId = ref(null);
const applyZones = ref([]);
const applyZoneIds = ref([]);
const applyProgress = ref({ show: false, running: false, pct: 0, msg: '' });

onMounted(async () => {
  loadTemplates();
  const res = await getCfAccountList({ pageSize: 100 });
  accounts.value = res?.data?.list || [];
});

async function loadTemplates() {
  loading.value = true;
  try {
    const res = await getDnsTemplateList({ pageSize: 100 });
    templates.value = res?.data?.list || [];
  } finally {
    loading.value = false;
  }
}

function parseRecords(jsonStr) {
  try { return JSON.parse(jsonStr) || []; } catch { return []; }
}

function textToRecords(text) {
  return text.trim().split('\n').filter(l => l.trim()).map(line => {
    const parts = line.trim().split(/\s+/);
    return { type: parts[0] || 'A', name: parts[1] || '', content: parts[2] || '', ttl: 1, proxied: parts[3] === 'true' };
  }).filter(r => r.name && r.content);
}

function recordsToText(records) {
  return records.map(r => `${r.type}  ${r.name}  ${r.content}${r.proxied ? '  true' : ''}`).join('\n');
}

function formatTs(ts) {
  if (!ts) return '-';
  return new Date(ts * 1000).toLocaleString('zh-CN', { hour12: false }).replace(/\//g, '-');
}

function handleAdd() {
  editingId.value = null;
  formData.value = { name: '', remark: '', recordsText: '' };
  editDialog.value = true;
}

function handleEdit(row) {
  editingId.value = row.id;
  formData.value = { name: row.name, remark: row.remark, recordsText: recordsToText(parseRecords(row.records)) };
  editDialog.value = true;
}

async function handleDelete(row) {
  await ElMessageBox.confirm(`确定删除模板「${row.name}」？`, '提示', { type: 'warning' });
  await deleteDnsTemplate(row.id);
  ElMessage.success('已删除');
  loadTemplates();
}

async function submitForm() {
  await formRef.value.validate();
  const records = textToRecords(formData.value.recordsText);
  if (!records.length) { ElMessage.warning('没有有效的 DNS 记录'); return; }
  const payload = { name: formData.value.name, remark: formData.value.remark, records: JSON.stringify(records) };
  if (editingId.value) {
    await updateDnsTemplate({ id: editingId.value, ...payload });
    ElMessage.success('更新成功');
  } else {
    await createDnsTemplate(payload);
    ElMessage.success('创建成功');
  }
  editDialog.value = false;
  loadTemplates();
}

function handleApply(row) {
  applyTemplate.value = row;
  applyAccountId.value = null;
  applyZones.value = [];
  applyZoneIds.value = [];
  applyProgress.value = { show: false, running: false, pct: 0, msg: '' };
  applyDialog.value = true;
}

async function onApplyAccountChange() {
  applyZones.value = [];
  applyZoneIds.value = [];
  if (!applyAccountId.value) return;
  const res = await getLocalZoneAll(applyAccountId.value);
  applyZones.value = (res?.data || []).map(z => ({ ...z, id: z.zoneId }));
}

async function submitApply() {
  if (!applyAccountId.value) { ElMessage.warning('请选择账号'); return; }
  if (!applyZoneIds.value.length) { ElMessage.warning('请选择至少一个域名'); return; }
  const records = parseRecords(applyTemplate.value.records);
  if (!records.length) { ElMessage.warning('模板中没有有效记录'); return; }
  const total = applyZoneIds.value.length;
  applyProgress.value = { show: true, running: true, pct: 0, msg: '开始处理...' };
  let ok = 0, fail = 0;
  for (let i = 0; i < applyZoneIds.value.length; i++) {
    const zoneId = applyZoneIds.value[i];
    const zoneName = applyZones.value.find(z => z.id === zoneId)?.name || zoneId;
    applyProgress.value.msg = `处理 ${zoneName} (${i + 1}/${total})`;
    try {
      const res = await batchCreateCfRecord({ accountId: applyAccountId.value, zoneId, records });
      ok += res?.data?.success || 0;
      fail += res?.data?.fail || 0;
    } catch { fail += records.length; }
    applyProgress.value.pct = Math.round(((i + 1) / total) * 100);
  }
  applyProgress.value = { show: true, running: false, pct: 100, msg: `完成：${total} 个域名，成功 ${ok} 条，失败 ${fail} 条` };
  ElMessage.success(applyProgress.value.msg);
}
</script>

<style scoped lang="scss">
.page-container { padding: 20px; }
.card-header { display: flex; align-items: center; gap: 12px; }
.title { font-size: 16px; font-weight: 600; }
.tip-text { color: #666; font-size: 13px; margin: 0 0 4px; }
.tip-example { background: #f5f5f5; padding: 8px 12px; border-radius: 4px; font-size: 12px; margin-bottom: 8px; }
.zone-actions { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.zone-count { color: #666; font-size: 13px; }
.records-preview { display: flex; flex-direction: column; gap: 6px; width: 100%; }
.record-row { display: flex; align-items: center; gap: 8px; font-size: 13px; }
.record-name { font-weight: 500; min-width: 80px; }
.record-content { color: #555; flex: 1; }
.progress-block { margin-top: 16px; }
.progress-msg { margin: 8px 0 0; color: #666; font-size: 13px; }
</style>
