<!-- WalletsList -->
<script lang="ts" setup>
import { computed, ref } from "vue";

const activeNames = ref([]);

const handleSave = () => {};
const handleDelete = () => {};
function queryCoinName() {}
</script>
<template>
  <van-collapse v-model="activeNames" accordion>
    <van-collapse-item
      v-for="(item, index) in [1, 2]"
      :key="index"
      :name="index"
      class="collapse_list"
    >
      <template #title>
        <div class="asset-title">
          <div class="detail-col">
            <p class="label">钱包</p>
            <span class="value">HYPE</span>
          </div>
          <div class="detail-col">
            <p class="label">币种</p>
            <span class="value">HYPE</span>
          </div>
          <div class="detail-col">
            <p class="label">资产</p>
            <span class="loss"> 0 ETH </span>
          </div>
          <div class="filters-item">
            <p class="label">Filters</p>
            <div class="content">
              <div
                v-for="(item, index) in [1, 2]"
                :key="index"
                :name="index"
                class="content-item"
              >
                <img src="/tag.png" alt="tag" />
                <span>增仓</span>
              </div>
            </div>
          </div>
        </div>
      </template>
      <div>
        <div class="asset-title">
          <div class="filters-item">
            <p class="label">Filters</p>
            <div class="content">
              <div
                v-for="(item, index) in [1, 2]"
                :key="index"
                :name="index"
                class="content-item"
              >
                <img src="/tag.png" alt="tag" />
                <span>增仓</span>
              </div>
            </div>
          </div>
        </div>
        <van-divider />
        <div class="flex gap-10">
          <div class="sent_button" @click="handleSave">保存</div>
          <div class="sent_button" @click="handleDelete">删除</div>
        </div>
      </div>
    </van-collapse-item>
  </van-collapse>
</template>

<style lang="scss">
.collapse_list {
  background-color: #f8f9fa;
  margin-bottom: 3px;
}

.asset-title {
  display: flex;
  width: 100%;

  .detail-col {
    width: 75px;

    .label {
      color: #d9d9d9;
    }

    .value {
      color: #000000;
    }

    .loss {
      color: #3a957f;
    }
  }
  .filters-item {
    flex-grow: 1;
    .label {
      color: #d9d9d9;
    }
    .content {
      display: flex;
      gap: 8px;
      color: #000000;
      
      .content-item {
        display: flex;
        gap: 3px;
        align-items: center;
      }
    }
  }
}

.sent_button {
  color: #46ccb9;
  cursor: pointer;
}
</style>
