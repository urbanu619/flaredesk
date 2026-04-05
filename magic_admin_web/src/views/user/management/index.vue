
<template>
	<div class="table-box">
		<div class="card table-search">
			<el-form ref="elSearchFormRef" :inline="true" size="small" :model="searchForm" class="demo-form-inline" @keyup.enter="onSubmit">
				<el-form-item label="UID">
					<el-input v-model="searchForm.uid" placeholder="UID" style="width:139px"></el-input>
				</el-form-item>
				<el-form-item label="邀请码">
					<el-input v-model="searchForm.code" placeholder="邀请码" style="width:139px"></el-input>
				</el-form-item>
				<el-form-item label="状态">
					<el-select v-model="searchForm.enable" placeholder="状态" style="width:139px">
						<el-option :label="'有效'" :value="1" />
						<el-option :label="'无效'" :value="0" />
					</el-select>
				</el-form-item>

				<el-form-item label="时间">
					<el-date-picker v-model="timeValue" type="daterange" range-separator="-" start-placeholder="开始时间" end-placeholder="结束时间" />
				</el-form-item>

				<div class="operation">
					<el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
					<el-button icon="refresh" @click="onResetSearch">重置</el-button>
				</div>
			</el-form>
		</div>

		<div class="card table-main">
			<div class="table-header">
				<div class="header-button-lf">
					<el-button type="success" @click="exportTable">导出</el-button>
				</div>
				<div class="header-button-ri">
					<el-button :icon="Refresh" circle @click="refresh" />
				</div>
			</div>

			<el-table ref="myTable" :data="tableData" border size="small">
				<el-table-column prop="id" label="ID" align="center" :width="80" />
				<el-table-column prop="uid" label="UID" align="center" :min-width="150" />
				<el-table-column prop="code" label="邀请码" align="center" :min-width="150" />
				<el-table-column prop="refCode" label="上级邀请码" align="center" :min-width="150" />
				<el-table-column prop="parentUid" label="上级交易所ID" align="center" :min-width="150" />

				<el-table-column prop="enable" label="账户状态" align="center" :min-width="120">
					<template #default="{ row }">
						<span>{{ valueToLabel(enableOptions, row.enable) }}</span>
					</template>
				</el-table-column>

				<el-table-column prop="lockWithdraw" label="锁定提现" align="center" :min-width="120">
					<template #default="{ row }">
						<span>{{ valueToLabel(lockOptions, row.lockWithdraw) }}</span>
					</template>
				</el-table-column>

				<el-table-column prop="isWhite" label="白名单" align="center" :min-width="120">
					<template #default="{ row }">
						<span>{{ row.isWhite === 1 ? '是' : '否' }}</span>
					</template>
				</el-table-column>

				<el-table-column prop="isReinvestment" label="复投" align="center" :min-width="120">
					<template #default="{ row }">
						<span>{{ row.isReinvestment ? '是' : '否' }}</span>
					</template>
				</el-table-column>

				<el-table-column prop="lastLoginTime" label="最近登录时间" align="center" :min-width="160">
					<template #default="{ row }">
						<span v-if="row.lastLoginTime">{{ formatUnix(row.lastLoginTime) }}</span>
					</template>
				</el-table-column>

				<el-table-column prop="lastLoginIp" label="登录IP" align="center" :min-width="150" />
				<el-table-column prop="lastLoginRemoteIp" label="Remote IP" align="center" :min-width="150" />

				<el-table-column prop="remark" label="备注" align="center" :min-width="220" />

				<el-table-column prop="updatedAt" label="更新时间" align="center" :min-width="160">
					<template #default="{ row }">
						<span v-if="row.updatedAt">{{ formatUnix(row.updatedAt) }}</span>
					</template>
				</el-table-column>

				<el-table-column fixed="right" label="操作" align="center" :min-width="150">
					<template #default="scope">
						<el-button type="primary" link @click="openDialog(scope.row)">编辑</el-button>
						<el-button type="warning" link @click="openSetParentDialog(scope.row)">设置上级</el-button>
						<el-button type="danger" link @click="handleResetParent(scope.row)">重置上级</el-button>
					</template>
				</el-table-column>
			</el-table>

			<Pagination :pageable="pageable" @handleCurrent="handleCurrent" />
		</div>
	</div>

	<!-- 导出字段过滤 -->
	<ExportFieldFiltering 
		ref="ExportFieldFilteringView" 
		:derived_field="derived_field" 
		@filesselectedfiles="filesselectedfiles"
	>
	</ExportFieldFiltering>

	<!-- 编辑用户状态弹窗 -->
	<el-dialog
		v-model="dialogEditVisible"
		title="编辑用户状态"
		width="520px"
		center
		destroy-on-close
	>
		<el-form
			ref="editFormRef"
			:model="editFormData"
			:rules="editFormRules"
			label-width="150px"
		>
			<el-form-item label="用户UID" prop="uid">
				<el-input v-model="editFormData.uid" disabled />
			</el-form-item>
			<el-form-item label="是否冻结账号" prop="enable">
				<el-radio-group v-model="editFormData.enable" size="small">
					<el-radio-button :value="false">正常</el-radio-button>
					<el-radio-button :value="true">冻结</el-radio-button>
				</el-radio-group>
			</el-form-item>
			<el-form-item label="是否锁定提现" prop="lockWithdraw">
				<el-radio-group v-model="editFormData.lockWithdraw" size="small">
					<el-radio-button :value="false">未锁定</el-radio-button>
					<el-radio-button :value="true">已锁定</el-radio-button>
				</el-radio-group>
			</el-form-item>
			<el-form-item label="是否根地址" prop="isRoot">
				<el-radio-group v-model="editFormData.isRoot" size="small">
					<el-radio-button :value="true">是</el-radio-button>
					<el-radio-button :value="false">否</el-radio-button>
				</el-radio-group>
			</el-form-item>
			<el-form-item label="是否零号线" prop="isZero">
				<el-radio-group v-model="editFormData.isZero" size="small">
					<el-radio-button :value="true">是</el-radio-button>
					<el-radio-button :value="false">否</el-radio-button>
				</el-radio-group>
			</el-form-item>
			<el-form-item label="赎回锁定" prop="lockRedeem">
				<el-radio-group v-model="editFormData.lockRedeem" size="small">
					<el-radio-button :value="false">正常</el-radio-button>
					<el-radio-button :value="true">锁定</el-radio-button>
				</el-radio-group>
			</el-form-item>
			<el-form-item label="锁定质押" prop="lockStake">
				<el-radio-group v-model="editFormData.lockStake" size="small">
					<el-radio-button :value="false">正常</el-radio-button>
					<el-radio-button :value="true">锁定不允许质押</el-radio-button>
				</el-radio-group>
			</el-form-item>
			<el-form-item label="锁定产出收益" prop="lockStakeProfit">
				<el-radio-group v-model="editFormData.lockStakeProfit" size="small">
					<el-radio-button :value="false">正常</el-radio-button>
					<el-radio-button :value="true">锁定不产出收益</el-radio-button>
				</el-radio-group>
			</el-form-item>
			<el-form-item label="锁定伞下提现" prop="subsLockWithdraw">
				<el-radio-group v-model="editFormData.subsLockWithdraw" size="small">
					<el-radio-button :value="false">未锁定</el-radio-button>
					<el-radio-button :value="true">锁定</el-radio-button>
				</el-radio-group>
			</el-form-item>
			<el-form-item label="自动领取静态收益" prop="unCollection">
				<el-radio-group v-model="editFormData.unCollection" size="small">
					<el-radio-button :value="false">正常</el-radio-button>
					<el-radio-button :value="true">自动领取</el-radio-button>
				</el-radio-group>
			</el-form-item>
			<el-form-item label="参与小区排行" prop="withoutFewRegionRanking">
				<el-radio-group v-model="editFormData.withoutFewRegionRanking" size="small">
					<el-radio-button :value="false">参与</el-radio-button>
					<el-radio-button :value="true">不参与</el-radio-button>
				</el-radio-group>
			</el-form-item>
			<el-form-item label="参与质押排行" prop="withoutStakeRanking">
				<el-radio-group v-model="editFormData.withoutStakeRanking" size="small">
					<el-radio-button :value="false">参与</el-radio-button>
					<el-radio-button :value="true">不参与</el-radio-button>
				</el-radio-group>
			</el-form-item>
		</el-form>
		<template #footer>
			<div class="dialog-footer">
				<el-button @click="dialogEditVisible = false">取消</el-button>
				<el-button type="primary" :loading="editLoading" @click="handleEditSubmit">
					确定
				</el-button>
			</div>
		</template>
	</el-dialog>

	<!-- 设置上级弹窗 -->
	<el-dialog
		v-model="dialogSetParentVisible"
		title="设置上级"
		width="420px"
		center
		destroy-on-close
	>
		<el-form
			ref="setParentFormRef"
			:model="setParentFormData"
			:rules="setParentFormRules"
			label-width="100px"
		>
			<el-form-item label="用户ID" prop="userId">
				<el-input v-model="setParentFormData.userId" disabled />
			</el-form-item>
			<el-form-item label="上级邀请码" prop="parentCode">
				<el-input v-model="setParentFormData.parentCode" placeholder="请输入上级邀请码" clearable />
			</el-form-item>
		</el-form>
		<template #footer>
			<div class="dialog-footer">
				<el-button @click="dialogSetParentVisible = false">取消</el-button>
				<el-button type="primary" :loading="setParentLoading" @click="handleSetParentSubmit">
					确定
				</el-button>
			</div>
		</template>
	</el-dialog>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import dayjs from "dayjs";
