<template>
    <div>
        <NavBar>
        </NavBar>
        <van-form ref="formRef">

        </van-form>
        <div class="note p-6">

        </div>
    </div>
</template>

<script setup lang="ts">
import NavBar from '@/components/navbar.vue';
// 定义 props 类型
interface Props {
    id: string // 因为路由参数永远是 string 类型
}

// 使用 withDefaults 或直接定义（Vue 3 <script setup> 自动解析 props）
const props = defineProps<Props>()

// 打印 id
console.log('User ID:', props.id)
</script>

<style scoped lang="less"></style>