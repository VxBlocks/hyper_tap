// utils/big-number-utils.ts
import Decimal from 'decimal.js';

export class BigNumberUtils {
  /**
   * 加法运算
   */
  static add(a: string | number, b: string | number): string {
    return new Decimal(a).add(new Decimal(b)).toString();
  }

  /**
   * 减法运算
   */
  static subtract(a: string | number, b: string | number): string {
    return new Decimal(a).sub(new Decimal(b)).toString();
  }

  /**
   * 乘法运算
   */
  static multiply(a: string | number, b: string | number): string {
    return new Decimal(a).mul(new Decimal(b)).toString();
  }

  /**
   * 除法运算
   */
  static divide(a: string | number, b: string | number): string {
    return new Decimal(a).div(new Decimal(b)).toString();
  }

  /**
   * 比较大小
   * 返回值: 1 (a > b), -1 (a < b), 0 (a == b)
   */
  static compare(a: string | number, b: string | number): number {
    return new Decimal(a).cmp(new Decimal(b));
  }

  /**
   * 格式化数字显示
   */
  static format(value: string | number, decimals: number = 2): string {
    return new Decimal(value).toFixed(decimals);
  }

  /**
   * 四舍五入
   */
  static round(value: string | number, decimals: number = 2): string {
    return new Decimal(value).toDecimalPlaces(decimals).toString();
  }

  /**
   * 向下取整
   */
  static floor(value: string | number, decimals: number = 0): string {
    if (Math.abs(Number(value)) < 0) {
      return new Decimal(value).toDecimalPlaces(6, Decimal.ROUND_DOWN).toString();
    }
    return new Decimal(value).toDecimalPlaces(decimals, Decimal.ROUND_DOWN).toString();
  }

  /**
   * 向上取整
   */
  static ceil(value: string | number, decimals: number = 0): string {
    if (Math.abs(Number(value)) < 0) {
      return new Decimal(value).toDecimalPlaces(6, Decimal.ROUND_DOWN).toString();
    }
    return new Decimal(value).toDecimalPlaces(decimals, Decimal.ROUND_UP).toString();
  }
  /**
   * 保留指定有效数字和最多小数位的四舍五入
   * @param value        原始数值
   * @param sigFigs      有效数字个数
   * @param maxDecimals  最多小数位数
   * @returns            处理后的数值
   */
  static roundToSignificantAndDecimal(value: number | string, sigFigs: number, maxDecimals: number): string {
    const absValue = new Decimal(value).abs();
    if (absValue.isZero()) return "0";

    const magnitude = absValue.logarithm(10).floor().toNumber();
    const scale = new Decimal(10).pow(sigFigs - magnitude - 1);
    const rounded = absValue.mul(scale).toNearest(1).div(scale);

    // 恢复原始符号
    const signedResult = Number(value) < 0 ? rounded.neg() : rounded;

    // 保留最多maxDecimals位小数
    return new Decimal(signedResult).toDecimalPlaces(maxDecimals).toString();
  }
}