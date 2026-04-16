// エラーメッセージ表示
// 使い方: if (error) return <ErrorMessage message={error} />;
interface Props {
  message: string;
}

export default function ErrorMessage({ message }: Props) {
  return (
    <div className="max-w-lg mx-auto mt-10 p-4 bg-red-50 border border-red-200 rounded-xl text-red-700 text-sm">
      {message}
    </div>
  );
}
