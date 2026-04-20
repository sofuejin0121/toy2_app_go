/** ブラウザタブ・<title> のサイト名（index.html と揃える） */
export const SITE_NAME = 'Chirp';

/** 「ページ名 - Chirp」形式 */
export function formatDocumentTitle(pageLabel: string): string {
  return `${pageLabel} - ${SITE_NAME}`;
}

/** 投稿スレッド用: 本文先頭をタイトル向けに短くする */
export function micropostTabTitle(content: string, maxLen = 48): string {
  const singleLine = content.replace(/\s+/g, ' ').trim();
  if (!singleLine) return '投稿';
  if (singleLine.length <= maxLen) return singleLine;
  return `${singleLine.slice(0, maxLen)}…`;
}
