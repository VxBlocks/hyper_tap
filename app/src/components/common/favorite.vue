<template>
    <button class="favorite-button" @click="toggleFavorite">
        <img v-if="isFavorite" src="~@/assets/icons/star_1.png" :alt="isFavorite ? '已收藏' : '未收藏'"
            class="favorite-icon" />
        <img v-else src="~@/assets/icons/star_0.png" :alt="isFavorite ? '已收藏' : '未收藏'"
            class="favorite-icon" />
    </button>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';

const isFavorite = ref(false); // 收藏状态

// 定义 props
const props = defineProps<{
  itemId: string;
  isFavorite: boolean;
}>();

onMounted(() => {
  // 在这里处理收藏状态的初始化逻辑
  isFavorite.value = props.isFavorite;
});

// 定义 emits
const emit = defineEmits<{
    (e: 'favorite-change', itemId: string, isFavorite: boolean): void  // 添加更新 visible 的事件
}>();;

function toggleFavorite() {
    isFavorite.value = !isFavorite.value;
    emit('favorite-change', props.itemId, isFavorite.value)
}
</script>

<style scoped>
.favorite-button {
    background: none;
    border: none;
    cursor: pointer;
    padding: 8px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    transition: transform 0.2s ease;
}

.favorite-button:hover {
    transform: scale(1.1);
}

.favorite-button:active {
    transform: scale(0.95);
}

.favorite-icon {
    width: 24px;
    height: 24px;
    object-fit: contain;
}
</style>