import { Refresh } from "@element-plus/icons-vue";
import { ElNotification, ElMessageBox } from "element-plus";
import Pagination from "@/components/Pangination/Pagination.vue";
import ExportFieldFiltering from "@/components/ExportFieldFiltering/index.vue";
import { valueToLabel } from "@/utils/index.ts";
import { downLoadFile } from "@/api/downloadFile.js";
import { magicUserList, userSetUser, userSetParent } from "@/api/modules/user.js";

const timeValue = ref([]);
const searchForm = ref({ uid: null, code: null, enable: null });
const tableData = ref([]);
const pageable = reactive({ pageNum: 1, pageSize: 30, total: 0 });

const ExportFieldFilteringView = ref(null);
const derived_field = ref([]);

const enableOptions = [
	{ value: 1, label: "有效" },
	{ value: 0, label: "无效" }
];
const lockOptions = [
	{ value: 0, label: "未锁定" },
	{ value: 1, label: "已锁定" }
];

// 编辑弹窗相关
const dialogEditVisible = ref(false);
const editFormRef = ref(null);
const editLoading = ref(false);
const originalData = ref(null); // 保存原始数据用于比较
const editFormData = ref({
	// 用户UID
	uid: "",
	// 是否冻结账号
	enable: false,
	// 是否根地址
	isRoot: false,
	// 是否零号线
	isZero: false,
	// 赎回锁定
	lockRedeem: false,
	// 0 正常 1锁定不允许质押
	lockStake: false,
	// 0 正常 1锁定不产出收益
	lockStakeProfit: false,
	// 是否锁定提现
	lockWithdraw: false,
	// 锁定伞下提现  0=未锁定  1=锁定
	subsLockWithdraw: false,
	// 0 正常 1 自动领取静态收益
	unCollection: false,
	// 0 参与 1不参与个人排行（小区）
	withoutFewRegionRanking: false,
	// 0 参与 1不参与个人排行（质押）
	withoutStakeRanking: false
});
const editFormRules = ref({
	uid: [{ required: true, message: "用户UID不能为空", trigger: "blur" }],
	enable: [{ required: true, message: "请选择是否冻结账号", trigger: "change" }],
	isRoot: [{ required: true, message: "请选择是否根地址", trigger: "change" }],
	isZero: [{ required: true, message: "请选择是否零号线", trigger: "change" }],
	lockRedeem: [{ required: true, message: "请选择是否赎回锁定", trigger: "change" }],
	lockStake: [{ required: true, message: "请选择是否锁定质押", trigger: "change" }],
	lockStakeProfit: [{ required: true, message: "请选择是否锁定产出收益", trigger: "change" }],
	lockWithdraw: [{ required: true, message: "请选择是否锁定提现", trigger: "change" }],
	subsLockWithdraw: [{ required: true, message: "请选择是否锁定伞下提现", trigger: "change" }],
	unCollection: [{ required: true, message: "请选择是否自动领取静态收益", trigger: "change" }],
	withoutFewRegionRanking: [{ required: true, message: "请选择是否参与小区排行", trigger: "change" }],
	withoutStakeRanking: [{ required: true, message: "请选择是否参与质押排行", trigger: "change" }]
});

