<!-- WalletsList -->
<script lang="ts" setup>
import { computed, ref } from "vue";

const activeNames = ref([]);

const handleSave = () => {};
const handleDelete = () => {};
function queryCoinName() {}
</script>
<template>
  <div>
    <div class="tag_title">
      <img src="/tag.png" alt="tag" />
      <span>增仓</span>
    </div>
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
              <p class="label">Name</p>
              <span class="value">HYPE</span>
            </div>
            <div class="detail-col">
              <p class="label">Type</p>
              <span class="value">HYPE</span>
            </div>
            <div class="detail-col">
              <p class="label">{{ $t("labels.quantity") }}</p>
              <span class="loss"> 0 ETH </span>
            </div>
            <div class="detail-col">
              <p class="label">价值</p>
              <span class="loss"> 0 ETH </span>
            </div>
            <div class="detail-col-button ">
              <div class="add_button">跟单</div>
            </div>
          </div>
        </template>
        <div></div>
      </van-collapse-item>
    </van-collapse>
  </div>
    <div>
    <div class="tag_title">
      <img src="/tag.png" alt="tag" />
      <span>仓位</span>
    </div>
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
              <p class="label">Name</p>
              <span class="value">HYPE</span>
            </div>
            <div class="detail-col">
              <p class="label">Type</p>
              <span class="value">HYPE</span>
            </div>
            <div class="detail-col">
              <p class="label">{{ $t("labels.quantity") }}</p>
              <span class="loss"> 0 ETH </span>
            </div>
            <div class="detail-col">
              <p class="label">价值</p>
              <span class="loss"> 0 ETH </span>
            </div>
            <div class="detail-col-button ">
              <div class="add_button">跟单</div>
            </div>
          </div>
        </template>
        <div></div>
      </van-collapse-item>
    </van-collapse>
  </div>
</template>

<style lang="scss">
.collapse_list {
  background-color: #f8f9fa;
  margin-bottom: 3px;
}
.tag_title {
  display: flex;
  gap: 5px;
  align-items: center;
  margin-top: 10px;
  margin-left: 5px;
}
.asset-title {
  display: flex;
  width: 100%;

  .detail-col {
    width: 15%;

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
  .detail-col-button {
    display: flex;
    align-items: center;
    width: 100px;
    justify-content: center;

    .add_button {
      border-radius: 10px;
      background-color: #3a957f;
      color:#FFFFFF;
      padding: 0px 10px;
      font-size:12px
    }
  }
}

.sent_button {
  color: #46ccb9;
  cursor: pointer;
}
</style>
