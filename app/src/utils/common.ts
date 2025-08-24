import { BigNumberUtils } from "./big-number-utils";

export function formatNumber(num: number | string) {
    if (!num) {
        return "0";
    }

    const numberValue = Number(num);

    // 检查是否为有效数字
    if (isNaN(numberValue)) {
        return "0";
    }

    // 转换为字符串以计算小数位数
    const numStr = numberValue.toString();
    const decimalIndex = numStr.indexOf('.');

    // 如果没有小数点，是整数
    if (decimalIndex === -1) {
        return numberValue.toLocaleString('en-US', {
            minimumFractionDigits: 0,
            maximumFractionDigits: 0,
        });
    }

    // 计算小数位数
    const decimalPlaces = numStr.length - decimalIndex - 1;

    return numberValue.toLocaleString('en-US', {
        minimumFractionDigits: decimalPlaces,
        maximumFractionDigits: decimalPlaces,
    });
}


export function calRoe(latestPrice: string, entryNtlPrice: string) {
    // 检查输入参数有效性
    const latest = parseFloat(latestPrice);
    const entry = parseFloat(entryNtlPrice);

    // 如果任一参数无效或入场价格为0，返回"0.00"
    if (isNaN(latest) || isNaN(entry) || entry === 0) {
        return "--";
    }
    const roe = BigNumberUtils.divide((BigNumberUtils.subtract(latestPrice, entryNtlPrice)), entry);

    return BigNumberUtils.round(BigNumberUtils.multiply(roe, 100));
}

export function formatNumberWithCommas(input: number | string): string {
  // 1. 清理输入：移除所有非数字、小数点、负号的字符（如已有逗号）
  const cleaned = input.toString().replace(/[^\d.-]/g, '');

  // 2. 检查是否为有效数字
  const num = parseFloat(cleaned);
  if (isNaN(num)) return '0'; // 无效输入时返回 0 或抛出错误

  // 3. 分离整数和小数部分
  const [integerPart, decimalPart] = num.toString().split('.');

  // 4. 格式化整数部分（添加千分位逗号）
  const formattedInteger = integerPart.replace(/\B(?=(\d{3})+(?!\d))/g, ',');

  // 5. 组合结果
  return decimalPart ? `${formattedInteger}.${decimalPart}` : formattedInteger;
}