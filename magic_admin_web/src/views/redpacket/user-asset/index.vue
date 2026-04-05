<template>
  <div class="user-asset-container">
    <ProTable
      ref="proTableRef"
      :columns="columns"
      :request-api="requestApi"
      :search-columns="searchColumns"
      :show-pagination="true"
    >
      <!-- Telegram用户名列 -->
      <template #telegramUsername="{ row }">
        <span v-if="row.telegramUsername" style="color: #409eff">@{{ row.telegramUsername }}</span>
        <span v-else style="color: #909399">-</span>
      </template>

      <!-- 余额列 -->
      <template #balance="{ row }">
        <span style="color: #67c23a; font-weight: bold">{{ formatAmount(row.balance) }}</span>
      </template>

      <!-- 冻结列 -->
      <template #frozen="{ row }">
        <span :style="{ color: row.frozen > 0 ? '#e6a23c' : '#909399' }">{{ formatAmount(row.frozen) }}</span>
      </template>

      <!-- 总资产列 -->
      <template #totalAmount="{ row }">
        <span style="font-weight: bold">{{ formatAmount(row.totalAmount) }}</span>
      </template>

      <!-- 操作列 -->
      <template #actions="{ row }">
        <el-button link type="primary" @click="handleView(row)">查看</el-button>
        <el-button link type="primary" @click="handleBillDetail(row)">明细</el-button>
      </template>
    </ProTable>

    <!-- 查看对话框 -->
    <el-dialog v-model="viewDialogVisible" title="资产详情" width="500px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="ID">{{ viewData.id }}</el-descriptions-item>
        <el-descriptions-item label="用户ID">{{ viewData.userId }}</el-descriptions-item>
        <el-descriptions-item label="Telegram">
          <span v-if="viewData.telegramUsername" style="color: #409eff">@{{ viewData.telegramUsername }}</span>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="币种">{{ viewData.symbol }}</el-descriptions-item>
        <el-descriptions-item label="可用余额">
          <span style="color: #67c23a; font-weight: bold">{{ formatAmount(viewData.balance) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="冻结资产">
          <span style="color: #e6a23c">{{ formatAmount(viewData.frozen) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="总资产">
          <span style="font-weight: bold">{{ formatAmount(viewData.totalAmount) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="版本号">{{ viewData.version }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatTimestamp(viewData.createdAt) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatTimestamp(viewData.updatedAt) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!-- 资产明细对话框 -->
    <el-dialog v-model="billDialogVisible" title="资产明细" width="900px">
      <div style="margin-bottom: 12px; color: #606266">
        用户ID: {{ billUserId }} | 币种: {{ billSymbol }}
      </div>
      <el-table :data="billList" border stripe v-loading="billLoading" max-height="500">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="businessName" label="业务类型" width="120">
          <template #default="{ row }">
            <el-tag size="small" :type="getBillTagType(row.businessName)">{{ row.businessName || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="变动金额" width="140">
          <template #default="{ row }">
            <span :style="{ color: parseFloat(row.amount) >= 0 ? '#67c23a' : '#f56c6c', fontWeight: 'bold' }">
              {{ parseFloat(row.amount) >= 0 ? '+' : '' }}{{ formatAmount(row.amount) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="beforeAmount" label="变动前" width="130">
          <template #default="{ row }">
            {{ formatAmount(row.beforeAmount) }}
          </template>
        </el-table-column>
        <el-table-column prop="afterAmount" label="变动后" width="130">
          <template #default="{ row }">
            {{ formatAmount(row.afterAmount) }}
          </template>
        </el-table-column>
        <el-table-column prop="describe" label="备注" min-width="200" show-overflow-tooltip />
        <el-table-column prop="createdAt" label="时间" width="170">
          <template #default="{ row }">
            {{ formatTimestamp(row.createdAt) }}
          </template>
        </el-table-column>
      </el-table>
      <div v-if="billList.length === 0 && !billLoading" style="text-align: center; padding: 40px; color: #909399">
        暂无流水记录
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import ProTable from '@/components/ProTable/index.vue';
import { getUserAssetList, getUserAssetBillList } from '@/api/modules/redpacket';

const proTableRef = ref();
const viewDialogVisible = ref(false);
const viewData = ref({});

// 资产明细
const billDialogVisible = ref(false);
const billLoading = ref(false);
const billList = ref([]);
const billUserId = ref('');
const billSymbol = ref('');

// 后端返回格式适配（MagicAssetService 返回 { list, paging, cols, url }）
const requestApi = async params => {
  const res = await getUserAssetList(params);
  if (res.data) {
    return { data: { list: res.data.list, paging: res.data.paging } };
  }
  return res;
};

// 表格列配置
const columns = [
  { prop: 'id', label: 'ID', width: 80 },
  { prop: 'userId', label: '用户ID', width: 100 },
  { prop: 'telegramUsername', label: 'Telegram', width: 150, slot: 'telegramUsername' },
  { prop: 'symbol', label: '币种', width: 100 },
  { prop: 'balance', label: '可用余额', minWidth: 150, slot: 'balance' },
  { prop: 'frozen', label: '冻结资产', width: 150, slot: 'frozen' },
  { prop: 'totalAmount', label: '总资产', width: 150, slot: 'totalAmount' }
];

// 搜索列配置
const searchColumns = [
  { prop: 'userId', label: '用户ID', type: 'input' },
  { prop: 'symbol', label: '币种', type: 'input' }
];

// 格式化金额
const formatAmount = val => {
  if (val === null || val === undefined) return '0.00';
  const num = typeof val === 'string' ? parseFloat(val) : val;
  if (isNaN(num)) return '0.00';
  return num.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
};

// 格式化时间戳
const formatTimestamp = ts => {
  if (!ts) return '-';
  return new Date(ts * 1000).toLocaleString('zh-CN');
};

// 业务类型标签颜色
const getBillTagType = name => {
  if (!name) return '';
  if (name.includes('收入') || name.includes('充值')) return 'success';
  if (name.includes('扣款') || name.includes('扣除')) return 'danger';
  if (name.includes('冻结')) return 'warning';
  if (name.includes('退款') || name.includes('解冻')) return '';
  return 'info';
};

// 查看
function handleView(row) {
  viewData.value = { ...row };
  viewDialogVisible.value = true;
}

// 资产明细
async function handleBillDetail(row) {
  billUserId.value = row.userId;
  billSymbol.value = row.symbol;
  billList.value = [];
  billDialogVisible.value = true;
  billLoading.value = true;

  try {
    const res = await getUserAssetBillList({
      userId: row.userId,
      symbol: row.symbol,
      pageSize: 100,
      pageNum: 1
    });
    if (res.data) {
      billList.value = res.data.list || res.data || [];
    }
  } catch (e) {
    console.error('获取资产明细失败', e);
  } finally {
    billLoading.value = false;
  }
}
</script>

<style scoped lang="scss">
.user-asset-container {
  padding: 20px;
}
</style>
