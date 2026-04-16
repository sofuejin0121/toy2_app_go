import type { Pagination as PaginationType } from '../types';

interface Props {
  pagination: PaginationType;
  onPageChange: (page: number) => void;
}

export default function Pagination({ pagination, onPageChange }: Props) {
  const { current_page, total_pages, has_prev, has_next } = pagination;
  if (total_pages <= 1) return null;

  const pages: (number | '...')[] = [];
  if (total_pages <= 7) {
    for (let i = 1; i <= total_pages; i++) pages.push(i);
  } else {
    pages.push(1);
    if (current_page > 3) pages.push('...');
    for (
      let i = Math.max(2, current_page - 1);
      i <= Math.min(total_pages - 1, current_page + 1);
      i++
    ) {
      pages.push(i);
    }
    if (current_page < total_pages - 2) pages.push('...');
    pages.push(total_pages);
  }

  return (
    <div className="flex items-center justify-center gap-1 mt-6">
      <button
        onClick={() => onPageChange(current_page - 1)}
        disabled={!has_prev}
        className="px-3 py-1.5 text-sm rounded border border-gray-200 text-gray-600 hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
      >
        ‹
      </button>
      {pages.map((p, i) =>
        p === '...' ? (
          <span key={`dot-${i}`} className="px-2 text-gray-400">
            …
          </span>
        ) : (
          <button
            key={p}
            onClick={() => onPageChange(p as number)}
            className={`px-3 py-1.5 text-sm rounded border transition-colors ${
              p === current_page
                ? 'bg-blue-600 text-white border-blue-600'
                : 'border-gray-200 text-gray-600 hover:bg-gray-50'
            }`}
          >
            {p}
          </button>
        ),
      )}
      <button
        onClick={() => onPageChange(current_page + 1)}
        disabled={!has_next}
        className="px-3 py-1.5 text-sm rounded border border-gray-200 text-gray-600 hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
      >
        ›
      </button>
    </div>
  );
}
