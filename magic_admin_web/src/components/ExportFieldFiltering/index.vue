<template>
  <!--角色菜单权限树-->
  <el-drawer v-model="dialogVisible" title="导出字段选择" size="30%" :close-on-click-modal="false">
    <div class="drawerView">
      <div class="drawerContent">
        <!--------------------->
        <el-form size="small">
          <el-form-item :label="item.comment" v-for="(item,index) in list" :key="index">
            <el-switch v-model="item.isSelect"/>
          </el-form-item>
        </el-form>
        <!--------------------------->
      </div>
      <div class="drawerFooter">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirm ">确定</el-button>
      </div>
    </div>
  </el-drawer>

</template>

<script setup>
import {ref} from "vue";
// dialog状态
const emit = defineEmits(["filesselectedfiles"])
const dialogVisible = ref(false);
const list = ref([])


const show = (data) => {
  list.value = data
  list.value.forEach((v) => {
    v.isSelect = true
  })
  console.log(list.value)
  dialogVisible.value = true

}
const confirm = () => {
  dialogVisible.value = false
  let str = ""
  let list_ = list.value.filter(item => item.isSelect)
  str = list_.map(item => item.field).join(",")
  emit("filesselectedfiles", str)
}
defineExpose({
  show
});
</script>
<style lang="scss" scoped>
.drawerView {
  display: flex;
  flex-direction: column;
  height: 100%;

  .drawerContent {
    flex: 1;
    overflow-y: auto;


    .list {
      margin-bottom: 5px;
      border-bottom: 1px dashed rgba(0, 0, 0, 0.5);

      .children {
        margin-bottom: 10px;

        .item {
          .children {
            .item {
            }
          }
        }
      }
    }


  }

  .drawerFooter {
    border-top: 1px solid #f1f1f1;
    padding-top: 10px;
    height: 40px;
    flex-shrink: 0;
    display: flex;
    justify-content: right;

  }
}


</style>
