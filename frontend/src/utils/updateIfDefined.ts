/**
 * SWR の mutate コールバックでよく使う「前のキャッシュが無ければそのまま、あれば更新」パターンを短く書く。
 *
 * 初心者が迷いやすい `(prev) => (prev ? { ...prev, x: y } : prev)` の「外側の括弧が何のためか」を避け、
 * `mutate((prev) => updateIfDefined(prev, (cur) => ({ ...cur, x: y })))` に統一する。
 *
 * 戻り値に null を含めない（SWR のキャッシュ型は通常 `T | undefined` のため）。
 */
export function updateIfDefined<T>(
  prev: T | undefined,
  updater: (current: T) => T,
): T | undefined {
  if (prev === undefined) return prev;
  return updater(prev);
}
