<template>
  <div class="page-container">
    <el-alert
      type="info"
      show-icon
      :closable="false"
      title="API Token 权限说明"
      description="若出现 Authentication error (10000)，说明当前账号里配置的 Token 未包含 SSL 签发权限。批量签发需要：Zone → SSL and Certificates → Edit；仅 DNS 编辑不够。请到 Cloudflare 编辑该 Token 勾选上述权限，并在「账号管理」中更新 Token。"
      style="margin-bottom: 12px"
    />
    <!-- 顶部操作区 -->
    <el-card class="selector-card">
      <el-row :gutter="12" align="middle">
        <el-col :span="5">
          <el-select v-model="selectedAccountId" placeholder="选择 CF 账号" @change="onAccountChange" style="width:100%">
            <el-option v-for="a in accounts" :key="a.id" :label="a.name" :value="a.id" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="validityDays" style="width:100%">
            <el-option label="有效期 90 天" :value="90" />
            <el-option label="有效期 1 年" :value="365" />
            <el-option label="有效期 2 年" :value="730" />
            <el-option label="有效期 15 年" :value="5475" />
          </el-select>
        </el-col>
        <el-col :span="15" style="display:flex;gap:8px;flex-wrap:wrap">
          <el-button type="primary" :disabled="!selectedAccountId" @click="loadZones">同步 Zone 列表</el-button>
          <el-button type="success" :disabled="selectedZoneIds.length === 0" :loading="generating" @click="handleBatchGenerate">
            批量生成证书 ({{ selectedZoneIds.length }} 个)
          </el-button>
          <el-button :disabled="!selectedAccountId" @click="selectAll">全选</el-button>
          <el-button @click="clearSelect">清空选择</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- Zone 列表 -->
    <el-card class="table-card">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span>域名列表（勾选后批量生成 Origin 证书）</span>
          <el-input v-model="searchKw" placeholder="搜索域名" clearable style="width:220px" />
        </div>
      </template>

      <el-table
        ref="tableRef"
        v-loading="loadingZones"
        :data="filteredZones"
        @selection-change="onSelectionChange"
        border
        stripe
        row-key="id"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column prop="name" label="域名" min-width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="plan_name" label="套餐" width="120" show-overflow-tooltip />
        <el-table-column label="证书结果" min-width="160">
          <template #default="{ row }">
            <template v-if="certResultMap[row.zone_id || row.id]">
              <el-tag type="success" size="small" v-if="certResultMap[row.zone_id || row.id].success">已生成</el-tag>
              <el-tooltip v-else :content="certResultMap[row.zone_id || row.id].error" placement="top">
                <el-tag type="danger" size="small">失败</el-tag>
              </el-tooltip>
            </template>
            <span v-else class="text-muted">—</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="table-footer" v-if="selectedZoneIds.length > 0">
        已选 {{ selectedZoneIds.length }} 个域名
      </div>
    </el-card>

    <!-- 生成结果摘要 -->
    <el-card v-if="lastResult" class="result-card">
      <template #header>
        <span>最近一次生成结果</span>
      </template>
      <el-descriptions :column="3" border size="small">
        <el-descriptions-item label="总计">{{ lastResult.total }}</el-descriptions-item>
        <el-descriptions-item label="成功">
          <el-tag type="success">{{ lastResult.success }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="失败">
          <el-tag type="danger">{{ lastResult.fail }}</el-tag>
        </el-descriptions-item>
      </el-descriptions>
      <div v-if="lastResult.fail > 0" style="margin-top:12px">
        <p style="font-weight:600;margin-bottom:6px">失败详情：</p>
        <el-tag
          v-for="r in lastResult.results.filter(r => !r.success)"
          :key="r.domain"
          type="danger"
          style="margin:4px"
        >{{ r.domain }}: {{ r.error }}</el-tag>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { getCfAccountList } from "@/api/modules/cloudflare";
import { getLocalZoneAll } from "@/api/modules/cloudflare";
import { batchGenerateCfCerts } from "@/api/modules/cloudflare";

const accounts = ref([]);
const selectedAccountId = ref(null);
const validityDays = ref(5475);
const zones = ref([]);
const loadingZones = ref(false);
const generating = ref(false);
const searchKw = ref("");
const selectedRows = ref([]);
const certResultMap = ref({});
const lastResult = ref(null);
const tableRef = ref(null);

const selectedZoneIds = computed(() => selectedRows.value.map(r => r.zone_id || r.id));

const filteredZones = computed(() => {
  if (!searchKw.value) return zones.value;
  return zones.value.filter(z => z.name.includes(searchKw.value));
});

onMounted(async () => {
  const res = await getCfAccountList({ page: 1, pageSize: 100 });
  accounts.value = res?.data?.list || res?.data || [];
});

async function onAccountChange() {
  zones.value = [];
  certResultMap.value = {};
  lastResult.value = null;
  selectedRows.value = [];
  await loadZones();
}

async function loadZones() {
  if (!selectedAccountId.value) return;
  loadingZones.value = true;
  try {
    const res = await getLocalZoneAll(selectedAccountId.value);
    zones.value = res?.data || [];
  } catch (e) {
    ElMessage.error("加载 Zone 列表失败");
  } finally {
    loadingZones.value = false;
  }
}

function onSelectionChange(rows) {
  selectedRows.value = rows;
}

function selectAll() {
  tableRef.value?.toggleAllSelection();
}

function clearSelect() {
  tableRef.value?.clearSelection();
}

async function handleBatchGenerate() {
  if (!selectedAccountId.value || selectedZoneIds.value.length === 0) return;

  generating.value = true;
  certResultMap.value = {};
  lastResult.value = null;

  const zonesPayload = selectedRows.value.map(r => ({
    zoneId: r.zone_id || r.id,
    domain: r.name
  }));

  try {
    // http.download 返回的 data 直接是 Blob（responseType: blob）
    const blob = await batchGenerateCfCerts({
      accountId: selectedAccountId.value,
      zones: zonesPayload,
      validityDays: validityDays.value,
      requestType: "origin-rsa"
    });

    // 判断是 zip 还是 JSON（全部失败时后端返回 JSON，但 axios 仍以 Blob 形式返回）
    if (blob.type === "application/zip") {
      triggerDownload(blob);
      ElMessage.success(`证书已生成，zip 开始下载`);
    } else {
      // 读取 Blob 内容，解析为 JSON
      const text = await blob.text();
      let json;
      try { json = JSON.parse(text); } catch { json = {}; }
      const inner = json.data || json;
      lastResult.value = {
        total: inner.total || zonesPayload.length,
        success: inner.success || 0,
        fail: inner.fail || zonesPayload.length,
        results: inner.results || []
      };
      for (const r of lastResult.value.results) {
        certResultMap.value[r.domain] = r;
      }
      ElMessage.warning("全部失败，请查看下方详情");
    }
  } catch (e) {
    ElMessage.error("请求失败: " + (e?.message || e));
  } finally {
    generating.value = false;
  }
}

function triggerDownload(blob) {
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `cf_certs_${Date.now()}.zip`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}
</script>

<style scoped lang="scss">
.page-container {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.selector-card,
.table-card,
.result-card {
  border-radius: 8px;
}

.table-footer {
  margin-top: 10px;
  color: #606266;
  font-size: 13px;
}

.text-muted {
  color: #c0c4cc;
}
</style>
