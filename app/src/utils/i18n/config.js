import { createI18n } from "vue-i18n";
// Chinese language pack
import zh from "./zh-CN.json";
// English language pack
import en from "./en-US.json";
// 准备语言包
const messages = {
  en,
  zh,
};

// 创建 i18n 实例
const i18n = createI18n({
  legacy: false, // 使用 Composition API 模式
  locale: localStorage.getItem("lang") || "en", // 从本地存储获取或默认英文
  fallbackLocale: "en", // 备用语言
  messages,
});

export default i18n;
