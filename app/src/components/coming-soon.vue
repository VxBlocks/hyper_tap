<template>
  <van-empty
    class="custom-image flex justify-center items-center"
    image="https://img.yzcdn.cn/vant/custom-empty-image.png"
    :description="description"
  />
</template>

<script lang="ts" setup>
  defineProps({
    description: {
      type: String,
      default: 'Coming Soon',
    },
  });
</script>

<style lang="scss">
  .custom-image .van-empty__image {
    width: 30px;
    height: 30px;
  }
</style>
