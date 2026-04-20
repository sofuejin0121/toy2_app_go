import { useEffect } from 'react';
import { formatDocumentTitle, SITE_NAME } from '../lib/documentTitle';

/**
 * ルートに応じた document.title を設定する。
 * pageLabel を省略・null・空ならサイト名のみ（「Chirp」）。
 */
export function useDocumentTitle(pageLabel?: string | null) {
  useEffect(() => {
    document.title =
      pageLabel && pageLabel.length > 0 ? formatDocumentTitle(pageLabel) : SITE_NAME;
  }, [pageLabel]);
}
