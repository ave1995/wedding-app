import { useRef, useState } from "react";
import { smallButtonBurst } from "../functions/success";

interface ButtonProps {
  label: string;
  type: ButtonType;
  onClickAsync: () => Promise<void>;
}

export default function Button({
  label,
  type,
  onClickAsync: onClickAsync,
}: ButtonProps) {
  const [loading, setLoading] = useState(false);
  const buttonRef = useRef<HTMLButtonElement>(null);

  const handleClick = async () => {
    if (loading) return;

    setLoading(true);
    try {
      await onClickAsync();
    } finally {
      setLoading(false);
      if (buttonRef.current) {
        smallButtonBurst(buttonRef.current);
      }
    }
  };

  return (
    <button
      ref={buttonRef}
      className={`border px-3 py-2 rounded-lg text-white 
        ${GetButtonColor(type)}
        ${GetButtonHoverColor(type)}
        hover:cursor-pointer`}
      onClick={handleClick}
    >
      <div className="relative flex items-center justify-center">
        <p
          className={`font-bold text-sm/6 ${
            loading ? "opacity-0" : "opacity-100"
          }`}
        >
          {label}
        </p>
        <svg
          className={`absolute size-5 text-white ${
            loading ? "opacity-100 animate-spin" : "opacity-0"
          }`}
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
        >
          <circle
            className="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            strokeWidth="4"
          ></circle>
          <path
            className="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          ></path>
        </svg>
      </div>
    </button>
  );
}

export const ButtonTypeEnum = {
  Basic: "Basic",
} as const;

type ButtonType = (typeof ButtonTypeEnum)[keyof typeof ButtonTypeEnum];

function GetButtonColor(type: ButtonType) {
  switch (type) {
    case ButtonTypeEnum.Basic:
      return "bg-pink-500";
    default:
      return "bg-pink-500";
  }
}

function GetButtonHoverColor(type: ButtonType) {
  switch (type) {
    case ButtonTypeEnum.Basic:
      return "hover:bg-pink-700";
    default:
      return "hover:bg-pink-700";
  }
}
