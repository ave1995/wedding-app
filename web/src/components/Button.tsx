import { useState } from "react";
// import { smallButtonBurst } from "../functions/success";
import Spinner from "./Spinner";

interface ButtonProps<T = void> {
  label: string;
  type?: ButtonType;
  onClick: () => Promise<T> | void;
}

export default function Button<T = void>({
  label,
  type,
  onClick,
}: ButtonProps<T>) {
  const [loading, setLoading] = useState(false);
  // const buttonRef = useRef<HTMLButtonElement>(null);

  const handleClick = async () => {
    if (loading) return;

    setLoading(true);
    try {
      await Promise.resolve(onClick());
    } finally {
      setLoading(false);
      // if (buttonRef.current) {
      //   smallButtonBurst(buttonRef.current);
      // }
    }
  };

  return (
    <button
      // ref={buttonRef}
      className={`w-full px-3 py-2 text-center border-b-4 rounded-2xl 
        ${GetButtonColor(type)}
        ${GetButtonHoverColor(type)}
        ${
          loading
            ? ""
            : "hover:cursor-pointer active:scale-95 active:shadow-sm transition-all duration-150 ease-in-out"
        }
        focus:outline-none font-semibold`}
      onClick={handleClick}
      disabled={loading}
    >
      <div className="relative flex items-center justify-center">
        <p
          className={`text-lg font-semibold ${
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
  Blue: "Blue",
} as const;

type ButtonType = (typeof ButtonTypeEnum)[keyof typeof ButtonTypeEnum];

function GetButtonColor(type?: ButtonType) {
  switch (type) {
    case ButtonTypeEnum.Basic:
      return "text-white bg-pink-500 hover:bg-pink-600 hover:border-pink-400 border-pink-300";
    case ButtonTypeEnum.Blue:
      return "border-2 text-gray-800 border-gray-300 ";
    default:
      return "text-white bg-pink-500 hover:bg-pink-600 hover:border-pink-400 border-pink-300";
  }
}

function GetButtonHoverColor(type?: ButtonType) {
  switch (type) {
    case ButtonTypeEnum.Basic:
      return "hover:bg-pink-700";
    case ButtonTypeEnum.Blue:
      return "hover:bg-gray-100 hover:border-gray-400";
    default:
      return "hover:bg-pink-700";
  }
}
