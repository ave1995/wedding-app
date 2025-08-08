interface InputTextProps {
  placeholder?: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

export default function InputText({ placeholder, value, onChange }: InputTextProps) {
  return (
    <input
      type="text"
      placeholder={placeholder}
      value={value}
      onChange={onChange}
      className="border-2 rounded-lg px-3 py-2 text-black"
    />
  );
}