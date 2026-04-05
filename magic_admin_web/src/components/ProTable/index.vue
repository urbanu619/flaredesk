<template>
  <div class="pro-table-container">
    <!-- 搜索表单 -->
    <div v-if="searchColumns && searchColumns.length > 0" class="search-form">
      <el-form :inline="true" :model="searchForm" @submit.prevent="handleSearch">
        <el-form-item
          v-for="column in searchColumns"
          :key="column.prop"
          :label="column.label"
        >
          <el-input
            v-if="!column.type || column.type === 'input'"
            v-model="searchForm[column.prop]"
            :placeholder="`请输入${column.label}`"
            clearable
          />
          <el-select
            v-else-if="column.type === 'select'"
            v-model="searchForm[column.prop]"
            :placeholder="`请选择${column.label}`"
            clearable
          >
            <el-option
              v-for="option in column.options"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 工具栏 -->
    <div v-if="toolbarButtons && toolbarButtons.length > 0" class="toolbar">
      <el-button
        v-for="button in toolbarButtons"
        :key="button.label"
        :type="button.type || 'primary'"
        :icon="button.icon"
        @click="button.onClick"
      >
        {{ button.label }}
      </el-button>
    </div>

    <!-- 表格 -->
    <el-table
      :data="tableData"
      :loading="loading"
      border
      stripe
      style="width: 100%"
      @selection-change="handleSelectionChange"
    >
      <el-table-column v-if="showSelection" type="selection" width="55" />
      <el-table-column
        v-for="column in columns"
        :key="column.prop"
        :prop="column.prop"
        :label="column.label"
        :width="column.width"
        :min-width="column.minWidth"
        :align="column.align || 'center'"
      >
        <template #default="scope">
          <slot
            v-if="column.slot"
            :name="column.slot"
            :row="scope.row"
            :column="column"
          />
          <span v-else>{{ scope.row[column.prop] }}</span>
        </template>
      </el-table-column>
      <el-table-column
        v-if="showActions"
        label="操作"
        :width="actionsWidth || 200"
        align="center"
        fixed="right"
      >
        <template #default="scope">
          <slot name="actions" :row="scope.row" />
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div v-if="showPagination" class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="totalCount"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted } from 'vue';

const props = defineProps({
  // 表格列配置
  columns: {
    type: Array,
    required: true,
  },
  // 搜索列配置
  searchColumns: {
    type: Array,
    default: () => [],
  },
  // 工具栏按钮
  toolbarButtons: {
    type: Array,
    default: () => [],
  },
  // 表格数据
  data: {
    type: Array,
    default: () => [],
  },
  // 是否显示选择框
  showSelection: {
    type: Boolean,
    default: false,
  },
  // 是否显示操作列
  showActions: {
    type: Boolean,
    default: true,
  },
  // 操作列宽度
  actionsWidth: {
    type: Number,
    default: 200,
  },
  // 是否显示分页
  showPagination: {
    type: Boolean,
    default: true,
  },
  // 加载状态
  loading: {
    type: Boolean,
    default: false,
  },
  // 请求函数
  requestApi: {
    type: Function,
    default: null,
  },
});

const emit = defineEmits(['search', 'selection-change', 'page-change']);

const tableData = ref(props.data);
const searchForm = reactive({});
const currentPage = ref(1);
const pageSize = ref(10);
const totalCount = ref(0);
const selectedRows = ref([]);

// 监听外部数据变化
watch(
  () => props.data,
  (newData) => {
    tableData.value = newData;
  }
);

// 搜索
const handleSearch = () => {
  currentPage.value = 1;
  fetchData();
};

// 重置
const handleReset = () => {
  Object.keys(searchForm).forEach((key) => {
    delete searchForm[key];
  });
  currentPage.value = 1;
  fetchData();
};

// 选择变化
const handleSelectionChange = (selection) => {
  selectedRows.value = selection;
  emit('selection-change', selection);
};

// 分页大小变化
const handleSizeChange = (val) => {
  pageSize.value = val;
  fetchData();
};

// 当前页变化
const handleCurrentChange = (val) => {
  currentPage.value = val;
  fetchData();
};

// 获取数据
const fetchData = async () => {
  if (props.requestApi) {
    const params = {
      ...searchForm,
      current: currentPage.value,
      pageSize: pageSize.value,
    };
    emit('search', params);
    try {
      const result = await props.requestApi(params);
      if (result && result.data) {
        tableData.value = result.data.list || result.data;
        // 从 paging 中提取总数
        if (result.data.paging && result.data.paging.total !== undefined) {
          totalCount.value = result.data.paging.total;
        }
      }
    } catch (error) {
      console.error('Failed to fetch data:', error);
    }
  } else {
    emit('search', {
      ...searchForm,
      current: currentPage.value,
      pageSize: pageSize.value,
    });
  }
};

// 刷新数据
const refresh = () => {
  fetchData();
};

// 暴露方法给父组件
defineExpose({
  refresh,
  searchForm,
  selectedRows,
});

onMounted(() => {
  if (props.requestApi) {
    fetchData();
  }
});
</script>

<style scoped lang="scss">
.pro-table-container {
  .search-form {
    margin-bottom: 16px;
    padding: 16px;
    background: #fff;
    border-radius: 4px;
  }

  .toolbar {
    margin-bottom: 16px;
    padding: 16px;
    background: #fff;
    border-radius: 4px;
  }

  .pagination {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