// 设置上级弹窗相关
const dialogSetParentVisible = ref(false);
const setParentFormRef = ref(null);
const setParentLoading = ref(false);
const setParentFormData = ref({
	userId: "",
	parentCode: ""
});
const setParentFormRules = ref({
	userId: [{ required: true, message: "用户ID不能为空", trigger: "blur" }],
	parentCode: [{ required: true, message: "请输入上级邀请码", trigger: "blur" }]
});

const formatUnix = (val) => {
	if (!val) return "";
	return dayjs.unix(val).format("YYYY-MM-DD HH:mm:ss");
};

onMounted(() => {
	getList();
});

const getList = async () => {
	tableData.value = [];
	const params = {
		current: pageable.pageNum,
		pageSize: pageable.pageSize,
		order: "id desc",
		beginTime: timeValue.value?.[0] ? dayjs(timeValue.value[0]).startOf("day").unix() : null,
		endTime: timeValue.value?.[1] ? dayjs(timeValue.value[1]).endOf("day").unix() : null,
		...searchForm.value
	};
	const res = await magicUserList(params);
	if (res.code === 200) {
		tableData.value = res.data.list;
		pageable.total = res.data.paging.total;

		// 设置可选择的导出字段
		if (res.data.cols) {
			derived_field.value = res.data.cols;
		}
	}
};

