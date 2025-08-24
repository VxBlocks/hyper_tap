<template>
  <div class="custom-progress-container">
    <!-- 进度条容器 -->
    <div class="progress-track">
      <!-- 进度线 -->
      <div 
        class="progress-line" 
        :style="{ width: percentage + '%' }"
      ></div>
      
      <!-- 刻度点 -->
      <div class="progress-ticks">
        <div class="tick" style="left: 0%"></div>
        <div class="tick" style="left: 25%"></div>
        <div class="tick" style="left: 50%"></div>
        <div class="tick" style="left: 75%"></div>
        <div class="tick" style="left: 100%"></div>
      </div>
      
      <!-- 当前位置指示器 -->
      <div 
        class="current-indicator" 
        :style="{ left: percentage + '%' }"
      ></div>
    </div>
    
    <!-- 百分比标签 -->
    <div class="percentage-label">
      {{ percentage }}%
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

// 进度百分比
const percentage = ref(100);

// 计算当前进度状态
const progressStatus = computed(() => {
  if (percentage.value === 0) return 'start';
  if (percentage.value === 100) return 'complete';
  return 'in-progress';
});
</script>

<style scoped>
.custom-progress-container {
  width: 100%;
  padding: 20px 0;
  position: relative;
}

.progress-track {
  position: relative;
  width: 100%;
  height: 8px;
  background-color: #f5f5f5;
  border-radius: 4px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.progress-line {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background-color: #3A957F;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.progress-ticks {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.tick {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 12px;
  height: 12px;
  background-color: #3A957F;
  border-radius: 50%;
  border: 2px solid #fff;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

/* 最后一个刻度点的特殊样式 */
.tick:last-child {
  border: 2px solid #3A957F;
  background-color: transparent;
  border-radius: 50%;
  width: 16px;
  height: 16px;
  transform: translateX(-50%) translateY(-50%);
}

.current-indicator {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 16px;
  height: 16px;
  background-color: #3A957F;
  border: 2px solid #fff;
  border-radius: 50%;
  z-index: 1;
  transition: left 0.3s ease;
}

/* 当进度为100%时的特殊样式 */
.current-indicator {
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0% {
    transform: translateY(-50%) scale(1);
  }
  50% {
    transform: translateY(-50%) scale(1.2);
  }
  100% {
    transform: translateY(-50%) scale(1);
  }
}

.percentage-label {
  position: absolute;
  right: 0;
  bottom: -20px;
  font-size: 16px;
  color: #3A957F;
  font-weight: 500;
}
</style>