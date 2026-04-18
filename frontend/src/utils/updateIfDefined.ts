/**
 * SWR の mutate コールバックでよく使う「前のキャッシュが無ければそのまま、あれば更新」パターンを短く書く。
 * `(prev) => (prev ? { ...prev, field: x } : prev)` の三項演算子を毎回書かなくてよくなる。
 * 戻り値に null を含めない（SWR のキャッシュ型は通常 `T | undefined` のため）。
 */
export function updateIfDefined<T>(
  prev: T | undefined,
  updater: (current: T) => T,
): T | undefined {
  if (prev === undefined) return prev;
  return updater(prev);
}
