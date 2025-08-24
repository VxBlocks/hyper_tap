<template>
  <div class="token_box">
    <div class="token_left">
      <div class="token_info">
        <img :src="queryIconUrl()" alt="star" onerror="this.src='/crypto/svg/color/BTC.svg'" class="token_icon" />
        <div>
          <div class="flex items-center">
            <span class="token_name">{{ data.symbol }}-{{ data.pair }}</span>
            <div class="token_tag">
              {{ data.leverage + "x" }}
            </div>
          </div>
          <div class="token_amount">{{ $t("labels.volume") }} $ {{ formatNumberWithCommas(data.volume) }}</div>
        </div>
      </div>
    </div>
    <div class="token_right">
      <div class="token_price">${{ formatNumberWithCommas(data.price) }}</div>
      <div :class="['token_rate', data.change.startsWith('-') ? 'loss' : 'profit']">
        {{ format_change(data.change) }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useHomeStore } from '@/store/modules/home';
import { formatNumberWithCommas } from '@/utils/common';
import { computed, ref } from 'vue';
const homeStore = useHomeStore()
const coinToAsset = computed(() => homeStore.getTokens.coin_to_asset)
const perps_coin_to_price = computed(() => homeStore.getTokens.perps_coin_to_price)

const props = defineProps({
  data: {
    type: Object,
    required: true,
    default: () => ({
      isCollected: false,
      symbol: "",
      price: "",
      pair: "",
      change: "",
      leverage: "",
    }),
  },
});
// 预先导入所有图标
function queryIconUrl() {
  return `/crypto/svg/color/${props.data.symbol}.svg`
}
function format_change(value: string): string {
  return value.startsWith("-") ? `${value}%` : `+${value}%`;
}
</script>

<style scoped>
.token_box {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 13px 0;

  .token_left {
    display: flex;
    gap: 5px;
    align-items: center;

    .collect_icon {
      width: 22px;
      height: 22px;
    }

    .token_info {
      display: flex;
      align-items: center;
      gap: 10px;

      .token_name {
        font-size: 14px;
      }

      .token_tag {
        margin-left: 13px;
        padding: 0px 8px;
        background-color: #c5ede7;
        color: #46ccb9;
        font-size: 10px;
        height: fit-content;
        border-radius: 4px;
      }

      .token_amount {
        color: #78817e;
        font-size: 11px;
      }
    }
  }

  .token_right {
    .token_price {
      font-size: 14px;
    }

    .token_rate {
      font-size: 11px;
      text-align: end;
    }

    .loss {
      color: #3a957f;
    }

    .profit {
      color: #f7435d;
    }
  }
}

.token_icon {
  width: 36px;
  height: 36px;
}
</style>
