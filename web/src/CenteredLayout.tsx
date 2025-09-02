import { type ReactNode } from "react";

interface CenteredLayoutProps {
  children: ReactNode;
}
const CenteredLayout = ({ children }: CenteredLayoutProps) => {
  return <div className="flex w-screen h-screen justify-center items-center">{children}</div>;
};

export default CenteredLayout;
