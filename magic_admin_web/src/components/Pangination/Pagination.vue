<template>
  <!-- 分页组件 -->
  <el-pagination
    :background="true"
    :current-page="props.pageable.current"
    :page-size="props.pageable.pageSize"
    :page-sizes="[30, 50, 100, 300, 500]"
    :total="props.pageable.total"
    size="small"
    layout="total, sizes, prev, pager, next, jumper"
    @size-change="handleSizeChange"
    @current-change="handleCurrentChange"
  ></el-pagination>
</template>

<script setup name="Pagination">
import { defineProps, defineEmits, nextTick } from 'vue';

const props = defineProps({
  pageable: {
    type: Object,
    required: true
  }
});

const emit = defineEmits(['handleCurrent']);

let isSizeChanging = false;

const handleSizeChange = (newSize) => {
  isSizeChanging = true;
  // 更新 pageSize，并重置 currentPage 为 1
  props.pageable.pageSize = newSize;
  props.pageable.current = 1;
  nextTick(() => {
    emit('handleCurrent', { ...props.pageable });
    isSizeChanging = false;
  });
};

const handleCurrentChange = (newPage) => {
  if (isSizeChanging) return; // 避免在 size-change 后立即触发
  props.pageable.current = newPage;
  emit('handleCurrent', { ...props.pageable });
};

</script>