const onSubmit = () => {
	pageable.pageNum = 1;
	getList();
};

const refresh = () => getList();

const onResetSearch = () => {
	searchForm.value.uid = null;
	searchForm.value.code = null;
	searchForm.value.enable = null;
	timeValue.value = [];
	getList();
};

const handleCurrent = (data) => {
	pageable.pageNum = data.current;
	pageable.pageSize = data.pageSize;
	getList();
};

const exportTable = () => {
	ExportFieldFilteringView.value.show(derived_field.value);
};

// 字段选择回调 - 导出
const filesselectedfiles = async (str) => {
	const params = {
		current: pageable.pageNum,
		pageSize: pageable.pageSize,
		order: "id desc",
		isExport: true,
		fields: str,
		beginTime: timeValue.value?.[0] ? dayjs(timeValue.value[0]).startOf("day").unix() : null,
		endTime: timeValue.value?.[1] ? dayjs(timeValue.value[1]).endOf("day").unix() : null,
		...searchForm.value
	};

	const res = await magicUserList(params);

	try {
		const downFile = {
			fileName: "magic_user",
			fileUrl: res.data.url
		};
		await downLoadFile(downFile);
		ElNotification.success("导出成功");
	} catch {
		ElNotification.error("导出失败");
	}
};

const openDialog = (row) => {
	// 将行数据转换为编辑表单格式（使用 true/false）
	const currentEnable = row.enable === 0; // enable=0 表示冻结，所以 true 表示冻结
	const currentLockWithdraw = row.lockWithdraw === 1 || row.lockWithdraw === true;
	const currentIsRoot = row.isRoot === 1 || row.isRoot === true;
	const currentIsZero = row.isZero === 1 || row.isZero === true;
	const currentLockRedeem = row.lockRedeem === 1 || row.lockRedeem === true;
	const currentLockStake = row.lockStake === 1 || row.lockStake === true;
	const currentLockStakeProfit = row.lockStakeProfit === 1 || row.lockStakeProfit === true;
	const currentSubsLockWithdraw = row.subsLockWithdraw === 1 || row.subsLockWithdraw === true;
	const currentUnCollection = row.unCollection === 1 || row.unCollection === true;
	const currentWithoutFewRegionRanking = row.withoutFewRegionRanking === 1 || row.withoutFewRegionRanking === true;
	const currentWithoutStakeRanking = row.withoutStakeRanking === 1 || row.withoutStakeRanking === true;
	
	// 保存原始数据（用于比较）
	originalData.value = {
		uid: String(row.uid || row.userId || ""),
		enable: currentEnable,
		lockWithdraw: currentLockWithdraw,
		isRoot: currentIsRoot,
		isZero: currentIsZero,
		lockRedeem: currentLockRedeem,
		lockStake: currentLockStake,
		lockStakeProfit: currentLockStakeProfit,
		subsLockWithdraw: currentSubsLockWithdraw,
		unCollection: currentUnCollection,
		withoutFewRegionRanking: currentWithoutFewRegionRanking,
		withoutStakeRanking: currentWithoutStakeRanking
	};
	
	// 将行数据映射到编辑表单
	editFormData.value = {
		uid: String(row.uid || row.userId || ""),
		enable: currentEnable,
		lockWithdraw: currentLockWithdraw,
		isRoot: currentIsRoot,
		isZero: currentIsZero,
		lockRedeem: currentLockRedeem,
		lockStake: currentLockStake,
		lockStakeProfit: currentLockStakeProfit,
		subsLockWithdraw: currentSubsLockWithdraw,
		unCollection: currentUnCollection,
		withoutFewRegionRanking: currentWithoutFewRegionRanking,
		withoutStakeRanking: currentWithoutStakeRanking
	};
	dialogEditVisible.value = true;
	// 清除上一次校验
	if (editFormRef.value) {
		editFormRef.value.clearValidate();
	}
};

