interface InputTextProps {
  placeholder?: string;
}

export default function InputText({ placeholder }: InputTextProps) {
  return (
    <input
      type="text"
      placeholder={placeholder}
      className="border rounded-lg px-3 py-2 text-black"
    />
  );
}