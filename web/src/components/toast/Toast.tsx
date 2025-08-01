import { useEffect } from "react";
import './Toast.css';

interface ToastProps {
  message: string;
  onClose: () => void;
  duration?: number;
}

export default function Toast({
  message,
  onClose,
  duration = 3000,
}: ToastProps) {
  useEffect(() => {
    const timer = setTimeout(onClose, duration);
    return () => clearTimeout(timer);
  }, [onClose, duration]);

  return (
    <div className="fixed top-5 border rounded-lg px-3 py-2 flex gap-2 bg-green-600 text-white font-bold text-sm/6 animate-slideDown">
      <p>{message}</p>
      <button
        onClick={onClose}
        className="text-xl font-bold leading-none hover:text-gray-200"
      >
        &times;
      </button>
    </div>
  );
}