// 提交编辑
const handleEditSubmit = async () => {
	if (!editFormRef.value) return;
	await editFormRef.value.validate(async (valid) => {
		if (!valid) return;
		editLoading.value = true;
		try {
			// 构建当前数据（使用 true/false）
			const currentData = {
				uid: editFormData.value.uid,
				enable: Boolean(editFormData.value.enable),
				lockWithdraw: Boolean(editFormData.value.lockWithdraw),
				isRoot: Boolean(editFormData.value.isRoot),
				isZero: Boolean(editFormData.value.isZero),
				lockRedeem: Boolean(editFormData.value.lockRedeem),
				lockStake: Boolean(editFormData.value.lockStake),
				lockStakeProfit: Boolean(editFormData.value.lockStakeProfit),
				subsLockWithdraw: Boolean(editFormData.value.subsLockWithdraw),
				unCollection: Boolean(editFormData.value.unCollection),
				withoutFewRegionRanking: Boolean(editFormData.value.withoutFewRegionRanking),
				withoutStakeRanking: Boolean(editFormData.value.withoutStakeRanking)
			};
			
			// 比较数据，只传递有变化的字段
			const params = {
				uid: currentData.uid // uid 必须传递
			};
			
			// 比较每个字段
			if (currentData.enable !== originalData.value.enable) {
				params.enable = currentData.enable;
			}
			if (currentData.lockWithdraw !== originalData.value.lockWithdraw) {
				params.lockWithdraw = currentData.lockWithdraw;
			}
			if (currentData.isRoot !== originalData.value.isRoot) {
				params.isRoot = currentData.isRoot;
			}
			if (currentData.isZero !== originalData.value.isZero) {
				params.isZero = currentData.isZero;
			}
			if (currentData.lockRedeem !== originalData.value.lockRedeem) {
				params.lockRedeem = currentData.lockRedeem;
			}
			if (currentData.lockStake !== originalData.value.lockStake) {
				params.lockStake = currentData.lockStake;
			}
			if (currentData.lockStakeProfit !== originalData.value.lockStakeProfit) {
				params.lockStakeProfit = currentData.lockStakeProfit;
			}
			if (currentData.subsLockWithdraw !== originalData.value.subsLockWithdraw) {
				params.subsLockWithdraw = currentData.subsLockWithdraw;
			}
			if (currentData.unCollection !== originalData.value.unCollection) {
				params.unCollection = currentData.unCollection;
			}
			if (currentData.withoutFewRegionRanking !== originalData.value.withoutFewRegionRanking) {
				params.withoutFewRegionRanking = currentData.withoutFewRegionRanking;
			}
			if (currentData.withoutStakeRanking !== originalData.value.withoutStakeRanking) {
				params.withoutStakeRanking = currentData.withoutStakeRanking;
			}
			
			// 如果没有数据变化，提示用户
			if (Object.keys(params).length === 1) {
				ElNotification.warning("没有数据变化");
				editLoading.value = false;
				return;
			}
			
			const res = await userSetUser(params);
			if (res.code === 200) {
				ElNotification.success("编辑成功");
				dialogEditVisible.value = false;
				getList();
			} else {
				ElNotification.error(res.msg || "编辑失败");
			}
		} catch (error) {
			ElNotification.error("编辑失败");
		} finally {
			editLoading.value = false;
		}
	});
};

