<template>
  <div class="h-screen flex flex-col overflow-hidden">
    <!-- 固定头部 -->
    <div class="sticky-header">
      <div class="relative bg-white px-4 pt-3 pb-1 border-b border-gray-200">
        <div class="flex items-center justify-between px-4">
          <div>
            <van-icon name="arrow-left" @click="goBack" />
          </div>
          <div class="flex w-full justify-center">
            Message Center
          </div>
        </div>
      </div>
      <div class="scroll-container">
        <div class="message_list" v-if="list && list.length > 0">
          <div v-for="data in list" :key="Number(data.id)" class="message_item">
            <div class="message_title">
              <p class="title">{{ data.event }}</p>
              <p>{{ formatTimestamp(Number(data.time)) }}</p>
            </div>
            <div class="content">{{ `User ${data.watchUserId} ${data.event} ${data.size} ${data.token}` }}</div>
          </div>
        </div>
        <Empty v-else />
      </div>
    </div>


  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue"; 
import { formatTimestamp } from "@/utils/time-utils";
import Empty from "@/components/empty.vue";
import { UserWatchMsg, UserWatchService } from "@/sdk/backend-sdk/userwatch/v1/userwatch_pb";
import createClient from "@/sdk/transport";
import { useAppKitAccount } from "@reown/appkit/vue";
import { useRouter } from "vue-router";
const router = useRouter()

function goBack() {
  router.back()
}

onMounted(() => {
  // 获取消息列表
  getMessageList();
});

async function getMessageList() {
  const userWatchClient = createClient(UserWatchService)
  const account = useAppKitAccount();

  const resp = await userWatchClient.listUserWatchMsg({ userId: account.value?.address });
  list.value = resp.results
}

const list = ref<UserWatchMsg[]>([]);
</script>

<style scoped>
.sticky-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: white;
}

.scroll-container {
  height: calc(100vh - 100px);
  /* 根据实际头部高度调整 */
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

/* Custom styles for Vant components */
.message_list {
  padding: 23px 28px 0 28px;

  .message_item {
    margin-bottom: 10px;

    .message_title {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
      font-size: 12px;

      .title {
        font-weight: 600;
      }
    }

    .content {
      color: #979797;
      font-size: 12px;
    }
  }
}
</style>
