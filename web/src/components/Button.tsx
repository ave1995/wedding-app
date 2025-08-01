import { useRef, useState } from "react";
import { smallButtonBurst } from "../functions/success";
import Spinner from "./spinner";

interface ButtonProps<T> {
  label: string;
  type: ButtonType;
  onClickAsync: () => Promise<T>;
}

export default function Button<T>({
  label,
  type,
  onClickAsync,
}: ButtonProps<T>) {
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
        <Spinner loading={loading} className="absolute size-5"></Spinner>
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