// 打开设置上级对话框
const openSetParentDialog = (row) => {
	setParentFormData.value = {
		// userId 使用列表的 ID 字段
		userId: String(row.id),
		// 用列表里的上级邀请码 refCode 作为回显
		parentCode: row.refCode ? String(row.refCode) : ""
	};
	dialogSetParentVisible.value = true;
	// 清除上一次校验
	if (setParentFormRef.value) {
		setParentFormRef.value.clearValidate();
	}
};

// 提交设置上级
const handleSetParentSubmit = async () => {
	if (!setParentFormRef.value) return;
	await setParentFormRef.value.validate(async (valid) => {
		if (!valid) return;
		setParentLoading.value = true;
		try {
			const params = {
				userId: setParentFormData.value.userId
			};
			// 如果输入了上级邀请码，则添加到参数中
			if (setParentFormData.value.parentCode && setParentFormData.value.parentCode.trim()) {
				params.parentCode = setParentFormData.value.parentCode.trim();
			}
			const res = await userSetParent(params);
			if (res.code === 200) {
				ElNotification.success("设置上级成功");
				dialogSetParentVisible.value = false;
				getList();
			} else {
				ElNotification.error(res.msg || "设置上级失败");
			}
		} catch (error) {
			ElNotification.error("设置上级失败");
		} finally {
			setParentLoading.value = false;
		}
	});
};

// 重置上级
const handleResetParent = async (row) => {
	// userId 使用列表的 ID 字段
	const userId = String(row.id);
	if (!userId) {
		ElNotification.error("用户ID为空，无法重置上级");
		return;
	}
	try {
		// 使用 Element Plus 确认弹窗
		await ElMessageBox.confirm(
			`确认重置用户 ${userId} 的上级吗？`,
			"提示",
			{
				type: "warning",
				confirmButtonText: "确定",
				cancelButtonText: "取消"
			}
		);

		const res = await userSetParent({ userId });
		if (res.code === 200) {
			ElNotification.success("重置上级成功");
			getList();
		} else {
			ElNotification.error(res.msg || "重置上级失败");
		}
	} catch (error) {
		// 用户取消确认不提示错误
		if (error !== "cancel" && error !== "close") {
			ElNotification.error("重置上级失败");
		}
	}
};
</script>

<style scoped>
.operation {
	display: inline-block;
	vertical-align: middle;
}
</style>

