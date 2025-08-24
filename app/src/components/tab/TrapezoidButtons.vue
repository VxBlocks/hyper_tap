<template>
    <div class="trapezoid-buttons-container">
        <div class="trapezoid-button left-button" :class="{ 'selected': leftSelected }" @click="selectLeft">
            <span class="button-text">{{ leftText }}</span>
        </div>
        <div class="trapezoid-button right-button" :class="{ 'selected': rightSelected }" @click="selectRight">
            <span class="button-text">{{ rightText }}</span>
        </div>
    </div>
</template>

<script setup lang="ts">
interface Props {
    leftText?: string;
    rightText?: string;
    initialSelected?: 'left' | 'right';
    clipSize?: number;
    minWidth?: number;
}

const props = withDefaults(defineProps < Props > (), {
    leftText: '选项1',
    rightText: '选项2',
    initialSelected: 'left',
    clipSize: 12,
    minWidth: 68
});

// 验证 initialSelected 值
if (!['left', 'right'].includes(props.initialSelected)) {
    throw new Error('initialSelected must be either "left" or "right"');
}

// 定义 emits
const emit = defineEmits < {
    (e: 'change', value: 'left' | 'right'): void;
}> ();

// 响应式数据
import { ref, watch } from 'vue';

const leftSelected = ref(props.initialSelected === 'left');
const rightSelected = ref(props.initialSelected === 'right');

// 监听 initialSelected 变化
watch(() => props.initialSelected, (newVal) => {
    leftSelected.value = newVal === 'left';
    rightSelected.value = newVal === 'right';
});

// 方法
const selectLeft = () => {
    leftSelected.value = true;
    rightSelected.value = false;
    emit('change', 'left');
};

const selectRight = () => {
    leftSelected.value = false;
    rightSelected.value = true;
    emit('change', 'right');
};
</script>

<style scoped>
.trapezoid-buttons-container {
    display: flex;
    position: relative;
    height: 22px;
    --clip-size: 10px;
    --text-clip: 16px;
    --gap-size: 8px;
    --selected-color: #3A957F;
    --unselected-color: #F0F0F0;
}

.trapezoid-button {
    position: relative;
    height: 100%;
    display: flex;
    align-items: center;
    min-width: var(--min-width, 68px);
    cursor: pointer;
    transition: all 0.3s ease;
    background-color: var(--unselected-color);
}

.trapezoid-button.selected {
    background-color: var(--selected-color);
    color: white;
}

/* 左侧按钮 - 剪掉右下角（上长下短） */
.left-button {
    font-size: 12px;
    clip-path: polygon(0% 0%, 100% 0%, calc(100% - var(--clip-size)) 100%, 0% 100%);
    padding-right: var(--text-clip);
    margin-right: calc(-1 * var(--clip-size) + var(--gap-size));
}

/* 右侧按钮 - 剪掉左上角（上短下长） */
.right-button {
    font-size: 12px;
    clip-path: polygon(var(--clip-size) 0%, 100% 0%, 100% 100%, 0% 100%);
    padding-left: var(--text-clip);
    margin-left: calc(-1 * var(--clip-size) + var(--gap-size));
}

.button-text {
    position: relative;
    z-index: 1;
    text-align: center;
    width: 100%;
}

/* 左侧按钮文字 - 基于底边（较短边）居中 */
.left-button .button-text {
    transform: translateX(calc(var(--clip-size) * 0.25));
    padding-left: calc(var(--clip-size) * 0.5);
}

/* 右侧按钮文字 - 基于上边（较短边）居中 */
.right-button .button-text {
    transform: translateX(calc(-1 * var(--clip-size) * 0.25));
    padding-right: calc(var(--clip-size) * 0.5);
}
</style>